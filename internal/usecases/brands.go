package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) getBrand(ctx context.Context, brandIn entity.Brand, seller *entity.MarketPlace) (*entity.Brand, error) {
	brand, err := s.brandRepo.SelectByTitleAndSeller(ctx, brandIn.Title, seller.ID)

	if errors.Is(err, ErrObjectNotFound) {
		brand, err = s.brandRepo.Insert(ctx, entity.Brand{
			ExternalID: brandIn.ExternalID,
			Title:      brandIn.Title,
			SellerID:   seller.ID,
		})
	}

	if err != nil {
		return nil, wrapErr(fmt.Errorf("oшибка при получении данных: %w", err))
	}

	return brand, nil
}
