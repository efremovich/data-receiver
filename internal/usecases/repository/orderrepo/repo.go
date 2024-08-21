package orderrepo

import (
	"context"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type OrderRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Order, error)
	SelectByCardIDAndDate(ctx context.Context, cardID int64, date time.Time) (*entity.Order, error)
	Insert(ctx context.Context, in entity.Order) (*entity.Order, error)
	UpdateExecOne(ctx context.Context, in entity.Order) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) OrderRepo
}

type repoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewOrderRepo(_ context.Context, db *postgresdb.DBConnection) (OrderRepo, error) {
	return &repoImpl{db: db}, nil
}

func (repo *repoImpl) SelectByID(ctx context.Context, id int64) (*entity.Order, error) {
	var result orderDB

	query := "SELECT * FROM shop.orders WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}

	return result.convertToEntityOrder(ctx), nil
}

func (repo *repoImpl) SelectByCardIDAndDate(ctx context.Context, cardID int64, date time.Time) (*entity.Order, error) {
	var result orderDB

	query := "SELECT * FROM shop.orders WHERE id = $1 and created_at = $2"

	err := repo.getReadConnection().Get(&result, query, cardID, date.Format("2006-01-02 00:00:00"))
	if err != nil {
		return nil, err
	}

	return result.convertToEntityOrder(ctx), nil
}

func (repo *repoImpl) Insert(ctx context.Context, in entity.Order) (*entity.Order, error) {
	dbModel := convertToDBOrder(ctx, in)

	query := `INSERT INTO shop.orders (
              external_id,
              price,
              status_id,
              direction,
              type,
              sale,
              created_at,
              seller_id, 
              card_id, 
              warehouse_id,
              region_id,
              price_size_id
            ) 
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query,
		dbModel.ExternalID,
		dbModel.Price,
		dbModel.StatusID,
		dbModel.Direction,
		dbModel.Type,
		dbModel.Sale,
		dbModel.CreatedAt,

		dbModel.SellerID,
		dbModel.CardID,
		dbModel.WarehouseID,
		dbModel.RegionID,
		dbModel.PriceSizeID,
	)
	if err != nil {
		return nil, err
	}

	in.ID = charIDWrap.ID.Int64

	return &in, nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, in entity.Order) error {
	dbModel := convertToDBOrder(ctx, in)

	query := `UPDATE shop.orders SET
              external_id= $1,
              price = $2,
              status_id = $3,
              direction = $4,
              type = $5,
              sale = $6,
             
              quantity = $7,
              created_at = $8,
              
              seller_id = $9, 
              card_id = $10, 
              warehouse_id = $11,
              region_id = $12,
              price_size_id = $13,
            WHERE id = $14`
	_, err := repo.getWriteConnection().ExecOne(query,
		dbModel.ExternalID,
		dbModel.Price,
		dbModel.StatusID,
		dbModel.Direction,
		dbModel.Type,
		dbModel.Sale,
		dbModel.Quantity,
		dbModel.CreatedAt,

		dbModel.SellerID,
		dbModel.CardID,
		dbModel.WarehouseID,
		dbModel.RegionID,
		dbModel.PriceSizeID,
		dbModel.ID,
	)
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

func (repo *repoImpl) WithTx(tx *postgresdb.Transaction) OrderRepo {
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
