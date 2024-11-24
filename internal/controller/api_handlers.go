package controller

import (
	"context"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	package_receiver "github.com/efremovich/data-receiver/pkg/data-receiver-service"
	"github.com/efremovich/data-receiver/pkg/logger"
)

func (gw *grpcGatewayServerImpl) ReceiveCard(ctx context.Context, in *package_receiver.ReceiveCardRequest) (*package_receiver.ReceiveCardResponse, error) {
	desc := entity.PackageDescription{}

	err := gw.core.ReceiveCards(ctx, desc)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (gw *grpcGatewayServerImpl) ReceiveWarehouse(ctx context.Context, in *package_receiver.ReceiveWarehouseRequest) (*package_receiver.ReceiveWarehouseResponse, error) {
	desc := entity.PackageDescription{}
	err := gw.core.ReceiveWarehouses(ctx, desc)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (gw *grpcGatewayServerImpl) ReceiveStock(ctx context.Context, in *package_receiver.ReceiveStockRequest) (*package_receiver.ReceiveStockResponse, error) {
	desc := entity.PackageDescription{}

	err := gw.core.ReceiveStocks(ctx, desc)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (gw *grpcGatewayServerImpl) autoupdate(ctx context.Context, upd time.Duration) {
	t := time.NewTicker(upd)

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			err := gw.update(ctx)
			if err != nil {
				logger.GetLoggerFromContext(ctx).Errorf("Ошибка получение карточек товара %s", err.Error())
			}
		}
	}
}

func (gw *grpcGatewayServerImpl) update(ctx context.Context) error {
	var (
		err  error
		desc entity.PackageDescription
	)

	// date := time.Date(2024, 7, 31, 0, 0, 0, 0, time.UTC)
	// date := time.Now()
	// daysToGet := 3
	// delay := 61
	desc = entity.PackageDescription{
		Limit:       100,
		Cursor:      "0",
		PackageType: entity.PackageTypeCard,
		Seller:      "ozon",
	}

	err = gw.core.ReceiveCards(ctx, desc)
	if err != nil {
		return err
	}

	// desc = entity.PackageDescription{
	// 	Limit:       100,
	// 	Cursor:      "0",
	// 	PackageType: entity.PackageTypeCard,
	// 	Seller:      "ozon",
	// }

	// err = gw.core.ReceiveWarehouses(ctx, desc)
	// if err != nil {
	// 	return err
	// }

	desc = entity.PackageDescription{
		PackageType: entity.PackageTypeStock,
		UpdatedAt:   time.Now(),
		Seller:      "ozon",
	}

	err = gw.core.ReceiveStocks(ctx, desc)
	if err != nil {
		return err
	}

	// desc = entity.PackageDescription{
	// 	PackageType: entity.PackageTypeOrder,
	// 	UpdatedAt:   date,
	// 	Seller:      "wb",
	// 	Limit:       daysToGet,
	// 	Delay:       61,
	// }

	// err = gw.core.ReceiveOrders(ctx, desc)
	// if err != nil {
	// 	return err
	// }

	// desc = entity.PackageDescription{
	// 	PackageType: entity.PackageTypeSale,
	// 	UpdatedAt:   date,
	// 	Seller:      "wb",
	// 	Limit:       daysToGet,
	// 	Delay:       delay,
	// }

	// err = gw.core.ReceiveSales(ctx, desc)
	// if err != nil {
	// 	return err
	// }

	return nil
}
