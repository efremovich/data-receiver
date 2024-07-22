package cardcharrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type CardCharRepo interface {
	SelectByID(ctx context.Context, id int64) (*entity.CardCharacteristic, error)
	SelectByCardIDAndCharID(ctx context.Context, cardID, charID int64) (*entity.CardCharacteristic, error)
	Insert(ctx context.Context, in entity.CardCharacteristic) (*entity.CardCharacteristic, error)
	UpdateExecOne(ctx context.Context, in entity.CardCharacteristic) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) CardCharRepo
}

type cardcharRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewCharRepo(_ context.Context, db *postgresdb.DBConnection) (CardCharRepo, error) {
	return &cardcharRepoImpl{db: db}, nil
}

func (repo *cardcharRepoImpl) SelectByID(ctx context.Context, id int64) (*entity.CardCharacteristic, error) {
	var result cardCharacteristicDB

	query := `select
              cc.id,
              cc.value as value,
              cc.card_id,
              cc.characteristic_id,
              c.title
            from
              shop.cards_characteristics cc
              left join shop.characteristics c on c.id = cc.characteristic_id 
            where
              cc.id = $1`

	err := repo.getReadConnection().Get(&result, query, id)
	if err != nil {
		return nil, err
	}
	return result.ConvertToEntityCardCharacteristic(ctx), nil
}

func (repo *cardcharRepoImpl) SelectByCardIDAndCharID(ctx context.Context, cardID, charID int64) (*entity.CardCharacteristic, error) {
	var result cardCharacteristicDB
	query := `select
              cc.id,
              cc.value as value,
              cc.card_id,
              cc.characteristic_id,
              c.title
            from
              shop.cards_characteristics cc
              left join shop.characteristics c on c.id = cc.characteristic_id 
            where
              cc.card_id = $1
              and cc.characteristic_id = $2`
	err := repo.getReadConnection().Get(&result, query, cardID, charID)
	if err != nil {
		return nil, err
	}
	return result.ConvertToEntityCardCharacteristic(ctx), nil
}

func (repo *cardcharRepoImpl) Insert(ctx context.Context, in entity.CardCharacteristic) (*entity.CardCharacteristic, error) {
	query := `INSERT INTO shop.cards_characteristics (card_id, value, characteristic_id) 
            VALUES ($1, $2, $3) RETURNING id`
	charIDWrap := repository.IDWrapper{}
	dbModel := convertToDBCardCharacteristic(ctx, in)

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, dbModel.CardID, dbModel.Value, dbModel.CharacteristicID)
	if err != nil {
		return nil, err
	}
	in.ID = charIDWrap.ID.Int64
	return &in, nil
}

func (repo *cardcharRepoImpl) UpdateExecOne(ctx context.Context, in entity.CardCharacteristic) error {
	dbModel := convertToDBCardCharacteristic(ctx, in)

	query := `UPDATE shop.characteristics SET card_id = $1, value = $2, characteristic_id = $3 WHERE id = $4`
	_, err := repo.getWriteConnection().ExecOne(query, dbModel.CardID, dbModel.Value, dbModel.CharacteristicID, dbModel.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *cardcharRepoImpl) Ping(ctx context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *cardcharRepoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *cardcharRepoImpl) WithTx(tx *postgresdb.Transaction) CardCharRepo {
	return &cardcharRepoImpl{db: repo.db, tx: tx}
}

func (repo *cardcharRepoImpl) getReadConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetReadConnection()
}

func (repo *cardcharRepoImpl) getWriteConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetWriteConnection()
}
