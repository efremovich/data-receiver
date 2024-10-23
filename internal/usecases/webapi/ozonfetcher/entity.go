package ozonfetcher

import "time"

type OzonFilter struct {
	Filter struct {
		OfferID    []string `json:"offer_id"`
		ProductID  []string `json:"product_id"`
		Visibility string   `json:"visibility"`
	} `json:"filter"`
	LastID string `json:"last_id"`
	Limit  int    `json:"limit"`
}

type ProductIdList struct {
	ProductID int    `json:"product_id"`
	OfferID   string `json:"offer_id"`
}
type ProductList struct {
	Result struct {
		Items  []ProductIdList `json:"items"`
		Total  int             `json:"total"`
		LastID string          `json:"last_id"`
	} `json:"result"`
}

type CardResponse struct {
	Result Result `json:"result"`
}
type Sources struct {
	IsEnabled bool   `json:"is_enabled"`
	Sku       int    `json:"sku"`
	Source    string `json:"source"`
}
type DiscountedStocks struct {
	Coming   int `json:"coming"`
	Present  int `json:"present"`
	Reserved int `json:"reserved"`
}
type Stocks struct {
	Coming   int `json:"coming"`
	Present  int `json:"present"`
	Reserved int `json:"reserved"`
}
type VisibilityDetails struct {
	HasPrice      bool `json:"has_price"`
	HasStock      bool `json:"has_stock"`
	ActiveProduct bool `json:"active_product"`
}
type ExternalIndexData struct {
	MinimalPrice         string  `json:"minimal_price"`
	MinimalPriceCurrency string  `json:"minimal_price_currency"`
	PriceIndexValue      float64 `json:"price_index_value"`
}
type OzonIndexData struct {
	MinimalPrice         string  `json:"minimal_price"`
	MinimalPriceCurrency string  `json:"minimal_price_currency"`
	PriceIndexValue      float64 `json:"price_index_value"`
}
type SelfMarketplacesIndexData struct {
	MinimalPrice         string  `json:"minimal_price"`
	MinimalPriceCurrency string  `json:"minimal_price_currency"`
	PriceIndexValue      float64 `json:"price_index_value"`
}
type PriceIndexes struct {
	ExternalIndexData         ExternalIndexData         `json:"external_index_data"`
	OzonIndexData             OzonIndexData             `json:"ozon_index_data"`
	PriceIndex                string                    `json:"price_index"`
	SelfMarketplacesIndexData SelfMarketplacesIndexData `json:"self_marketplaces_index_data"`
}
type Status struct {
	State            string    `json:"state"`
	StateFailed      string    `json:"state_failed"`
	ModerateStatus   string    `json:"moderate_status"`
	DeclineReasons   []string  `json:"decline_reasons"`
	ValidationState  string    `json:"validation_state"`
	StateName        string    `json:"state_name"`
	StateDescription string    `json:"state_description"`
	IsFailed         bool      `json:"is_failed"`
	IsCreated        bool      `json:"is_created"`
	StateTooltip     string    `json:"state_tooltip"`
	ItemErrors       []string  `json:"item_errors"`
	StateUpdatedAt   time.Time `json:"state_updated_at"`
}

