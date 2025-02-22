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
	"github.com/efremovich/data-receiver/pkg/logger"
	"github.com/efremovich/data-receiver/pkg/metrics"
)

type SupplyOrderList struct {
	SupplyOrderID     []int `json:"supply_order_id"`
	LastSupplyOrderID int   `json:"last_supply_order_id"`
}

func (ozon *apiClientImp) GetStocks(ctx context.Context, desc entity.PackageDescription) ([]entity.StockMeta, error) {
	supplyList, err := getSupplyList(ctx, marketPlaceAPIURL, ozon.clientID, ozon.apiKey, ozon.metric)
	if err != nil {
		return nil, err
	}

	supplyData, err := getSupplyDataList(ctx, marketPlaceAPIURL, ozon.clientID, ozon.apiKey, ozon.metric, supplyList)
	if err != nil {
		return nil, err
	}

	stockResponce, err := getSupplyBundle(ctx, marketPlaceAPIURL, ozon.clientID, ozon.apiKey, ozon.metric, supplyData)
	if err != nil {
		return nil, err
	}

	var stockMetaList []entity.StockMeta

	for _, elem := range stockResponce.StocksMeta {
		stockMeta := entity.StockMeta{}

		stockDate, err := time.Parse("02.01.2006", elem.CreationDate)
		if err != nil {
			stockDate = time.Now()

			logger.GetLoggerFromContext(ctx).Warnf("Не удалось распознать дату остатка %s", elem.CreationDate)
		}

		stockMeta.Stock = entity.Stock{
			Quantity:  elem.Quantity,
			CreatedAt: stockDate,
		}

		stockMeta.Seller2Card = entity.Seller2Card{
			ExternalID: int64(elem.ProductID),
		}

		stockMeta.Barcode = entity.Barcode{
			ExternalID: int64(elem.ProductID),
			Barcode:    elem.Barcode,
		}

		stockMeta.Warehouse = entity.Warehouse{
			Title:      elem.Warehouse.Name,
			Address:    elem.Warehouse.Address,
			ExternalID: elem.Warehouse.WarehouseID,
			TypeName:   "FBO",
		}

		vendorID, vendorCode, currSize := getMetaFromVendorID(elem.OfferID)
		stockMeta.Size = entity.Size{
			TechSize: currSize,
			Title:    currSize,
		}
		stockMeta.Card = entity.Card{
			VendorCode: vendorCode,
			VendorID:   vendorID,
		}

		for _, cardMeta := range stockResponce.CardMeta {
			if cardMeta.ID == elem.ProductID {
				price, err := strconv.ParseFloat(cardMeta.Price, 32)
				if err != nil {
					logger.GetLoggerFromContext(ctx).Warnf("не удалось преобразовать в число %s", &cardMeta.Price)
				}

				oldPrice, err := strconv.ParseFloat(cardMeta.OldPrice, 32)
				if err != nil {
					logger.GetLoggerFromContext(ctx).Warnf("не удалось преобразовать в число %s", &cardMeta.Price)
				}

				marketingPrice, err := strconv.ParseFloat(cardMeta.MarketingPrice, 32)
				if err != nil {
					logger.GetLoggerFromContext(ctx).Warnf("не удалось преобразовать в число %s", &cardMeta.Price)
				}

				stockMeta.PriceSize = entity.PriceSize{
					Price:        float32(price),
					Discount:     float32(oldPrice - price),
					SpecialPrice: float32(marketingPrice),
				}
			}
		}

		stockMetaList = append(stockMetaList, stockMeta)
	}

	return stockMetaList, nil
}

