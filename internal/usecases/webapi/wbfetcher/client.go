package wbfetcher

import (
	"context"
	"net/http"
	"time"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/entity"
)

const SellerType = "wb"

type ExtAPIFetcher interface {
	GetCards(ctx context.Context, desc entity.PackageDescription) ([]entity.Card, error)
	GetStocks(ctx context.Context, desc entity.PackageDescription) ([]entity.StockMeta, error)
	GetWarehouses(ctx context.Context) ([]entity.Warehouse, error)
	GetOrders(ctx context.Context, desc entity.PackageDescription) ([]entity.Order, error)
	GetSales(ctx context.Context, desc entity.PackageDescription) ([]entity.Sale, error)

	Ping(ctx context.Context) error
}

func New(_ context.Context, cfg config.Seller) ExtAPIFetcher {
	timeout := time.Second * time.Duration(cfg.ProcessTimeoutSeconds)

	c := &http.Client{
		Timeout: timeout,
	}
	client := &wbAPIclientImp{
		client:    c,
		token:     cfg.Token,
		addr:      cfg.URL,
		addrStat:  cfg.URLStat,
		tokenStat: cfg.TokenStat,
	}

	return client
}

type wbAPIclientImp struct {
	client    *http.Client
	addr      string
	addrStat  string
	token     string
	tokenStat string
}