type Items struct {
	ID                    int               `json:"id"`
	Name                  string            `json:"name"`
	OfferID               string            `json:"offer_id"`
	IsArchived            bool              `json:"is_archived,omitempty"`
	IsAutoarchived        bool              `json:"is_autoarchived,omitempty"`
	Barcode               string            `json:"barcode"`
	Barcodes              []string          `json:"barcodes,omitempty"`
	BuyboxPrice           string            `json:"buybox_price"`
	DescriptionCategoryID int               `json:"description_category_id,omitempty"`
	TypeID                int               `json:"type_id,omitempty"`
	CreatedAt             time.Time         `json:"created_at"`
	Images                []string          `json:"images"`
	CurrencyCode          string            `json:"currency_code,omitempty"`
	MarketingPrice        string            `json:"marketing_price"`
	MinPrice              string            `json:"min_price"`
	OldPrice              string            `json:"old_price"`
	Price                 string            `json:"price"`
	RecommendedPrice      string            `json:"recommended_price"`
	Sources               []Sources         `json:"sources"`
	HasDiscountedItem     bool              `json:"has_discounted_item,omitempty"`
	IsDiscounted          bool              `json:"is_discounted,omitempty"`
	DiscountedStocks      DiscountedStocks  `json:"discounted_stocks,omitempty"`
	State                 string            `json:"state"`
	Stocks                Stocks            `json:"stocks"`
	Errors                []string          `json:"errors"`
	UpdatedAt             time.Time         `json:"updated_at"`
	Vat                   string            `json:"vat"`
	Visible               bool              `json:"visible"`
	VisibilityDetails     VisibilityDetails `json:"visibility_details"`
	PriceIndexes          PriceIndexes      `json:"price_indexes,omitempty"`
	Images360             []string          `json:"images360"`
	IsKgt                 bool              `json:"is_kgt"`
	ColorImage            string            `json:"color_image"`
	PrimaryImage          string            `json:"primary_image"`
	Status                Status            `json:"status"`
	CategoryID            int               `json:"category_id,omitempty"`
	PriceIndex            string            `json:"price_index,omitempty"`
}
type Result struct {
	Items []Items `json:"items"`
}
type Attributes struct {
	Result []Attribute `json:"result"`
}
type Attribute struct {
	ID                 int    `json:"id"`
	AttributeComplexID int    `json:"attribute_complex_id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Type               string `json:"type"`
	IsCollection       bool   `json:"is_collection"`
	IsRequired         bool   `json:"is_required"`
	IsAspect           bool   `json:"is_aspect"`
	MaxValueCount      int    `json:"max_value_count"`
	GroupName          string `json:"group_name"`
	GroupID            int    `json:"group_id"`
	DictionaryID       int    `json:"dictionary_id"`
	CategoryDependent  bool   `json:"category_dependent"`
}

type CatMeta struct {
	CategoryID string
	TypeID     string
}
type AttibutesMeta struct {
	Result []AttributeMeta `json:"result"`
	Total  int             `json:"total"`
	LastID string          `json:"last_id"`
}

type AttributeMeta struct {
	ID            int    `json:"id"`
	Barcode       string `json:"barcode"`
	CategoryID    int    `json:"category_id"`
	Name          string `json:"name"`
	OfferID       string `json:"offer_id"`
	Height        int    `json:"height"`
	Depth         int    `json:"depth"`
	Width         int    `json:"width"`
	DimensionUnit string `json:"dimension_unit"`
	Weight        int    `json:"weight"`
	WeightUnit    string `json:"weight_unit"`
	Images        []struct {
		FileName string `json:"file_name"`
		Default  bool   `json:"default"`
		Index    int    `json:"index"`
	} `json:"images"`
	Images360  []any `json:"images360"`
	PdfList    []any `json:"pdf_list"`
	Attributes []struct {
		AttrMeta    Attribute `json:"attr_meta"`
		AttributeID int       `json:"attribute_id"`
		ComplexID   int       `json:"complex_id"`
		Values      []struct {
			DictionaryValueID int    `json:"dictionary_value_id"`
			Value             string `json:"value"`
		} `json:"values"`
	} `json:"attributes"`
	ComplexAttributes []any  `json:"complex_attributes"`
	ColorImage        string `json:"color_image"`
	LastID            string `json:"last_id"`
}

type Categories struct {
	TopCategories []Category `json:"result"`
}
type Category struct {
	TypeName              string     `json:"type_name"`
	TypeID                int        `json:"type_id"`
	DescriptionCategoryID int        `json:"description_category_id"`
	CategoryName          string     `json:"category_name"`
	Disabled              bool       `json:"disabled"`
	Children              []Category `json:"children"`
	ParentID              int        `json:"parent_id"`
}
