package entity

import "time"

type Sale struct {
	ID         int64
	ExternalID string // Уникальный номер продажи

	Price      float32 // Цена без скидки
	DiscountP  float32 // Скидка продавца
	DiscountS  float32 // Скидка на маркетпрейсе
	FinalPrice float32 // Конечная цена
	Type       string  // Тип продажи
	ForPay     float32 // Сумма к перечислению

	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time

	Status    *Status
	Region    *Region
	Warehouse *Warehouse
	Seller    *Seller
	Card      *Card
	PriceSize *PriceSize
	Barcode   *Barcode
	Order     *Order
	Size      *Size
}
