package warehousetyperepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type warehousetypeDB struct {
	ID    int64  `db:"id"`
	Title string `db:"name"`
}

func convertToDBWarehousetype(_ context.Context, in entity.WarehouseType) *warehousetypeDB {
	return &warehousetypeDB{
		ID:    in.ID,
		Title: in.Title,
	}
}

func (c warehousetypeDB) convertToEntityWarehousetype(_ context.Context) *entity.WarehouseType {
	return &entity.WarehouseType{
		ID:    c.ID,
		Title: c.Title,
	}
}
