package odincfetcer

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
	"github.com/efremovich/data-receiver/pkg/metrics"
)

const (
	marketPlaceAPIURL = "http://94.199.4.215:1281/rbb_cut/"
)

func New(_ context.Context, cfg config.Config, metrics metrics.Collector) []webapi.ExtAPIFetcher {
	timeout := time.Second * time.Duration(cfg.ProcessTimeoutSeconds)

	c := &http.Client{
		Timeout: timeout,
	}

	clients := []webapi.ExtAPIFetcher{}

	for _, mpConfig := range cfg.MarketPlaces {
		if mpConfig.Type == string(entity.OdinAss) {
			marketPlace := entity.MarketPlace{
				ExternalID: mpConfig.ID,
				Title:      mpConfig.Name,
				IsEnabled:  true,
				Type:       entity.OdinAss,
			}
			cred := strings.Split(mpConfig.Token, ":")
			client := &apiClientImp{
				client:      c,
				login:       cred[0],
				password:    cred[1],
				marketPlace: marketPlace,

				metric: metrics,
			}
			clients = append(clients, client)
		}
	}

	return clients
}

type apiClientImp struct {
	client   *http.Client
	login    string
	password string

	marketPlace entity.MarketPlace

	metric metrics.Collector
}

func (odinc *apiClientImp) GetMarketPlace() entity.MarketPlace {
	return odinc.marketPlace
}
