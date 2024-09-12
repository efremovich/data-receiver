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

func (s *receiverCoreServiceImpl) ReceiveSales(ctx context.Context, desc entity.PackageDescription) aerror.AError {
	client := s.apiFetcher[desc.Seller]

	salesMetaList, err := client.GetSales(ctx, desc)
	if err != nil {
		return aerror.New(ctx, entity.GetDataFromExSources, err, "ошибка получение данные из внешнего источника %s в БД: %s ", "", err.Error())
	}

	var notFoundElements int

	for _, sale := range salesMetaList {
		seller, err := s.sellerRepo.SelectByTitle(ctx, desc.Seller)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении seller %s в БД.", desc.Seller)
		}

		wb2card, err := s.wb2cardrepo.SelectByNmid(ctx, sale.Card.ExternalID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении wb2card %s в БД.", desc.Seller)
		}

		if wb2card == nil {
			odincClient := s.apiFetcher["odinc"]
			query := make(map[string]string)
			query["barcode"] = sale.Barcode.Barcode
			query["article"] = sale.Card.VendorCode

			pkg := entity.PackageDescription{
				Query: query,
			}

			cardlist, err := odincClient.GetCards(ctx, pkg)
			if err != nil {
				return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении Seller %s в БД.", desc.Seller)
			}

			if len(cardlist) == 0 {
				notFoundElements++

				alogger.InfoFromCtx(ctx, "Не найден элемент ID: %+v", sale.Card.VendorCode)

				continue
			}

			for _, card := range cardlist {
				card.ExternalID = sale.Card.ExternalID
				// Brands
				brand, err := s.getBrand(ctx, card.Brand, seller)
				if err != nil {
					return err
				}

				card.Brand = *brand

				err = s.setWb2Card(ctx, &card)
				if err != nil {
					notFoundElements++

					alogger.ErrorFromCtx(ctx, "Ошибка записи элемента в базу %s", err.Error())

					continue
				}

				wb2card = &entity.Wb2Card{
					NMID:   card.ExternalID,
					KTID:   0,
					NMUUID: "",
					CardID: card.ID,
				}
				_, aerr := s.wb2cardrepo.Insert(ctx, *wb2card)

				if aerr != nil {
					notFoundElements++

					alogger.ErrorFromCtx(ctx, "Ошибка записи элемента в базу %s", aerr.Error())

					continue
				}
			}
		}

		// Warehouse
		warehouse, err := s.warehouserepo.SelectBySellerIDAndTitle(ctx, seller.ID, sale.Warehouse.Title)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении warehouserepo %s в БД.", sale.Warehouse.Title)
		}

		if warehouse == nil {
			notFoundElements++
			continue
		}

		// Barcode
		barcode, err := s.barcodeRepo.SelectByBarcode(ctx, sale.Barcode.Barcode)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении barcode %s в БД.", sale.Barcode.Barcode)
		}

		if barcode == nil {
			notFoundElements++
			continue
		}

		// PriceSize
		priceSize, err := s.pricesizerepo.SelectByID(ctx, barcode.PriceSizeID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении barcode %d в БД.", barcode.PriceSizeID)
		}

		if priceSize == nil {
			notFoundElements++
			continue
		}

		// Region
		country, err := s.countryrepo.SelectByName(ctx, sale.Region.Country.Name)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении country %s в БД.", sale.Region.Country.Name)
		}

		if country == nil {
			country, err = s.countryrepo.Insert(ctx, sale.Region.Country)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении country в БД.")
			}
		}

		district, err := s.districtrepo.SelectByName(ctx, sale.Region.District.Name)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении district %s в БД.", sale.Region.District.Name)
		}

		if district == nil {
			district, err = s.districtrepo.Insert(ctx, sale.Region.District)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении district в БД.")
			}
		}

		region, err := s.regionrepo.SelectByName(ctx, sale.Region.RegionName, district.ID, country.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении region %s в БД.", sale.Region.RegionName)
		}

		if region == nil {
			_, err := s.regionrepo.Insert(ctx, &entity.Region{
				RegionName: sale.Region.RegionName,
				District:   *district,
				Country:    *country,
			})
			if err != nil {
				fmt.Printf("Давай разберемся почему ошибка: Имя региона: %s district: %d country %d", sale.Region.RegionName, district.ID, country.ID)
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении district в БД.")
			}
		}

		orderData, err := s.orderrepo.SelectByExternalID(ctx, sale.Order.ExternalID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении orderData %s в БД.", sale.Order.ExternalID)
		}

		if orderData == nil {
			notFoundElements++
			continue
		}

		saleData, err := s.salerepo.SelectByCardIDAndDate(ctx, wb2card.CardID, desc.UpdatedAt)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении saleData %d в БД.", wb2card.CardID)
		}

		if saleData == nil {
			_, err = s.salerepo.Insert(ctx, entity.Sale{
				ExternalID: sale.ExternalID,

				Price:      sale.Price,
				DiscountP:  sale.DiscountP,
				FinalPrice: sale.FinalPrice,
				Type:       sale.Type,
				ForPay:     sale.ForPay,

				Quantity:  sale.Quantity,
				CreatedAt: sale.CreatedAt,

				Order:  orderData,
				Seller: seller,
				Card: &entity.Card{
					ID: wb2card.CardID,
				},
				Warehouse: warehouse,
				Region:    region,
				PriceSize: priceSize,
			})
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении sale в БД.")
			}
		} else {
			saleData.UpdatedAt = time.Now()

			err = s.salerepo.UpdateExecOne(ctx, saleData)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при обновлении sale в БД.")
			}
		}
	}

	alogger.InfoFromCtx(ctx, "Загружена информация о продажах всего: %d из них не найдено %d", len(salesMetaList), notFoundElements)

	if desc.Limit > 0 {
		p := entity.PackageDescription{
			PackageType: entity.PackageTypeSale,

			UpdatedAt: desc.UpdatedAt.Add(-24 * time.Hour),
			Limit:     desc.Limit - 1,
			Seller:    desc.Seller,
		}

		err = s.brokerPublisher.SendPackage(ctx, &p)
		if err != nil {
			return aerror.New(ctx, entity.BrokerSendErrorID, err, "Ошибка постановки задачи в очередь")
		}

		alogger.InfoFromCtx(ctx, "Создана очередь для получения продаж на %s", p.UpdatedAt.Format("02.01.2006"))
	}

	return nil
}
