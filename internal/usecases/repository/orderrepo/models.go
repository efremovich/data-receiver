package orderrepo

import (
	"context"
	"database/sql"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
)

type orderDB struct {
	ID           int64        `db:"id"`
	ExtID        string       `db:"ext_id"`
	Price        float32      `db:"price"`
	Discount     float32      `db:"discount"`
	SpecialPrice float32      `db:"special_price"`
	Status       string       `db:"status"`
	Direction    string       `db:"direction"`
	Type         string       `db:"type"`
	CreatedAt    sql.NullTime `db:"created_at"`
	UpdatedAt    sql.NullTime `db:"updated_at"`

	SellerID    int64 `db:"seller_id"`
	CardID      int64 `db:"card_id"`
	WarehouseID int64 `db:"warehouse_id"`
}

func convertToDBOrder(_ context.Context, in entity.Order) *orderDB {
	return &orderDB{
		ID:           in.ID,
		ExtID:        in.ExtID,
		Price:        in.Price,
		Discount:     in.Discount,
		SpecialPrice: in.SpecialPrice,
		Status:       in.Status,
		Direction:    in.Direction,
		Type:         in.Type,
		CreatedAt:    repository.TimeToNullInt(in.CreatedAt),
		UpdatedAt:    repository.TimeToNullInt(in.UpdatedAt),

		SellerID:    in.SellerID,
		WarehouseID: in.WarehouseID,
		CardID:      in.CardID,
	}
}

func (c orderDB) convertToEntityOrder(_ context.Context) *entity.Order {
	return &entity.Order{
		ID:           c.ID,
		ExtID:        c.ExtID,
		Price:        c.Price,
		Discount:     c.Discount,
		SpecialPrice: c.SpecialPrice,
		Status:       c.Status,
		Type:         c.Type,
		Direction:    c.Direction,
		CreatedAt:    repository.NullTimeToTime(c.CreatedAt),
		UpdatedAt:    repository.NullTimeToTime(c.UpdatedAt),
		WarehouseID:  c.WarehouseID,
		SellerID:     c.SellerID,
		CardID:       c.CardID,
	}
}