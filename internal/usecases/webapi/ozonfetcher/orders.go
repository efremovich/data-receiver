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
)

const cancelledStatus string = "cancelled"

//nolint:dupl // похожий метод есть и в sale.go но они задублированны не случайно
func (ozon *apiClientImp) GetOrders(ctx context.Context, desc entity.PackageDescription) ([]entity.Order, error) {
	// Загружаем все заказы со всеми возможными статусами
	// Возможные статусы:
	//    awaiting_packaging — ожидает упаковки,
	//    awaiting_deliver — ожидает отгрузки,
	//    delivering — доставляется,
	//    delivered — доставлено,
	//    cancelled — отменено.
	ordersResponse, err := ozon.getOrersList(ctx, desc, "")
	if err != nil {
		return nil, err
	}

	skus := []int{}

	for _, elem := range ordersResponse.Result {
		for _, product := range elem.Products {
			skus = append(skus, product.Sku)
		}
	}

	productInfo, err := ozon.getProductInfoOnSKU(ctx, skus)
	if err != nil {
		return nil, fmt.Errorf("ошибка получение подробной информации о товаре %w", err)
	}

	var orders []entity.Order

	for _, elem := range ordersResponse.Result {
		warehouse := entity.Warehouse{}
		warehouse.Title = elem.AnalyticsData.WarehouseName
		warehouse.ExternalID = elem.AnalyticsData.WarehouseID

		status := entity.Status{
			Name: elem.Status,
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
						Price:                fData.Price,
						PriceWithoutDiscount: fData.OldPrice,
						PriceFinal:           fData.Payout,
						UpdatedAt:            time.Now(),
					}
				}
			}

			order := entity.Order{}
			order.ExternalID = strconv.Itoa(elem.OrderID)
			order.Price = priceSize.Price
			order.Type = elem.Status
			order.Quantity = product.Quantity

			order.IsCancel = elem.Status == cancelledStatus

			order.CreatedAt = elem.CreatedAt

			order.Size = &size
			order.PriceSize = &priceSize
			order.Status = &status
			order.Warehouse = &warehouse
			order.Card = &card
			order.Barcode = &barcode
			order.Region = &region

			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (ozon *apiClientImp) getOrersList(ctx context.Context, desc entity.PackageDescription, status string) (OrderRespose, error) {
	url := fmt.Sprintf("%s%s", marketPlaceAPIURL, fboPostingListMethod)

	startDate := desc.UpdatedAt.Truncate(24 * time.Hour)

	filter := OrderFilter{}
	filter.Dir = "desc"
	filter.Limit = requestItemLimit
	filter.With.AnalyticsData = true
	filter.With.FinancialData = true
	filter.Filter.Since = startDate
	filter.Filter.To = startDate.Add(24 * time.Hour)
	filter.Filter.Status = status

	bodyData, err := json.Marshal(filter)
	if err != nil {
		return OrderRespose{}, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", fboPostingListMethod, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyData))
	if err != nil {
		return OrderRespose{}, fmt.Errorf("%s: ошибка создания запроса: %w", fboPostingListMethod, err)
	}

	for k, v := range ozonHeaders[ozon.marketPlace.ExternalID] {
		req.Header.Set(k, v)
	}

	resp, err := ozon.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return OrderRespose{}, fmt.Errorf("%s: ошибка выполнения запроса: %s", fboPostingListMethod, err.Error())
	}

	defer resp.Body.Close()

	var response OrderRespose
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return OrderRespose{}, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", fboPostingListMethod, err)
	}

	return response, nil
}
