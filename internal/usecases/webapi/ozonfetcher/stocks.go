package ozonfetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/httputil"
	"github.com/efremovich/data-receiver/pkg/metrics"
)

type SupplyOrderList struct {
	SupplyOrderID     []int `json:"supply_order_id"`
	LastSupplyOrderID int   `json:"last_supply_order_id"`
}

func (ozon *ozonAPIclientImp) GetStocks(ctx context.Context, desc entity.PackageDescription) ([]entity.StockMeta, error) {
	supplyList, err := getSupplyList(ctx, ozon.baseURL, ozon.clientID, ozon.apiKey, ozon.metric)
	if err != nil {
		return nil, err
	}

	supplyData, err := getSupplyDataList(ctx, ozon.baseURL, ozon.clientID, ozon.apiKey, ozon.metric, supplyList)
	if err != nil {
		return nil, err
	}

	_, err = getSupplyBundle(ctx, ozon.baseURL, ozon.clientID, ozon.apiKey, ozon.metric, supplyData)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func getSupplyBundle(ctx context.Context, baseURL, clientID, apiKey string, metric metrics.Collector, supplyData *SupplyData) ([]BundleItems, error) {
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

	items := []BundleItems{}

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
			}

			endOfList = !response.HasNext

			for _, item := range response.Items {
				for _, warehouse := range supplyData.Warehouses {
					if warehouse.WarehouseID == order.DropoffWarehouseID {
						item.Warehouse = warehouse
					}
				}

				item.CreationDate = order.CreationDate
				items = append(items, item)
			}

			time.Sleep(3 * time.Second)
		}
	}

	return items, nil
}

func getSupplyDataList(ctx context.Context, baseURL, clientID, apiKey string, metric metrics.Collector, supplyList []int) (*SupplyData, error) {
	timeout := time.Second * time.Duration(30)

	methodName := "/v2/supply-order/get"

	url := fmt.Sprintf("%s%s", baseURL, methodName)

	type SupplyFilterList struct {
		OrderIDs []int `json:"order_ids"`
	}

	filter := SupplyFilterList{}

	for i := 0; i < 50; i++ {
		filter.OrderIDs = append(filter.OrderIDs, supplyList[i])
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

	var response SupplyData
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
	}

	return &response, nil
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
