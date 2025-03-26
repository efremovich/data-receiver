package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) setPvz(ctx context.Context, in *entity.Pvz) (*entity.Pvz, error) {
	pvz, err := s.pvzrepo.SelectByOfficeID(ctx, in.OfficeID)
	if errors.Is(err, ErrObjectNotFound) {
		pvz, err = s.pvzrepo.Insert(ctx, *in)
	}

	if err != nil {
		return nil, wrapErr(fmt.Errorf("ошибка при получении данных: %w", err))
	}

	return pvz, nil
}
