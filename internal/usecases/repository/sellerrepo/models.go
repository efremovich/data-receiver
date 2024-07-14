package sellerrepo

import (
	"context"
	"database/sql"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
)

type sellerDB struct {
	ID       int64          `db:"id"`
	Title    string         `db:"title"`
	IsEnable bool           `db:"is_enable"`
	ExtID    sql.NullString `db:"ext_id"`
}

func convertToDBSeller(_ context.Context, in entity.Seller) *sellerDB {
	return &sellerDB{
		ID:       in.ID,
		Title:    in.Title,
		IsEnable: in.IsEnable,
	}
}

func (c sellerDB) convertToEntitySeller(_ context.Context) *entity.Seller {
	return &entity.Seller{
		ID:       c.ID,
		Title:    c.Title,
		IsEnable: c.IsEnable,
		ExtID:    repository.NullStringToString(c.ExtID),
	}
}
