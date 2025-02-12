package offerfeedrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

var ErrObjectNotFound = entity.ErrObjectNotFound

type OfferRepo interface {
	GetOffers(ctx context.Context) ([]*entity.Offer, error)
	GetStocks(ctx context.Context) ([]*entity.Inventory, error)
	Ping(ctx context.Context) error
	BeginTX(ctx context.Context) (postgresdb.Transaction, error)
	WithTx(*postgresdb.Transaction) OfferRepo
}

type offerRepoImpl struct {
	db *postgresdb.DBConnection
	tx *postgresdb.Transaction
}

func NewOfferRepo(_ context.Context, db *postgresdb.DBConnection) (OfferRepo, error) {
	return &offerRepoImpl{db: db}, nil
}

func (repo *offerRepoImpl) GetStocks(ctx context.Context) ([]*entity.Inventory, error) {
}
func (repo *offerRepoImpl) GetOffers(ctx context.Context) ([]*entity.Offer, error) {
	var results []offerDB
	query := `
  WITH ranked_stocks AS (
    SELECT
        s.card_id,
        s.created_at,
        s.quantity,
        w.seller_id,
        b2.barcode,
        ps.price,
        ps.special_price,
        ROW_NUMBER() OVER (PARTITION BY s.card_id, s.quantity, w.seller_id ORDER BY created_at DESC) AS rn
    FROM
        shop.stocks s
    LEFT JOIN shop.warehouses w ON w.id = s.warehouse_id
    LEFT JOIN shop.barcodes b2 ON b2.id = s.barcode_id
    LEFT JOIN shop.price_sizes ps ON ps.id = s.barcode_id
    WHERE
        created_at <= NOW()
),
ranked_pictures AS (
    SELECT
        mf.card_id,
        mf.link,
        ROW_NUMBER() OVER (PARTITION BY mf.card_id ORDER BY mf.id) AS rn
    FROM
        shop.media_files mf
    WHERE
        mf.type_id = 1
        and mf.link ~ '1.webp'
)
SELECT
    c.id,
    c.vendor_id AS group_id,
    CASE
        WHEN s.quantity != 0 OR s.quantity IS NOT NULL AND s.card_id = c.id AND s.seller_id = sc.seller_id THEN true
        ELSE false
    END AS available,
    c.vendor_code AS vendor_code,
    c.title AS name,
    ARRAY_AGG(DISTINCT sc.external_id) AS market_id,
    b.title AS vendor,
    ARRAY_AGG(rp.link) AS picture, -- Используем подзапрос для одной картинки
    ARRAY_AGG(DISTINCT cc2.category_id) AS category_id,
    c.vendor_id AS similar,
    s.barcode,
    s.price,
    s.special_price AS old_price,
    c.description
FROM
    shop.cards c
LEFT JOIN shop.card_categories cc ON c.id = cc.card_id
LEFT JOIN shop.categories c2 ON cc.category_id = c2.id
LEFT JOIN shop.seller2cards sc ON sc.card_id = c.id
LEFT JOIN shop.brands b ON b.id = c.brand_id
LEFT JOIN ranked_pictures rp ON rp.card_id = c.id AND rp.rn = 1 -- Выбираем только первую картинку
LEFT JOIN ranked_stocks s ON s.card_id = c.id AND s.rn = 1
LEFT JOIN shop.card_categories cc2 ON cc2.card_id = c.id
GROUP BY
    c.id,
    s.quantity,
    b.title,
    s.card_id,
    sc.external_id,
    s.seller_id,
    sc.id,
    cc2.category_id,
    s.barcode,
    s.price,
    s.special_price;`
	err := repo.getReadConnection().Select(&results, query)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при получении данных фида предложений: %w", err)
	}

	var offers []*entity.Offer
	for _, result := range results {
		offer := result.ConvertToEntityOffer(ctx)
		offers = append(offers, offer)
	}

	return offers, nil
}

func (repo *offerRepoImpl) Ping(_ context.Context) error {
	return repo.getReadConnection().Ping()
}

func (repo *offerRepoImpl) BeginTX(ctx context.Context) (postgresdb.Transaction, error) {
	return repo.db.GetReadConnection().BeginTX(ctx)
}

func (repo *offerRepoImpl) WithTx(tx *postgresdb.Transaction) OfferRepo {
	return &offerRepoImpl{db: repo.db, tx: tx}
}

func (repo *offerRepoImpl) getReadConnection() postgresdb.QueryExecutor {
	if repo.tx != nil {
		return *repo.tx
	}

	return repo.db.GetReadConnection()
}
