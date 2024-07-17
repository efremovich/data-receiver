-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.price_history (
    id serial NOT NULL,
    updated_at timestamp without time zone DEFAULT now(),
    price_size_id integer NOT NULL,
    price numeric(10,2) NOT NULL,
    discount numeric(10,2) NOT NULL
);
ALTER TABLE shop.price_history OWNER TO shop_user_rw;

ALTER TABLE ONLY shop.price_history
    ADD CONSTRAINT price_history_pkey PRIMARY KEY (id);
ALTER TABLE ONLY shop.price_history
    ADD CONSTRAINT price_history_sizes_fk FOREIGN KEY (price_size_id) REFERENCES shop.price_sizes(id);

CREATE INDEX price_history_price_size_id_idx ON shop.price_history (price_size_id);

COMMENT ON TABLE shop.price_history IS 'История цен';
COMMENT ON COLUMN shop.price_history.id IS 'Идентификатор';
COMMENT ON COLUMN shop.price_history.updated_at IS 'Дата обновления';
COMMENT ON COLUMN shop.price_history.price_size_id IS 'Размер';
COMMENT ON COLUMN shop.price_history.price IS 'Цена';
COMMENT ON COLUMN shop.price_history.discount IS 'Скидка';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX price_history_price_size_id_idx;
DROP TABLE shop.price_history;
-- +goose StatementEnd
