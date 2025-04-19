-- + goose Up
-- +goose StatementBegin
-- Создание таблицы promotions
-- Таблица рекламных кампаний
CREATE TABLE shop.promotions (
    id SERIAL PRIMARY KEY, -- Уникальный идентификатор кампании
    external_id BIGINT NOT NULL, -- Внешний идентификатор кампании (из API Wildberries)
    name VARCHAR(255) NOT NULL, -- Название кампании
    type INT NOT NULL, -- Тип кампании (1-реклама, 2-продвижение)
    status INT NOT NULL, -- Статус кампании (1-активная, 2-заблокированная)
    change_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Дата и время изменения кампании
    create_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- Дата и время создания кампании
    date_start TIMESTAMP WITH TIME ZONE NOT NULL, -- Дата и время начала кампании
    date_end TIMESTAMP WITH TIME ZONE NOT NULL, -- Дата и время окончания кампании
    views INT DEFAULT 0, -- Общее количество просмотров
    clicks INT DEFAULT 0, -- Общее количество кликов
    ctr FLOAT DEFAULT 0.0, -- Click-through rate (кликабельность) в процентах
    cpc FLOAT DEFAULT 0.0, -- Средняя стоимость клика (Cost Per Click)
    spent FLOAT DEFAULT 0.0, -- Общий бюджет кампании
    orders INT DEFAULT 0, -- Количество заказов
    cr FLOAT DEFAULT 0.0, -- Conversion Rate (конверсия в заказы)
    shks INT DEFAULT 0, -- Количество уникальных товаров в заказах
    order_amount FLOAT DEFAULT 0.0, -- Общая сумма заказов
    seller_id INT NOT NULL -- ID продавца
);

-- Устанавливаем владельца таблицы
ALTER TABLE shop.promotions OWNER TO erp_db_usr;

ALTER TABLE ONLY shop.promotions ADD CONSTRAINT promotions_seller_id_fkey 
    FOREIGN KEY (seller_id) REFERENCES shop.sellers(id);

-- Добавление индексов
CREATE INDEX idx_promotions_external_id ON shop.promotions (external_id);
CREATE INDEX idx_promotions_type ON shop.promotions (type);
CREATE INDEX idx_promotions_status ON shop.promotions (status);
CREATE INDEX idx_promotions_date_start ON shop.promotions (date_start);
CREATE INDEX idx_promotions_date_end ON shop.promotions (date_end);

-- Добавление комментариев к колонкам
COMMENT ON COLUMN shop.promotions.id IS 'Уникальный идентификатор кампании';
COMMENT ON COLUMN shop.promotions.external_id IS 'Внешний идентификатор кампании (из API Wildberries)';
COMMENT ON COLUMN shop.promotions.name IS 'Название кампании';
COMMENT ON COLUMN shop.promotions.type IS 'Тип кампании (1-реклама, 2-продвижение)';
COMMENT ON COLUMN shop.promotions.status IS 'Статус кампании (1-активная, 2-заблокированная)';
COMMENT ON COLUMN shop.promotions.change_time IS 'Дата и время изменения кампании';
COMMENT ON COLUMN shop.promotions.create_time IS 'Дата и время создания кампании';
COMMENT ON COLUMN shop.promotions.date_start IS 'Дата и время начала кампании';
COMMENT ON COLUMN shop.promotions.date_end IS 'Дата и время окончания кампании';
COMMENT ON COLUMN shop.promotions.views IS 'Общее количество просмотров';
COMMENT ON COLUMN shop.promotions.clicks IS 'Общее количество кликов';
COMMENT ON COLUMN shop.promotions.ctr IS 'Click-through rate (кликабельность) в процентах';
COMMENT ON COLUMN shop.promotions.cpc IS 'Средняя стоимость клика (Cost Per Click)';
COMMENT ON COLUMN shop.promotions.spent IS 'Общий бюджет кампании';
COMMENT ON COLUMN shop.promotions.orders IS 'Количество заказов';
COMMENT ON COLUMN shop.promotions.cr IS 'Conversion Rate (конверсия в заказы)';
COMMENT ON COLUMN shop.promotions.shks IS 'Количество уникальных товаров в заказах';
COMMENT ON COLUMN shop.promotions.order_amount IS 'Общая сумма заказов';

