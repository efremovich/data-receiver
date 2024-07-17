package categoryrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type categoryDB struct {
	ID         int64  `db:"id"`
	Title      string `db:"title"`
	SellerID   int64  `db:"seller_id"`
	ExternalID int64 `db:"external_id"`
}

func convertToDBCategory(_ context.Context, in entity.Category) *categoryDB {
	return &categoryDB{
		ID:         in.ID,
		Title:      in.Title,
		SellerID:   in.SellerID,
		ExternalID: in.ExternalID,
	}
}

func (c categoryDB) convertToEntityCategory(_ context.Context) *entity.Category {
	return &entity.Category{
		ID:         c.ID,
		Title:      c.Title,
		SellerID:   c.SellerID,
		ExternalID: c.ExternalID,
	}
}
