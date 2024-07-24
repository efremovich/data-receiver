package warehouserepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type warehouseDB struct {
	ID         int64  `db:"id"`
	ExternalID int64  `db:"external_id"`
	Title      string `db:"name"`
	Address    string `db:"address"`
	SellerID   int64  `db:"seller_id"`
	TypeID     int64  `db:"warehouse_type_id"`
}

func convertToDBWarehouse(_ context.Context, in entity.Warehouse) *warehouseDB {
	return &warehouseDB{
		ID:         in.ID,
		ExternalID: in.ExternalID,
		Title:      in.Title,
		Address:    in.Address,
		TypeID:     in.TypeID,
		SellerID:   in.SellerID,
	}
}

func (c warehouseDB) convertToEntityWarehouse(_ context.Context) *entity.Warehouse {
	return &entity.Warehouse{
		ID:         c.ID,
		ExternalID: c.ExternalID,
		Title:      c.Title,
		Address:    c.Address,
		TypeID:     c.TypeID,
		SellerID:   c.SellerID,
	}
}
