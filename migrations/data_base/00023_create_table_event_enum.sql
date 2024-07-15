-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.event_enum (
    event_enum_id integer NOT NULL,
    event_desc text NOT NULL
);
ALTER TABLE shop_dev.event_enum OWNER TO shop_user_rw;
ALTER TABLE ONLY shop_dev.event_enum
    ADD CONSTRAINT event_enum_pkey PRIMARY KEY (event_enum_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop_dev.event_enum;
-- +goose StatementEnd
