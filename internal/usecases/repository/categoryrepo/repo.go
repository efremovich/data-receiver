package categoryrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type CategoryRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Category, error)
	SelectByCardID(ctx context.Context, CardID int64) ([]*entity.Category, error)
	SelectByPriceID(ctx context.Context, CardID int64) ([]*entity.Category, error)
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

	query := "SELECT * FROM categories WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}
	return result.convertToEntityCategory(ctx), nil
}

func (repo *charRepoImpl) SelectByCardID(ctx context.Context, cardID int64) ([]*entity.Category, error) {
	var result []categoryDB

	query := "SELECT * FROM categories WHERE card_id = $1"

	err := repo.getReadConnection().Select(&result, query, cardID)
	if err != nil {
		return nil, err
	}

	var resEntity []*entity.Category
	for _, v := range result {
		resEntity = append(resEntity, v.convertToEntityCategory(ctx))
	}
	return resEntity, nil
}

func (repo *charRepoImpl) SelectByPriceID(ctx context.Context, cardID int64) ([]*entity.Category, error) {
	var result []categoryDB

	query := "SELECT * FROM categories WHERE price_id = $1"

	err := repo.getReadConnection().Select(&result, query, cardID)
	if err != nil {
		return nil, err
	}

	var resEntity []*entity.Category
	for _, v := range result {
		resEntity = append(resEntity, v.convertToEntityCategory(ctx))
	}
	return resEntity, nil
}

func (repo *charRepoImpl) Insert(ctx context.Context, in entity.Category) (*entity.Category, error) {
	query := `INSERT INTO categories (card_id, title, seller_id) 
            VALUES ($1, $2, $3) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, in.CardID, in.Title, in.SellerID)
	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *charRepoImpl) UpdateExecOne(ctx context.Context, in entity.Category) error {
	dbModel := convertToDBCategory(ctx, in)

	query := `UPDATE categorys SET card_id = $1, title = $2, seller_id = $3 WHERE id = $4`
	_, err := repo.getWriteConnection().ExecOne(query, dbModel.CardID, dbModel.Title, dbModel.SellerID, dbModel.ID)
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
