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
		CategoryIDs:  c.CategoryIDs,
		Pictures:     c.Pictures,
		MarketIDs:    c.MarketIDs,
		Params:       []entity.Param{},
		Badges:       []entity.Badge{},
	}
}

type stockDB struct {
	ID        int64           `db:"id"`
	Quantity  float32         `db:"quantity"`
	Price     sql.NullFloat64 `db:"price"`
	OldPrice  sql.NullFloat64 `db:"old_price"`
	StorageID int64           `db:"storage_id"`
	SellerID  int64           `db:"seller_id"`
}

type storageDB struct {
	ID         int64          `db:"id"`
	Name       string         `db:"name"`
	City       string         `db:"city"`
	Type       sql.NullString `db:"type"`
	Address    string         `db:"address"`
	Lat        string         `db:"lat"`
	Lon        string         `db:"lon"`
	Region     string         `db:"region"`
	WorkTime   string         `db:"work_time"`
	Phone      string         `db:"phone"`
	Icon       string         `db:"icon"`
	SellerID   int64          `db:"seller_id"`
	SellerName string         `db:"seller_name"`
}

func (c stockDB) ConvertToEntityStock(_ context.Context) *entity.OfferStock {
	return &entity.OfferStock{
		ID:        c.ID,
		Quantity:  c.Quantity,
		Price:     repository.NullFloatToFloat(c.Price),
		OldPrice:  repository.NullFloatToFloat(c.OldPrice),
		StorageID: c.StorageID,
		SellerID:  c.SellerID,
	}
}

func (c storageDB) ConvertToEntityStorage(_ context.Context) *entity.OfferStorage {
	return &entity.OfferStorage{
		ID:       c.ID,
		Name:     c.Name,
		City:     c.City,
		Type:     repository.NullStringToString(c.Type),
		Address:  c.Address,
		Lat:      c.Lat,
		Lon:      c.Lon,
		Region:   c.Region,
		WorkTime: c.WorkTime,
		Phone:    c.Phone,
		Icon:     c.Icon,
		Seller: entity.SellerStock{
			SellerID:   c.SellerID,
			SellerName: c.Name,
		},
	}
}
