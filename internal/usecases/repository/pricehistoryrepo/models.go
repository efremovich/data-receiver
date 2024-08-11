package pricehistoryrepo

import (
	"context"
	"database/sql"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
)

type priceHistoryDB struct {
	ID          int64        `db:"id"`
	Price       float32      `db:"price"`
	Discount    float32      `db:"discount"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
	PriceSizeID int64        `db:"price_size_id"`
}

func convertToDBPriceHistory(_ context.Context, in entity.PriceHistory) *priceHistoryDB {
	return &priceHistoryDB{
		ID:          in.ID,
		Price:       in.Price,
		Discount:    in.Discount,
		PriceSizeID: in.PriceSizeID,
		UpdatedAt:   repository.TimeToNullInt(in.UpdatedAt),
	}
}

func (c priceHistoryDB) convertToEntityPriceHistory(_ context.Context) *entity.PriceHistory {
	return &entity.PriceHistory{
		ID:          c.ID,
		Price:       c.Price,
		Discount:    c.Discount,
		PriceSizeID: c.PriceSizeID,
		UpdatedAt:   c.UpdatedAt.Time,
	}
}
