package ozonfetcher

import (
	"context"
	"fmt"
	"strconv"

	"github.com/efremovich/data-receiver/internal/entity"
)

const noInfo string = "нет информации"

func (ozon *apiClientImp) GetSaleReport(ctx context.Context, desc entity.PackageDescription) ([]entity.SaleReport, error) {
	var saleReports []entity.SaleReport
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

	if len(saleResponse.Result) == 0 {
		return nil, nil
	}
	skus := []int{}

	for _, elem := range saleResponse.Result {
		for _, product := range elem.FinancialData.Products {
			skus = append(skus, product.ProductID)
		}
	}

	productInfo, err := ozon.getProductInfoOnSKU(ctx, skus)
	if err != nil {
		return nil, fmt.Errorf("ошибка получение подробной информации о товаре %w", err)
	}

	for _, elem := range saleResponse.Result {
		for _, product := range elem.FinancialData.Products {
			saleReport := entity.SaleReport{}

			saleReport.ExternalID = strconv.Itoa(elem.OrderID)
			saleReport.Quantity = 1
			saleReport.RetailPrice = product.OldPrice
			saleReport.SalePercent = int(product.TotalDiscountPercent)
			saleReport.CommissionPercent = float32(product.CommissionPercent)
			saleReport.RetailPriceWithdiscRub = float32(product.OldPrice)
			saleReport.DeliveryAmount = 1
			saleReport.DeliveryCost = float32(product.ItemServices.MarketplaceServiceItemDelivToCustomer)
			saleReport.ReturnAmount = 0
			saleReport.PvzReward = 0
			saleReport.SellerReward = float32(product.CommissionAmount)
			saleReport.SellerRewardWithNds = float32(product.CommissionAmount)

			saleReport.DateFrom = elem.CreatedAt
			saleReport.DateTo = elem.InProcessAt
			saleReport.CreateReportDate = elem.CreatedAt
			saleReport.OrderDate = elem.CreatedAt
			saleReport.SaleDate = elem.InProcessAt
			saleReport.TransactionDate = elem.InProcessAt

			saleReport.DocType = "Продажа"
			saleReport.AcquiringBank = elem.AnalyticsData.PaymentTypeGroupName

			pvz := entity.Pvz{}
			pvz.OfficeID = 0
			pvz.OfficeName = noInfo
			pvz.SupplierName = noInfo
			pvz.SupplierID = 0
			pvz.SupplierINN = noInfo
			saleReport.Pvz = &pvz

			barcode := entity.Barcode{}

			barcodes := productInfo[product.ProductID].Barcodes
			for _, b := range barcodes {
				barcode.Barcode = b
			}

			saleReport.Barcode = &barcode

			offerID := productInfo[product.ProductID].OfferID
			vendorID, vendorCode, currSize := getMetaFromVendorID(offerID)

			size := entity.Size{}
			size.Title = currSize
			size.TechSize = currSize

			saleReport.Size = &size

			card := entity.Card{}
			card.VendorCode = vendorCode
			card.VendorID = vendorID
			card.ExternalID = int64(productInfo[product.ProductID].ID)

			saleReport.Card = &card

			order := entity.Order{}
			order.ExternalID = strconv.Itoa(elem.OrderID)

			saleReport.Order = &order

			warehouse := entity.Warehouse{}
			warehouse.Title = elem.AnalyticsData.WarehouseName
			warehouse.ExternalID = elem.AnalyticsData.WarehouseID

			saleReport.Warehouse = &warehouse

			saleReports = append(saleReports, saleReport)
		}
	}

	return saleReports, nil
}
