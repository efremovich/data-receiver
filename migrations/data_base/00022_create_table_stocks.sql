-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.stocks (
    stock_id serial NOT NULL,
    quantity numeric(10,0) NOT NULL,
    warehouse_id integer NOT NULL,
    card_id integer NOT NULL,
    barcode_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);
ALTER TABLE shop_dev.stocks OWNER TO shop_user_rw;

ALTER TABLE ONLY shop_dev.stocks
    ADD CONSTRAINT stocks_pkey PRIMARY KEY (stock_id);
ALTER TABLE ONLY shop_dev.stocks
    ADD CONSTRAINT stocks_card_id_fkey FOREIGN KEY (card_id) REFERENCES shop_dev.cards(card_id);
ALTER TABLE ONLY shop_dev.stocks
    ADD CONSTRAINT stocks_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES shop_dev.warehouses(warehouse_id);
ALTER TABLE ONLY shop_dev.stocks
    ADD CONSTRAINT stocks_fk FOREIGN KEY (barcode_id) REFERENCES shop_dev.barcodes(barcode_id);

CREATE INDEX stocks_barcode_idx ON shop_dev.stocks USING btree (barcode_id);
CREATE INDEX stocks_card_id_idx ON shop_dev.stocks USING btree (card_id);
CREATE INDEX stocks_created_at_idx ON shop_dev.stocks USING btree (created_at);
CREATE INDEX stocks_updated_at_idx ON shop_dev.stocks USING btree (updated_at);
CREATE INDEX stocks_warehouse_id_idx ON shop_dev.stocks USING btree (warehouse_id);


COMMENT ON TABLE shop_dev.stocks IS 'Складские остатки';
COMMENT ON COLUMN shop_dev.stocks.stock_id IS 'Идентификатор';
COMMENT ON COLUMN shop_dev.stocks.quantity IS 'Количество';
COMMENT ON COLUMN shop_dev.stocks.warehouse_id IS 'Идентификатор склада';
COMMENT ON COLUMN shop_dev.stocks.card_id IS 'Идентификатор номенклатуры';
COMMENT ON COLUMN shop_dev.stocks.barcode_id IS 'Штрихкод';
COMMENT ON COLUMN shop_dev.stocks.created_at IS 'Дата создания';
COMMENT ON COLUMN shop_dev.stocks.updated_at IS 'Дата обновления';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX stocks_barcode_idx;
DROP INDEX stocks_card_id_idx;
DROP INDEX stocks_created_at_idx;
DROP INDEX stocks_updated_at_idx;
DROP INDEX stocks_warehouse_id_idx;
DROP TABLE shop_dev.stocks;
-- +goose StatementEnd
