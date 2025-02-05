package entity

import "encoding/xml"

type YMLCatalog struct {
	XMLName xml.Name `xml:"yml_catalog"`
	Date    string   `xml:"date,attr"`
	Shop    Shop     `xml:"shop"`
}

type Shop struct {
	Name       string         `xml:"name"`
	Company    string         `xml:"company"`
	URL        string         `xml:"url"`
	Categories []FeedCategory `xml:"categories>category"`
	Offers     []*Offer       `xml:"offers>offer"`
}

type FeedCategory struct {
	ID       string `xml:"id,attr"`
	ParentID string `xml:"parentId,attr,omitempty"`
	Picture  string `xml:"picture,attr,omitempty"`
	Value    string `xml:",chardata"`
}

type Offer struct {
	ID           int64    `xml:"id,attr"`
	Available    bool     `xml:"available,attr"`
	GroupID      string   `xml:"group_id,attr"`
	Name         string   `xml:"name"`
	Similar      string   `xml:"similar"`
	Price        float32  `xml:"price"`
	Barcode      string   `xml:"barcode"`
	URL          string   `xml:"url"`
	VendorCode   string   `xml:"vendorCode"`
	Sort         int      `xml:"sort"`
	Vendor       string   `xml:"vendor"`
	Rating       float32  `xml:"rating"`
	ReviewsCount int      `xml:"reviews_count"`
	Description  string   `xml:"description"`
	OldPrice     float32  `xml:"oldprice"`
	CategoryIDs  []int64  `xml:"category_ids"`
	Pictures     []string `xml:"picture"`
	MarketIDs    []int64  `xml:"market_ids"`
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
