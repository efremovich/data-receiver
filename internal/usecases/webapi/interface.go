package webapi

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type ExtAPIFetcher interface {
	GetCards(ctx context.Context, desc entity.PackageDescription) ([]entity.Card, error)
	GetStocks(ctx context.Context, desc entity.PackageDescription) ([]entity.StockMeta, error)
	GetWarehouses(ctx context.Context) ([]entity.Warehouse, error)
	GetOrders(ctx context.Context, desc entity.PackageDescription) ([]entity.Order, error)
	GetSales(ctx context.Context, desc entity.PackageDescription) ([]entity.Sale, error)
	GetSaleReport(ctx context.Context, desc entity.PackageDescription) ([]entity.SaleReport, error)
	GetCosts(ctx context.Context, desc entity.PackageDescription) ([]entity.Cost, error)
	GetPromotion(ctx context.Context, desc entity.PackageDescription) ([]*entity.Promotion, error)

	Ping(ctx context.Context) error

	GetMarketPlace() entity.MarketPlace
}
