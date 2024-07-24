package wbfetcher

import (
	"context"
	"net/http"
	"time"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/entity"
)

const SellerType = "wb"

type ExtApiFetcher interface {
	GetCards(ctx context.Context, desc entity.PackageDescription) ([]entity.Card, error)
	GetStocks(ctx context.Context, desc entity.PackageDescription) ([]entity.StockMeta, error)

	Ping(ctx context.Context) error
}

func New(_ context.Context, cfg config.Seller) ExtApiFetcher {
	timeout := time.Second * time.Duration(cfg.ProcessTimeoutSeconds)

	c := &http.Client{
		Timeout: timeout,
	}
	client := &wbAPIclientImp{client: c, token: cfg.Token, addr: cfg.URL, tokenStat: cfg.TokenStat}

	return client
}

type wbAPIclientImp struct {
	client *http.Client
	addr   string
	token  string
  tokenStat string
}

