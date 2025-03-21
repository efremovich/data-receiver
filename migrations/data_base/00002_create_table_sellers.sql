-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.sellers (
	id serial NOT NULL,
	title text NOT NULL,
	is_enabled boolean DEFAULT TRUE,
	external_id text
);

ALTER TABLE shop.sellers OWNER TO erp_db_usr;

ALTER TABLE ONLY shop.sellers
	ADD CONSTRAINT sellers_pkey PRIMARY KEY (id);

COMMENT ON TABLE shop.sellers IS 'Продавцы';

COMMENT ON COLUMN shop.sellers.id IS 'Идентификатор';

COMMENT ON COLUMN shop.sellers.title IS 'Наименование продавца';

COMMENT ON COLUMN shop.sellers.is_enabled IS 'Признак активности';

COMMENT ON COLUMN shop.sellers.external_id IS 'Внешний идентификатор продавца';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE shop.sellers;

-- +goose StatementEnd
