package entity

import "time"

// Promotion представляет рекламную кампанию с агрегированными метриками.
type Promotion struct {
	ID             int64            // Уникальный идентификатор кампании
	ExternalID     int64            // Внешний идентификатор кампании (из API Wildberries)
	Name           string           // Название кампании
	Type           int              // Тип кампании (1-реклама, 2-продвижение)
	Status         int              // Статус кампании (1-активная, 2-заблокированная)
	ChangeTime     time.Time        // Дата и время изменения кампании
	CreateTime     time.Time        // Дата и время создания кампании
	DateStart      time.Time        // Дата и время начала кампании
	DateEnd        time.Time        // Дата и время окончания кампании
	Views          int              // Общее количество просмотров
	Clicks         int              // Общее количество кликов
	CTR            float64          // Click-through rate (кликабельность) в процентах
	CPC            float64          // Средняя стоимость клика (Cost Per Click)
	Spent          float64          // Общий бюджет кампании
	Orders         int              // Количество заказов
	CR             float64          // Conversion Rate (конверсия в заказы)
	SHKs           int              // Количество уникальных товаров в заказах
	OrderAmount    float64          // Общая сумма заказов
	PromotionStats []PromotionStats // Статистика за конкретный день
	SellerID       int64            // ID продавца
}

// PromotionStats содержит детализированную статистику за конкретный день.
type PromotionStats struct {
	ID          int64     // Уникальный идентификатор кампании
	Date        time.Time // Дата статистики
	Views       int       // Просмотры за день
	Clicks      int       // Клики за день
	CTR         float64   // CTR за день
	CPC         float64   // CPC за день
	Spent       float64   // Затраты за день
	Orders      int       // Заказы за день
	CR          float64   // Конверсия за день
	SHKs        int       // Товары в заказах за день
	OrderAmount float64   // Сумма заказов за день
	AppType     int       // Тип платформы (1-сайт, 32-Android, 64-iOS)

	PromotionID    int64 // Идентификатор рекламной компании.
	CardID         int64 // Карточка товара
	SellerID       int64 // ID продавца
	CardExternalID int64 // Внешний код карточки товара
}
