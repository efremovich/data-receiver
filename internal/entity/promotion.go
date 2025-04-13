package entity

import "time"

type Promotion struct {
	ID          int64     // Уникальный идентификатор кампании в системе
	ExternalID  int64     // Внешний идентификатор кампании (из API Wildberries)
	Name        string    // Название рекламной кампании
	DateStart   time.Time // Дата и время начала кампании
	DateEnd     time.Time // Дата и время окончания кампании
	Type        string    // Тип кампании (например, "auto" или "manual")
	Status      string    // Текущий статус кампании (например, "active", "paused")
	AdvertType  string    // Тип рекламы (например, "search", "card")
	BudgetSpent float64   // Потраченный бюджет кампании
	BudgetTotal float64   // Общий бюджет кампании
	CreatedAt   time.Time // Дата и время создания записи
	UpdatedAt   time.Time // Дата и время последнего обновления записи
	SellerID    int64     // Идентификатор продавца
}

type PromotionInfo struct {
	Promotion            // Основная информация о кампании
	CardID       int64   // Идентификатор карточки товара
	NmID         int64   // Артикул товара (номенклатура Wildberries)
	Subject      string  // Категория товара
	Brand        string  // Бренд товара
	Clicks       int     // Количество кликов по рекламе
	Views        int     // Количество показов рекламы
	CTR          float64 // Click-Through Rate (отношение кликов к показам)
	Orders       int     // Количество заказов
	CostPerClick float64 // Средняя стоимость клика
	CostPerOrder float64 // Средняя стоимость заказа
}

type PromotionStats struct {
	PromotionID      int64     // Идентификатор связанной кампании
	ItemID           int64     // Идентификатор банера
	BrandID          int64     // Бренд
	CategoryName     string    // Название категории
	AdvertType       string    // Тип рекламы
	Place            int       // Место на странице
	Views            int       // Количество показов
	Clicks           int       // Количество кликов
	ConversionRate   float64   // Conversion rate количество заказов к общему количество посещений медиакaмпании
	ClickThroughRate float64   // Click-Through Rate показатель кликабельности
	DateFrom         time.Time // Время начала размещения
	DateTo           time.Time // Время завершения размещения
	SubjectName      string    // Родительская категория
	Atbs             int       // Количество добавления товаров к корзину
	Orders           int       // Количество заказов
	Price            float64   // Стоимость размещения
	CPC              float64   // Cost Per Click стоимость клика
	Status           int       // Статус медиакампании
	Expence          int       // Стоимость размещения банера
	Cr1              float64   // Отношение количества добавлений в корзину к количеству кликов
	Cr2              int       // Отношение количества заказов к количеству добавлений в корзину
}

// Campaign представляет рекламную кампанию с агрегированными метриками
type Campaign struct {
	ID           int           // Уникальный идентификатор кампании
	ExternalID   int64         // Внешний идентификатор кампании (из API Wildberries)
	DateStart    time.Time     // Дата и время начала кампании
	DateEnd      time.Time     // Дата и время окончания кампании
	Views        int           // Общее количество просмотров
	Clicks       int           // Общее количество кликов
	CTR          float64       // Click-through rate (кликабельность) в процентах
	CPC          float64       // Средняя стоимость клика (Cost Per Click)
	Spent        float64       // Общий бюджет кампании
	AddToCarts   int           // Количество добавлений в корзину
	Orders       int           // Количество заказов
	CR           int           // Conversion Rate (конверсия в заказы)
	SKUs         int           // Количество уникальных товаров в заказах
	OrderAmount  float64       // Общая сумма заказов
	DailyStats   []DayStat     // Детальная статистика по дням
	BoosterStats []BoosterStat // Данные о позициях в поиске
}

// DayStat содержит детализированную статистику за конкретный день
type DayStat struct {
	Date        time.Time      // Дата статистики
	Views       int            // Просмотры за день
	Clicks      int            // Клики за день
	CTR         float64        // CTR за день
	CPC         float64        // CPC за день
	Spent       float64        // Затраты за день
	AddToCarts  int            // Добавления в корзину за день
	Orders      int            // Заказы за день
	CR          int            // Конверсия за день
	SKUs        int            // Товары в заказах за день
	OrderAmount float64        // Сумма заказов за день
	Platforms   []PlatformStat // Статистика по платформам
}

// PlatformStat содержит метрики для конкретной платформы (сайт, моб. приложение)
type PlatformStat struct {
	Views       int           // Просмотры на платформе
	Clicks      int           // Клики на платформе
	CTR         float64       // CTR на платформе
	CPC         float64       // CPC на платформе
	Spent       float64       // Затраты на платформе
	AddToCarts  int           // Добавления в корзину с платформы
	Orders      int           // Заказы с платформы
	CR          int           // Конверсия с платформы
	SKUs        int           // Товары в заказах с платформы
	OrderAmount float64       // Сумма заказов с платформы
	Products    []ProductStat // Статистика по товарам на платформе
	AppType     int           // Тип платформы (1-сайт, 32-Android, 64-iOS)
}

// ProductStat содержит метрики для конкретного товара
type ProductStat struct {
	Views       int     // Просмотры товара
	Clicks      int     // Клики по товару
	CTR         float64 // CTR товара
	CPC         float64 // CPC товара
	Spent       float64 // Затраты на товар
	AddToCarts  int     // Добавления товара в корзину
	Orders      int     // Заказы товара
	CR          int     // Конверсия товара
	SKUs        int     // Количество заказанного товара
	OrderAmount float64 // Сумма заказов товара
	Name        string  // Название товара
	NmID        int64   // Артикул товара (Wildberries)
}

// BoosterStat содержит данные о позиции товара в поисковой выдаче
type BoosterStat struct {
	Date        time.Time // Дата измерения позиции
	NmID        int       // Артикул товара
	AvgPosition int       // Средняя позиция в поиске
}
