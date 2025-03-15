package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
	"github.com/efremovich/data-receiver/pkg/alogger"
	"golang.org/x/sync/errgroup"
)

func (s *receiverCoreServiceImpl) ReceiveSales(ctx context.Context, desc entity.PackageDescription) error {
	clients := s.apiFetcher[desc.Seller]

	g, gCtx := errgroup.WithContext(ctx)

	for _, c := range clients {
		client := c

		g.Go(func() error {
			return s.receiveAndSaveSales(gCtx, client, desc)
		})
	}

	// Ждем завершения всех горутин и проверяем наличие ошибок
	if err := g.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			alogger.WarnFromCtx(ctx, "Операция была отменена: %v", err)
			return nil
		}
		return fmt.Errorf("ошибка при обработке клиентов: %w", err)
	}
	if desc.Limit > 0 {
		p := entity.PackageDescription{
			PackageType: entity.PackageTypeSale,

			UpdatedAt: desc.UpdatedAt.Add(-24 * time.Hour),
			Limit:     desc.Limit - 1,
			Seller:    desc.Seller,
			Delay:     desc.Delay,
		}

		alogger.InfoFromCtx(ctx, "постановка задачи в очередь %d", desc.Limit)

		err := s.brokerPublisher.SendPackage(ctx, &p)
		if err != nil {
			return fmt.Errorf("ошибка постановки задачи в очередь: %w", err)
		}

		alogger.InfoFromCtx(ctx, "Создана очередь для получения продаж на %s", p.UpdatedAt.Format("02.01.2006"))
	} else {
		alogger.InfoFromCtx(ctx, "Все элементы обработаны")
	}

	return nil
}

func (s *receiverCoreServiceImpl) receiveAndSaveSales(ctx context.Context, client webapi.ExtAPIFetcher, desc entity.PackageDescription) error {
	var notFoundElements int

	salesMetaList, err := client.GetSales(ctx, desc)
	if err != nil {
		return fmt.Errorf("ошибка получение данных о продажах из внешнего источника %s, %w", desc.Seller, err)
	}

	alogger.InfoFromCtx(ctx, "Получение данных об продажах за %s", desc.UpdatedAt.Format("02.01.2006"))

	seller, err := s.getSeller(ctx, client.GetMarketPlace())
	if err != nil {
		return wrapErr(fmt.Errorf("ошибка получения данных о продавце %s модуль sales:%w", desc.Seller, err))
	}

	for _, meta := range salesMetaList {

		meta.Seller = seller
		// Проверим есть ли товар в базе, в случае отсутствия запросим его в 1с
		_, err = s.getSeller2Card(ctx, meta.Card.ExternalID, seller.ID)
		// Получаем ошибку что такой записи нет, поищем карточку товара в 1с
		if errors.Is(err, ErrObjectNotFound) {
			query := make(map[string]string)
			query["barcode"] = meta.Barcode.Barcode
			query["article"] = meta.Card.VendorID

			in := entity.PackageDescription{
				PackageType: desc.PackageType,
				Seller:      "odinc",
				Query:       query,
			}

			err := s.ReceiveCards(ctx, in)
			if err != nil {
				return err
			}
		} else if err != nil {
			return wrapErr(fmt.Errorf("ошибка получения данных отсутствует связь между продавцом и товаром %s модуль stocks:%w", desc.Seller, err))
		}

		card, err := s.getCardByVendorID(ctx, meta.Card.VendorID)
		if errors.Is(err, ErrObjectNotFound) {
			// Нам не удалось получить запись, значит данные по этому товару исчезли, пропускаем загрузку остатков
			// TODO Писать эти данные в jaeger
			notFoundElements++
			continue
		} else if err != nil {
			return err
		}
		// Проверим и создадим связь продавца и товара
		seller2card := entity.Seller2Card{
			CardID:   card.ID,
			SellerID: seller.ID,
		}
		meta.Card = card
		_, err = s.setSeller2Card(ctx, seller2card)
		if err != nil {
			return err
		}

		// Size
		size, err := s.getSizeByTitle(ctx, meta.Size.TechSize)
		if err != nil {
			return err
		}

		// Добавление данных
		// PriceSize
		meta.PriceSize.CardID = card.ID
		meta.PriceSize.SizeID = size.ID

		priceSize, err := s.setPriceSize(ctx, *meta.PriceSize)
		if err != nil {
			return err
		}
		meta.PriceSize = priceSize

		// Barcode
		meta.Barcode.SellerID = seller.ID
		meta.Barcode.PriceSizeID = priceSize.ID

		_, err = s.setBarcode(ctx, *meta.Barcode)
		if err != nil {
			return err
		}

		// Warehouse
		meta.Warehouse.SellerID = seller.ID
		warehouse, err := s.setWarehouse(ctx, meta.Warehouse)

		if err != nil {
			return wrapErr(fmt.Errorf("ошибка получения данных по складам хранения модуль stock:%w", err))
		}

		meta.Warehouse = warehouse

		// Region
		country, err := s.setCountry(ctx, meta.Region.Country)
		if err != nil {
			return err
		}

		meta.Region.Country.ID = country.ID

		district, err := s.setDistrict(ctx, meta.Region.District)
		if err != nil {
			return err
		}

		meta.Region.District.ID = district.ID

		region, err := s.setRegion(ctx, *meta.Region)
		if err != nil {
			return err
		}

		meta.Region = region

		// Order
		order, err := s.getOrderByExternalID(ctx, meta.Order.ExternalID)
		if errors.Is(err, ErrObjectNotFound) {
			// Нам не удалось получить запись, значит данные по этому товару исчезли, пропускаем загрузку остатков
			// TODO Писать эти данные в jaeger
			notFoundElements++
			continue
		} else if err != nil {
			return err
		}

		meta.Order = order

		_, err = s.setSale(ctx, &meta)

		if err != nil {
			return err
		}
	}

	alogger.InfoFromCtx(ctx, "Загружена информация о продажах всего: %d из них не найдено %d", len(salesMetaList), notFoundElements)

	return nil
}

func (s *receiverCoreServiceImpl) setSale(ctx context.Context, in *entity.Sale) (*entity.Sale, error) {
	sale, err := s.salerepo.SelectByExternalID(ctx, in.ExternalID, in.CreatedAt)
	if errors.Is(err, ErrObjectNotFound) {
		sale, err = s.salerepo.Insert(ctx, *in)
	}

	if err != nil {
		return nil, err
	}

	return sale, nil
}
