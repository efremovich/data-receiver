-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.seller2cards (
	id serial NOT NULL,
	external_id integer NOT NULL,
	"int" integer NOT NULL,
	nmUUID text NOT NULL,
	created_at timestamp without time zone DEFAULT NOW(),
	updated_at timestamp without time zone DEFAULT NOW(),
	card_id integer NOT NULL,
	seller_id integer NOT NULL
);

ALTER TABLE shop.seller2cards OWNER TO erp_db_usr;

ALTER TABLE ONLY shop.seller2cards
	ADD CONSTRAINT seller2cards_pkey PRIMARY KEY (id);

ALTER TABLE ONLY shop.seller2cards
	ADD CONSTRAINT seller2cards_card_id_fkey FOREIGN KEY (card_id) REFERENCES shop.cards (id);

ALTER TABLE shop.seller2cards
	ADD CONSTRAINT seller2cards_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES shop.sellers (id);

CREATE INDEX seller2cards_card_id_idx ON shop.seller2cards USING btree (card_id);

CREATE INDEX seller2cards_seller_id_idx ON shop.seller2cards USING btree (seller_id);

CREATE INDEX seller2cards_external_id_idx ON shop.seller2cards USING btree (external_id);

COMMENT ON TABLE shop.seller2cards IS 'Товары WB';

COMMENT ON COLUMN shop.seller2cards.id IS 'Уникальный идентификатор';

COMMENT ON COLUMN shop.seller2cards.external_id IS 'Внешний код карточки';

COMMENT ON COLUMN shop.seller2cards. "int" IS 'Идентификатор КТ';

COMMENT ON COLUMN shop.seller2cards.nmUUID IS 'Внуттренний технический идентификатор товара';

COMMENT ON COLUMN shop.seller2cards.created_at IS 'Дата создания';

COMMENT ON COLUMN shop.seller2cards.updated_at IS 'Дата обновления';

COMMENT ON COLUMN shop.seller2cards.card_id IS 'Идентификатор номенклатуры';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX seller2cards_seller_id_idx;

DROP INDEX seller2cards_card_id_idx;

DROP TABLE shop.seller2cards;

-- +goose StatementEnd
