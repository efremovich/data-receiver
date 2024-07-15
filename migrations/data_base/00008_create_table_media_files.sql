-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.media_files (
    media_file_id serial NOT NULL,
    link text NOT NULL,
    card_id integer NOT NULL
);
ALTER TABLE shop_dev.media_files OWNER TO shop_user_rw;

ALTER TABLE ONLY shop_dev.media_files
    ADD CONSTRAINT media_files_pkey PRIMARY KEY (media_file_id);
ALTER TABLE ONLY shop_dev.media_files
    ADD CONSTRAINT media_files_card_id_fkey FOREIGN KEY (card_id) REFERENCES shop_dev.cards(card_id);

CREATE INDEX media_files_card_id_idx ON shop_dev.media_files USING btree (card_id);

COMMENT ON TABLE shop_dev.media_files IS 'Медиафайлы';
COMMENT ON COLUMN shop_dev.media_files.media_file_id IS 'Идентификатор';
COMMENT ON COLUMN shop_dev.media_files.link IS 'Ссылка на медиафайл';
COMMENT ON COLUMN shop_dev.media_files.card_id IS 'Идентификатор номенклатуры';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX media_files_card_id_idx;
DROP TABLE shop_dev.media_files;
-- +goose StatementEnd
