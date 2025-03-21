-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.brands (
	id serial NOT NULL,
	title text NOT NULL,
	seller_id integer NOT NULL
);

ALTER TABLE shop.brands OWNER TO erp_db_usr;

ALTER TABLE ONLY shop.brands
	ADD CONSTRAINT brands_pkey PRIMARY KEY (id);

ALTER TABLE ONLY shop.brands
	ADD CONSTRAINT brands_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES shop.sellers (id);

CREATE INDEX brands_seller_id_idx ON shop.brands USING btree (seller_id);

COMMENT ON TABLE shop.brands IS 'Бренды';

COMMENT ON COLUMN shop.brands.id IS 'Идентификатор';

COMMENT ON COLUMN shop.brands.title IS 'Наименование бренда';

COMMENT ON COLUMN shop.brands.seller_id IS 'Идентификатор продавца';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX brands_seller_id_idx;

DROP TABLE shop.brands;

-- +goose StatementEnd
