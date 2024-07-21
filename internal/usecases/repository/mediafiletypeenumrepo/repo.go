package mediafiletypeenumrepo 

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type MediaFileTyepEnumRepo interface {
	SelectByID(ctx context.Context, typeID int64) (*entity.MediaFileTypeEnum, error)
	SelectByType(ctx context.Context, typeName string) (*entity.MediaFileTypeEnum, error)
	Insert(ctx context.Context, in entity.MediaFileTypeEnum) (*entity.MediaFileTypeEnum, error)
	UpdateExecOne(ctx context.Context, in entity.MediaFileTypeEnum) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) MediaFileTyepEnumRepo
}

type mediafileTypeEnumRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewMediaFileTypeEnumRepo(_ context.Context, db *postgresdb.DBConnection) (MediaFileTyepEnumRepo, error) {
	return &mediafileTypeEnumRepoImpl{db: db}, nil
}

func (repo *mediafileTypeEnumRepoImpl) SelectByID(ctx context.Context, mediaFileTypeEnumID int64) (*entity.MediaFileTypeEnum, error) {
	var result mediaFileTypeEnumDB

	query := `SELECT * FROM shop.media_files_types_enum WHERE id = $1`

	err := repo.getReadConnection().Get(&result, query, mediaFileTypeEnumID)
	if err != nil {
		return nil, err
	}

	return result.ConvertToEntityMediaFileTypeEnum(ctx), nil
}

func (repo *mediafileTypeEnumRepoImpl) SelectByType(ctx context.Context, typeName string) (*entity.MediaFileTypeEnum, error) {
	var result mediaFileTypeEnumDB

	query := `SELECT id, "type" FROM shop.media_files_types_enum WHERE type = $1`

	err := repo.getReadConnection().Get(&result, query, typeName)
	if err != nil {
		return nil, err
	}

	return result.ConvertToEntityMediaFileTypeEnum(ctx), nil
}

func (repo *mediafileTypeEnumRepoImpl) Insert(ctx context.Context, in entity.MediaFileTypeEnum) (*entity.MediaFileTypeEnum, error) {
	query := `INSERT INTO shop.media_files_types_enum ("type") 
            VALUES ($1) RETURNING id`
	mediafileIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&mediafileIDWrap, query, in.Type)
	if err != nil {
		return nil, err
	}
	in.ID = mediafileIDWrap.ID.Int64
	return &in, nil
}

func (repo *mediafileTypeEnumRepoImpl) UpdateExecOne(ctx context.Context, in entity.MediaFileTypeEnum) error {
	dbModel := convertToDBMediaFileTypeEnum(ctx, in)

	query := `UPDATE shop.media_files_types_enum SET "type" = $1 WHERE id = $2`
	_, err := repo.getWriteConnection().ExecOne(query,dbModel.Type,dbModel.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *mediafileTypeEnumRepoImpl) Ping(ctx context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *mediafileTypeEnumRepoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *mediafileTypeEnumRepoImpl) WithTx(tx *postgresdb.Transaction) MediaFileTyepEnumRepo {
	return &mediafileTypeEnumRepoImpl{db: repo.db, tx: tx}
}

func (repo *mediafileTypeEnumRepoImpl) getReadConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetReadConnection()
}

func (repo *mediafileTypeEnumRepoImpl) getWriteConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetWriteConnection()
}
