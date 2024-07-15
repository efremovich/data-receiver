-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.brands (
    brand_id serial NOT NULL,
    title text NOT NULL,
    seller_id integer NOT NULL
);
ALTER TABLE shop_dev.brands OWNER TO shop_user_rw;
ALTER TABLE ONLY shop_dev.brands
    ADD CONSTRAINT brands_pkey PRIMARY KEY (brand_id);
ALTER TABLE ONLY shop_dev.brands
    ADD CONSTRAINT brands_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES shop_dev.sellers(seller_id);

CREATE INDEX brands_seller_id_idx ON shop_dev.brands USING btree (seller_id);

COMMENT ON TABLE shop_dev.brands IS 'Бренды';
COMMENT ON COLUMN shop_dev.brands.brand_id IS 'Идентификатор';
COMMENT ON COLUMN shop_dev.brands.title IS 'Наименование бренда';
COMMENT ON COLUMN shop_dev.brands.seller_id IS 'Идентификатор продавца';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX brands_seller_id_idx;
DROP TABLE shop_dev.brands;
-- +goose StatementEnd
