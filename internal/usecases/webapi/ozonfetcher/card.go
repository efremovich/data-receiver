package ozonfetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/efremovich/data-receiver/pkg/metrics"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/httputil"
)

func (ozon *ozonAPIclientImp) GetCards(ctx context.Context, desc entity.PackageDescription) ([]entity.Card, error) {
	cardsIDs, err := getCardList(ctx, ozon.baseURL, ozon.clientID, ozon.apiKey, desc.Limit, ozon.timeout, ozon.metric)
	if err != nil {
		return nil, err
	}

	cardsMeta, err := getCardsMeta(ctx, ozon.baseURL, ozon.clientID, ozon.apiKey, desc.Limit, ozon.timeout, ozon.metric, cardsIDs)
	if err != nil {
		return nil, err
	}

	categoryIDsMap := make(map[int]int)

	for _, card := range cardsMeta {
		categoryIDsMap[card.TypeID] = card.DescriptionCategoryID
	}

	categoriesMap, err := getCategory(ctx, ozon.baseURL, ozon.clientID, ozon.apiKey, desc.Limit, ozon.timeout, ozon.metric)
	if err != nil {
		return nil, err
	}

	attributes, err := getAttributeList(ctx, ozon.baseURL, ozon.clientID, ozon.apiKey, desc.Limit, ozon.timeout, ozon.metric, categoryIDsMap)
	if err != nil {
		return nil, err
	}

	attibutesMeta, err := getAttributeMetaList(ctx, ozon.baseURL, ozon.clientID, ozon.apiKey, desc.Limit, ozon.timeout, ozon.metric)
	if err != nil {
		return nil, err
	}

	attr := make(map[int]Attribute)
	for _, attribute := range attributes {
		attr[attribute.ID] = attribute
	}

	var cardsList []entity.Card

	const brandID = 31

	for _, in := range attibutesMeta {
		characteristics := []*entity.CardCharacteristic{}
		categories := []*entity.Category{}
		dimension := entity.Dimension{}
		sizes := []*entity.Size{}
		barcodes := []*entity.Barcode{}
		mediaFile := []*entity.MediaFile{}
		brand := entity.Brand{}
		description := ""

		// Char и Brand
		for _, char := range in.Attributes {
			switch {
			case char.AttributeID == brandID:
				// Brand
				brand.ExternalID = int64(char.AttributeID)
				for _, charVal := range char.Values {
					brand.Title = charVal.Value
				}
			case attr[char.AttributeID].Name == "Аннотация":
				// Description
				for _, charVal := range char.Values {
					description = charVal.Value
				}
			default:
				// Characteristics
				charValues := []string{}
				for _, charVal := range char.Values {
					charValues = append(charValues, charVal.Value)
				}

				characteristic := entity.CardCharacteristic{
					Title: attr[char.AttributeID].Name,
					Value: charValues,
				}
				characteristics = append(characteristics, &characteristic)
			}
		}

		if cat, ok := categoriesMap[in.DescriptionCategoryID]; ok {
			categories = append(categories, &entity.Category{
				Title:      cat.CategoryName,
				ExternalID: int64(in.DescriptionCategoryID),
			})
		}

		dimension.Width = in.Width
		dimension.Height = in.Height
		dimension.Length = in.Depth

		vendorCode, vendorID, currSize := getMetaFromVendorID(in.OfferID)
		sizes = append(sizes, &entity.Size{
			TechSize: currSize,
			Title:    currSize,
		})

		barcodes = append(barcodes, &entity.Barcode{
			Barcode: in.Barcode,
		})

		for _, media := range in.Images {
			mediaFile = append(mediaFile, &entity.MediaFile{
				Link:   media.FileName,
				TypeID: 1,
			})
		}

		card := entity.Card{
			ExternalID:      int64(in.ID),
			VendorID:        vendorID,
			VendorCode:      vendorCode,
			Title:           in.Name,
			Description:     description,
			CreatedAt:       time.Now(),
			Brand:           brand,
			Dimension:       dimension,
			Characteristics: characteristics,
			Categories:      categories,
			Sizes:           sizes,
			Barcodes:        barcodes,
			MediaFile:       mediaFile,
		}
		cardsList = append(cardsList, card)
	}

	return cardsList, nil
}

