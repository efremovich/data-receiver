package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) setDimension(ctx context.Context, card *entity.Card) (*entity.Dimension, error) {
	dimension, err := s.dimensionrepo.SelectByCardID(ctx, card.ID)

	if errors.Is(err, ErrObjectNotFound) {
		dimension, err = s.dimensionrepo.Insert(ctx, entity.Dimension{
			Width:   card.Dimension.Width,
			Height:  card.Dimension.Height,
			Length:  card.Dimension.Length,
			IsVaild: card.Dimension.IsVaild,
			CardID:  card.ID,
		})
	}

	if err != nil {
		return nil, wrapErr(fmt.Errorf("Ошибка при получении данных из бд: %w", err))
	}
	return dimension, nil
}
