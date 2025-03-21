-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.sales (
	id serial NOT NULL,
	external_id text NOT NULL,
	price numeric(10, 2) NOT NULL,
	discount numeric(10, 2) NOT NULL,
	final_price numeric(10, 2) NOT NULL,
	type text,
	for_pay numeric(10, 2),
	quantity integer NOT NULL,
	created_at timestamp without time zone DEFAULT NOW(),
	updated_at timestamp without time zone DEFAULT NOW(),
	order_id integer NOT NULL,
	seller_id integer NOT NULL,
	card_id integer NOT NULL,
	warehouse_id integer NOT NULL,
	region_id integer,
	price_size_id integer NOT NULL
);

ALTER TABLE shop.sales OWNER TO erp_db_usr;

ALTER TABLE ONLY shop.sales
	ADD CONSTRAINT sales_pkey PRIMARY KEY (id);

ALTER TABLE ONLY shop.sales
	ADD CONSTRAINT sales_card_id_fkey FOREIGN KEY (card_id) REFERENCES shop.cards (id);

ALTER TABLE ONLY shop.sales
	ADD CONSTRAINT sales_price_sizes_fk FOREIGN KEY (price_size_id) REFERENCES shop.price_sizes (id);

ALTER TABLE ONLY shop.sales
	ADD CONSTRAINT sales_regions_fk FOREIGN KEY (region_id) REFERENCES shop.regions (id);

ALTER TABLE ONLY shop.sales
	ADD CONSTRAINT sales_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES shop.sellers (id);

ALTER TABLE ONLY shop.sales
	ADD CONSTRAINT sales_warehouse_idfkey FOREIGN KEY (warehouse_id) REFERENCES shop.warehouses (id);

ALTER TABLE ONLY shop.sales
	ADD CONSTRAINT sales_order_idfkey FOREIGN KEY (order_id) REFERENCES shop.orders (id);

CREATE INDEX sales_card_id_idx ON shop.sales USING btree (card_id);

CREATE INDEX sales_created_at_idx ON shop.sales USING btree (created_at);

CREATE INDEX sales_external_id_idx ON shop.sales USING btree (external_id);

CREATE INDEX sales_seller_id_idx ON shop.sales USING btree (seller_id);

CREATE INDEX sales_updated_at_idx ON shop.sales USING btree (updated_at);

CREATE INDEX sales_warehouse_id_idx ON shop.sales USING btree (warehouse_id);

CREATE INDEX sales_order_id_idx ON shop.sales USING btree (order_id);

COMMENT ON TABLE shop.sales IS 'Продажи';

COMMENT ON COLUMN shop.sales.id IS 'Идентификатор';

COMMENT ON COLUMN shop.sales.external_id IS 'Внешний идентификатор';

COMMENT ON COLUMN shop.sales.price IS 'Цена';

COMMENT ON COLUMN shop.sales.discount IS 'Скидка';

COMMENT ON COLUMN shop.sales.final_price IS 'Конечная цена';

COMMENT ON COLUMN shop.sales.type IS 'Тип';

COMMENT ON COLUMN shop.sales.for_pay IS 'Сумма оплаты';

COMMENT ON COLUMN shop.sales.quantity IS 'Количество';

COMMENT ON COLUMN shop.sales.created_at IS 'Дата создания';

COMMENT ON COLUMN shop.sales.updated_at IS 'Дата обновления';

COMMENT ON COLUMN shop.sales.order_id IS 'Идентификатор заказа';

COMMENT ON COLUMN shop.sales.seller_id IS 'Идентификатор продавца';

COMMENT ON COLUMN shop.sales.card_id IS 'Идентификатор карты';

COMMENT ON COLUMN shop.sales.warehouse_id IS 'Идентификатор склада';

COMMENT ON COLUMN shop.sales.region_id IS 'Идентификатор региона';

COMMENT ON COLUMN shop.sales.price_size_id IS 'Идентификатор размера цены';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX sales_card_id_idx;

DROP INDEX sales_created_at_idx;

DROP INDEX sales_external_id_idx;

DROP INDEX sales_seller_id_idx;

DROP INDEX sales_updated_at_idx;

DROP INDEX sales_warehouse_id_idx;

DROP INDEX sales_order_id_idx;

DROP TABLE shop.sales;

-- +goose StatementEnd
