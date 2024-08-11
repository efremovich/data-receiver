-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.statuses (
    id serial NOT NULL,
    name text NOT NULL
);
ALTER TABLE shop.statuses OWNER TO erp_db_usr;

ALTER TABLE ONLY shop.statuses
    ADD CONSTRAINT statuses_name_key UNIQUE (name);
ALTER TABLE ONLY shop.statuses
    ADD CONSTRAINT statuses_pkey PRIMARY KEY (id);

COMMENT ON TABLE shop.statuses IS 'Статусы заказа';
COMMENT ON COLUMN shop.statuses.id IS 'Уникальный идентификатор';
COMMENT ON COLUMN shop.statuses.name IS 'Наименование статуса';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop.statuses;
-- +goose StatementEnd
