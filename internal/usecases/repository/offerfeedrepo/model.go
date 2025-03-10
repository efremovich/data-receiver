package offerfeedrepo

import (
	"context"
	"database/sql"
	"strings"

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

type vkFeedDB struct {
	Title       string          `db:"title"`
	VendorCode  string          `db:"code"`
	Description string          `db:"description"`
	MediaLinks  string          `db:"media_links"`
	Subject     sql.NullString  `db:"subject"`
	Color       sql.NullString  `db:"color"`
	Gender      sql.NullString  `db:"gender"`
	Price       sql.NullFloat64 `db:"price"`
	Size        sql.NullString  `db:"size"`
	ExternalID  sql.NullInt64   `db:"external_id"`
	SellerName  string          `db:"seller_name"`
}

func (c vkFeedDB) ConvertToEntityVKCard(_ context.Context) *entity.VKCard {
	trimBraces := strings.Trim(c.MediaLinks, "{}")
	mediaLinks := strings.Split(trimBraces, ",")

	// Опционально: удаление пробелов вокруг ссылок
	for i, link := range mediaLinks {
		mediaLinks[i] = strings.TrimSpace(link)
	}
	return &entity.VKCard{
		VendorCode:  c.VendorCode,
		Subject:     repository.NullStringToString(c.Subject),
		Color:       repository.NullStringToString(c.Color),
		Title:       c.Title,
		Gender:      repository.NullStringToString(c.Gender),
		Description: c.Description,
		MediaLinks:  mediaLinks,
		Price:       repository.NullFloatToFloat(c.Price),
		MaxSize:     repository.NullStringToString(c.Size),
		ExternalID:  repository.NullIntToInt(c.ExternalID),
		SellerName:  c.SellerName,
	}
}
