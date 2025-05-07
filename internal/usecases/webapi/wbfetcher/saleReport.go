package wbfetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (wb *apiClientImp) GetSaleReport(ctx context.Context, desc entity.PackageDescription) ([]entity.SaleReport, error) {
	saleReportResponces, err := wb.getSaleReportResoponse(ctx, desc)
	if err != nil {
		return nil, err
	}

	var saleReports []entity.SaleReport

	for _, elem := range saleReportResponces {
		saleReport := entity.SaleReport{}

		saleReport.ExternalID = strconv.Itoa(elem.RrdID)
		saleReport.Quantity = float32(elem.Quantity)
		saleReport.RetailPrice = elem.RetailPrice
		saleReport.RetailAmount = elem.RetailAmount
		saleReport.SalePercent = int(elem.SalePercent)
		saleReport.CommissionPercent = elem.CommissionPercent
		saleReport.RetailPriceWithdiscRub = float32(elem.RetailPriceWithdiscRub)
		saleReport.DeliveryAmount = float32(elem.DeliveryAmount)
		saleReport.DeliveryCost = float32(elem.DeliveryRub)
		saleReport.ReturnAmount = float32(elem.ReturnAmount)
		saleReport.PvzReward = float32(elem.PpvzReward)
		saleReport.SellerReward = float32(elem.PpvzVw)
		saleReport.SellerRewardWithNds = float32(elem.PpvzVwNds)

		saleReport.DateFrom = convertrepotDate(elem.DateFrom)
		saleReport.DateTo = convertrepotDate(elem.DateTo)
		saleReport.CreateReportDate = convertrepotDate(elem.CreateDt)
		saleReport.OrderDate = elem.OrderDt
		saleReport.SaleDate = elem.SaleDt
		saleReport.TransactionDate = convertrepotDate(elem.RrDt)

		saleReport.SAName = elem.SaName
		saleReport.BonusTypeName = elem.BonusTypeName
		saleReport.Penalty = float32(elem.Penalty)
		saleReport.AdditionalPayment = float32(elem.AdditionalPayment)
		saleReport.AcquiringFee = float32(elem.AcquiringFee)
		saleReport.AcquiringPercent = float32(elem.AcquiringPercent)
		saleReport.AcquiringBank = elem.AcquiringBank
		saleReport.DocType = elem.DocTypeName
		saleReport.SupplierOperName = elem.SupplierOperName

		saleReport.SiteCountry = elem.SiteCountry
		saleReport.KIZ = elem.Kiz
		saleReport.StorageFee = float32(elem.StorageFee)
		saleReport.Deduction = float32(elem.Deduction)
		saleReport.Acceptance = float32(elem.Acceptance)

		pvz := entity.Pvz{}
		pvz.OfficeID = elem.PpvzOfficeID
		pvz.OfficeName = elem.PpvzOfficeName
		pvz.SupplierName = elem.PpvzSupplierName
		pvz.SupplierID = elem.PpvzSupplierID
		pvz.SupplierINN = elem.PpvzInn

		saleReport.Pvz = &pvz

		seller := wb.marketPlace
		saleReport.Seller = &seller

		barcode := entity.Barcode{}
		barcode.Barcode = elem.Barcode
		barcode.ExternalID = int64(elem.ShkID)
		barcode.SellerID = seller.ID

		saleReport.Barcode = &barcode

		size := entity.Size{}
		size.Title = elem.TSName
		size.TechSize = elem.TSName

		saleReport.Size = &size

		card := entity.Card{}
		card.ExternalID = int64(elem.NmID)
		card.VendorID = elem.SaName

		saleReport.Card = &card

		order := entity.Order{}
		order.ExternalID = elem.Srid

		saleReport.Order = &order

		warehouse := entity.Warehouse{}
		warehouse.Title = elem.OfficeName
		warehouse.SellerID = seller.ID

		saleReport.Warehouse = &warehouse

		saleReports = append(saleReports, saleReport)
	}

	return saleReports, nil
}

func (wb *apiClientImp) getSaleReportResoponse(ctx context.Context, desc entity.PackageDescription) ([]SaleReportResponce, error) {
	var saleReportResponces []SaleReportResponce

	startOfDay := time.Date(desc.UpdatedAt.Year(), desc.UpdatedAt.Month(), desc.UpdatedAt.Day(), 0, 0, 0, 0, desc.UpdatedAt.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Nanosecond)

	// Форматируем в RFC3339
	startOfDayRFC3339 := startOfDay.Format(time.RFC3339)
	endOfDayRFC3339 := endOfDay.Format(time.RFC3339)
	urlValue := url.Values{}
	urlValue.Set("dateFrom", startOfDayRFC3339)
	urlValue.Set("dateTo", endOfDayRFC3339)
	urlValue.Set("limit", strconv.Itoa(saleReportResponseLimit))
	urlValue.Set("rrdid", "") // Начальное значение

	run := true
	for run {
		reqURL := fmt.Sprintf("%s%s?%s", statisticAPIURL, reportDetailByPeriodMethod, urlValue.Encode())

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка создания запроса: %s", reportDetailByPeriodMethod, err.Error())
		}

		req.Header.Set("Authorization", wb.token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("accept", "Application/json")

		resp, err := wb.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка отправки запроса: %s", reportDetailByPeriodMethod, err.Error())
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s: сервер ответил: %d", reportDetailByPeriodMethod, resp.StatusCode)
		}

		defer resp.Body.Close()

		var response []SaleReportResponce
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %s", reportDetailByPeriodMethod, err.Error())
		}

		saleReportResponces = append(saleReportResponces, response...)

		if len(response) == 0 {
			run = false
		} else {
			// Указатель на последную строку
			rriID := response[len(response)-1].RrdID
			urlValue.Set("rrdid", strconv.Itoa(rriID))
		}
	}

	return saleReportResponces, nil
}

func convertrepotDate(date string) time.Time {
	result, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}
	}

	return result
}
