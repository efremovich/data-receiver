package warehousetyperepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

var ErrObjectNotFound = entity.ErrObjectNotFound

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
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по id %d в таблице warehouse_type: %w", id, err)
	}

	return result.convertToEntityWarehousetype(ctx), nil
}

func (repo *repoImpl) SelectByTitle(ctx context.Context, title string) (*entity.WarehouseType, error) {
	var result warehousetypeDB

	query := "SELECT * FROM shop.warehouse_types WHERE name = $1"

	err := repo.getReadConnection().Get(&result, query, title)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по имени %s в таблице warehouse_type: %w", title, err)
	}

	return result.convertToEntityWarehousetype(ctx), nil
}

func (repo *repoImpl) Insert(_ context.Context, in entity.WarehouseType) (*entity.WarehouseType, error) {
	charIDWrap := repository.IDWrapper{}

	query := "INSERT INTO shop.warehouse_types(id, name) VALUES ($1, $2) RETURNING id"

	err := repo.getWriteConnection().Get(&charIDWrap, query, in.ID, in.Title)
	if err != nil {
		return nil, fmt.Errorf("ошибка при вставке данных %s в таблицу brands: %w", in.Title, err)
	}

	in.ID = charIDWrap.ID.Int64

	return &in, nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, in entity.WarehouseType) error {
	dbModel := convertToDBWarehousetype(ctx, in)

	query := "UPDATE shop.warehouse_types SET name = $1 WHERE id = $2"

	_, err := repo.getWriteConnection().Exec(query, in.Title, dbModel.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении данных %s в таблицу brands: %w", in.Title, err)
	}

	return nil
}

func (repo *repoImpl) Ping(_ context.Context) error {
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
