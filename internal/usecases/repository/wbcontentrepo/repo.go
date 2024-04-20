package wbcontentrepo

import (
	"context"
	"sync"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/webapi/wbfetcher"
)

type WBContentRepo interface {
	GetCard(ctx context.Context) []entity.Card
	Ping(ctx context.Context) error
}

type wbContentRepoImpl struct {
	client wbfetcher.WildberriesFetcher
	cards  []entity.Card
	mu     sync.RWMutex
}

func NewWBContentRepo(ctx context.Context, client wbfetcher.WildberriesFetcher) (WBContentRepo, error) {
	repo := &wbContentRepoImpl{
		client: client,
		mu:     sync.RWMutex{},
	}

	err := repo.update(ctx)

	return repo, err
}

func (wb *wbContentRepoImpl) GetCard(ctx context.Context) []entity.Card {
	wb.mu.RLock()
	defer wb.mu.RUnlock()
	cards := wb.cards
	return cards
}

func (wb *wbContentRepoImpl) Ping(ctx context.Context) error {
	return wb.client.Ping(ctx)
}

func (wb *wbContentRepoImpl) update(ctx context.Context) error {
	cards, err := wb.client.GetCards(ctx)
	if err != nil {
		return err
	}

	wb.mu.Lock()
	defer wb.mu.Unlock()
	wb.cards = cards

	return nil
}
