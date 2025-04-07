package odincfetcer

type Card struct {
	VendorID    string `json:"vendor_id"`
	VendorCode  string `json:"vendor_code"`
	Title       string `json:"title"`
	ExternalID  int64  `json:"external_id"`
	Description string `json:"description"`
	Category    struct {
		ID    string `json:"ID"`
		Title string `json:"title"`
	} `json:"category"`
	Length  float32 `json:"length"`
	Width   float32 `json:"width"`
	Height  float32 `json:"height"`
	Barcode string  `json:"barcode"`
	Brand   string  `json:"brand"`
	Size    string  `json:"size"`
}

type Cost struct {
	VendorCode string `json:"vendor_code"`
	Title      string `json:"title"`
	SizeCode   string `json:"size_code"`
	SizeTitle  string `json:"size_title"`
	CostPrice  string `json:"cost_price"`
}
