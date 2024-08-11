package entity

import "time"

type Stock struct {
	ID               int64
	Quantity         int
	InWayToClient    int
	InWayFromClient  int
	InWayToWarehouse int
	CreatedAt        time.Time
	UpdatedAt        time.Time

	SizeID      int64
	BarcodeID   int64
	WarehouseID int64
	CardID      int64
	SellerID    int64
}

type StockMeta struct {
	Stock     Stock
	PriceSize PriceSize
	Barcode   Barcode
	Warehouse Warehouse
	Wb2Card   Wb2Card
	Ozon2Card Ozon2Card
	Size      Size
}
