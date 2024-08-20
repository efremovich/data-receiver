package statusrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type StatusRepo interface {
	SelectByName(ctx context.Context, name string) (*entity.Status, error)
	Insert(ctx context.Context, in entity.Status) (*entity.Status, error)
	UpdateExecOne(ctx context.Context, in entity.Status) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) StatusRepo
}

type repoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewStatusRepo(_ context.Context, db *postgresdb.DBConnection) (StatusRepo, error) {
	return &repoImpl{db: db}, nil
}

func (repo *repoImpl) SelectByName(ctx context.Context, name string) (*entity.Status, error) {
	var result statusDB
	query := "SELECT * FROM shop.statuses WHERE name = $1"

	err := repo.getReadConnection().Get(&result, query, name)
	if err != nil {
		return nil, err
	}
	return result.convertToEntityStatus(ctx), nil
}

func (repo *repoImpl) Insert(ctx context.Context, in entity.Status) (*entity.Status, error) {
	dbModel := convertToDBStatus(ctx, in)
	query := "INSERT INTO shop.statuses (name) VALUES ($1) RETURNING id"
	charIDWrap := repository.IDWrapper{}
	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, dbModel.Name)

	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, in entity.Status) error {
	dbModel := convertToDBStatus(ctx, in)
	query := "UPDATE shop.statuses SET name = $1 WHERE id = $2"
	_, err := repo.getWriteConnection().Exec(query, dbModel.Name, in.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *repoImpl) Ping(ctx context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *repoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *repoImpl) WithTx(tx *postgresdb.Transaction) StatusRepo {
	return &repoImpl{db: repo.db, tx: tx}
}

func (repo *repoImpl) getReadConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}
	return repo.db.GetReadConnection()
}

func (repo *repoImpl) getWriteConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}
	return repo.db.GetWriteConnection()
}
