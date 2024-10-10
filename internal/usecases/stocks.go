package usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	aerror "github.com/efremovich/data-receiver/pkg/aerror"
	"github.com/efremovich/data-receiver/pkg/alogger"
)

func (s *receiverCoreServiceImpl) ReceiveStocks(ctx context.Context, desc entity.PackageDescription) error {
	client := s.apiFetcher[desc.Seller]

	stockMetaList, err := client.GetStocks(ctx, desc)
	if err != nil {
		return fmt.Errorf("ошибка получение данные из внешнего источника %s, %s", desc.Seller, err)
	}

	alogger.InfoFromCtx(ctx, "Получение данных об остатках за %s", desc.UpdatedAt.Format("02.01.2006"))

	var notFoundElements int

	for _, stock := range stockMetaList {
		seller, err := s.getSeller(ctx, desc.Seller)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		select2card, err := s.getSeller2Card(ctx, stock.Seller2Card.ExternalID, seller.ID)
		if err != nil {
			return err
		}
		// Получаем ошибку что такой записи нет, поицем карточку товара в 1с
		if errors.Is(err, ErrObjectNotFound) {
			query := make(map[string]string)
			query["barcode"] = stock.Barcode.Barcode
			query["article"] = stock.SupplierArticle

			desc := entity.PackageDescription{
				Seller: "1c",
				Query:  query,
			}
			err := s.ReceiveCards(ctx, desc)
			if err != nil {
				return err
			}
		}

		card, err := s.getCardByExternalID(ctx, select2card.ExternalID, seller.ID)
		if errors.Is(err, ErrObjectNotFound) {
			// Нам не удалось получить запись в 1с, значит данные по этому товару исчезли, пропускаем загрузку остатков
			// TODO Писать эти данные в jaeger
			notFoundElements++
			continue
		} else if err != nil {
			return err
		}
		stock.PriceSize.CardID = card.ID

		size, err := s.getSizeByTitle(ctx, stock.Size.TechSize)
		if err != nil {
			return err
		}
		stock.PriceSize.SizeID = size.ID

		_, err = s.setPriceSize(ctx, stock.PriceSize)
		if err != nil {
			return err
		}

		barcode, err := s.setBarcode(ctx, stock.Barcode)
		if err != nil {
			return err
		}

		stock.Warehouse.SellerID = seller.ID
		warehouse, err := s.setWarehouse(ctx, &stock.Warehouse)
		if err != nil {
			return err
		}

		stockData, err := s.stockrepo.SelectByBarcode(ctx, barcode.ID, desc.UpdatedAt)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении priceSize %s в БД.", "wb")
		}

		if stockData == nil {
			_, err = s.stockrepo.Insert(ctx, entity.Stock{
				Quantity:    stock.Stock.Quantity,
				BarcodeID:   barcode.ID,
				WarehouseID: warehouse.ID,
				CardID:      card.ID,
				SellerID:    seller.ID,
				CreatedAt:   desc.UpdatedAt,
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
			stockData.UpdatedAt = time.Now()

			err = s.stockrepo.UpdateExecOne(ctx, *stockData)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении stock в БД.")
			}
		}
	}

	alogger.InfoFromCtx(ctx, "Загружена информация о остатке всего: %d из них не найдено %d", len(stockMetaList), notFoundElements)

	if desc.Limit > 0 {
		p := entity.PackageDescription{
			PackageType: entity.PackageTypeStock,

			UpdatedAt: desc.UpdatedAt.Add(-24 * time.Hour),
			Limit:     desc.Limit - 1,
			Seller:    desc.Seller,
		}

		err = s.brokerPublisher.SendPackage(ctx, &p)
		if err != nil {
			return fmt.Errorf("Ошибка постановки задачи в очередь: %w", err)
		}

		alogger.InfoFromCtx(ctx, "Создана очередь на получение остатков на дату: %s", p.UpdatedAt.Format("02.01.2006"))
	} else {
		alogger.InfoFromCtx(ctx, "Все элементы обработаны")
	}

	return nil
}
