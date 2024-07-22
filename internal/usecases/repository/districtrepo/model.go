package districtrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type districtDB struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func convertToDBDistrict(_ context.Context, in entity.District) *districtDB {
	return &districtDB{
		ID:   in.ID,
		Name: in.Name,
	}
}

func (c districtDB) ConvertToEntityDistrict(_ context.Context) *entity.District {
	return &entity.District{
		ID:   c.ID,
		Name: c.Name,
	}
}
