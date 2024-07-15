-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.cards_characteristics (
    card_characteristic_id serial NOT NULL,
    card_id integer NOT NULL,
    characteristic_id integer NOT NULL,
    value text NOT NULL
);
ALTER TABLE shop_dev.cards_characteristics OWNER TO shop_user_rw;

ALTER TABLE ONLY shop_dev.cards_characteristics
    ADD CONSTRAINT cards_characteristics_pkey PRIMARY KEY (card_characteristic_id);
ALTER TABLE ONLY shop_dev.cards_characteristics
    ADD CONSTRAINT cards_characteristics_fk FOREIGN KEY (card_id) REFERENCES shop_dev.cards(card_id);

CREATE INDEX cards_characteristics_card_id_idx ON shop_dev.cards_characteristics (card_id);

COMMENT ON TABLE shop_dev.cards_characteristics IS 'Таблица связи между карточкой товара и характеристиками';
COMMENT ON COLUMN shop_dev.cards_characteristics.card_characteristic_id IS 'Уникальный идентификатор';
COMMENT ON COLUMN shop_dev.cards_characteristics.card_id IS 'Карточка товара';
COMMENT ON COLUMN shop_dev.cards_characteristics.characteristic_id IS 'Характеристика товара';
COMMENT ON COLUMN shop_dev.cards_characteristics.value IS 'Значение характеристики';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX cards_characteristics_card_id_idx;
DROP TABLE shop_dev.cards_characteristics;
-- +goose StatementEnd
