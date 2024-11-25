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
		category, err := s.categoryRepo.SelectByTitle(ctx, cat.Title, seller.ID)
		if errors.Is(err, ErrObjectNotFound) {
			category, err = s.categoryRepo.Insert(ctx, entity.Category{
				Title:      cat.Title,
				ExternalID: cat.ExternalID,
				SellerID:   seller.ID,
			})
		}

		if err != nil {
			return nil, wrapErr(fmt.Errorf("ошибка при получении данных из бд: %w", err))
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (s *receiverCoreServiceImpl) setCardCategories(ctx context.Context, cardID int64, categoryIDs []*entity.Category) error {
	for _, category := range categoryIDs {
		_, err := s.cardcategoryrepo.SelectByCardIDAndCategoryID(ctx, cardID, category.ID)
		if errors.Is(err, ErrObjectNotFound) {
			_, err = s.cardcategoryrepo.Insert(ctx, entity.CardCategory{
				CardID:     cardID,
				CategoryID: category.ID,
			})
		}

		if err != nil {
			return wrapErr(fmt.Errorf("ошибка при получении данных из бд: %w", err))
		}
	}
	return nil
}
