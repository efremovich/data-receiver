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
			return s.receiveAndSaveOrders(gCtx, client, desc)
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
}

func (s *receiverCoreServiceImpl) receiveAndSaveOrders(ctx context.Context, client webapi.ExtAPIFetcher, desc entity.PackageDescription) error {
	var notFoundElements int

	startTime := time.Now()
	rootSpan, spanCtx := jaeger.StartSpan(ctx, "receiveAndSaveOrders")

	defer rootSpan.Finish()

	rootSpan.SetTag("seller", desc.Seller)
	rootSpan.SetTag("date", desc.UpdatedAt.Format("02.01.2006"))

	ordersSpan, ctx := jaeger.StartSpan(spanCtx, "GetOrders")
	defer ordersSpan.Finish()

	ordersMetaList, err := client.GetOrders(ctx, desc)
	if err != nil {
		rootSpan.SetTag("error", true)
		return fmt.Errorf("ошибка при получение данных из внешнего источника %s: %w", desc.Seller, err)
	}

	s.metricsCollector.AddReceiveReqestTime(time.Since(startTime), "orders", "receive")
	alogger.InfoFromCtx(ctx, "Получение данных о заказах за %s", desc.UpdatedAt.Format("02.01.2006"))

	sellerSpan, ctx := jaeger.StartSpan(spanCtx, "getSeller")
	defer sellerSpan.Finish()

	seller, err := s.getSeller(ctx, client.GetMarketPlace())
	if err != nil {
		return wrapErr(fmt.Errorf("ошибка получения данных о продавце %s модуль sales:%w", desc.Seller, err))
	}

	for iter, meta := range ordersMetaList {
		if iter > 5 {
			continue
		}

		meta.Seller = seller
		orderSpan, orderCtx := jaeger.StartSpan(spanCtx, fmt.Sprintf("ProcessOrder-%d", iter))

		err = s.processSingleOrder(orderCtx, &meta)
		if err != nil {
			orderSpan.SetTag("error", true)
			orderSpan.Finish()
		}

		ordersSpan.Finish()
	}

	s.metricsCollector.AddReceiveReqestTime(time.Since(startTime), "orders", "write")
	alogger.InfoFromCtx(ctx, "Загружена информация о заказах всего: %d из них не найдено %d", len(ordersMetaList), notFoundElements)

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

	spanCardId, ctx := jaeger.StartSpan(ctx, "getCardByVendorID")
	card, err := s.getCardByVendorID(ctx, meta.Card.VendorID)

	if errors.Is(err, ErrObjectNotFound) {
		return nil
	} else if err != nil {
		return err
	}

	spanCardId.Finish()
	// Проверим и создадим связь продавца и товара
	seller2card := entity.Seller2Card{
		CardID:   card.ID,
		SellerID: meta.Seller.ID,
	}

	spanSeller2Card, ctx := jaeger.StartSpan(ctx, "setSeller2Card")

	_, err = s.setSeller2Card(ctx, seller2card)
	if err != nil {
		return err
	}

	meta.Card = card

	spanSeller2Card.Finish()

	// Size
	spanSize, ctx := jaeger.StartSpan(ctx, "setSize")

	size, err := s.setSize(ctx, meta.Size)
	if err != nil {
		return err
	}

	spanSize.Finish()

	// PriceSize
	meta.PriceSize.CardID = card.ID
	meta.PriceSize.SizeID = size.ID

	spanPriceSize, ctx := jaeger.StartSpan(ctx, "setPriceSize")

	priceSize, err := s.setPriceSize(ctx, *meta.PriceSize)
	if err != nil {
		return err
	}

	meta.PriceSize = priceSize

	spanPriceSize.Finish()

	// Barcode
	meta.Barcode.SellerID = meta.Seller.ID
	meta.Barcode.PriceSizeID = priceSize.ID

	spanSetBarcode, ctx := jaeger.StartSpan(ctx, "setBarcode")

	_, err = s.setBarcode(ctx, *meta.Barcode)
	if err != nil {
		return err
	}

	spanSetBarcode.Finish()

	// Warehouse
	meta.Warehouse.SellerID = meta.Seller.ID

	spanSetWarehouse, ctx := jaeger.StartSpan(ctx, "setWarehouse")

	warehouse, err := s.setWarehouse(ctx, meta.Warehouse)
	if err != nil {
		return wrapErr(fmt.Errorf("ошибка получения данных по складам хранения модуль orders:%w", err))
	}

	spanSetWarehouse.Finish()

	meta.Warehouse = warehouse

	// Status
	spanSetStatus, ctx := jaeger.StartSpan(ctx, "setStatus")

	status, err := s.setStatus(ctx, meta.Status)
	if err != nil {
		return wrapErr(fmt.Errorf("ошибка получения данных по статусу заказа модуль orders:%w", err))
	}

	meta.Status = status

	spanSetStatus.Finish()

	// Region
	spanSetCountry, ctx := jaeger.StartSpan(ctx, "setCountry")

	country, err := s.setCountry(ctx, meta.Region.Country)
	if err != nil {
		return err
	}

	meta.Region.Country.ID = country.ID

	spanSetCountry.Finish()

	spanSetDistrict, ctx := jaeger.StartSpan(ctx, "setDistrict")

	district, err := s.setDistrict(ctx, meta.Region.District)
	if err != nil {
		return err
	}

	meta.Region.District.ID = district.ID

	spanSetDistrict.Finish()

	spanSetRegion, ctx := jaeger.StartSpan(ctx, "setRegion")

	region, err := s.setRegion(ctx, *meta.Region)
	if err != nil {
		return err
	}

	meta.Region = region

	spanSetRegion.Finish()

	spanSetOrder, ctx := jaeger.StartSpan(ctx, "setOrder")

	_, err = s.setOrder(ctx, meta)
	if err != nil {
		return err
	}

	spanSetOrder.Finish()

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
