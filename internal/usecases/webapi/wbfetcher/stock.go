package wbfetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/efremovich/data-receiver/internal/entity"
)

type StrockResponce struct {
	LastChangeDate  string `json:"lastChangeDate"`
	WarehouseName   string `json:"warehouseName"`
	SupplierArticle string `json:"supplierArticle"`
	NmID            int    `json:"nmId"`
	Barcode         string `json:"barcode"`
	Quantity        int    `json:"quantity"`
	InWayToClient   int    `json:"inWayToClient"`
	InWayFromClient int    `json:"inWayFromClient"`
	QuantityFull    int    `json:"quantityFull"`
	Category        string `json:"category"`
	Subject         string `json:"subject"`
	Brand           string `json:"brand"`
	TechSize        string `json:"techSize"`
	Price           int    `json:"Price"`
	Discount        int    `json:"Discount"`
	IsSupply        bool   `json:"isSupply"`
	IsRealization   bool   `json:"isRealization"`
	SCCode          string `json:"SCCode"`
}

func (wb *apiClientImp) GetStocks(ctx context.Context, desc entity.PackageDescription) ([]entity.StockMeta, error) {
	const methodName = "/api/v1/supplier/stocks"

	urlValue := url.Values{}
	urlValue.Set("dateFrom", desc.UpdatedAt.Format("2006-01-02"))

	reqURL := fmt.Sprintf("%s%s?%s", statisticApiURL, methodName, urlValue.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)

	if err != nil {
		return nil, fmt.Errorf("%s: ошибка создания запроса: %w", methodName, err)
	}

	req.Header.Set("Authorization", wb.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	resp, err := wb.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: ошибка отправки запроса: %w", methodName, err)
	}
	defer resp.Body.Close()

	var stockResponce []StrockResponce
	if err := json.NewDecoder(resp.Body).Decode(&stockResponce); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
	}

	var stockMetaList []entity.StockMeta

	for _, elem := range stockResponce {
		stockMeta := entity.StockMeta{}

		stockMeta.Stock = entity.Stock{
			Quantity:        elem.Quantity,
			InWayToClient:   elem.InWayToClient,
			InWayFromClient: elem.InWayToClient,
			CreatedAt:       desc.UpdatedAt,
		}

		stockMeta.PriceSize = entity.PriceSize{
			Price:        float32(elem.Price),
			Discount:     float32(elem.Discount),
			SpecialPrice: 0,
		}

		stockMeta.Seller2Card = entity.Seller2Card{
			ExternalID: int64(elem.NmID),
		}

		stockMeta.Barcode = entity.Barcode{
			Barcode: elem.Barcode,
		}

		stockMeta.Warehouse = entity.Warehouse{
			Title: elem.WarehouseName,
		}

		stockMeta.Size = entity.Size{
			TechSize: elem.TechSize,
			Title:    elem.TechSize,
		}

		vendorID := elem.SupplierArticle
		if reVendorCode.MatchString(elem.SupplierArticle) {
			vendorID = reVendorCode.FindString(elem.SupplierArticle)
		}

		stockMeta.Card = entity.Card{
			VendorID: vendorID,
		}

		stockMetaList = append(stockMetaList, stockMeta)
	}

	return stockMetaList, nil
}
