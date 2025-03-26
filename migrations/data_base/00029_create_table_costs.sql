-- + goose Up
-- +goose StatementBegin
-- Создаем основную таблицу costs как секционированную по created_at
CREATE TABLE shop.costs (
    id SERIAL NOT NULL,
    card_id INTEGER, -- ID карточки товара
    cost NUMERIC(10, 2), -- Себестоимость товара
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(), -- Дата создания записи
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() -- Дата обновления записи
) PARTITION BY RANGE (created_at);

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
COMMENT ON COLUMN shop.costs.cost IS 'Себестоимость товара с точностью до двух знаков после запятой';
COMMENT ON COLUMN shop.costs.created_at IS 'Дата создания записи, определяет актуальность себестоимости и используется для секционирования';
COMMENT ON COLUMN shop.costs.updated_at IS 'Дата последнего обновления записи';

-- Создаем индексы для оптимизации запросов
CREATE INDEX costs_card_id_idx ON shop.costs (card_id);
CREATE INDEX costs_created_at_idx ON shop.costs (created_at);

CREATE OR REPLACE FUNCTION shop.create_costs_partitions()
RETURNS VOID AS $$
DECLARE
    partition_date TIMESTAMP;
    partition_name TEXT;
    start_date TEXT;
    end_date TEXT;
BEGIN
    FOR i IN 0..11 LOOP -- Создаем партиции на год вперед
        partition_date := date_trunc('month', current_date + (i * interval '1 month'));
        partition_name := 'shop.costs_' || to_char(partition_date, 'YYYY_MM');
        start_date := to_char(partition_date, 'YYYY-MM-DD');
        end_date := to_char(partition_date + interval '1 month', 'YYYY-MM-DD');

        EXECUTE format('
            CREATE TABLE IF NOT EXISTS %I PARTITION OF shop.costs
            FOR VALUES FROM (%L) TO (%L)',
            partition_name, start_date, end_date
        );

        -- Устанавливаем владельца для новой партиции
        EXECUTE format('ALTER TABLE %I OWNER TO erp_db_usr', partition_name);
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Выполняем функцию для создания партиций
SELECT shop.create_costs_partitions();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE shop.costs;
-- +goose StatementEnd
