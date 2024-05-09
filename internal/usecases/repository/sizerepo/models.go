package sizerepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type sizeDB struct {
	ID       int64  `db:"id"`
	TechSize string `db:"tech_size"`
	Title    string `db:"title"`
	PriceID  int64  `db:"price_id"`
	CardID   int64  `db:"card_id"`
}

func convertToDBSize(_ context.Context, in entity.Size) *sizeDB {
	return &sizeDB{
		ID:       in.ID,
		TechSize: in.TechSize,
		Title:    in.Title,
		PriceID:  in.PriceID,
		CardID:   in.CardID,
	}
}

func (c sizeDB) convertToEntitySize(_ context.Context) *entity.Size {
	return &entity.Size{
		ID:       c.ID,
		TechSize: c.TechSize,
		Title:    c.Title,
		PriceID:  c.PriceID,
		CardID:   c.CardID,
	}
}
