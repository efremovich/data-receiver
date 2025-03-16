package ozonfetcher

import (
	"context"
	"fmt"
	"strconv"

	"github.com/efremovich/data-receiver/internal/entity"
)

//nolint:dupl // похожий метод есть и в order.go но они задублированны не случайно
func (ozon *apiClientImp) GetSales(ctx context.Context, desc entity.PackageDescription) ([]entity.Sale, error) {
	// Загружаем только только продажи со статусом delivered
	// Возможные статусы:
	//    awaiting_packaging — ожидает упаковки,
	//    awaiting_deliver — ожидает отгрузки,
	//    delivering — доставляется,
	//    delivered — доставлено,
	//    cancelled — отменено.
	saleResponse, err := ozon.getOrersList(ctx, desc, "delivered")
	if err != nil {
		return nil, err
	}

	skus := []int{}

	for _, elem := range saleResponse.Result {
		for _, product := range elem.Products {
			skus = append(skus, product.Sku)
		}
	}

	productInfo, err := ozon.getProductInfoOnSKU(ctx, skus)
	if err != nil {
		return nil, fmt.Errorf("ошибка получение подробной информации о товаре %w", err)
	}

	var sales []entity.Sale
	for _, elem := range saleResponse.Result {
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
						Price:        fData.Price,
						Discount:     fData.TotalDiscountValue,
						SpecialPrice: fData.OldPrice,
					}
				}
			}

			sale := entity.Sale{}
			sale.ExternalID = strconv.Itoa(elem.OrderID)
			sale.Price = priceSize.Price
			sale.Type = elem.Status
			sale.Quantity = product.Quantity

			sale.CreatedAt = elem.CreatedAt

			sale.Size = &size
			sale.PriceSize = &priceSize
			sale.Status = &status
			sale.Warehouse = &warehouse
			sale.Card = &card
			sale.Barcode = &barcode
			sale.Region = &region

			sales = append(sales, sale)
		}

	}
	return sales, nil
}
