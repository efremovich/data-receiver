package ozonfetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
)

//nolint:funlen // TODO: Разбить метод на более мелкие части
func (ozon *apiClientImp) GetCards(ctx context.Context, _ entity.PackageDescription) ([]entity.Card, error) {
	cardsID, err := ozon.getCardList(ctx)
	if err != nil {
		return nil, err
	}

	cardsMeta, err := ozon.getCardMetaOnProductID(ctx, cardsID)
	if err != nil {
		return nil, err
	}

	categoryIDsMap := make(map[int]int)

	for _, card := range cardsMeta {
		categoryIDsMap[card.TypeID] = card.DescriptionCategoryID
	}

	attributes, err := ozon.getDescriptionCategoryAttr(ctx, categoryIDsMap)
	if err != nil {
		return nil, err
	}

	attibutesMeta, err := ozon.getProductsInfoAttrs(ctx)
	if err != nil {
		return nil, err
	}

	categoriesMap, err := ozon.getCardCategory(ctx)
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
				Link:   media,
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

func (ozon *apiClientImp) getDescriptionCategoryAttr(ctx context.Context, categoryIDsMap map[int]int) ([]Attribute, error) {
	items := []Attribute{}

	url := fmt.Sprintf("%s%s", marketPlaceAPIURL, descriptionCategoryAttrMethod)

	type Filter struct {
		CategoryID int    `json:"description_category_id"`
		Language   string `json:"language"`
		TypeID     int    `json:"type_id"`
	}

	filter := Filter{}

	for key, value := range categoryIDsMap {
		filter.CategoryID = value
		filter.Language = "DEFAULT"
		filter.TypeID = key

		bodyData, err := json.Marshal(filter)
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", descriptionCategoryAttrMethod, err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyData))
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка создания запроса: %w", descriptionCategoryAttrMethod, err)
		}

		for k, v := range ozonHeaders[ozon.marketPlace.ExternalID] {
			req.Header.Set(k, v)
		}

		resp, err := ozon.client.Do(req)

		if err != nil || resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s: ошибка выполнения запроса: %w", descriptionCategoryAttrMethod, err)
		}

		defer resp.Body.Close()

		var response Attributes
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", descriptionCategoryAttrMethod, err)
		}

		items = append(items, response.Result...)
	}

	return items, nil
}

func (ozon *apiClientImp) getProductsInfoAttrs(ctx context.Context) ([]AttributeMeta, error) {
	items := []AttributeMeta{}

	url := fmt.Sprintf("%s%s", marketPlaceAPIURL, productInfoAttrMetod)

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
		Limit: requestItemLimit,
	}

	total := 0
	run := true

	for run {
		bodyData, err := json.Marshal(filter)
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", productInfoAttrMetod, err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyData))
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка создания запроса: %w", productInfoAttrMetod, err)
		}

		for k, v := range ozonHeaders[ozon.marketPlace.ExternalID] {
			req.Header.Set(k, v)
		}

		resp, err := ozon.client.Do(req)

		if err != nil || resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s: ошибка выполнения запроса: %d", productInfoAttrMetod, resp.StatusCode)
		}

		defer resp.Body.Close()

		var response RespAttributeMeta
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", productInfoAttrMetod, err)
		}

		filter.LastID = response.LastID
		total += len(response.Result)

		items = append(items, response.Result...)

		if total == response.Total {
			run = false
		}
	}

	return items, nil
}

