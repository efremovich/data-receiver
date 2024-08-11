package dimensionrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type DimensionsRepo interface {
	SelectByCardID(ctx context.Context, cardID int64) (*entity.Dimension, error)
	Insert(ctx context.Context, in entity.Dimension) (*entity.Dimension, error)
	UpdateExecOne(ctx context.Context, in entity.Dimension) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) DimensionsRepo
}

type charRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewDimensionRepo(_ context.Context, db *postgresdb.DBConnection) (DimensionsRepo, error) {
	return &charRepoImpl{db: db}, nil
}

func (repo *charRepoImpl) SelectByCardID(ctx context.Context, cardID int64) (*entity.Dimension, error) {
	var result dimensionDB

	query := "SELECT * FROM shop.dimensions WHERE card_id = $1"

	err := repo.getReadConnection().Get(&result, query, cardID)
	if err != nil {
		return nil, err
	}
	return result.convertToEntityDimension(ctx), nil
}

func (repo *charRepoImpl) Insert(ctx context.Context, in entity.Dimension) (*entity.Dimension, error) {
	query := `INSERT INTO shop.dimensions (length, width, height, is_valid, card_id) 
            VALUES ($1, $2, $3, $4, $5) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, in.Length, in.Width, in.Height, in.IsVaild, in.CardID)
	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *charRepoImpl) UpdateExecOne(ctx context.Context, in entity.Dimension) error {
	dbModel := convertToDBDimension(ctx, in)

	query := `UPDATE shop.dimensions SET length = $1, width = $2, height = $3, is_valid = $4, card_id = $5 WHERE id = $6`
	_, err := repo.getWriteConnection().ExecOne(query, dbModel.Length, dbModel.Width, dbModel.Height, dbModel.IsValid, dbModel.CardID, dbModel.ID)
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

func (repo *charRepoImpl) WithTx(tx *postgresdb.Transaction) DimensionsRepo {
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
