-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.sellers (
	id SERIAL PRIMARY KEY,
  title VARCHAR NOT NULL,
  is_enable BOOLEAN DEFAULT TRUE,
  ext_id VARCHAR
);
COMMENT ON TABLE sellers is 'Продавцы';
COMMENT ON COLUMN sellers.id is 'Идентификатор';
COMMENT ON COLUMN sellers.title is 'Наименование продавца';
COMMENT ON COLUMN sellers.is_enable is 'Признак активности';
COMMENT ON COLUMN sellers.ext_id is 'Внешний идентификатор';

INSERT INTO public.sellers (title) VALUES ('wb'), ('ozon'), ('1c');

CREATE TABLE public.brands (
  id SERIAL PRIMARY KEY,
  title VARCHAR NOT NULL,
  seller_id INTEGER NOT NULL,
  CONSTRAINT brands_seller_id_fkey FOREIGN KEY ("seller_id") REFERENCES public.sellers("id")
);
CREATE INDEX brands_seller_id_idx ON brands(seller_id);

COMMENT ON TABLE brands is 'Бренды';
COMMENT ON COLUMN brands.id is 'Идентификатор';
COMMENT ON COLUMN brands.title is 'Наименование бренда';
COMMENT ON COLUMN brands.seller_id is 'Идентификатор продавца';

CREATE TABLE public.categories (
  id SERIAL PRIMARY KEY,
  title VARCHAR NOT NULL,
  seller_id INTEGER NOT NULL,
  CONSTRAINT categories_seller_id_fkey FOREIGN KEY ("seller_id") REFERENCES public.sellers("id")
);
CREATE INDEX categories_seller_id_idx ON categories(seller_id);

COMMENT ON TABLE categories is 'Категории товаров';
COMMENT ON COLUMN categories.id is 'Идентификатор';
COMMENT ON COLUMN categories.title is 'Наименование категории';
COMMENT ON COLUMN categories.seller_id is 'Идентификатор продавца';

CREATE TABLE public.cards (
	id SERIAL PRIMARY KEY,
  vendor_id VARCHAR,
  vendor_code VARCHAR,
  title VARCHAR NOT NULL,
  description TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  brand_id INTEGER,  
  category_id INTEGER 
);
CREATE INDEX cards_vendor_code_idx ON cards(vendor_code);
CREATE INDEX cards_vendor_id_idx ON cards(vendor_id);
CREATE INDEX cards_title_idx ON cards(title);
CREATE INDEX cards_created_at_idx ON cards(created_at);
CREATE INDEX cards_updated_at_idx ON cards(updated_at);

COMMENT ON TABLE cards is 'Товары';
COMMENT ON COLUMN cards.id is 'Внутренний идентификатор';
COMMENT ON COLUMN cards.vendor_id is 'Внутренный код товара (из 1с)';
COMMENT ON COLUMN cards.vendor_code is 'Артикул (из 1с)';
COMMENT ON COLUMN cards.title is 'Наименование номенклатуры';
COMMENT ON COLUMN cards.description is 'Описание номенклатуры';
COMMENT ON COLUMN cards.created_at is 'Дата создания';
COMMENT ON COLUMN cards.updated_at is 'Дата обновления';

CREATE TABLE public.dimensions(
  id SERIAL PRIMARY KEY,
  width INTEGER,
  height INTEGER,
  length INTEGER,

  card_id SERIAL NOT NULL,
  CONSTRAINT dimensions_card_id_fkey FOREIGN KEY ("card_id") REFERENCES public.cards("id")
);
CREATE INDEX dimensions_card_id_idx ON dimensions(card_id);

COMMENT ON TABLE dimensions is 'Размеры';
COMMENT ON COLUMN dimensions.id is 'Идентификатор';
COMMENT ON COLUMN dimensions.width is 'Ширина';
COMMENT ON COLUMN dimensions.height is 'Высота';
COMMENT ON COLUMN dimensions.length is 'Длина';
COMMENT ON COLUMN dimensions.card_id is 'Идентификатор номенклатуры';

CREATE TABLE public.characteristics (
	id SERIAL PRIMARY KEY,
	title VARCHAR NOT NULL,
  CONSTRAINT characteristics_card_id_fkey FOREIGN KEY ("card_id") REFERENCES public.cards("id")
);
CREATE INDEX characteristics_card_id_idx ON characteristics(card_id);

COMMENT ON TABLE characteristics is 'Характеристики';
COMMENT ON COLUMN characteristics.id is 'Идентификатор';
COMMENT ON COLUMN characteristics.title is 'Наименование';
COMMENT ON COLUMN characteristics.value is 'Значение';
COMMENT ON COLUMN characteristics.card_id is 'Идентификатор номенклатуры';

