package countryrepo

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

type CountryRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Country, error)
	SelectByName(ctx context.Context, name string) (*entity.Country, error)
	Insert(ctx context.Context, in entity.Country) (*entity.Country, error)
	UpdateExecOne(ctx context.Context, in entity.Country) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) CountryRepo
}

type countryRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewCountryRepo(_ context.Context, db *postgresdb.DBConnection) (CountryRepo, error) {
	return &countryRepoImpl{db: db}, nil
}

func (repo *countryRepoImpl) SelectByID(ctx context.Context, id int64) (*entity.Country, error) {
	var result countryDB

	query := "SELECT id, name FROM shop.countries WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по внешнему ID %d в таблице countries: %w", id, err)
	}

	return result.ConvertToEntityCountry(ctx), nil
}

func (repo *countryRepoImpl) SelectByName(ctx context.Context, name string) (*entity.Country, error) {
	var result countryDB

	query := "SELECT id, name FROM shop.countries WHERE name = $1"

	err := repo.getReadConnection().Get(&result, query, name)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по имени %s в таблице countries: %w", name, err)
	}

	return result.ConvertToEntityCountry(ctx), nil
}

func (repo *countryRepoImpl) Insert(ctx context.Context, in entity.Country) (*entity.Country, error) {
	query := `INSERT INTO shop.countries (name) 
            VALUES ($1) RETURNING id`
	countryIDWrap := repository.IDWrapper{}
	dbModel := convertToDBCountry(ctx, in)

	err := repo.getWriteConnection().QueryAndScan(&countryIDWrap, query, dbModel.Name)
	if err != nil {
		return nil, fmt.Errorf("ошибка вставки данных в таблицу countries %s: %w", in.Name, err)
	}

	in.ID = countryIDWrap.ID.Int64

	return &in, nil
}

func (repo *countryRepoImpl) UpdateExecOne(ctx context.Context, in entity.Country) error {
	query := `UPDATE shop.countries SET name = $1 WHERE id = $2`
	dbModel := convertToDBCountry(ctx, in)

	_, err := repo.getWriteConnection().Exec(query, dbModel.Name, dbModel.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления данных в таблице countries %s: %w", in.Name, err)
	}

	return nil
}

func (repo *countryRepoImpl) Ping(_ context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *countryRepoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *countryRepoImpl) WithTx(tx *postgresdb.Transaction) CountryRepo {
	return &countryRepoImpl{db: repo.db, tx: tx}
}

func (repo *countryRepoImpl) getReadConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetReadConnection()
}

func (repo *countryRepoImpl) getWriteConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetWriteConnection()
}
