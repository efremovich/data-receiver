package charrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type characteristicDB struct {
	ID    int64  `db:"id"`
	Title string `db:"title"`
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
	}
}

