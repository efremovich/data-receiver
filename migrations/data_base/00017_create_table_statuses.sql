-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.statuses (
    status_id serial NOT NULL,
    name text NOT NULL
);
ALTER TABLE shop_dev.statuses OWNER TO shop_user_rw;

ALTER TABLE ONLY shop_dev.statuses
    ADD CONSTRAINT statuses_name_key UNIQUE (name);
ALTER TABLE ONLY shop_dev.statuses
    ADD CONSTRAINT statuses_pkey PRIMARY KEY (status_id);

COMMENT ON TABLE shop_dev.statuses IS 'Статусы заказа';
COMMENT ON COLUMN shop_dev.statuses.status_id IS 'Уникальный идентификатор';
COMMENT ON COLUMN shop_dev.statuses.name IS 'Наименование статуса';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop_dev.statuses;
-- +goose StatementEnd
