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

type SaleReportResponce struct {
	RealizationreportID      int       `json:"realizationreport_id"`
	DateFrom                 string    `json:"date_from"`
	DateTo                   string    `json:"date_to"`
	CreateDt                 string    `json:"create_dt"`
	CurrencyName             string    `json:"currency_name"`
	SuppliercontractCode     any       `json:"suppliercontract_code"`
	RrdID                    int       `json:"rrd_id"`
	GiID                     int       `json:"gi_id"`
	DlvPrc                   float64   `json:"dlv_prc"`
	FixTariffDateFrom        string    `json:"fix_tariff_date_from"`
	FixTariffDateTo          string    `json:"fix_tariff_date_to"`
	SubjectName              string    `json:"subject_name"`
	NmID                     int       `json:"nm_id"`
	BrandName                string    `json:"brand_name"`
	SaName                   string    `json:"sa_name"`
	TSName                   string    `json:"ts_name"`
	Barcode                  string    `json:"barcode"`
	DocTypeName              string    `json:"doc_type_name"`
	Quantity                 int       `json:"quantity"`
	RetailPrice              float32   `json:"retail_price"`
	RetailAmount             float32   `json:"retail_amount"`
	SalePercent              float32   `json:"sale_percent"`
	CommissionPercent        float32   `json:"commission_percent"`
	OfficeName               string    `json:"office_name"`
	SupplierOperName         string    `json:"supplier_oper_name"`
	OrderDt                  time.Time `json:"order_dt"`
	SaleDt                   time.Time `json:"sale_dt"`
	RrDt                     string    `json:"rr_dt"`
	ShkID                    int       `json:"shk_id"`
	RetailPriceWithdiscRub   float64   `json:"retail_price_withdisc_rub"`
	DeliveryAmount           int       `json:"delivery_amount"`
	ReturnAmount             int       `json:"return_amount"`
	DeliveryRub              float64   `json:"delivery_rub"`
	GiBoxTypeName            string    `json:"gi_box_type_name"`
	ProductDiscountForReport float64   `json:"product_discount_for_report"`
	SupplierPromo            int       `json:"supplier_promo"`
	Rid                      int64     `json:"rid"`
	PpvzSppPrc               float64   `json:"ppvz_spp_prc"`
	PpvzKvwPrcBase           float64   `json:"ppvz_kvw_prc_base"`
	PpvzKvwPrc               float64   `json:"ppvz_kvw_prc"`
	SupRatingPrcUp           int       `json:"sup_rating_prc_up"`
	IsKgvpV2                 float64   `json:"is_kgvp_v2"`
	PpvzSalesCommission      float64   `json:"ppvz_sales_commission"`
	PpvzForPay               float64   `json:"ppvz_for_pay"`
	PpvzReward               float64   `json:"ppvz_reward"`
	AcquiringFee             float64   `json:"acquiring_fee"`
	AcquiringPercent         float64   `json:"acquiring_percent"`
	PaymentProcessing        string    `json:"payment_processing"`
	AcquiringBank            string    `json:"acquiring_bank"`
	PpvzVw                   float64   `json:"ppvz_vw"`
	PpvzVwNds                float64   `json:"ppvz_vw_nds"`
	PpvzOfficeName           string    `json:"ppvz_office_name"`
	PpvzOfficeID             int       `json:"ppvz_office_id"`
	PpvzSupplierID           int       `json:"ppvz_supplier_id"`
	PpvzSupplierName         string    `json:"ppvz_supplier_name"`
	PpvzInn                  string    `json:"ppvz_inn"`
	DeclarationNumber        string    `json:"declaration_number"`
	BonusTypeName            string    `json:"bonus_type_name"`
	StickerID                string    `json:"sticker_id"`
	SiteCountry              string    `json:"site_country"`
	SrvDbs                   bool      `json:"srv_dbs"`
	Penalty                  float64   `json:"penalty"`
	AdditionalPayment        int       `json:"additional_payment"`
	RebillLogisticCost       float64   `json:"rebill_logistic_cost"`
	RebillLogisticOrg        string    `json:"rebill_logistic_org"`
	StorageFee               float64   `json:"storage_fee"`
	Deduction                int       `json:"deduction"`
	Acceptance               int       `json:"acceptance"`
	AssemblyID               int64     `json:"assembly_id"`
	Kiz                      string    `json:"kiz"`
	Srid                     string    `json:"srid"`
	ReportType               int       `json:"report_type"`
	IsLegalEntity            bool      `json:"is_legal_entity"`
	TrbxID                   string    `json:"trbx_id"`
}

const LIMIT = "100000" // Максимальное количество строк отчета, возвращаемых методом. Не может быть более 100000.

func (wb *apiClientImp) GetSaleReport(ctx context.Context, desc entity.PackageDescription) ([]entity.SaleReport, error) {
	const methodName = "/api/v5/supplier/reportDetailByPeriod"
	startOfDay := time.Date(desc.UpdatedAt.Year(), desc.UpdatedAt.Month(), desc.UpdatedAt.Day(), 0, 0, 0, 0, desc.UpdatedAt.Location())

	endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Nanosecond)

	// Форматируем в RFC3339
	startOfDayRFC3339 := startOfDay.Format(time.RFC3339)
	endOfDayRFC3339 := endOfDay.Format(time.RFC3339)
	rrdid := desc.GetCursor()
	urlValue := url.Values{}
	urlValue.Set("dateFrom", startOfDayRFC3339)
	urlValue.Set("dateTo", endOfDayRFC3339)
	urlValue.Set("limit", LIMIT)
	urlValue.Set("flag", "1")
	urlValue.Set("rrdid", rrdid) // Начальное значение

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

	var saleReportResponces []SaleReportResponce
	if err := json.NewDecoder(resp.Body).Decode(&saleReportResponces); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %s", methodName, err.Error())
	}

	var saleReports []entity.SaleReport

	for _, elem := range saleReportResponces {
		saleReport := entity.SaleReport{}

		saleReport.ExternalID = strconv.Itoa(elem.RrdID)
		saleReport.Quantity = float32(elem.Quantity)
		saleReport.RetailPrice = elem.RetailPrice
		saleReport.SalePercent = int(elem.SalePercent)
		saleReport.CommissionPercent = elem.CommissionPercent
		saleReport.RetailPriceWithdiscRub = float32(elem.RetailPriceWithdiscRub)
		saleReport.DeliveryAmount = float32(elem.DeliveryAmount)
		saleReport.DeliveryCost = float32(elem.DeliveryRub)
		saleReport.ReturnAmount = float32(elem.ReturnAmount)
		saleReport.PvzReward = float32(elem.PpvzReward)
		saleReport.SellerReward = float32(elem.PpvzVw)
		saleReport.SellerRewardWithNds = float32(elem.PpvzVwNds)

		saleReport.DateFrom = convertrepotDate(elem.DateFrom, "2006-01-02")
		saleReport.DateTo = convertrepotDate(elem.DateTo, "2006-01-02")
		saleReport.CreateReportDate = convertrepotDate(elem.CreateDt, "2006-01-02")
		saleReport.OrderDate = elem.OrderDt
		saleReport.SaleDate = elem.SaleDt
		saleReport.TransactionDate = convertrepotDate(elem.RrDt, "2006-01-02")

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

func convertrepotDate(date string, layout string) time.Time {
	result, err := time.Parse(layout, date)
	if err != nil {
		return time.Time{}
	}
	return result
}
