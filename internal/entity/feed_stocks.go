package entity

type Inventory struct {
	Storages     []*Storage `xml:"storages>storage"`
	Availability []*Offer   `xml:"availability>offer"`
}

type Storage struct {
	ID       int64       `xml:"id,attr"`
	Name     string      `xml:"name"`
	City     string      `xml:"city"`
	Type     string      `xml:"type"`
	Address  string      `xml:"address"`
	Lat      string      `xml:"lat"`
	Lon      string      `xml:"lon"`
	Region   string      `xml:"region"`
	WorkTime string      `xml:"work_time"`
	Phone    string      `xml:"phone"`
	Icon     string      `xml:"icon"`
	Seller   SellerStock `xml:"seller"`
}

type SellerStock struct {
	SellerID   int64  `xml:"id,attr"`
	SellerName string `xml:"name,attr"`
}

type OfferStock struct {
	ID       int64   `xml:"id,attr"`
	Storage  int64   `xml:"storage,attr"`
	Quantity float32 `xml:"quantity,attr"`
	Price    float32 `xml:"price,attr"`
	OldPrice float32 `xml:"old_price,attr"`
}
