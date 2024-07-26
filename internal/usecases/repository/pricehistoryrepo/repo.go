package pricehistoryrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type PriceHistoryRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.PriceHistory, error)
	SelectByPriceID(ctx context.Context, CardID int64) ([]*entity.PriceHistory, error)
	Insert(ctx context.Context, in entity.PriceHistory) (*entity.PriceHistory, error)
	UpdateExecOne(ctx context.Context, in entity.PriceHistory) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) PriceHistoryRepo
}

type charRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewPriceHistoryRepo(_ context.Context, db *postgresdb.DBConnection) (PriceHistoryRepo, error) {
	return &charRepoImpl{db: db}, nil
}

func (repo *charRepoImpl) SelectByID(ctx context.Context, id int64) (*entity.PriceHistory, error) {
	var result priceHistoryDB

	query := "SELECT * FROM shop.price_sizes WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}
	return result.convertToEntityPriceHistory(ctx), nil
}

func (repo *charRepoImpl) SelectByPriceID(ctx context.Context, priceID int64) ([]*entity.PriceHistory, error) {
	var result []priceHistoryDB

	query := "SELECT * FROM shop.price_history WHERE price_id = $1"

	err := repo.getReadConnection().Select(&result, query, priceID)
	if err != nil {
		return nil, err
	}

	var resEntity []*entity.PriceHistory
	for _, v := range result {
		resEntity = append(resEntity, v.convertToEntityPriceHistory(ctx))
	}
	return resEntity, nil
}

func (repo *charRepoImpl) Insert(ctx context.Context, in entity.PriceHistory) (*entity.PriceHistory, error) {
	query := `INSERT INTO shop.price_history (price, discount, price_size_id, updated_at)
            VALUES ($1, $2, $3, now()) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, in.Price, in.Discount, in.PriceSizeID)
	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *charRepoImpl) UpdateExecOne(ctx context.Context, in entity.PriceHistory) error {
	dbModel := convertToDBPriceHistory(ctx, in)

	query := `UPDATE shop.price_sizes SET 
            price = $1, discount = $2, price_size_id = $3, updated_at = now()
            WHERE id = $4`
	_, err := repo.getWriteConnection().ExecOne(query, dbModel.Price, dbModel.Discount, dbModel.PriceSizeID, dbModel.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *charRepoImpl) Ping(ctx context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *charRepoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *charRepoImpl) WithTx(tx *postgresdb.Transaction) PriceHistoryRepo {
	return &charRepoImpl{db: repo.db, tx: tx}
}

func (repo *charRepoImpl) getReadConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetReadConnection()
}

func (repo *charRepoImpl) getWriteConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetWriteConnection()
}
