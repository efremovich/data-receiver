package cardcategory

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

type CardCategory interface {
	SelectByID(ctx context.Context, id int64) (*entity.CardCategory, error)
	SelectByCardIDAndSellerID(ctx context.Context, cardID, sellerID int64) (*entity.CardCategory, error)
	Insert(ctx context.Context, in entity.CardCategory) error
	UpdateExecOne(ctx context.Context, in entity.CardCategory) error
	Ping(ctx context.Context) error
}

type cardCategoryImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewCardCategory(_ context.Context, db *postgresdb.DBConnection) (CardCategory, error) {
	return &cardCategoryImpl{db: db}, nil
}

func (repo *cardCategoryImpl) SelectByID(ctx context.Context, id int64) (*entity.CardCategory, error) {
	var result cardCategoryDB

	query := `SELECT * FROM shop.cards_categories WHERE id = $1`

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по ID %d в таблице card_categories: %w", id, err)
	}

	return result.convertToEntityCardCategory(ctx), nil
}

func (repo *cardCategoryImpl) SelectByCardIDAndSellerID(ctx context.Context, cardID, sellerID int64) (*entity.CardCategory, error) {
	var result cardCategoryDB

	query := `SELECT * FROM shop.cards_categories WHERE card_id = $1 AND category_id = $2`

	err := repo.getReadConnection().Get(&result, query, cardID, sellerID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по title %d и sellerID %d в таблице card_categories: %w", cardID, sellerID, err)
	}

	return result.convertToEntityCardCategory(ctx), nil
}

func (repo *cardCategoryImpl) Insert(ctx context.Context, in entity.CardCategory) error {
	dbModel := convertToDB(ctx, in)
	query := `INSERT INTO shop.cards_categories (card_id, category_id)
  VALUES ($1, $2) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, dbModel.CardID, dbModel.CategoryID)
	if err != nil {
		return fmt.Errorf("ошибка при вставке данных в таблицу cards_categories: %w", err)
	}

	return nil
}

func (repo *cardCategoryImpl) UpdateExecOne(ctx context.Context, in entity.CardCategory) error {
	dbModel := convertToDB(ctx, in)
	query := `UPDATE shop.cards_categories SET card_id = $1, category_id = $2 WHERE id = $3`

	_, err := repo.getWriteConnection().ExecOne(query, dbModel.CardID, dbModel.CategoryID, dbModel.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении данных в таблицу cards_categories: %w", err)
	}

	return nil
}

func (repo *cardCategoryImpl) Ping(_ context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *cardCategoryImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *cardCategoryImpl) WithTx(tx *postgresdb.Transaction) CardCategory {
	return &cardCategoryImpl{db: repo.db, tx: tx}
}

func (repo *cardCategoryImpl) getReadConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetReadConnection()
}

func (repo *cardCategoryImpl) getWriteConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetWriteConnection()
}