CREATE TABLE public.card_characteristics (
	id SERIAL PRIMARY KEY,
  characteristics_id integer,
	value TEXT[] NOT NULL,
  card_id INTEGER NOT NULL,
  CONSTRAINT characteristics_card_id_fkey FOREIGN KEY ("card_id") REFERENCES public.cards("id")
);
CREATE INDEX card_characteristics_card_id_idx ON card_characteristics(card_id);

COMMENT ON TABLE card_characteristics is 'Характеристики';
COMMENT ON COLUMN card_characteristics.id is 'Идентификатор';
COMMENT ON COLUMN card_characteristics.title is 'Наименование';
COMMENT ON COLUMN card_characteristics.value is 'Значение';
COMMENT ON COLUMN card_characteristics.card_id is 'Идентификатор номенклатуры';

CREATE TABLE public.media_files (
	id SERIAL PRIMARY KEY,
	link VARCHAR NOT NULL,
	card_id INTEGER NOT NULL,
  CONSTRAINT media_files_card_id_fkey FOREIGN KEY ("card_id") REFERENCES public.cards("id")
);
CREATE INDEX media_files_card_id_idx ON media_files(card_id);

COMMENT ON TABLE media_files is 'Медиафайлы';
COMMENT ON COLUMN media_files.id is 'Идентификатор';
COMMENT ON COLUMN media_files.link is 'Ссылка на медиафайл';
COMMENT ON COLUMN media_files.card_id is 'Идентификатор номенклатуры';

