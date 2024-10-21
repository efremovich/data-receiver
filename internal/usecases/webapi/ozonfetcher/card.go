package ozonfetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/efremovich/data-receiver/pkg/metrics"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/httputil"
)

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
	Result []AttibuteMeta `json:"result"`
	Total  int            `json:"total"`
	LastID string         `json:"last_id"`
}

type AttibuteMeta struct {
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
		AttributeID int `json:"attribute_id"`
		ComplexID   int `json:"complex_id"`
		Values      []struct {
			DictionaryValueID int    `json:"dictionary_value_id"`
			Value             string `json:"value"`
		} `json:"values"`
	} `json:"attributes"`
	ComplexAttributes []any  `json:"complex_attributes"`
	ColorImage        string `json:"color_image"`
	LastID            string `json:"last_id"`
}

func (ozon *ozonAPIclientImp) GetCards(ctx context.Context, desc entity.PackageDescription) ([]entity.Card, error) {
	// cardsIDs, err := getCardList(ctx, ozon.baseURL, ozon.clientID, ozon.apiKey, desc.Limit, ozon.timeout, ozon.metric)
	// if err != nil {
	// 	return nil, err
	// }

	// cardsMeta, err := getCardsMeta(ctx, ozon.baseURL, ozon.clientID, ozon.apiKey, desc.Limit, ozon.timeout, ozon.metric, cardsIDs)
	// if err != nil {
	// 	return nil, err
	// }

	// categoryIDsMap := make(map[int]int)

	// for _, card := range cardsMeta {
	// 	categoryIDsMap[card.DescriptionCategoryID] = card.TypeID
	// }

	// attributes, err := getAttributeList(ctx, ozon.baseURL, ozon.clientID, ozon.apiKey, desc.Limit, ozon.timeout, ozon.metric, categoryIDsMap)
	// if err != nil {
	// 	return nil, err
	// }

	attibutesMeta, err := getAttributeMetaList(ctx, ozon.baseURL, ozon.clientID, ozon.apiKey, desc.Limit, ozon.timeout, ozon.metric)
	if err != nil {
		return nil, err
	}
	fmt.Println(len(attibutesMeta))

	var cardsList []entity.Card

	// for _, in := range cardsMeta {
	// 	card := entity.Card{
	// 		ID:              0,
	// 		ExternalID:      int64(in.ID),
	// 		VendorID:        "",
	// 		VendorCode:      "",
	// 		Title:           in.Name,
	// 		Description:     "",
	// 		CreatedAt:       time.Now(),
	// 		Brand:           entity.Brand{},
	// 		Dimension:       entity.Dimension{},
	// 		Characteristics: []*entity.CardCharacteristic{},
	// 		Categories:      []*entity.Category{},
	// 		Sizes:           []*entity.Size{},
	// 		Barcodes:        []*entity.Barcode{},
	// 		MediaFile:       []*entity.MediaFile{},
	// 	}
	// 	cardsList = append(cardsList, card)
	// }

	return cardsList, nil
}

func getAttributeMetaList(ctx context.Context, baseURL, clientID, apiKey string, limit int, timeout time.Duration, metric metrics.Collector) ([]AttibuteMeta, error) {
	items := []AttibuteMeta{}

	methodName := "/v3/products/info/attributes"

	url := fmt.Sprintf("%s%s", baseURL, methodName)

	type Filter struct {
		Filter struct {
			ProductID  []string `json:"product_id"`
			Visibility string   `json:"visibility"`
		} `json:"filter"`
		Limit   int    `json:"limit"`
		LastID  string `json:"last_id"`
		SortDir string `json:"sort_dir"`
	}

	filter := Filter{
		Limit: limit,
	}

	headers := make(map[string]string)
	headers["Client-Id"] = clientID
	headers["Api-Key"] = apiKey
	headers["Content-Type"] = "application/json"

	total := 0
	run := true

	for run {
		bodyData, err := json.Marshal(filter)
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", methodName, err)
		}

		code, resp, err := httputil.SendHTTPRequest(http.MethodPost, url, bodyData, headers, "", "", timeout)
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка выполнения запроса: %w", methodName, err)
		}

		if code != http.StatusOK {
			return nil, fmt.Errorf("%s: ошибка выполнения запроса: %s", methodName, resp)
		}

		var response AttibutesMeta
		if err := json.Unmarshal(resp, &response); err != nil {
			return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
		}

		filter.LastID = response.LastID
		total += len(response.Result)

		items = append(items, response.Result...)

		if total >= response.Total {
			run = false
		}
	}

	return items, nil
}

