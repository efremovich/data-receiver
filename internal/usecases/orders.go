package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
	"github.com/efremovich/data-receiver/pkg/alogger"
	"github.com/efremovich/data-receiver/pkg/jaeger"
)

func (s *receiverCoreServiceImpl) ReceiveOrders(ctx context.Context, desc entity.PackageDescription) error {
	clients := s.apiFetcher[desc.Seller]

	group, gCtx := errgroup.WithContext(ctx)

	for _, c := range clients {
		client := c

		group.Go(func() error {
			if err := s.receiveAndSaveOrders(gCtx, client, desc); err != nil {
				return err
			}

			if desc.Limit > 0 {
				packet := entity.PackageDescription{
					PackageType: entity.PackageTypeOrder,

					UpdatedAt: desc.UpdatedAt.Add(-24 * time.Hour),
					Limit:     desc.Limit - 1,
					Seller:    desc.Seller,
					Delay:     desc.Delay,
				}

				err := s.brokerPublisher.SendPackage(ctx, &packet)
				if err != nil {
					return fmt.Errorf("ошибка постановки задачи в очередь: %w", err)
				}

				alogger.InfoFromCtx(ctx, "Создана очередь для получения заказов на %s", packet.UpdatedAt.Format("02.01.2006"))
			} else {
				alogger.InfoFromCtx(ctx, "Все элементы обработаны")
			}
			return nil
		})
	}

	// Ждем завершения всех горутин и проверяем наличие ошибок
	if err := group.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			alogger.WarnFromCtx(ctx, "Операция была отменена: %v", err)
			return nil
		}
		return fmt.Errorf("ошибка при обработке клиентов: %w", err)
	}

	alogger.InfoFromCtx(ctx, "постановка задачи в очередь %d", desc.Limit)

	return nil
}

func (s *receiverCoreServiceImpl) receiveAndSaveOrders(ctx context.Context, client webapi.ExtAPIFetcher, desc entity.PackageDescription) error {
	var notFoundElements int

	startTime := time.Now()
	rootSpan, ctx := jaeger.StartSpan(ctx, "receiveAndSaveOrders")

	ordersSpan, ctx := jaeger.StartSpan(ctx, "GetOrders from "+client.GetMarketPlace().Title+" "+desc.UpdatedAt.Format("02.01.2006"))

	ordersMetaList, err := client.GetOrders(ctx, desc)
	if err != nil {
		rootSpan.SetTag("error", true)
		return fmt.Errorf("ошибка при получение данных из внешнего источника %s: %w", desc.Seller, err)
	}

	ordersSpan.Finish()

	s.metricsCollector.AddReceiveReqestTime(time.Since(startTime), "orders", "receive")
	alogger.InfoFromCtx(ctx, "Получение данных о заказах за %s", desc.UpdatedAt.Format("02.01.2006"))

	sellerSpan, sellerCtx := jaeger.StartSpan(ctx, "getSeller")

	seller, err := s.getSeller(sellerCtx, client.GetMarketPlace())
	if err != nil {
		return wrapErr(fmt.Errorf("ошибка получения данных о продавце %s модуль sales:%w", desc.Seller, err))
	}

	sellerSpan.Finish()

	orders := make([]*entity.Order, 0, len(ordersMetaList))
	chunkSize := 1000 // Размер чанка

	for _, meta := range ordersMetaList {
		meta.Seller = seller

		// err = s.processSingleOrder(ctx, &meta)
		err = s.processUpdatePriceOrder(ctx, &meta)
		if err != nil {
			rootSpan.SetTag("error", true)
			return err
		}

		orders = append(orders, &meta)

		// Если накопилось 100 заказов - отправляем и очищаем срез
		if len(orders) == chunkSize {
			updateOrderSpan, _ := jaeger.StartSpan(ctx, "updateOrder Чанк 1000")
			defer updateOrderSpan.Finish()
			// err = s.updateOrders(ctx, orders)
			// if err != nil {
			// 	rootSpan.SetTag("error", true)
			// 	return err
			// }

			orders = orders[:0] // Очищаем срез (но сохраняем capacity)
		}
	}

	// Отправляем оставшиеся заказы (если есть)
	if len(orders) > 0 {
		updateOrderSpan, _ := jaeger.StartSpan(ctx, fmt.Sprintf("UpdateOrder финал чанк %d", len(orders)))
		defer updateOrderSpan.Finish()

		// err = s.updateOrders(ctx, orders)
		// if err != nil {
		// 	return err
		// }
	}

	// Если не нужно делить или деление невозможно, обрабатываем все заказы сразу
	s.metricsCollector.AddReceiveReqestTime(time.Since(startTime), "orders", "write")
	alogger.InfoFromCtx(ctx, "Загружена информация о заказах всего: %d из них не найдено %d", len(ordersMetaList), notFoundElements)

	rootSpan.Finish()

	return nil
}

