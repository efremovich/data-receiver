package wbfetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/efremovich/data-receiver/internal/entity"
)

type OrdersResponce struct {
	Date            string `json:"date"`
	LastChangeDate  string `json:"lastChangeDate"`
	WarehouseName   string `json:"warehouseName"`
	CountryName     string `json:"countryName"`
	OblastOkrugName string `json:"oblastOkrugName"`
	RegionName      string `json:"regionName"`
	SupplierArticle string `json:"supplierArticle"`
	NmID            int    `json:"nmId"`
	Barcode         string `json:"barcode"`
	Category        string `json:"category"`
	Subject         string `json:"subject"`
	Brand           string `json:"brand"`
	TechSize        string `json:"techSize"`
	IncomeID        int    `json:"incomeID"`
	IsSupply        bool   `json:"isSupply"`
	IsRealization   bool   `json:"isRealization"`
	TotalPrice      float32    `json:"totalPrice"`
	DiscountPercent float32    `json:"discountPercent"`
	Spp             int    `json:"spp"`
	FinishedPrice   float32    `json:"finishedPrice"`
	PriceWithDisc   float32    `json:"priceWithDisc"`
	IsCancel        bool   `json:"isCancel"`
	CancelDate      string `json:"cancelDate"`
	OrderType       string `json:"orderType"`
	Sticker         string `json:"sticker"`
	GNumber         string `json:"gNumber"`
	Srid            string `json:"srid"`
}

func (wb *wbAPIclientImp) GetOrders(ctx context.Context, desc entity.PackageDescription) ([]entity.OrderMeta, error) {
	const methodName = "/api/v1/supplier/orders"

	urlValue := url.Values{}
	urlValue.Set("dateFrom", desc.UpdatedAt.Format("2006-01-02"))
	urlValue.Set("flag", "0")

	reqUrl := fmt.Sprintf("%s%s?%s", wb.addrStat, methodName, urlValue.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
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

	var orderResponce []OrdersResponce
	if err := json.NewDecoder(resp.Body).Decode(&orderResponce); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
	}
	var ordersMeta []entity.OrderMeta

	for _, elem := range orderResponce {
		warehouse := entity.Warehouse{}
		warehouse.Title = elem.WarehouseName

		barcode := entity.Barcode{}
		barcode.Barcode = elem.Barcode

		card := entity.Card{}
		card.ExternalID = int64(elem.NmID)
		card.VendorCode = elem.SupplierArticle

		order := entity.Order{}
		order.ExternalID = elem.Srid
		order.Price = elem.TotalPrice
		order.Discount = elem.DiscountPercent
		order.SpecialPrice = elem.FinishedPrice
		order.Status = elem.OrderType

		orderMeta := entity.OrderMeta{
			Warehouse: warehouse,
			Barcode:   barcode,
			Card:      card,
			Order:     order,
		}
		ordersMeta = append(ordersMeta, orderMeta)
	}
	return ordersMeta, nil
}
