package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) setSizes(ctx context.Context, card *entity.Card) ([]*entity.Size, error) {
	sizes := []*entity.Size{}
	for _, elem := range card.Sizes {
		size, err := s.sizerepo.SelectByTitle(ctx, elem.Title)

		if errors.Is(err, ErrObjectNotFound) {
			size, err = s.sizerepo.Insert(ctx, entity.Size{
				ExternalID: elem.ExternalID,
				TechSize:   elem.TechSize,
				Title:      elem.Title,
			})
		}

		if err != nil {
			return nil, wrapErr(fmt.Errorf("Ошибка при получении данных: %w", err))
		}
		sizes = append(sizes, size)
	}

	return sizes, nil
}

func (s *receiverCoreServiceImpl) getSizeByTitle(ctx context.Context, title string) (*entity.Size, error) {
	size, err := s.sizerepo.SelectByTitle(ctx, title)
	if err != nil {
		return nil, err
	}

	return size, err
}
