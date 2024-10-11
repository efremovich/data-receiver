package app

import (
	"context"
	"fmt"
	"os"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/controller"
	"github.com/efremovich/data-receiver/internal/usecases"
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
	"github.com/efremovich/data-receiver/internal/usecases/repository/salerepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/seller2cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sizerepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/statusrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/stockrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/warehouserepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/warehousetyperepo"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
	"github.com/efremovich/data-receiver/internal/usecases/webapi/odincfetcer"
	"github.com/efremovich/data-receiver/internal/usecases/webapi/wbfetcher"
	"github.com/efremovich/data-receiver/pkg/alogger"
	"github.com/efremovich/data-receiver/pkg/broker/brokerconsumer"
	"github.com/efremovich/data-receiver/pkg/broker/brokerpublisher"
	"github.com/efremovich/data-receiver/pkg/metrics"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type Application struct {
	Gateway controller.GrpcGatewayServer
}

func New(ctx context.Context, conf config.Config) (*Application, error) {
	alogger.SetDefaultConfig(&alogger.Config{
		Level:  alogger.Level(conf.LogLevel),
		Output: os.Stdout,
	})

	// Сборщик метрик.
	metricsCollector, err := metrics.NewMetricCollector(conf.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании сборщика метрик: %s", err.Error())
	}

	// Брокер издатель.
	brokerPublisher, err := brokerpublisher.NewBrokerPublisher(ctx, conf.BrokerPublisherURL, true)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании брокера издателя: %w", err)
	}

	if err := brokerPublisher.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка при подключении к брокеру издателю: %w ", err)
	}

	// Инициализация брокера потребителя.
	brokerConsumer, err := brokerconsumer.NewBrokerConsumer(ctx, conf.BrokerConsumerURL, true)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании брокера потребителя: %w", err)
	}

	// Подключение к БД.
	conn, err := postgresdb.New(ctx, conf.PGWriterConn, conf.PGReaderConn)
	if err != nil {
		return nil, err
	}

	apiFetcher := make(map[string]webapi.ExtAPIFetcher)
	// TODO Завернем клиентов всех маркетплейсов в мапу
	apiFetcher["wb"] = wbfetcher.New(ctx, conf.Seller.WB)
	apiFetcher["odinc"] = odincfetcer.New(ctx, conf.Seller.OdinC)

	// Репозиторий Cards.
	cardRepo, err := cardrepo.NewCardRepo(ctx, conn)
	if err != nil {
		return nil, err
	}
	// Репозиторий Size.
	sizeRepo, err := sizerepo.NewSizeRepo(ctx, conn)
	if err != nil {
		return nil, err
	}
	// Репозиторий Seller
	sellerRepo, err := sellerrepo.NewSellerRepo(ctx, conn)
	if err != nil {
		return nil, err
	}
	// Репозиторий Brand
	brandRepo, err := brandrepo.NewBrandRepo(ctx, conn)
	if err != nil {
		return nil, err
	}
	// Репозиторий Dimension
	dimensionRepo, err := dimensionrepo.NewDimensionRepo(ctx, conn)
	if err != nil {
		return nil, err
	}
	// Репозитозий Categories
	categoryRepo, err := categoryrepo.NewCategoryRepo(ctx, conn)
	if err != nil {
		return nil, err
	}
	// Репозиторий Characteristic
	charRepo, err := charrepo.NewCharRepo(ctx, conn)
	if err != nil {
		return nil, err
	}
	// Репозиторий Barcode
	barcodeRepo, err := barcoderepo.NewBarcodeRepo(ctx, conn)
	if err != nil {
		return nil, err
	}
	// Репозиторий CardCharacteristic
	cardcharrepo, err := cardcharrepo.NewCharRepo(ctx, conn)
	if err != nil {
		return nil, err
	}
	// Репозиторий WarhouseType
	warehouseTypeRepo, err := warehousetyperepo.NewWarehouseTypeRepo(ctx, conn)
	if err != nil {
		return nil, err
	}
	// Репозиторий Mediafile
	mediafileRepo, err := mediafilerepo.NewMediaFileRepo(ctx, conn)
	if err != nil {
		return nil, err
	}
	// Репозиторий Warhouse
	warehouseRepo, err := warehouserepo.NewWarehouseRepo(ctx, conn)
	if err != nil {
		return nil, err
	}
	// Репозиторий PriceSize
	priceSizeRepo, err := pricerepo.NewPriceRepo(ctx, conn)
	if err != nil {
		return nil, err
	}
	// Репозиторий Stock
	stockRepo, err := stockrepo.NewStockRepo(ctx, conn)
	if err != nil {
		return nil, err
	}
	// Репозиторий seller2Card
	seller2carRepo, err := seller2cardrepo.NewWb2CardRepo(ctx, conn)
	if err != nil {
		return nil, err
	}
	// Репозиторий Order
	orderRepo, err := orderrepo.NewOrderRepo(ctx, conn)
	if err != nil {
		return nil, err
	}

	// Репозиторий Status
	statusRepo, err := statusrepo.NewStatusRepo(ctx, conn)
	if err != nil {
		return nil, err
	}

	// Репозиторий Country
	countryRepo, err := countryrepo.NewCountryRepo(ctx, conn)
	if err != nil {
		return nil, err
	}

	// Репозиторий Region
	regionRepo, err := regionrepo.NewRegionRepo(ctx, conn)
	if err != nil {
		return nil, err
	}

	// Репозиторий District
	districtRepo, err := districtrepo.NewDistrictRepo(ctx, conn)
	if err != nil {
		return nil, err
	}

	// Репозиторий Sale
	saleRepo, err := salerepo.NewSaleRepo(ctx, conn)
	if err != nil {
		return nil, err
	}
	// Основной бизнес-сервис.
	packageReceiverCoreService := usecases.NewPackageReceiverService(
		conf,

		cardRepo,
		sizeRepo,
		sellerRepo,
		brandRepo,
		charRepo,
		cardcharrepo,
		barcodeRepo,
		categoryRepo,
		dimensionRepo,
		mediafileRepo,
		priceSizeRepo,
		stockRepo,
		seller2carRepo,
		orderRepo,
		statusRepo,
		countryRepo,
		regionRepo,
		districtRepo,
		saleRepo,

		warehouseRepo,
		warehouseTypeRepo,

		brokerPublisher,
		apiFetcher,
		metricsCollector)

	gw, err := controller.NewGatewayServer(ctx, conf, packageReceiverCoreService, metricsCollector, brokerConsumer)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании gateway сервиса: %s", err.Error())
	}

	return &Application{
		Gateway: gw,
	}, nil
}

func (a *Application) Start(ctx context.Context) error {
	err := a.Gateway.Start(ctx)
	if err != nil {
		return fmt.Errorf("ошибка при работе gateway сервиса: %s", err.Error())
	}

	return nil
}
