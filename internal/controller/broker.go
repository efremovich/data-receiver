package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/alogger"
	anats "github.com/efremovich/data-receiver/pkg/anats"
)

// Создание подписчиков на очереди.
func (gw *grpcGatewayServerImpl) makeBrokerSubscribers(ctx context.Context) error {
	// Очередь создания пакета и отправки его в package-sender.
	optCreatePackageAndSend := anats.SubscribeOptions{
		Workers:           gw.cfg.Queue.Workers,
		MaxDeliver:        gw.cfg.Queue.MaxDeliver,
		AckWaitSeconds:    gw.cfg.Queue.AckWaitSeconds,
		NakTimeoutSeconds: gw.cfg.Queue.NakTimeoutSeconds,
		MaxAckPending:     gw.cfg.Queue.MaxAckPending,
	}

	err := gw.brokerConsumer.SubcriveToGetCards(ctx, gw.handlerForCreateCard, optCreatePackageAndSend)
	if err != nil {
		return fmt.Errorf("ошибка создания подписки для очереди приёма документов - %s", err.Error())
	}

	return nil
}

// Обработчик события создания пакета и отправки в сервис package-sender.
func (gw *grpcGatewayServerImpl) handlerForCreateCard(ctx context.Context, desc entity.PackageDescription, retry int, isLastRetry bool) anats.MessageResultEnum {
	start := time.Now()
	alogger.DebugFromCtx(ctx, fmt.Sprintf("начало обработки сообщения %d", desc.Cursor), nil, nil, false)
	switch desc.PackageType {
	case entity.PackageTypeCard:

		aerr := gw.core.ReceiveCards(ctx, desc)
		if aerr != nil {
			alogger.ErrorFromCtx(ctx, fmt.Sprintf("ошибка обработки пакета %d: %s", desc.Cursor, aerr.DeveloperMessage()), aerr, nil, false)
			if aerr.IsCritical() {
				return anats.MessageResultEnumFatalError
			}

			return anats.MessageResultEnumTempError
		}
	case entity.PackageTypeStock:

		aerr := gw.core.ReceiveStocks(ctx, desc)
		if aerr != nil {
			alogger.ErrorFromCtx(ctx, fmt.Sprintf("ошибка обработки пакета %d: %s", desc.Cursor, aerr.DeveloperMessage()), aerr, nil, false)
			if aerr.IsCritical() {
				return anats.MessageResultEnumFatalError
			}

			return anats.MessageResultEnumTempError
		}
	}

	alogger.InfoFromCtx(ctx, fmt.Sprintf("окончание обработки создания пакета %d. время - %.3fs", desc.Cursor, time.Since(start).Seconds()), nil, nil, false)

	return anats.MessageResultEnumSuccess
}
