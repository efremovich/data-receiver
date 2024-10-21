package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) setCategory(ctx context.Context, card *entity.Card, seller *entity.Seller) ([]*entity.Category, error) {
	categories := []*entity.Category{}

	for _, cat := range card.Categories {
		category, err := s.categoryRepo.SelectByTitle(ctx, cat.Title)
		if errors.Is(err, ErrObjectNotFound) {
			category, err = s.categoryRepo.Insert(ctx, entity.Category{
				Title:      cat.Title,
				ExternalID: cat.ExternalID,
				CardID:     card.ID,
				SellerID:   seller.ID,
				ParentID:   0,
			})
		}

		if err != nil {
			return nil, wrapErr(fmt.Errorf("ошибка при получении данных из бд: %w", err))
		}

		categories = append(categories, category)
	}

	return categories, nil
}
