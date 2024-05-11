package warehouserepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type WarehouseRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Warehouse, error)
	SelectByCardID(ctx context.Context, CardID int64) ([]*entity.Warehouse, error)
	SelectByPriceID(ctx context.Context, CardID int64) ([]*entity.Warehouse, error)
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

func (repo *repoImpl) SelectByCardID(ctx context.Context, cardID int64) ([]*entity.Warehouse, error) {
	var result []warehouseDB

	query := "SELECT * FROM warehouses WHERE card_id = $1"

	err := repo.getReadConnection().Select(&result, query, cardID)
	if err != nil {
		return nil, err
	}

	var resEntity []*entity.Warehouse
	for _, v := range result {
		resEntity = append(resEntity, v.convertToEntityWarehouse(ctx))
	}
	return resEntity, nil
}

func (repo *repoImpl) SelectByPriceID(ctx context.Context, cardID int64) ([]*entity.Warehouse, error) {
	var result []warehouseDB

	query := "SELECT * FROM warehouses WHERE price_id = $1"

	err := repo.getReadConnection().Select(&result, query, cardID)
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
	query := `INSERT INTO warehouses (title, tech_warehouse, card_id, price_id) 
            VALUES ($1, $2, $3, $4) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, in.Title, in.TechWarehouse, in.CardID, in.PriceID)
	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, in entity.Warehouse) error {
	dbModel := convertToDBWarehouse(ctx, in)

	query := `UPDATE warehouses SET title = $1, tech_warehouse = $2, card_id = $3, price_id = $4 WHERE id = $5`
	_, err := repo.getWriteConnection().ExecOne(query, dbModel.Title, dbModel.TechWarehouse, dbModel.CardID, dbModel.PriceID, dbModel.ID)
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
