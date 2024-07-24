package sizerepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type SizeRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.Size, error)
	SelectByTitle(ctx context.Context, title string) (*entity.Size, error)
	SelectByTechSize(ctx context.Context, techSize string) (*entity.Size, error)
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

	query := "SELECT * FROM shop.sizes WHERE id = $1"

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}
	return result.convertToEntitySize(ctx), nil
}

func (repo *charRepoImpl) SelectByTitle(ctx context.Context, title string) (*entity.Size, error) {
	var result sizeDB

	query := "SELECT * FROM shop.sizes WHERE name = $1"

	err := repo.getReadConnection().Get(&result, query, title)
	if err != nil {
		return nil, err
	}
	return result.convertToEntitySize(ctx), nil
}

func (repo *charRepoImpl) SelectByTechSize(ctx context.Context, techSize string) (*entity.Size, error) {
	var result sizeDB

	query := "SELECT * FROM shop.sizes WHERE tech_size = $1"

	err := repo.getReadConnection().Get(&result, query, techSize)
	if err != nil {
		return nil, err
	}
	return result.convertToEntitySize(ctx), nil
}

func (repo *charRepoImpl) Insert(ctx context.Context, in entity.Size) (*entity.Size, error) {
	query := `INSERT INTO shop.sizes (name, tech_size) 
            VALUES ($1, $2) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, in.Title, in.TechSize)
	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *charRepoImpl) UpdateExecOne(ctx context.Context, in entity.Size) error {
	dbModel := convertToDBSize(ctx, in)

	query := `UPDATE shop.sizes SET name = $1, tech_size = $2 WHERE id = $3`
	_, err := repo.getWriteConnection().ExecOne(query, dbModel.Title, dbModel.TechSize, dbModel.ID)
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
