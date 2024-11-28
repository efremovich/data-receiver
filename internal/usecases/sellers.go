package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) getSeller(ctx context.Context, sellerTitle string) (*entity.Seller, error) {
	seller, err := s.sellerRepo.SelectByTitle(ctx, sellerTitle)
	if errors.Is(err, ErrObjectNotFound) {
		seller, err = s.sellerRepo.Insert(ctx, entity.Seller{
			Title:     sellerTitle,
			IsEnabled: true,
		})
	}

	if err != nil {
		return nil, wrapErr(fmt.Errorf("ошибка при получении данных из бд: %w", err))
	}

	return seller, nil
}
