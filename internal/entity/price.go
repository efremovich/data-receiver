package entity

import "time"

type PriceSize struct {
	ID           int64
	Price        float32 // Значения цены
	Discount     float32 // Сумма скидки
	SpecialPrice float32 // Специальная цена
	UpdatedAt    time.Time

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
