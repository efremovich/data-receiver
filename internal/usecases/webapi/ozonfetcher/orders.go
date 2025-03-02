package ozonfetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/httputil"
	"github.com/efremovich/data-receiver/pkg/metrics"
)

func (o *apiClientImp) GetOrders(ctx context.Context, desc entity.PackageDescription) ([]entity.Order, error) {
	const methodName = "/v2/posting/fbo/list"

	timeout := time.Second * time.Duration(30)
	url := fmt.Sprintf("%s%s", marketPlaceAPIURL, methodName)

	filter := OrderFilter{}
	filter.Dir = "desc"
	filter.Limit = 1000
	filter.With.AnalyticsData = true
	filter.With.FinancialData = true

	startDate := desc.UpdatedAt.Truncate(24 * time.Hour)
	filter.Filter.Since = startDate
	filter.Filter.To = startDate.Add(24 * time.Hour)

	headers := make(map[string]string)
	headers["Client-Id"] = o.clientID
	headers["Api-Key"] = o.apiKey
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

	var response OrderRespose
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
	}

	skus := []int{}

	for _, elem := range response.Result {
		for _, product := range elem.Products {
			skus = append(skus, product.Sku)
		}
	}
	productInfo, err := getProductInfo(ctx, marketPlaceAPIURL, o.clientID, o.apiKey, o.metric, skus)
	if err != nil {
		return nil, fmt.Errorf("ошибка получение подробной информации о товаре %w", err)
	}

	var orders []entity.Order
	for _, elem := range response.Result {
		warehouse := entity.Warehouse{}
		warehouse.Title = elem.AnalyticsData.WarehouseName
		warehouse.ExternalID = elem.AnalyticsData.WarehouseID

		status := entity.Status{
			Name: elem.Status,
		}

		seller := entity.MarketPlace{
			Title:      "ozon",
			ExternalID: o.clientID,
		}

		region := entity.Region{
			RegionName: "Неопределенно",
			District: entity.District{
				Name: "Неопределенно",
			},
			Country: entity.Country{
				Name: "Россия",
			},
		}

		for _, product := range elem.Products {
			barcode := entity.Barcode{}
			vendorID, vendorCode, currSize := getMetaFromVendorID(product.OfferID)
			card := entity.Card{
				VendorCode: vendorCode,
				VendorID:   vendorID,
				ExternalID: int64(productInfo[product.Sku].ID),
			}

			size := entity.Size{
				TechSize: currSize,
				Title:    currSize,
			}

			barcodes := productInfo[product.Sku].Barcodes
			for _, b := range barcodes {
				barcode.Barcode = b
			}

			priceSize := entity.PriceSize{}
			for _, fData := range elem.FinancialData.Products {
				if fData.ProductID == product.Sku {
					priceSize = entity.PriceSize{
						Price:        fData.Price,
						Discount:     fData.TotalDiscountValue,
						SpecialPrice: fData.OldPrice,
					}
				}
			}

			order := entity.Order{}
			order.ExternalID = strconv.Itoa(elem.OrderID)
			order.Price = priceSize.Price
			order.Type = elem.Status
			order.Quantity = product.Quantity

			order.CreatedAt = elem.CreatedAt

			order.Size = &size
			order.PriceSize = &priceSize
			order.Status = &status
			order.Warehouse = &warehouse
			order.Seller = &seller
			order.Card = &card
			order.Barcode = &barcode
			order.Region = &region

			orders = append(orders, order)
		}

	}
	return orders, nil
}

func getProductInfo(ctx context.Context, baseURL, clientID, apiKey string, metric metrics.Collector, skus []int) (map[int]ItemsResponse, error) {
	const methodName = "/v3/product/info/list"

	timeout := time.Second * time.Duration(30)
	url := fmt.Sprintf("%s%s", baseURL, methodName)

	type f struct {
		Sku []int `json:"sku"`
	}

	filter := f{
		Sku: skus,
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

	var response ProductInfoResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
	}

	result := make(map[int]ItemsResponse, len(response.Items))
	for _, elem := range response.Items {
		for _, sku := range elem.Sources {
			result[sku.Sku] = elem
		}

		if len(elem.Sources) > 1 {
			fmt.Println("ZZZ Не должно быть более 1")
		}
	}
	return result, nil
}
