package cardrepo

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

type CardRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Card, error)
	SelectByExternalID(ctx context.Context, id int64) (*entity.Card, error)
	SelectByVendorID(ctx context.Context, vendorID string) (*entity.Card, error)
	SelectByVendorCode(ctx context.Context, vendorCode string) (*entity.Card, error)
	SelectByTitle(ctx context.Context, title string) (*entity.Card, error)

	Insert(ctx context.Context, in entity.Card) (*entity.Card, error)
	UpdateExecOne(ctx context.Context, in entity.Card) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) CardRepo
}

type cardRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewCardRepo(_ context.Context, db *postgresdb.DBConnection) (CardRepo, error) {
	return &cardRepoImpl{db: db}, nil
}

func (repo *cardRepoImpl) SelectByID(ctx context.Context, id int64) (*entity.Card, error) {
	var result cardDB

	query := "SELECT * FROM shop.cards WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по id %d в таблице cards: %w", id, err)
	}

	card := result.ConvertToEntityCard(ctx)

	return card, nil
}

func (repo *cardRepoImpl) SelectByExternalID(ctx context.Context, externalID int64) (*entity.Card, error) {
	var result cardDB

	query := "SELECT * FROM shop.cards WHERE external_id = $1"

	err := repo.getReadConnection().Get(&result, query, externalID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по external_id %d в таблице cards: %w", externalID, err)
	}

	card := result.ConvertToEntityCard(ctx)

	return card, nil
}

func (repo *cardRepoImpl) SelectByVendorID(ctx context.Context, vendorID string) (*entity.Card, error) {
	var result cardDB

	query := "SELECT * FROM shop.cards WHERE vendor_id = $1"

	err := repo.getReadConnection().Get(&result, query, vendorID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по vendor_id %s в таблице cards: %w", vendorID, err)
	}

	return result.ConvertToEntityCard(ctx), nil
}

func (repo *cardRepoImpl) SelectByVendorCode(ctx context.Context, vendorCode string) (*entity.Card, error) {
	var result cardDB

	query := "SELECT * FROM shop.cards WHERE vendor_code = $1"

	err := repo.getReadConnection().Get(&result, query, vendorCode)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по vendor_code %s в таблице cards: %w", vendorCode, err)
	}

	return result.ConvertToEntityCard(ctx), nil
}

func (repo *cardRepoImpl) SelectByTitle(ctx context.Context, title string) (*entity.Card, error) {
	var result cardDB

	query := "SELECT * FROM shop.cards WHERE title = $1"

	err := repo.getReadConnection().Get(&result, query, title)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по title %s в таблице cards: %w", title, err)
	}

	return result.ConvertToEntityCard(ctx), nil
}

func (repo *cardRepoImpl) Insert(_ context.Context, in entity.Card) (*entity.Card, error) {
	query := `INSERT INTO shop.cards (vendor_id, vendor_code, title, description, brand_id, created_at) 
            VALUES ($1, $2, $3, $4, $5, now()) RETURNING id`
	cardIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&cardIDWrap, query, in.VendorID, in.VendorCode, in.Title, in.Description, in.Brand.ID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при вставке данных в таблицу cards: %w", err)
	}

	in.ID = cardIDWrap.ID.Int64

	return &in, nil
}

func (repo *cardRepoImpl) UpdateExecOne(ctx context.Context, in entity.Card) error {
	dbModel := convertToDBCard(ctx, in)

	query := `UPDATE shop.cards SET vendor_id = $1, vendor_code = $2, title = $3, description = $4, updated_at = NOW() WHERE id = $5`

	_, err := repo.getWriteConnection().ExecOne(query, in.VendorID, in.VendorCode, in.Title, in.Description, dbModel.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении данных в таблицу cards: %w", err)
	}

	return nil
}

func (repo *cardRepoImpl) Ping(_ context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *cardRepoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *cardRepoImpl) WithTx(tx *postgresdb.Transaction) CardRepo {
	return &cardRepoImpl{db: repo.db, tx: tx}
}

func (repo *cardRepoImpl) getReadConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetReadConnection()
}

func (repo *cardRepoImpl) getWriteConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetWriteConnection()
}
