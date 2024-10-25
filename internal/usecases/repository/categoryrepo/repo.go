package categoryrepo

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

type CategoryRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Category, error)
	SelectBySellerID(ctx context.Context, sellerID int64) ([]*entity.Category, error)
	SelectByTitle(ctx context.Context, title string, sellerID int64) (*entity.Category, error)
	Insert(ctx context.Context, in entity.Category) (*entity.Category, error)
	UpdateExecOne(ctx context.Context, in entity.Category) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) CategoryRepo
}

type charRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewCategoryRepo(_ context.Context, db *postgresdb.DBConnection) (CategoryRepo, error) {
	return &charRepoImpl{db: db}, nil
}

func (repo *charRepoImpl) SelectByID(ctx context.Context, id int64) (*entity.Category, error) {
	var result categoryDB

	query := "SELECT * FROM shop.categories WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по id %d в таблице categories: %w", id, err)
	}

	return result.convertToEntityCategory(ctx), nil
}

func (repo *charRepoImpl) SelectBySellerID(ctx context.Context, sellerID int64) ([]*entity.Category, error) {
	var result []categoryDB

	query := "SELECT * FROM shop.categories WHERE seller_id = $1"

	err := repo.getReadConnection().Select(&result, query, sellerID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по ID продавца %d в таблице categories: %w", sellerID, err)
	}

	var resEntity []*entity.Category
	for _, v := range result {
		resEntity = append(resEntity, v.convertToEntityCategory(ctx))
	}

	return resEntity, nil
}

func (repo *charRepoImpl) SelectByTitle(ctx context.Context, title string, sellerID int64) (*entity.Category, error) {
	var result categoryDB

	query := "SELECT * FROM shop.categories WHERE title = $1 and seller_id = $2"

	err := repo.getReadConnection().Get(&result, query, title, sellerID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по имени категории %s таблице categories: %w", title, err)
	}

	return result.convertToEntityCategory(ctx), nil
}

func (repo *charRepoImpl) Insert(_ context.Context, in entity.Category) (*entity.Category, error) {
	query := `INSERT INTO shop.categories (title, seller_id, card_id, external_id, parent_id)
            VALUES ($1, $2, $3, $4, $5) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, in.Title, in.SellerID, in.CardID, in.ExternalID, in.ParentID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при вставке данных в %s продавца %d таблицу categories: %w", in.Title, in.SellerID, err)
	}

	in.ID = charIDWrap.ID.Int64

	return &in, nil
}

func (repo *charRepoImpl) UpdateExecOne(ctx context.Context, in entity.Category) error {
	dbModel := convertToDBCategory(ctx, in)

	query := `UPDATE shop.categories SET title = $1, seller_id = $2, card_id = $3, external_id = $4, parent_id = $5 WHERE id = $6`

	_, err := repo.getWriteConnection().ExecOne(query, dbModel.Title, dbModel.SellerID, dbModel.CardID, dbModel.ExternalID, dbModel.ParentID, dbModel.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении данных в %s продавца %d таблицу categories: %w", in.Title, in.SellerID, err)
	}

	return nil
}

func (repo *charRepoImpl) Ping(_ context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *charRepoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *charRepoImpl) WithTx(tx *postgresdb.Transaction) CategoryRepo {
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
