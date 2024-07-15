-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.sizes (
    size_id serial NOT NULL,
    tech_size text NOT NULL,
    name text NOT NULL
);
ALTER TABLE shop_dev.sizes OWNER TO shop_user_rw;

ALTER TABLE ONLY shop_dev.sizes
    ADD CONSTRAINT sizes_pk PRIMARY KEY (size_id);

COMMENT ON TABLE shop_dev.sizes IS 'Размеры';
COMMENT ON COLUMN shop_dev.sizes.size_id IS 'Уникальный идентификатор';
COMMENT ON COLUMN shop_dev.sizes.tech_size IS 'Техническое обозначение товара';
COMMENT ON COLUMN shop_dev.sizes.name IS 'Наименование размера';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop_dev.sizes;
-- +goose StatementEnd
