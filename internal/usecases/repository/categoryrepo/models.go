package categoryrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type categoryDB struct {
	ID       int64  `db:"id"`
	Title    string `db:"title"`
	CardID   int64  `db:"card_id"`
	SellerID int64  `db:"seller_id"`
}

func convertToDBCategory(_ context.Context, in entity.Category) *categoryDB {
	return &categoryDB{
		ID:       in.ID,
		Title:    in.Title,
		CardID:   in.CardID,
		SellerID: in.SellerID,
	}
}

func (c categoryDB) convertToEntityCategory(_ context.Context) *entity.Category {
	return &entity.Category{
		ID:       c.ID,
		Title:    c.Title,
		CardID:   c.CardID,
		SellerID: c.SellerID,
	}
}
