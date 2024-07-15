-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.warehouse_types (
    warehouse_type_id serial NOT NULL,
    name text NOT NULL
);
ALTER TABLE shop_dev.warehouse_types OWNER TO shop_user_rw;

ALTER TABLE ONLY shop_dev.warehouse_types
    ADD CONSTRAINT warehouse_types_name_key UNIQUE (name);
ALTER TABLE ONLY shop_dev.warehouse_types
    ADD CONSTRAINT warehouse_types_pkey PRIMARY KEY (warehouse_type_id);

COMMENT ON TABLE shop_dev.warehouse_types IS 'Типы складов';
COMMENT ON COLUMN shop_dev.warehouse_types.warehouse_type_id IS 'Уникальный идентификатор';
COMMENT ON COLUMN shop_dev.warehouse_types.name IS 'Наименование склада';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop_dev.warehouse_types;
-- +goose StatementEnd
