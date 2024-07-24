package wbfetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

type StockRequestData struct {
	DateFrom time.Time `json:"dateFrom"`
}

func (wb *wbAPIclientImp) GetStocks(ctx context.Context, desc entity.PackageDescription) ([]entity.StockMeta, error) {
	const methodName = "/api/v1/supplier/stocks"
	requestData, err := json.Marshal(StockRequestData{DateFrom: time.Now()})
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", methodName, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", wb.addr, methodName), bytes.NewReader(requestData))
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка создания запроса: %w", methodName, err)
	}
	req.Header.Set("Authorization", wb.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	resp, err := wb.client.Do(req)
	if err != nil && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: ошибка отправки запроса: %w", methodName, err)
	}

	var stockResponce []StrockResponce
	if err := json.NewDecoder(resp.Body).Decode(&stockResponce); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
	}
	var stockMetaList []entity.StockMeta

	for _, elem := range stockResponce {
		stockMeta := entity.StockMeta{}

		stockMeta.Stock = entity.Stock{
			Quantity:         elem.Quantity,
			InWayToClient:    elem.InWayToClient,
			InWayFromClient:  elem.InWayToClient,
			Barcode:          elem.Barcode,
		}

		stockMeta.PriceSize = entity.PriceSize{
			Price:        float32(elem.Price),
			Discount:     float32(elem.Discount),
			SpecialPrice: 0,
		}

    stockMeta.Wb2Card = entity.Wb2Card{
    	NMID:      int64(elem.NmID),
    }
		stockMetaList = append(stockMetaList, stockMeta)
	}
	return stockMetaList, nil
}
