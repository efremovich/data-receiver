-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.tp_status_enum (
    id SERIAL PRIMARY KEY,
    tp_status_desc VARCHAR NOT NULL
);
INSERT INTO public.tp_status_enum (tp_status_desc) VALUES ('new'), ('success'), ('failed'), ('failed_internal');

CREATE TABLE public.tp_event_enum (
    id SERIAL PRIMARY KEY,
    tp_event_desc VARCHAR NOT NULL
);
INSERT INTO public.tp_event_enum (tp_event_desc) VALUES ('CREATED'), ('SUCCESS'), ('GOT_AGAIN'), ('REPROCESS'), ('ERROR'), ('SEND_TASK_NEXT');

CREATE TABLE public.transport_package (
	id BIGSERIAL PRIMARY KEY,
    name VARCHAR(40) NOT NULL,
	is_receipt BOOLEAN,
    sender_operator_code VARCHAR(40),
    receipt_url text NOT NULL,
    tp_status_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
	CONSTRAINT tp_status_id_fkey FOREIGN KEY (tp_status_id) REFERENCES public.tp_status_enum(id)
);

CREATE INDEX idx_tp_name ON public.transport_package (name);
CREATE INDEX idx_tp_sender_operator_code ON public.transport_package (sender_operator_code);
CREATE INDEX idx_tp_receipt_url ON public.transport_package (receipt_url);
CREATE INDEX idx_tp_created_at ON public.transport_package (created_at);
CREATE INDEX idx_tp_updated_at ON public.transport_package (updated_at);

COMMENT ON COLUMN transport_package.sender_operator_code is 'Код оператора-отправителя пакета';
COMMENT ON COLUMN transport_package.receipt_url is 'URL для получения технологической квитанции';

CREATE TABLE public.tp_error (
    tp_id BIGINT PRIMARY KEY,
	error_text text,
	error_code varchar(150),
	CONSTRAINT tp_error_id_fkey FOREIGN KEY (tp_id) REFERENCES public.transport_package(id)
);
CREATE INDEX idx_tp_error_text ON public.tp_error (error_text);
CREATE INDEX idx_tp_error_code ON public.tp_error (error_code);

CREATE TABLE public.tp_event (
	id BIGSERIAL PRIMARY KEY,
	tp_id BIGINT NOT NULL,
	event_type_id int NOT NULL,
	description text,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
	CONSTRAINT event_tp_id_fkey FOREIGN KEY (tp_id) REFERENCES public.transport_package(id),
	CONSTRAINT event_type_id_fkey FOREIGN KEY (event_type_id) REFERENCES public.tp_event_enum(id)
);
CREATE INDEX idx_tp_id_event ON public.tp_event (tp_id);
CREATE INDEX idx_event_created_at ON public.tp_event (created_at);

CREATE TABLE public.tp_directory (
	id BIGSERIAL PRIMARY KEY,
    tp_id BIGINT,
	name VARCHAR(50),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
	CONSTRAINT tp_content_id_fkey FOREIGN KEY (tp_id) REFERENCES public.transport_package(id)
);
CREATE UNIQUE INDEX idx_directory_unique ON public.tp_directory (tp_id, name);
CREATE INDEX idx_directory_name ON public.tp_directory (name);
CREATE INDEX idx_directory_created_at ON public.tp_directory (created_at);

CREATE TABLE public.tp_document (
	id BIGSERIAL PRIMARY KEY,
    directory_id BIGINT,
	name VARCHAR(50),
	CONSTRAINT tp_content_id_fkey FOREIGN KEY (directory_id) REFERENCES public.tp_directory(id)
);
CREATE UNIQUE INDEX idx_document_unique ON public.tp_document (directory_id, name);
CREATE INDEX idx_document_name ON public.tp_document (name);


-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
-- не будем сносить всю бд
SELECT 1;
-- +goose StatementEnd