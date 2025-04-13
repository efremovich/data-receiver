package controller

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/efremovich/data-receiver/internal/entity"
	package_receiver "github.com/efremovich/data-receiver/pkg/data-receiver-service"
	"github.com/efremovich/data-receiver/pkg/logger"
)

func (gw *grpcGatewayServerImpl) OfferFeed(ctx context.Context, _ *emptypb.Empty) (*package_receiver.OfferFeedResponse, error) {
	responce, err := gw.core.OfferFeed(ctx)

	offerResp := package_receiver.OfferFeedResponse{
		Body: string(responce),
	}

	return &offerResp, err
}

func (gw *grpcGatewayServerImpl) runTask(ctx context.Context) {
	err := gw.receiveCostFrom1C(ctx)
	if err != nil {
		logger.GetLoggerFromContext(ctx).Errorf("ошибка при получении отчета о продажах:%s", err.Error())
	}
}

type Task struct {
	Name     string
	Interval string
	Function func(ctx context.Context) error
}

func (gw *grpcGatewayServerImpl) scheduleTasks(ctx context.Context) {
	c := cron.New()
	tasks := []Task{
		{"Загрузка товарных позиций wildberries", "0 12 * * *", gw.receiveCardsWB}, // Каждый день в 12
		{"Загрузка товарных позиций ozon", "05 12 * * *", gw.receiveCardsOzon},

		{"Загрузка складов wildberries", "15 12 * * *", gw.receiveWarehousesWB},

		{"Загрузка остатков wildberries", "30 13 * * *", gw.receiveStocksWB},
		{"Загрузка остатков ozon", "0 13 * * *", gw.receiveStocksOzon},

		{"Загрузка заказов wildberries", "30 18 * * *", gw.receiveOrdersWB},
		{"Загрузка заказов ozon", "0 16 * * *", gw.receiveOrdersOzon},

		{"Загрузка продаж wildberries", "30 19 * * *", gw.receiveSalesWB},
		{"Загрузка продаж ozon", "30 19 * * *", gw.receiveSalesOzon},

		{"Загрузка отчета по продажам wildberries", "30 19 * * *", gw.receiveSaleReportWB},
		{"Загрузка отчета по продажам ozon", "30 19 * * *", gw.receiveSaleReportOzon},

		{"Загрузка себестоимости товара из 1с", "30 17 * * *", gw.receiveCostFrom1C},
	}

	for _, task := range tasks {
		_, err := c.AddFunc(task.Interval, func() {
			if err := task.Function(ctx); err != nil {
				logger.GetLoggerFromContext(ctx).Errorf("задание %s завершилось с ошибкой: %v", task.Name, err)
			} else {
				logger.GetLoggerFromContext(ctx).Infof("задача %s успешно завершена", task.Name)
			}
		})
		if err != nil {
			logger.GetLoggerFromContext(ctx).Errorf("Failed to schedule task %s: %v", task.Name, err)
		}
	}

	c.Start()
	select {} // Block forever
}

func (gw *grpcGatewayServerImpl) receiveCardsWB(ctx context.Context) error {
	limit := 100
	desc := entity.PackageDescription{
		Limit:       limit,
		Cursor:      "0",
		PackageType: entity.PackageTypeCard,
		Seller:      entity.Wildberries,
	}

	return gw.core.ReceiveCards(ctx, desc)
}

func (gw *grpcGatewayServerImpl) receiveCardsOzon(ctx context.Context) error {
	limit := 100
	desc := entity.PackageDescription{
		Limit:       limit,
		Cursor:      "0",
		PackageType: entity.PackageTypeCard,
		Seller:      entity.Ozon,
	}

	return gw.core.ReceiveCards(ctx, desc)
}

func (gw *grpcGatewayServerImpl) receiveWarehousesWB(ctx context.Context) error {
	limit := 100
	descWarehouse := entity.PackageDescription{
		Limit:       limit,
		Cursor:      "0",
		PackageType: entity.PackageTypeCard,
		Seller:      entity.Wildberries,
	}

	return gw.core.ReceiveWarehouses(ctx, descWarehouse)
}

func (gw *grpcGatewayServerImpl) receiveStocksWB(ctx context.Context) error {
	daysToGet := 5 // Количество дней для загрузки
	descStocks := entity.PackageDescription{
		PackageType: entity.PackageTypeStock,
		UpdatedAt:   time.Now(),
		Limit:       daysToGet,
		Seller:      entity.Wildberries,
	}

	return gw.core.ReceiveStocks(ctx, descStocks)
}