func getAttributeList(ctx context.Context, baseURL, clientID, apiKey string, limit int, timeout time.Duration, metric metrics.Collector, categoryIDsMap map[int]int) ([]Attribute, error) {
	items := []Attribute{}

	methodName := "/v1/description-category/attribute"

	url := fmt.Sprintf("%s%s", baseURL, methodName)

	type Filter struct {
		CategoryID int    `json:"description_category_id"`
		Language   string `json:"language"`
		TypeID     int    `json:"type_id"`
	}

	filter := Filter{}

	headers := make(map[string]string)
	headers["Client-Id"] = clientID
	headers["Api-Key"] = apiKey
	headers["Content-Type"] = "application/json"

	for key, value := range categoryIDsMap {
		filter.CategoryID = key
		filter.Language = "DEFAULT"
		filter.TypeID = value

		bodyData, err := json.Marshal(filter)
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", methodName, err)
		}

		code, resp, err := httputil.SendHTTPRequest(http.MethodPost, url, bodyData, headers, "", "", timeout)
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка выполнения запроса: %w", methodName, err)
		}

		if code != http.StatusOK {
			return nil, fmt.Errorf("%s: ошибка выполнения запроса: %s", methodName, resp)
		}

		var response Attributes
		if err := json.Unmarshal(resp, &response); err != nil {
			return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
		}

		items = append(items, response.Result...)
	}

	return items, nil
}

func getCardList(ctx context.Context, baseURL, clientID, apiKey string, limit int, timeout time.Duration, metric metrics.Collector) ([]ProductIdList, error) {
	offerIDs := []ProductIdList{}

	methodName := "/v2/product/list"

	url := fmt.Sprintf("%s%s", baseURL, methodName)
	filter := OzonFilter{
		LastID: "",
		Limit:  limit,
	}

	bodyData, err := json.Marshal(filter)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", methodName, err)
	}

	headers := make(map[string]string)
	headers["Client-Id"] = clientID
	headers["Api-Key"] = apiKey
	headers["Content-Type"] = "application/json"

	total := 0
	run := true

	for run {
		code, resp, err := httputil.SendHTTPRequest(http.MethodPost, url, bodyData, headers, "", "", timeout)
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка выполнения запроса: %w", methodName, err)
		}

		if code != http.StatusOK {
			return nil, fmt.Errorf("%s: ошибка выполнения запроса: %s", methodName, resp)
		}

		var productListResponse ProductList
		if err := json.Unmarshal(resp, &productListResponse); err != nil {
			return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
		}

		filter.LastID = productListResponse.Result.LastID
		total += len(productListResponse.Result.Items)

		offerIDs = append(offerIDs, productListResponse.Result.Items...)

		if total >= productListResponse.Result.Total {
			run = false
		}
	}

	return offerIDs, nil
}

func getCardsMeta(ctx context.Context, baseURL, clientID, apiKey string, limit int, timeout time.Duration, metric metrics.Collector, productIDList []ProductIdList) ([]Items, error) {
	items := []Items{}

	methodName := "/v2/product/info/list"

	url := fmt.Sprintf("%s%s", baseURL, methodName)

	type Filter struct {
		OfferID   []string `json:"offer_id"`
		ProductID []string `json:"product_id"`
		Sku       []string `json:"sku"`
	}

	filter := Filter{}

	headers := make(map[string]string)
	headers["Client-Id"] = clientID
	headers["Api-Key"] = apiKey
	headers["Content-Type"] = "application/json"

	offerID := []string{}
	for _, id := range productIDList {
		offerID = append(offerID, id.OfferID)
		filter.OfferID = offerID

		if len(filter.OfferID) == limit { // Лимит в 1000 штук
			bodyData, err := json.Marshal(filter)
			if err != nil {
				return nil, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", methodName, err)
			}

			code, resp, err := httputil.SendHTTPRequest(http.MethodPost, url, bodyData, headers, "", "", timeout)
			if err != nil {
				return nil, fmt.Errorf("%s: ошибка выполнения запроса: %w", methodName, err)
			}

			if code != http.StatusOK {
				return nil, fmt.Errorf("%s: ошибка выполнения запроса: %s", methodName, resp)
			}

			var productListResponse CardResponse
			if err := json.Unmarshal(resp, &productListResponse); err != nil {
				return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
			}

			items = append(items, productListResponse.Result.Items...)

			offerID = []string{}
		}
	}

	return items, nil
}
