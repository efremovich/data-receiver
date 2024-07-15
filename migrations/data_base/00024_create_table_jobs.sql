-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop_dev.jobs (
    job_id serial8 NOT NULL,
    pub text NOT NULL,
    status text NOT NULL,
    event_type_id integer NOT NULL,
    description text,
    created_at timestamp without time zone DEFAULT now()
);
ALTER TABLE shop_dev.jobs OWNER TO shop_user_rw;
ALTER TABLE shop_dev.jobs
    ADD CONSTRAINT jobs_pkey PRIMARY KEY (job_id);
ALTER TABLE shop_dev.jobs
    ADD CONSTRAINT jobs_event_type_id_fkey FOREIGN KEY (event_type_id) REFERENCES shop_dev.event_enum(event_enum_id);

CREATE INDEX jobs_created_at_idx ON shop_dev.jobs USING btree (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX jobs_created_at_idx;
DROP TABLE shop_dev.jobs;
-- +goose StatementEnd
