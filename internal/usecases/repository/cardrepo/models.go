package cardrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
)

type cardDB struct {
	ID          int64          `db:"id"`
	VendorID    string         `db:"vendor_id"`
	VendorCode  sql.NullString `db:"vendor_code"`
	Title       string         `db:"title"`
	Description sql.NullString `db:"description"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	BrandID     int64          `db:"brand_id"`
}

func convertToDBCard(_ context.Context, in entity.Card) *cardDB {
	return &cardDB{
		ID:          in.ID,
		VendorID:    in.VendorID,
		VendorCode:  repository.StringToNullString(in.VendorCode),
		Title:       in.Title,
		Description: repository.StringToNullString(in.Description),
		BrandID:     in.Brand.ID,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}
}

func (c cardDB) ConvertToEntityCard(_ context.Context) *entity.Card {
	return &entity.Card{
		ID:          c.ID,
		VendorID:    c.VendorID,
		VendorCode:  repository.NullStringToString(c.VendorCode),
		Title:       c.Title,
		Description: repository.NullStringToString(c.Description),
		Brand:       entity.Brand{ID: c.BrandID},
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
