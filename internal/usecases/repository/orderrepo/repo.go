package orderrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type OrderRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Order, error)
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

	query := "SELECT * FROM orders WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}
	return result.convertToEntityOrder(ctx), nil
}

func (repo *repoImpl) Insert(ctx context.Context, in entity.Order) (*entity.Order, error) {
	dbModel := convertToDBOrder(ctx, in)

	query := `INSERT INTO orders (
              ext_id,
              price,
              quantity,
              discount,
              special_price,
              status, 
              type,
              warehouse_id,
              seller_id, 
              card_id, 
              direction,
              created_at
            ) 
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, now()) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query,
		dbModel.ExtID,
		dbModel.Price,
		dbModel.Quantity,
		dbModel.Discount,
		dbModel.SpecialPrice,
		dbModel.Status,
		dbModel.Type,
		dbModel.WarehouseID,
		dbModel.SellerID,
		dbModel.CardID,
		dbModel.Direction,
	)
	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, in entity.Order) error {
	dbModel := convertToDBOrder(ctx, in)

	query := `UPDATE orders SET
            ext_id = $1, 
            price = $2,
            quantity = $3,
            discount = $4, 
            special_price = $5,
            status = $6, 
            type = $7,
            warehouse_id = $8,
            seller_id = $9, 
            card_id = $10,
            direction = $11,
            updated_at = now()
            WHERE id = $12`
	_, err := repo.getWriteConnection().ExecOne(query,
		dbModel.ExtID,
		dbModel.Price,
		dbModel.Quantity,
		dbModel.Discount,
		dbModel.SpecialPrice,
		dbModel.Status,
		dbModel.Type,
		dbModel.WarehouseID,
		dbModel.SellerID,
		dbModel.CardID,
		dbModel.Direction,
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
