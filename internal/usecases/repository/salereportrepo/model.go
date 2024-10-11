package salereportrepo

import (
	"context"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
)

type saleReportDB struct {
	ID                      int64     `db:"id"`
	ExternalID              string    `db:"external_id"`
	UpdatedAt               time.Time `db:"updated_at"`
	ContractNuber           string    `db:"contract_nuber"`
	PartID                  int       `db:"part_id"`
	DocType                 string    `db:"doc_type"`
	Quantity                float32   `db:"quantity"`
	RetailPrice             float32   `db:"retail_price"`
	FinisPrice              float32   `db:"finis_price"`
	RetailSum               float32   `db:"retail_sum"`
	SalePercent             int       `db:"sale_percent"`
	CommissionPercent       float32   `db:"commission_percent"`
	WarehouseName           string    `db:"warehouse_name"`
	OrderDate               time.Time `db:"order_date"`
	SaleDate                time.Time `db:"sale_date"`
	OperationName           string    `db:"operation_name"`
	DeliveryAmount          float32   `db:"delivery_amount"`
	ReturnAmoun             float32   `db:"return_amoun"`
	DeliveryCost            float32   `db:"delivery_cost"`
	PackageType             string    `db:"package_type"`
	ProductDiscount         float32   `db:"product_discount"`
	BuyerDiscount           float32   `db:"buyer_discount"`
	BaseRatioDiscount       float32   `db:"base_ratio_discount"`
	TotalRatioDiscount      float32   `db:"total_ratio_discount"`
	ReductionRatingRatio    float32   `db:"reduction_rating_ratio"`
	ReductionPromotionRatio float32   `db:"reduction_promotion_ratio"`
	RewardRatio             float32   `db:"reward_ratio"`
	ForPay                  float32   `db:"for_pay"`
	Reward                  float32   `db:"reward"`
	AcquiringFee            float32   `db:"acquiring_fee"`
	AcquiringPercent        float32   `db:"acquiring_percent"`
	AcquiringBank           string    `db:"acquiring_bank"`
	SellerReward            float32   `db:"seller_reward"`
	OfficeID                int       `db:"office_id"`
	OfficeName              string    `db:"office_name"`
	SupplierID              int       `db:"supplier_id"`
	SupplierName            string    `db:"supplier_name"`
	SupplierINN             string    `db:"supplier_inn"`
	DeclarationNumber       string    `db:"declaration_number"`
	BonusTypeName           string    `db:"bonus_type_name"`
	StickerID               string    `db:"sticker_id"`
	CountryOfSale           string    `db:"country_of_sale"`
	Penalty                 float32   `db:"penalty"`
	AdditionalPayment       float32   `db:"additional_payment"`
	RebillLogisticCost      float32   `db:"rebill_logistic_cost"`
	RebillLogisticOrg       string    `db:"rebill_logistic_org"`
	KIZ                     string    `db:"kiz"`
	StorageFee              float32   `db:"storage_fee"`
	Deduction               float32   `db:"deduction"`
	Acceptance              float32   `db:"acceptance"`
	OrderID                 int64     `db:"order_id"`
	SellerID                int64     `db:"seller_id"`
	CardID                  int64     `db:"card_id"`
	Barcode                 string    `db:"barcode"`
}

func convertDBToSaleReport(_ context.Context, in *entity.SaleReport) *saleReportDB {
	return &saleReportDB{
		ID:                in.ID,
		ExternalID:        in.ExternalID,
		UpdatedAt:         in.UpdatedAt,
		ContractNuber:     in.ContractNumber,
		PartID:            in.PartID,
		DocType:           in.DocType,
		Quantity:          in.Quantity,
		RetailPrice:       in.RetailPrice,
		FinisPrice:        in.FinisPrice,
		RetailSum:         in.RetailSum,
		SalePercent:       in.SalePercent,
		CommissionPercent: in.CommissionPercent,
		WarehouseName:     in.WarehouseName,
		OrderDate:         in.OrderDate,
		SaleDate:          in.SaleDate,
		OperationName:     in.OperationName,
		DeliveryAmount:    in.DeliveryAmount,
		ReturnAmoun:       in.ReturnAmoun,
		DeliveryCost:      in.DeliveryCost,
		PackageType:       in.PackageType,
		ProductDiscount:   in.ProductDiscount,

		BuyerDiscount:           in.Pvz.BuyerDiscount,
		BaseRatioDiscount:       in.Pvz.BaseRatioDiscount,
		TotalRatioDiscount:      in.Pvz.TotalRatioDiscount,
		ReductionRatingRatio:    in.Pvz.ReductionRatingRatio,
		ReductionPromotionRatio: in.Pvz.ReductionPromotionRatio,
		RewardRatio:             in.Pvz.RewardRatio,
		ForPay:                  in.Pvz.ForPay,
		Reward:                  in.Pvz.Reward,
		AcquiringFee:            in.Pvz.AcquiringFee,
		AcquiringPercent:        in.Pvz.AcquiringPercent,
		AcquiringBank:           in.Pvz.AcquiringBank,
		SellerReward:            in.Pvz.SellerReward,
		OfficeID:                in.Pvz.OfficeID,
		OfficeName:              in.Pvz.OfficeName,
		SupplierID:              in.Pvz.SupplierID,
		SupplierName:            in.Pvz.SupplierName,
		SupplierINN:             in.Pvz.SupplierINN,
		DeclarationNumber:       in.Pvz.DeclarationNumber,
		BonusTypeName:           in.Pvz.BonusTypeName,
		StickerID:               in.Pvz.StickerID,
		CountryOfSale:           in.Pvz.CountryOfSale,
		Penalty:                 in.Pvz.Penalty,
		AdditionalPayment:       in.Pvz.AdditionalPayment,
		RebillLogisticCost:      in.Pvz.RebillLogisticCost,
		RebillLogisticOrg:       in.Pvz.RebillLogisticOrg,
		KIZ:                     in.Pvz.KIZ,
		StorageFee:              in.Pvz.StorageFee,
		Deduction:               in.Pvz.Deduction,

		OrderID:  in.Order.ID,
		SellerID: in.Seller.ID,
		CardID:   in.Card.ID,
		Barcode:  in.Barcode.Barcode,
	}
}

func (s saleReportDB) convertToEntitySaleReport(_ context.Context) *entity.SaleReport {
	return &entity.SaleReport{
		ID:                s.ID,
		ExternalID:        s.ExternalID,
		UpdatedAt:         s.UpdatedAt,
		ContractNumber:    s.ContractNuber,
		PartID:            s.PartID,
		DocType:           s.DocType,
		Quantity:          s.Quantity,
		RetailPrice:       s.RetailPrice,
		FinisPrice:        s.FinisPrice,
		RetailSum:         s.RetailSum,
		SalePercent:       s.SalePercent,
		CommissionPercent: s.CommissionPercent,
		WarehouseName:     s.WarehouseName,
		OrderDate:         s.OrderDate,
		SaleDate:          s.SaleDate,
		OperationName:     s.OperationName,
		DeliveryAmount:    s.DeliveryAmount,
		ReturnAmoun:       s.ReturnAmoun,
		DeliveryCost:      s.DeliveryCost,
	}
}
