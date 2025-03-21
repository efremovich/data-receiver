-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.sizes (
	id serial NOT NULL,
	tech_size text NOT NULL,
	name text NOT NULL
);

ALTER TABLE shop.sizes OWNER TO erp_db_usr;

ALTER TABLE ONLY shop.sizes
	ADD CONSTRAINT sizes_pk PRIMARY KEY (id);

COMMENT ON TABLE shop.sizes IS 'Размеры';

COMMENT ON COLUMN shop.sizes.id IS 'Уникальный идентификатор';

COMMENT ON COLUMN shop.sizes.tech_size IS 'Техническое обозначение товара';

COMMENT ON COLUMN shop.sizes.name IS 'Наименование размера';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE shop.sizes;

-- +goose StatementEnd
