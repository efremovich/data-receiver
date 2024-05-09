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
	BrandID     int            `db:"brand_id"`
}


type categoryDB struct {
	ID    int64  `db:"id"`
	Title string `db:"title"`

	CardID   int64 `db:"card_id"`
	SellerID int64 `db:"seller_id"`
}

type BarcodeDB struct {
	Barcode string `db:"barcode"`

	SizeID   int64 `db:"size_id"`
	SellerID int64 `db:"seller_id"`
}

func convertToBarcodeDB(_ context.Context, in entity.Barcode) *BarcodeDB {
	return &BarcodeDB{
		Barcode: in.Barcode,

		SizeID:   in.SizeID,
		SellerID: in.SellerID,
	}
}

func (c BarcodeDB) ConvertToEntityBarcode(_ context.Context) *entity.Barcode {
	return &entity.Barcode{
		Barcode: c.Barcode,

		SizeID:   c.SizeID,
		SellerID: c.SellerID,
	}
}

func convertToDBCard(_ context.Context, in entity.Card) *cardDB {
	return &cardDB{
		ID:          in.ID,
		VendorID:    in.VendorID,
		VendorCode:  repository.StringToNullString(in.VendorCode),
		Title:       in.Title,
		Description: repository.StringToNullString(in.Description),
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
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}


func convertToDBCategories(_ context.Context, in entity.Category) *categoryDB {
	return &categoryDB{
		ID:    in.ID,
		Title: in.Title,

		CardID:   in.CardID,
		SellerID: in.SellerID,
	}
}

func (c categoryDB) ConvertToEntityCategory(_ context.Context) *entity.Category {
	return &entity.Category{
		ID:    c.ID,
		Title: c.Title,

		CardID:   c.CardID,
		SellerID: c.SellerID,
	}
}
