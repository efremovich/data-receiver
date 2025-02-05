package offerfeedrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type offerDB struct {
	ID          int64    `db:"id"`
	Available   bool     `db:"available"`
	GroupID     string   `db:"group_id"`
	Name        string   `db:"name"`
	Similar     string   `db:"similar"`
	Price       float32  `db:"price"`
	OldPrice    float32  `db:"old_price"`
	Barcode     string   `db:"barcode"`
	VendorCode  string   `db:"vendor_code"`
	MarketIDs   []int64  `db:"market_id"`
	Vendor      string   `db:"vendor"`
	Pictures    []string `db:"picture"`
	CategoryIDs []int64  `db:"category_id"`
	Description string   `db:"description"`
}

func (c offerDB) ConvertToEntityOffer(_ context.Context) *entity.Offer {
	return &entity.Offer{
		ID:           c.ID,
		Available:    c.Available,
		GroupID:      c.GroupID,
		Name:         c.Name,
		Similar:      c.Similar,
		Price:        c.Price,
		Barcode:      c.Barcode,
		URL:          "",
		VendorCode:   c.VendorCode,
		Sort:         1,
		Vendor:       c.Vendor,
		Rating:       0,
		ReviewsCount: 0,
		Description:  c.Description,
		OldPrice:     c.OldPrice,
		CategoryIDs:  c.CategoryIDs,
		Pictures:     c.Pictures,
		MarketIDs:    c.MarketIDs,
		Params:       []entity.Param{},
		Badges:       []entity.Badge{},
	}
}
