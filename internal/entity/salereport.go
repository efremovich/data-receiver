package entity

import "time"

type SaleReport struct {
	ExternalID    string
	DateFrom      string
	DateTo        string
	UpdatedAt     time.Time
	ContractNuber string
	PartID        int
	DocType       string

	Quantity          float32
	RetailPrice       float32
	FinisPrice        float32
	RetailSum         float32
	SalePercent       int
	CommissionPercent float32
	OfficeName        string
	OrderDate         time.Time
	SaleDate          time.Time
	DeliveryAmount    float32
	ReturnAmoun       float32
	DeliveryCost      float32
	PuckageType       string
	ProductDiscon     float32

	Pvz Pvz

	Order   Order
	Barcode Barcode
}

type Pvz struct {
	BuyerSale     float32
	BuyerSaleBase float32
	TotalSale     float32
}
