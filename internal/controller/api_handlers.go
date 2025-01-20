package controller

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/logger"
)

func (gw *grpcGatewayServerImpl) runTask(ctx context.Context) {
	gw.receiveCardsWB(ctx)
	gw.receiveCardsOzon(ctx)
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

		{"Загрузка остатков ozon", "0 13 * * *", gw.receiveStocksOzon},
		{"Загрузка остатков wildberries", "30 13 * * *", gw.receiveStocksWB},

		{"Загрузка заказов wildberries", "30 18 * * *", gw.receiveOrdersWB},
		{"Загрузка заказов ozon", "0 16 * * *", gw.receiveOrdersOzon},

		{"Загрузка продаж wildberries", "30 19 * * *", gw.receiveSalesWB},
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
		Seller:      "wb",
	}

	return gw.core.ReceiveCards(ctx, desc)
}

func (gw *grpcGatewayServerImpl) receiveCardsOzon(ctx context.Context) error {
	limit := 100
	desc := entity.PackageDescription{
		Limit:       limit,
		Cursor:      "0",
		PackageType: entity.PackageTypeCard,
		Seller:      "ozon",
	}

	return gw.core.ReceiveCards(ctx, desc)
}

func (gw *grpcGatewayServerImpl) receiveWarehousesWB(ctx context.Context) error {
	limit := 100
	descWarehouse := entity.PackageDescription{
		Limit:       limit,
		Cursor:      "0",
		PackageType: entity.PackageTypeCard,
		Seller:      "wb",
	}

	return gw.core.ReceiveWarehouses(ctx, descWarehouse)

}

func (gw *grpcGatewayServerImpl) receiveStocksWB(ctx context.Context) error {
	// TODO Перенести в конфиг
	daysToGet := 5 // Количество дней для загрузки
	descStocks := entity.PackageDescription{
		PackageType: entity.PackageTypeStock,
		UpdatedAt:   time.Now(),
		Limit:       daysToGet,
		Seller:      "wb",
	}

	return gw.core.ReceiveStocks(ctx, descStocks)
}

func (gw *grpcGatewayServerImpl) receiveStocksOzon(ctx context.Context) error {

	descOzonStock := entity.PackageDescription{
		PackageType: entity.PackageTypeStock,
		UpdatedAt:   time.Now(),
		Seller:      "ozon",
	}

	return gw.core.ReceiveStocks(ctx, descOzonStock)

}

func (gw *grpcGatewayServerImpl) receiveOrdersWB(ctx context.Context) error {
	// TODO Перенести в конфиг
	daysToGet := 5 // Количество дней для загрузки
	delay := 61    // Количество секунд задержки перед следующим запросом
	descOrderOzon := entity.PackageDescription{
		PackageType: entity.PackageTypeOrder,
		UpdatedAt:   time.Now(),
		Seller:      "ozon",
		Limit:       daysToGet,
		Delay:       delay,
	}

	return gw.core.ReceiveOrders(ctx, descOrderOzon)
}

func (gw *grpcGatewayServerImpl) receiveOrdersOzon(ctx context.Context) error {
	// TODO Перенести в конфиг
	daysToGet := 5 // Количество дней для загрузки
	delay := 61    // Количество секунд задержки перед следующим запросом
	descOrderOzon := entity.PackageDescription{
		PackageType: entity.PackageTypeOrder,
		UpdatedAt:   time.Now(),
		Seller:      "wb",
		Limit:       daysToGet,
		Delay:       delay,
	}

	return gw.core.ReceiveOrders(ctx, descOrderOzon)
}

func (gw *grpcGatewayServerImpl) receiveSalesWB(ctx context.Context) error {
	// TODO Перенести в конфиг
	daysToGet := 5 // Количество дней для загрузки
	delay := 61    // Количество секунд задержки перед следующим запросом

	descDescription := entity.PackageDescription{
		PackageType: entity.PackageTypeSale,
		UpdatedAt:   time.Now(),
		Seller:      "wb",
		Limit:       daysToGet,
		Delay:       delay,
	}

	return gw.core.ReceiveSales(ctx, descDescription)
}
