-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.media_files (
    id serial NOT NULL,
    link text NOT NULL,
    card_id integer NOT NULL
);
ALTER TABLE shop.media_files OWNER TO shop_user_rw;

ALTER TABLE ONLY shop.media_files
    ADD CONSTRAINT media_files_pkey PRIMARY KEY (id);
ALTER TABLE ONLY shop.media_files
    ADD CONSTRAINT media_files_card_id_fkey FOREIGN KEY (card_id) REFERENCES shop.cards(id);

CREATE INDEX media_files_card_id_idx ON shop.media_files USING btree (card_id);

COMMENT ON TABLE shop.media_files IS 'Медиафайлы';
COMMENT ON COLUMN shop.media_files.id IS 'Идентификатор';
COMMENT ON COLUMN shop.media_files.link IS 'Ссылка на медиафайл';
COMMENT ON COLUMN shop.media_files.card_id IS 'Идентификатор номенклатуры';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX media_files_card_id_idx;
DROP TABLE shop.media_files;
-- +goose StatementEnd
