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

	attrs := make(map[string]interface{})
	attrs["количество данных"] = len(salesMetaList)
	attrs["seller"] = desc.Seller

	var notFoundElements int

	for _, sale := range salesMetaList {
		wb2card, err := s.wb2cardrepo.SelectByNmid(ctx, sale.Card.ExternalID)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении wb2card %s в БД.", "wb")
		}
		// TODO В случае отсутствия в Wb2Card - добавлять в него
		if wb2card == nil {
			notFoundElements++
			continue
		}

		seller, err := s.sellerRepo.SelectByTitle(ctx, desc.Seller)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении seller %s в БД.", "wb")
		}
		// Warehouse
		warehouse, err := s.warehouserepo.SelectBySellerIDAndTitle(ctx, seller.ID, sale.Warehouse.Title)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении warehouserepo %s в БД.", "wb")
		}

		if warehouse == nil {
			notFoundElements++
			continue
		}

		// Barcode
		barcode, err := s.barcodeRepo.SelectByBarcode(ctx, sale.Barcode.Barcode)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении barcode %s в БД.", "wb")
		}

		if barcode == nil {
			notFoundElements++
			continue
		}

		// PriceSize
		priceSize, err := s.pricesizerepo.SelectByID(ctx, barcode.PriceSizeID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении barcode %s в БД.", "wb")
		}

		if priceSize == nil {
			notFoundElements++
			continue
		}

		// Region
		country, err := s.countryrepo.SelectByName(ctx, sale.Region.Country.Name)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении country %s в БД.", "wb")
		}
		if country == nil {
			country, err = s.countryrepo.Insert(ctx, sale.Region.Country)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении country в БД.")
			}
		}

		district, err := s.districtrepo.SelectByName(ctx, sale.Region.District.Name)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении district %s в БД.", "wb")
		}
		if district == nil {
			district, err = s.districtrepo.Insert(ctx, sale.Region.District)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении district в БД.")
			}
		}

		region, err := s.regionrepo.SelectByName(ctx, sale.Region.RegionName, district.ID, country.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении region %s в БД.", "wb")
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
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении orderData %s в БД.", "wb")
		}
		if orderData == nil {
			notFoundElements++
			continue
		}

		saleData, err := s.salerepo.SelectByCardIDAndDate(ctx, wb2card.CardID, desc.UpdatedAt)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении saleData %s в БД.", "wb")
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

	attrs["не найденных элементов"] = notFoundElements
	alogger.InfoFromCtx(ctx, "Загружена информация о остатке %s", attrs)

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

		attrs["дата остатков"] = p.UpdatedAt.Format("02.01.2006")
		alogger.InfoFromCtx(ctx, "Создана очередь stocs, limit:%s", attrs)
	}

	return nil
}
