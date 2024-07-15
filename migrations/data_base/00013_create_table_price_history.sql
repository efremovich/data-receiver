-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.price_history (
    price_history_id serial NOT NULL,
    updated_at timestamp without time zone DEFAULT now(),
    price_size_id integer NOT NULL,
    price numeric(10,2) NOT NULL,
    discount numeric(10,2) NOT NULL
);
ALTER TABLE shop_dev.price_history OWNER TO shop_user_rw;

ALTER TABLE ONLY shop_dev.price_history
    ADD CONSTRAINT price_history_pkey PRIMARY KEY (price_history_id);
ALTER TABLE ONLY shop_dev.price_history
    ADD CONSTRAINT price_history_sizes_fk FOREIGN KEY (price_size_id) REFERENCES shop_dev.price_sizes(price_size_id);

CREATE INDEX price_history_price_size_id_idx ON shop_dev.price_history (price_size_id);

COMMENT ON TABLE shop_dev.price_history IS 'История цен';
COMMENT ON COLUMN shop_dev.price_history.price_history_id IS 'Идентификатор';
COMMENT ON COLUMN shop_dev.price_history.updated_at IS 'Дата обновления';
COMMENT ON COLUMN shop_dev.price_history.price_size_id IS 'Размер';
COMMENT ON COLUMN shop_dev.price_history.price IS 'Цена';
COMMENT ON COLUMN shop_dev.price_history.discount IS 'Скидка';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX price_history_price_size_id_idx;
DROP TABLE shop_dev.price_history;
-- +goose StatementEnd
