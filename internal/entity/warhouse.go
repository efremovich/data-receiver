package entity

type Warehouse struct {
	ID      int64
	ExtID   int64
	Name    string
	Address string
	Type    string

	Seller Seller
}
