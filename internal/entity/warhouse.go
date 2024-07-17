package entity

type Warehouse struct {
	ID       int64
	ExternalID    string
	Title    string
	Address  string
	Type     string

  SellerID int64
}
