package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) setCharacterisitc(ctx context.Context, card *entity.Card) ([]*entity.CardCharacteristic, error) {
	cardCharacteristics := []*entity.CardCharacteristic{}

	for _, elem := range card.Characteristics {
		char, err := s.charRepo.SelectByTitle(ctx, elem.Title)
		if errors.Is(err, ErrObjectNotFound) {
			char, err = s.charRepo.Insert(ctx, entity.Characteristic{
				Title: elem.Title,
			})
		}

		if err != nil {
			return nil, wrapErr(fmt.Errorf("ошибка при получении данных: %w", err))
		}

		cardCharacteristic, err := s.cardCharRepo.SelectByCardIDAndCharID(ctx, card.ID, char.ID)

		if errors.Is(err, ErrObjectNotFound) {
			cardCharacteristic, err = s.cardCharRepo.Insert(ctx, entity.CardCharacteristic{
				Value:            elem.Value,
				Title:            elem.Title,
				CharacteristicID: char.ID,
				CardID:           card.ID,
			})
		}

		if err != nil {
			return nil, wrapErr(fmt.Errorf("ошибка при получении данных: %w", err))
		}

		cardCharacteristics = append(cardCharacteristics, cardCharacteristic)
	}

	return cardCharacteristics, nil
}
