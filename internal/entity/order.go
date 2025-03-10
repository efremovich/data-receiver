package entity

import "time"

type Order struct {
	ID         int64
	ExternalID string

	Price     float32
	Type      string
	Direction string
	Sale      float32

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
