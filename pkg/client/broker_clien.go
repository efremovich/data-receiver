package client

import (
	"context"

	"github.com/efremovich/data-receiver/pkg/broker/brokerconsumer"
)

// Broker клиент к сервису package-creator.
type PackageCreatorBrokerClient interface {
	// Создание пакета и отправка его в package-sender.
	CreatePackageAndSend(ctx context.Context, seller string, cursor int) error
}

type packageCreatorBrokerClientImpl struct {
	cl brokerconsumer.BrokerConsumer
}

func NewPackageCreatorBrokerClient(ctx context.Context, url string) (PackageCreatorBrokerClient, error) {
	cl, err := brokerconsumer.NewBrokerConsumer(ctx, []string{url}, false)
	if err != nil {
		return nil, err
	}

	return &packageCreatorBrokerClientImpl{
		cl: cl,
	}, nil
}

func (c *packageCreatorBrokerClientImpl) CreatePackageAndSend(ctx context.Context, seller string, cursor int) error {
	err := c.cl.PublishMessageToGetCard(ctx, seller, cursor)
	if err != nil {
		return err
	}

	return nil
}
