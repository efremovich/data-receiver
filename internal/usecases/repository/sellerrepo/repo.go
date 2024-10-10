package sellerrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

var ErrObjectNotFound = entity.ErrObjectNotFound

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

	query := "SELECT * FROM shop.sellers WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по id %d в таблице seller: %w", id, err)
	}

	return result.convertToEntitySeller(ctx), nil
}

func (repo *repoImpl) SelectByTitle(ctx context.Context, title string) (*entity.Seller, error) {
	var result sellerDB

	query := "SELECT * FROM shop.sellers WHERE title = $1"

	err := repo.getReadConnection().Get(&result, query, title)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по имени %s в таблице sellers: %w", title, err)
	}

	return result.convertToEntitySeller(ctx), nil
}

func (repo *repoImpl) Insert(_ context.Context, in entity.Seller) (*entity.Seller, error) {
	charIDWrap := repository.IDWrapper{}
	query := `INSERT INTO shop.sellers (title, is_enabled, external_id) 
            VALUES ($1, $2, $3) RETURNING id`

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, in.Title, in.IsEnabled, in.ExternalID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при вставке данных в %s таблицу seller: %w", in.Title, err)
	}

	in.ID = charIDWrap.ID.Int64

	return &in, nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, in entity.Seller) error {
	dbModel := convertToDBSeller(ctx, in)

	query := `UPDATE shop.sellers SET 
            title = $1, is_enabled = $2, external_id = $3
            WHERE id = $4`
	_, err := repo.getWriteConnection().ExecOne(query, in.Title, in.IsEnabled, in.ExternalID, dbModel.ID)

	if err != nil {
		return fmt.Errorf("ошибка при обновлении данных в %s таблицу seller: %w", in.Title, err)
	}

	return nil
}

func (repo *repoImpl) Ping(_ context.Context) error {
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
