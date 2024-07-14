package controller

import (
	"context"

	package_receiver "github.com/efremovich/data-receiver/pkg/data-receiver-service"
)

func (gw *grpcGatewayServerImpl) ReceiveCard(ctx context.Context, in *package_receiver.ReceiveCardRequest) (*package_receiver.ReceiveCardResponse, error) {
	err := gw.core.ReceiveCards(ctx, 0)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
