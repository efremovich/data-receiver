package brandrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type brandDB struct {
	ID       int64  `db:"brand_id"`
	Title    string `db:"title"`
	SellerID int64  `db:"seller_id"`
}

func convertToDBBrand(_ context.Context, in entity.Brand) *brandDB {
	return &brandDB{
		ID:       in.ID,
		Title:    in.Title,
		SellerID: in.SellerID,
	}
}

func (c brandDB) convertToEntityBrand(_ context.Context) *entity.Brand {
	return &entity.Brand{
		ID:       c.ID,
		Title:    c.Title,
		SellerID: c.SellerID,
	}
}
