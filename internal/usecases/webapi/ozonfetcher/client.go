package ozonfetcher

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
	"github.com/efremovich/data-receiver/pkg/metrics"
)

const (
	marketPlaceAPIURL string = "https://api-seller.ozon.ru"

	// products.
	descriptionCategoryAttrMethod string = "/v1/description-category/attribute"
	productInfoListMethod         string = "/v3/product/info/list"
	productInfoAttrMetod          string = "/v4/product/info/attributes"
	descriptionCategoryMethod     string = "/v1/description-category/tree"
	productListMethod             string = "/v3/product/list"
	// orders and sales.
	fboPostingListMethod string = "/v2/posting/fbo/list"
	// stocks.
	supplyOrderListMethod    string = "/v2/supply-order/list"
	supplyOrderGetMethod     string = "/v2/supply-order/get"
	suppolyOrderBundleMethod string = "/v1/supply-order/bundle"

	requestItemLimit int = 1000
)

var ozonHeaders map[string]map[string]string

func New(_ context.Context, cfg config.Config, metrics metrics.Collector) []webapi.ExtAPIFetcher {
	ozonHeaders = make(map[string]map[string]string)

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

			headers := make(map[string]string)
			headers["Client-Id"] = client.clientID
			headers["Api-Key"] = client.apiKey
			headers["Content-Type"] = "application/json"

			ozonHeaders[marketPlace.ExternalID] = headers

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

func (wb *apiClientImp) GetCosts(_ context.Context, _ entity.PackageDescription) ([]entity.Cost, error) {
	return nil, fmt.Errorf("не реализовано")
}

func (ozon *apiClientImp) GetMarketPlace() entity.MarketPlace {
	return ozon.marketPlace
}

func (ozon *apiClientImp) Ping(_ context.Context) error {
	return nil
}
