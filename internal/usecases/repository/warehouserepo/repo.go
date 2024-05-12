package warehouserepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type WarehouseRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Warehouse, error)
	SelectBySellerID(ctx context.Context, sellerID int64) ([]*entity.Warehouse, error)
	Insert(ctx context.Context, in entity.Warehouse) (*entity.Warehouse, error)
	UpdateExecOne(ctx context.Context, in entity.Warehouse) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) WarehouseRepo
}

type repoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewWarehouseRepo(_ context.Context, db *postgresdb.DBConnection) (WarehouseRepo, error) {
	return &repoImpl{db: db}, nil
}

func (repo *repoImpl) SelectByID(ctx context.Context, id int64) (*entity.Warehouse, error) {
	var result warehouseDB

	query := "SELECT * FROM warehouses WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}
	return result.convertToEntityWarehouse(ctx), nil
}

func (repo *repoImpl) SelectBySellerID(ctx context.Context, sellerID int64) ([]*entity.Warehouse, error) {
	var result []warehouseDB

	query := "SELECT * FROM warehouses WHERE seller_id = $1"

	err := repo.getReadConnection().Select(&result, query, sellerID)
	if err != nil {
		return nil, err
	}

	var resEntity []*entity.Warehouse
	for _, v := range result {
		resEntity = append(resEntity, v.convertToEntityWarehouse(ctx))
	}
	return resEntity, nil
}

func (repo *repoImpl) Insert(ctx context.Context, in entity.Warehouse) (*entity.Warehouse, error) {
	dbModel := convertToDBWarehouse(ctx, in)

	query := `INSERT INTO warehouses (ext_id, title, address, type, seller_id) 
            VALUES ($1, $2, $3, $4, $5) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, dbModel.ExtID, dbModel.Title, dbModel.Address, dbModel.Type, dbModel.SellerID)
	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, in entity.Warehouse) error {
	dbModel := convertToDBWarehouse(ctx, in)

	query := `UPDATE warehouses SET 
            ext_id = $1, title = $2, address = $3, type = $4, seller_id = $5 
            WHERE id = $6`
	_, err := repo.getWriteConnection().ExecOne(query, dbModel.ExtID, dbModel.Title, dbModel.Address, dbModel.Type, dbModel.SellerID, dbModel.ID)
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

func (repo *repoImpl) WithTx(tx *postgresdb.Transaction) WarehouseRepo {
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
