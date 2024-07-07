package controller

import (
	"context"

	package_receiver "github.com/efremovich/data-receiver/pkg/data-receiver-service"
)

func (gw *grpcGatewayServerImpl) GetCard(ctx context.Context, in *package_receiver.GetCardRequest) (*package_receiver.GetCardResponse, error) {
	err := gw.core.ReceiveData(ctx)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