func getMetaFromVendorID(offerID string) (string, string, string) {
	var vendorID, vendorCode, vendorSize string
	// Артикул, код, размер "RBB-061/00-0014881/58"
	vendorData := strings.Split(offerID, "/")
	if len(vendorData) == 2 {
		vendorID = vendorData[0]
		vendorSize = vendorData[1]
	}

	if len(vendorData) == 3 {
		vendorCode = vendorData[0]
		vendorID = vendorData[1]
		vendorSize = vendorData[2]
	}

	if len(vendorData) == 1 {
		vendorCode = vendorData[0]
		vendorID = vendorData[0]
		vendorSize = "One size"
	}

	return vendorID, vendorCode, vendorSize
}

func getAttributeMetaList(_ context.Context, baseURL, clientID, apiKey string, limit int, timeout time.Duration, metric metrics.Collector) ([]AttributeMeta, error) {
	items := []AttributeMeta{}

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

	// TODO подумать на счет лимита. Озон позволяет брать по 1000
	limit *= 10
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

		if limit != response.Total {
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
		filter.CategoryID = value
		filter.Language = "DEFAULT"
		filter.TypeID = key

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

func getCardList(ctx context.Context, baseURL, clientID, apiKey string, limit int, timeout time.Duration, metric metrics.Collector) ([]int, error) {
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

	productIDs := []int{}
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
	for _, elem := range offerIDs {
		productIDs = append(productIDs, elem.ProductID)
	}
	return productIDs, nil
}

func getCardsMeta(ctx context.Context, baseURL, clientID, apiKey string, limit int, timeout time.Duration, metric metrics.Collector, productIDList []int) ([]Items, error) {
	items := []Items{}

	methodName := "/v2/product/info/list"

	url := fmt.Sprintf("%s%s", baseURL, methodName)

	type Filter struct {
		OfferID   []string `json:"offer_id"`
		ProductID []int    `json:"product_id"`
		Sku       []string `json:"sku"`
	}

	filter := Filter{}

	headers := make(map[string]string)
	headers["Client-Id"] = clientID
	headers["Api-Key"] = apiKey
	headers["Content-Type"] = "application/json"

	chunkSize := 50
	chunks := chunkIntSlice(productIDList, chunkSize)

	for _, chunk := range chunks {
		filter.ProductID = chunk
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
	}

	return items, nil
}

func getCategory(ctx context.Context, baseURL, clientID, apiKey string, limit int, timeout time.Duration, metric metrics.Collector) (map[int]Category, error) {
	methodName := "/v1/description-category/tree"

	url := fmt.Sprintf("%s%s", baseURL, methodName)

	type Filter struct {
		Language string `json:"language"`
	}

	filter := Filter{
		Language: "DEFAULT",
	}

	headers := make(map[string]string)
	headers["Client-Id"] = clientID
	headers["Api-Key"] = apiKey
	headers["Content-Type"] = "application/json"

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

	var response Categories
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
	}

	categories := make(map[int]Category)
	for _, topCat := range response.TopCategories {
		flattenCategories(topCat, 0, categories)
	}

	return categories, nil
}

func flattenCategories(cat Category, parentID int, categories map[int]Category) {
	category := Category{
		TypeName:              cat.TypeName,
		TypeID:                cat.TypeID,
		DescriptionCategoryID: cat.DescriptionCategoryID,
		CategoryName:          cat.CategoryName,
		Disabled:              cat.Disabled,
		Children:              []Category{},
		ParentID:              parentID,
	}

	categories[cat.DescriptionCategoryID] = category

	for _, child := range cat.Children {
		flattenCategories(child, cat.DescriptionCategoryID, categories)
	}
}
