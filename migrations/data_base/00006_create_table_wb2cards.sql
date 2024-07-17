-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.wb2cards (
    id serial NOT NULL,
    nmID integer NOT NULL,
    "int" integer NOT NULL,
    nmUUID text NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    card_id integer NOT NULL
);

ALTER TABLE shop.wb2cards OWNER TO shop_user_rw;

ALTER TABLE ONLY shop.wb2cards
    ADD CONSTRAINT wb2cards_pkey PRIMARY KEY (id);
ALTER TABLE ONLY shop.wb2cards
    ADD CONSTRAINT wb2cards_card_id_fkey FOREIGN KEY (card_id) REFERENCES shop.cards(id);

CREATE INDEX wb2cards_card_id_idx ON shop.wb2cards USING btree (card_id);

COMMENT ON TABLE shop.wb2cards IS 'Товары WB';
COMMENT ON COLUMN shop.wb2cards.id IS 'Уникальный идентификатор';
COMMENT ON COLUMN shop.wb2cards.nmID IS 'Артикул WB';
COMMENT ON COLUMN shop.wb2cards."int" IS 'Идентификатор КТ';
COMMENT ON COLUMN shop.wb2cards.nmUUID IS 'Внуттренний технический идентификатор товара';
COMMENT ON COLUMN shop.wb2cards.created_at IS 'Дата создания';
COMMENT ON COLUMN shop.wb2cards.updated_at IS 'Дата обновления';
COMMENT ON COLUMN shop.wb2cards.card_id IS 'Идентификатор номенклатуры';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX wb2cards_card_id_idx;
DROP TABLE shop.wb2cards;
-- +goose StatementEnd
