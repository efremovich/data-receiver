package usecases

import (
	"context"

	conf "github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/barcoderepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/brandrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardcharrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/categoryrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/charrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/countryrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/dimensionrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/districtrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/mediafilerepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/orderrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/pricerepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/regionrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sizerepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/statusrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/stockrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/warehouserepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/warehousetyperepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/wb2cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/webapi/wbfetcher"
	"github.com/efremovich/data-receiver/pkg/broker/brokerpublisher"
	"github.com/efremovich/data-receiver/pkg/metrics"

	aerror "github.com/efremovich/data-receiver/pkg/aerror"
)

type ReceiverCoreService interface {
	ReceiveCards(ctx context.Context, desc entity.PackageDescription) aerror.AError
	ReceiveWarehouses(ctx context.Context) aerror.AError
	ReceiveStocks(ctx context.Context, desc entity.PackageDescription) aerror.AError
	ReceiveOrders(ctx context.Context, desc entity.PackageDescription) aerror.AError

	PingDB(ctx context.Context) error
	PingNATS(_ context.Context) error
}

type receiverCoreServiceImpl struct {
	cfg conf.Config
	// Блок репозиотриев
	sellerRepo    sellerrepo.SellerRepo
	cardRepo      cardrepo.CardRepo
	sizerepo      sizerepo.SizeRepo
	brandRepo     brandrepo.BrandRepo
	charRepo      charrepo.CharRepo
	cardCharRepo  cardcharrepo.CardCharRepo
	barcodeRepo   barcoderepo.BarcodeRepo
	categoryRepo  categoryrepo.CategoryRepo
	dimensionrepo dimensionrepo.DimensionsRepo
	mediafilerepo mediafilerepo.MediaFileRepo
	pricesizerepo pricerepo.PriceRepo
	stockrepo     stockrepo.StockRepo
	orderrepo     orderrepo.OrderRepo
	statusrepo    statusrepo.StatusRepo
	countryrepo   countryrepo.CountryRepo
	regionrepo    regionrepo.RegoinRepo
	districtrepo  districtrepo.DistrictRepo

	wb2cardrepo wb2cardrepo.Wb2CardRepo

	// Блок остатков
	warehouserepo     warehouserepo.WarehouseRepo
	warehousetyperepo warehousetyperepo.WarehouseTypeRepo

	brokerPublisher  brokerpublisher.BrokerPublisher
	apiFetcher       map[string]wbfetcher.ExtApiFetcher
	metricsCollector metrics.Collector
}

func NewPackageReceiverService(
	cfg conf.Config,

	cardRepo cardrepo.CardRepo,
	sizerepo sizerepo.SizeRepo,
	sellerRepo sellerrepo.SellerRepo,
	brandRepo brandrepo.BrandRepo,
	charRepo charrepo.CharRepo,
	cardcharRepo cardcharrepo.CardCharRepo,
	barcoderepo barcoderepo.BarcodeRepo,
	categoryRepo categoryrepo.CategoryRepo,
	dimensionrepo dimensionrepo.DimensionsRepo,
	mediafilerepo mediafilerepo.MediaFileRepo,
	pricesizerepo pricerepo.PriceRepo,
	stockrepo stockrepo.StockRepo,
	wb2cardrepo wb2cardrepo.Wb2CardRepo,
	orderrepo orderrepo.OrderRepo,
	statusrepo statusrepo.StatusRepo,
	countryrepo countryrepo.CountryRepo,
	regionrepo regionrepo.RegoinRepo,
	districtrepo districtrepo.DistrictRepo,

	warehouserepo warehouserepo.WarehouseRepo,
	warehousetyperepo warehousetyperepo.WarehouseTypeRepo,

	brokerPublisher brokerpublisher.BrokerPublisher,
	apiFetcher map[string]wbfetcher.ExtApiFetcher,
	metricsCollector metrics.Collector,
) ReceiverCoreService {
	service := receiverCoreServiceImpl{
		cfg: cfg,

		cardRepo:      cardRepo,
		sizerepo:      sizerepo,
		sellerRepo:    sellerRepo,
		brandRepo:     brandRepo,
		charRepo:      charRepo,
		cardCharRepo:  cardcharRepo,
		barcodeRepo:   barcoderepo,
		categoryRepo:  categoryRepo,
		dimensionrepo: dimensionrepo,
		mediafilerepo: mediafilerepo,
		pricesizerepo: pricesizerepo,
		stockrepo:     stockrepo,
		wb2cardrepo:   wb2cardrepo,
		orderrepo:     orderrepo,
		statusrepo:    statusrepo,
		countryrepo:   countryrepo,
		regionrepo:    regionrepo,
		districtrepo:  districtrepo,

		warehouserepo:     warehouserepo,
		warehousetyperepo: warehousetyperepo,

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
