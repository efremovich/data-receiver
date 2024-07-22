package wb2cardrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type Wb2CardRepo interface {
	SelectByNmid(ctx context.Context, nmid int64) (*entity.Wb2Card, error)
	SelectByCardID(ctx context.Context, cardID int64) (*entity.Wb2Card, error)
	SelectByNMUUID(ctx context.Context, nmUUID string) (*entity.Wb2Card, error)
	SelectByKTID(ctx context.Context, ktID int) (*entity.Wb2Card, error)
	Insert(ctx context.Context, in entity.Wb2Card) (*entity.Wb2Card, error)
	UpdateExecOne(ctx context.Context, in entity.Wb2Card) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) Wb2CardRepo
}

type wb2cardRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewWb2CardRepo(_ context.Context, db *postgresdb.DBConnection) (Wb2CardRepo, error) {
	return &wb2cardRepoImpl{db: db}, nil
}

func (repo *wb2cardRepoImpl) SelectByNmid(ctx context.Context, nmid int64) (*entity.Wb2Card, error) {
	var result wb2cardDB
	query := "SELECT id, nmid, int, nmuuid, created_at, updated_at, card_id FROM shop.wb2cards WHERE nmid = $1"
	err := repo.getReadConnection().Get(&result, query, nmid)
	if err != nil {
		return nil, err
	}
	return result.ConvertToEntityWb2Card(ctx), nil
}

func (repo *wb2cardRepoImpl) SelectByKTID(ctx context.Context, ktID int) (*entity.Wb2Card, error) {
	var result wb2cardDB
	query := "SELECT id, nmid, int, nmuuid, created_at, updated_at, card_id FROM shop.wb2cards WHERE int = $1"
	err := repo.getReadConnection().Get(&result, query, ktID)
	if err != nil {
		return nil, err
	}
	return result.ConvertToEntityWb2Card(ctx), nil
}

func (repo *wb2cardRepoImpl) SelectByCardID(ctx context.Context, cardID int64) (*entity.Wb2Card, error) {
	var result wb2cardDB
	query := "SELECT id, nmid, int, nmuuid, created_at, updated_at, card_id FROM shop.wb2cards WHERE card_id = $1"
	err := repo.getReadConnection().Get(&result, query, cardID)
	if err != nil {
		return nil, err
	}
	return result.ConvertToEntityWb2Card(ctx), nil
}

func (repo *wb2cardRepoImpl) SelectByNMUUID(ctx context.Context, nmUUID string) (*entity.Wb2Card, error) {
	var result wb2cardDB
	query := "SELECT id, nmid, int, nmuuid, created_at, updated_at, card_id FROM shop.wb2cards WHERE nmuuid = $1"
	err := repo.getReadConnection().Get(&result, query, nmUUID)
	if err != nil {
		return nil, err
	}
	return result.ConvertToEntityWb2Card(ctx), nil
}

func (repo *wb2cardRepoImpl) Insert(ctx context.Context, in entity.Wb2Card) (*entity.Wb2Card, error) {
	query := `INSERT INTO shop.wb2cards (nmid, int, nmuuid, created_at, updated_at, card_id) 
            VALUES ($1, $2, $3, now(), now(), $4) RETURNING id`
	wb2cardIDWrap := repository.IDWrapper{}
	dbModel := convertToDBWb2Card(ctx, in)
	err := repo.getWriteConnection().QueryAndScan(&wb2cardIDWrap, query, dbModel.NMID, dbModel.KTID, dbModel.NMUUID, dbModel.CardID)
	if err != nil {
		return nil, err
	}
	in.ID = wb2cardIDWrap.ID.Int64
	return &in, nil
}

func (repo *wb2cardRepoImpl) UpdateExecOne(ctx context.Context, in entity.Wb2Card) error {
	dbModel := convertToDBWb2Card(ctx, in)
	query := `UPDATE shop.wb2cards SET nmid = $1, int = $2, nmuuid = $3, updated_at = now(), card_id = $4 WHERE id = $5`
	_, err := repo.getWriteConnection().Exec(query, dbModel.NMID, dbModel.KTID, dbModel.NMUUID, dbModel.CardID, dbModel.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *wb2cardRepoImpl) Ping(ctx context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *wb2cardRepoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *wb2cardRepoImpl) WithTx(tx *postgresdb.Transaction) Wb2CardRepo {
	return &wb2cardRepoImpl{db: repo.db, tx: tx}
}

func (repo *wb2cardRepoImpl) getReadConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetReadConnection()
}

func (repo *wb2cardRepoImpl) getWriteConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetWriteConnection()
}
