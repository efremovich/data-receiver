package usecases

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) getPriceSize(ctx context.Context, cardID, sizeID int64) (*entity.PriceSize, error) {
	priceSize, err := s.pricesizerepo.SelectByCardIDAndSizeID(ctx, cardID, sizeID)
	// if err != nil && !errors.Is(err, sql.ErrNoRows) {
	if err != nil {
		// TODO Обратока ошибок во всех функциях
		return nil, err
	}

	return priceSize, err
}

func (s *receiverCoreServiceImpl) setPriceSize(ctx context.Context, in entity.PriceSize) (*entity.PriceSize, error) {
	priceSize, err := s.pricesizerepo.Insert(ctx, in)
	if err != nil {
		// TODO Обратока ошибок во всех функциях
		return nil, err
	}

	return priceSize, err
}
