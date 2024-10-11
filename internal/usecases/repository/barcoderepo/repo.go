package barcoderepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

type BarcodeRepo interface {
	SelectByBarcode(ctx context.Context, barcode string) (*entity.Barcode, error)
	Insert(ctx context.Context, in entity.Barcode) (*entity.Barcode, error)
	UpdateExecOne(ctx context.Context, in entity.Barcode) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) BarcodeRepo
}

type charRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewBarcodeRepo(_ context.Context, db *postgresdb.DBConnection) (BarcodeRepo, error) {
	return &charRepoImpl{db: db}, nil
}

func (repo *charRepoImpl) SelectByBarcode(ctx context.Context, barcode string) (*entity.Barcode, error) {
	var result barcodeDB

	query := "SELECT * FROM shop.barcodes WHERE barcode = $1"

	err := repo.getReadConnection().Get(&result, query, barcode)
	if err != nil {
		return nil, err
	}

	return result.convertToEntityBarcode(ctx), nil
}

func (repo *charRepoImpl) Insert(ctx context.Context, in entity.Barcode) (*entity.Barcode, error) {
	query := `INSERT INTO shop.barcodes (barcode, price_size_id, seller_id) 
            VALUES ($1, $2, $3) RETURNING id`
	charIDWrap := repository.IDWrapper{}

	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query, in.Barcode, in.PriceSizeID, in.SellerID)
	if err != nil {
		return nil, err
	}

	in.ID = charIDWrap.ID.Int64

	return &in, nil
}

func (repo *charRepoImpl) UpdateExecOne(ctx context.Context, in entity.Barcode) error {
	dbModel := convertToDBBarcode(ctx, in)

	query := `UPDATE shop.barcodes SET price_size_id = $1, seller_id = $2 WHERE barcode = $3`

	_, err := repo.getWriteConnection().ExecOne(query, dbModel.PriceSizeID, dbModel.SellerID, dbModel.Barcode)
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

func (repo *charRepoImpl) WithTx(tx *postgresdb.Transaction) BarcodeRepo {
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
