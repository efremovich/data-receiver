package usecases

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/operatorrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/tprepo"
	"github.com/efremovich/data-receiver/internal/usecases/webapi/storage"
	"github.com/efremovich/data-receiver/pkg/broker"
	"github.com/efremovich/data-receiver/pkg/metrics"

	aerror "github.com/efremovich/data-receiver/pkg/aerror"
)

type ReceiverCoreService interface {
	ReceivePackage(ctx context.Context, tpName string, tpBytes []byte, responseURL string) ([]byte, aerror.AError) // обработка нового водящего пакета

	GetTP(ctx context.Context, tp string, document string) (*entity.TransportPackage, []entity.TpEvent, []*entity.TpDirectory, error)

	PingDB(ctx context.Context) error
	PingStorage(ctx context.Context) error
	PingNATS(_ context.Context) error
	PingOperator(ctx context.Context) error
}

type receiverCoreServiceImpl struct {
	tpRepo           tprepo.TransportPackageRepo
	brokerNats       broker.NATS
	operatorRepo     operatorrepo.OperatorRepo
	storage          storage.Storage
	metricsCollector metrics.Collector
}

func NewPackageReceiverService(opR operatorrepo.OperatorRepo, tpR tprepo.TransportPackageRepo,
	nats broker.NATS, storage storage.Storage, metricsCollector metrics.Collector) ReceiverCoreService {
	service := receiverCoreServiceImpl{
		operatorRepo:     opR,
		brokerNats:       nats,
		storage:          storage,
		tpRepo:           tpR,
		metricsCollector: metricsCollector,
	}

	return &service
}

func (s *receiverCoreServiceImpl) GetTP(ctx context.Context, tp string, document string) (*entity.TransportPackage, []entity.TpEvent, []*entity.TpDirectory, error) {
	var (
		err   error
		tpRes *entity.TransportPackage
	)

	if tp != "" {
		if !strings.Contains(tp, ".cms") {
			tp += ".cms"
		}

		tpRes, err = s.tpRepo.SelectByName(ctx, tp)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil, err
		}
	}

	if tpRes == nil && document != "" {
		for _, ext := range []string{"", ".bin", ".xml", ".p7s"} {
			tpRes, err = s.tpRepo.SelectByDocument(ctx, document+ext)
			if err == nil {
				break
			}

			if !errors.Is(err, sql.ErrNoRows) {
				return nil, nil, nil, err
			}
		}
	}

	if tpRes == nil {
		return nil, nil, nil, nil
	}

	events, err := s.tpRepo.SelectEvents(tpRes.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	fileStructure, err := s.tpRepo.SelectFileStructure(ctx, tpRes.ID)
	if err != nil {
		return nil, nil, nil, err
	}

	return tpRes, events, fileStructure, nil
}

func (s *receiverCoreServiceImpl) PingDB(ctx context.Context) error {
	return s.tpRepo.Ping(ctx)
}

func (s *receiverCoreServiceImpl) PingStorage(ctx context.Context) error {
	return s.storage.Ping(ctx)
}

func (s *receiverCoreServiceImpl) PingNATS(_ context.Context) error {
	return s.brokerNats.Ping()
}

func (s *receiverCoreServiceImpl) PingOperator(ctx context.Context) error {
	return s.operatorRepo.Ping(ctx)
}
