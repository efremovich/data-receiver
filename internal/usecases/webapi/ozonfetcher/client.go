package ozonfetcher

import (
	"context"
	"time"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
	"github.com/efremovich/data-receiver/pkg/metrics"
)

const SellerType = "ozon"

func New(_ context.Context, cfg config.SellerOZON, metrics metrics.Collector) []webapi.ExtAPIFetcher {
	timeout := time.Second * time.Duration(cfg.ProcessTimeoutSeconds)

	configs := []webapi.ExtAPIFetcher{}

	for i := 0; i < len(cfg.APIKey); i++ {
		cfg := &ozonAPIclientImp{
			baseURL:  cfg.URL,
			apiKey:   cfg.APIKey[i],
			clientID: cfg.ClientID[i],
			timeout:  timeout,
		}
		configs = append(configs, cfg)
	}

	return configs
}

type ozonAPIclientImp struct {
	baseURL  string
	apiKey   string
	clientID string
	timeout  time.Duration

	metric metrics.Collector
}

func (odinc *ozonAPIclientImp) Ping(ctx context.Context) error {
	return nil
}