CREATE TABLE public.prices (
	id SERIAL PRIMARY KEY,
	price NUMERIC(10,2),
	discount NUMERIC(10,2),
  special_price NUMERIC(10,2),
  seller_id INTEGER NOT NULL,
  CONSTRAINT prices_seller_id_fkey FOREIGN KEY ("seller_id") REFERENCES public.sellers("id"),
	card_id INTEGER NOT NULL,
  CONSTRAINT prices_card_id_fkey FOREIGN KEY ("card_id") REFERENCES public.cards("id"),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
CREATE INDEX prices_card_id_idx ON prices(card_id);
CREATE INDEX prices_seller_id_idx ON prices(seller_id);
CREATE INDEX prices_created_at_idx ON prices(created_at);

COMMENT ON TABLE prices is 'Цены';
COMMENT ON COLUMN prices.id is 'Идентификатор';
COMMENT ON COLUMN prices.price is 'Цена';
COMMENT ON COLUMN prices.discount is 'Скидка';
COMMENT ON COLUMN prices.special_price is 'Спецпредложение';
COMMENT ON COLUMN prices.seller_id is 'Идентификатор продавца';
COMMENT ON COLUMN prices.card_id is 'Идентификатор номенклатуры';
COMMENT ON COLUMN prices.created_at is 'Дата создания';

-- CREATE TABLE public.price_history (
-- id SERIAL PRIMARY KEY,
-- updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
-- card_id SERIAL NOT NULL,
-- CONSTRAINT price_history_card_id_fkey FOREIGN KEY ("card_id") REFERENCES
-- public.cards("id")
-- );
-- CREATE INDEX price_history_card_id_idx ON price_history(card_id);
-- COMMENT ON TABLE price_history is 'История цен';
-- COMMENT ON COLUMN price_history.id is 'Идентификатор';
-- COMMENT ON COLUMN price_history.updated_at is 'Дата обновления';
-- COMMENT ON COLUMN price_history.card_id is 'Идентификатор номенклатуры'; 
CREATE TABLE public.sizes (
	id SERIAL PRIMARY KEY,
	tech_size VARCHAR(40) NOT NULL,
	title text NOT NULL,
	card_id INTEGER NOT NULL,
  CONSTRAINT sizes_card_id_fkey FOREIGN KEY ("card_id") REFERENCES public.cards("id"),
  price_id INTEGER NOT NULL,
  CONSTRAINT sizes_price_id_fkey FOREIGN KEY ("price_id") REFERENCES public.prices("id")
);
CREATE INDEX sizes_card_id_idx ON sizes(card_id);
CREATE INDEX sizes_price_id_idx ON sizes(price_id);

COMMENT ON TABLE sizes is 'Размеры';
COMMENT ON COLUMN sizes.id is 'Идентификатор';
COMMENT ON COLUMN sizes.tech_size is 'Технический обозначение размера';
COMMENT ON COLUMN sizes.title is 'Наименование размера';
COMMENT ON COLUMN sizes.card_id is 'Идентификатор номенклатуры';
COMMENT ON COLUMN sizes.price_id is 'Идентификатор цены';

CREATE TABLE public.barcodes(
  barcode VARCHAR(128) PRIMARY KEY,
  size_id INTEGER NOT NULL,
  CONSTRAINT barcodes_size_id_fkey FOREIGN KEY ("size_id") REFERENCES public.sizes("id"),
  seller_id INTEGER NOT NULL,
  CONSTRAINT barcodes_seller_fkey FOREIGN KEY ("seller_id") REFERENCES public.sellers("id")
);
CREATE INDEX barcodes_size_id_idx ON barcodes(size_id);
CREATE INDEX barcodes_seller_id_idx ON barcodes(seller_id);

COMMENT ON TABLE barcodes is 'Штрихкоды';
COMMENT ON COLUMN barcodes.barcode is 'Штрихкод';
COMMENT ON COLUMN barcodes.size_id is 'Идентификатор размера';
COMMENT ON COLUMN barcodes.seller_id is 'Идентификатор продавца';

CREATE TABLE public.wb2cards(
  nmID INTEGER PRIMARY KEY,
  int INTEGER NOT NULL,
  nmUUID VARCHAR NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  card_id INTEGER NOT NULL,
  CONSTRAINT wb2cards_card_id_fkey FOREIGN KEY ("card_id") REFERENCES public.cards("id")
);
CREATE INDEX wb2cards_card_id_idx ON wb2cards(card_id);

COMMENT ON TABLE wb2cards is 'Товары WB';
COMMENT ON COLUMN wb2cards.nmID is 'Артикул WB';
COMMENT ON COLUMN wb2cards.int is 'Идентификатор КТ';
COMMENT ON COLUMN wb2cards.nmUUID is 'Внуттренний технический идентификатор товара';
COMMENT ON COLUMN wb2cards.created_at is 'Дата создания';
COMMENT ON COLUMN wb2cards.updated_at is 'Дата обновления';
COMMENT ON COLUMN wb2cards.card_id is 'Идентификатор номенклатуры';

CREATE TABLE public.ozon2cards(
  id INTEGER PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  card_id INTEGER NOT NULL,
  CONSTRAINT ozon2cards_card_id_fkey FOREIGN KEY ("card_id") REFERENCES public.cards("id")
);
CREATE INDEX ozon2cards_card_id_idx ON ozon2cards(card_id);

COMMENT ON TABLE ozon2cards is 'Товары Ozon';
COMMENT ON COLUMN ozon2cards.id is 'Идентификатор товара на ozon';
COMMENT ON COLUMN ozon2cards.created_at is 'Дата создания';
COMMENT ON COLUMN ozon2cards.updated_at is 'Дата обновления';


CREATE TABLE public.warehouses(
  id SERIAL PRIMARY KEY,
  ext_id VARCHAR NOT NULL,
  title VARCHAR NOT NULL,
  address VARCHAR,
  type VARCHAR,
  seller_id SERIAL,
  CONSTRAINT warehouses_seller_id_fkey FOREIGN KEY ("seller_id") REFERENCES public.sellers("id")
);
CREATE INDEX warehouses_seller_id_idx ON warehouses(seller_id);

COMMENT ON TABLE warehouses is 'Склады';
COMMENT ON COLUMN warehouses.id is 'Идентификатор';
COMMENT ON COLUMN warehouses.ext_id is 'Внешний идентификатор';
COMMENT ON COLUMN warehouses.title is 'Наименование склада';
COMMENT ON COLUMN warehouses.address is 'Адрес склада';
COMMENT ON COLUMN warehouses.type is 'Тип склада';
COMMENT ON COLUMN warehouses.seller_id is 'Идентификатор продавца';

CREATE TABLE public.orders(
  id SERIAL PRIMARY KEY,
  ext_id VARCHAR NOT NULL,
  price NUMERIC(10,2) NOT NULL,
  quantity INTEGER NOT NULL,
  discount NUMERIC(10,2),
  special_price NUMERIC(10,2),
  status VARCHAR NOT NULL,
  direction VARCHAR NOT NULL,
  type VARCHAR,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),

  warehouse_id INTEGER NOT NULL,
  CONSTRAINT orders_warehouse_idfkey FOREIGN KEY ("warehouse_id") REFERENCES public.warehouses("id"),
  seller_id INTEGER NOT NULL,
  CONSTRAINT orders_seller_id_fkey FOREIGN KEY ("seller_id") REFERENCES public.sellers("id"),
  card_id INTEGER NOT NULL,
  CONSTRAINT orders_card_id_fkey FOREIGN KEY ("card_id") REFERENCES public.cards("id")
);
CREATE INDEX orders_card_id_idx ON orders(card_id);
CREATE INDEX orders_warehouse_id_idx ON orders(warehouse_id);
CREATE INDEX orders_ext_id_idx ON orders(ext_id);
CREATE INDEX orders_seller_id_idx ON orders(seller_id);
CREATE INDEX orders_created_at_idx ON orders(created_at);
CREATE INDEX orders_updated_at_idx ON orders(updated_at);

