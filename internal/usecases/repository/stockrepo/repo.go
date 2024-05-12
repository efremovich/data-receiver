package stockrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type StockRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Stock, error)
	SelectBySellerID(ctx context.Context, sellerID int64) ([]*entity.Stock, error)
	Insert(ctx context.Context, in entity.Stock) (*entity.Stock, error)
	UpdateExecOne(ctx context.Context, in entity.Stock) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) StockRepo
}

type repoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewStockRepo(_ context.Context, db *postgresdb.DBConnection) (StockRepo, error) {
	return &repoImpl{db: db}, nil
}

func (repo *repoImpl) SelectByID(ctx context.Context, id int64) (*entity.Stock, error) {
	var result stockDB

	query := "SELECT * FROM stocks WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}
	return result.convertToEntityStock(ctx), nil
}

func (repo *repoImpl) SelectBySellerID(ctx context.Context, sellerID int64) ([]*entity.Stock, error) {
	var result []stockDB

	query := "SELECT * FROM stocks WHERE seller_id = $1"

	err := repo.getReadConnection().Select(&result, query, sellerID)
	if err != nil {
		return nil, err
	}

	var resEntity []*entity.Stock
	for _, v := range result {
		resEntity = append(resEntity, v.convertToEntityStock(ctx))
	}
	return resEntity, nil
}

func (repo *repoImpl) Insert(ctx context.Context, in entity.Stock) (*entity.Stock, error) {
	dbModel := convertToDBStock(ctx, in)

	query := `INSERT INTO stocks (
              quantity, 
              in_way_to_client,
              in_way_from_client,
              size_id,
              barcode,
              warehouse_id,
              card_id,
              seller_id,
              created_at
            ) 
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now()) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query,
		dbModel.Quantity,
		dbModel.InWayToClient,
		dbModel.InWayFromClient,
		dbModel.SizeID,
		dbModel.Barcode,
		dbModel.WarehouseID,
		dbModel.CardID,
		dbModel.SellerID,
	)
	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, in entity.Stock) error {
	dbModel := convertToDBStock(ctx, in)

	query := `UPDATE stocks SET 
              quantity = $1, 
              in_way_to_client = $2,
              in_way_from_client = $3,
              size_id = $4,
              barcode = $5,
              warehouse_id = $6,
              card_id = $7,
              seller_id = $8,
              updated_at = now()
            WHERE id = $9`
	_, err := repo.getWriteConnection().ExecOne(query,
		dbModel.Quantity,
		dbModel.InWayToClient,
		dbModel.InWayFromClient,
		dbModel.SizeID,
		dbModel.Barcode,
		dbModel.WarehouseID,
		dbModel.CardID,
		dbModel.SellerID,
		dbModel.ID)
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

func (repo *repoImpl) WithTx(tx *postgresdb.Transaction) StockRepo {
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
