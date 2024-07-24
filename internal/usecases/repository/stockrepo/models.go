package stockrepo

import (
	"context"
	"database/sql"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
)

type stockDB struct {
	ID               int64        `db:"id"`
	Quantity         int          `db:"quantity"`
	InWayToClient    int          `db:"in_way_to_client"`
	InWayFromClient  int          `db:"in_way_from_client"`
	InWayToWarehouse int          `db:"in_way_to_warehouse"`
	CreatedAt        sql.NullTime `db:"created_at"`
	UpdatedAt        sql.NullTime `db:"updated_at"`
	SizeID           int64        `db:"size_id"`
	BarcodeID        int64        `db:"barcode_id"`
	WarehouseID      int64        `db:"warehouse_id"`
	CardID           int64        `db:"card_id"`
	SellerID         int64        `db:"seller_id"`
}

func convertToDBStock(_ context.Context, in entity.Stock) *stockDB {
	return &stockDB{
		ID:              in.ID,
		Quantity:        in.Quantity,
		InWayToClient:   in.InWayToClient,
		InWayFromClient: in.InWayFromClient,
		CreatedAt:       repository.TimeToNullInt(in.CreatedAt),
		UpdatedAt:       repository.TimeToNullInt(in.CreatedAt),
		SizeID:          in.SizeID,
		BarcodeID:       in.BarcodeID,
		WarehouseID:     in.WarehouseID,
		CardID:          in.CardID,
		SellerID:        in.SellerID,
	}
}

func (c stockDB) convertToEntityStock(_ context.Context) *entity.Stock {
	return &entity.Stock{
		ID:              c.ID,
		Quantity:        c.Quantity,
		InWayToClient:   c.InWayToClient,
		InWayFromClient: c.InWayFromClient,
		CreatedAt:       repository.NullTimeToTime(c.CreatedAt),
		UpdatedAt:       repository.NullTimeToTime(c.UpdatedAt),
		SizeID:          c.SizeID,
		BarcodeID:       c.BarcodeID,
		WarehouseID:     c.WarehouseID,
		CardID:          c.CardID,
		SellerID:        c.SellerID,
	}
}
