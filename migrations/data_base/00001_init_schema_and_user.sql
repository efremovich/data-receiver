-- +goose no transaction
-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS shop_dev;
CREATE USER shop_user_rw PASSWORD '${DB_PASSWORD}';
ALTER SCHEMA shop_dev OWNER TO shop_user_rw;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP USER IF EXISTS shop_user_rw;
DROP SCHEMA shop_dev;
-- +goose StatementEnd
