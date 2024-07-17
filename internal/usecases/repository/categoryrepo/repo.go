package categoryrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type CategoryRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Category, error)
	SelectBySellerID(ctx context.Context, sellerID int64) ([]*entity.Category, error)
	SelectByTitle(ctx context.Context, title string) (*entity.Category, error)
	Insert(ctx context.Context, in entity.Category) (*entity.Category, error)
	UpdateExecOne(ctx context.Context, in entity.Category) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) CategoryRepo
}

type charRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewCategoryRepo(_ context.Context, db *postgresdb.DBConnection) (CategoryRepo, error) {
	return &charRepoImpl{db: db}, nil
}

func (repo *charRepoImpl) SelectByID(ctx context.Context, id int64) (*entity.Category, error) {
	var result categoryDB

	query := "SELECT * FROM shop.categories WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}
	return result.convertToEntityCategory(ctx), nil
}

func (repo *charRepoImpl) SelectBySellerID(ctx context.Context, sellerID int64) ([]*entity.Category, error) {
	var result []categoryDB

	query := "SELECT * FROM shop.categories WHERE seller_id = $1"

	err := repo.getReadConnection().Select(&result, query, sellerID)
	if err != nil {
		return nil, err
	}

	var resEntity []*entity.Category
	for _, v := range result {
		resEntity = append(resEntity, v.convertToEntityCategory(ctx))
	}
	return resEntity, nil
}

func (repo *charRepoImpl) SelectByTitle(ctx context.Context, title string) (*entity.Category, error) {
	var result categoryDB

	query := "SELECT * FROM shop.categories WHERE title = $1"

	err := repo.getReadConnection().Select(&result, query, title)
	if err != nil {
		return nil, err
	}

	return result.convertToEntityCategory(ctx), nil
}

func (repo *charRepoImpl) Insert(ctx context.Context, in entity.Category) (*entity.Category, error) {
	query := `INSERT INTO shop.categories (external_id, title, seller_id) 
            VALUES ($1, $2, $3) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, in.ExternalID, in.Title, in.SellerID)
	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *charRepoImpl) UpdateExecOne(ctx context.Context, in entity.Category) error {
	dbModel := convertToDBCategory(ctx, in)

	query := `UPDATE shop.categories SET title = $1, seller_id = $2 WHERE id = $3`
	_, err := repo.getWriteConnection().ExecOne(query, dbModel.Title, dbModel.SellerID, dbModel.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *charRepoImpl) Ping(ctx context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *charRepoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *charRepoImpl) WithTx(tx *postgresdb.Transaction) CategoryRepo {
	return &charRepoImpl{db: repo.db, tx: tx}
}

func (repo *charRepoImpl) getReadConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetReadConnection()
}

func (repo *charRepoImpl) getWriteConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetWriteConnection()
}
