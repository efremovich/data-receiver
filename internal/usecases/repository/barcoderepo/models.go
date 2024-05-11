package barcoderepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type barcodeDB struct {
	Barcode  string `db:"barcode"`
	SizeID   int64  `db:"size_id"`
	SellerID int64  `db:"seller_id"`
}

func convertToDBBarcode(_ context.Context, in entity.Barcode) *barcodeDB {
	return &barcodeDB{
		Barcode:  in.Barcode,
		SizeID:   in.SizeID,
		SellerID: in.SellerID,
	}
}

func (c barcodeDB) convertToEntityBarcode(_ context.Context) *entity.Barcode {
	return &entity.Barcode{
		Barcode:  c.Barcode,
		SizeID:   c.SizeID,
		SellerID: c.SellerID,
	}
}
