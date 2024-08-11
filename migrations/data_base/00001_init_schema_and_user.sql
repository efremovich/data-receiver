-- +goose no transaction
-- +goose Up
-- +goose StatementBegin
DO $$ 
BEGIN 
   IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'erp_db_usr') THEN 
       CREATE ROLE erp_db_usr PASSWORD '${DB_USER_RW_PASSWORD}'; 
   END IF; 
END $$;

CREATE SCHEMA IF NOT EXISTS shop;
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'shop_user_rw') THEN
        CREATE ROLE shop_user_rw PASSWORD '${DB_USER_RW_PASSWORD}';
    END IF;
END $$;

ALTER SCHEMA shop OWNER TO erp_db_usr;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP USER IF EXISTS erp_db_usr;
DROP SCHEMA shop;
-- +goose StatementEnd
