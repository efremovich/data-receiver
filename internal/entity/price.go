package entity

import "time"

type Price struct {
	ID           int64
	Price        float32 // Значения цены
	Discount     float32 // Сумма скидки
	SpecialPrice float32 // Специальная цена
	CreatedAt    time.Time
	UpdatedAt    time.Time

	SellerID int64
	CardID   int64
}

// TODO Подумать, нужна ли нам эта структура
type PriceHistory struct {
	ID        int64
	CreatedAt time.Time

	PriceID int64
}
