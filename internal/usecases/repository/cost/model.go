package costrepo

import (
	"context"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
)

type costsDB struct {
	ID        int64     `db:"id"`
	CardID    int64     `db:"card_id"`
	Amount    float64   `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func convertDBToCost(_ context.Context, in *entity.Cost) *costsDB {
	return &costsDB{
		CardID:    in.CardID,
		Amount:    in.Amount,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}

func (s costsDB) convertToEntityCost(_ context.Context) *entity.Cost {
	return &entity.Cost{
		ID:        s.ID,
		CardID:    s.CardID,
		Amount:    s.Amount,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
