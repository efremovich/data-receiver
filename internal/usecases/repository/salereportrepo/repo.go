package salereportrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

var ErrObjectNotFound = entity.ErrObjectNotFound

type SaleReportRepo interface {
	SelectByExternalID(ctx context.Context, externalID string, date time.Time) (*entity.SaleReport, error)
	Insert(ctx context.Context, in entity.SaleReport) (*entity.SaleReport, error)
	UpdateExecOne(ctx context.Context, in *entity.SaleReport) error

	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) SaleReportRepo
}

type repoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewSaleReportRepo(_ context.Context, db *postgresdb.DBConnection) (SaleReportRepo, error) {
	return &repoImpl{db: db}, nil
}

func (repo *repoImpl) SelectByExternalID(ctx context.Context, externalID string, date time.Time) (*entity.SaleReport, error) {
	var result saleReportDB

	query := `
            select
              id,
              external_id,
              updated_at,
              quantity,
              retail_price,
              return_amoun,
              sale_percent,
              commission_percent,
              retail_price_withdisc_rub,
              delivery_amount,
              return_amount,
              delivery_cost,
              pvz_reward,
              seller_reward,
              seller_reward_with_nds,
              date_from,
              date_to,
              create_report_date,
              order_date,
              sale_date,
              transaction_date,
              sa_name,
              bonus_type_name,
              penalty,
              additional_payment,
              acquiring_fee,
              acquiring_percent,
              acquiring_bank,
              doc_type,
              supplier_oper_name,
              site_country,
              kiz,
              storage_fee,
              deduction,
              acceptance,
              pvz_id,
              barcode,
              size_id,
              card_id,
              order_id,
              warehouse_id,
              seller_id
            from
              shop.sale_reports sr
            where
              external_id = $1
              and sale_date = $2
  `
	err := repo.getReadConnection().Get(&result, query, externalID, date.Format("2006-01-02 15:04:05"))
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при поиске данных по ID %s в таблице sale: %w", externalID, err)
	}

	return result.convertToEntitySaleReport(ctx), nil
}

func (repo *repoImpl) Insert(ctx context.Context, in entity.SaleReport) (*entity.SaleReport, error) {
	dbModel := convertDBToSaleReport(ctx, &in)
	dbModel.UpdatedAt = time.Now()
	query := `
    INSERT INTO shop.sale_reports (
        external_id,
        updated_at,
        quantity,
        retail_price,
        return_amoun,
        sale_percent,
        commission_percent,
        retail_price_withdisc_rub,
        delivery_amount,
        return_amount,
        delivery_cost,
        pvz_reward,
        seller_reward,
        seller_reward_with_nds,
        date_from,
        date_to,
        create_report_date,
        order_date,
        sale_date,
        transaction_date,
        sa_name,
        bonus_type_name,
        penalty,
        additional_payment,
        acquiring_fee,
        acquiring_percent,
        acquiring_bank,
        doc_type,
        supplier_oper_name,
        site_country,
        kiz,
        storage_fee,
        deduction,
        acceptance,
        pvz_id,
        barcode,
        size_id,
        card_id,
        order_id,
        warehouse_id,
        seller_id
    )
    VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,$41) RETURNING id;
  `
	charIDWrap := repository.IDWrapper{}
	err := repo.getWriteConnection().QueryAndScan(&charIDWrap, query,
		dbModel.ExternalID,
		dbModel.UpdatedAt,
		dbModel.Quantity,
		dbModel.RetailPrice,
		dbModel.ReturnAmoun,
		dbModel.SalePercent,
		dbModel.CommissionPercent,
		dbModel.RetailPriceWithdiscRub,
		dbModel.DeliveryAmount,
		dbModel.ReturnAmount,
		dbModel.DeliveryCost,
		dbModel.PvzReward,
		dbModel.SellerReward,
		dbModel.SellerRewardWithNds,
		dbModel.DateFrom,
		dbModel.DateTo,
		dbModel.CreateReportDate,
		dbModel.OrderDate,
		dbModel.SaleDate,
		dbModel.TransactionDate,
		dbModel.SAName,
		dbModel.BonusTypeName,
		dbModel.Penalty,
		dbModel.AdditionalPayment,
		dbModel.AcquiringFee,
		dbModel.AcquiringPercent,
		dbModel.AcquiringBank,
		dbModel.DocType,
		dbModel.SupplierOperName,
		dbModel.SiteCountry,
		dbModel.KIZ,
		dbModel.StorageFee,
		dbModel.Deduction,
		dbModel.Acceptance,
		dbModel.PvzID,
		dbModel.Barcode,
		dbModel.SizeID,
		dbModel.CardID,
		dbModel.OrderID,
		dbModel.WarehouseID,
		dbModel.SellerID,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка при вставке данных в таблицу sale_reports: %w", err)
	}

	in.ID = charIDWrap.ID.Int64

	return dbModel.convertToEntitySaleReport(ctx), nil
}

