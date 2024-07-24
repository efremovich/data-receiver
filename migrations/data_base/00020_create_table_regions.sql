-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.regions (
    id serial NOT NULL,
    region_name text NOT NULL,
    district_id integer NOT NULL,
    country_id integer
);
ALTER TABLE shop.regions OWNER TO shop_user_rw;

ALTER TABLE ONLY shop.regions
    ADD CONSTRAINT regions_pk PRIMARY KEY (id);
ALTER TABLE ONLY shop.regions
    ADD CONSTRAINT regions_unique UNIQUE (region_name);
ALTER TABLE ONLY shop.regions
    ADD CONSTRAINT regions_countries_fk FOREIGN KEY (country_id) REFERENCES shop.countries(id);
ALTER TABLE ONLY shop.regions
    ADD CONSTRAINT regions_districts_fk FOREIGN KEY (district_id) REFERENCES shop.districts(id);

COMMENT ON TABLE shop.regions IS 'Регионы';
COMMENT ON COLUMN shop.regions.id IS 'Уникальный идентификатор';
COMMENT ON COLUMN shop.regions.region_name IS 'Наименование региона';
COMMENT ON COLUMN shop.regions.district_id IS 'Округ';
COMMENT ON COLUMN shop.regions.country_id IS 'Страна';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop.regions;
-- +goose StatementEnd
