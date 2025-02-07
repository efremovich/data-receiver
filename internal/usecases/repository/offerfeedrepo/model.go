package offerfeedrepo

import (
	"context"
	"database/sql"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
)

type offerDB struct {
	ID          int64           `db:"id"`
	Available   bool            `db:"available"`
	GroupID     string          `db:"group_id"`
	Name        string          `db:"name"`
	Similar     string          `db:"similar"`
	Price       sql.NullFloat64 `db:"price"`
	OldPrice    sql.NullFloat64 `db:"old_price"`
	Barcode     sql.NullString  `db:"barcode"`
	VendorCode  string          `db:"vendor_code"`
	MarketIDs   string          `db:"market_id"`
	Vendor      string          `db:"vendor"`
	Pictures    string          `db:"picture"`
	CategoryIDs string          `db:"category_id"`
	Description string          `db:"description"`
}

func (c offerDB) ConvertToEntityOffer(_ context.Context) *entity.Offer {
	return &entity.Offer{
		ID:           c.ID,
		Available:    c.Available,
		GroupID:      c.GroupID,
		Name:         c.Name,
		Similar:      c.Similar,
		Price:        repository.NullFloatToFloat(c.Price),
		Barcode:      repository.NullStringToString(c.Barcode),
		URL:          "",
		VendorCode:   c.VendorCode,
		Sort:         1,
		Vendor:       c.Vendor,
		Rating:       0,
		ReviewsCount: 0,
		Description:  c.Description,
		OldPrice:     repository.NullFloatToFloat(c.OldPrice),
		// CategoryIDs:  c.CategoryIDs,
		// Pictures:     c.Pictures,
		// MarketIDs:    c.MarketIDs,
		Params: []entity.Param{},
		Badges: []entity.Badge{},
	}
}
