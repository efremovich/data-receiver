-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.warehouses (
    id serial NOT NULL,
    external_id integer NOT NULL,
    name text NOT NULL,
    address text,
    warehouse_type_id integer,
    seller_id integer NOT NULL
);
ALTER TABLE shop.warehouses OWNER TO shop_user_rw;

ALTER TABLE ONLY shop.warehouses
    ADD CONSTRAINT warehouses_pkey PRIMARY KEY (id);
ALTER TABLE ONLY shop.warehouses
    ADD CONSTRAINT warehouse_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES shop.sellers(id);
ALTER TABLE ONLY shop.warehouses
    ADD CONSTRAINT warehouses_fk FOREIGN KEY (warehouse_type_id) REFERENCES shop.warehouse_types(id);

CREATE INDEX warehouse_seller_id_idx ON shop.warehouses USING btree (seller_id);

COMMENT ON TABLE shop.warehouses IS 'Склады';
COMMENT ON COLUMN shop.warehouses.id IS 'Идентификатор';
COMMENT ON COLUMN shop.warehouses.external_id IS 'Внешний идентификатор';
COMMENT ON COLUMN shop.warehouses.name IS 'Наименование склада';
COMMENT ON COLUMN shop.warehouses.address IS 'Адрес склада';
COMMENT ON COLUMN shop.warehouses.warehouse_type_id IS 'Тип склада';
COMMENT ON COLUMN shop.warehouses.seller_id IS 'Идентификатор продавца';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX warehouse_seller_id_idx;
DROP TABLE shop.warehouses;
-- +goose StatementEnd
