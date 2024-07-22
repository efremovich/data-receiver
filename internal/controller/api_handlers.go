package controller

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	package_receiver "github.com/efremovich/data-receiver/pkg/data-receiver-service"
)

func (gw *grpcGatewayServerImpl) ReceiveCard(ctx context.Context, in *package_receiver.ReceiveCardRequest) (*package_receiver.ReceiveCardResponse, error) {
	desc := entity.PackageDescription{
		Cursor: 0,
		Limit:  100,
	}
	err := gw.core.ReceiveCards(ctx, desc)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
