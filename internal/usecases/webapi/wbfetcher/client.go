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
	cardListMethod    string = "/content/v2/get/cards/list?locale=ru"

	reportDetaioByPeriodMethod string = "/api/v5/supplier/reportDetailByPeriod"

	saleReportResponseLimit int = 100000 // Максимальное количество строк отчета, возвращаемых методом. Не может быть более 100000.
	cardRequestLimit        int = 100
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

func (wb *apiClientImp) GetMarketPlace() entity.MarketPlace {
	return wb.marketPlace
}
