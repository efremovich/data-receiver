package entity

import "time"

type Order struct {
	ID         int64
	ExternalID string

	Price                float64 // Цена на сайте
	PriceWithoutDiscount float64 // Цена без скидок 4800 руб.
	PriceFinal           float64 // Реальная цена по которой была продажи 428.28

	Type      string
	Direction string
	Sale      float64
	IsCancel  bool

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
	Size      *Size
}
