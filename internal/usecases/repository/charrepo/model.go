package charrepo

import (
	"context"
	"strings"

	"github.com/efremovich/data-receiver/internal/entity"
)

type characteristicDB struct {
	ID    int64  `db:"id"`
	Title string `db:"title"`
}

type cardCharacteristicDB struct {
	ID               int64  `db:"id"`
	Value            string `db:"value"`
	CharacteristicID int64  `db:"characteristic_id"`
	CardID           int64  `db:"card_id"`
}

func convertToDBCharacteristic(_ context.Context, in entity.Characteristic) *characteristicDB {
	return &characteristicDB{
		ID:    in.ID,
		Title: in.Title,
	}
}

func (c characteristicDB) ConvertToEntityCharacteristic(_ context.Context) *entity.Characteristic {
	return &entity.Characteristic{
		ID:    c.ID,
		Title: c.Title,
		// Value: strings.Split(c.Value, ","),

		// CardID: c.CardID,
	}
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
		CardID:           c.CardID,
	}
}
