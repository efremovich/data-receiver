package salereportrepo

import (
	"context"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
)

type saleReportDB struct {
	ID                     int64     `db:"id"`
	ExternalID             string    `db:"external_id"`               // Уникальный идентификатор заказа.
	UpdatedAt              time.Time `db:"updated_at"`                // Дата обновления данных
	Quantity               float32   `db:"quantity"`                  // Количество
	RetailPrice            float32   `db:"retail_price"`              // Цена розничная
	ReturnAmoun            float32   `db:"return_amoun"`              // Количество возвратов
	SalePercent            int       `db:"sale_percent"`              // Процент скидки
	CommissionPercent      float32   `db:"commission_percent"`        // Процент комиссии
	RetailPriceWithdiscRub float32   `db:"retail_price_withdisc_rub"` // Цена розничная с учетом скидок в рублях.
	DeliveryAmount         float32   `db:"delivery_amount"`           // Количество доставок
	ReturnAmount           float32   `db:"return_amount"`             // Количество возвратов
	DeliveryCost           float32   `db:"delivery_cost"`             // Стоимость доставки
	PvzReward              float32   `db:"pvz_reward"`                // Возмещение за выдачу на ПВЗ
	SellerReward           float32   `db:"seller_reward"`             // Возмещение марекетплейса без НДС
	SellerRewardWithNds    float32   `db:"seller_reward_with_nds"`    // Возмещение марекетплейса с НДС

	DateFrom         time.Time `db:"date_from"`          // Дата начала отчета
	DateTo           time.Time `db:"date_to"`            // Дата окончания отчета
	CreateReportDate time.Time `db:"create_report_date"` // Дата создания отчета
	OrderDate        time.Time `db:"order_date"`         // Дата заказа
	SaleDate         time.Time `db:"sale_date"`          // Дата продажи
	TransactionDate  time.Time `db:"transaction_date"`   // Дата транзакции

	SAName            string  `db:"sa_name"`            // Артикул продавца TODO Проверить нужен ли.
	BonusTypeName     string  `db:"bonus_type_name"`    // Штрафы или доплаты
	Penalty           float32 `db:"penalty"`            // Штрафы
	AdditionalPayment float32 `db:"additional_payment"` // Доплаты
	AcquiringFee      float32 `db:"acquiring_fee"`      // Возмещение издержек по эквайрингу. Издержки WB за услуги эквайринга: вычитаются из вознаграждения WB и не влияют на доход продавца
	AcquiringPercent  float32 `db:"acquiring_percent"`  // Размер комиссии за эквайринг без НДС, %
	AcquiringBank     string  `db:"acquiring_bank"`     // Банк экварйрер
	DocType           string  `db:"doc_type"`           // Тип документа
	SupplierOperName  string  `db:"supplier_oper_name"` // Обоснование оплаты
	SiteCountry       string  `db:"site_country"`       // Страна сайта
	KIZ               string  `db:"kiz"`                // Код маркировки товара
	StorageFee        float32 `db:"storage_fee"`        // Стоимость хранения
	Deduction         float32 `db:"deduction"`          // Прочие удержания и выплаты
	Acceptance        float32 `db:"acceptance"`         // Стоимость платной приемки

	PvzID       int64  `db:"pvz_id"`
	Barcode     string `db:"barcode"`
	SizeID      int64  `db:"size_id"`
	CardID      int64  `db:"card_id"`
	OrderID     int64  `db:"order_id"`
	WarehouseID int64  `db:"warehouse_id"`
	SellerID    int64  `db:"seller_id"`
}

