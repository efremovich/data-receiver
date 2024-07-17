package barcoderepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type barcodeDB struct {
	ID       int64  `db:"id"`
	Barcode  string `db:"barcode"`
	SizeID   int64  `db:"size_id"`
	SellerID int64  `db:"seller_id"`
}

func convertToDBBarcode(_ context.Context, in entity.Barcode) *barcodeDB {
	return &barcodeDB{
		ID:       in.ID,
		Barcode:  in.Barcode,
		SizeID:   in.SizeID,
		SellerID: in.SellerID,
	}
}

func (c barcodeDB) convertToEntityBarcode(_ context.Context) *entity.Barcode {
	return &entity.Barcode{
		ID:       c.ID,
		Barcode:  c.Barcode,
		SizeID:   c.SizeID,
		SellerID: c.SellerID,
	}
}
