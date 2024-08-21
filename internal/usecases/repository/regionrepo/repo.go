package regionrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type RegoinRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Region, error)
	SelectByName(ctx context.Context, regionName string, countryID, districtID int64) (*entity.Region, error)
	Insert(ctx context.Context, region *entity.Region) (*entity.Region, error)
	UpdateExecOne(ctx context.Context, region *entity.Region) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) RegoinRepo
}

type repoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewRegionRepo(_ context.Context, db *postgresdb.DBConnection) (RegoinRepo, error) {
	return &repoImpl{db: db}, nil
}

func (repo *repoImpl) SelectByID(ctx context.Context, id int64) (*entity.Region, error) {
	var result regionDB

	query := "SELECT * FROM shop.regions WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}

	return result.convertToEntityRegion(ctx), nil
}

func (repo *repoImpl) SelectByName(ctx context.Context, regionName string, districtID, countryID int64) (*entity.Region, error) {
	var result regionDB

	query := "SELECT * FROM shop.regions WHERE country_id = $1 and region_name = $2 and district_id = $3"

	err := repo.getReadConnection().Get(&result, query, countryID, regionName, districtID)
	if err != nil {
		return nil, err
	}

	return result.convertToEntityRegion(ctx), nil
}

func (repo *repoImpl) Insert(ctx context.Context, region *entity.Region) (*entity.Region, error) {
	dbModel := convertToDBRegion(ctx, region)
	query := "INSERT INTO shop.regions (country_id, region_name, district_id) VALUES ($1, $2, $3) RETURNING id"
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, dbModel.CountryID, dbModel.RegionName, dbModel.DistrictID)
	if err != nil {
		return nil, err
	}

	region.ID = charIDWrap.ID.Int64

	return region, nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, region *entity.Region) error {
	dbModel := convertToDBRegion(ctx, region)
	query := "UPDATE shop.regions SET country_id = $1, region_name = $2, district_id = $3 WHERE id = $4"

	_, err := repo.getWriteConnection().Exec(query, dbModel.CountryID, dbModel.RegionName, dbModel.DistrictID, region.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repoImpl) Ping(ctx context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *repoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *repoImpl) WithTx(tx *postgresdb.Transaction) RegoinRepo {
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
