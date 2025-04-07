-- + goose Up
-- +goose StatementBegin
-- Создаем основную таблицу costs как секционированную по created_at
CREATE TABLE shop.costs (
    id SERIAL NOT NULL,
    card_id INTEGER, -- ID карточки товара
    amount NUMERIC(10, 2), -- Себестоимость товара
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(), -- Дата создания записи
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() -- Дата обновления записи
);

-- Устанавливаем владельца таблицы
ALTER TABLE shop.costs OWNER TO erp_db_usr;

-- Добавляем первичный ключ
ALTER TABLE ONLY shop.costs ADD CONSTRAINT costs_pkey PRIMARY KEY (id);

-- Добавляем внешний ключ на таблицу shop.cards
ALTER TABLE ONLY shop.costs ADD CONSTRAINT costs_card_id_fkey 
    FOREIGN KEY (card_id) REFERENCES shop.cards(id);

-- Добавляем комментарии к колонкам
COMMENT ON TABLE shop.costs IS 'Таблица для хранения себестоимости товаров, секционированная по дате создания';
COMMENT ON COLUMN shop.costs.id IS 'Уникальный идентификатор записи';
COMMENT ON COLUMN shop.costs.card_id IS 'Идентификатор карточки товара из таблицы shop.cards';
COMMENT ON COLUMN shop.costs.amount IS 'Себестоимость товара с точностью до двух знаков после запятой';
COMMENT ON COLUMN shop.costs.created_at IS 'Дата создания записи, определяет актуальность себестоимости и используется для секционирования';
COMMENT ON COLUMN shop.costs.updated_at IS 'Дата последнего обновления записи';

-- Создаем индексы для оптимизации запросов
CREATE INDEX costs_card_id_idx ON shop.costs (card_id);
CREATE INDEX costs_created_at_idx ON shop.costs (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop.costs;
-- +goose StatementEnd
