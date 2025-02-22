package entity

import "time"

type SaleReport struct {
	ID                int64
	ExternalID        string    // Уникальный идентификатор заказа.
	UpdatedAt         time.Time // Дата обновления данных
	ContractNumber    string    // Договор
	PartID            int       // Номер поставки
	DocType           string    // Тип документа
	Quantity          float32   // Количество
	RetailPrice       float32   // Цена розничная
	FinisPrice        float32   //
	RetailSum         float32   // Сумма продажи (возврата)
	SalePercent       int       // Процент скидки
	CommissionPercent float32   // Процент комиссии
	WarehouseName     string    // Наименование склада
	OrderDate         time.Time // Дата заказа
	SaleDate          time.Time // Дата продажи
	OperationName     string    // Наименование операции
	DeliveryAmount    float32   // Количество доставок
	ReturnAmoun       float32   // Количество возвратов
	DeliveryCost      float32   // Стоимость доставки
	PackageType       string    // Тип упаковки
	ProductDiscount   float32   // Финальная скидка
	Pvz               Pvz
	Card              Card
	Order             Order
	Barcode           Barcode
	Seller            MarketPlace
}

// кВВ - коэффициент вознаграждения вайлдерис.
type Pvz struct {
	BuyerDiscount           float32 // Скидка постоянного покупателя
	BaseRatioDiscount       float32 // Размер кВВ без НДС, % базовый
	TotalRatioDiscount      float32 // Итоговый кВВ без НДС, %
	ReductionRatingRatio    float32 // Размер снижения кВВ из-за рейтинга
	ReductionPromotionRatio float32 // Размер снижения кВВ из-за акции
	RewardRatio             float32 // Вознаграждение с продаж до вычета услуг поверенного, без НДС
	ForPay                  float32 // К перечислению продавцу за реализованный товар
	Reward                  float32 // Возмещение за выдачу и возврат товаров на ПВЗ
	AcquiringFee            float32 // Возмещение издержек по эквайрингу. Издержки WB за услуги эквайринга: вычитаются из вознаграждения WB и не влияют на доход продавца
	AcquiringPercent        float32 // Размер комиссии за эквайринг без НДС, %
	AcquiringBank           string  // Банк экварйрер
	SellerReward            float32 // Вознаграждение маркетплейса без НДС
	OfficeID                int     // ID офиса
	OfficeName              string  // Наименование офиса
	SupplierID              int     // ID поставщика
	SupplierName            string  // Наименование поставщика
	SupplierINN             string  // ИНН поставщика
	DeclarationNumber       string  // Номер таможенной декларации
	BonusTypeName           string  // Штрафы или доплаты
	StickerID               string  // Цифровое значение стикера, который клеится на товар в процессе сборки заказа по схеме "Маркетплейс"
	CountryOfSale           string  // Страна продажи
	Penalty                 float32 // Штрафы
	AdditionalPayment       float32 // Доплаты
	RebillLogisticCost      float32 // Стоимость возмещения издержек перевозки
	RebillLogisticOrg       string  // Организация перевозки
	KIZ                     string  // Код маркировки товара
	StorageFee              float32 // Стоимость хранения
	Deduction               float32 // Прочие удержания и выплаты
	Acceptance              float32 // Стоимость платной приемки
}
