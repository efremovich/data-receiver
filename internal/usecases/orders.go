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

func (s *receiverCoreServiceImpl) ReceiveOrders(ctx context.Context, desc entity.PackageDescription) error {
	clients := s.apiFetcher[desc.Seller]

	g, gCtx := errgroup.WithContext(ctx)

	for _, c := range clients {
		client := c

		g.Go(func() error {
			return s.receiveAndSaveOrders(gCtx, client, desc)
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

	alogger.InfoFromCtx(ctx, "постановка задачи в очередь %d", desc.Limit)

	if desc.Limit > 0 {
		p := entity.PackageDescription{
			PackageType: entity.PackageTypeOrder,

			UpdatedAt: desc.UpdatedAt.Add(-24 * time.Hour),
			Limit:     desc.Limit - 1,
			Seller:    desc.Seller,
			Delay:     desc.Delay,
		}

		err := s.brokerPublisher.SendPackage(ctx, &p)
		if err != nil {
			return fmt.Errorf("ошибка постановки задачи в очередь: %w", err)
		}

		alogger.InfoFromCtx(ctx, "Создана очередь для получения заказов на %s", p.UpdatedAt.Format("02.01.2006"))
	} else {
		alogger.InfoFromCtx(ctx, "Все элементы обработаны")
	}

	return nil
}

func (s *receiverCoreServiceImpl) receiveAndSaveOrders(ctx context.Context, client webapi.ExtAPIFetcher, desc entity.PackageDescription) error {
	var notFoundElements int

	startTime := time.Now()

	ordersMetaList, err := client.GetOrders(ctx, desc)
	if err != nil {
		return fmt.Errorf("ошибка при получение данных из внешнего источника %s: %w", desc.Seller, err)
	}

	s.metricsCollector.AddReceiveReqestTime(time.Since(startTime), "orders", "receive")
	alogger.InfoFromCtx(ctx, "Получение данных о заказах за %s", desc.UpdatedAt.Format("02.01.2006"))

	seller, err := s.getSeller(ctx, client.GetMarketPlace())
	if err != nil {
		return wrapErr(fmt.Errorf("ошибка получения данных о продавце %s модуль sales:%w", desc.Seller, err))
	}

	for _, meta := range ordersMetaList {

		meta.Seller = seller

		// Проверим есть ли товар в базе, в случае отсутствия запросим его в 1с
		_, err = s.getSeller2Card(ctx, meta.Card.ExternalID, seller.ID)
		// Получаем ошибку что такой записи нет, поищем карточку товара в 1с
		if errors.Is(err, ErrObjectNotFound) {
			query := make(map[string]string)
			query["barcode"] = meta.Barcode.Barcode
			query["article"] = meta.Card.VendorID

			descOdinAss := entity.PackageDescription{
				Seller: "odinc",
				Query:  query,
			}

			err := s.ReceiveCards(ctx, descOdinAss)
			if err != nil {
				return err
			}
		} else if err != nil {
			return wrapErr(fmt.Errorf("ошибка получения данных отсутствует связь между продавцом и товаром модуль stocks:%w", err))
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

		_, err = s.setSeller2Card(ctx, seller2card)
		if err != nil {
			return err
		}

		meta.Card = card

		// Size
		size, err := s.setSize(ctx, meta.Size)
		if err != nil {
			return err
		}

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
			return wrapErr(fmt.Errorf("ошибка получения данных по складам хранения модуль orders:%w", err))
		}

		meta.Warehouse = warehouse

		// Status
		status, err := s.setStatus(ctx, meta.Status)
		if err != nil {
			return wrapErr(fmt.Errorf("ошибка получения данных по статусу заказа модуль orders:%w", err))
		}

		meta.Status = status

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

		_, err = s.setOrder(ctx, &meta)
		if err != nil {
			return err
		}
	}

	s.metricsCollector.AddReceiveReqestTime(time.Since(startTime), "orders", "write")
	alogger.InfoFromCtx(ctx, "Загружена информация о заказах всего: %d из них не найдено %d", len(ordersMetaList), notFoundElements)

	return nil
}

func (s *receiverCoreServiceImpl) getOrderByExternalID(ctx context.Context, externalID string) (*entity.Order, error) {
	order, err := s.orderrepo.SelectByExternalID(ctx, externalID)
	if err != nil && errors.Is(err, entity.ErrObjectNotFound) {
		return &entity.Order{}, nil
	} else if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *receiverCoreServiceImpl) setOrder(ctx context.Context, in *entity.Order) (*entity.Order, error) {
	order, err := s.orderrepo.SelectByExternalID(ctx, in.ExternalID)
	if errors.Is(err, ErrObjectNotFound) {
		order, err = s.orderrepo.Insert(ctx, *in)
	}

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *receiverCoreServiceImpl) setStatus(ctx context.Context, in *entity.Status) (*entity.Status, error) {
	status, err := s.statusrepo.SelectByName(ctx, in.Name)
	if errors.Is(err, ErrObjectNotFound) {
		status, err = s.statusrepo.Insert(ctx, *in)
	}

	if err != nil {
		return nil, err
	}
	return status, nil
}
