-- + goose Up
-- +goose StatementBegin
CREATE TABLE shop.pvzs (
	id serial NOT NULL,
	external_id text, -- ID поставщика
	office_name text, -- Наименование офиса
	office_id integer, -- ID офиса
	supplier_name text, -- Наименование поставщика
	supplier_id integer, -- ID поставщика
	supplier_inn text -- ИНН поставщика
);

ALTER TABLE shop.pvzs OWNER to erp_db_usr;

ALTER TABLE ONLY shop.pvzs
	ADD CONSTRAINT pvzs_pkey PRIMARY KEY (id);

-- Комментарии к столбцам
COMMENT ON TABLE shop.pvzs IS 'Пункты выдачи заказов';

COMMENT ON COLUMN shop.pvzs.id IS 'Идентификатор';

COMMENT ON COLUMN shop.pvzs.external_id IS 'ID поставщика';

COMMENT ON COLUMN shop.pvzs.office_name IS 'Наименование офиса';

COMMENT ON COLUMN shop.pvzs.office_id IS 'ID офиса';

COMMENT ON COLUMN shop.pvzs.supplier_name IS 'Наименование поставщика';

COMMENT ON COLUMN shop.pvzs.supplier_id IS 'ID поставщика';

COMMENT ON COLUMN shop.pvzs.supplier_inn IS 'ИНН поставщика';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE shop.pvzs;

-- +goose StatementEnd
