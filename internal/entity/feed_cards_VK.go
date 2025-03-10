package entity

type VKCard struct {
	VendorCode  string   `json:"code"`
	Subject     string   `json:"subject"`
	Color       string   `json:"color"`
	Title       string   `json:"title"`
	Gender      string   `json:"gender"`
	Description string   `json:"description"`
	MediaLinks  []string `json:"media_links"`
	Price       float64  `json:"price"`
	MaxSize     string   `json:"size"`
	ExternalID  int64    `json:"external_id"`
	SellerName  string   `json:"seller_name"`
}

type ResponseVKCard struct {
	Cards  []*VKCard `json:"cards"`
	Total  int       `json:"total"`
	Cursor string    `json:"cursor"`
}

type VkCardsFeedParams struct {
	Limit  int    `json:"limit"`
	Cursor string `json:"cursor"`
	Filter string `json:"filter"`
}
