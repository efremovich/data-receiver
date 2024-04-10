package operatorrepo

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/webapi/operatorfetcher"
	"github.com/efremovich/data-receiver/pkg/logger"
)

const updateIntervalDefault = time.Minute * 10

type OperatorRepo interface {
	// GetOperators возвращает список операторов
	GetOperators(ctx context.Context) ([]entity.Operator, error)
	// GetOperatorsMap возвращает операторов мапой, где ключ thumb
	GetOperatorsMap(ctx context.Context) (map[string]entity.Operator, error)
	// Ping вернет nil, если репозиторий доступен, либо ошибку
	Ping(ctx context.Context) error
}

type operatorRepoImpl struct {
	client       operatorfetcher.OperatorFetcher
	operators    []entity.Operator
	operatorsMap map[string]entity.Operator
	mu           sync.RWMutex
}

func NewOperatorRepo(ctx context.Context, opclient operatorfetcher.OperatorFetcher) (OperatorRepo, error) {
	repo := &operatorRepoImpl{
		client:       opclient,
		operators:    nil,
		mu:           sync.RWMutex{},
		operatorsMap: make(map[string]entity.Operator),
	}

	// при создании подтянем список один раз, чтобы удостовериться,
	// что он точно не пустой
	err := repo.update(ctx)

	// запустим фоновое обновление
	go repo.autoupdate(ctx, updateIntervalDefault)

	return repo, err
}

func (r *operatorRepoImpl) GetOperators(_ context.Context) ([]entity.Operator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ops := r.operators

	return ops, nil
}

func (r *operatorRepoImpl) GetOperatorsMap(_ context.Context) (map[string]entity.Operator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := r.operatorsMap

	return res, nil
}

func (r *operatorRepoImpl) Ping(ctx context.Context) error {
	return r.client.Ping(ctx)
}

// Autoupdate обновляет список автоматически по тикеру с интервалом upd.
func (r *operatorRepoImpl) autoupdate(ctx context.Context, upd time.Duration) {
	t := time.NewTicker(upd)

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			err := r.update(ctx)
			if err != nil {
				logger.GetLoggerFromContext(ctx).Errorf("ошибка обновления списка операторов: %s", err.Error())
			}
		}
	}
}

func (r *operatorRepoImpl) update(ctx context.Context) error {
	ops, err := r.client.GetOperatorList(ctx)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.operators = ops

	for _, o := range ops {
		for _, thumb := range o.Thumbs {
			r.operatorsMap[strings.ToLower(thumb)] = o
		}
	}

	return nil
}