func (ozon *apiClientImp) getCardList(ctx context.Context) ([]int, error) {
	offerIDs := []ProductIDList{}

	url := fmt.Sprintf("%s%s", marketPlaceAPIURL, productListMethod)
	filter := OzonFilter{
		LastID: "",
		Limit:  requestItemLimit,
	}

	bodyData, err := json.Marshal(filter)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", productListMethod, err)
	}

	total := 0
	run := true

	productIDs := []int{}

	for run {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyData))
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка создания запроса: %w", productListMethod, err)
		}

		for k, v := range ozonHeaders[ozon.marketPlace.ExternalID] {
			req.Header.Set(k, v)
		}

		resp, err := ozon.client.Do(req)

		if err != nil || resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s: ошибка выполнения запроса: %d", productListMethod, resp.StatusCode)
		}

		defer resp.Body.Close()

		var response ProductList
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", productListMethod, err)
		}

		filter.LastID = response.Result.LastID
		total += len(response.Result.Items)

		offerIDs = append(offerIDs, response.Result.Items...)

		if total >= response.Result.Total {
			run = false
		}
	}

	for _, elem := range offerIDs {
		productIDs = append(productIDs, elem.ProductID)
	}

	return productIDs, nil
}

func (ozon *apiClientImp) getCardMetaOnProductID(ctx context.Context, productIDList []int) ([]Items, error) {
	items := []Items{}

	url := fmt.Sprintf("%s%s", marketPlaceAPIURL, productInfoListMethod)

	type Filter struct {
		OfferID   []string `json:"offer_id"`
		ProductID []int    `json:"product_id"`
		Sku       []int    `json:"sku"`
	}

	filter := Filter{}

	chunks := chunkIntSlice(productIDList, requestItemLimit)

	for _, chunk := range chunks {
		filter.ProductID = chunk

		bodyData, err := json.Marshal(filter)
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", productInfoListMethod, err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyData))
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка создания запроса: %w", productInfoListMethod, err)
		}

		for k, v := range ozonHeaders[ozon.marketPlace.ExternalID] {
			req.Header.Set(k, v)
		}

		resp, err := ozon.client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s: ошибка выполнения запроса: %d", productInfoListMethod, resp.StatusCode)
		}
		defer resp.Body.Close()

		var response CardResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", productInfoListMethod, err)
		}

		items = append(items, response.Result.Items...)
	}

	return items, nil
}

func (ozon *apiClientImp) getProductInfoOnSKU(ctx context.Context, skus []int) (map[int]ItemsResponse, error) {
	url := fmt.Sprintf("%s%s", marketPlaceAPIURL, productInfoListMethod)

	type f struct {
		Sku []int `json:"sku"`
	}

	filter := f{
		Sku: skus,
	}

	bodyData, err := json.Marshal(filter)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", productInfoListMethod, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyData))
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка создания запроса: %w", productInfoListMethod, err)
	}

	for k, v := range ozonHeaders[ozon.marketPlace.ExternalID] {
		req.Header.Set(k, v)
	}
	resp, err := ozon.client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: ошибка выполнения запроса: %d", productInfoListMethod, resp.StatusCode)
	}

	defer resp.Body.Close()

	var response ProductInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", productInfoListMethod, err)
	}

	result := make(map[int]ItemsResponse, len(response.Items))

	for _, elem := range response.Items {
		for _, sku := range elem.Sources {
			result[sku.Sku] = elem
		}
	}

	return result, nil
}

func (ozon *apiClientImp) getCardCategory(ctx context.Context) (map[int]Category, error) {
	url := fmt.Sprintf("%s%s", marketPlaceAPIURL, descriptionCategoryMethod)

	type Filter struct {
		Language string `json:"language"`
	}

	filter := Filter{
		Language: "DEFAULT",
	}

	bodyData, err := json.Marshal(filter)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", descriptionCategoryMethod, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyData))
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка создания запроса: %w", descriptionCategoryMethod, err)
	}

	for k, v := range ozonHeaders[ozon.marketPlace.ExternalID] {
		req.Header.Set(k, v)
	}
	resp, err := ozon.client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: ошибка выполнения запроса: %d", descriptionCategoryMethod, resp.StatusCode)
	}
	defer resp.Body.Close()

	var response Categories
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", descriptionCategoryMethod, err)
	}

	categories := make(map[int]Category)
	for _, topCat := range response.TopCategories {
		flattenCategories(topCat, 0, categories)
	}

	return categories, nil
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
