-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS "data_seller_enums" (
	"id" integer NOT NULL,
	"isEnable" boolean NOT NULL,
	"name" text NOT NULL,
	PRIMARY KEY ("id")
);

INSERT INTO public.data_seller_enums (is_enable, name) VALUES (true, "wb"), (true, "ozon"), (true, "1c");

CREATE TABLE IF NOT EXISTS "goods_cards" (
	"id" integer NOT NULL,
	"group_id" integer NOT NULL,
	"subject_id" bigint NOT NULL,
	"vendor_code" text NOT NULL,
	"subject_name" text NOT NULL,
	"brand" text NOT NULL,
	"title" text NOT NULL,
	"description" text NOT NULL,
	"seller_id" integer NOT NULL,
	"created_at"  NOT NULL,
	"updated_at"  NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "characteristics" (
	"id" integer NOT NULL,
	"name" text NOT NULL,
	"value" text NOT NULL,
	"card_id" bigint NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "media_files" (
	"id" integer NOT NULL,
	"link" text NOT NULL,
	"card_id" bigint NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "prices" (
	"id" integer NOT NULL,
	"card_id" bigint NOT NULL,
	"pris" bigint NOT NULL,
	"price" double precision NOT NULL,
	"pric" bigint NOT NULL,
	"discount" double precision NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "sizes" (
	"id" integer NOT NULL,
	"card_id" bigint NOT NULL,
	"techSize" text NOT NULL,
	"title" text NOT NULL,
	"barcode" text NOT NULL,
	PRIMARY KEY ("id")
);

ALTER TABLE "data_seller_enums" ADD CONSTRAINT "data_seller_enums_fk0" FOREIGN KEY ("id") REFERENCES "goods_cards"("seller_id");

ALTER TABLE "characteristics" ADD CONSTRAINT "characteristics_fk3" FOREIGN KEY ("card_id") REFERENCES "goods_cards"("id");
ALTER TABLE "media_files" ADD CONSTRAINT "media_files_fk2" FOREIGN KEY ("card_id") REFERENCES "goods_cards"("id");
ALTER TABLE "prices" ADD CONSTRAINT "prices_fk1" FOREIGN KEY ("card_id") REFERENCES "goods_cards"("id");
ALTER TABLE "sizes" ADD CONSTRAINT "sizes_fk1" FOREIGN KEY ("card_id") REFERENCES "goods_cards"("id");

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
select 1
;
-- +goose StatementEnd

