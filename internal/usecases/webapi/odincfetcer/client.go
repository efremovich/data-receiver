package odincfetcer

import (
	"context"
	"net/http"
	"time"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
)

const SellerType = "odinc"

func New(_ context.Context, cfg config.SellerOdinC) []webapi.ExtAPIFetcher {
	timeout := time.Second * time.Duration(cfg.ProcessTimeoutSeconds)
	c := &http.Client{
		Timeout: timeout,
	}
	clients := []webapi.ExtAPIFetcher{}
	client := &odincAPIclientImp{
		client:   c,
		addr:     cfg.URL,
		login:    cfg.Login,
		password: cfg.Password,
	}
	clients = append(clients, client)

	return clients
}

type odincAPIclientImp struct {
	client   *http.Client
	addr     string
	login    string
	password string
}
