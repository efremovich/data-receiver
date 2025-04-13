-- Создаем последовательность для ID кампаний (если нужно)
CREATE SEQUENCE IF NOT EXISTS campaign_id_seq;

-- Таблица кампаний
CREATE TABLE IF NOT EXISTS campaigns (
    id INTEGER PRIMARY KEY DEFAULT nextval('campaign_id_seq'),
    views INTEGER NOT NULL DEFAULT 0,
    clicks INTEGER NOT NULL DEFAULT 0,
    ctr NUMERIC(10, 2) NOT NULL DEFAULT 0,
    cpc NUMERIC(10, 2) NOT NULL DEFAULT 0,
    spent NUMERIC(10, 2) NOT NULL DEFAULT 0,
    add_to_carts INTEGER NOT NULL DEFAULT 0,
    orders INTEGER NOT NULL DEFAULT 0,
    cr INTEGER NOT NULL DEFAULT 0,
    skus INTEGER NOT NULL DEFAULT 0,
    order_amount NUMERIC(10, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица для хранения дат кампаний (отношение многие-ко-многим)
CREATE TABLE IF NOT EXISTS campaign_dates (
    campaign_id INTEGER REFERENCES campaigns(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    PRIMARY KEY (campaign_id, date)
);

-- Таблица дневной статистики
CREATE TABLE IF NOT EXISTS day_stats (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaigns(id) ON DELETE CASCADE,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    views INTEGER NOT NULL DEFAULT 0,
    clicks INTEGER NOT NULL DEFAULT 0,
    ctr NUMERIC(10, 2) NOT NULL DEFAULT 0,
    cpc NUMERIC(10, 2) NOT NULL DEFAULT 0,
    spent NUMERIC(10, 2) NOT NULL DEFAULT 0,
    add_to_carts INTEGER NOT NULL DEFAULT 0,
    orders INTEGER NOT NULL DEFAULT 0,
    cr INTEGER NOT NULL DEFAULT 0,
    skus INTEGER NOT NULL DEFAULT 0,
    order_amount NUMERIC(10, 2) NOT NULL DEFAULT 0,
    UNIQUE (campaign_id, date)
);

-- Таблица статистики по платформам
CREATE TABLE IF NOT EXISTS platform_stats (
    id SERIAL PRIMARY KEY,
    day_stat_id INTEGER REFERENCES day_stats(id) ON DELETE CASCADE,
    views INTEGER NOT NULL DEFAULT 0,
    clicks INTEGER NOT NULL DEFAULT 0,
    ctr NUMERIC(10, 2) NOT NULL DEFAULT 0,
    cpc NUMERIC(10, 2) NOT NULL DEFAULT 0,
    spent NUMERIC(10, 2) NOT NULL DEFAULT 0,
    add_to_carts INTEGER NOT NULL DEFAULT 0,
    orders INTEGER NOT NULL DEFAULT 0,
    cr INTEGER NOT NULL DEFAULT 0,
    skus INTEGER NOT NULL DEFAULT 0,
    order_amount NUMERIC(10, 2) NOT NULL DEFAULT 0,
    app_type INTEGER NOT NULL -- 1 - сайт, 32 - Android, 64 - IOS
);

-- Таблица статистики по товарам
CREATE TABLE IF NOT EXISTS product_stats (
    id SERIAL PRIMARY KEY,
    platform_stat_id INTEGER REFERENCES platform_stats(id) ON DELETE CASCADE,
    views INTEGER NOT NULL DEFAULT 0,
    clicks INTEGER NOT NULL DEFAULT 0,
    ctr NUMERIC(10, 2) NOT NULL DEFAULT 0,
    cpc NUMERIC(10, 2) NOT NULL DEFAULT 0,
    spent NUMERIC(10, 2) NOT NULL DEFAULT 0,
    add_to_carts INTEGER NOT NULL DEFAULT 0,
    orders INTEGER NOT NULL DEFAULT 0,
    cr INTEGER NOT NULL DEFAULT 0,
    skus INTEGER NOT NULL DEFAULT 0,
    order_amount NUMERIC(10, 2) NOT NULL DEFAULT 0,
    name VARCHAR(255) NOT NULL,
    nm_id BIGINT NOT NULL
);

-- Таблица статистики по позициям в поиске (booster stats)
CREATE TABLE IF NOT EXISTS booster_stats (
    id SERIAL PRIMARY KEY,
    campaign_id INTEGER REFERENCES campaigns(id) ON DELETE CASCADE,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    nm_id INTEGER NOT NULL,
    avg_position INTEGER NOT NULL,
    UNIQUE (campaign_id, date, nm_id)
);

-- Индексы для ускорения запросов
CREATE INDEX IF NOT EXISTS idx_campaign_dates_date ON campaign_dates(date);
CREATE INDEX IF NOT EXISTS idx_day_stats_date ON day_stats(date);
CREATE INDEX IF NOT EXISTS idx_product_stats_nm_id ON product_stats(nm_id);
CREATE INDEX IF NOT EXISTS idx_booster_stats_nm_id ON booster_stats(nm_id);

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_campaign_timestamp
BEFORE UPDATE ON campaigns
FOR EACH ROW EXECUTE FUNCTION update_timestamp();
