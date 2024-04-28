package entity

import "time"

type Order struct {
	ID          int64
	ExtID       int64
	Status      string
	Description string
	Type        string
	Sale        float32
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Price     Price
	Warehouse Warehouse
  Seller Seller
}
