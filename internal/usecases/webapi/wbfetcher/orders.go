package wbfetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
)

type OrdersResponce struct {
	Date            string  `json:"date"`
	LastChangeDate  string  `json:"lastChangeDate"`
	WarehouseName   string  `json:"warehouseName"`
	CountryName     string  `json:"countryName"`
	OblastOkrugName string  `json:"oblastOkrugName"`
	RegionName      string  `json:"regionName"`
	SupplierArticle string  `json:"supplierArticle"`
	NmID            int     `json:"nmId"`
	Barcode         string  `json:"barcode"`
	Category        string  `json:"category"`
	Subject         string  `json:"subject"`
	Brand           string  `json:"brand"`
	TechSize        string  `json:"techSize"`
	IncomeID        int     `json:"incomeID"`
	IsSupply        bool    `json:"isSupply"`
	IsRealization   bool    `json:"isRealization"`
	TotalPrice      float32 `json:"totalPrice"`
	DiscountPercent float32 `json:"discountPercent"`
	Spp             int     `json:"spp"`
	FinishedPrice   float32 `json:"finishedPrice"`
	PriceWithDisc   float32 `json:"priceWithDisc"`
	IsCancel        bool    `json:"isCancel"`
	CancelDate      string  `json:"cancelDate"`
	OrderType       string  `json:"orderType"`
	Sticker         string  `json:"sticker"`
	GNumber         string  `json:"gNumber"`
	Srid            string  `json:"srid"`
}

func (wb *apiClientImp) GetOrders(ctx context.Context, desc entity.PackageDescription) ([]entity.Order, error) {
	const methodName = "/api/v1/supplier/orders"

	urlValue := url.Values{}
	urlValue.Set("dateFrom", desc.UpdatedAt.Format("2006-01-02"))
	urlValue.Set("flag", "1")

	reqURL := fmt.Sprintf("%s%s?%s", statisticApiURL, methodName, urlValue.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка создания запроса: %s", methodName, err.Error())
	}

	req.Header.Set("Authorization", wb.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	resp, err := wb.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка отправки запроса: %s", methodName, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: сервер ответил: %d", methodName, resp.StatusCode)
	}

	defer resp.Body.Close()

	var orderResponce []OrdersResponce
	if err := json.NewDecoder(resp.Body).Decode(&orderResponce); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %s", methodName, err.Error())
	}

	// На ВБ 1 товар одна строка с данными
	// Попробуем получить дату заказа
	orders := fillOrderStruct(orderResponce)

	return orders, nil
}

func fillOrderStruct(orderResponce []OrdersResponce) []entity.Order {
	var orders []entity.Order

	for _, elem := range orderResponce {
		warehouse := entity.Warehouse{}
		warehouse.Title = elem.WarehouseName

		barcode := entity.Barcode{}
		barcode.Barcode = elem.Barcode

		vendorID := elem.SupplierArticle
		if reVendorCode.MatchString(elem.SupplierArticle) {
			vendorID = reVendorCode.FindString(elem.SupplierArticle)
		}

		card := entity.Card{}
		card.ExternalID = int64(elem.NmID)
		card.VendorID = vendorID

		status := entity.Status{
			Name: elem.OrderType,
		}

		region := entity.Region{
			RegionName: elem.RegionName,
			District: entity.District{
				Name: elem.RegionName,
			},
			Country: entity.Country{
				Name: elem.CountryName,
			},
		}

		seller := entity.MarketPlace{
			Title: "wb",
		}

		priceSize := entity.PriceSize{
			Price:        elem.FinishedPrice,
			Discount:     elem.DiscountPercent,
			SpecialPrice: elem.TotalPrice,
		}

		size := entity.Size{
			TechSize: elem.TechSize,
			Title:    elem.TechSize,
		}

		order := entity.Order{}
		order.ExternalID = elem.Srid
		order.Price = elem.TotalPrice
		order.Type = elem.OrderType
		order.Quantity = 1

		order.CreatedAt, _ = time.Parse("2006-01-02T15:04:05", elem.Date)

		order.Size = &size
		order.PriceSize = &priceSize
		order.Status = &status
		order.Region = &region
		order.Warehouse = &warehouse
		order.Seller = &seller
		order.Card = &card
		order.Barcode = &barcode

		orders = append(orders, order)
	}

	return orders
}