func (gw *grpcGatewayServerImpl) receiveStocksOzon(ctx context.Context) error {
	descOzonStock := entity.PackageDescription{
		PackageType: entity.PackageTypeStock,
		UpdatedAt:   time.Now(),
		Seller:      entity.Ozon,
	}

	return gw.core.ReceiveStocks(ctx, descOzonStock)
}

func (gw *grpcGatewayServerImpl) receiveOrdersWB(ctx context.Context) error {
	daysToGet := 30 // Количество дней для загрузки
	delay := 61     // Количество секунд задержки перед следующим запросом
	descOrderOzon := entity.PackageDescription{
		PackageType: entity.PackageTypeOrder,
		UpdatedAt:   time.Now(),
		Seller:      entity.Wildberries,
		Limit:       daysToGet,
		Delay:       delay,
	}

	return gw.core.ReceiveOrders(ctx, descOrderOzon)
}

func (gw *grpcGatewayServerImpl) receiveOrdersOzon(ctx context.Context) error {
	daysToGet := 30 // Количество дней для загрузки
	delay := 61     // Количество секунд задержки перед следующим запросом
	descOrderOzon := entity.PackageDescription{
		PackageType: entity.PackageTypeOrder,
		UpdatedAt:   time.Now(),
		Seller:      entity.Ozon,
		Limit:       daysToGet,
		Delay:       delay,
	}

	return gw.core.ReceiveOrders(ctx, descOrderOzon)
}

func (gw *grpcGatewayServerImpl) receiveSalesWB(ctx context.Context) error {
	daysToGet := 30 // Количество дней для загрузки
	delay := 61     // Количество секунд задержки перед следующим запросом

	descDescription := entity.PackageDescription{
		PackageType: entity.PackageTypeSale,
		UpdatedAt:   time.Now(),
		Seller:      entity.Wildberries,
		Limit:       daysToGet,
		Delay:       delay,
	}

	return gw.core.ReceiveSales(ctx, descDescription)
}

func (gw *grpcGatewayServerImpl) receiveSalesOzon(ctx context.Context) error {
	daysToGet := 30 // Количество дней для загрузки
	delay := 61     // Количество секунд задержки перед следующим запросом

	descDescription := entity.PackageDescription{
		PackageType: entity.PackageTypeSale,
		UpdatedAt:   time.Now(),
		Seller:      entity.Ozon,
		Limit:       daysToGet,
		Delay:       delay,
	}

	return gw.core.ReceiveSales(ctx, descDescription)
}

func (gw *grpcGatewayServerImpl) receiveSaleReportWB(ctx context.Context) error {
	daysToGet := 30 // Количество дней для загрузки
	delay := 61     // Количество секунд задержки перед следующим запросом
	startDate := time.Now()
	// startDate := time.Date(2025, 03, 01, 0, 0, 0, 0, time.Local)
	descDescription := entity.PackageDescription{
		PackageType: entity.PackageTypeSaleReports,
		UpdatedAt:   startDate,
		Seller:      entity.Wildberries,
		Limit:       daysToGet,
		Delay:       delay,
	}

	return gw.core.ReceiveSaleReport(ctx, descDescription)
}

func (gw *grpcGatewayServerImpl) receiveSaleReportOzon(ctx context.Context) error {
	daysToGet := 30 // Количество дней для загрузки
	delay := 61     // Количество секунд задержки перед следующим запросом
	startDate := time.Now()
	descDescription := entity.PackageDescription{
		PackageType: entity.PackageTypeSaleReports,
		UpdatedAt:   startDate,
		Seller:      entity.Ozon,
		Limit:       daysToGet,
		Delay:       delay,
	}

	return gw.core.ReceiveSaleReport(ctx, descDescription)
}

func (gw *grpcGatewayServerImpl) receiveCostFrom1C(ctx context.Context) error {
	daysToGet := 30
	delay := 3
	startDate := time.Now()
	descDescription := entity.PackageDescription{
		PackageType: entity.PackageTypeCostFrom1C,
		UpdatedAt:   startDate,
		Seller:      entity.OdinAss,
		Limit:       daysToGet,
		Delay:       delay,
	}
	return gw.core.ReceiveCostFrom1C(ctx, descDescription)
}
