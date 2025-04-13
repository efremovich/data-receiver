package costrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

var ErrObjectNotFound = entity.ErrObjectNotFound

type CostRepo interface {
	SelectByCardIDAndDate(ctx context.Context, cardID int64, date time.Time) (*entity.Cost, error)
	Insert(ctx context.Context, in entity.Cost) (*entity.Cost, error)
	UpdateExecOne(ctx context.Context, in entity.Cost) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(t *postgresdb.Transaction) CostRepo
}

type costRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewCostRepo(_ context.Context, db *postgresdb.DBConnection) (CostRepo, error) {
	return &costRepoImpl{db: db}, nil
}

func (repo *costRepoImpl) SelectByCardIDAndDate(ctx context.Context, cardID int64, date time.Time) (*entity.Cost, error) {
	var result costsDB

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `SELECT
							id,
							card_id,
							amount,
							created_at,
							updated_at
						FROM
							shop.costs
						WHERE
							card_id = $1 and created_at between $2 and $3;`

	err := repo.getReadConnection().Get(&result, query, cardID, startOfDay, endOfDay)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данные по cardID %d и дате %s", cardID, startOfDay)
	}

	return result.convertToEntityCost(ctx), nil
}

func (repo *costRepoImpl) Insert(_ context.Context, in entity.Cost) (*entity.Cost, error) {
	startOfDay := time.Date(in.CreatedAt.Year(), in.CreatedAt.Month(), in.CreatedAt.Day(), 0, 0, 0, 0, time.UTC)

	query := `INSERT INTO shop.costs (card_id, amount, created_at, updated_at)
	VALUES ($1, $2, $3, $4) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query,
		in.CardID,
		in.Amount,
		startOfDay,
		time.Now(),
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка при вставке данных в таблицу cost: %w", err)
	}

	in.ID = charIDWrap.ID.Int64

	return &in, nil
}

func (repo *costRepoImpl) UpdateExecOne(ctx context.Context, in entity.Cost) error {
	dbModel := convertDBToCost(ctx, &in)
	startOfDay := time.Date(in.CreatedAt.Year(), in.CreatedAt.Month(), in.CreatedAt.Day(), 0, 0, 0, 0, time.UTC)

	query := `UPDATE shop.costs SET card_id = $1, amount=$2, created_at = $3, updated_at = $4 WHERE id = $5`

	_, err := repo.getWriteConnection().ExecOne(query,
		dbModel.CardID,
		dbModel.Amount,
		startOfDay,
		time.Now(),
		dbModel.ID,
	)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении данных в таблицу costs: %w", err)
	}

	return nil
}

func (repo *costRepoImpl) Ping(_ context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *costRepoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *costRepoImpl) WithTx(tx *postgresdb.Transaction) CostRepo {
	return &costRepoImpl{db: repo.db, tx: tx}
}

func (repo *costRepoImpl) getReadConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetReadConnection()
}

func (repo *costRepoImpl) getWriteConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetWriteConnection()
}
