package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) setPriceSize(ctx context.Context, income entity.PriceSize) (*entity.PriceSize, error) {
	priceSize, err := s.pricesizerepo.SelectByCardIDAndSizeID(ctx, income.CardID, income.SizeID)
	if priceSize.Price != income.Price ||
		priceSize.Discount != income.Discount ||
		priceSize.SpecialPrice != income.SpecialPrice {
		err = s.pricesizerepo.UpdateExecOne(ctx, income)
		if err != nil {
			return nil, fmt.Errorf("ошибка вставки данных о ценовой базе модуль pricesizes:%w", err)
		}
	}

	if errors.Is(err, ErrObjectNotFound) {
		priceSize, err = s.pricesizerepo.Insert(ctx, income)
		if err != nil {
			return nil, fmt.Errorf("ошибка вставки данных о ценовой базе модуль pricesizes:%w", err)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных о ценовой базе модуль pricesizes:%w", err)
	}

	return priceSize, nil
}
