-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.orders (
    order_id serial NOT NULL,
    ext_id integer NOT NULL,
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
ALTER TABLE shop_dev.orders OWNER TO shop_user_rw;

ALTER TABLE ONLY shop_dev.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (order_id);
ALTER TABLE ONLY shop_dev.orders
    ADD CONSTRAINT orders_card_id_fkey FOREIGN KEY (card_id) REFERENCES shop_dev.cards(card_id);
ALTER TABLE ONLY shop_dev.orders
    ADD CONSTRAINT orders_fk FOREIGN KEY (status_id) REFERENCES shop_dev.statuses(status_id);
ALTER TABLE ONLY shop_dev.orders
    ADD CONSTRAINT orders_price_sizes_fk FOREIGN KEY (price_size_id) REFERENCES shop_dev.price_sizes(price_size_id);
ALTER TABLE ONLY shop_dev.orders
    ADD CONSTRAINT orders_regions_fk FOREIGN KEY (region_id) REFERENCES shop_dev.regions(region_id);
ALTER TABLE ONLY shop_dev.orders
    ADD CONSTRAINT orders_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES shop_dev.sellers(seller_id);
ALTER TABLE ONLY shop_dev.orders
    ADD CONSTRAINT orders_warehouse_idfkey FOREIGN KEY (warehouse_id) REFERENCES shop_dev.warehouses(warehouse_id);


CREATE INDEX orders_card_id_idx ON shop_dev.orders USING btree (card_id);
CREATE INDEX orders_created_at_idx ON shop_dev.orders USING btree (created_at);
CREATE INDEX orders_ext_id_idx ON shop_dev.orders USING btree (ext_id);
CREATE INDEX orders_seller_id_idx ON shop_dev.orders USING btree (seller_id);
CREATE INDEX orders_updated_at_idx ON shop_dev.orders USING btree (updated_at);
CREATE INDEX orders_warehouse_id_idx ON shop_dev.orders USING btree (warehouse_id);

COMMENT ON TABLE shop_dev.orders IS 'Заказы';
COMMENT ON COLUMN shop_dev.orders.order_id IS 'Идентификатор';
COMMENT ON COLUMN shop_dev.orders.ext_id IS 'Внешний идентификатор';
COMMENT ON COLUMN shop_dev.orders.price IS 'Цена';
COMMENT ON COLUMN shop_dev.orders.warehouse_id IS 'Идентификатор склада';
COMMENT ON COLUMN shop_dev.orders.status_id IS 'Статус заказа';
COMMENT ON COLUMN shop_dev.orders.direction IS 'Направление заказа';
COMMENT ON COLUMN shop_dev.orders.type IS 'Тип заказа';
COMMENT ON COLUMN shop_dev.orders.sale IS 'Скидка (СПП)';
COMMENT ON COLUMN shop_dev.orders.card_id IS 'Идентификатор номенклатуры';
COMMENT ON COLUMN shop_dev.orders.seller_id IS 'Идентификатор продавца';
COMMENT ON COLUMN shop_dev.orders.region_id IS 'Регион заказа';
COMMENT ON COLUMN shop_dev.orders.price_size_id IS 'Цена по размеру';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX orders_card_id_idx;
DROP INDEX orders_created_at_idx;
DROP INDEX orders_ext_id_idx;
DROP INDEX orders_seller_id_idx;
DROP INDEX orders_updated_at_idx;
DROP INDEX orders_warehouse_id_idx;
DROP TABLE shop_dev.orders;
-- +goose StatementEnd
