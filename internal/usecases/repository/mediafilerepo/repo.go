package mediafilerepo

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

type MediaFileRepo interface {
	SelectByCardID(ctx context.Context, cardID int64, link string) ([]*entity.MediaFile, error)
	Insert(ctx context.Context, in entity.MediaFile) (*entity.MediaFile, error)
	UpdateExecOne(ctx context.Context, in entity.MediaFile) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) MediaFileRepo
}

type mediafileRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewMediaFileRepo(_ context.Context, db *postgresdb.DBConnection) (MediaFileRepo, error) {
	return &mediafileRepoImpl{db: db}, nil
}

func (repo *mediafileRepoImpl) SelectByCardID(ctx context.Context, id int64, link string) ([]*entity.MediaFile, error) {
	var result []mediaFileDB

	query := "SELECT id, link, card_id, type_id FROM shop.media_files WHERE card_id = $1 and link = $2"

	err := repo.getReadConnection().Select(&result, query, id, link)
	if err != nil && errors.Is(err, sql.ErrNoRows) || len(result) == 0 {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по id %d в таблице mediaFiles: %w", id, err)
	}

	var resEntity []*entity.MediaFile
	for _, v := range result {
		resEntity = append(resEntity, v.ConvertToEntityMediaFile(ctx))
	}

	return resEntity, nil
}

func (repo *mediafileRepoImpl) Insert(_ context.Context, in entity.MediaFile) (*entity.MediaFile, error) {
	query := `INSERT INTO shop.media_files (link, card_id, type_id) 
            VALUES ($1, $2, $3) RETURNING id`
	mediafileIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&mediafileIDWrap, query, in.Link, in.CardID, in.TypeID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при вставке данных в таблицу mediaFiles: %w", err)
	}

	in.ID = mediafileIDWrap.ID.Int64

	return &in, nil
}

func (repo *mediafileRepoImpl) UpdateExecOne(ctx context.Context, in entity.MediaFile) error {
	dbModel := convertToDBMediaFile(ctx, in)

	query := `UPDATE shop.media_files SET link = $1, card_id = $2, type_id = $3  WHERE id = $4`

	_, err := repo.getWriteConnection().ExecOne(query, dbModel.Link, dbModel.CardID, dbModel.TypeID, dbModel.ID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении данных в таблицу mediaFiles: %w", err)
	}

	return nil
}

func (repo *mediafileRepoImpl) Ping(_ context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *mediafileRepoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *mediafileRepoImpl) WithTx(tx *postgresdb.Transaction) MediaFileRepo {
	return &mediafileRepoImpl{db: repo.db, tx: tx}
}

func (repo *mediafileRepoImpl) getReadConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetReadConnection()
}

func (repo *mediafileRepoImpl) getWriteConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetWriteConnection()
}