func convertDBToSaleReport(_ context.Context, in *entity.SaleReport) *saleReportDB {
	return &saleReportDB{
		ExternalID:             in.ExternalID,
		UpdatedAt:              in.UpdatedAt,
		Quantity:               in.Quantity,
		RetailPrice:            in.RetailPrice,
		ReturnAmoun:            in.ReturnAmoun,
		SalePercent:            in.SalePercent,
		CommissionPercent:      in.CommissionPercent,
		RetailPriceWithdiscRub: in.RetailPriceWithdiscRub,
		DeliveryAmount:         in.DeliveryAmount,
		ReturnAmount:           in.ReturnAmount,
		DeliveryCost:           in.DeliveryCost,
		PvzReward:              in.PvzReward,
		SellerReward:           in.SellerReward,
		SellerRewardWithNds:    in.SellerRewardWithNds,
		DateFrom:               in.DateFrom,
		DateTo:                 in.DateTo,
		CreateReportDate:       in.CreateReportDate,
		OrderDate:              in.OrderDate,
		SaleDate:               in.SaleDate,
		TransactionDate:        in.TransactionDate,
		SAName:                 in.SAName,
		BonusTypeName:          in.BonusTypeName,
		Penalty:                in.Penalty,
		AdditionalPayment:      in.AdditionalPayment,
		AcquiringFee:           in.AcquiringFee,
		AcquiringPercent:       in.AcquiringPercent,
		AcquiringBank:          in.AcquiringBank,
		DocType:                in.DocType,
		SupplierOperName:       in.SupplierOperName,
		SiteCountry:            in.SiteCountry,
		KIZ:                    in.KIZ,
		StorageFee:             in.StorageFee,
		Deduction:              in.Deduction,
		Acceptance:             in.Acceptance,
		PvzID:                  in.Pvz.ID,
		Barcode:                in.Barcode.Barcode,
		SizeID:                 in.Size.ID,
		CardID:                 in.Card.ID,
		OrderID:                in.Order.ID,
		WarehouseID:            in.Warehouse.ID,
		SellerID:               in.Seller.ID,
	}
}

func (s saleReportDB) convertToEntitySaleReport(_ context.Context) *entity.SaleReport {
	return &entity.SaleReport{
		ID:                     s.ID,
		ExternalID:             s.ExternalID,
		UpdatedAt:              s.UpdatedAt,
		Quantity:               s.Quantity,
		RetailPrice:            s.RetailPrice,
		ReturnAmoun:            s.ReturnAmoun,
		SalePercent:            s.SalePercent,
		CommissionPercent:      s.CommissionPercent,
		RetailPriceWithdiscRub: s.RetailPriceWithdiscRub,
		DeliveryAmount:         s.DeliveryAmount,
		ReturnAmount:           s.ReturnAmount,
		DeliveryCost:           s.DeliveryCost,
		PvzReward:              s.PvzReward,
		SellerReward:           s.SellerReward,
		SellerRewardWithNds:    s.SellerRewardWithNds,
		DateFrom:               s.DateFrom,
		DateTo:                 s.DateTo,
		CreateReportDate:       s.CreateReportDate,
		OrderDate:              s.OrderDate,
		SaleDate:               s.SaleDate,
		TransactionDate:        s.TransactionDate,
		SAName:                 s.SAName,
		BonusTypeName:          s.BonusTypeName,
		Penalty:                s.Penalty,
		AdditionalPayment:      s.AdditionalPayment,
		AcquiringFee:           s.AcquiringFee,
		AcquiringPercent:       s.AcquiringPercent,
		AcquiringBank:          s.AcquiringBank,
		DocType:                s.DocType,
		SupplierOperName:       s.SupplierOperName,
		SiteCountry:            s.SiteCountry,
		KIZ:                    s.KIZ,
		StorageFee:             s.StorageFee,
		Deduction:              s.Deduction,
		Acceptance:             s.Acceptance,
		Pvz: &entity.Pvz{
			ID: s.PvzID,
		},
		Barcode: &entity.Barcode{
			Barcode: s.Barcode,
		},
		Size: &entity.Size{
			ID: s.SizeID,
		},
		Card: &entity.Card{
			ID: s.CardID,
		},
		Order: &entity.Order{
			ID: s.OrderID,
		},
		Warehouse: &entity.Warehouse{
			ID: s.WarehouseID,
		},
		Seller: &entity.MarketPlace{
			ID: s.SellerID,
		},
	}

}
