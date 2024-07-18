-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.media_files_types_enum (
    id serial NOT NULL,
    "type" text NOT NULL
);

ALTER TABLE shop.media_files_types_enum OWNER TO shop_user_rw;
ALTER TABLE ONLY shop.media_files_types_enum
    ADD CONSTRAINT media_file_types_enum_pkey PRIMARY KEY (id);

COMMENT ON TABLE shop.media_files_types_enum IS 'Медиафайлы';
COMMENT ON COLUMN shop.media_files_types_enum.id IS 'Идентификатор';
COMMENT ON COLUMN shop.media_files_types_enum.type IS 'Тип медиа файла';

INSERT INTO shop.media_files_types_enum (type) VALUES ('IMAGE'), ('VIDEO');

-------------------------------------------------------------------------

CREATE TABLE shop.media_files (
    id serial NOT NULL,
    link text NOT NULL,
    type_id integer NOT NULL,
    card_id integer NOT NULL
);
ALTER TABLE shop.media_files OWNER TO shop_user_rw;

ALTER TABLE ONLY shop.media_files
    ADD CONSTRAINT media_files_pkey PRIMARY KEY (id);
ALTER TABLE ONLY shop.media_files
    ADD CONSTRAINT media_files_card_id_fkey FOREIGN KEY (card_id) REFERENCES shop.cards(id);
ALTER TABLE ONLY shop.media_files
    ADD CONSTRAINT media_files_type_id_fkey FOREIGN KEY (type_id) REFERENCES shop.media_files_types_enum(id);

CREATE INDEX media_files_card_id_idx ON shop.media_files USING btree (card_id);

COMMENT ON TABLE shop.media_files IS 'Медиафайлы';
COMMENT ON COLUMN shop.media_files.id IS 'Идентификатор';
COMMENT ON COLUMN shop.media_files.link IS 'Ссылка на медиафайл';
COMMENT ON COLUMN shop.media_files.card_id IS 'Идентификатор номенклатуры';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop.media_files_types_enum

DROP INDEX media_files_card_id_idx;
DROP TABLE shop.media_files;
-- +goose StatementEnd


