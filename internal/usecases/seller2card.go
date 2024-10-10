package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) setSeller2Card(ctx context.Context, cardID, externalID, sellerID int64) (*entity.Seller2Card, error) {
	seller2Card, err := s.seller2cardrepo.SelectByExternalID(ctx, externalID)
	if errors.Is(err, ErrObjectNotFound) {
		seller2Card, err = s.seller2cardrepo.Insert(ctx, entity.Seller2Card{
			ExternalID: externalID,
			CardID:     cardID,
			SellerID:   sellerID,
		})
	}

	if err != nil {
		return nil, wrapErr(fmt.Errorf("ошибка при получении данных: %w", err))
	}

	return seller2Card, nil
}

func (s *receiverCoreServiceImpl) getSeller2Card(ctx context.Context, externalID, sellerID int64) (*entity.Seller2Card, error) {
	seller2Card, err := s.seller2cardrepo.SelectByExternalID(ctx, externalID)
	if err != nil {
		return nil, wrapErr(fmt.Errorf("ошибка при получении данных: %w", err))
	}
	return seller2Card, err
}
