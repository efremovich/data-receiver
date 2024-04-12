-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.data_seller_enums (
	id SERIAL PRIMARY KEY,
	name VARCHAR NOT NULL
);
INSERT INTO public.data_seller_enums (name) VALUES ('wb'), ('ozon'), ('1c');
COMMENT ON COLUMN data_seller_enums.name is 'Наименование поставщика данных';

CREATE TABLE public.goods_cards (
	id SERIAL PRIMARY KEY,
  vendor_id SERIAL NOT NULL,
	group_id SERIAL NOT NULL,
	subject_id SERIAL NOT NULL,
	subject_name SERIAL NOT NULL,
	vendor_code VARCHAR NOT NULL,
	brand TEXT NOT NULL,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	seller_id SERIAL NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
	CONSTRAINT data_seller_enums_fkey FOREIGN KEY (seller_id) REFERENCES public.data_seller_enums(id)
);

COMMENT ON COLUMN goods_cards.id is 'Внутренний идентификатор';
COMMENT ON COLUMN goods_cards.vendor_id is 'Идентификатор продавца';

CREATE TABLE public.characteristics (
	id SERIAL PRIMARY KEY,
	name VARCHAR NOT NULL,
	value VARCHAR NOT NULL,
	card_id SERIAL NOT NULL,
  CONSTRAINT characteristics_fkey FOREIGN KEY ("card_id") REFERENCES public.goods_cards("id")
);

CREATE TABLE public.media_files (
	id SERIAL PRIMARY KEY,
	link VARCHAR NOT NULL,
	card_id SERIAL NOT NULL,
  CONSTRAINT media_files_fkey FOREIGN KEY ("card_id") REFERENCES public.goods_cards("id")
);

CREATE TABLE public.prices (
	id SERIAL PRIMARY KEY,
	card_id SERIAL NOT NULL,
	price NUMERIC(10,2),
	discount NUMERIC(10,2),
  CONSTRAINT prices_fkey FOREIGN KEY ("card_id") REFERENCES public.goods_cards("id")
);

CREATE TABLE public.sizes (
	id SERIAL PRIMARY KEY,
	card_id SERIAL NOT NULL,
	techSize VARCHAR(40) NOT NULL,
	title text NOT NULL,
	barcode VARCHAR(128) NOT NULL,
  CONSTRAINT sizes_fkey FOREIGN KEY ("card_id") REFERENCES public.goods_cards("id")
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- не будем сносить всю бд
select 1
;
-- +goose StatementEnd

