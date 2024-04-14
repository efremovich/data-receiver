package usecases

import (
	"context"

	"github.com/efremovich/data-receiver/internal/usecases/repository/tprepo"
	"github.com/efremovich/data-receiver/pkg/broker"
	"github.com/efremovich/data-receiver/pkg/metrics"

	aerror "github.com/efremovich/data-receiver/pkg/aerror"
)

type ReceiverCoreService interface {
	ReceiveData(ctx context.Context) aerror.AError

	PingDB(ctx context.Context) error
	PingNATS(_ context.Context) error
}

type receiverCoreServiceImpl struct {
	tpRepo           tprepo.TransportPackageRepo
	brokerNats       broker.NATS
	metricsCollector metrics.Collector
}

func NewPackageReceiverService(tpR tprepo.TransportPackageRepo,nats broker.NATS, metricsCollector metrics.Collector) ReceiverCoreService {
	service := receiverCoreServiceImpl{
		brokerNats:       nats,
		metricsCollector: metricsCollector,
		tpRepo:           tpR,
	}

	return &service
}

func (s *receiverCoreServiceImpl) PingNATS(_ context.Context) error {
	return s.brokerNats.Ping()
}

func (s *receiverCoreServiceImpl) PingDB(ctx context.Context) error {
	return s.tpRepo.Ping(ctx)
}
