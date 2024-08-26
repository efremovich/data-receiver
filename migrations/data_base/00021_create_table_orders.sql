-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.orders (
    id serial NOT NULL,
    external_id text NOT NULL,
    price numeric(10,2) NOT NULL,
    warehouse_id integer NOT NULL,
    status_id integer,
    direction text,
    type text,
    sale numeric(10,2),
    card_id integer NOT NULL,
    seller_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    region_id integer,
    price_size_id integer NOT NULL
);
ALTER TABLE shop.orders OWNER TO erp_db_usr;

ALTER TABLE ONLY shop.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);
ALTER TABLE ONLY shop.orders
    ADD CONSTRAINT orders_card_id_fkey FOREIGN KEY (card_id) REFERENCES shop.cards(id);
ALTER TABLE ONLY shop.orders
    ADD CONSTRAINT orders_fk FOREIGN KEY (status_id) REFERENCES shop.statuses(id);
ALTER TABLE ONLY shop.orders
    ADD CONSTRAINT orders_price_sizes_fk FOREIGN KEY (price_size_id) REFERENCES shop.price_sizes(id);
ALTER TABLE ONLY shop.orders
    ADD CONSTRAINT orders_regions_fk FOREIGN KEY (region_id) REFERENCES shop.regions(id);
ALTER TABLE ONLY shop.orders
    ADD CONSTRAINT orders_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES shop.sellers(id);
ALTER TABLE ONLY shop.orders
    ADD CONSTRAINT orders_warehouse_idfkey FOREIGN KEY (warehouse_id) REFERENCES shop.warehouses(id);


CREATE INDEX orders_card_id_idx ON shop.orders USING btree (card_id);
CREATE INDEX orders_created_at_idx ON shop.orders USING btree (created_at);
CREATE INDEX orders_external_id_idx ON shop.orders USING btree (external_id);
CREATE INDEX orders_seller_id_idx ON shop.orders USING btree (seller_id);
CREATE INDEX orders_updated_at_idx ON shop.orders USING btree (updated_at);
CREATE INDEX orders_warehouse_id_idx ON shop.orders USING btree (warehouse_id);

COMMENT ON TABLE shop.orders IS 'Заказы';
COMMENT ON COLUMN shop.orders.id IS 'Идентификатор';
COMMENT ON COLUMN shop.orders.external_id IS 'Внешний идентификатор';
COMMENT ON COLUMN shop.orders.price IS 'Цена';
COMMENT ON COLUMN shop.orders.warehouse_id IS 'Идентификатор склада';
COMMENT ON COLUMN shop.orders.status_id IS 'Статус заказа';
COMMENT ON COLUMN shop.orders.direction IS 'Направление заказа';
COMMENT ON COLUMN shop.orders.type IS 'Тип заказа';
COMMENT ON COLUMN shop.orders.sale IS 'Скидка (СПП)';
COMMENT ON COLUMN shop.orders.card_id IS 'Идентификатор номенклатуры';
COMMENT ON COLUMN shop.orders.seller_id IS 'Идентификатор продавца';
COMMENT ON COLUMN shop.orders.region_id IS 'Регион заказа';
COMMENT ON COLUMN shop.orders.price_size_id IS 'Цена по размеру';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX orders_card_id_idx;
DROP INDEX orders_created_at_idx;
DROP INDEX orders_external_id_idx;
DROP INDEX orders_seller_id_idx;
DROP INDEX orders_updated_at_idx;
DROP INDEX orders_warehouse_id_idx;
DROP TABLE shop.orders;
-- +goose StatementEnd
