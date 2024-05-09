package pricerepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type priceDB struct {
	ID           int64   `db:"id"`
	Price        float32 `db:"price"`
	Discount     float32 `db:"discount"`
	SpecialPrice float32 `db:"special_price"`
	CardID       int64   `db:"card_id"`
	SellerID     int64   `db:"seller_id"`
}

func convertToDBPrice(_ context.Context, in entity.Price) *priceDB {
	return &priceDB{
		ID:           in.ID,
		Price:        in.Price,
		Discount:     in.Discount,
		SpecialPrice: in.SpecialPrice,
		CardID:       in.CardID,
		SellerID:     in.SellerID,
	}
}

func (c priceDB) convertToEntityPrice(_ context.Context) *entity.Price {
	return &entity.Price{
		ID:           c.ID,
		Price:        c.Price,
		Discount:     c.Discount,
		SpecialPrice: c.SpecialPrice,
		CardID:   c.CardID,
		SellerID: c.SellerID,
	}
}
