package barcoderepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type barcodeDB struct {
	ID          int64  `db:"id"`
	Barcode     string `db:"barcode"`
	SellerID    int64  `db:"seller_id"`
	PriceSizeID int64  `db:"price_size_id"`
}

func convertToDBBarcode(_ context.Context, in entity.Barcode) *barcodeDB {
	return &barcodeDB{
		ID:          in.ID,
		Barcode:     in.Barcode,
		SellerID:    in.SellerID,
		PriceSizeID: in.PriceSizeID,
	}
}

func (c barcodeDB) convertToEntityBarcode(_ context.Context) *entity.Barcode {
	return &entity.Barcode{
		ID:          c.ID,
		Barcode:     c.Barcode,
		PriceSizeID: c.PriceSizeID,
		SellerID:    c.SellerID,
	}
}
