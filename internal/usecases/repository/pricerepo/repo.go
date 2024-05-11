package pricerepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type PriceRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Price, error)
	SelectByCardID(ctx context.Context, CardID int64) ([]*entity.Price, error)
	SelectByPriceID(ctx context.Context, CardID int64) ([]*entity.Price, error)
	Insert(ctx context.Context, in entity.Price) (*entity.Price, error)
	UpdateExecOne(ctx context.Context, in entity.Price) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) PriceRepo
}

type charRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewPriceRepo(_ context.Context, db *postgresdb.DBConnection) (PriceRepo, error) {
	return &charRepoImpl{db: db}, nil
}

func (repo *charRepoImpl) SelectByID(ctx context.Context, id int64) (*entity.Price, error) {
	var result priceDB

	query := "SELECT * FROM prices WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}
	return result.convertToEntityPrice(ctx), nil
}

func (repo *charRepoImpl) SelectByCardID(ctx context.Context, cardID int64) ([]*entity.Price, error) {
	var result []priceDB

	query := "SELECT * FROM prices WHERE card_id = $1"

	err := repo.getReadConnection().Select(&result, query, cardID)
	if err != nil {
		return nil, err
	}

	var resEntity []*entity.Price
	for _, v := range result {
		resEntity = append(resEntity, v.convertToEntityPrice(ctx))
	}
	return resEntity, nil
}

func (repo *charRepoImpl) SelectByPriceID(ctx context.Context, priceID int64) ([]*entity.Price, error) {
	var result []priceDB

	query := "SELECT * FROM prices WHERE price_id = $1"

	err := repo.getReadConnection().Select(&result, query, priceID)
	if err != nil {
		return nil, err
	}

	var resEntity []*entity.Price
	for _, v := range result {
		resEntity = append(resEntity, v.convertToEntityPrice(ctx))
	}
	return resEntity, nil
}

func (repo *charRepoImpl) Insert(ctx context.Context, in entity.Price) (*entity.Price, error) {
	query := `INSERT INTO prices (price, discount, special_price, seller_id, card_id) 
            VALUES ($1, $2, $3, $4, $5) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, in.Price, in.Discount, in.SpecialPrice, in.SellerID, in.CardID)
	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *charRepoImpl) UpdateExecOne(ctx context.Context, in entity.Price) error {
	dbModel := convertToDBPrice(ctx, in)

	query := `UPDATE prices SET 
            price = $1, discount = $2, special_price = $3, seller_id = $4, card_id = $5
            WHERE id = $6`
	_, err := repo.getWriteConnection().ExecOne(query, dbModel.Price, dbModel.Discount, dbModel.SpecialPrice, dbModel.SellerID, dbModel.CardID, dbModel.ID)
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

func (repo *charRepoImpl) WithTx(tx *postgresdb.Transaction) PriceRepo {
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