-- Таблица детализированной статистики по дням
CREATE TABLE shop.promotion_stats (
    id SERIAL PRIMARY KEY, -- Уникальный идентификатор кампании
    date DATE NOT NULL, -- Дата статистики
    views INT DEFAULT 0, -- Просмотры за день
    clicks INT DEFAULT 0, -- Клики за день
    ctr FLOAT DEFAULT 0.0, -- CTR за день
    cpc FLOAT DEFAULT 0.0, -- CPC за день
    spent FLOAT DEFAULT 0.0, -- Затраты за день
    orders INT DEFAULT 0, -- Заказы за день
    cr FLOAT DEFAULT 0.0, -- Конверсия за день
    shks INT DEFAULT 0, -- Товары в заказах за день
    order_amount FLOAT DEFAULT 0.0, -- Сумма заказов за день
    app_type INT NOT NULL, -- Тип платформы (1-сайт, 32-Android, 64-iOS)
    promotion_id INT NOT NULL, -- Идентификатор рекламной компании
    card_id INT NOT NULL, -- Карточка товара
    seller_id INT NOT NULL -- ID продавца
);

-- Устанавливаем владельца таблицы
ALTER TABLE shop.promotion_stats OWNER TO erp_db_usr;

ALTER TABLE ONLY shop.promotion_stats ADD CONSTRAINT promotion_stats_card_id_fkey 
    FOREIGN KEY (card_id) REFERENCES shop.cards(id);

ALTER TABLE ONLY shop.promotion_stats ADD CONSTRAINT promotion_stats_seller_id_fkey 
    FOREIGN KEY (seller_id) REFERENCES shop.sellers(id);

ALTER TABLE ONLY shop.promotion_stats ADD CONSTRAINT promotion_stats_promotion_id_fkey 
    FOREIGN KEY (promotion_id) REFERENCES shop.promotions(id);

CREATE INDEX idx_shop_promotion_stats_date ON shop.promotion_stats (date);
CREATE INDEX idx_shop_promotion_stats_promotion_id ON shop.promotion_stats (promotion_id);
CREATE INDEX idx_shop_promotion_stats_app_type ON shop.promotion_stats (app_type);
CREATE INDEX idx_shop_promotion_stats_card_id ON shop.promotion_stats (card_id);

COMMENT ON COLUMN shop.promotion_stats.date IS 'Дата статистики';
COMMENT ON COLUMN shop.promotion_stats.views IS 'Просмотры за день';
COMMENT ON COLUMN shop.promotion_stats.clicks IS 'Клики за день';
COMMENT ON COLUMN shop.promotion_stats.ctr IS 'CTR за день';
COMMENT ON COLUMN shop.promotion_stats.cpc IS 'CPC за день';
COMMENT ON COLUMN shop.promotion_stats.spent IS 'Затраты за день';
COMMENT ON COLUMN shop.promotion_stats.orders IS 'Заказы за день';
COMMENT ON COLUMN shop.promotion_stats.cr IS 'Конверсия за день';
COMMENT ON COLUMN shop.promotion_stats.shks IS 'Товары в заказах за день';
COMMENT ON COLUMN shop.promotion_stats.order_amount IS 'Сумма заказов за день';
COMMENT ON COLUMN shop.promotion_stats.app_type IS 'Тип платформы (1-сайт, 32-Android, 64-iOS)';
COMMENT ON COLUMN shop.promotion_stats.promotion_id IS 'Идентификатор рекламной компании';
COMMENT ON COLUMN shop.promotion_stats.card_id IS 'Карточка товара';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Удаление таблицы shop.promotion_stats
DROP TABLE IF EXISTS shop.promotion_stats CASCADE;

-- Удаление таблицы promotions
DROP TABLE IF EXISTS shop.promotions CASCADE;
-- +goose StatementEnd
