package seller2cardrepo

import (
	"context"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
)

type seller2cardDB struct {
	ID         int64     `db:"id"`
	ExternalID int64     `db:"external_id"` // Артикул WB
	KTID       int       `db:"int"`         // Идентификатор карточки товара
	NMUUID     string    `db:"nmuuid"`      // Внутренний идентификатор
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	CardID     int64     `db:"card_id"`
	SellerID   int64     `db:"seller_id"`
}

func convertToDBWb2Card(_ context.Context, in entity.Seller2Card) *seller2cardDB {
	return &seller2cardDB{
		ID:         in.ID,
		ExternalID: in.ExternalID,
		KTID:       in.KTID,
		NMUUID:     in.NMUUID,
		CreatedAt:  in.CreatedAt,
		UpdatedAt:  in.UpdatedAt,
		CardID:     in.CardID,
		SellerID:   in.SellerID,
	}
}

func (c seller2cardDB) ConvertToEntityWb2Card(_ context.Context) *entity.Seller2Card {
	return &entity.Seller2Card{
		ID:         c.ID,
		ExternalID: c.ExternalID,
		KTID:       c.KTID,
		NMUUID:     c.NMUUID,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
		CardID:     c.CardID,
		SellerID:   c.SellerID,
	}
}
