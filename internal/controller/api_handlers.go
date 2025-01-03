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
				logger.GetLoggerFromContext(ctx).Errorf("ошибка получение карточек товара %s", err.Error())
			}
		}
	}
}

func (gw *grpcGatewayServerImpl) update(ctx context.Context) error {
	var err error
	// desc entity.PackageDescription

	// date := time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC)
	date := time.Now()
	daysToGet := 15
	// d2tToGet := 90
	delay := 61
	descOzon := entity.PackageDescription{
		Limit:       100,
		Cursor:      "0",
		PackageType: entity.PackageTypeCard,
		Seller:      "ozon",
	}

	err = gw.core.ReceiveCards(ctx, descOzon)
	if err != nil {
		return err
	}

	descCard := entity.PackageDescription{
		Limit:       100,
		Cursor:      "0",
		PackageType: entity.PackageTypeCard,
		Seller:      "wb",
	}

	err = gw.core.ReceiveCards(ctx, descCard)
	if err != nil {
		return err
	}

	descWarehouse := entity.PackageDescription{
		Limit:       100,
		Cursor:      "0",
		PackageType: entity.PackageTypeCard,
		Seller:      "wb",
	}

	err = gw.core.ReceiveWarehouses(ctx, descWarehouse)
	if err != nil {
		return err
	}

	descOzonStock := entity.PackageDescription{
		PackageType: entity.PackageTypeStock,
		UpdatedAt:   time.Now(),
		Seller:      "ozon",
	}

	err = gw.core.ReceiveStocks(ctx, descOzonStock)
	if err != nil {
		return err
	}

	descStocks := entity.PackageDescription{
		PackageType: entity.PackageTypeStock,
		UpdatedAt:   time.Now(),
		Limit:       daysToGet,
		Seller:      "wb",
	}

	err = gw.core.ReceiveStocks(ctx, descStocks)
	if err != nil {
		return err
	}

	descOrderWb := entity.PackageDescription{
		PackageType: entity.PackageTypeOrder,
		UpdatedAt:   date,
		Seller:      "ozon",
		Limit:       daysToGet,
		Delay:       delay,
	}

	err = gw.core.ReceiveOrders(ctx, descOrderWb)
	if err != nil {
		return err
	}

	descOrderOzon := entity.PackageDescription{
		PackageType: entity.PackageTypeOrder,
		UpdatedAt:   date,
		Seller:      "ozon",
		Limit:       daysToGet,
		Delay:       delay,
	}

	err = gw.core.ReceiveOrders(ctx, descOrderOzon)
	if err != nil {
		return err
	}
	descDescription := entity.PackageDescription{
		PackageType: entity.PackageTypeSale,
		UpdatedAt:   date,
		Seller:      "wb",
		Limit:       daysToGet,
		Delay:       delay,
	}

	err = gw.core.ReceiveSales(ctx, descDescription)
	if err != nil {
		return err
	}

	return nil
}
