package charrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type CharRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Characteristic, error)
	SelectByTitle(ctx context.Context, title string) (*entity.Characteristic, error)
	Insert(ctx context.Context, in entity.Characteristic) (*entity.Characteristic, error)
	UpdateExecOne(ctx context.Context, in entity.Characteristic) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) CharRepo
}

type charRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewCharRepo(_ context.Context, db *postgresdb.DBConnection) (CharRepo, error) {
	return &charRepoImpl{db: db}, nil
}

func (repo *charRepoImpl) SelectByID(ctx context.Context, id int64) (*entity.Characteristic, error) {
	var result characteristicDB

	query := "SELECT id, title FROM shop.characteristics WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}
	return result.ConvertToEntityCharacteristic(ctx), nil
}

func (repo *charRepoImpl) SelectByTitle(ctx context.Context, title string) (*entity.Characteristic, error) {
	var result characteristicDB

	query := "SELECT id, title FROM shop.characteristics WHERE title = $1"

	err := repo.getReadConnection().Get(&result, query, title)
	if err != nil {
		return nil, err
	}
	return result.ConvertToEntityCharacteristic(ctx), nil
}

func (repo *charRepoImpl) Insert(ctx context.Context, in entity.Characteristic) (*entity.Characteristic, error) {
	query := `INSERT INTO shop.characteristics (title) 
            VALUES ($1) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, in.Title)
	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *charRepoImpl) UpdateExecOne(ctx context.Context, in entity.Characteristic) error {
	dbModel := convertToDBCharacteristic(ctx, in)

	query := `UPDATE shop.characteristics SET title = $1 WHERE id = $2`
	_, err := repo.getWriteConnection().ExecOne(query, dbModel.Title, dbModel.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *charRepoImpl) Ping(ctx context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *charRepoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *charRepoImpl) WithTx(tx *postgresdb.Transaction) CharRepo {
	return &charRepoImpl{db: repo.db, tx: tx}
}

func (repo *charRepoImpl) getReadConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetReadConnection()
}

func (repo *charRepoImpl) getWriteConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetWriteConnection()
}
