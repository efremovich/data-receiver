package sizerepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type SizeRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Size, error)
	SelectByCardID(ctx context.Context, CardID int64) ([]*entity.Size, error)
	SelectByPriceID(ctx context.Context, CardID int64) ([]*entity.Size, error)
	Insert(ctx context.Context, in entity.Size) (*entity.Size, error)
	UpdateExecOne(ctx context.Context, in entity.Size) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) SizeRepo
}

type charRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewSizeRepo(_ context.Context, db *postgresdb.DBConnection) (SizeRepo, error) {
	return &charRepoImpl{db: db}, nil
}

func (repo *charRepoImpl) SelectByID(ctx context.Context, id int64) (*entity.Size, error) {
	var result sizeDB

	query := "SELECT * FROM sizes WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}
	return result.convertToEntitySize(ctx), nil
}

func (repo *charRepoImpl) SelectByCardID(ctx context.Context, cardID int64) ([]*entity.Size, error) {
	var result []sizeDB

	query := "SELECT * FROM sizes WHERE card_id = $1"

	err := repo.getReadConnection().Select(&result, query, cardID)
	if err != nil {
		return nil, err
	}

	var resEntity []*entity.Size
	for _, v := range result {
		resEntity = append(resEntity, v.convertToEntitySize(ctx))
	}
	return resEntity, nil
}

func (repo *charRepoImpl) SelectByPriceID(ctx context.Context, cardID int64) ([]*entity.Size, error) {
	var result []sizeDB

	query := "SELECT * FROM sizes WHERE price_id = $1"

	err := repo.getReadConnection().Select(&result, query, cardID)
	if err != nil {
		return nil, err
	}

	var resEntity []*entity.Size
	for _, v := range result {
		resEntity = append(resEntity, v.convertToEntitySize(ctx))
	}
	return resEntity, nil
}

func (repo *charRepoImpl) Insert(ctx context.Context, in entity.Size) (*entity.Size, error) {
	query := `INSERT INTO sizes (card_id, title, tech_size, price_id) 
            VALUES ($1, $2, $3, $4) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, in.CardID, in.Title, in.TechSize, in.PriceID)
	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *charRepoImpl) UpdateExecOne(ctx context.Context, in entity.Size) error {
	dbModel := convertToDBSize(ctx, in)

	query := `UPDATE sizes SET card_id = $1, title = $2, tech_size = $3, price_id= $5 WHERE id = $4`
	_, err := repo.getWriteConnection().ExecOne(query, dbModel.CardID, dbModel.Title, dbModel.TechSize, dbModel.ID)
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

func (repo *charRepoImpl) WithTx(tx *postgresdb.Transaction) SizeRepo {
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
