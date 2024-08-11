-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.dimensions (
    id serial NOT NULL,
    length integer NOT NULL,
    width integer NOT NULL,
    height integer NOT NULL,
    is_valid bool,
    card_id integer NOT NULL
);
ALTER TABLE shop.dimensions OWNER TO erp_db_usr;
ALTER TABLE shop.dimensions
    ADD CONSTRAINT dimensions_pkey PRIMARY KEY (id);
ALTER TABLE shop.dimensions
    ADD CONSTRAINT dimention_card_id_fkey FOREIGN KEY (card_id) REFERENCES shop.cards(id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop.dimensions;
-- +goose StatementEnd
