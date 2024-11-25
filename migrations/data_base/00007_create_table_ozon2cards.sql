-- Отказ от использования этой таблицы

-- +goose Up
-- +goose StatementBegin
-- CREATE TABLE shop.ozon2cards (
--     id serial NOT NULL,
--     created_at timestamp without time zone DEFAULT now(),
--     updated_at timestamp without time zone DEFAULT now(),
--     card_id integer NOT NULL
-- );
-- ALTER TABLE shop.ozon2cards OWNER TO erp_db_usr;

-- ALTER TABLE ONLY shop.ozon2cards
--     ADD CONSTRAINT ozon2cards_pkey PRIMARY KEY (id);
-- ALTER TABLE ONLY shop.ozon2cards
--     ADD CONSTRAINT ozon2cards_card_id_fkey FOREIGN KEY (card_id) REFERENCES shop.cards(id);

-- CREATE INDEX ozon2cards_card_id_idx ON shop.ozon2cards USING btree (card_id);

-- COMMENT ON TABLE shop.ozon2cards IS 'Товары Ozon';
-- COMMENT ON COLUMN shop.ozon2cards.id IS 'Уникальный идентификатор';
-- COMMENT ON COLUMN shop.ozon2cards.id IS 'Идентификатор товара на ozon';
-- COMMENT ON COLUMN shop.ozon2cards.created_at IS 'Дата создания';
-- COMMENT ON COLUMN shop.ozon2cards.updated_at IS 'Дата обновления';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DROP INDEX ozon2cards_card_id_idx;
-- DROP TABLE shop.ozon2cards;
-- +goose StatementEnd
