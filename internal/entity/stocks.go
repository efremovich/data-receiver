package entity

import "time"

type Stocks struct {
	ID        int64
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time

	BarcodeID   int64
	WarehouseID int64
	CardID      int64
	SellerID    int64
}
