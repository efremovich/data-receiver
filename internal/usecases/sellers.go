package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) getSeller(ctx context.Context, marketPlace entity.MarketPlace) (*entity.MarketPlace, error) {
	seller, err := s.sellerRepo.SelectByTitle(ctx, marketPlace.Title)
	if errors.Is(err, ErrObjectNotFound) {
		seller, err = s.sellerRepo.Insert(ctx, marketPlace)
	}

	if err != nil {
		return nil, wrapErr(fmt.Errorf("ошибка при получении данных из бд: %w", err))
	}

	return seller, nil
}
