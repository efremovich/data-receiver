package stockrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

var ErrObjectNotFound = entity.ErrObjectNotFound

type StockRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Stock, error)
	SelectByBarcode(ctx context.Context, barcodeID int64, dateFrom time.Time) (*entity.Stock, error)
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

	query := "SELECT * FROM shop.stocks WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по id %d в таблице stocks: %w", id, err)
	}

	return result.convertToEntityStock(ctx), nil
}

func (repo *repoImpl) SelectByBarcode(ctx context.Context, barcodeID int64, dateFrom time.Time) (*entity.Stock, error) {
	var result stockDB

	query := "SELECT * FROM shop.stocks WHERE barcode_id = $1 and created_at = $2"

	err := repo.getReadConnection().Get(&result, query, barcodeID, dateFrom.Format("2006-01-02 00:00:00"))
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по id %d и дате %s в таблице stocks: %w", barcodeID, dateFrom.Format("02.01.2006"), err)
	}

	return result.convertToEntityStock(ctx), nil
}

func (repo *repoImpl) SelectBySellerID(ctx context.Context, sellerID int64) ([]*entity.Stock, error) {
	var result []stockDB

	query := "SELECT * FROM shop.stocks WHERE seller_id = $1"

	err := repo.getReadConnection().Select(&result, query, sellerID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по продавцу %d в таблице stocks: %w", sellerID, err)
	}

	var resEntity []*entity.Stock
	for _, v := range result {
		resEntity = append(resEntity, v.convertToEntityStock(ctx))
	}

	return resEntity, nil
}

func (repo *repoImpl) Insert(ctx context.Context, in entity.Stock) (*entity.Stock, error) {
	dbModel := convertToDBStock(ctx, in)

	query := `INSERT INTO shop.stocks (
              quantity, 
              barcode_id,
              warehouse_id,
              card_id,
              created_at
            ) 
            VALUES ($1, $2, $3, $4, $5) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query,
		dbModel.Quantity,
		dbModel.BarcodeID,
		dbModel.WarehouseID,
		dbModel.CardID,
		dbModel.CreatedAt.Time.Format("2006-01-02 00:00:00"),
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка при вставке данных в таблицу stocks: %w", err)
	}

	in.ID = charIDWrap.ID.Int64

	return &in, nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, in entity.Stock) error {
	dbModel := convertToDBStock(ctx, in)

	query := `UPDATE shop.stocks SET 
              quantity = $1, 
              barcode_id = $2,
              warehouse_id = $3,
              card_id = $4,
              updated_at = now()
            WHERE id = $5`
	_, err := repo.getWriteConnection().ExecOne(query,
		dbModel.Quantity,
		dbModel.BarcodeID,
		dbModel.WarehouseID,
		dbModel.CardID,
		dbModel.ID)

	if err != nil {
		return fmt.Errorf("ошибка при обновлении данных в таблицу stocks: %w", err)
	}

	return nil
}

func (repo *repoImpl) Ping(_ context.Context) error {
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
