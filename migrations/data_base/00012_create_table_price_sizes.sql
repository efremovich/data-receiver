-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.price_sizes (
	id serial NOT NULL,
	card_id integer NOT NULL,
	size_id integer NOT NULL,
	price numeric(10, 2) NOT NULL,
	price_with_discount numeric(10, 2) NOT NULL,
	price_finish numeric(10, 2) NOT NULL,
	updated_at timestamp NOT NULL
);

ALTER TABLE shop.price_sizes OWNER TO erp_db_usr;

ALTER TABLE ONLY shop.price_sizes
	ADD CONSTRAINT price_sizes_pkey PRIMARY KEY (id);

ALTER TABLE ONLY shop.price_sizes
	ADD CONSTRAINT price_sizes_cards_fk FOREIGN KEY (card_id) REFERENCES shop.cards (id);

ALTER TABLE ONLY shop.price_sizes
	ADD CONSTRAINT price_sizes_sizes_fk FOREIGN KEY (size_id) REFERENCES shop.sizes (id);

CREATE INDEX price_sizes_card_id_idx ON shop.price_sizes (card_id);

CREATE INDEX price_sizes_size_id_idx ON shop.price_sizes (size_id);

COMMENT ON TABLE shop.price_sizes IS 'Размеры';

COMMENT ON COLUMN shop.price_sizes.id IS 'Уникальный идентификатор';

COMMENT ON COLUMN shop.price_sizes.card_id IS 'Идентификатор номенклатуры';

COMMENT ON COLUMN shop.price_sizes.size_id IS 'Размер';

COMMENT ON COLUMN shop.price_sizes.price IS 'Цена';

COMMENT ON COLUMN shop.price_sizes.price_finish IS 'Окончательная цена';

COMMENT ON COLUMN shop.price_sizes.price_with_discount IS 'Цена без скидки';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX price_sizes_card_id_idx;

DROP INDEX price_sizes_size_id_idx;

DROP TABLE shop.price_sizes;

-- +goose StatementEnd