COMMENT ON TABLE orders is 'Заказы';
COMMENT ON COLUMN orders.id is 'Идентификатор';
COMMENT ON COLUMN orders.ext_id is 'Внешний идентификатор';
COMMENT ON COLUMN orders.price is 'Цена';
COMMENT ON COLUMN orders.quantity is 'Количество';
COMMENT ON COLUMN orders.discount is 'Скидка';
COMMENT ON COLUMN orders.special_price is 'Спеццена';
COMMENT ON COLUMN orders.warehouse_id is 'Идентификатор склада';
COMMENT ON COLUMN orders.status is 'Статус заказа';
COMMENT ON COLUMN orders.direction is 'Направление заказа';
COMMENT ON COLUMN orders.type is 'Тип заказа';
COMMENT ON COLUMN orders.card_id is 'Идентификатор номенклатуры';
COMMENT ON COLUMN orders.seller_id is 'Идентификатор продавца';

CREATE TABLE public.stocks(
  id SERIAL PRIMARY KEY,
  quantity INTEGER NOT NULL,
  in_way_to_client INTEGER,
  in_way_from_client INTEGER,
  in_way_to_warehouse INTEGER,

  warehouse_id INTEGER NOT NULL,
  CONSTRAINT stocks_warehouse_id_fkey FOREIGN KEY ("warehouse_id") REFERENCES public.warehouses("id"),
  card_id INTEGER NOT NULL,
  CONSTRAINT stocks_card_id_fkey FOREIGN KEY ("card_id") REFERENCES public.cards("id"),
  barcode VARCHAR NOT NULL,
  CONSTRAINT stocks_barcode_fkey FOREIGN KEY ("barcode") REFERENCES public.barcodes("barcode"),
  seller_id INTEGER NOT NULL,
  CONSTRAINT stocks_seller_id_fkey FOREIGN KEY ("seller_id") REFERENCES public.sellers("id"),
  size_id INTEGER NOT NULL,
  CONSTRAINT stocks_size_id_fkey FOREIGN KEY ("size_id") REFERENCES public.sizes("id"),

  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
CREATE INDEX stocks_card_id_idx ON stocks(card_id);
CREATE INDEX stocks_warehouse_id_idx ON stocks(warehouse_id);
CREATE INDEX stocks_barcode_idx ON stocks(barcode);
CREATE INDEX stocks_size_id_idx ON stocks(size_id);
CREATE INDEX stocks_created_at_idx ON stocks(created_at);
CREATE INDEX stocks_updated_at_idx ON stocks(updated_at);

COMMENT ON TABLE stocks is 'Складские остатки';
COMMENT ON COLUMN stocks.id is 'Идентификатор';
COMMENT ON COLUMN stocks.quantity is 'Количество на складе';
COMMENT ON COLUMN stocks.in_way_to_client is 'Количество в пути к клиенту';
COMMENT ON COLUMN stocks.in_way_from_client is 'Количество в пути от клиента';
COMMENT ON COLUMN stocks.warehouse_id is 'Идентификатор склада';
COMMENT ON COLUMN stocks.card_id is 'Идентификатор номенклатуры';
COMMENT ON COLUMN stocks.barcode is 'Штрихкод';

CREATE TABLE public.event_enum (
    id SERIAL PRIMARY KEY,
    event_desc VARCHAR NOT NULL
);
INSERT INTO public.event_enum (event_desc) VALUES ('CREATED'), ('SUCCESS'), ('GOT_AGAIN'), ('ERROR'), ('SEND_TASK_NEXT');

CREATE TABLE public.jobs(
  id BIGSERIAL PRIMARY KEY,
  pub VARCHAR(128) NOT NULL,
  status VARCHAR(128) NOT NULL,
  event_type_id INTEGER NOT NULL,
  CONSTRAINT jobs_event_type_id_fkey FOREIGN KEY ("event_type_id") REFERENCES public.event_enum("id"), 
  description VARCHAR, 
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
CREATE INDEX jobs_created_at_idx ON jobs(created_at);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- не будем сносить всю бд
select 1
;
-- +goose StatementEnd

