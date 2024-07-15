-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.warehouses (
    warehouse_id serial NOT NULL,
    ext_id integer NOT NULL,
    name text NOT NULL,
    address text,
    warehouse_type_id integer,
    seller_id integer NOT NULL
);
ALTER TABLE shop_dev.warehouses OWNER TO shop_user_rw;

ALTER TABLE ONLY shop_dev.warehouses
    ADD CONSTRAINT warehouses_pkey PRIMARY KEY (warehouse_id);
ALTER TABLE ONLY shop_dev.warehouses
    ADD CONSTRAINT warehouse_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES shop_dev.sellers(seller_id);
ALTER TABLE ONLY shop_dev.warehouses
    ADD CONSTRAINT warehouses_fk FOREIGN KEY (warehouse_type_id) REFERENCES shop_dev.warehouse_types(warehouse_type_id);

CREATE INDEX warehouse_seller_id_idx ON shop_dev.warehouses USING btree (seller_id);

COMMENT ON TABLE shop_dev.warehouses IS 'Склады';
COMMENT ON COLUMN shop_dev.warehouses.warehouse_id IS 'Идентификатор';
COMMENT ON COLUMN shop_dev.warehouses.ext_id IS 'Внешний идентификатор';
COMMENT ON COLUMN shop_dev.warehouses.name IS 'Наименование склада';
COMMENT ON COLUMN shop_dev.warehouses.address IS 'Адрес склада';
COMMENT ON COLUMN shop_dev.warehouses.warehouse_type_id IS 'Тип склада';
COMMENT ON COLUMN shop_dev.warehouses.seller_id IS 'Идентификатор продавца';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX warehouse_seller_id_idx;
DROP TABLE shop_dev.warehouses;
-- +goose StatementEnd
