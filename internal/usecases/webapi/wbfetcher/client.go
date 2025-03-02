package wbfetcher

import (
	"context"
	"net/http"
	"time"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
	"github.com/efremovich/data-receiver/pkg/metrics"
)

const (
	marketPlaceAPIURL string = "https://marketplace-api.wildberries.ru"
	contentAPIURL     string = "https://content-api.wildberries.ru"
	statisticAPIURL   string = "https://statistics-api.wildberries.ru"
)

func New(_ context.Context, cfg config.Config, metrics metrics.Collector) []webapi.ExtAPIFetcher {
	timeout := time.Second * time.Duration(cfg.ProcessTimeoutSeconds)

	c := &http.Client{
		Timeout: timeout,
	}

	clients := []webapi.ExtAPIFetcher{}

	for _, mpConfig := range cfg.MarketPlaces {
		if mpConfig.Type == string(entity.Wildberries) {
			marketPlace := entity.MarketPlace{
				ExternalID: mpConfig.ID,
				Title:      mpConfig.Name,
				IsEnabled:  true,
				Type:       entity.Wildberries,
			}
			client := &apiClientImp{
				client:      c,
				token:       mpConfig.Token,
				marketPlace: marketPlace,

				metric: metrics,
			}
			clients = append(clients, client)
		}
	}
	return clients
}

type apiClientImp struct {
	client *http.Client
	token  string

	marketPlace entity.MarketPlace

	metric metrics.Collector
}

func (c *apiClientImp) GetMarketPlace() entity.MarketPlace {
	return c.marketPlace
}
