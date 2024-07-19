package categoryrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type categoryDB struct {
	ID         int64  `db:"id"`
	Title      string `db:"title"`
	SellerID   int64  `db:"seller_id"`
	CardID     int64  `db:"card_id"`
	ExternalID int64  `db:"external_id"`
	ParentID   int64  `db:"parent_id"`
}

func convertToDBCategory(_ context.Context, in entity.Category) *categoryDB {
	return &categoryDB{
		ID:         in.ID,
		Title:      in.Title,
		SellerID:   in.SellerID,
		CardID:     in.CardID,
		ExternalID: in.ExternalID,
		ParentID:   in.ParentID,
	}
}

func (c categoryDB) convertToEntityCategory(_ context.Context) *entity.Category {
	return &entity.Category{
		ID:         c.ID,
		Title:      c.Title,
		SellerID:   c.SellerID,
		CardID:     c.CardID,
		ExternalID: c.ExternalID,
		ParentID:   c.ParentID,
	}
}
