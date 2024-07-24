package cardcharrepo

import (
	"context"
	"strings"

	"github.com/efremovich/data-receiver/internal/entity"
)

type cardCharacteristicDB struct {
	ID               int64  `db:"id"`
	Value            string `db:"value"`
	CharacteristicID int64  `db:"characteristic_id"`
	CardID           int64  `db:"card_id"`
	Title            string `db:"title"`
}

func convertToDBCardCharacteristic(_ context.Context, in entity.CardCharacteristic) *cardCharacteristicDB {
	return &cardCharacteristicDB{
		ID:               in.ID,
		Value:            strings.Join(in.Value, ","),
		CharacteristicID: in.CharacteristicID,
		CardID:           in.CardID,
	}
}

func (c cardCharacteristicDB) ConvertToEntityCardCharacteristic(_ context.Context) *entity.CardCharacteristic {
	return &entity.CardCharacteristic{
		ID:               c.ID,
		Value:            strings.Split(c.Value, ","),
		CharacteristicID: c.CharacteristicID,
		Title:            c.Title,
		CardID:           c.CardID,
	}
}
