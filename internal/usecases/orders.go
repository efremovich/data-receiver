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

func (s *receiverCoreServiceImpl) ReceiveOrders(ctx context.Context, desc entity.PackageDescription) error {
	client := s.apiFetcher[desc.Seller]

	ordersMetaList, err := client.GetOrders(ctx, desc)
	if err != nil {
		return fmt.Errorf("ошибка при получение данных из внешнего источника %s: %w", desc.Seller, err)
	}

	var notFoundElements int

	for _, order := range ordersMetaList {
		seller, err := s.sellerRepo.SelectByTitle(ctx, desc.Seller)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении seller %s в БД.", desc.Seller)
		}

		wb2card, err := s.wb2cardrepo.SelectByNmid(ctx, order.Card.ExternalID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении wb2card %s в БД.", desc.Seller)
		}
		if wb2card == nil {
			odincClient := s.apiFetcher["odinc"]
			query := make(map[string]string)
			query["barcode"] = order.Barcode.Barcode
			query["article"] = order.Card.VendorCode

			pkg := entity.PackageDescription{
				Query: query,
			}

			cardlist, err := odincClient.GetCards(ctx, pkg)
			if err != nil {
				return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении Seller %s в БД.", desc.Seller)
			}

			if len(cardlist) == 0 {
				notFoundElements++

				alogger.InfoFromCtx(ctx, "Не найден элемент ID: %+v", order.Card.VendorCode)

				continue
			}

			for _, card := range cardlist {
				card.ExternalID = order.Card.ExternalID
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

				wb2card = &entity.Seller2Card{
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
		warehouse, err := s.warehouserepo.SelectBySellerIDAndTitle(ctx, seller.ID, order.Warehouse.Title)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении warehouserepo %s в БД.", desc.Seller)
		}

		if warehouse == nil {
			notFoundElements++

			alogger.InfoFromCtx(ctx, "Не найден склад: %s", order.Warehouse.Title)

			continue
		}

		size, err := s.sizerepo.SelectByTechSize(ctx, order.Size.TechSize)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении size %s в БД.", "wb")
		}

		if size == nil {
			size, err = s.sizerepo.Insert(ctx, entity.Size{
				TechSize: order.Size.TechSize,
				Title:    order.Size.TechSize,
			})
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении size в БД.")
			}
		}

		priceSize, err := s.pricesizerepo.SelectByCardIDAndSizeID(ctx, wb2card.CardID, size.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении цены %s в БД.", desc.Seller)
		}

		if priceSize == nil {
			priceSize, err = s.pricesizerepo.Insert(ctx, entity.PriceSize{
				Price:        order.PriceSize.Price,
				Discount:     order.PriceSize.Discount,
				SpecialPrice: order.PriceSize.SpecialPrice,
				CardID:       wb2card.CardID,
				SizeID:       size.ID,
				UpdatedAt:    order.CreatedAt,
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

		// Barcode
		barcode, err := s.barcodeRepo.SelectByBarcode(ctx, order.Barcode.Barcode)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении barcode %s в БД.", desc.Seller)
		}

		if barcode == nil {
			_, err := s.barcodeRepo.Insert(ctx, entity.Barcode{
				Barcode:     order.Barcode.Barcode,
				ExternalID:  0,
				PriceSizeID: priceSize.ID,
				SellerID:    seller.ID,
			})
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении barcode %s в БД.", order.Barcode.Barcode)
			}
		}

		// Status
		status, err := s.statusrepo.SelectByName(ctx, order.Status.Name)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении status %s в БД.", desc.Seller)
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
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении country %s в БД.", desc.Seller)
		}

		if country == nil {
			country, err = s.countryrepo.Insert(ctx, order.Region.Country)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении country в БД.")
			}
		}

		district, err := s.districtrepo.SelectByName(ctx, order.Region.District.Name)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении district %s в БД.", desc.Seller)
		}

		if district == nil {
			district, err = s.districtrepo.Insert(ctx, order.Region.District)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении district в БД.")
			}
		}

		region, err := s.regionrepo.SelectByName(ctx, order.Region.RegionName, district.ID, country.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении region %s в БД.", desc.Seller)
		}

		if region == nil {
			region, err = s.regionrepo.Insert(ctx, &entity.Region{
				RegionName: order.Region.RegionName,
				District:   *district,
				Country:    *country,
			})
			if err != nil {
				fmt.Printf("Давай разберемся почему ошибка: Имя региона: %s district: %d country %d", order.Region.RegionName, district.ID, country.ID)
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении district в БД.")
			}
		}

		orderData, err := s.orderrepo.SelectByCardIDAndDate(ctx, wb2card.CardID, desc.UpdatedAt)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении orderData %s в БД.", desc.Seller)
		}

		if orderData == nil {
			_, err = s.orderrepo.Insert(ctx, entity.Order{
				ExternalID: order.ExternalID,
				Price:      order.Price,
				Type:       order.Type,
				Direction:  order.Direction,
				Sale:       order.Sale,
				Quantity:   1,
				Status:     status,
				Region:     region,
				Warehouse:  warehouse,
				Seller:     seller,
				PriceSize:  priceSize,
				Card: &entity.Card{
					ID: wb2card.CardID,
				},
				CreatedAt: order.CreatedAt,
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

	alogger.InfoFromCtx(ctx, "Загружена информация о заказах всего: %d из них не найдено %d", len(ordersMetaList), notFoundElements)

	if desc.Limit > 0 {
		p := entity.PackageDescription{
			PackageType: entity.PackageTypeOrder,

			UpdatedAt: desc.UpdatedAt.Add(-24 * time.Hour),
			Limit:     desc.Limit - 1,
			Seller:    desc.Seller,
			Delay:     desc.Delay,
		}

		err = s.brokerPublisher.SendPackage(ctx, &p)
		if err != nil {
			return aerror.New(ctx, entity.BrokerSendErrorID, err, "Ошибка постановки задачи в очередь")
		}

		alogger.InfoFromCtx(ctx, "Создана очередь для получения заказов на %s", p.UpdatedAt.Format("02.01.2006"))
	}

	return nil
}
