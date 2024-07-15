-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.categories (
    category_id serial NOT NULL,
    title text NOT NULL,
    seller_id integer NOT NULL
);
ALTER TABLE shop_dev.categories OWNER TO shop_user_rw;
ALTER TABLE ONLY shop_dev.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (category_id);
ALTER TABLE ONLY shop_dev.categories
    ADD CONSTRAINT categories_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES shop_dev.sellers(seller_id);

CREATE INDEX categories_seller_id_idx ON shop_dev.categories USING btree (seller_id);

COMMENT ON TABLE shop_dev.categories IS 'Категории товаров';
COMMENT ON COLUMN shop_dev.categories.category_id IS 'Идентификатор';
COMMENT ON COLUMN shop_dev.categories.title IS 'Наименование категории';
COMMENT ON COLUMN shop_dev.categories.seller_id IS 'Идентификатор продавца';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX categories_seller_id_idx;
DROP TABLE shop_dev.categories;
-- +goose StatementEnd
