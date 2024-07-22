-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.cards_characteristics (
    id serial NOT NULL,
    card_id integer NOT NULL,
    characteristic_id integer NOT NULL,
    value text NOT NULL
);
ALTER TABLE shop.cards_characteristics OWNER TO shop_user_rw;

ALTER TABLE ONLY shop.cards_characteristics
    ADD CONSTRAINT cards_characteristics_pkey PRIMARY KEY (id);
ALTER TABLE ONLY shop.cards_characteristics
    ADD CONSTRAINT cards_characteristics_fk FOREIGN KEY (card_id) REFERENCES shop.cards(id);
ALTER TABLE ONLY shop.cards_characteristics
    ADD CONSTRAINT characteristics_id_fk FOREIGN KEY (characteristic_id) REFERENCES shop.characteristics(id);

CREATE INDEX cards_characteristics_card_id_idx ON shop.cards_characteristics (card_id);

COMMENT ON TABLE shop.cards_characteristics IS 'Таблица связи между карточкой товара и характеристиками';
COMMENT ON COLUMN shop.cards_characteristics.id IS 'Уникальный идентификатор';
COMMENT ON COLUMN shop.cards_characteristics.card_id IS 'Карточка товара';
COMMENT ON COLUMN shop.cards_characteristics.characteristic_id IS 'Характеристика товара';
COMMENT ON COLUMN shop.cards_characteristics.value IS 'Значение характеристики';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX cards_characteristics_card_id_idx;
DROP TABLE shop.cards_characteristics;
-- +goose StatementEnd
