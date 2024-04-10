package tprepo

import (
	"context"
	"crypto/rand"
	"database/sql"
	"math/big"
	"sync"
	"time"

	postgres "github.com/efremovich/data-receiver/pkg/postgresdb"

	"github.com/efremovich/data-receiver/internal/entity"
)

func NewTpRepoMock() TransportPackageRepo {
	return tpRepoMockImpl{
		store:     make(map[int64]entity.TransportPackage),
		events:    make(map[int64][]entity.TpEvent),
		storeDirs: make(map[int64][]*entity.TpDirectory),
		m:         &sync.Mutex{}}
}

type tpRepoMockImpl struct {
	store     map[int64]entity.TransportPackage
	storeDirs map[int64][]*entity.TpDirectory
	events    map[int64][]entity.TpEvent
	m         *sync.Mutex
}

func (repo tpRepoMockImpl) SelectByID(_ context.Context, id int64) (*entity.TransportPackage, error) {
	repo.m.Lock()
	defer repo.m.Unlock()

	res, ok := repo.store[id]
	if !ok {
		return nil, sql.ErrNoRows
	}

	return &res, nil
}

func (repo tpRepoMockImpl) SelectByName(_ context.Context, name string) (*entity.TransportPackage, error) {
	for _, v := range repo.store {
		if v.Name == name {
			return &v, nil
		}
	}

	return nil, sql.ErrNoRows
}

func (repo tpRepoMockImpl) SelectByDocument(_ context.Context, doc string) (*entity.TransportPackage, error) {
	var (
		tpID int64
		ok   bool
	)

	for id, content := range repo.storeDirs {
		for _, dir := range content {
			_, ok = dir.Files[doc]
			if ok {
				tpID = id
			}
		}
	}

	if tpID == 0 {
		return nil, sql.ErrNoRows
	}

	tp, ok := repo.store[tpID]
	if !ok {
		return nil, sql.ErrNoRows
	}

	return &tp, nil
}

func (repo tpRepoMockImpl) Insert(_ context.Context, name string, receiptURL string) (*entity.TransportPackage, error) {
	repo.m.Lock()
	defer repo.m.Unlock()

	maxID := 101
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(maxID)))
	id := r.Int64()

	newTP := entity.TransportPackage{
		ID:         id,
		Name:       name,
		ReceiptURL: receiptURL,
		Status:     entity.TpStatusEnumNew,
	}
	repo.store[id] = newTP

	return &newTP, nil
}

func (repo tpRepoMockImpl) UpdateExecOne(_ context.Context, tp entity.TransportPackage) error {
	repo.m.Lock()
	defer repo.m.Unlock()

	res, ok := repo.store[tp.ID]
	if !ok {
		return sql.ErrNoRows
	}

	res.Status = tp.Status
	res.ErrorCode = tp.ErrorCode
	res.ErrorText = tp.ErrorText
	res.Origin = tp.Origin
	res.IsReceipt = tp.IsReceipt
	repo.store[res.ID] = res

	return nil
}

func (repo tpRepoMockImpl) SaveFileStructure(_ context.Context, tpID int64, fileStructure []*entity.TpDirectory) error {
	repo.m.Lock()
	defer repo.m.Unlock()

	repo.storeDirs[tpID] = fileStructure

	return nil
}

func (repo tpRepoMockImpl) SelectFileStructure(_ context.Context, tpID int64) ([]*entity.TpDirectory, error) {
	repo.m.Lock()
	defer repo.m.Unlock()

	res := repo.storeDirs[tpID]

	return res, nil
}

func (repo tpRepoMockImpl) AddNewEvent(tpID int64, event entity.TpEventTypeEnum, desc string) error {
	repo.m.Lock()
	defer repo.m.Unlock()

	_, ok := repo.events[tpID]
	if !ok {
		repo.events[tpID] = []entity.TpEvent{}
	}

	repo.events[tpID] = append(repo.events[tpID], entity.TpEvent{TpID: tpID, Description: desc, EventType: event, CreatedAt: time.Now()})

	return nil
}

func (repo tpRepoMockImpl) SelectEvents(tpID int64) ([]entity.TpEvent, error) {
	repo.m.Lock()
	defer repo.m.Unlock()

	res := repo.events[tpID]

	return res, nil
}

func (repo tpRepoMockImpl) Ping(_ context.Context) error {
	repo.m.Lock()
	defer repo.m.Unlock()

	return nil
}

func (repo tpRepoMockImpl) BeginTX(_ context.Context) (postgres.Transaction, error) {
	return postgres.TransactionMock{}, nil
}

func (repo tpRepoMockImpl) WithTx(_ *postgres.Transaction) TransportPackageRepo {
	return repo
}
