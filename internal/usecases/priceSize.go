package usecases

import (
	"context"
	"errors"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) setPriceSize(ctx context.Context, in entity.PriceSize) (*entity.PriceSize, error) {
	priceSize, err := s.pricesizerepo.SelectByCardIDAndSizeID(ctx, in.CardID, in.SizeID)
	if errors.Is(err, ErrObjectNotFound) {
		priceSize, err = s.pricesizerepo.Insert(ctx, in)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return priceSize, err
}
