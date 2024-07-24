-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.categories (
    id serial NOT NULL,
    title text NOT NULL,
    seller_id integer NOT NULL,
    card_id integer NOT NULL,
    external_id integer,
    parent_id integer 
);
ALTER TABLE shop.categories OWNER TO shop_user_rw;
ALTER TABLE ONLY shop.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);
ALTER TABLE ONLY shop.categories
    ADD CONSTRAINT categories_card_id_fkey FOREIGN KEY (card_id) REFERENCES shop.cards(id);
-- ALTER TABLE ONLY shop.categories
    -- ADD CONSTRAINT categories_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES shop.sellers(id);

CREATE INDEX categories_seller_id_idx ON shop.categories USING btree (seller_id);

COMMENT ON TABLE shop.categories IS 'Категории товаров';
COMMENT ON COLUMN shop.categories.id IS 'Идентификатор';
COMMENT ON COLUMN shop.categories.title IS 'Наименование категории';
COMMENT ON COLUMN shop.categories.seller_id IS 'Идентификатор продавца';
COMMENT ON COLUMN shop.categories.card_id IS 'Идентификатор карточки товара';
COMMENT ON COLUMN shop.categories.external_id IS 'Внешний идентификатор категории продавца';
COMMENT ON COLUMN shop.categories.parent_id IS 'Родительская катетегория';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX categories_seller_id_idx;
DROP TABLE shop.categories;
-- +goose StatementEnd


