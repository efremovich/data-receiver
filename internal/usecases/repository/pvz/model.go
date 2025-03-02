package pvzrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type pvzDB struct {
	ID           int64  `db:"id"`
	OfficeName   string `db:"office_name"`   // Наименование поставщика
	OfficeID     int    `db:"office_id"`     // ID офиса
	SupplierName string `db:"supplier_name"` // Наименование поставщика
	SupplierID   int    `db:"supplier_id"`   // ID поставщика
	SupplierINN  string `db:"supplier_inn"`  // ИНН поставщика
}

func convertToEntityPvz(_ context.Context, c *entity.Pvz) *pvzDB {
	return &pvzDB{
		ID:           c.ID,
		OfficeName:   c.OfficeName,
		OfficeID:     c.OfficeID,
		SupplierName: c.SupplierName,
		SupplierID:   c.SupplierID,
		SupplierINN:  c.SupplierINN,
	}
}
func (c pvzDB) convertToEntityPvz(_ context.Context) *entity.Pvz {
	return &entity.Pvz{
		ID:           c.ID,
		OfficeName:   c.OfficeName,
		OfficeID:     c.OfficeID,
		SupplierName: c.SupplierName,
		SupplierID:   c.SupplierID,
		SupplierINN:  c.SupplierINN,
	}
}
