package cardcategoryrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type cardCategoryDB struct {
	ID         int64 `db:"id"`
	CardID     int64 `db:"card_id"`
	CategoryID int64 `db:"category_id"`
}

func convertToDB(_ context.Context, in entity.CardCategory) *cardCategoryDB {
	return &cardCategoryDB{
		ID:         in.ID,
		CardID:     in.CardID,
		CategoryID: in.CategoryID,
	}
}

func (c cardCategoryDB) convertToEntityCardCategory(_ context.Context) *entity.CardCategory {
	return &entity.CardCategory{
		ID:         c.ID,
		CardID:     c.CardID,
		CategoryID: c.CategoryID,
	}
}
