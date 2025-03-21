package salerepo

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

type SaleRepo interface {
	SelectByExternalID(ctx context.Context, externalID string, date time.Time) (*entity.Sale, error)
	SelectByCardIDAndDate(ctx context.Context, cardID int64, date time.Time) (*entity.Sale, error)
	SelectByOrderID(ctx context.Context, orderID int64) (*entity.Sale, error)
	Insert(ctx context.Context, in entity.Sale) (*entity.Sale, error)
	UpdateExecOne(ctx context.Context, in *entity.Sale) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) SaleRepo
}

type repoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewSaleRepo(_ context.Context, db *postgresdb.DBConnection) (SaleRepo, error) {
	return &repoImpl{db: db}, nil
}

func (repo *repoImpl) SelectByExternalID(ctx context.Context, externalID string, date time.Time) (*entity.Sale, error) {
	var result saleDB

	query := `SELECT
              external_id,
              price,
              type,
              final_price,
              for_pay,
              quantity,
              created_at,
              order_id,
              seller_id,
              card_id,
              warehouse_id,
              region_id,
              price_size_id
            FROM shop.sales WHERE external_id = $1 and created_at = $2`

	err := repo.getReadConnection().Get(&result, query, externalID, date.Format("2006-01-02 15:04:05"))
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по ID %s в таблице sale: %w", externalID, err)
	}

	return result.convertToEntitySale(ctx), nil
}

func (repo *repoImpl) SelectByCardIDAndDate(ctx context.Context, cardID int64, date time.Time) (*entity.Sale, error) {
	var result saleDB

	query := `SELECT 
              external_id, 
              price,
              type,
              final_price,
              for_pay,
              quantity, 
              created_at, 
              order_id, 
              seller_id, 
              card_id, 
              warehouse_id, 
              region_id,
              price_size_id 
            FROM shop.sales WHERE card_id = $1 and created_at = $2`

	err := repo.getReadConnection().Get(&result, query, cardID, date.Format("2006-01-02 15:04:05"))
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по ID %d в таблице sale: %w", cardID, err)
	}

	return result.convertToEntitySale(ctx), nil
}

func (repo *repoImpl) SelectByOrderID(ctx context.Context, orderID int64) (*entity.Sale, error) {
	var result saleDB

	query := "SELECT * FROM shop.sales WHERE order_id = $1"

	err := repo.getReadConnection().Get(&result, query, orderID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по orderID %d в таблице sale: %w", orderID, err)
	}

	return result.convertToEntitySale(ctx), nil
}

func (repo *repoImpl) Insert(ctx context.Context, in entity.Sale) (*entity.Sale, error) {
	dbModel := convertToDBSale(ctx, &in)

	query := `INSERT INTO shop.sales(
              external_id,
              price,
              discount,
              type,
              final_price,
              for_pay,
              
              quantity,
              created_at,
              updated_at,
              order_id,
              seller_id,
              card_id, 
              warehouse_id,
              region_id,
              price_size_id
  )
 VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`

	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query,
		dbModel.ExternalID,
		dbModel.Price,
		dbModel.Discount,
		dbModel.Type,

		dbModel.FinalPrice,
		dbModel.ForPay,

		dbModel.Quantity,
		dbModel.CreatedAt,
		time.Now(),

		dbModel.OrderID,
		dbModel.SellerID,
		dbModel.CardID,
		dbModel.WarehouseID,
		dbModel.RegionID,
		dbModel.PriceSizeID,
	)

	if err != nil {
		return nil, err
	}

	return dbModel.convertToEntitySale(ctx), nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, in *entity.Sale) error {
	dbModel := convertToDBSale(ctx, in)
	query := `UPDATE shop.sales SET
              external_id = $1,
              price = $2,
              discount = $3,
              type = $4,
              final_price = $5,
              for_pay = $6,
              
              quantity = $7,
              updated_at = now()
              
              WHERE id = $8`

	_, err := repo.getWriteConnection().Exec(query,
		dbModel.ExternalID,
		dbModel.Price,
		dbModel.Discount,
		dbModel.Type,
		dbModel.FinalPrice,
		dbModel.ForPay,
		dbModel.Quantity,
		dbModel.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repoImpl) Ping(_ context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *repoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *repoImpl) WithTx(tx *postgresdb.Transaction) SaleRepo {
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
