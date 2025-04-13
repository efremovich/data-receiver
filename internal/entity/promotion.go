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
