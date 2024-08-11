package dimensionrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type dimensionDB struct {
	ID      int64 `db:"id"`
	Length  int   `db:"length"`
	Width   int   `db:"width"`
	Height  int   `db:"height"`
	IsValid bool  `db:"is_valid"`
	CardID  int64 `db:"card_id"`
}

func convertToDBDimension(_ context.Context, in entity.Dimension) *dimensionDB {
	return &dimensionDB{
		ID:      in.ID,
		Length:  in.Length,
		Width:   in.Width,
		Height:  in.Height,
		IsValid: in.IsVaild,
		CardID:  in.CardID,
	}
}

func (c dimensionDB) convertToEntityDimension(_ context.Context) *entity.Dimension {
	return &entity.Dimension{
		ID:      c.ID,
		Length:  c.Length,
		Width:   c.Width,
		Height:  c.Height,
		IsVaild: c.IsValid,
		CardID:  c.CardID,
	}
}
