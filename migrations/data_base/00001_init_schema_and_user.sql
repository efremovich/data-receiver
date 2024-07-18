-- +goose no transaction
-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS shop;
CREATE USER shop_user_rw PASSWORD '${USER_RW_PASSWORD}';
ALTER SCHEMA shop OWNER TO shop_user_rw;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP USER IF EXISTS shop_user_rw;
DROP SCHEMA shop;
-- +goose StatementEnd
