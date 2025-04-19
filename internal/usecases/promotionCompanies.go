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

	groups, gCtx := errgroup.WithContext(ctx)

	for _, c := range clients {
		client := c

		groups.Go(func() error {
			return s.receivePromotionCompanies(gCtx, client, desc)
		})
	}

	// Ждем завершения всех горутин и проверяем наличие ошибок
	if err := groups.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			alogger.WarnFromCtx(ctx, "Операция была отменена: %v", err)
			return nil
		}
		return fmt.Errorf("ошибка при обработке клиентов: %w", err)
	}
	return nil
}

func (s *receiverCoreServiceImpl) receivePromotionCompanies(ctx context.Context, client webapi.ExtAPIFetcher, desc entity.PackageDescription) error {
	promotionMetaList, err := client.GetPromotion(ctx, desc)
	if err != nil {
		return fmt.Errorf("ошибка получение данных о рекламных компаниях из внешнего источника %s, %s", desc.Seller, err.Error())
	}

	seller, err := s.getSeller(ctx, client.GetMarketPlace())
	if err != nil {
		return wrapErr(fmt.Errorf("ошибка получения данных о продавце %s модуль sales:%w", desc.Seller, err))
	}

	for _, meta := range promotionMetaList {
		meta.SellerID = seller.ID

		promotion, err := s.setPromotion(ctx, meta)
		if err != nil {
			return err
		}

		for _, stats := range meta.PromotionStats {
			stats.PromotionID = meta.ID

			// Проверим и создадим связь продавца и товара
			seller2card := entity.Seller2Card{
				CardID:   stats.CardExternalID,
				SellerID: seller.ID,
			}

			card, err := s.setSeller2Card(ctx, seller2card)
			if err != nil {
				return err
			}

			stats.CardID = card.ID

			stats.PromotionID = promotion.ID
			stats.SellerID = seller.ID

			_, err = s.setPromotionStats(ctx, &stats)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *receiverCoreServiceImpl) setPromotion(ctx context.Context, in *entity.Promotion) (*entity.Promotion, error) {
	promotion, err := s.promotionrepo.SelectByExternalID(ctx, in.ExternalID)
	if errors.Is(err, ErrObjectNotFound) {
		promotion, err = s.promotionrepo.Insert(ctx, *in)
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка при вставке данных в таблицу Promotion: %w", err)
	}

	return promotion, nil
}

func (s *receiverCoreServiceImpl) setPromotionStats(ctx context.Context, in *entity.PromotionStats) (*entity.PromotionStats, error) {
	promotionStats, err := s.promotionstatsrepo.SelectByPromotionID(ctx, in.PromotionID, in.CardID, int64(in.AppType))
	if errors.Is(err, ErrObjectNotFound) {
		promotionStats, err = s.promotionstatsrepo.Insert(ctx, *in)
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка при вставке данных в таблицу PromotionStrats: %w", err)
	}

	return promotionStats, nil
}
