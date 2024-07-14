package sellerrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type SellerRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Seller, error)
	SelectByTitle(ctx context.Context, title string) (*entity.Seller, error)
	Insert(ctx context.Context, in entity.Seller) (*entity.Seller, error)
	UpdateExecOne(ctx context.Context, in entity.Seller) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) SellerRepo
}

type repoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewSellerRepo(_ context.Context, db *postgresdb.DBConnection) (SellerRepo, error) {
	return &repoImpl{db: db}, nil
}

func (repo *repoImpl) SelectByID(ctx context.Context, id int64) (*entity.Seller, error) {
	var result sellerDB

	query := "SELECT * FROM sellers WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}
	return result.convertToEntitySeller(ctx), nil
}

func (repo *repoImpl) SelectByTitle(ctx context.Context, title string) (*entity.Seller, error) {
	var result sellerDB

	query := "SELECT * FROM sellers WHERE title = $1"

	err := repo.getReadConnection().Get(&result, query, title)
	if err != nil {
		return nil, err
	}

	return result.convertToEntitySeller(ctx), nil
}

func (repo *repoImpl) Insert(ctx context.Context, in entity.Seller) (*entity.Seller, error) {
	charIDWrap := repository.IDWrapper{}
	query := `INSERT INTO sellers (title, is_enable, ext_id) 
            VALUES ($1, $2, $3) RETURNING id`

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, in.Title, in.IsEnable, in.ExtID)
	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, in entity.Seller) error {
	dbModel := convertToDBSeller(ctx, in)

	query := `UPDATE sellers SET 
            title = $1, is_enable = $2, ext_id = $3
            WHERE id = $4`
	_, err := repo.getWriteConnection().ExecOne(query, in.Title, in.IsEnable, in.ExtID, dbModel.ID)
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

func (repo *repoImpl) WithTx(tx *postgresdb.Transaction) SellerRepo {
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
