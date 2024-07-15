-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.districts (
    district_id serial NOT NULL,
    name text NOT NULL
);
ALTER TABLE shop_dev.districts OWNER TO shop_user_rw;

ALTER TABLE ONLY shop_dev.districts
    ADD CONSTRAINT districts_pk PRIMARY KEY (district_id);
ALTER TABLE ONLY shop_dev.districts
    ADD CONSTRAINT districts_unique UNIQUE (name);

COMMENT ON TABLE shop_dev.districts IS 'Справочник округов';
COMMENT ON COLUMN shop_dev.districts.district_id IS 'Уникальный идентификатор';
COMMENT ON COLUMN shop_dev.districts.name IS 'Наименование округа';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop_dev.districts;
-- +goose StatementEnd
