-- +goose Up
-- +goose StatementBegin

-- Создание таблицы sale_reports
CREATE TABLE sale_reports (
    id serial NOT NULL,
    external_id text, -- Уникальный идентификатор заказа.
    updated_at timestamp, -- Дата обновления данных
    quantity numeric(10, 2), -- Количество
    retail_price numeric(10, 2), -- Цена розничная
    retail_amount numeric(10, 2), -- Цена продажи
    sale_percent integer, -- Процент скидки
    commission_percent numeric(10, 2), -- Процент комиссии
    retail_price_withdisc_rub numeric(10, 2), -- Цена розничная с учетом скидок в рублях.
    delivery_amount numeric(10, 2), -- Количество доставок
    return_amount numeric(10, 2), -- Количество возвратов
    delivery_cost numeric(10, 2), -- Стоимость доставки
    pvz_reward numeric(10, 2), -- Возмещение за выдачу на ПВЗ
    seller_reward numeric(10, 2), -- Возмещение марекетплейса без НДС
    seller_reward_with_nds numeric(10, 2), -- Возмещение марекетплейса с НДС
    date_from timestamp, -- Дата начала отчета
    date_to timestamp, -- Дата окончания отчета
    create_report_date timestamp, -- Дата создания отчета
    order_date timestamp, -- Дата заказа
    sale_date timestamp, -- Дата продажи
    transaction_date timestamp, -- Дата транзакции
    sa_name text, -- Артикул продавца
    bonus_type_name text, -- Штрафы или доплаты
    penalty numeric(10, 2), -- Штрафы
    additional_payment numeric(10, 2), -- Доплаты
    acquiring_fee numeric(10, 2), -- Возмещение издержек по эквайрингу
    acquiring_percent numeric(10, 2), -- Размер комиссии за эквайринг без НДС, %
    acquiring_bank text, -- Банк экварйрер
    doc_type text, -- Тип документа
    supplier_oper_name text, -- Обоснование оплаты
    site_country text, -- Страна сайта
    kiz text, -- Код маркировки товара
    storage_fee numeric(10, 2), -- Стоимость хранения
    deduction numeric(10, 2), -- Прочие удержания и выплаты
    acceptance numeric(10, 2), -- Стоимость платной приемки
    pvz_id serial, -- ID пункта выдачи
    barcode text, -- Штрихкод
    size_id serial, -- ID размера
    card_id serial, -- ID карты
    order_id serial, -- ID заказа
    warehouse_id serial, -- ID склада
    seller_id serial -- ID продавца
);

ALTER TABLE shop.sale_reports OWNER TO erp_db_usr;

ALTER TABLE ONLY shop.sale_reports
	ADD CONSTRAINT salereports_pkey PRIMARY KEY (id);

ALTER TABLE ONLY shop.sale_reports
	ADD CONSTRAINT sale_reports_card_id_fkey FOREIGN KEY (card_id) REFERENCES shop.cards (id);

ALTER TABLE ONLY shop.sale_reports
	ADD CONSTRAINT sale_reports_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES shop.sellers (id);

ALTER TABLE ONLY shop.sale_reports
	ADD CONSTRAINT sale_reports_order_id_fkey FOREIGN KEY (order_id) REFERENCES shop.orders (id);

ALTER TABLE ONLY shop.sale_reports
	ADD CONSTRAINT sale_reports_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES shop.warehouses (id);

ALTER TABLE ONLY shop.sale_reports
	ADD CONSTRAINT sale_reports_pvz_id_fkey FOREIGN KEY (pvz_id) REFERENCES shop.pvzs (id);

ALTER TABLE ONLY shop.sale_reports
	ADD CONSTRAINT sale_reports_size_id_fkey FOREIGN KEY (size_id) REFERENCES shop.sizes (id);

CREATE INDEX sale_reports_card_id_idx ON shop.sale_reports USING btree (card_id);

CREATE INDEX sale_reports_order_id_idx ON shop.sale_reports USING btree (order_id);

CREATE INDEX sale_reports_updated_at_idx ON shop.sale_reports USING btree (updated_at);

CREATE INDEX sale_reports_sale_date_idx ON shop.sale_reports USING btree (sale_date);

CREATE INDEX sale_reports_external_id_idx ON shop.sale_reports USING btree (external_id);

