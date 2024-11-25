package entity

type Warehouse struct {
	ID         int64
	ExternalID int64
	Title      string
	Address    string
	TypeName   string

	TypeID   int64
	SellerID int64
}

type WarehouseType struct {
	ID    int64
	Title string
}
