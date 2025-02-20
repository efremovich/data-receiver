-- +goose Up
-- +goose StatementBegin

CREATE TABLE shop.salereports (
    id serial NOT NULL, -- Идентификатор
    external_id text NOT NULL, -- Уникальный идентификатор заказа.
    card_id integer NOT NULL, -- Идентификатор товара
    seller_id integer NOT NULL, -- Иденификатор продавца
    barcode text, -- Штрихкод
    order_date timestamp, -- Дата заказа
    sale_date timestamp, -- Дата продажи
    report_id text -- ID отчета
    report_date_from timestamp -- Отчетный период с 
    report_date_to timestamp -- Отчетный период по
    order_id integer, -- Идентификатор заказа
    sale_id integer, -- Идентификатор продажи
    updated_at timestamp NOT NULL, -- Дата обновления данных
    contract_number text, -- Договор
    part_id integer NOT NULL, -- Номер поставки
    doc_type text, -- Тип документа
    quantity numeric(10,2), -- Количество
    retail_price numeric(10,2), -- Цена розничная
    finis_price numeric(10,2), -- Цена розничная с учетом согласованной скидки 
    retail_sum numeric(10,2), -- Сумма продажи (возврата)
    sale_percent integer, -- Процент скидки
    commission_percent numeric(10,2), -- Процент комиссии
    warehouse_name text, -- Наименование склада
    operation_name text, -- Наименование операции
    delivery_amount numeric(10,2), -- Количество доставок
    return_amount numeric(10,2), -- Количество возвратов
    delivery_cost numeric(10,2), -- Стоимость доставки
    package_type text, -- Тип упаковки
    product_discount numeric(10,2), -- Финальная скидка
    buyer_discount numeric(10,2), -- Скидка постоянного покупателя
    base_ratio_discount numeric(10,2), -- Размер кВВ без НДС, % базовый
    total_ratio_discount numeric(10,2), -- Итоговый кВВ без НДС, %
    reduction_rating_ratio numeric(10,2), -- Размер снижения кВВ из-за рейтинга
    reduction_promotion_ratio numeric(10,2), -- Размер снижения кВВ из-за акции
    reward_ratio numeric(10,2), -- Вознаграждение с продаж до вычета услуг поверенного, без НДС
    for_pay numeric(10,2), -- К перечислению продавцу за реализованный товар
    reward numeric(10,2), -- Возмещение за выдачу и возврат товаров на ПВЗ
    acquiring_fee numeric(10,2), -- Возмещение издержек по эквайрингу.
    acquiring_percent numeric(10,2), -- Размер комиссии за эквайринг без НДС, %
    acquiring_bank text, -- Банк экварйрер
    seller_reward numeric(10,2), -- Вознаграждение маркетплейса без НДС
    office_id integer, -- ID офиса
    office_name text, -- Наименование офиса
    supplier_id integer, -- ID поставщика
    supplier_name text, -- Наименование поставщика
    supplier_inn text, -- ИНН поставщика
    declaration_number text, -- Номер таможенной декларации
    bonus_type_name text, -- Штрафы или доплаты
    sticker_id text, -- Цифровое значение стикера
    country_of_sale text, -- Страна продажи
    penalty numeric(10,2), -- Штрафы
    additional_payment numeric(10,2), -- Доплаты
    rebill_logistic_cost numeric(10,2), -- Стоимость возмещения издержек перевозки
    rebill_logistic_org text, -- Организация перевозки
    kiz text, -- Код маркировки товара
    storage_fee numeric(10,2), -- Стоимость хранения
    deduction numeric(10,2), -- Прочие удержания и выплаты
    acceptance numeric(10,2) -- Стоимость платной приемки
);

ALTER TABLE shop.salereports OWNER TO erp_db_usr;
ALTER TABLE ONLY shop.salereports
    ADD CONSTRAINT salereports_pkey PRIMARY KEY (id);
ALTER TABLE ONLY shop.salereports
    ADD CONSTRAINT salereports_card_id_fkey FOREIGN KEY (card_id) REFERENCES shop.cards(id);
ALTER TABLE ONLY shop.salereports
    ADD CONSTRAINT salereports_seller_id_fkey FOREIGN KEY (seller_id) REFERENCES shop.sellers(id);
ALTER TABLE ONLY shop.salereports
    ADD CONSTRAINT salereports_order_id_fkey FOREIGN KEY (order_id) REFERENCES shop.orders(id);
ALTER TABLE ONLY shop.salereports
    ADD CONSTRAINT salereports_sale_idfkey FOREIGN KEY (sale_id) REFERENCES shop.sales(id);

CREATE INDEX salereports_card_id_idx ON shop.salereports USING btree (card_id);
CREATE INDEX salereports_order_id_idx ON shop.salereports USING btree (order_id);
CREATE INDEX salereports_sale_id_idx ON shop.salereports USING btree (sale_id);
CREATE INDEX salereports_updated_at_idx ON shop.salereports USING btree (updated_at);
CREATE INDEX salereports_sale_date_idx ON shop.salereports USING btree (sale_date);
CREATE INDEX salereports_order_date_idx ON shop.salereports USING btree (order_date);
CREATE INDEX salereports_external_id_idx ON shop.salereports USING btree (external_id);

