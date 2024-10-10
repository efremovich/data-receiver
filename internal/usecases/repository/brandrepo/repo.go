package brandrepo

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

type BrandRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Brand, error)
	SelectByTitleAndSeller(ctx context.Context, title string, sellerID int64) (*entity.Brand, error)
	Insert(ctx context.Context, in entity.Brand) (*entity.Brand, error)
	UpdateExecOne(ctx context.Context, in entity.Brand) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) BrandRepo
}

type repoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewBrandRepo(_ context.Context, db *postgresdb.DBConnection) (BrandRepo, error) {
	return &repoImpl{db: db}, nil
}

func (repo *repoImpl) SelectByID(ctx context.Context, id int64) (*entity.Brand, error) {
	var result brandDB

	query := "SELECT * FROM shop.brands WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по id %d в таблице brands: %w", id, err)
	}

	return result.convertToEntityBrand(ctx), nil
}

func (repo *repoImpl) SelectByTitleAndSeller(ctx context.Context, title string, sellerID int64) (*entity.Brand, error) {
	var result brandDB

	query := "SELECT * FROM shop.brands WHERE title = $1 and seller_id = $2"

	err := repo.getReadConnection().Get(&result, query, title, sellerID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по имени %s продавца %d в таблице brands: %w", title, sellerID, err)
	}

	return result.convertToEntityBrand(ctx), nil
}

func (repo *repoImpl) Insert(_ context.Context, in entity.Brand) (*entity.Brand, error) {
	charIDWrap := repository.IDWrapper{}

	query := `INSERT INTO shop.brands (title, seller_id) 
            VALUES ($1, $2)
            ON CONFLICT (id) DO NOTHING
            RETURNING id`

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, in.Title, in.SellerID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при вставке данных %s в таблицу brands: %w", in.Title, err)
	}

	in.ID = charIDWrap.ID.Int64

	return &in, nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, in entity.Brand) error {
	dbModel := convertToDBBrand(ctx, in)

	query := `UPDATE shop.brands SET 
            title = $1, seller_id = $2
            WHERE id = $3`
	_, err := repo.getWriteConnection().ExecOne(query, in.Title, in.SellerID, dbModel.ID)

	if err != nil {
		return fmt.Errorf("ошибка при обновлении данных в %s таблицу brands: %w", in.Title, err)
	}

	return nil
}

func (repo *repoImpl) Ping(_ context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *repoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *repoImpl) WithTx(tx *postgresdb.Transaction) BrandRepo {
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
