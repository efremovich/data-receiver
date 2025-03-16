package ozonfetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/logger"
)

type SupplyOrderList struct {
	SupplyOrderID     []int `json:"supply_order_id"`
	LastSupplyOrderID int   `json:"last_supply_order_id"`
}

func (ozon *apiClientImp) GetStocks(ctx context.Context, _ entity.PackageDescription) ([]entity.StockMeta, error) {
	supplyList, err := ozon.getSupplyList(ctx)
	if err != nil {
		return nil, err
	}

	supplyData, err := ozon.getSupplyDataList(ctx, supplyList)
	if err != nil {
		return nil, err
	}

	stockResponce, err := ozon.getSupplyBundle(ctx, supplyData)
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
					logger.GetLoggerFromContext(ctx).Warnf("не удалось преобразовать в число %s", cardMeta.Price)
				}

				oldPrice, err := strconv.ParseFloat(cardMeta.OldPrice, 32)
				if err != nil {
					logger.GetLoggerFromContext(ctx).Warnf("не удалось преобразовать в число %s", cardMeta.Price)
				}

				marketingPrice, err := strconv.ParseFloat(cardMeta.MarketingPrice, 32)
				if err != nil {
					logger.GetLoggerFromContext(ctx).Warnf("не удалось преобразовать в число %s", cardMeta.Price)
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

func (ozon *apiClientImp) getSupplyBundle(ctx context.Context, supplyData *SupplyData) (*StocksMeta, error) {
	url := fmt.Sprintf("%s%s", marketPlaceAPIURL, suppolyOrderBundleMethod)

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
				return nil, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", suppolyOrderBundleMethod, err)
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyData))
			if err != nil {
				return nil, fmt.Errorf("%s: ошибка создания запроса: %w", supplyOrderListMethod, err)
			}

			for k, v := range ozonHeaders[ozon.marketPlace.ExternalID] {
				req.Header.Set(k, v)
			}
			resp, err := ozon.client.Do(req)

			if err != nil || resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("%s: ошибка выполнения запроса: %d", suppolyOrderBundleMethod, resp.StatusCode)
			}
			defer resp.Body.Close()

			var response SupplyBundleData
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", suppolyOrderBundleMethod, err)
			}

			if response.HasNext {
				filter.LastID = response.LastID
			} else {
				filter.LastID = ""
			}

			endOfList = !response.HasNext

			stockMeta, productIDList = collectStocks(response, supplyData, order, stockMeta, productIDList)
		}
	}

	cardsMeta, err := ozon.getCardMetaOnProductID(ctx, productIDList)
	if err != nil {
		return nil, err
	}

	stocksMeta.StocksMeta = stockMeta
	stocksMeta.CardMeta = cardsMeta

	return &stocksMeta, nil
}

func collectStocks(response SupplyBundleData, supplyData *SupplyData, order Orders, stockMeta []BundleItems, productIDList []int) ([]BundleItems, []int) {
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
	return stockMeta, productIDList
}

func (ozon *apiClientImp) getSupplyDataList(ctx context.Context, supplyList []int) (*SupplyData, error) {
	url := fmt.Sprintf("%s%s", marketPlaceAPIURL, supplyOrderGetMethod)

	type SupplyFilterList struct {
		OrderIDs []int `json:"order_ids"`
	}

	chunkSize := 50
	chunks := chunkIntSlice(supplyList, chunkSize)

	var response SupplyData

	for _, chunk := range chunks {
		filter := SupplyFilterList{
			OrderIDs: chunk,
		}

		bodyData, err := json.Marshal(filter)
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", supplyOrderGetMethod, err)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyData))
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка создания запроса: %w", supplyOrderGetMethod, err)
		}

		for k, v := range ozonHeaders[ozon.marketPlace.ExternalID] {
			req.Header.Set(k, v)
		}
		resp, err := ozon.client.Do(req)

		if err != nil || resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s: ошибка выполнения запроса: %d", supplyOrderGetMethod, resp.StatusCode)
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", supplyOrderGetMethod, err)
		}

		response.Orders = append(response.Orders, response.Orders...)
		response.Warehouses = append(response.Warehouses, response.Warehouses...)
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

func (ozon *apiClientImp) getSupplyList(ctx context.Context) ([]int, error) {

	url := fmt.Sprintf("%s%s", marketPlaceAPIURL, supplyOrderListMethod)

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
	filter.Paging.Limit = requestItemLimit

	supplyOrderID := []int{}
	endOfList := false

	for !endOfList {
		bodyData, err := json.Marshal(filter)
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", supplyOrderListMethod, err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyData))
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка создания запроса: %w", supplyOrderListMethod, err)
		}

		for k, v := range ozonHeaders[ozon.marketPlace.ExternalID] {
			req.Header.Set(k, v)
		}
		resp, err := ozon.client.Do(req)

		if err != nil || resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s: ошибка выполнения запроса: %s", supplyOrderListMethod, err.Error())
		}
		defer resp.Body.Close()

		var response SupplyOrderList

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", supplyOrderListMethod, err)
		}

		supplyOrderID = append(supplyOrderID, response.SupplyOrderID...)
		filter.Paging.FromSupplyOrderID = response.LastSupplyOrderID

		if len(response.SupplyOrderID) < filter.Paging.Limit {
			endOfList = true
		}
	}

	return supplyOrderID, nil
}
