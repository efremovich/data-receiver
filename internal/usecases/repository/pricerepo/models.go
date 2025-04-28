package pricerepo

import (
	"context"
	"database/sql"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
)

type priceSizeDB struct {
	ID                   int64        `db:"id"`
	CardID               int64        `db:"card_id"`
	SizeID               int64        `db:"size_id"`
	Price                float64      `db:"price"`
	PriceWithoutDiscount float64      `db:"price_without_discount"`
	PriceFinish          float64      `db:"price_finish"`
	UpdatedAt            sql.NullTime `db:"updated_at"`
}

func convertToDBPrice(_ context.Context, income *entity.PriceSize) *priceSizeDB {
	return &priceSizeDB{
		ID:                   income.ID,
		CardID:               income.CardID,
		SizeID:               income.SizeID,
		Price:                income.Price,
		PriceWithoutDiscount: income.PriceWithoutDiscount,
		PriceFinish:          income.PriceFinish,
		UpdatedAt:            repository.TimeToNullInt(income.UpdatedAt),
	}
}

func (c priceSizeDB) convertToEntityPrice(_ context.Context) *entity.PriceSize {
	return &entity.PriceSize{
		ID:                   c.ID,
		Price:                c.Price,
		PriceWithoutDiscount: c.PriceWithoutDiscount,
		PriceFinish:          c.PriceFinish,
		CardID:               c.CardID,
		SizeID:               c.SizeID,
		UpdatedAt:            c.UpdatedAt.Time,
	}
}
