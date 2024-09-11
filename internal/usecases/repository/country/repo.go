package country

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type CountryRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Country, error)
	SelectByName(ctx context.Context, name string) (*entity.Country, error)
	Insert(ctx context.Context, in entity.Country) (*entity.Country, error)
	UpdateExecOne(ctx context.Context, in entity.Country) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) CountryRepo
}

type repoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewCountryRepo(_ context.Context, db *postgresdb.DBConnection) (CountryRepo, error) {
	return &repoImpl{db: db}, nil
}

func (repo *repoImpl) SelectByID(ctx context.Context, id int64) (*entity.Country, error) {
	var result countryDB

	query := "SELECT * FROM shop.countries WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}

	return result.convertToEntityCountry(ctx), nil
}

func (repo *repoImpl) SelectByName(ctx context.Context, name string) (*entity.Country, error) {
	var result countryDB

	query := "SELECT * FROM shop.countries WHERE name = $1"

	err := repo.getReadConnection().Get(&result, query, name)
	if err != nil {
		return nil, err
	}

	return result.convertToEntityCountry(ctx), nil
}

func (repo *repoImpl) Insert(ctx context.Context, in entity.Country) (*entity.Country, error) {
	query := `INSERT INTO shop.countries (name) 
              VALUES ($1) RETURNING id`
	countryIDWrap := repository.IDWrapper{}
	dbModel := convertToDBCountry(ctx, &in)

	err := repo.getWriteConnection().QueryAndScan(&countryIDWrap, query, dbModel.Name)
	if err != nil {
		return nil, err
	}

	in.ID = countryIDWrap.ID.Int64

	return &in, nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, in entity.Country) error {
	dbModel := convertToDBCountry(ctx, &in)
	query := "UPDATE shop.countries SET name = $1 WHERE id = $2"

	_, err := repo.getWriteConnection().Exec(query, dbModel.Name, dbModel.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repoImpl) Ping(_ context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *repoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *repoImpl) WithTx(tx *postgresdb.Transaction) CountryRepo {
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
