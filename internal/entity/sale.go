package entity

type Sale struct {
	ID       int64
	ExtID    string // Уникальный номер продажи

  Country  string // Страна
	Area     string // Область / край
	District string // Регион
	Sity     string // Город

	Price      float32 // Цена без скидки
	DiscountP  float32 // Скидка продавца
	DiscountS  float32 // Скидка на маркетпрейсе
	FinalPrice float32 // Конечная цена
	Type       string  // Тип продажи

	Barcode     string // Штрихкод
	WarehouseID int64  // ID склада
}
