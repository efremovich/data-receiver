package app

import (
	"context"
	"fmt"
	"os"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/controller"
	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases"
	"github.com/efremovich/data-receiver/internal/usecases/repository/tprepo"
	"github.com/efremovich/data-receiver/pkg/alogger"
	"github.com/efremovich/data-receiver/pkg/broker"
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

	entity.AddUserErrorMessages()

	// Cборщик метрик.
	metricsCollector, err := metrics.NewMetricCollector(conf.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании сборщика метрик: %s", err.Error())
	}

	// Брокер сообщений.
	natsClient, err := broker.NewNats(ctx, conf, true)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании подключения к NATS: %s", err.Error())
	}

	// Подключение к БД.
	conn, err := postgresdb.New(ctx, conf.PGWriterConn, conf.PGReaderConn)
	if err != nil {
		return nil, err
	}

	// Репозиторий ТП.
	tpRepo, err := tprepo.NewTransportPackageRepo(ctx, conn)
	if err != nil {
		return nil, err
	}

	// Основной бизнес-сервис.
	packageReceiverCoreService := usecases.NewPackageReceiverService(tpRepo, natsClient, metricsCollector)

	gw, err := controller.NewGatewayServer(conf.Gateway, packageReceiverCoreService, metricsCollector)
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
