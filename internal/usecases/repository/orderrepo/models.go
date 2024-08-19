package orderrepo

import (
	"context"
	"database/sql"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
)

type orderDB struct {
	ID         int64   `db:"id"`
	ExternalID string  `db:"external_id"`
	Price      float32 `db:"price"`
	StatusID   int64   `db:"status_id"`
	Direction  string  `db:"direction"`
	Type       string  `db:"type"`
	Sale       float32 `db:"sale"`

	Quantity  int          `db:"quantity"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`

	SellerID    int64 `db:"seller_id"`
	CardID      int64 `db:"card_id"`
	WarehouseID int64 `db:"warehouse_id"`
	RegionID    int64 `db:"region_id"`
	PriceSizeID int64 `db:"price_size_id"`
}

func convertToDBOrder(_ context.Context, in entity.Order) *orderDB {
	return &orderDB{
		ID:         in.ID,
		ExternalID: in.ExternalID,
		Price:      in.Price,
		Quantity:   in.Quantity,
		Direction:  in.Direction,
		Type:       in.Type,
		CreatedAt:  repository.TimeToNullInt(in.CreatedAt),
		UpdatedAt:  repository.TimeToNullInt(in.UpdatedAt),

		StatusID:    in.Status.ID,
		SellerID:    in.Seller.ID,
		WarehouseID: in.Warehouse.ID,
		RegionID:    in.Region.ID,
		CardID:      in.Card.ID,
		PriceSizeID: in.PriceSize.ID,
	}
}

func (c orderDB) convertToEntityOrder(_ context.Context) *entity.Order {
	return &entity.Order{
		ID:         c.ID,
		ExternalID: c.ExternalID,
		Price:      c.Price,
		Type:       c.Type,
		Direction:  c.Direction,
		Sale:       c.Sale,
		Quantity:   c.Quantity,
		CreatedAt:  c.CreatedAt.Time,
		UpdatedAt:  c.UpdatedAt.Time,
		Status: &entity.Status{
			ID: c.StatusID,
		},
		Region: &entity.Region{
			ID: c.RegionID,
		},
		Warehouse: &entity.Warehouse{
			ID: c.WarehouseID,
		},
		Seller: &entity.Seller{
			ID: c.SellerID,
		},
		Card: &entity.Card{
			ID: c.CardID,
		},
		PriceSize: &entity.PriceSize{
			ID: c.PriceSizeID,
		},
	}
}
