-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.regions (
    region_id serial NOT NULL,
    region_name text NOT NULL,
    district_id integer NOT NULL,
    country_id integer
);
ALTER TABLE shop_dev.regions OWNER TO shop_user_rw;

ALTER TABLE ONLY shop_dev.regions
    ADD CONSTRAINT regions_pk PRIMARY KEY (region_id);
ALTER TABLE ONLY shop_dev.regions
    ADD CONSTRAINT regions_unique UNIQUE (region_name);
ALTER TABLE ONLY shop_dev.regions
    ADD CONSTRAINT regions_countries_fk FOREIGN KEY (country_id) REFERENCES shop_dev.countries(country_id);
ALTER TABLE ONLY shop_dev.regions
    ADD CONSTRAINT regions_districts_fk FOREIGN KEY (district_id) REFERENCES shop_dev.districts(district_id);

COMMENT ON TABLE shop_dev.regions IS 'Регионы';
COMMENT ON COLUMN shop_dev.regions.region_id IS 'Уникальный идентификатор';
COMMENT ON COLUMN shop_dev.regions.region_name IS 'Наименование региона';
COMMENT ON COLUMN shop_dev.regions.district_id IS 'Округ';
COMMENT ON COLUMN shop_dev.regions.country_id IS 'Страна';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop_dev.regions;
-- +goose StatementEnd
