package promotionstatsrepo

import (
	"context"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
)

type promotionStatsDB struct {
	ID          int64     `db:"id"`
	Date        time.Time `db:"date"`
	Views       int       `db:"views"`
	Clicks      int       `db:"clicks"`
	CTR         float64   `db:"ctr"`
	CPC         float64   `db:"cpc"`
	Spent       float64   `db:"spent"`
	Orders      int       `db:"orders"`
	CR          float64   `db:"cr"`
	SHKs        int       `db:"shks"`
	OrderAmount float64   `db:"order_amount"`
	AppType     int       `db:"app_type"`
	PromotionID int64     `db:"promotion_id"`
	CardID      int64     `db:"card_id"`
	SellerID    int64     `db:"seller_id"`
}

func convertToDBPromotionStats(_ context.Context, income entity.PromotionStats) *promotionStatsDB {
	return &promotionStatsDB{
		ID:          income.ID,
		Date:        income.Date,
		Views:       income.Views,
		Clicks:      income.Clicks,
		CTR:         income.CTR,
		CPC:         income.CPC,
		Spent:       income.Spent,
		Orders:      income.Orders,
		CR:          income.CR,
		SHKs:        income.SHKs,
		OrderAmount: income.OrderAmount,
		AppType:     income.AppType,
		PromotionID: income.PromotionID,
		CardID:      income.CardID,
		SellerID:    income.SellerID,
	}
}

func (c promotionStatsDB) convertToEntityPromotionStats(_ context.Context) *entity.PromotionStats {
	return &entity.PromotionStats{
		ID:          c.ID,
		Date:        c.Date,
		Views:       c.Views,
		Clicks:      c.Clicks,
		CTR:         c.CTR,
		CPC:         c.CPC,
		Spent:       c.Spent,
		Orders:      c.Orders,
		CR:          c.CR,
		SHKs:        c.SHKs,
		OrderAmount: c.OrderAmount,
		AppType:     c.AppType,
		PromotionID: c.PromotionID,
		CardID:      c.CardID,
		SellerID:    c.SellerID,
	}
}
