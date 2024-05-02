package cardrepo

import (
	"context"
	"crypto/rand"
	"database/sql"
	"math/big"
	"sync"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func NewCardRepoMock() CardRepo {
	return cardRepoMockImpl{m: &sync.Mutex{}}
}

type cardRepoMockImpl struct {
	store map[int64]entity.Card
	m     *sync.Mutex
}

func (repo cardRepoMockImpl) SelectByID(_ context.Context, id int64) (*entity.Card, error) {
	repo.m.Lock()
	defer repo.m.Unlock()
	res, ok := repo.store[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return &res, nil
}

func (repo cardRepoMockImpl) SelectByVendorID(_ context.Context, vendorID string) (*entity.Card, error) {
	for _, v := range repo.store {
		if v.VendorID == vendorID {
			return &v, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (repo cardRepoMockImpl) SelectByTitle(_ context.Context, title string) (*entity.Card, error) {
	for _, v := range repo.store {
		if v.Title == title {
			return &v, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (repo cardRepoMockImpl) SelectByVendorCode(_ context.Context, vendorCode string) (*entity.Card, error) {
	for _, v := range repo.store {
		if v.VendorCode == vendorCode {
			return &v, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (repo cardRepoMockImpl) Insert(_ context.Context, in entity.Card) (*entity.Card, error) {
	repo.m.Lock()
	defer repo.m.Unlock()

	maxID := 101
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(maxID)))
	id := r.Int64()

	newCard := entity.Card{
		ID:          id,
		VendorID:    in.VendorID,
		VendorCode:  in.VendorCode,
		Title:       in.Title,
		Description: in.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	repo.store[id] = newCard

	return &newCard, nil
}

func (repo cardRepoMockImpl) UpdateExecOne(_ context.Context, card entity.Card) error {
	repo.m.Lock()
	defer repo.m.Unlock()

	res, ok := repo.store[card.ID]
	if !ok {
		return sql.ErrNoRows
	}
	res.Title = card.Title
	res.VendorCode = card.VendorCode
	res.VendorID = card.VendorID
	res.Description = card.Description

	repo.store[res.ID] = res
	return nil
}

func (repo cardRepoMockImpl) Ping(_ context.Context) error {
	repo.m.Lock()
	defer repo.m.Unlock()

	return nil
}

func (repo cardRepoMockImpl) BeginTX(_ context.Context) (postgresdb.Transaction, error) {
	return postgresdb.TransactionMock{}, nil
}

func (repo cardRepoMockImpl) WithTx(_ *postgresdb.Transaction) CardRepo {
	return repo
}
