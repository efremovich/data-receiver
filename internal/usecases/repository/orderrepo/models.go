package orderrepo

import (
	"context"
	"database/sql"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
)

type orderDB struct {
	ID                   int64   `db:"id"`
	ExternalID           string  `db:"external_id"`
	Price                float64 `db:"price"`
	PriceWithoutDiscount float64 `db:"price_without_discount"`
	PriceFinal           float64 `db:"price_final"`
	StatusID             int64   `db:"status_id"`
	Direction            string  `db:"direction"`
	Type                 string  `db:"type"`
	Sale                 float64 `db:"sale"`
	IsCancel             bool    `db:"is_cancel"`

	Quantity  int          `db:"quantity"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`

	SellerID    int64 `db:"seller_id"`
	CardID      int64 `db:"card_id"`
	WarehouseID int64 `db:"warehouse_id"`
	RegionID    int64 `db:"region_id"`
	PriceSizeID int64 `db:"price_size_id"`
}

func convertToDBOrder(_ context.Context, income *entity.Order) *orderDB {
	return &orderDB{
		ID:                   income.ID,
		ExternalID:           income.ExternalID,
		Price:                income.Price,
		PriceWithoutDiscount: income.PriceWithoutDiscount,
		PriceFinal:           income.PriceFinal,

		Quantity:  income.Quantity,
		Direction: income.Direction,
		Type:      income.Type,
		IsCancel:  income.IsCancel,
		CreatedAt: repository.TimeToNullInt(income.CreatedAt),
		UpdatedAt: repository.TimeToNullInt(income.UpdatedAt),

		StatusID:    income.Status.ID,
		SellerID:    income.Seller.ID,
		WarehouseID: income.Warehouse.ID,
		RegionID:    income.Region.ID,
		CardID:      income.Card.ID,
		PriceSizeID: income.PriceSize.ID,
	}
}

func (c orderDB) convertToEntityOrder(_ context.Context) *entity.Order {
	return &entity.Order{
		ID:                   c.ID,
		ExternalID:           c.ExternalID,
		Price:                c.Price,
		PriceWithoutDiscount: c.PriceWithoutDiscount,
		PriceFinal:           c.PriceFinal,
		Type:                 c.Type,
		IsCancel:             c.IsCancel,
		Direction:            c.Direction,
		Sale:                 c.Sale,
		Quantity:             c.Quantity,
		CreatedAt:            c.CreatedAt.Time,
		UpdatedAt:            c.UpdatedAt.Time,
		Status: &entity.Status{
			ID: c.StatusID,
		},
		Region: &entity.Region{
			ID: c.RegionID,
		},
		Warehouse: &entity.Warehouse{
			ID: c.WarehouseID,
		},
		Seller: &entity.MarketPlace{
			ID: c.SellerID,
		},
		Card: &entity.Card{
			ID: c.CardID,
		},
		PriceSize: &entity.PriceSize{
			ID: c.PriceSizeID,
		},
	}
}

// Добавляем метод для конвертации среза заказов.
func convertToDBOrders(ctx context.Context, orders []*entity.Order) []*orderDB {
	dbOrders := make([]*orderDB, len(orders))
	for i, order := range orders {
		dbOrders[i] = convertToDBOrder(ctx, order)
	}

	return dbOrders
}
