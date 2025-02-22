package ozonfetcher

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
	marketPlaceAPIURL = "https://api-seller.ozon.ru"
)

func New(_ context.Context, cfg config.Config, metrics metrics.Collector) []webapi.ExtAPIFetcher {
	timeout := time.Second * time.Duration(cfg.ProcessTimeoutSeconds)

	c := &http.Client{
		Timeout: timeout,
	}

	clients := []webapi.ExtAPIFetcher{}

	for _, mpConfig := range cfg.MarketPlaces {
		if mpConfig.Type == string(entity.Ozon) {
			marketPlace := entity.MarketPlace{
				ExternalID: mpConfig.ID,
				Title:      mpConfig.Name,
				IsEnabled:  true,
				Type:       entity.OdinAss,
			}
			cred := strings.Split(mpConfig.Token, ":")
			client := &apiClientImp{
				client:      c,
				clientID:    cred[0],
				apiKey:      cred[1],
				marketPlace: marketPlace,

				timeout: timeout,

				metric: metrics,
			}
			clients = append(clients, client)
		}
	}

	return clients
}

type apiClientImp struct {
	client   *http.Client
	apiKey   string
	clientID string

	timeout     time.Duration
	marketPlace entity.MarketPlace

	metric metrics.Collector
}

func (c *apiClientImp) GetMarketPlace() entity.MarketPlace {
	return c.marketPlace
}

func (odinc *apiClientImp) Ping(ctx context.Context) error {
	return nil
}
