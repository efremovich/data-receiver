package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
	"github.com/efremovich/data-receiver/pkg/alogger"
	"golang.org/x/sync/errgroup"
)

func (s *receiverCoreServiceImpl) ReceiveCostFrom1C(ctx context.Context, desc entity.PackageDescription) error {

	clients := s.apiFetcher[desc.Seller]

	g, gCtx := errgroup.WithContext(ctx)

	for _, c := range clients {
		client := c

		g.Go(func() error {
			return s.receiveAndSaveCost(gCtx, client, desc)
		})
	}

	if err := g.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			alogger.WarnFromCtx(ctx, "Операция была отменена: %v", err)
			return nil
		}
		return fmt.Errorf("ошибка при обработке клиентов: %w", err)
	}

	alogger.InfoFromCtx(ctx, "постановка задачи в очередь %d", desc.Limit)

	if desc.Limit > 0 {
		p := entity.PackageDescription{
			PackageType: entity.PackageTypeOrder,

			UpdatedAt: desc.UpdatedAt.Add(-24 * time.Hour),
			Limit:     desc.Limit - 1,
			Seller:    desc.Seller,
			Delay:     desc.Delay,
		}

		err := s.brokerPublisher.SendPackage(ctx, &p)
		if err != nil {
			return fmt.Errorf("ошибка постановки задачи в очередь: %w", err)
		}

		alogger.InfoFromCtx(ctx, "Создана очередь для получения заказов на %s", p.UpdatedAt.Format("02.01.2006"))
	} else {
		alogger.InfoFromCtx(ctx, "Все элементы обработаны")
	}
	return nil
}

func (s *receiverCoreServiceImpl) receiveAndSaveCost(ctx context.Context, client webapi.ExtAPIFetcher, desc entity.PackageDescription) error {
	costs, err := client.GetCosts(ctx, desc)
	if err != nil {
		return fmt.Errorf("ошибка получение данные из внешнего источника %s, %w", desc.Seller, err)
	}

	// Seller
	seller, err := s.getSeller(ctx, client.GetMarketPlace())
	if err != nil {
		return err
	}

	alogger.InfoFromCtx(ctx, "Начали загрузку карточек товара %d от %s", len(costs), seller.Title)

	for _, elem := range costs {
		externalID := strings.TrimSpace(elem.ExternalID)

		card, err := s.cardRepo.SelectByVendorCode(ctx, externalID)
		if err != nil {
			// TODO: обработать ошибку.
			continue
		}

		elem.CardID = card.ID

		err = s.setCost(ctx, elem)
		if err != nil {
			return err
		}
	}
	alogger.InfoFromCtx(ctx, "Задание загрузка карточек товаров успешно завешена количество: %d, маркетплейс %s", len(costs), seller.Title)
	return nil
}

func (s *receiverCoreServiceImpl) setCost(ctx context.Context, in entity.Cost) error {
	_, err := s.costrepo.SelectByCardIDAndDate(ctx, in.CardID, in.CreatedAt)
	if errors.Is(err, ErrObjectNotFound) {
		_, err = s.costrepo.Insert(ctx, in)
		if err != nil {
			return wrapErr(fmt.Errorf("ошибка сохранения записи себестоимости: #w", err))
		}
	}
	return nil
}
