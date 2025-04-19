package promotionrepo

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

type PromotionRepo interface {
	SelectByExternalID(ctx context.Context, externalID int64) (*entity.Promotion, error)
	Insert(ctx context.Context, promotion entity.Promotion) (*entity.Promotion, error)
	Update(ctx context.Context, promotion entity.Promotion) error
	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(tx *postgresdb.Transaction) PromotionRepo
}

type promotionRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewPromotionRepo(_ context.Context, db *postgresdb.DBConnection) (PromotionRepo, error) {
	return &promotionRepoImpl{db: db}, nil
}

func (r *promotionRepoImpl) SelectByExternalID(ctx context.Context, externalID int64) (*entity.Promotion, error) {
	var result promotionDB

	query := `
		SELECT id, external_id, name, type, status, change_time, create_time, 
		       date_start, date_end, views, clicks, ctr, cpc, spent, 
		       orders, cr, shks, order_amount, seller_id
		FROM shop.promotions 
		WHERE external_id = $1`

	err := r.getReadConnection().Get(&result, query, externalID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по id %d в таблице Promotion: %w", externalID, err)
	}

	return result.convertToEntityPromotion(ctx), nil
}

func (r *promotionRepoImpl) Insert(_ context.Context, promotion entity.Promotion) (*entity.Promotion, error) {
	query := `
		INSERT INTO shop.promotions (
			external_id, name, type, status, change_time, create_time, 
			date_start, date_end, views, clicks, ctr, cpc, spent, 
			orders, cr, shks, order_amount, seller_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		) RETURNING id`

	promotionIDWrap := repository.IDWrapper{}
	err := r.getWriteConnection().QueryAndScan(&promotionIDWrap, query,
		promotion.ExternalID,
		promotion.Name,
		promotion.Type,
		promotion.Status,
		promotion.ChangeTime,
		promotion.CreateTime,
		promotion.DateStart,
		promotion.DateEnd,
		promotion.Views,
		promotion.Clicks,
		promotion.CTR,
		promotion.CPC,
		promotion.Spent,
		promotion.Orders,
		promotion.CR,
		promotion.SHKs,
		promotion.OrderAmount,
		promotion.SellerID,
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка при вставке данных в таблицу Promotion: %w", err)
	}

	promotion.ID = promotionIDWrap.ID.Int64

	return &promotion, nil
}

func (r *promotionRepoImpl) Update(ctx context.Context, promotion entity.Promotion) error {
	dbModel := convertToDBPromotion(ctx, promotion)
	query := `
		UPDATE promotions SET
			external_id = $1,
			name = $2,
			type = $3,
			status = $4,
			change_time = $5,
			date_start = $6,
			date_end = $7,
			views = $8,
			clicks = $9,
			ctr = $10,
			cpc = $11,
			spent = $12,
			orders = $13,
			cr = $14,
			shks = $15,
			order_amount = $16
			seller_id = $17,
		WHERE id = $18`

	_, err := r.getReadConnection().ExecOne(query,
		dbModel.ExternalID,
		dbModel.Name,
		dbModel.Type,
		dbModel.Status,
		dbModel.ChangeTime,
		dbModel.DateStart,
		dbModel.DateEnd,
		dbModel.Views,
		dbModel.Clicks,
		dbModel.CTR,
		dbModel.CPC,
		dbModel.Spent,
		dbModel.Orders,
		dbModel.CR,
		dbModel.SHKs,
		dbModel.OrderAmount,
		dbModel.SellerID,
		dbModel.ID,
	)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении данных в таблицу Promotion: %w", err)
	}

	return nil
}

func (r *promotionRepoImpl) Ping(_ context.Context) error {
	return r.getReadConnection().Ping()
}

func (r *promotionRepoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return r.db.GetReadConnection().BeginTX(ctx)
}

func (r *promotionRepoImpl) WithTx(tx *postgresdb.Transaction) PromotionRepo {
	return &promotionRepoImpl{db: r.db, tx: tx}
}

func (r *promotionRepoImpl) getReadConnection() postgresdb.QueryExecutor {
	if r.tx != nil {
		return *r.tx
	}

	return r.db.GetReadConnection()
}

func (r *promotionRepoImpl) getWriteConnection() postgresdb.QueryExecutor {
	if r.tx != nil {
		return *r.tx
	}

	return r.db.GetWriteConnection()
}
