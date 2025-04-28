package salerepo

import (
	"context"
	"database/sql"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
)

type saleDB struct {
	ID         int64  `db:"id"`
	ExternalID string `db:"external_id"`

	Price      float64 `db:"price"`
	Discount   float64 `db:"discount"`
	FinalPrice float64 `db:"final_price"`
	Type       string  `db:"type"`
	ForPay     float64 `db:"for_pay"`

	Quantity  int          `db:"quantity"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`

	OrderID     int64 `db:"order_id"`
	SellerID    int64 `db:"seller_id"`
	CardID      int64 `db:"card_id"`
	WarehouseID int64 `db:"warehouse_id"`
	RegionID    int64 `db:"region_id"`
	PriceSizeID int64 `db:"price_size_id"`
}

func convertToDBSale(_ context.Context, income *entity.Sale) *saleDB {
	return &saleDB{
		ID:         income.ID,
		ExternalID: income.ExternalID,
		Price:      income.Price,
		Discount:   income.DiscountP,
		FinalPrice: income.FinalPrice,
		Type:       income.Type,
		ForPay:     income.ForPay,
		Quantity:   income.Quantity,
		CreatedAt:  repository.TimeToNullInt(income.CreatedAt),
		UpdatedAt:  repository.TimeToNullInt(income.UpdatedAt),

		OrderID:     income.Order.ID,
		SellerID:    income.Seller.ID,
		CardID:      income.Card.ID,
		WarehouseID: income.Warehouse.ID,
		RegionID:    income.Region.ID,
		PriceSizeID: income.PriceSize.ID,
	}
}

func (s saleDB) convertToEntitySale(_ context.Context) *entity.Sale {
	return &entity.Sale{
		ID:         s.ID,
		ExternalID: s.ExternalID,
		Price:      s.Price,
		DiscountP:  s.Discount,
		FinalPrice: s.FinalPrice,
		Type:       s.Type,
		ForPay:     s.ForPay,
		Quantity:   s.Quantity,
		CreatedAt:  s.CreatedAt.Time,
		UpdatedAt:  s.UpdatedAt.Time,

		Order:     &entity.Order{ID: s.OrderID},
		Seller:    &entity.MarketPlace{ID: s.SellerID},
		Card:      &entity.Card{ID: s.CardID},
		Warehouse: &entity.Warehouse{ID: s.WarehouseID},
		Region:    &entity.Region{ID: s.RegionID},
		PriceSize: &entity.PriceSize{ID: s.PriceSizeID},
	}
}
