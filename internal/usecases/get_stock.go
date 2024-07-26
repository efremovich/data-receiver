package usecases

import (
	"context"
	"database/sql"
	"errors"

	"github.com/efremovich/data-receiver/internal/entity"
	aerror "github.com/efremovich/data-receiver/pkg/aerror"
)

func (s *receiverCoreServiceImpl) ReceiveStocks(ctx context.Context) aerror.AError {
	client := s.apiFetcher["wb"]
	stockMetaList, err := client.GetStocks(ctx)
	if err != nil {
		return aerror.New(ctx, entity.GetDataFromExSources, err, "ошибка получение данные из внешнего источника %s в БД: %s ", "", err.Error())
	}

	for _, stock := range stockMetaList {
		wb2card, err := s.wb2cardrepo.SelectByNmid(ctx, stock.Wb2Card.NMID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении wb2card %s в БД.", "wb")
		}
		// TODO В случае отсутствия в Wb2Card - добавлять в него
		if wb2card == nil {
			continue
		}

		seller, err := s.sellerRepo.SelectByTitle(ctx, "wb")
		if err != nil {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении Seller %s в БД.", "wb")
		}

		card, err := s.cardRepo.SelectByID(ctx, wb2card.CardID)
		if err != nil {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении card %s в БД.", "wb")
		}

		size, err := s.sizerepo.SelectByTechSize(ctx, stock.Size.TechSize)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении size %s в БД.", "wb")
		}
		if size == nil {
      size, err = s.sizerepo.Insert(ctx, entity.Size{
				TechSize: stock.Size.TechSize,
				Title:    stock.Size.TechSize,
			})
      if err != nil{
        return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении size в БД.")
      }
		}

		priceSize, err := s.pricesizerepo.SelectByCardIDAndSizeID(ctx, card.ID, size.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении priceSize %s в БД.", "wb")
		}
		if priceSize == nil {
			priceSize, err = s.pricesizerepo.Insert(ctx, entity.PriceSize{
				Price:        stock.PriceSize.Price,
				Discount:     stock.PriceSize.Discount,
				SpecialPrice: stock.PriceSize.SpecialPrice,
				CardID:       card.ID,
				SizeID:       size.ID,
			})
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении priceSize в БД.")
			}
		} else {
			err = s.pricesizerepo.UpdateExecOne(ctx, *priceSize)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении priceSize в БД.")
			}
		}
		barcode, err := s.barcodeRepo.SelectByBarcode(ctx, stock.Barcode.Barcode)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении barcode %s в БД.", "wb")
		}
		if barcode == nil {
			barcode, err = s.barcodeRepo.Insert(ctx, entity.Barcode{
				Barcode:     stock.Barcode.Barcode,
				ExternalID:  0,
				PriceSizeID: priceSize.ID,
				SellerID:    seller.ID,
			})
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении barcode %s в БД.", stock.Barcode.Barcode)
			}
		}

		warehouse, err := s.warehouserepo.SelectBySellerIDAndTitle(ctx, seller.ID, stock.Warehouse.Title)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении warehouse %s в БД.", "wb")
		}
		if warehouse == nil {
			warehouse, err = s.warehouserepo.Insert(ctx, entity.Warehouse{
				Title:      stock.Warehouse.Title,
				ExternalID: 0,
				Address:    stock.Warehouse.Address,
				TypeID:     1,
				SellerID:   seller.ID,
			})
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении warehouse %s в БД.", stock.Warehouse.Title)
			}
		}

		stockData, err := s.stockrepo.SelectByBarcode(ctx, barcode.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении priceSize %s в БД.", "wb")
		}
		if stockData == nil {
			stockData, err = s.stockrepo.Insert(ctx, entity.Stock{
				Quantity:    stock.Stock.Quantity,
				BarcodeID:   barcode.ID,
				WarehouseID: warehouse.ID,
				CardID:      card.ID,
				SellerID:    seller.ID,
			})
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении stock в БД.")
			}
		} else {
			stockData.Quantity = stock.Stock.Quantity
			stockData.BarcodeID = barcode.ID
			stockData.WarehouseID = warehouse.ID
			stockData.CardID = card.ID
			stockData.SellerID = seller.ID

			err = s.stockrepo.UpdateExecOne(ctx, *stockData)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении stock в БД.")
			}
		}

	}
	return nil
}
