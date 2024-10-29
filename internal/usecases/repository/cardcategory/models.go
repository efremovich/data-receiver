package cardcategory

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type cardCategoryDB struct {
	ID         int64
	CardID     int64
	CategoryID int64
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
