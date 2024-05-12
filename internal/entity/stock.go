package entity

import "time"

type Stock struct {
	ID int64
	Quantity        int
	InWayToClient   int
	InWayFromClient int
	CreatedAt       time.Time
	UpdatedAt       time.Time

	SizeID      int64
	Barcode     string
	WarehouseID int64
	CardID      int64
	SellerID    int64
}
