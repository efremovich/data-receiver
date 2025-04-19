package promotionrepo

import (
	"context"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
)

type promotionDB struct {
	ID          int64     `db:"id"`
	ExternalID  int64     `db:"external_id"`
	Name        string    `db:"name"`
	Type        int       `db:"type"`
	Status      int       `db:"status"`
	ChangeTime  time.Time `db:"change_time"`
	CreateTime  time.Time `db:"create_time"`
	DateStart   time.Time `db:"date_start"`
	DateEnd     time.Time `db:"date_end"`
	Views       int       `db:"views"`
	Clicks      int       `db:"clicks"`
	CTR         float64   `db:"ctr"`
	CPC         float64   `db:"cpc"`
	Spent       float64   `db:"spent"`
	Orders      int       `db:"orders"`
	CR          float64   `db:"cr"`
	SHKs        int       `db:"shks"`
	OrderAmount float64   `db:"order_amount"`
	SellerID    int64     `db:"seller_id"`
}

func convertToDBPromotion(_ context.Context, in entity.Promotion) *promotionDB {
	return &promotionDB{
		ID:          in.ID,
		ExternalID:  in.ExternalID,
		Name:        in.Name,
		Type:        in.Type,
		Status:      in.Status,
		ChangeTime:  in.ChangeTime,
		CreateTime:  in.CreateTime,
		DateStart:   in.DateStart,
		DateEnd:     in.DateEnd,
		Views:       in.Views,
		Clicks:      in.Clicks,
		CTR:         in.CTR,
		CPC:         in.CPC,
		Spent:       in.Spent,
		Orders:      in.Orders,
		CR:          in.CR,
		SHKs:        in.SHKs,
		OrderAmount: in.OrderAmount,
		SellerID:    in.SellerID,
	}
}

func (c promotionDB) convertToEntityPromotion(_ context.Context) *entity.Promotion {
	return &entity.Promotion{
		ID:          c.ID,
		ExternalID:  c.ExternalID,
		Name:        c.Name,
		Type:        c.Type,
		Status:      c.Status,
		ChangeTime:  c.ChangeTime,
		CreateTime:  c.CreateTime,
		DateStart:   c.DateStart,
		DateEnd:     c.DateEnd,
		Views:       c.Views,
		Clicks:      c.Clicks,
		CTR:         c.CTR,
		CPC:         c.CPC,
		Spent:       c.Spent,
		Orders:      c.Orders,
		CR:          c.CR,
		SHKs:        c.SHKs,
		OrderAmount: c.OrderAmount,
		SellerID:    c.SellerID,
	}
}
