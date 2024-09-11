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
	WarehouseName   string  `json:"warehouseName"`
	NmID            int     `json:"nmId"`
	Barcode         string  `json:"barcode"`
	Quantity        int     `json:"quantity"`
	InWayToClient   int     `json:"inWayToClient"`
	InWayFromClient int     `json:"inWayFromClient"`
	Category        string  `json:"category"`
	Brand           string  `json:"brand"`
	TechSize        string  `json:"techSize"`
	Price           float32 `json:"Price"`
	Discount        float32 `json:"Discount"`
}

func (wb *wbAPIclientImp) GetStocks(ctx context.Context, desc entity.PackageDescription) ([]entity.StockMeta, error) {
	const methodName = "/api/v1/supplier/stocks"

	urlValue := url.Values{}
	urlValue.Set("dateFrom", desc.UpdatedAt.Format("2006-01-02"))

	reqURL := fmt.Sprintf("%s%s?%s", wb.addrStat, methodName, urlValue.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)

	if err != nil {
		return nil, fmt.Errorf("%s: ошибка создания запроса: %w", methodName, err)
	}

	req.Header.Set("Authorization", wb.tokenStat)
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
		}

		stockMeta.PriceSize = entity.PriceSize{
			Price:        float32(elem.Price),
			Discount:     float32(elem.Discount),
			SpecialPrice: 0,
		}

		stockMeta.Wb2Card = entity.Wb2Card{
			NMID: int64(elem.NmID),
		}

		stockMeta.Barcode = entity.Barcode{
			Barcode: elem.Barcode,
		}

		stockMeta.Warehouse = entity.Warehouse{
			Title: elem.WarehouseName,
		}

		stockMeta.Size = entity.Size{
			TechSize: elem.TechSize,
		}
		stockMetaList = append(stockMetaList, stockMeta)
	}

	return stockMetaList, nil
}
