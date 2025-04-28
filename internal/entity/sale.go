package entity

import "time"

type Sale struct {
	ID         int64
	ExternalID string // Уникальный номер продажи

	Price      float64 // Цена без скидки
	DiscountP  float64 // Скидка продавца
	DiscountS  float64 // Скидка на маркетпрейсе
	FinalPrice float64 // Конечная цена
	Type       string  // Тип продажи
	ForPay     float64 // Сумма к перечислению

	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time

	Status    *Status
	Region    *Region
	Warehouse *Warehouse
	Seller    *MarketPlace
	Card      *Card
	PriceSize *PriceSize
	Barcode   *Barcode
	Order     *Order
	Size      *Size
}
