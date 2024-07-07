package usecases

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	aerror "github.com/efremovich/data-receiver/pkg/aerror"
	"github.com/efremovich/data-receiver/pkg/client"
)

func (s *receiverCoreServiceImpl) ReceiveData(ctx context.Context) aerror.AError {
	brokerClient, err := client.NewPackageCreatorBrokerClient(ctx, "nats://localhost:4222")
	if err != nil {
		return aerror.New(ctx, entity.BrokerSendErrorID, err, "Ошибка инициализации клиента к брокеру")
	}

	err = brokerClient.CreatePackageAndSend(ctx, "wb", 0)
	if err != nil {
		return aerror.New(ctx, entity.BrokerSendErrorID, err, "Ошибка публикации сообщения в брокер")
	}
	return nil
}
