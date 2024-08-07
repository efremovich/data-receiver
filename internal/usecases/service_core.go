package usecases

import (
	"context"

	conf "github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/brandrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardcharrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/categoryrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/charrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/wb2cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/webapi/wbfetcher"
	"github.com/efremovich/data-receiver/pkg/broker/brokerpublisher"
	"github.com/efremovich/data-receiver/pkg/metrics"

	aerror "github.com/efremovich/data-receiver/pkg/aerror"
)

type ReceiverCoreService interface {
	ReceiveCards(ctx context.Context, desc entity.PackageDescription) aerror.AError

	PingDB(ctx context.Context) error
	PingNATS(_ context.Context) error
}

type receiverCoreServiceImpl struct {
	cfg conf.Config
	// Блок репозиотриев
	sellerRepo   sellerrepo.SellerRepo
	cardRepo     cardrepo.CardRepo
	brandRepo    brandrepo.BrandRepo
	charRepo     charrepo.CharRepo
  cardCharRepo cardcharrepo.CardCharRepo
	categoryRepo categoryrepo.CategoryRepo
	wb2cardrepo  wb2cardrepo.Wb2CardRepo

	brokerPublisher  brokerpublisher.BrokerPublisher
	apiFetcher       map[string]wbfetcher.ExtApiFetcher
	metricsCollector metrics.Collector
}

func NewPackageReceiverService(
	cfg conf.Config,

	cardRepo cardrepo.CardRepo,
	sellerRepo sellerrepo.SellerRepo,
	brandRepo brandrepo.BrandRepo,
	charRepo charrepo.CharRepo,
	cardcharRepo cardcharrepo.CardCharRepo,
	categoryRepo categoryrepo.CategoryRepo,
	wb2cardrepo wb2cardrepo.Wb2CardRepo,

	brokerPublisher brokerpublisher.BrokerPublisher,
	apiFetcher map[string]wbfetcher.ExtApiFetcher,
	metricsCollector metrics.Collector,
) ReceiverCoreService {
	service := receiverCoreServiceImpl{
		cfg: cfg,

		cardRepo:     cardRepo,
		sellerRepo:   sellerRepo,
		brandRepo:    brandRepo,
		charRepo:     charRepo,
    cardCharRepo: cardcharRepo,
		categoryRepo: categoryRepo,
		wb2cardrepo:  wb2cardrepo,

		brokerPublisher:  brokerPublisher,
		apiFetcher:       apiFetcher,
		metricsCollector: metricsCollector,
	}

	return &service
}

func (s *receiverCoreServiceImpl) PingNATS(_ context.Context) error {
	return s.brokerPublisher.Ping()
}

func (s *receiverCoreServiceImpl) PingDB(ctx context.Context) error {
	return s.cardRepo.Ping(ctx)
}
