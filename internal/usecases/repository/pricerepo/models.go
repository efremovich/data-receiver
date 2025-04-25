package pricerepo

import (
	"context"
	"database/sql"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
)

type priceSizeDB struct {
	ID           int64        `db:"id"`
	CardID       int64        `db:"card_id"`
	SizeID       int64        `db:"size_id"`
	Price        float32      `db:"price"`
	Discount     float32      `db:"discount"`
	SpecialPrice float32      `db:"special_price"`
	UpdatedAt    sql.NullTime `db:"updated_at"`
}

func convertToDBPrice(_ context.Context, income *entity.PriceSize) *priceSizeDB {
	return &priceSizeDB{
		ID:           income.ID,
		CardID:       income.CardID,
		SizeID:       income.SizeID,
		Price:        income.Price,
		Discount:     income.Discount,
		SpecialPrice: income.SpecialPrice,
		UpdatedAt:    repository.TimeToNullInt(income.UpdatedAt),
	}
}

func (c priceSizeDB) convertToEntityPrice(_ context.Context) *entity.PriceSize {
	return &entity.PriceSize{
		ID:           c.ID,
		Price:        c.Price,
		Discount:     c.Discount,
		SpecialPrice: c.SpecialPrice,
		CardID:       c.CardID,
		SizeID:       c.SizeID,
		UpdatedAt:    c.UpdatedAt.Time,
	}
}
