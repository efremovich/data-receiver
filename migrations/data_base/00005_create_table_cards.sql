-- +goose no transaction
-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.cards (
    card_id serial NOT NULL,
    vendor_id text NOT NULL,
    vendor_code text NOT NULL,
    title text NOT NULL,
    description text,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    brand_id integer NOT NULL,
    category_id integer NOT NULL
);
ALTER TABLE shop_dev.cards OWNER TO shop_user_rw;

ALTER TABLE ONLY shop_dev.cards
    ADD CONSTRAINT cards_pkey PRIMARY KEY (card_id);
ALTER TABLE ONLY shop_dev.cards
    ADD CONSTRAINT cards_brand_id_fkey FOREIGN KEY (brand_id) REFERENCES shop_dev.brands(brand_id);
ALTER TABLE ONLY shop_dev.cards
    ADD CONSTRAINT cards_category_id_fkey FOREIGN KEY (category_id) REFERENCES shop_dev.categories(category_id);

CREATE INDEX cards_created_at_idx ON shop_dev.cards USING btree (created_at);
CREATE INDEX cards_title_idx ON shop_dev.cards USING btree (title);
CREATE INDEX cards_updated_at_idx ON shop_dev.cards USING btree (updated_at);
CREATE INDEX cards_vendor_code_idx ON shop_dev.cards USING btree (vendor_code);
CREATE INDEX cards_vendor_id_idx ON shop_dev.cards USING btree (vendor_id);

COMMENT ON TABLE shop_dev.cards IS 'Товары';
COMMENT ON COLUMN shop_dev.cards.card_id IS 'Внутренний идентификатор';
COMMENT ON COLUMN shop_dev.cards.vendor_id IS 'Внутренный код товара (из 1с)';
COMMENT ON COLUMN shop_dev.cards.vendor_code IS 'Артикул (из 1с)';
COMMENT ON COLUMN shop_dev.cards.title IS 'Наименование номенклатуры';
COMMENT ON COLUMN shop_dev.cards.description IS 'Описание номенклатуры';
COMMENT ON COLUMN shop_dev.cards.created_at IS 'Дата создания';
COMMENT ON COLUMN shop_dev.cards.updated_at IS 'Дата обновления';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX cards_created_at_idx;
DROP INDEX cards_title_idx;
DROP INDEX cards_updated_at_idx;
DROP INDEX cards_vendor_code_idx;
DROP INDEX cards_vendor_id_idx;
DROP TABLE shop_dev.cards;
-- +goose StatementEnd
