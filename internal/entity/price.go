package entity

import "time"

type PriceSize struct {
	ID                   int64
	Price                float64 // Цена со скидкой 480 руб.
	PriceWithoutDiscount float64 // Цена без скидок 4800 руб.
	PriceFinish          float64 // Реальная цена по которой была продажи 428.28
	UpdatedAt            time.Time

	CardID int64
	SizeID int64
}

type PriceHistory struct {
	ID           int64
	Price        float32 // Значения цены
	Discount     float32 // Сумма скидки
	SpecialPrice float32 // Специальная цена
	UpdatedAt    time.Time

	PriceSizeID int64
}
