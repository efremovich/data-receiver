-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.barcodes (
    barcode_id serial NOT NULL,
    barcode text NOT NULL,
    seller_id integer NOT NULL,
    price_size_id integer NOT NULL
);
ALTER TABLE shop_dev.barcodes OWNER TO shop_user_rw;

ALTER TABLE ONLY shop_dev.barcodes
    ADD CONSTRAINT barcodes_pkey PRIMARY KEY (barcode_id);
ALTER TABLE ONLY shop_dev.barcodes
    ADD CONSTRAINT barcodes_barcode_key UNIQUE (barcode);
ALTER TABLE ONLY shop_dev.barcodes
    ADD CONSTRAINT barcodes_price_size_id_fkey FOREIGN KEY (price_size_id) REFERENCES shop_dev.price_sizes(price_size_id);
ALTER TABLE ONLY shop_dev.barcodes
    ADD CONSTRAINT barcodes_seller_fkey FOREIGN KEY (seller_id) REFERENCES shop_dev.sellers(seller_id);

CREATE INDEX barcodes_seller_id_idx ON shop_dev.barcodes USING btree (seller_id);

COMMENT ON TABLE shop_dev.barcodes IS 'Штрихкоды';
COMMENT ON COLUMN shop_dev.barcodes.barcode_id IS 'Уникальный идентификатор';
COMMENT ON COLUMN shop_dev.barcodes.barcode IS 'Штрихкод';
COMMENT ON COLUMN shop_dev.barcodes.seller_id IS 'Идентификатор продавца';
COMMENT ON COLUMN shop_dev.barcodes.price_size_id IS 'Размер';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX barcodes_seller_id_idx;
DROP TABLE shop_dev.barcodes;
-- +goose StatementEnd
