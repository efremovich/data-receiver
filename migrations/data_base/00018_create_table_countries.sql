-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.countries (
    id serial NOT NULL,
    name text NOT NULL
);
ALTER TABLE shop.countries OWNER TO erp_db_usr;

ALTER TABLE ONLY shop.countries
    ADD CONSTRAINT countries_pk PRIMARY KEY (id);
ALTER TABLE ONLY shop.countries
    ADD CONSTRAINT countries_unique UNIQUE (name);

COMMENT ON TABLE shop.countries IS 'Справочник стран';
COMMENT ON COLUMN shop.countries.id IS 'Уникальный идентификатор';
COMMENT ON COLUMN shop.countries.name IS 'Наименование страны';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop.countries;
-- +goose StatementEnd
