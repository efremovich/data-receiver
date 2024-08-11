-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.stocks (
    id serial NOT NULL,
    quantity numeric(10,0) NOT NULL,
    warehouse_id integer NOT NULL,
    card_id integer NOT NULL,
    barcode_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);
ALTER TABLE shop.stocks OWNER TO erp_db_usr;

ALTER TABLE ONLY shop.stocks
    ADD CONSTRAINT stocks_pkey PRIMARY KEY (id);
ALTER TABLE ONLY shop.stocks
    ADD CONSTRAINT stocks_card_id_fkey FOREIGN KEY (card_id) REFERENCES shop.cards(id);
ALTER TABLE ONLY shop.stocks
    ADD CONSTRAINT stocks_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES shop.warehouses(id);
ALTER TABLE ONLY shop.stocks
    ADD CONSTRAINT stocks_fk FOREIGN KEY (barcode_id) REFERENCES shop.barcodes(id);

CREATE INDEX stocks_barcode_idx ON shop.stocks USING btree (barcode_id);
CREATE INDEX stocks_card_id_idx ON shop.stocks USING btree (card_id);
CREATE INDEX stocks_created_at_idx ON shop.stocks USING btree (created_at);
CREATE INDEX stocks_updated_at_idx ON shop.stocks USING btree (updated_at);
CREATE INDEX stocks_warehouse_id_idx ON shop.stocks USING btree (warehouse_id);


COMMENT ON TABLE shop.stocks IS 'Складские остатки';
COMMENT ON COLUMN shop.stocks.id IS 'Идентификатор';
COMMENT ON COLUMN shop.stocks.quantity IS 'Количество';
COMMENT ON COLUMN shop.stocks.warehouse_id IS 'Идентификатор склада';
COMMENT ON COLUMN shop.stocks.card_id IS 'Идентификатор номенклатуры';
COMMENT ON COLUMN shop.stocks.barcode_id IS 'Штрихкод';
COMMENT ON COLUMN shop.stocks.created_at IS 'Дата создания';
COMMENT ON COLUMN shop.stocks.updated_at IS 'Дата обновления';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX stocks_barcode_idx;
DROP INDEX stocks_card_id_idx;
DROP INDEX stocks_created_at_idx;
DROP INDEX stocks_updated_at_idx;
DROP INDEX stocks_warehouse_id_idx;
DROP TABLE shop.stocks;
-- +goose StatementEnd