COMMENT ON COLUMN shop.sale_reports.id IS 'Уникальный идентификатор отчета';

COMMENT ON COLUMN shop.sale_reports.external_id IS 'Уникальный идентификатор заказа';

COMMENT ON COLUMN shop.sale_reports.updated_at IS 'Дата обновления данных';

COMMENT ON COLUMN shop.sale_reports.quantity IS 'Количество';

COMMENT ON COLUMN shop.sale_reports.retail_price IS 'Цена розничная';

COMMENT ON COLUMN shop.sale_reports.retail_amount IS 'Цена продажи';

COMMENT ON COLUMN shop.sale_reports.sale_percent IS 'Процент скидки';

COMMENT ON COLUMN shop.sale_reports.commission_percent IS 'Процент комиссии';

COMMENT ON COLUMN shop.sale_reports.retail_price_withdisc_rub IS 'Цена розничная с учетом скидок в рублях';

COMMENT ON COLUMN shop.sale_reports.delivery_amount IS 'Количество доставок';

COMMENT ON COLUMN shop.sale_reports.return_amount IS 'Количество возвратов';

COMMENT ON COLUMN shop.sale_reports.delivery_cost IS 'Стоимость доставки';

COMMENT ON COLUMN shop.sale_reports.pvz_reward IS 'Возмещение за выдачу на ПВЗ';

COMMENT ON COLUMN shop.sale_reports.seller_reward IS 'Возмещение марекетплейса без НДС';

COMMENT ON COLUMN shop.sale_reports.seller_reward_with_nds IS 'Возмещение марекетплейса с НДС';

COMMENT ON COLUMN shop.sale_reports.date_from IS 'Дата начала отчета';

COMMENT ON COLUMN shop.sale_reports.date_to IS 'Дата окончания отчета';

COMMENT ON COLUMN shop.sale_reports.create_report_date IS 'Дата создания отчета';

COMMENT ON COLUMN shop.sale_reports.order_date IS 'Дата заказа';

COMMENT ON COLUMN shop.sale_reports.sale_date IS 'Дата продажи';

COMMENT ON COLUMN shop.sale_reports.transaction_date IS 'Дата транзакции';

COMMENT ON COLUMN shop.sale_reports.sa_name IS 'Артикул продавца';

COMMENT ON COLUMN shop.sale_reports.bonus_type_name IS 'Штрафы или доплаты';

COMMENT ON COLUMN shop.sale_reports.penalty IS 'Штрафы';

COMMENT ON COLUMN shop.sale_reports.additional_payment IS 'Доплаты';

COMMENT ON COLUMN shop.sale_reports.acquiring_fee IS 'Возмещение издержек по эквайрингу';

COMMENT ON COLUMN shop.sale_reports.acquiring_percent IS 'Размер комиссии за эквайринг без НДС, %';

COMMENT ON COLUMN shop.sale_reports.acquiring_bank IS 'Банк экварйрер';

COMMENT ON COLUMN shop.sale_reports.doc_type IS 'Тип документа';

COMMENT ON COLUMN shop.sale_reports.supplier_oper_name IS 'Обоснование оплаты';

COMMENT ON COLUMN shop.sale_reports.site_country IS 'Страна сайта';

COMMENT ON COLUMN shop.sale_reports.kiz IS 'Код маркировки товара';

COMMENT ON COLUMN shop.sale_reports.storage_fee IS 'Стоимость хранения';

COMMENT ON COLUMN shop.sale_reports.deduction IS 'Прочие удержания и выплаты';

COMMENT ON COLUMN shop.sale_reports.acceptance IS 'Стоимость платной приемки';

COMMENT ON COLUMN shop.sale_reports.pvz_id IS 'ID пункта выдачи';

COMMENT ON COLUMN shop.sale_reports.barcode IS 'Штрихкод';

COMMENT ON COLUMN shop.sale_reports.size_id IS 'ID размера';

COMMENT ON COLUMN shop.sale_reports.card_id IS 'ID карты';

COMMENT ON COLUMN shop.sale_reports.order_id IS 'ID заказа';

COMMENT ON COLUMN shop.sale_reports.warehouse_id IS 'ID склада';

COMMENT ON COLUMN shop.sale_reports.seller_id IS 'ID продавца';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX sales_card_id_idx;

DROP TABLE shop.salereports;
-- +goose StatementEnd
