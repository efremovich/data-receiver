-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.categories (
    id serial NOT NULL,
    title text NOT NULL,
    seller_id integer NOT NULL,
    external_id integer NOT NULL
);
ALTER TABLE shop.categories OWNER TO erp_db_usr;
ALTER TABLE ONLY shop.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);
ALTER TABLE ONLY shop.categories
    ADD CONSTRAINT categories_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES shop.sellers(id);

CREATE INDEX categories_seller_id_idx ON shop.categories USING btree (seller_id);

COMMENT ON TABLE shop.categories IS 'Категории товаров';
COMMENT ON COLUMN shop.categories.id IS 'Идентификатор';
COMMENT ON COLUMN shop.categories.title IS 'Наименование категории';
COMMENT ON COLUMN shop.categories.seller_id IS 'Идентификатор продавца';

CREATE TABLE shop.card_categories (
  id serial NOT NULL,
  card_id integer NOT NULL,
  category_id integer NOT NULL
);

ALTER TABLE shop.card_categories OWNER TO erp_db_usr;
ALTER TABLE ONLY shop.card_categories
    ADD CONSTRAINT card_categories_pkey PRIMARY KEY (id);
ALTER TABLE ONLY shop.card_categories
    ADD CONSTRAINT card_categories_card_id_fkey FOREIGN KEY (card_id) REFERENCES shop.cards(id);
ALTER TABLE ONLY shop.card_categories
    ADD CONSTRAINT card_categories_category_id_fkey FOREIGN KEY (category_id) REFERENCES shop.categories(id);

CREATE INDEX card_categories_card_id_idx ON shop.card_categories USING btree (card_id);
CREATE INDEX card_categories_category_id_idx ON shop.card_categories USING btree (category_id);

COMMENT ON TABLE shop.card_categories IS 'Категории карточки';
COMMENT ON COLUMN shop.card_categories.id IS 'Идентификатор категории';
COMMENT ON COLUMN shop.card_categories.card_id IS 'Идентификатор карточки';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX categories_seller_id_idx;
DROP TABLE shop.categories;

DROP INDEX card_categories_card_id_idx;
DROP INDEX card_categories_category_id_idx;
DROP TABLE shop.card_categories;
-- +goose StatementEnd



