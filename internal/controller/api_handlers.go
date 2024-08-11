package controller

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	package_receiver "github.com/efremovich/data-receiver/pkg/data-receiver-service"
)

func (gw *grpcGatewayServerImpl) ReceiveCard(ctx context.Context, in *package_receiver.ReceiveCardRequest) (*package_receiver.ReceiveCardResponse, error) {
	desc := entity.PackageDescription{
		Cursor:      0,
		Limit:       100,
		PackageType: entity.PackageTypeCard,
		Seller:      in.Seller,
		Query: map[string]string{
			"seller": in.Seller,
		},
	}
	err := gw.core.ReceiveCards(ctx, desc)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (gw *grpcGatewayServerImpl) ReceiveWarehouse(ctx context.Context, in *package_receiver.ReceiveWarehouseRequest) (*package_receiver.ReceiveWarehouseResponse, error) {
	err := gw.core.ReceiveWarehouses(ctx)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (gw *grpcGatewayServerImpl) ReceiveStock(ctx context.Context, in *package_receiver.ReceiveStockRequest) (*package_receiver.ReceiveStockResponse, error) {
	desc := entity.PackageDescription{
		PackageType: entity.PackageTypeCard,
		Seller:      in.Seller,
		Query: map[string]string{
			"dateFrom": in.DateFrom,
		},
	}

  err := gw.core.ReceiveStocks(ctx, desc)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
