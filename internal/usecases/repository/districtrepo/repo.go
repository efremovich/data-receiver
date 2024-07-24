package districtrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type DistrictRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.District, error)
	SelectByName(ctx context.Context, name string) (*entity.District, error)
	Insert(ctx context.Context, in entity.District) (*entity.District, error)
	UpdateExecOne(ctx context.Context, in entity.District) error
	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) DistrictRepo
}

type districtRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewDistrictRepo(_ context.Context, db *postgresdb.DBConnection) (DistrictRepo, error) {
	return &districtRepoImpl{db: db}, nil
}

func (repo *districtRepoImpl) SelectByID(ctx context.Context, id int64) (*entity.District, error) {
	var result districtDB
	query := "SELECT id, name FROM shop.districts WHERE id = $1"
	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}
	return result.ConvertToEntityDistrict(ctx), nil
}

func (repo *districtRepoImpl) SelectByName(ctx context.Context, name string) (*entity.District, error) {
	var result districtDB
	query := "SELECT id, name FROM shop.districts WHERE name = $1"
	err := repo.getReadConnection().Get(&result, query, name)
	if err != nil {
		return nil, err
	}
	return result.ConvertToEntityDistrict(ctx), nil
}

func (repo *districtRepoImpl) Insert(ctx context.Context, in entity.District) (*entity.District, error) {
	query := `INSERT INTO shop.districts (name) 
            VALUES ($1) RETURNING id`
	districtIDWrap := repository.IDWrapper{}
	dbModel := convertToDBDistrict(ctx, in)
	err := repo.getWriteConnection().QueryAndScan(&districtIDWrap, query, dbModel.Name)
	if err != nil {
		return nil, err
	}
	in.ID = districtIDWrap.ID.Int64
	return &in, nil
}

func (repo *districtRepoImpl) UpdateExecOne(ctx context.Context, in entity.District) error {
	dbModel := convertToDBDistrict(ctx, in)
	query := `UPDATE shop.districts SET name = $1 WHERE id = $2`
	_, err := repo.getWriteConnection().Exec(query, dbModel.Name, dbModel.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *districtRepoImpl) Ping(ctx context.Context) error {
	return repo.getReadConnection().Ping()
}

func (repo *districtRepoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *districtRepoImpl) WithTx(tx *postgresdb.Transaction) DistrictRepo {
	return &districtRepoImpl{db: repo.db, tx: tx}
}

func (repo *districtRepoImpl) getReadConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetReadConnection()
}

func (repo *districtRepoImpl) getWriteConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetWriteConnection()
}
