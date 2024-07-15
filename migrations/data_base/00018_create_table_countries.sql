-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.countries (
    country_id serial NOT NULL,
    name text NOT NULL
);
ALTER TABLE shop_dev.countries OWNER TO shop_user_rw;

ALTER TABLE ONLY shop_dev.countries
    ADD CONSTRAINT countries_pk PRIMARY KEY (country_id);
ALTER TABLE ONLY shop_dev.countries
    ADD CONSTRAINT countries_unique UNIQUE (name);

COMMENT ON TABLE shop_dev.countries IS 'Справочник стран';
COMMENT ON COLUMN shop_dev.countries.country_id IS 'Уникальный идентификатор';
COMMENT ON COLUMN shop_dev.countries.name IS 'Наименование страны';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop_dev.countries;
-- +goose StatementEnd
