-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.warehouse_types (
    id serial NOT NULL,
    name text NOT NULL
);
ALTER TABLE shop.warehouse_types OWNER TO shop_user_rw;

ALTER TABLE ONLY shop.warehouse_types
    ADD CONSTRAINT warehouse_types_name_key UNIQUE (name);
ALTER TABLE ONLY shop.warehouse_types
    ADD CONSTRAINT warehouse_types_pkey PRIMARY KEY (id);

COMMENT ON TABLE shop.warehouse_types IS 'Типы складов';
COMMENT ON COLUMN shop.warehouse_types.id IS 'Уникальный идентификатор';
COMMENT ON COLUMN shop.warehouse_types.name IS 'Наименование склада';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop.warehouse_types;
-- +goose StatementEnd
