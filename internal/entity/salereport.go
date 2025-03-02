package entity

import "time"

type SaleReport struct {
	ID         int64
	ExternalID string    // Уникальный идентификатор заказа.
	UpdatedAt  time.Time // Дата обновления данных

	Quantity               float32   // Количество
	RetailPrice            float32   // Цена розничная
	ReturnAmoun            float32   // Количество возвратов
	SalePercent            int       // Процент скидки
	CommissionPercent      float32   // Процент комиссии
	RetailPriceWithdiscRub float32   // Цена розничная с учетом скидок в рублях.
	DeliveryAmount         float32   // Количество доставок
	ReturnAmount           float32   // Количество возвратов
	DeliveryCost           float32   // Стоимость доставки
	PvzReward              float32   // Возмещение за выдачу на ПВЗ
	SellerReward           float32   // Возмещение марекетплейса без НДС
	SellerRewardWithNds    float32   // Возмещение марекетплейса с НДС
	DateFrom               time.Time // Дата начала отчета
	DateTo                 time.Time // Дата окончания отчета
	CreateReportDate       time.Time // Дата создания отчета
	OrderDate              time.Time // Дата заказа
	SaleDate               time.Time // Дата продажи
	TransactionDate        time.Time // Дата транзакции

	SAName            string  // Артикул продавца TODO Проверить нужен ли.
	BonusTypeName     string  // Штрафы или доплаты
	Penalty           float32 // Штрафы
	AdditionalPayment float32 // Доплаты
	AcquiringFee      float32 // Возмещение издержек по эквайрингу. Издержки WB за услуги эквайринга: вычитаются из вознаграждения WB и не влияют на доход продавца
	AcquiringPercent  float32 // Размер комиссии за эквайринг без НДС, %
	AcquiringBank     string  // Банк экварйрер
	DocType           string  // Тип документа
	SupplierOperName  string  // Обоснование оплаты

	SiteCountry string  // Страна сайта
	KIZ         string  // Код маркировки товара
	StorageFee  float32 // Стоимость хранения
	Deduction   float32 // Прочие удержания и выплаты
	Acceptance  float32 // Стоимость платной приемки

	Pvz       *Pvz
	Barcode   *Barcode
	Size      *Size
	Card      *Card
	Order     *Order
	Warehouse *Warehouse
	Seller    *MarketPlace
}

// Пункт выдачи заказов
type Pvz struct {
	ID           int64
	OfficeName   string // Наименование поставщика
	OfficeID     int    // ID офиса
	SupplierName string // Наименование поставщика
	SupplierID   int    // ID поставщика
	SupplierINN  string // ИНН поставщика

}
