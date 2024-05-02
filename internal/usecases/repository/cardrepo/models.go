package cardrepo

import (
	"context"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
)

type cardDB struct {
	ID          int64     `db:"id"`
	VendorID    string    `db:"vendor_id"`
	VendorCode  string    `db:"vendor_code"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func convertToDBCard(_ context.Context, in entity.Card) *cardDB {
	return &cardDB{
		ID:          in.ID,
		VendorID:    in.VendorID,
		VendorCode:  in.VendorCode,
		Title:       in.Title,
		Description: in.Description,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}
}

func (c cardDB) ConvertToEntityCard(_ context.Context) *entity.Card {
	return &entity.Card{
		ID:          c.ID,
		VendorID:    c.VendorID,
		VendorCode:  c.VendorCode,
		Title:       c.Title,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