func (s *receiverCoreServiceImpl) processSingleOrder(ctx context.Context, meta *entity.Order) error {
	singleOrderSpan, ctx := jaeger.StartSpan(ctx, "processSingleOrder")
	defer singleOrderSpan.Finish()
	// Проверим есть ли товар в базе, в случае отсутствия запросим его в 1с
	_, err := s.getSeller2Card(ctx, meta.Card.ExternalID, meta.Seller.ID)

	// Получаем ошибку что такой записи нет, поищем карточку товара в 1с
	if errors.Is(err, ErrObjectNotFound) {
		query := make(map[string]string)
		query["barcode"] = meta.Barcode.Barcode
		query["article"] = meta.Card.VendorID

		descOdinAss := entity.PackageDescription{
			Seller: "odinc",
			Query:  query,
		}

		span, ctx := jaeger.StartSpan(ctx, "getCardFromOdinAs")

		err := s.ReceiveCards(ctx, descOdinAss)
		if err != nil {
			return err
		}

		span.Finish()
	} else if err != nil {
		return wrapErr(fmt.Errorf("ошибка получения данных отсутствует связь между продавцом и товаром модуль stocks:%w", err))
	}

	card, err := s.getCardByVendorID(ctx, meta.Card.VendorID)

	if errors.Is(err, ErrObjectNotFound) {
		return nil
	} else if err != nil {
		return err
	}

	// Проверим и создадим связь продавца и товара
	seller2card := entity.Seller2Card{
		CardID:   card.ID,
		SellerID: meta.Seller.ID,
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

	priceSize, err := s.setPriceSize(ctx, meta.PriceSize)
	if err != nil {
		return err
	}

	meta.PriceSize = priceSize

	// Barcode
	meta.Barcode.SellerID = meta.Seller.ID
	meta.Barcode.PriceSizeID = priceSize.ID

	_, err = s.setBarcode(ctx, *meta.Barcode)
	if err != nil {
		return err
	}

	// Warehouse
	meta.Warehouse.SellerID = meta.Seller.ID

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

	return nil
}

func (s *receiverCoreServiceImpl) processUpdatePriceOrder(ctx context.Context, meta *entity.Order) error {
	singleOrderSpan, ctx := jaeger.StartSpan(ctx, "processUpdatePriceOrder")
	defer singleOrderSpan.Finish()

	order, err := s.getOrderByExternalID(ctx, meta.ExternalID)
	if errors.Is(err, ErrObjectNotFound) {
		return nil
	} else if err != nil {
		return err
	}

	// Size
	size, err := s.setSize(ctx, meta.Size)
	if err != nil {
		return err
	}

	// PriceSize
	meta.PriceSize.CardID = order.Card.ID
	meta.PriceSize.SizeID = size.ID

	priceSize, err := s.setPriceSize(ctx, meta.PriceSize)
	if err != nil {
		return err
	}

	meta.PriceSize = priceSize

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

func (s *receiverCoreServiceImpl) setOrder(ctx context.Context, income *entity.Order) (*entity.Order, error) {
	order, err := s.orderrepo.SelectByExternalID(ctx, income.ExternalID)
	if errors.Is(err, ErrObjectNotFound) {
		order, err = s.orderrepo.Insert(ctx, income)
	}

	if err != nil {
		return nil, err
	}

	if order.Price != income.Price {
		err = s.orderrepo.UpdateExecOne(ctx, order)
		if err != nil {
			return nil, err
		}
	}

	return order, nil
}
func (s *receiverCoreServiceImpl) updateOrders(ctx context.Context, newOrders []*entity.Order) error {
	// 1. Проверяем существующие заказы и фильтруем
	toInsert, err := s.orderrepo.SelectByExternalIDAndCheckPrice(ctx, newOrders)
	if err != nil {
		return fmt.Errorf("ошибка проверки заказов: %w", err)
	}

	// 2. Если есть что вставлять - выполняем пакетную вставку
	if len(toInsert) > 0 {
		_, err = s.orderrepo.InsertBatch(ctx, toInsert)
		if err != nil {
			return fmt.Errorf("ошибка вставки заказов: %w", err)
		}
	}

	return nil
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
