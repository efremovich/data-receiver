package seller2cardrepo

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

type Seller2CardRepo interface {
	SelectByExternalID(ctx context.Context, externalID, sellerID int64) (*entity.Seller2Card, error)
	SelectByCardID(ctx context.Context, cardID int64) (*entity.Seller2Card, error)
	SelectByNMUUID(ctx context.Context, nmUUID string) (*entity.Seller2Card, error)
	SelectByKTID(ctx context.Context, ktID int) (*entity.Seller2Card, error)
	Insert(ctx context.Context, in entity.Seller2Card) (*entity.Seller2Card, error)
	UpdateExecOne(ctx context.Context, in entity.Seller2Card) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) Seller2CardRepo
}

type seller2cardRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewWb2CardRepo(_ context.Context, db *postgresdb.DBConnection) (Seller2CardRepo, error) {
	return &seller2cardRepoImpl{db: db}, nil
}

func (repo *seller2cardRepoImpl) SelectByExternalID(ctx context.Context, externalID, sellerID int64) (*entity.Seller2Card, error) {
	var result seller2cardDB

	query := "SELECT id, external_id, int, nmuuid, created_at, updated_at, card_id FROM shop.seller2cards WHERE external_id = $1 and seller_id = $2"

	err := repo.getReadConnection().Get(&result, query, externalID, sellerID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по внешнему ID %d в таблице seller2cards: %w", externalID, err)
	}

	return result.ConvertToEntityWb2Card(ctx), nil
}

func (repo *seller2cardRepoImpl) SelectByKTID(ctx context.Context, ktID int) (*entity.Seller2Card, error) {
	var result seller2cardDB

	query := "SELECT id, external_id, int, nmuuid, created_at, updated_at, card_id FROM shop.seller2cards WHERE int = $1"

	err := repo.getReadConnection().Get(&result, query, ktID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по KTID %d в таблице seller2cards: %w", ktID, err)
	}

	return result.ConvertToEntityWb2Card(ctx), nil
}

func (repo *seller2cardRepoImpl) SelectByCardID(ctx context.Context, cardID int64) (*entity.Seller2Card, error) {
	var result seller2cardDB

	query := "SELECT id, external_id, int, nmuuid, created_at, updated_at, card_id FROM shop.seller2cards WHERE card_id = $1"

	err := repo.getReadConnection().Get(&result, query, cardID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по id карточки %d в таблице seller2cards: %w", cardID, err)
	}

	return result.ConvertToEntityWb2Card(ctx), nil
}

func (repo *seller2cardRepoImpl) SelectByNMUUID(ctx context.Context, nmUUID string) (*entity.Seller2Card, error) {
	var result seller2cardDB

	query := "SELECT id, external_id, int, nmuuid, created_at, updated_at, card_id FROM shop.seller2cards WHERE nmuuid = $1"

	err := repo.getReadConnection().Get(&result, query, nmUUID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по NMUUID %s в таблице seller2cards: %w", nmUUID, err)
	}

	return result.ConvertToEntityWb2Card(ctx), nil
}

func (repo *seller2cardRepoImpl) Insert(ctx context.Context, in entity.Seller2Card) (*entity.Seller2Card, error) {
	query := `INSERT INTO shop.seller2cards (external_id, int, nmuuid, created_at, updated_at, card_id, seller_id) 
            VALUES ($1, $2, $3, now(), now(), $4, $5) RETURNING id`

	wb2cardIDWrap := repository.IDWrapper{}
	dbModel := convertToDBWb2Card(ctx, in)

	err := repo.getWriteConnection().QueryAndScan(&wb2cardIDWrap, query,
		dbModel.ExternalID,
		dbModel.KTID,
		dbModel.NMUUID,
		dbModel.CardID,
		dbModel.SellerID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при вставке данных в таблицу seller2cards: %w", err)
	}

	in.ID = wb2cardIDWrap.ID.Int64

	return &in, nil
}

func (repo *seller2cardRepoImpl) UpdateExecOne(ctx context.Context, in entity.Seller2Card) error {
	dbModel := convertToDBWb2Card(ctx, in)
	query := `UPDATE shop.seller2cards SET external_id = $1, int = $2, nmuuid = $3, updated_at = now(), card_id = $4 WHERE id = $5`

	_, err := repo.getWriteConnection().Exec(query, dbModel.ExternalID, dbModel.KTID, dbModel.NMUUID, dbModel.CardID, dbModel.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении данных в таблицу seller2cards: %w", err)
	}

	return nil
}

func (repo *seller2cardRepoImpl) Ping(_ context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *seller2cardRepoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *seller2cardRepoImpl) WithTx(tx *postgresdb.Transaction) Seller2CardRepo {
	return &seller2cardRepoImpl{db: repo.db, tx: tx}
}

func (repo *seller2cardRepoImpl) getReadConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetReadConnection()
}

func (repo *seller2cardRepoImpl) getWriteConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetWriteConnection()
}
