package pvzrepo

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

type PvzRepo interface {
	SelectByOfficeID(ctx context.Context, officeID int) (*entity.Pvz, error)
	Insert(ctx context.Context, in entity.Pvz) (*entity.Pvz, error)
	UpdateExecOne(ctx context.Context, in *entity.Pvz) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) PvzRepo
}

type repoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewPvzRepo(_ context.Context, db *postgresdb.DBConnection) (PvzRepo, error) {
	return &repoImpl{db: db}, nil
}

func (repo *repoImpl) SelectByOfficeID(ctx context.Context, officeID int) (*entity.Pvz, error) {
	var result pvzDB

	query := `
  SELECT
      id,
      office_name,
      office_id,
      supplier_name,
      supplier_id,
      supplier_inn
    FROM shop.pvzs
    WHERE office_id = $1
  `
	err := repo.getReadConnection().Get(&result, query, officeID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по ID %d в таблице pvz: %w", officeID, err)
	}

	return result.convertToEntityPvz(ctx), nil
}

func (repo *repoImpl) Insert(ctx context.Context, pvz entity.Pvz) (*entity.Pvz, error) {
	dbModel := convertToEntityPvz(ctx, &pvz)
	query := `
  INSERT INTO shop.pvzs (
    office_name,
    office_id,
    supplier_name,
    supplier_id,
    supplier_inn
  ) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
  ) RETURNING id
  `
	pvzIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&pvzIDWrap, query, dbModel.OfficeName, dbModel.OfficeID, dbModel.SupplierName, dbModel.SupplierID, dbModel.SupplierINN)
	if err != nil {
		return nil, fmt.Errorf("ошибка вставки данных в таблицу pvz %s в таблице pvz %w", pvz.OfficeName, err)
	}
	pvz.ID = pvzIDWrap.ID.Int64

	return &pvz, nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, pvz *entity.Pvz) error {
	dbModel := convertToEntityPvz(ctx, pvz)
	query := `
  UPDATE pvz
  SET office_name = $1,
      office_id = $2,
      supplier_name = $3,
      supplier_id = $4,
      supplier_inn = $5
  WHERE id = $6
  `
	_, err := repo.getWriteConnection().Exec(query, dbModel.OfficeName, dbModel.OfficeID, dbModel.SupplierName, dbModel.SupplierID, dbModel.SupplierINN, pvz.ID)
	return err
}

func (repo *repoImpl) Ping(_ context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *repoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *repoImpl) WithTx(tx *postgresdb.Transaction) PvzRepo {
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
