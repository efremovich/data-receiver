package wbfetcher

import (
	"context"
	"net/http"
	"time"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
)

const SellerType = "wb"

func New(_ context.Context, cfg config.Seller) webapi.ExtAPIFetcher {
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