COMMENT ON COLUMN shop.salereports.external_id IS 'Уникальный идентификатор заказа';
COMMENT ON COLUMN shop.salereports.card_id IS 'ИД товара';
COMMENT ON COLUMN shop.salereports.saller_id IS 'ИД маркетплейса';
COMMENT ON COLUMN shop.salereports.barcode IS 'Штрихкод';
COMMENT ON COLUMN shop.salereports.name IS 'Наименование';
COMMENT ON COLUMN shop.salereports.order_id IS 'ИД заказа';
COMMENT ON COLUMN shop.salereports.sale_id IS 'ИД продажи';
COMMENT ON COLUMN shop.salereports.updated_at IS 'Дата обновления данных';
COMMENT ON COLUMN shop.salereports.contract_number IS 'Договор';
COMMENT ON COLUMN shop.salereports.part_id IS 'Номер поставки';
COMMENT ON COLUMN shop.salereports.doc_type IS 'Тип документа';
COMMENT ON COLUMN shop.salereports.quantity IS 'Количество';
COMMENT ON COLUMN shop.salereports.retail_price IS 'Цена розничная';
COMMENT ON COLUMN shop.salereports.finis_price IS 'Цена розничная с учетом согласованной скидки';
COMMENT ON COLUMN shop.salereports.retail_sum IS 'Сумма продажи (возврата)';
COMMENT ON COLUMN shop.salereports.sale_percent IS 'Процент скидки';
COMMENT ON COLUMN shop.salereports.commission_percent IS 'Процент комиссии';
COMMENT ON COLUMN shop.salereports.warehouse_name IS 'Наименование склада';
COMMENT ON COLUMN shop.salereports.order_date IS 'Дата заказа';
COMMENT ON COLUMN shop.salereports.sale_date IS 'Дата продажи';
COMMENT ON COLUMN shop.salereports.operation_name IS 'Наименование операции';
COMMENT ON COLUMN shop.salereports.delivery_amount IS 'Количество доставок';
COMMENT ON COLUMN shop.salereports.return_amount IS 'Количество возвратов';
COMMENT ON COLUMN shop.salereports.delivery_cost IS 'Стоимость доставки';
COMMENT ON COLUMN shop.salereports.package_type IS 'Тип упаковки';
COMMENT ON COLUMN shop.salereports.product_discount IS 'Финальная скидка';
COMMENT ON COLUMN shop.salereports.buyer_discount IS 'Скидка постоянного покуптеля';
COMMENT ON COLUMN shop.salereports.base_ratio_discount IS 'Размер кВВ без НДС, % базовый';
COMMENT ON COLUMN shop.salereports.total_ratio_discount IS 'Итоговый кВВ без НДС, %';
COMMENT ON COLUMN shop.salereports.reduction_rating_ratio IS 'Размер снижения кВВ из-за рейтинга';
COMMENT ON COLUMN shop.salereports.reduction_promotion_ratio IS 'Размер снижения кВВ из-за акции';
COMMENT ON COLUMN shop.salereports.reward_ratio IS 'Вознаграждение с продаж до вычета услуг поверенного, без НДС';
COMMENT ON COLUMN shop.salereports.for_pay IS 'К перечислению продавцу за реализованный товар';
COMMENT ON COLUMN shop.salereports.reward IS 'Возмещение за выдачу и возврат товаров на ПВЗ';
COMMENT ON COLUMN shop.salereports.acquiring_fee IS 'Возмещение издержек по эквайрингу';
COMMENT ON COLUMN shop.salereports.acquiring_percent IS 'Размер комиссии за эквайринг без НДС, %';
COMMENT ON COLUMN shop.salereports.acquiring_bank IS 'Банк экварйрер';
COMMENT ON COLUMN shop.salereports.seller_reward IS 'Вознаграждение маркетплейса без НДС';
COMMENT ON COLUMN shop.salereports.office_id IS 'ID офиса';
COMMENT ON COLUMN shop.salereports.office_name IS 'Наименование офиса';
COMMENT ON COLUMN shop.salereports.supplier_id IS 'ID поставщика';
COMMENT ON COLUMN shop.salereports.supplier_name IS 'Наименование поставщика';
COMMENT ON COLUMN shop.salereports.supplier_inn IS 'ИНН поставщика';
COMMENT ON COLUMN shop.salereports.declaration_number IS 'Номер таможенной декларации';
COMMENT ON COLUMN shop.salereports.bonus_type_name IS 'Штрафы или доплаты';
COMMENT ON COLUMN shop.salereports.sticker_id IS 'Цифровое значение стикера';
COMMENT ON COLUMN shop.salereports.country_of_sale IS 'Страна продажи';
COMMENT ON COLUMN shop.salereports.penalty IS 'Штрафы';
COMMENT ON COLUMN shop.salereports.additional_payment IS 'Доплаты';
COMMENT ON COLUMN shop.salereports.rebill_logistic_cost IS 'Стоимость возмещения издержек перевозки';
COMMENT ON COLUMN shop.salereports.rebill_logistic_org IS 'Организация перевозки';
COMMENT ON COLUMN shop.salereports.kiz IS 'Код маркировки товара';
COMMENT ON COLUMN shop.salereports.storage_fee IS 'Стоимость хранения';
COMMENT ON COLUMN shop.salereports.deduction IS 'Прочие удержания и выплаты';
COMMENT ON COLUMN shop.salereports.acceptance IS 'Стоимость платной приемки';

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX sales_card_id_idx;
DROP TABLE shop.salereports;
-- +goose StatementEnd