func getSupplyBundle(ctx context.Context, baseURL, clientID, apiKey string, metric metrics.Collector, supplyData *SupplyData) (*StocksMeta, error) {
	timeout := time.Second * time.Duration(30)

	methodName := "/v1/supply-order/bundle"

	url := fmt.Sprintf("%s%s", baseURL, methodName)

	type BundleFiter struct {
		BundleIDs []string `json:"bundle_ids"`
		IsAsc     bool     `json:"is_asc"`
		Limit     int      `json:"limit"`
		LastID    string   `json:"last_id"`
	}

	filter := BundleFiter{
		BundleIDs: []string{},
		IsAsc:     false,
		Limit:     100,
	}

	headers := make(map[string]string)
	headers["Client-Id"] = clientID
	headers["Api-Key"] = apiKey
	headers["Content-Type"] = "application/json"

	stockMeta := []BundleItems{}
	productIDList := []int{}
	stocksMeta := StocksMeta{}

	for _, order := range supplyData.Orders {
		bundles := []string{}
		for _, bandleID := range order.Supplies {
			bundles = append(bundles, bandleID.BundleID)
		}

		filter.BundleIDs = bundles

		endOfList := false
		for !endOfList {
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

			var response SupplyBundleData
			if err := json.Unmarshal(resp, &response); err != nil {
				return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
			}

			if response.HasNext {
				filter.LastID = response.LastID
			} else {
				filter.LastID = ""
			}

			endOfList = !response.HasNext

			for _, item := range response.Items {
				for _, warehouse := range supplyData.Warehouses {
					if warehouse.WarehouseID == order.DropoffWarehouseID {
						item.Warehouse = warehouse
					}
				}

				item.CreationDate = order.CreationDate
				stockMeta = append(stockMeta, item)
				productIDList = append(productIDList, item.ProductID)
			}

			time.Sleep(3 * time.Second)
		}
	}

	cardsMeta, err := getCardsMeta(ctx, baseURL, clientID, apiKey, 1000, timeout, metric, productIDList)
	if err != nil {
		return nil, err
	}

	stocksMeta.StocksMeta = stockMeta
	stocksMeta.CardMeta = cardsMeta

	return &stocksMeta, nil
}

func getSupplyDataList(ctx context.Context, baseURL, clientID, apiKey string, metric metrics.Collector, supplyList []int) (*SupplyData, error) {
	timeout := time.Second * time.Duration(30)

	methodName := "/v2/supply-order/get"

	url := fmt.Sprintf("%s%s", baseURL, methodName)

	type SupplyFilterList struct {
		OrderIDs []int `json:"order_ids"`
	}

	headers := make(map[string]string)
	headers["Client-Id"] = clientID
	headers["Api-Key"] = apiKey
	headers["Content-Type"] = "application/json"

	chunkSize := 50
	// TODO Срез из 100 для теста. В бою удалть.
	chunks := chunkIntSlice(supplyList, chunkSize)

	var response SupplyData

	for _, chunk := range chunks {
		filter := SupplyFilterList{
			OrderIDs: chunk,
		}

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

		var sd SupplyData
		if err := json.Unmarshal(resp, &sd); err != nil {
			return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
		}

		response.Orders = append(response.Orders, sd.Orders...)
		response.Warehouses = append(response.Warehouses, sd.Warehouses...)
	}

	return &response, nil
}

func chunkIntSlice(slice []int, chunkSize int) [][]int {
	var chunks [][]int

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func getSupplyList(ctx context.Context, baseURL, clientID, apiKey string, metric metrics.Collector) ([]int, error) {
	timeout := time.Second * time.Duration(30)

	methodName := "/v2/supply-order/list"

	url := fmt.Sprintf("%s%s", baseURL, methodName)

	type Filter struct {
		States []string `json:"states"`
	}

	type Paging struct {
		FromSupplyOrderID int `json:"from_supply_order_id"`
		Limit             int `json:"limit"`
	}

	type SupplyFilter struct {
		Filter Filter `json:"filter"`
		Paging Paging `json:"paging"`
	}

	filter := SupplyFilter{}

	filter.Filter.States = []string{"ORDER_STATE_COMPLETED"}
	filter.Paging.Limit = 100

	headers := make(map[string]string)
	headers["Client-Id"] = clientID
	headers["Api-Key"] = apiKey
	headers["Content-Type"] = "application/json"

	supplyOrderID := []int{}
	endOfList := false

	for !endOfList {
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

		var response SupplyOrderList
		if err := json.Unmarshal(resp, &response); err != nil {
			return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
		}

		supplyOrderID = append(supplyOrderID, response.SupplyOrderID...)
		filter.Paging.FromSupplyOrderID = response.LastSupplyOrderID

		if len(response.SupplyOrderID) < filter.Paging.Limit {
			endOfList = true
		}
	}

	return supplyOrderID, nil
}
