package wbfetcher

import (
	"context"
	"net/http"
	"time"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
)

const SellerType = "wb"

func New(_ context.Context, cfg config.SellerWB) []webapi.ExtAPIFetcher {
	timeout := time.Second * time.Duration(cfg.ProcessTimeoutSeconds)

	c := &http.Client{
		Timeout: timeout,
	}

	clients := []webapi.ExtAPIFetcher{}

	for i := 0; i < len(cfg.Token); i++ {
		client := &wbAPIclientImp{
			client:          c,
			token:           cfg.Token[i],
			addr:            cfg.URL,
			addrContent:     cfg.URLContent,
			addrMarketPlace: cfg.URLMarketPlace,
			addrStat:        cfg.URLStat,
			tokenStat:       cfg.TokenStat[i],
		}
		clients = append(clients, client)
	}

	return clients
}

type wbAPIclientImp struct {
	client          *http.Client
	addr            string
	addrContent     string
	addrMarketPlace string
	addrStat        string
	token           string
	tokenStat       string
}
