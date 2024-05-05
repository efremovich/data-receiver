package charrepo

import (
	"context"
	"crypto/rand"
	"database/sql"
	"math/big"
	"sync"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func NewCharRepoMock() CharRepo {
	return charRepoMockImpl{
		store: make(map[int64]entity.Characteristic),
		m:     &sync.Mutex{},
	}
}

type charRepoMockImpl struct {
	store map[int64]entity.Characteristic
	m     *sync.Mutex
}

func (repo charRepoMockImpl) SelectByID(_ context.Context, id int64) (*entity.Characteristic, error) {
	repo.m.Lock()
	defer repo.m.Unlock()

	res, ok := repo.store[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return &res, nil
}

func (repo charRepoMockImpl) SelectByCardID(_ context.Context, cardID int64) (*entity.Characteristic, error) {
	repo.m.Lock()
	defer repo.m.Unlock()

	for _, v := range repo.store {
		if v.CardID == cardID {
			return &v, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (repo charRepoMockImpl) Insert(_ context.Context, in entity.Characteristic) (*entity.Characteristic, error) {
	repo.m.Lock()
	defer repo.m.Unlock()

	maxID := 101
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(maxID)))
	id := r.Int64()

	newChar := entity.Characteristic{
		ID:     id,
		Title:  in.Title,
		Value:  in.Value,
		CardID: in.CardID,
	}
	repo.store[id] = newChar

	return &newChar, nil
}

func (repo charRepoMockImpl) UpdateExecOne(_ context.Context, in entity.Characteristic) error {
	repo.m.Lock()
	defer repo.m.Unlock()

	res, ok := repo.store[in.ID]
	if !ok {
		return sql.ErrNoRows
	}
	res.Title = in.Title
	res.Value = in.Value
	res.CardID = in.CardID

	repo.store[res.ID] = res

	return nil
}

func (repo charRepoMockImpl) Ping(_ context.Context) error {
	repo.m.Lock()
	defer repo.m.Unlock()

	return nil
}

func (repo charRepoMockImpl) BeginTX(_ context.Context) (postgresdb.Transaction, error) {
	return postgresdb.TransactionMock{}, nil
}

func (repo charRepoMockImpl) WithTx(_ *postgresdb.Transaction) CharRepo {
	return repo
}
