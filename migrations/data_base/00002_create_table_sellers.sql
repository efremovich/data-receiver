-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.sellers (
    seller_id serial NOT NULL,
    title text NOT NULL,
    is_enabled boolean DEFAULT true,
    ext_id text
);
ALTER TABLE shop_dev.sellers OWNER TO shop_user_rw;
ALTER TABLE ONLY shop_dev.sellers
    ADD CONSTRAINT sellers_pkey PRIMARY KEY (seller_id);

COMMENT ON TABLE shop_dev.sellers IS 'Продавцы';
COMMENT ON COLUMN shop_dev.sellers.seller_id IS 'Идентификатор';
COMMENT ON COLUMN shop_dev.sellers.title IS 'Наименование продавца';
COMMENT ON COLUMN shop_dev.sellers.is_enabled IS 'Признак активности';
COMMENT ON COLUMN shop_dev.sellers.ext_id IS 'Внешний идентификатор';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop_dev.sellers;
-- +goose StatementEnd
