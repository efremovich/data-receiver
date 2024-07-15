-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.characteristics (
    characteristic_id serial NOT NULL,
    title text NOT NULL
);
ALTER TABLE shop_dev.characteristics OWNER TO shop_user_rw;

ALTER TABLE ONLY shop_dev.characteristics
    ADD CONSTRAINT characteristics_pkey PRIMARY KEY (characteristic_id);
ALTER TABLE ONLY shop_dev.characteristics
    ADD CONSTRAINT characteristics_title_key UNIQUE (title);

COMMENT ON TABLE shop_dev.characteristics IS 'Характеристики';
COMMENT ON COLUMN shop_dev.characteristics.characteristic_id IS 'Идентификатор';
COMMENT ON COLUMN shop_dev.characteristics.title IS 'Наименование';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
