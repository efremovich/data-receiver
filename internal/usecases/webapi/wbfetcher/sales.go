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

type SalesResponse struct {
	Date              string  `json:"date"`
	LastChangeDate    string  `json:"lastChangeDate"`
	WarehouseName     string  `json:"warehouseName"`
	CountryName       string  `json:"countryName"`
	OblastOkrugName   string  `json:"oblastOkrugName"`
	RegionName        string  `json:"regionName"`
	SupplierArticle   string  `json:"supplierArticle"`
	NmID              int     `json:"nmId"`
	Barcode           string  `json:"barcode"`
	Category          string  `json:"category"`
	Subject           string  `json:"subject"`
	Brand             string  `json:"brand"`
	TechSize          string  `json:"techSize"`
	IncomeID          int     `json:"incomeID"`
	IsSupply          bool    `json:"isSupply"`
	IsRealization     bool    `json:"isRealization"`
	TotalPrice        float32 `json:"totalPrice"`
	DiscountPercent   float32 `json:"discountPercent"`
	Spp               int     `json:"spp"`
	PaymentSaleAmount float32 `json:"paymentSaleAmount"`
	ForPay            float32 `json:"forPay"`
	FinishedPrice     float32 `json:"finishedPrice"`
	PriceWithDisc     float32 `json:"priceWithDisc"`
	SaleID            string  `json:"saleID"`
	OrderType         string  `json:"orderType"`
	Sticker           string  `json:"sticker"`
	GNumber           string  `json:"gNumber"`
	Srid              string  `json:"srid"`
}

func (wb *apiClientImp) GetSales(ctx context.Context, desc entity.PackageDescription) ([]entity.Sale, error) {
	const methodName = "/api/v1/supplier/sales"

	urlValue := url.Values{}
	urlValue.Set("dateFrom", desc.UpdatedAt.Format("2006-01-02"))
	urlValue.Set("flag", "1")

	reqURL := fmt.Sprintf("%s%s?%s", statisticAPIURL, methodName, urlValue.Encode())

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

	var saleResponce []SalesResponse
	if err := json.NewDecoder(resp.Body).Decode(&saleResponce); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %s", methodName, err.Error())
	}

	// TODO Уточнить какую скидку ставить
	// sale.DiscountP = 0
	// Попробуем получить дату заказа
	sales := fillSaleStruct(saleResponce)

	return sales, nil
}

func fillSaleStruct(saleResponce []SalesResponse) []entity.Sale {
	var sales []entity.Sale

	for _, elem := range saleResponce {
		warehouse := entity.Warehouse{}
		warehouse.Title = elem.WarehouseName

		barcode := entity.Barcode{}
		barcode.Barcode = elem.Barcode

		vendorCode := elem.SupplierArticle
		if reVendorCode.MatchString(elem.SupplierArticle) {
			vendorCode = reVendorCode.FindString(elem.SupplierArticle)
		}

		card := entity.Card{}
		card.ExternalID = int64(elem.NmID)
		card.VendorCode = vendorCode
		card.VendorID = vendorCode

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
		size := entity.Size{
			TechSize: elem.TechSize,
			Title:    elem.TechSize,
		}

		priceSize := entity.PriceSize{
			Price:        elem.TotalPrice,
			Discount:     elem.DiscountPercent,
			SpecialPrice: elem.PriceWithDisc,
		}

		order := &entity.Order{
			ExternalID: elem.Srid,
		}
		sale := entity.Sale{}
		sale.ExternalID = elem.SaleID
		sale.Order = order
		sale.Price = elem.TotalPrice
		sale.Type = elem.OrderType
		sale.DiscountP = elem.DiscountPercent
		sale.ForPay = elem.ForPay
		sale.FinalPrice = elem.FinishedPrice

		sale.CreatedAt, _ = time.Parse("2006-01-02T15:04:05", elem.Date)

		sale.Status = &status
		sale.Region = &region
		sale.Warehouse = &warehouse
		sale.Card = &card
		sale.Barcode = &barcode
		sale.Size = &size
		sale.PriceSize = &priceSize

		sales = append(sales, sale)
	}

	return sales
}
