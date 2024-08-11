-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.event_enum (
    id serial NOT NULL,
    event_desc text NOT NULL
);
ALTER TABLE shop.event_enum OWNER TO erp_db_usr;
ALTER TABLE ONLY shop.event_enum
    ADD CONSTRAINT event_enum_pkey PRIMARY KEY (id);

INSERT INTO shop.event_enum (event_desc) VALUES ('CREATED'), ('SUCCESS'), ('GOT_AGAIN'), ('REPROCESS'), ('ERROR'), ('SEND_TASK_NEXT');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop.event_enum;
-- +goose StatementEnd
