package entity

type Shop struct {
	Name       string         `xml:"name"`
	Company    string         `xml:"company"`
	URL        string         `xml:"url"`
	Categories []FeedCategory `xml:"categories>category"`
	Offers     []Offer        `xml:"offers>offer"`
}

type FeedCategory struct {
	ID       string `xml:"id,attr"`
	ParentID string `xml:"parentId,attr,omitempty"`
	Picture  string `xml:"picture,attr,omitempty"`
	Value    string `xml:",chardata"`
}

type Offer struct {
	ID           string   `xml:"id,attr"`
	Available    bool     `xml:"available,attr"`
	GroupID      string   `xml:"group_id,attr"`
	Name         string   `xml:"name"`
	Similar      string   `xml:"similar"`
	Price        int      `xml:"price"`
	Barcode      string   `xml:"barcode"`
	URL          string   `xml:"url"`
	VendorCode   string   `xml:"vendorCode"`
	Sort         int      `xml:"sort"`
	Vendor       string   `xml:"vendor"`
	Rating       float64  `xml:"rating"`
	ReviewsCount int      `xml:"reviews_count"`
	Description  string   `xml:"description"`
	OldPrice     int      `xml:"oldprice"`
	CategoryIDs  []string `xml:"categoryId"`
	Pictures     []string `xml:"picture"`
	Params       []Param  `xml:"param"`
	Badges       []Badge  `xml:"badge"`
}

type Param struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
	Hex   string `xml:"hex,attr,omitempty"`
}

type Badge struct {
	BgColor   string `xml:"bgColor,attr"`
	TextColor string `xml:"textColor,attr"`
	Value     string `xml:",chardata"`
}
