package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
	"github.com/efremovich/data-receiver/pkg/alogger"
	"golang.org/x/sync/errgroup"
)

func (s *receiverCoreServiceImpl) ReceivePromotionCompanies(ctx context.Context, desc entity.PackageDescription) error {

	clients := s.apiFetcher[desc.Seller]

	g, gCtx := errgroup.WithContext(ctx)

	for _, c := range clients {
		client := c

		g.Go(func() error {
			return s.receivePromotionCompanies(gCtx, client, desc)
		})
	}

	// Ждем завершения всех горутин и проверяем наличие ошибок
	if err := g.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			alogger.WarnFromCtx(ctx, "Операция была отменена: %v", err)
			return nil
		}
		return fmt.Errorf("ошибка при обработке клиентов: %w", err)
	}
	return nil
}

func (s *receiverCoreServiceImpl) receivePromotionCompanies(ctx context.Context, client webapi.ExtAPIFetcher, desc entity.PackageDescription) error {
	_, err := client.GetPromotion(ctx, desc)
	if err != nil {
		return fmt.Errorf("ошибка получение данных о рекламных компаниях из внешнего источника %s, %s", desc.Seller, err.Error())
	}
	return nil
}
