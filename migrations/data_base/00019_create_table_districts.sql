-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.districts (
    id serial NOT NULL,
    name text NOT NULL
);
ALTER TABLE shop.districts OWNER TO erp_db_usr;

ALTER TABLE ONLY shop.districts
    ADD CONSTRAINT districts_pk PRIMARY KEY (id);
ALTER TABLE ONLY shop.districts
    ADD CONSTRAINT districts_unique UNIQUE (name);

COMMENT ON TABLE shop.districts IS 'Справочник округов';
COMMENT ON COLUMN shop.districts.id IS 'Уникальный идентификатор';
COMMENT ON COLUMN shop.districts.name IS 'Наименование округа';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop.districts;
-- +goose StatementEnd
