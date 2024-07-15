-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.price_sizes (
    price_size_id serial NOT NULL,
    card_id integer NOT NULL,
    size_id integer NOT NULL,
    price numeric(10,2) NOT NULL,
    discount numeric(10,2) NOT NULL,
    updated_at timestamp NOT NULL
);
ALTER TABLE shop_dev.price_sizes OWNER TO shop_user_rw;

ALTER TABLE ONLY shop_dev.price_sizes
    ADD CONSTRAINT price_sizes_pkey PRIMARY KEY (price_size_id);
ALTER TABLE ONLY shop_dev.price_sizes
    ADD CONSTRAINT price_sizes_cards_fk FOREIGN KEY (card_id) REFERENCES shop_dev.cards(card_id);
ALTER TABLE ONLY shop_dev.price_sizes
    ADD CONSTRAINT price_sizes_sizes_fk FOREIGN KEY (size_id) REFERENCES shop_dev.sizes(size_id);

CREATE INDEX price_sizes_card_id_idx ON shop_dev.price_sizes (card_id);
CREATE INDEX price_sizes_size_id_idx ON shop_dev.price_sizes (size_id);


COMMENT ON TABLE shop_dev.price_sizes IS 'Размеры';
COMMENT ON COLUMN shop_dev.price_sizes.price_size_id IS 'Уникальный идентификатор';
COMMENT ON COLUMN shop_dev.price_sizes.card_id IS 'Идентификатор номенклатуры';
COMMENT ON COLUMN shop_dev.price_sizes.size_id IS 'Размер';
COMMENT ON COLUMN shop_dev.price_sizes.price IS 'Цена';
COMMENT ON COLUMN shop_dev.price_sizes.discount IS 'Скидка';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX price_sizes_card_id_idx;
DROP INDEX price_sizes_size_id_idx;
DROP TABLE shop_dev.price_sizes;
-- +goose StatementEnd
