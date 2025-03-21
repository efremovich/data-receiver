-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.characteristics (
	id serial NOT NULL,
	title text NOT NULL
);

ALTER TABLE shop.characteristics OWNER TO erp_db_usr;

ALTER TABLE ONLY shop.characteristics
	ADD CONSTRAINT characteristics_pkey PRIMARY KEY (id);

ALTER TABLE ONLY shop.characteristics
	ADD CONSTRAINT characteristics_title_key UNIQUE (title);

COMMENT ON TABLE shop.characteristics IS 'Характеристики';

COMMENT ON COLUMN shop.characteristics.id IS 'Идентификатор';

COMMENT ON COLUMN shop.characteristics.title IS 'Наименование';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
	'down SQL query';

-- +goose StatementEnd
