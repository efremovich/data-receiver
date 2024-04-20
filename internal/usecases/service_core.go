package usecases

import (
	"context"

	conf "github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/wbcontentrepo"
	"github.com/efremovich/data-receiver/pkg/broker/brokerpublisher"
	"github.com/efremovich/data-receiver/pkg/metrics"

	aerror "github.com/efremovich/data-receiver/pkg/aerror"
)

type ReceiverCoreService interface {
	ReceiveCards(ctx context.Context, sellerTitle string) aerror.AError

	PingDB(ctx context.Context) error
	PingNATS(_ context.Context) error
}

type receiverCoreServiceImpl struct {
	cfg              conf.Config
	sellerRepo       sellerrepo.SellerRepo
	cardRepo         cardrepo.CardRepo
	brokerPublisher  brokerpublisher.BrokerPublisher
	metricsCollector metrics.Collector
	wbContentRepo    wbcontentrepo.WBContentRepo
}

func NewPackageReceiverService(cfg conf.Config, wbContentRepo wbcontentrepo.WBContentRepo,
	cardR cardrepo.CardRepo, sellerR sellerrepo.SellerRepo, brokerPublisher brokerpublisher.BrokerPublisher, metricsCollector metrics.Collector,
) ReceiverCoreService {
	service := receiverCoreServiceImpl{
		cfg:              cfg,
		brokerPublisher:  brokerPublisher,
		metricsCollector: metricsCollector,
		cardRepo:         cardR,
		sellerRepo:       sellerR,
		wbContentRepo:    wbContentRepo,
	}

	return &service
}

func (s *receiverCoreServiceImpl) PingNATS(_ context.Context) error {
	return s.brokerPublisher.Ping()
}

func (s *receiverCoreServiceImpl) PingDB(ctx context.Context) error {
	return s.cardRepo.Ping(ctx)
}
