package wb2cardrepo

import (
	"context"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
)

type wb2cardDB struct {
	ID        int64     `db:"id"`
	NMID      int64     `db:"nmid"`   // Артикул WB
	KTID      int       `db:"int"`    // Идентификатор карточки товара
	NMUUID    string    `db:"nmuuid"` // Внутренний идентификатор
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	CardID    int64     `db:"card_id"`
}

func convertToDBWb2Card(_ context.Context, in entity.Wb2Card) *wb2cardDB {
	return &wb2cardDB{
		ID:        in.ID,
		NMID:      in.NMID,
		KTID:      in.KTID,
		NMUUID:    in.NMUUID,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
		CardID:    in.CardID,
	}
}

func (c wb2cardDB) ConvertToEntityWb2Card(_ context.Context) *entity.Wb2Card {
	return &entity.Wb2Card{
		ID:        c.ID,
		NMID:      c.NMID,
		KTID:      c.KTID,
		NMUUID:    c.NMUUID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		CardID:    c.CardID,
	}
}
