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

func (s *receiverCoreServiceImpl) ReceiveOrders(ctx context.Context, desc entity.PackageDescription) aerror.AError {
	client := s.apiFetcher[desc.Seller]

	ordersMetaList, err := client.GetOrders(ctx, desc)
	if err != nil {
		return aerror.New(ctx, entity.GetDataFromExSources, err, "ошибка получение данные из внешнего источника %s в БД: %s ", "", err.Error())
	}

	attrs := make(map[string]interface{})
	attrs["количество данных"] = len(ordersMetaList)
	attrs["seller"] = desc.Seller

	var notFoundElements int

	for i, order := range ordersMetaList {
		wb2card, err := s.wb2cardrepo.SelectByNmid(ctx, order.Card.ExternalID)

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
		warehouse, err := s.warehouserepo.SelectBySellerIDAndTitle(ctx, seller.ID, order.Warehouse.Title)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении warehouserepo %s в БД.", "wb")
		}

		if warehouse == nil {
			notFoundElements++
			continue
		}

		// Barcode
		barcode, err := s.barcodeRepo.SelectByBarcode(ctx, order.Barcode.Barcode)
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

		// Status
		status, err := s.statusrepo.SelectByName(ctx, order.Status.Name)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении status %s в БД.", "wb")
		}
		if status == nil {
			status, err = s.statusrepo.Insert(ctx, *order.Status)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении stock в БД.")
			}
		}

		// Region
		country, err := s.countryrepo.SelectByName(ctx, order.Region.Country.Name)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении country %s в БД.", "wb")
		}
		if country == nil {
			country, err = s.countryrepo.Insert(ctx, order.Region.Country)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении country в БД.")
			}
		}

		district, err := s.districtrepo.SelectByName(ctx, order.Region.District.Name)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении district %s в БД.", "wb")
		}
		if district == nil {
			district, err = s.districtrepo.Insert(ctx, order.Region.District)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении district в БД.")
			}
		}

		region, err := s.regionrepo.SelectByName(ctx, order.Region.RegionName, district.ID, country.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении region %s в БД.", "wb")
		}
		if region == nil {
			region, err = s.regionrepo.Insert(ctx, &entity.Region{
				RegionName: order.Region.RegionName,
				District:   *district,
				Country:    *country,
			})
			if err != nil {
				fmt.Println(i)
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении district в БД.")
			}
		}

		orderData, err := s.orderrepo.SelectByCardIDAndDate(ctx, wb2card.CardID, desc.UpdatedAt)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении orderData %s в БД.", "wb")
		}
		if orderData == nil {
			_, err = s.orderrepo.Insert(ctx, entity.Order{
				ExternalID: order.ExternalID,
				Price:      order.Price,
				Type:       order.Type,
				Direction:  order.Direction,
				Sale:       order.Sale,
				Quantity:   1,

				Status:    status,
				Region:    region,
				Warehouse: warehouse,
				Seller:    seller,
				PriceSize: priceSize,
				Card: &entity.Card{
					ID: wb2card.CardID,
				},
				CreatedAt: desc.UpdatedAt,
			})
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении stock в БД.")
			}
		} else {
			orderData.UpdatedAt = time.Now()
			err = s.orderrepo.UpdateExecOne(ctx, *orderData)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении stock в БД.")
			}
		}
	}

	attrs["не найденных элементов"] = notFoundElements
	alogger.InfoFromCtx(ctx, "Загружена информация о остатке", nil, attrs, false)

	if desc.Limit > 0 {
		p := entity.PackageDescription{
			PackageType: entity.PackageTypeOrder,

			UpdatedAt: desc.UpdatedAt.Add(-24 * time.Hour),
			Limit:     desc.Limit - 1,
			Seller:    desc.Seller,
		}
		err = s.brokerPublisher.SendPackage(ctx, &p)
		if err != nil {
			return aerror.New(ctx, entity.BrokerSendErrorID, err, "Ошибка постановки задачи в очередь")
		}

		attrs["дата остатков"] = p.UpdatedAt.Format("02.01.2006")
		alogger.InfoFromCtx(ctx, "Создана очередь stocs, limit:%1", nil, attrs, false)
	}

	return nil
}
