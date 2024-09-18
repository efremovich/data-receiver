package brokerconsumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/alogger"
	anats "github.com/efremovich/data-receiver/pkg/anats"
)

type BrokerConsumer interface {
	// Проверка работы брокера
	Ping() error
	// Подписка на сообщения пакета и отправка
	SubcriveToGetCards(ctx context.Context, h handlerForCreatePackageAndSend, opt anats.SubscribeOptions) error
	// Публикация сообщения
	PublishMessageToGetCard(ctx context.Context, seller string, cursor int) error
}

// Имплементация брокера потребителя.
type brokerConsumerImpl struct {
	nats anats.NatsClient
}

// Инициализация брокера потребителя.
func NewBrokerConsumer(ctx context.Context, urls []string, updateStream bool) (BrokerConsumer, error) {
	cfg := anats.NatsClientConfig{
		Urls:               urls,
		StreamName:         PackageCreatorStreamName,
		Subjects:           []string{SubjectForGetCards},
		CreateUpdateStream: updateStream,
	}

	cl, err := anats.NewNatsClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &brokerConsumerImpl{nats: cl}, nil
}

// Подписаться на создания пакетов и их отправки.
func (b *brokerConsumerImpl) SubcriveToGetCards(ctx context.Context, h handlerForCreatePackageAndSend, opt anats.SubscribeOptions) error {
	return b.nats.Subscribe(ctx, CardCreatorConsumer, SubjectForGetCards,
		func(ctx context.Context, msg anats.Message) anats.MessageResultEnum {
			var desc entity.PackageDescription

			err := json.Unmarshal(msg.GetData(), &desc)
			if err != nil {
				msg := fmt.Sprintf("ошибка десериализации сообщения из брокера %s: %s", SubjectForGetCards, err.Error())
				alogger.ErrorFromCtx(ctx, msg, err, nil, false)

				return anats.MessageResultEnumFatalError
			}

			return h(ctx, desc, msg.GetRetryCount(), msg.IsLastAttempt())
		}, opt)
}

func (b *brokerConsumerImpl) PublishMessageToGetCard(ctx context.Context, seller string, cursor int) error {
	task := Task{
		Seller: seller,
		Cursor: cursor,
	}

	msgBytes, err := json.Marshal(&task)
	if err != nil {
		return err
	}

	return b.nats.PublishMessageDupe(ctx, SubjectForGetCards, msgBytes, "card")
}

func (b *brokerConsumerImpl) Ping() error {
	return b.nats.Ping()
}
