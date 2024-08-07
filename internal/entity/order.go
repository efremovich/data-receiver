package entity

import "time"

type Order struct {
	ID           int64
	ExternalID        string
	Price        float32
	Quantity     int
	Discount     float32
	SpecialPrice float32
	Status       string
	Type         string
	Direction    string
	CreatedAt    time.Time
	UpdatedAt    time.Time

	WarehouseID int64
	SellerID    int64
	CardID      int64
}
