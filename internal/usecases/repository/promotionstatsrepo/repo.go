package promotionstatsrepo

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

type PromotionStatsRepo interface {
	SelectByPromotionID(ctx context.Context, promotionID, cardID, appType int64) (*entity.PromotionStats, error)
	Insert(ctx context.Context, promotion entity.PromotionStats) (*entity.PromotionStats, error)
	Update(ctx context.Context, promotion entity.PromotionStats) error
	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) PromotionStatsRepo
}

type promotionRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewPromotionStatsRepo(_ context.Context, db *postgresdb.DBConnection) (PromotionStatsRepo, error) {
	return &promotionRepoImpl{db: db}, nil
}

func (r *promotionRepoImpl) SelectByPromotionID(ctx context.Context, promotionID, cardID, appType int64) (*entity.PromotionStats, error) {
	var result promotionStatsDB

	query := `
		SELECT "date", "views", clicks, ctr, cpc, spent, orders, cr, shks, order_amount, app_type, promotion_id, card_id
		FROM shop.promotion_stats
		WHERE promotion_id = $1 and card_id = $2 and app_type = $3`

	err := r.getReadConnection().Get(&result, query, promotionID, cardID, appType)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по id %d в таблице Promotion: %w", promotionID, err)
	}

	return result.convertToEntityPromotionStats(ctx), nil
}

func (r *promotionRepoImpl) Insert(_ context.Context, promotion entity.PromotionStats) (*entity.PromotionStats, error) {
	query := `
		INSERT INTO shop.promotion_stats
			("date", "views", clicks, ctr, cpc, spent, orders, cr, shks, order_amount, app_type, promotion_id, card_id, seller_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`

	promotionIDWrap := repository.IDWrapper{}
	err := r.getWriteConnection().QueryAndScan(&promotionIDWrap, query,
		promotion.Date,
		promotion.Views,
		promotion.Clicks,
		promotion.CTR,
		promotion.CPC,
		promotion.Spent,
		promotion.Orders,
		promotion.CR,
		promotion.SHKs,
		promotion.OrderAmount,
		promotion.AppType,
		promotion.PromotionID,
		promotion.CardID,
		promotion.SellerID,
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка при вставке данных в таблицу Promotion: %w", err)
	}

	promotion.ID = promotionIDWrap.ID.Int64

	return &promotion, nil
}

func (r *promotionRepoImpl) Update(ctx context.Context, promotionStats entity.PromotionStats) error {
	dbModel := convertToDBPromotionStats(ctx, promotionStats)
	query := `
		UPDATE shop.promotion_stats
		SET views=$1, clicks=$2, ctr=$3, cpc=$4, spent=$5, orders=$6, cr=$7, shks=$8, order_amount=$9, app_type=$10, promotion_id=$11, card_id=$12, date = $13
		WHERE id = $14`

	_, err := r.getReadConnection().ExecOne(query,
		dbModel.Views,
		dbModel.Clicks,
		dbModel.CTR,
		dbModel.CPC,
		dbModel.Spent,
		dbModel.Orders,
		dbModel.CR,
		dbModel.SHKs,
		dbModel.OrderAmount,
		dbModel.AppType,
		dbModel.PromotionID,
		dbModel.CardID,
		dbModel.Date,
		dbModel.ID,
		dbModel.SellerID,
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

func (r *promotionRepoImpl) WithTx(tx *postgresdb.Transaction) PromotionStatsRepo {
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
