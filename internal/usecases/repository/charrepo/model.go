package charrepo

import (
	"context"
	"strings"

	"github.com/efremovich/data-receiver/internal/entity"
)

type characteristicDB struct {
	ID    int64  `db:"id"`
	Title string `db:"title"`
	Value string `db:"value"`

	CardID int64 `db:"card_id"`
}

func convertToDBCharacteristic(_ context.Context, in entity.Characteristic) *characteristicDB {
	return &characteristicDB{
		ID:    in.ID,
		Title: in.Title,
		Value: strings.Join(in.Value, ","),

		CardID: in.CardID,
	}
}

func (c characteristicDB) ConvertToEntityCharacteristic(_ context.Context) *entity.Characteristic {
	return &entity.Characteristic{
		ID:    c.ID,
		Title: c.Title,
		Value: strings.Split(c.Value, ","),

		CardID: c.CardID,
	}
}
