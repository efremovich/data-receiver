package offerfeedrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

var ErrObjectNotFound = entity.ErrObjectNotFound

type OfferRepo interface {
	GetOffers(ctx context.Context) ([]*entity.Offer, error)
	GetStocks(ctx context.Context) (*entity.Inventory, error)
	GetCardsVkFeed(ctx context.Context, params entity.VkCardsFeedParams) ([]*entity.VKCard, error)
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

func (repo *offerRepoImpl) GetStocks(ctx context.Context) (*entity.Inventory, error) {
	var (
		stockDB   []stockDB
		storageDB []storageDB
	)

	result := entity.Inventory{}

	stockQuery := `
            with ranked_stocks as (
              select
                row_number() over (partition by s.card_id,	s.quantity,	w.seller_id order by created_at desc) as rn,
                s.card_id as id,
                s.quantity,
                ps.price,
                ps.special_price as old_price,
                w.seller_id,
                w.id as storage_id
              from
                shop.stocks s
              left join shop.warehouses w on
                w.id = s.warehouse_id
              left join shop.price_sizes ps on
                ps.id = s.barcode_id
              where
                created_at <= NOW()
            )
            select id, quantity, price, old_price, storage_id, seller_id from ranked_stocks rn
            where rn.rn = 1
limit 100
  `
	storageQuery := `
            select
              w.seller_id,
              seller.title as seller_name,
              w.id as id,
              w.name as name,
              '' as city,
              wt.name as type,
              w.address as address,
              '' as lat,
              '' as lon,
              '' as region,
              '' as work_time,
              '' as phone,
              '' as icon
            from
              shop.warehouses w
            left join shop.sellers seller on
              seller.id = w.seller_id
            left join shop.warehouse_types wt on
              wt.id = w.id
  limit 100
  `

	err := repo.getReadConnection().Select(&stockDB, stockQuery)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при получении данных фида предложений: %w", err)
	}

	for _, resStock := range stockDB {
		stock := resStock.ConvertToEntityStock(ctx)
		result.Availability = append(result.Availability, stock)
	}

	err = repo.getReadConnection().Select(&storageDB, storageQuery)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при получении данных фида предложений: %w", err)
	}

	for _, resStorage := range storageDB {
		storage := resStorage.ConvertToEntityStorage(ctx)
		result.Storages = append(result.Storages, storage)
	}

	return &result, nil
}

func (repo *offerRepoImpl) GetOffers(ctx context.Context) ([]*entity.Offer, error) {
	var results []offerDB
	query := `
  with ranked_stocks as (
    select
      s.card_id,
      s.created_at,
      s.quantity,
      w.seller_id,
      b2.barcode,
      ps.price,
      ps.special_price,
      row_number() over (partition by s.card_id,
      s.quantity,
      w.seller_id
    order by
      created_at desc) as rn
    from
      shop.stocks s
    left join shop.warehouses w on
      w.id = s.warehouse_id
    left join shop.barcodes b2 on
      b2.id = s.barcode_id
    left join shop.price_sizes ps on
      ps.id = s.barcode_id
    where
      created_at <= NOW()
    )
    select
      c.id,
      c.vendor_id as group_id,
      case
        when 
      s.quantity != 0
        or s.quantity is not null
        and s.card_id = c.id
        and s.seller_id = sc.seller_id 
    then true
        else false
      end as available,
      c.vendor_code as vendor_code,
      c.title as name,
      array_agg(distinct sc.external_id) as market_id,
      b.title as vendor,
      array_agg(distinct mf.link) as picture,
      array_agg(distinct cc2.category_id) as category_id,
      c.vendor_id as similar,
      s.barcode,
      s.price,
      s.special_price as old_price,
      c.description
    from
      shop.cards c
    left join shop.card_categories cc on
      c.id = cc.card_id
    left join shop.categories c2 on
      cc.category_id = c2.id
    left join shop.seller2cards sc on
      sc.card_id = c.id
    left join shop.brands b on
      b.id = c.brand_id
    left join shop.media_files mf on
      mf.card_id = c.id
      and mf.type_id = 1
    left join ranked_stocks s on
      s.card_id = c.id
      and s.rn = 1
    left join shop.card_categories cc2 on
      cc2.card_id = c.id
    group by
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
      s.special_price
      Limit 100` // TODO: Убрать лимит. Для теста
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

func (repo *offerRepoImpl) GetCardsVkFeed(ctx context.Context, params entity.VkCardsFeedParams) ([]*entity.VKCard, error) {
	const limit = 1000
	conditions := []string{}
	whereCondition := ""

	if params.Limit == 0 {
		params.Limit = limit
	}

	limitCondititon := fmt.Sprintf("LIMIT %d", params.Limit)

	if params.Cursor != "" {
		conditions = append(conditions, fmt.Sprintf("vendor_code > '%s'", params.Cursor))
	}

	if params.Filter != "" {
		filters := strings.Split(params.Filter, ",")
		if len(filters) > 0 {
			conditions = append(conditions, fmt.Sprintf("vendor_code IN ('%s')", strings.Join(filters, "','")))
		}
	}

	if len(conditions) > 0 {
		whereCondition = " WHERE " + strings.Join(conditions, " AND ")
	}

	var results []vkFeedDB
	query := fmt.Sprintf(`
          WITH filtered_brands AS (
            SELECT id, title
            FROM shop.brands
            WHERE title ~ 'LARETTO|LRTT'
          )
          SELECT 
            card.vendor_code AS code,
            c.title AS subject,
            char_color.value AS color,
            card.title,
            char_gender.value AS gender,
            card.description,
            array_agg(DISTINCT mf.link) AS media_links,
            MAX(ps.price) AS price, 
            COALESCE(MIN(s."name")::text, '') || ' - ' || COALESCE(MAX(s."name")::text, '') AS size,
            sc.external_id as external_id,
            COALESCE(sl.title, '') as seller_name
          FROM shop.cards card
          JOIN filtered_brands b ON b.id = card.brand_id
          LEFT JOIN shop.cards_characteristics char_color ON char_color.card_id = card.id AND char_color.characteristic_id = 11
          LEFT JOIN shop.cards_characteristics char_gender ON char_gender.card_id = card.id AND char_gender.characteristic_id = 8
          LEFT JOIN shop.media_files mf ON mf.card_id = card.id AND mf.type_id = 1
          LEFT JOIN shop.card_categories cc ON cc.card_id = card.id 
          LEFT JOIN shop.categories c ON c.id = cc.category_id
          LEFT JOIN shop.price_sizes ps ON ps.card_id = card.id
          LEFT JOIN shop.sizes s ON s.id = ps.size_id 
          LEFT JOIN shop.seller2cards sc ON sc.card_id = card.id
          LEFT JOIN shop.sellers sl on sl.id = sc.seller_id 
          %s
          GROUP BY card.vendor_code, c.title, char_color.value, card.title, char_gender.value, card.description, b.title, sc.external_id, sc.seller_id, sl.title 
          ORDER BY code ASC 
          %s`, whereCondition, limitCondititon)

	err := repo.getReadConnection().Select(&results, query)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, fmt.Errorf("ошибка при получении данных фида предложений: %w", err)
	}

	var vkCards []*entity.VKCard
	for _, result := range results {
		vkCard := result.ConvertToEntityVKCard(ctx)
		vkCards = append(vkCards, vkCard)
	}
	return vkCards, nil
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
