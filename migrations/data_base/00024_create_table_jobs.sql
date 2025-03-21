-- +goose Up
-- +goose StatementBegin
CREATE TABLE shop.jobs (
	id serial NOT NULL,
	pub text NOT NULL,
	status text NOT NULL,
	event_type_id integer NOT NULL,
	description text,
	created_at timestamp without time zone DEFAULT NOW()
);

ALTER TABLE shop.jobs OWNER TO erp_db_usr;

ALTER TABLE shop.jobs
	ADD CONSTRAINT jobs_pkey PRIMARY KEY (id);

ALTER TABLE shop.jobs
	ADD CONSTRAINT jobs_event_type_id_fkey FOREIGN KEY (event_type_id) REFERENCES shop.event_enum (id);

CREATE INDEX jobs_created_at_idx ON shop.jobs USING btree (created_at);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX jobs_created_at_idx;

DROP TABLE shop.jobs;

-- +goose StatementEnd
