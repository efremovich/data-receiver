package controller

import (
	"context"
	"errors"
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
func (gw *grpcGatewayServerImpl) handlerForCreateCard(ctx context.Context, desc entity.PackageDescription, _ int, _ bool) anats.MessageResultEnum {
	start := time.Now()

	switch desc.PackageType {
	case entity.PackageTypeCard:
		err := gw.core.ReceiveCards(ctx, desc)
		if err != nil {
			alogger.ErrorFromCtx(ctx, "ошибка обработки пакета %s: %s", desc.Cursor, err.Error())

			if errors.Is(err, entity.ErrPermanent) {
				return anats.MessageResultEnumFatalError
			}

			return anats.MessageResultEnumTempError
		}
	case entity.PackageTypeStock:
		err := gw.core.ReceiveStocks(ctx, desc)
		if err != nil {
			alogger.ErrorFromCtx(ctx, "ошибка обработки пакета %s: %s", desc.Cursor, err.Error())

			if errors.Is(err, entity.ErrPermanent) {
				return anats.MessageResultEnumFatalError
			}

			return anats.MessageResultEnumTempError
		}

	case entity.PackageTypeSale:
		err := gw.core.ReceiveSales(ctx, desc)
		if err != nil {
			alogger.ErrorFromCtx(ctx, "ошибка обработки пакета %s: %s", desc.Cursor, err.Error())

			if errors.Is(err, entity.ErrPermanent) {
				return anats.MessageResultEnumFatalError
			}

			return anats.MessageResultEnumTempError
		}
	case entity.PackageTypeOrder:
		err := gw.core.ReceiveOrders(ctx, desc)
		if err != nil {
			alogger.ErrorFromCtx(ctx, "ошибка обработки пакета %s: %s", desc.Cursor, err.Error())

			if errors.Is(err, entity.ErrPermanent) {
				return anats.MessageResultEnumFatalError
			}

			return anats.MessageResultEnumTempError
		}
	case entity.PackageTypeSaleReports:
		err := gw.core.ReceiveSaleReport(ctx, desc)
		if err != nil {
			alogger.ErrorFromCtx(ctx, "ошибка обработки пакета %s: %s", desc.Cursor, err.Error())

			if errors.Is(err, entity.ErrPermanent) {
				return anats.MessageResultEnumFatalError
			}

			return anats.MessageResultEnumTempError
		}
	}

	if desc.Delay != 0 {
		elapsed := time.Since(start).Seconds()
		remainingDelay := float64(desc.Delay) - elapsed
		time.Sleep(time.Duration(remainingDelay * float64(time.Second)))
	}

	alogger.InfoFromCtx(ctx, "окончание обработки создания пакета время - %.3fs", time.Since(start).Seconds())

	return anats.MessageResultEnumSuccess
}
