package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) setSize(ctx context.Context, in *entity.Size) (*entity.Size, error) {
	size, err := s.sizerepo.SelectByTitle(ctx, in.Title)

	if errors.Is(err, ErrObjectNotFound) {
		size, err = s.sizerepo.Insert(ctx, *in)
	}

	if err != nil {
		return nil, wrapErr(fmt.Errorf("Ошибка при получении данных: %w", err))
	}

	return size, nil
}

func (s *receiverCoreServiceImpl) getSizeByTitle(ctx context.Context, title string) (*entity.Size, error) {
	size, err := s.sizerepo.SelectByTitle(ctx, title)
	if err != nil {
		return nil, err
	}

	return size, err
}
