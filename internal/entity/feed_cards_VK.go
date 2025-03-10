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
