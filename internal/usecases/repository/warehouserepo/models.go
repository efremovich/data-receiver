package warehouserepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type warehouseDB struct {
	ID       int64  `db:"id"`
	ExternalID    string `db:"external_id"`
	Title    string `db:"title"`
	Address  string `db:"address"`
	Type     string `db:"type"`
	SellerID int64  `db:"seller_id"`
}

func convertToDBWarehouse(_ context.Context, in entity.Warehouse) *warehouseDB {
	return &warehouseDB{
		ID:       in.ID,
		ExternalID:    in.ExternalID,
		Title:    in.Title,
		Address:  in.Address,
		Type:     in.Type,
		SellerID: in.SellerID,
	}
}

func (c warehouseDB) convertToEntityWarehouse(_ context.Context) *entity.Warehouse {
	return &entity.Warehouse{
		ID:       c.ID,
		ExternalID:    c.ExternalID,
		Title:    c.Title,
		Address:  c.Address,
		Type:     c.Type,
		SellerID: c.SellerID,
	}
}
