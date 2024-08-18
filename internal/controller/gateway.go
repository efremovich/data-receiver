// Package v1 implements routing paths. Each services in own file.
package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/valyala/fasthttp/reuseport"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/efremovich/data-receiver/pkg/alogger"
	"github.com/efremovich/data-receiver/pkg/broker/brokerconsumer"

	conf "github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/controller/middleware"
	"github.com/efremovich/data-receiver/internal/usecases"
	desc "github.com/efremovich/data-receiver/pkg/data-receiver-service"
	"github.com/efremovich/data-receiver/pkg/metrics"
)

const updateIntervalDefault = time.Hour * 10

type GrpcGatewayServer interface {
	Start(ctx context.Context) error
	gracefulStop()
}

type grpcGatewayServerImpl struct {
	httpServer *fiber.App
	grpcServer *grpc.Server
	cfg        conf.Config

	metricsCollector metrics.Collector

	core usecases.ReceiverCoreService

	brokerConsumer brokerconsumer.BrokerConsumer

	desc.UnimplementedCardReceiverServer
}

func NewGatewayServer(ctx context.Context, cfg conf.Config, packageReceiver usecases.ReceiverCoreService, metricsCollector metrics.Collector, broker brokerconsumer.BrokerConsumer) (GrpcGatewayServer, error) {
	gwmux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := desc.RegisterCardReceiverHandlerFromEndpoint(
		context.Background(),
		gwmux,
		cfg.Gateway.GRPC.Host+":"+cfg.Gateway.GRPC.Port,
		opts,
	)
	if err != nil {
		return nil, err
	}

	gateway := &grpcGatewayServerImpl{
		cfg:              cfg,
		core:             packageReceiver,
		metricsCollector: metricsCollector,
		brokerConsumer:   broker,
	}

	router := newRouter(gwmux, cfg.Gateway, gateway.CardReceiveV1Handler, metricsCollector)

	interceptors := grpc.ChainUnaryInterceptor(
		alogger.UnaryTraceIdInterceptor,
		middleware.MetricInterceptor(metricsCollector),
		middleware.AuthInterceptor(cfg.Gateway.AuthToken),
	)

	grpcServer := grpc.NewServer(
		interceptors,
	)

	reflection.Register(grpcServer)

	gateway.httpServer = router
	gateway.grpcServer = grpcServer

	desc.RegisterCardReceiverServer(grpcServer, gateway)

	err = gateway.update(ctx)
	if err != nil {
		return nil, err
	}
	gateway.autoupdate(ctx, updateIntervalDefault)

	return gateway, nil
}

func (gw *grpcGatewayServerImpl) Start(ctx context.Context) error {
	// Запуск брокера потребителя.
	err := gw.makeBrokerSubscribers(ctx)
	if err != nil {
		return fmt.Errorf("ошибка при запуске брокера потребителя для создания пакетов: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)

	// GRPC
	g.Go(func() error {
		adr := gw.cfg.Gateway.GRPC.Host + ":" + gw.cfg.Gateway.GRPC.Port

		grpcListener, err := reuseport.Listen("tcp4", adr)
		if err != nil {
			return fmt.Errorf("ошибка при запуске GRPC сервера: %w", err)
		}
		alogger.InfoFromCtx(ctx, "запуск GRPC сервера на "+adr, nil, nil, false)
		defer alogger.InfoFromCtx(ctx, "GRPC сервер остановлен", nil, nil, false)

		err = gw.grpcServer.Serve(grpcListener)
		if err != nil {
			return fmt.Errorf("ошибка при запуске GRPC сервера: %w", err)
		}
		return nil
	})

	// REST
	g.Go(func() error {
		adr := gw.cfg.Gateway.HTTP.Host + ":" + gw.cfg.Gateway.HTTP.Port

		alogger.InfoFromCtx(ctx, "запуск HTTP сервера на "+adr, nil, nil, false)
		defer alogger.InfoFromCtx(ctx, "HTTP сервер остановлен", nil, nil, false)

		err := gw.httpServer.Listen(adr)
		if err != nil {
			return fmt.Errorf("ошибка при запуске http сервера: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		gw.grpcServer.GracefulStop()
		err := gw.httpServer.Shutdown()

		return err
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("ошибка: %w", err)
	}

	return nil
}

func (gw *grpcGatewayServerImpl) gracefulStop() {
	gw.grpcServer.GracefulStop()
	_ = gw.httpServer.Shutdown()

	gracefulStopWaitMillisecond := 100
	time.Sleep(time.Millisecond * time.Duration(gracefulStopWaitMillisecond))
}