func (repo *repoImpl) UpdateExecOne(ctx context.Context, in *entity.SaleReport) error {
	dbModel := convertDBToSaleReport(ctx, in)
	query := `
          UPDATE shop.sale_reports
          SET 
              id = $1,
              external_id = $2,
              updated_at = $3,
              quantity = $4,
              retail_price = $5,
              return_amoun = $6,
              sale_percent = $7,
              commission_percent = $8,
              retail_price_withdisc_rub = $9,
              delivery_amount = $10,
              return_amount = $11,
              delivery_cost = $12,
              pvz_reward = $13,
              seller_reward = $14,
              seller_reward_with_nds = $15,
              date_from = $16,
              date_to = $17,
              create_report_date = $18,
              order_date = $19,
              sale_date = $20,
              transaction_date = $21,
              sa_name = $22,
              bonus_type_name = $23,
              penalty = $24,
              additional_payment = $25,
              acquiring_fee = $26,
              acquiring_percent = $27,
              acquiring_bank = $28,
              doc_type = $29,
              supplier_oper_name = $30,
              site_country = $31,
              kiz = $32,
              storage_fee = $33,
              deduction = $34,
              acceptance = $35,
              pvz_id = $36,
              barcode = $37,
              size_id = $38,
              card_id = $39,
              order_id = $40,
              warehouse_id = $41,
              seller_id = $42
          WHERE 
              id = $43;
  `

	_, err := repo.getWriteConnection().Exec(query,
		dbModel.ID,
		dbModel.ExternalID,
		time.Now(),
		dbModel.Quantity,
		dbModel.RetailPrice,
		dbModel.ReturnAmoun,
		dbModel.SalePercent,
		dbModel.CommissionPercent,
		dbModel.RetailPriceWithdiscRub,
		dbModel.DeliveryAmount,
		dbModel.ReturnAmount,
		dbModel.DeliveryCost,
		dbModel.PvzReward,
		dbModel.SellerReward,
		dbModel.SellerRewardWithNds,
		dbModel.DateFrom,
		dbModel.DateTo,
		dbModel.CreateReportDate,
		dbModel.OrderDate,
		dbModel.SaleDate,
		dbModel.TransactionDate,
		dbModel.SAName,
		dbModel.BonusTypeName,
		dbModel.Penalty,
		dbModel.AdditionalPayment,
		dbModel.AcquiringFee,
		dbModel.AcquiringPercent,
		dbModel.AcquiringBank,
		dbModel.DocType,
		dbModel.SupplierOperName,
		dbModel.SiteCountry,
		dbModel.KIZ,
		dbModel.StorageFee,
		dbModel.Deduction,
		dbModel.Acceptance,
		dbModel.PvzID,
		dbModel.Barcode,
		dbModel.SizeID,
		dbModel.CardID,
		dbModel.OrderID,
		dbModel.WarehouseID,
		dbModel.SellerID,
		dbModel.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *repoImpl) Ping(_ context.Context) error {
	return repo.getWriteConnection().Ping()
}

func (repo *repoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *repoImpl) WithTx(tx *postgresdb.Transaction) SaleReportRepo {
	return &repoImpl{db: repo.db, tx: tx}
}

func (repo *repoImpl) getReadConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetReadConnection()
}

func (repo *repoImpl) getWriteConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetWriteConnection()
}
