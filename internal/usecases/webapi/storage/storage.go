package storage

import (
	"bytes"
	"context"
	"time"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/pkg/metrics"
	ggc "github.com/efremovich/data-receiver/pkg/protos/generic_client"
	"github.com/efremovich/data-receiver/pkg/protos/storage"
)

type Storage interface {
	SaveFile(ctx context.Context, fileName string, data []byte) error
	GetFile(ctx context.Context, fileName string) ([]byte, error)
	Ping(ctx context.Context) error
}

func NewStorageClient(_ context.Context, cfg config.Storage, serviceName string, metricsCollector metrics.Collector) (Storage, error) {
	if cfg.UseMockStorage {
		return NewMockStorageClient(cfg.URL)
	}

	timeout := time.Second * time.Duration(cfg.TimeoutSeconds)

	c := ggc.Config{
		Addr:               cfg.URL,
		Token:              cfg.Token,
		Timeout:            timeout,
		InsecureSkipVerify: true,
	}

	storageClient, err := storage.NewConn(serviceName, c)
	if err != nil {
		return nil, err
	}

	return &storageImpl{
		storageClient:    storageClient,
		metricsCollector: metricsCollector,
	}, nil
}

type storageImpl struct {
	storageClient *storage.Conn

	metricsCollector metrics.Collector
}

func (s *storageImpl) SaveFile(_ context.Context, fileName string, data []byte) error {
	startTime := time.Now()
	defer func() { s.metricsCollector.AddSaveStorageTime(time.Since(startTime)) }()

	in := storage.SaveFileInput{
		CustomId: fileName,
		File:     bytes.NewBuffer(data),
	}

	_, err := s.storageClient.SaveFile(in)
	if err != nil {
		return err
	}

	return nil
}

func (s *storageImpl) GetFile(_ context.Context, fileName string) ([]byte, error) {
	in := storage.GetFileInput{
		CustomId: fileName,
	}

	res, err := s.storageClient.GetFile(in)
	if err != nil {
		return nil, err
	}

	return res.File.Bytes(), nil
}

func (s *storageImpl) Ping(ctx context.Context) error {
	return s.storageClient.Ping(ctx)
}
