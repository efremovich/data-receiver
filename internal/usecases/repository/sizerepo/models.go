package sizerepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type sizeDB struct {
	ID       int64  `db:"id"`
	TechSize string `db:"tech_size"`
	Title    string `db:"name"`
}

func convertToDBSize(_ context.Context, in entity.Size) *sizeDB {
	return &sizeDB{
		ID:       in.ID,
		TechSize: in.TechSize,
		Title:    in.Title,
	}
}

func (c sizeDB) convertToEntitySize(_ context.Context) *entity.Size {
	return &entity.Size{
		ID:       c.ID,
		TechSize: c.TechSize,
		Title:    c.Title,
	}
}
