package brandrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type BrandRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Brand, error)
	SelectByTitle(ctx context.Context, CardID int64) ([]*entity.Brand, error)
	Insert(ctx context.Context, in entity.Brand) (*entity.Brand, error)
	UpdateExecOne(ctx context.Context, in entity.Brand) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) BrandRepo
}

type repoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewBrandRepo(_ context.Context, db *postgresdb.DBConnection) (BrandRepo, error) {
	return &repoImpl{db: db}, nil
}

func (repo *repoImpl) SelectByID(ctx context.Context, id int64) (*entity.Brand, error) {
	var result brandDB

	query := "SELECT * FROM brands WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}
	return result.convertToEntityBrand(ctx), nil
}

func (repo *repoImpl) SelectByTitle(ctx context.Context, cardID int64) ([]*entity.Brand, error) {
	var result []brandDB

	query := "SELECT * FROM brands WHERE title = $1"

	err := repo.getReadConnection().Select(&result, query, cardID)
	if err != nil {
		return nil, err
	}

	var resEntity []*entity.Brand
	for _, v := range result {
		resEntity = append(resEntity, v.convertToEntityBrand(ctx))
	}
	return resEntity, nil
}

func (repo *repoImpl) Insert(ctx context.Context, in entity.Brand) (*entity.Brand, error) {
	query := `INSERT INTO brands (title seller_id) 
            VALUES ($1, $2) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, in.Title, in.SellerID)
	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, in entity.Brand) error {
	dbModel := convertToDBBrand(ctx, in)

	query := `UPDATE brands SET 
            title = $1 seller_id = $2
            WHERE id = $3`
	_, err := repo.getWriteConnection().ExecOne(query, in.Title, in.SellerID, dbModel.ID)
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

func (repo *repoImpl) WithTx(tx *postgresdb.Transaction) BrandRepo {
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
