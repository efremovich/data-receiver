package charrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type characteristicDB struct {
	ID    int64
	Title string
	Value []string

	CardID int64
}

func convertToDBCharacteristic(_ context.Context, in entity.Characteristic) *characteristicDB {
	return &characteristicDB{
		ID:    in.ID,
		Title: in.Title,
		Value: in.Value,

		CardID: in.CardID,
	}
}

func (c characteristicDB) ConvertToEntityCharacteristic(_ context.Context) *entity.Characteristic {
	return &entity.Characteristic{
		ID:    c.ID,
		Title: c.Title,
		Value: c.Value,

		CardID: c.CardID,
	}
}
