package warehousetyperepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type WarehouseTypeRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.WarehouseType, error)
	SelectByTitle(ctx context.Context, title string) (*entity.WarehouseType, error)
	Insert(ctx context.Context, in entity.WarehouseType) (*entity.WarehouseType, error)
	UpdateExecOne(ctx context.Context, in entity.WarehouseType) error
	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) WarehouseTypeRepo
}

type repoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewWarehouseTypeRepo(_ context.Context, db *postgresdb.DBConnection) (WarehouseTypeRepo, error) {
	return &repoImpl{db: db}, nil
}

func (repo *repoImpl) SelectByID(ctx context.Context, id int64) (*entity.WarehouseType, error) {
	var result warehousetypeDB
	query := "SELECT * FROM shop.warehouse_types WHERE id = $1"
	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}
	return result.convertToEntityWarehousetype(ctx), nil
}

func (repo *repoImpl) SelectByTitle(ctx context.Context, title string) (*entity.WarehouseType, error) {
	var result warehousetypeDB
	query := "SELECT * FROM shop.warehouse_types WHERE name = $1"
	err := repo.getReadConnection().Get(&result, query, title)
	if err != nil {
		return nil, err
	}
	return result.convertToEntityWarehousetype(ctx), nil
}

func (repo *repoImpl) Insert(ctx context.Context, in entity.WarehouseType) (*entity.WarehouseType, error) {
	query := "INSERT INTO shop.warehouse_types(id, name) VALUES ($1, $2) RETURNING id"
	charIDWrap := repository.IDWrapper{}
	err := repo.getWriteConnection().Get(&charIDWrap, query, in.ID, in.Title)
	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, in entity.WarehouseType) error {
	query := "UPDATE shop.warehouse_types SET name = $2 WHERE id = $1"
	_, err := repo.getWriteConnection().Exec(query, in.ID, in.Title)
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

func (repo *repoImpl) WithTx(tx *postgresdb.Transaction) WarehouseTypeRepo {
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
