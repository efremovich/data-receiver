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

type SaleReportResponce struct {
	RealizationreportID      int         `json:"realizationreport_id"`
	DateFrom                 string      `json:"date_from"`
	DateTo                   string      `json:"date_to"`
	CreateDt                 string      `json:"create_dt"`
	CurrencyName             string      `json:"currency_name"`
	SuppliercontractCode     interface{} `json:"suppliercontract_code"`
	RrdID                    int         `json:"rrd_id"`
	GiID                     int         `json:"gi_id"`
	SubjectName              string      `json:"subject_name"`
	NmID                     int         `json:"nm_id"`
	BrandName                string      `json:"brand_name"`
	SaName                   string      `json:"sa_name"`
	TsName                   string      `json:"ts_name"`
	Barcode                  string      `json:"barcode"`
	DocTypeName              string      `json:"doc_type_name"`
	Quantity                 int         `json:"quantity"`
	RetailPrice              int         `json:"retail_price"`
	RetailAmount             int         `json:"retail_amount"`
	SalePercent              int         `json:"sale_percent"`
	CommissionPercent        float64     `json:"commission_percent"`
	OfficeName               string      `json:"office_name"`
	SupplierOperName         string      `json:"supplier_oper_name"`
	OrderDt                  time.Time   `json:"order_dt"`
	SaleDt                   time.Time   `json:"sale_dt"`
	RrDt                     string      `json:"rr_dt"`
	ShkID                    int         `json:"shk_id"`
	RetailPriceWithdiscRub   float64     `json:"retail_price_withdisc_rub"`
	DeliveryAmount           int         `json:"delivery_amount"`
	ReturnAmount             int         `json:"return_amount"`
	DeliveryRub              int         `json:"delivery_rub"`
	GiBoxTypeName            string      `json:"gi_box_type_name"`
	ProductDiscountForReport float64     `json:"product_discount_for_report"`
	SupplierPromo            int         `json:"supplier_promo"`
	Rid                      int64       `json:"rid"`
	PpvzSppPrc               float64     `json:"ppvz_spp_prc"`
	PpvzKvwPrcBase           float64     `json:"ppvz_kvw_prc_base"`
	PpvzKvwPrc               float64     `json:"ppvz_kvw_prc"`
	SupRatingPrcUp           int         `json:"sup_rating_prc_up"`
	IsKgvpV2                 int         `json:"is_kgvp_v2"`
	PpvzSalesCommission      float64     `json:"ppvz_sales_commission"`
	PpvzForPay               float64     `json:"ppvz_for_pay"`
	PpvzReward               int         `json:"ppvz_reward"`
	AcquiringFee             float64     `json:"acquiring_fee"`
	AcquiringPercent         float64     `json:"acquiring_percent"`
	AcquiringBank            string      `json:"acquiring_bank"`
	PpvzVw                   float64     `json:"ppvz_vw"`
	PpvzVwNds                float64     `json:"ppvz_vw_nds"`
	PpvzOfficeID             int         `json:"ppvz_office_id"`
	PpvzOfficeName           string      `json:"ppvz_office_name"`
	PpvzSupplierID           int         `json:"ppvz_supplier_id"`
	PpvzSupplierName         string      `json:"ppvz_supplier_name"`
	PpvzInn                  string      `json:"ppvz_inn"`
	DeclarationNumber        string      `json:"declaration_number"`
	BonusTypeName            string      `json:"bonus_type_name"`
	StickerID                string      `json:"sticker_id"`
	SiteCountry              string      `json:"site_country"`
	Penalty                  float64     `json:"penalty"`
	AdditionalPayment        int         `json:"additional_payment"`
	RebillLogisticCost       float64     `json:"rebill_logistic_cost"`
	RebillLogisticOrg        string      `json:"rebill_logistic_org"`
	Kiz                      string      `json:"kiz"`
	StorageFee               float64     `json:"storage_fee"`
	Deduction                int         `json:"deduction"`
	Acceptance               int         `json:"acceptance"`
	Srid                     string      `json:"srid"`
	ReportType               int         `json:"report_type"`
}

const LIMIT = "100000" // Максимальное количество строк отчета, возвращаемых методом. Не может быть более 100000.
// https://openapi.wildberries.ru/statistics/api/ru/#tag/Statistika/paths/~1api~1v5~1supplier~1reportDetailByPeriod/get
func (wb *apiClientImp) GetSaleReport(ctx context.Context, desc entity.PackageDescription) ([]entity.SaleReport, error) {
	const methodName = "/api/v5/supplier/reportDetailByPeriod"

	rrdid := desc.Cursor
	urlValue := url.Values{}
	urlValue.Set("dateFrom", desc.UpdatedAt.Format("2006-01-02 00:00:00"))
	urlValue.Set("dateTo", desc.UpdatedAt.Format("2006-01-02 23:59:59"))
	urlValue.Set("limit", LIMIT)
	urlValue.Set("flag", "1")
	urlValue.Set("rrdid", rrdid) // Начальное значение

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

	var saleReportResponces []SaleReportResponce
	if err := json.NewDecoder(resp.Body).Decode(&saleReportResponces); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %s", methodName, err.Error())
	}

	var saleReports []entity.SaleReport

	for _, elem := range saleReportResponces {
		warehouse := entity.Warehouse{}
		warehouse.Title = elem.OfficeName

		barcode := entity.Barcode{}
		barcode.Barcode = elem.Barcode

		card := entity.Card{}
		card.ExternalID = int64(elem.NmID)
		card.VendorCode = elem.SaName

		saleReport := entity.SaleReport{}
		saleReport.Card = card
		saleReport.Barcode = barcode

		saleReports = append(saleReports, saleReport)

	}

	return saleReports, nil
}
