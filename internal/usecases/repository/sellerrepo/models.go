package sellerrepo

import (
	"context"
	"database/sql"

	"github.com/efremovich/data-receiver/internal/entity"
 	"github.com/efremovich/data-receiver/internal/usecases/repository"
)

type sellerDB struct {
	ID         int64          `db:"id"`
	Title      string         `db:"title"`
	IsEnabled  bool           `db:"is_enabled"`
	ExternalID sql.NullString `db:"external_id"`
}

func convertToDBSeller(_ context.Context, in entity.Seller) *sellerDB {
	return &sellerDB{
		ID:        in.ID,
		Title:     in.Title,
		IsEnabled: in.IsEnabled,
	}
}

func (c sellerDB) convertToEntitySeller(_ context.Context) *entity.Seller {
	return &entity.Seller{
		ID:         c.ID,
		Title:      c.Title,
		IsEnabled:  c.IsEnabled,
		ExternalID: repository.NullStringToString(c.ExternalID),
	}
}
