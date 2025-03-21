-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.barcodes (
	id serial NOT NULL,
	barcode text NOT NULL,
	seller_id integer NOT NULL,
	price_size_id integer NOT NULL
);

ALTER TABLE shop.barcodes OWNER TO erp_db_usr;

ALTER TABLE ONLY shop.barcodes
	ADD CONSTRAINT barcodes_pkey PRIMARY KEY (id);

ALTER TABLE ONLY shop.barcodes
	ADD CONSTRAINT barcodes_barcode_key UNIQUE (barcode);

ALTER TABLE ONLY shop.barcodes
	ADD CONSTRAINT barcodes_price_size_id_fkey FOREIGN KEY (price_size_id) REFERENCES shop.price_sizes (id);

ALTER TABLE ONLY shop.barcodes
	ADD CONSTRAINT barcodes_seller_fkey FOREIGN KEY (seller_id) REFERENCES shop.sellers (id);

CREATE INDEX barcodes_seller_id_idx ON shop.barcodes USING btree (seller_id);

COMMENT ON TABLE shop.barcodes IS 'Штрихкоды';

COMMENT ON COLUMN shop.barcodes.id IS 'Уникальный идентификатор';

COMMENT ON COLUMN shop.barcodes.barcode IS 'Штрихкод';

COMMENT ON COLUMN shop.barcodes.seller_id IS 'Идентификатор продавца';

COMMENT ON COLUMN shop.barcodes.price_size_id IS 'Размер';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX barcodes_seller_id_idx;

DROP TABLE shop.barcodes;

-- +goose StatementEnd
