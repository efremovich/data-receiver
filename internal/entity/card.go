package entity

import "time"

type Card struct {
	ID              int64             // id в бд
	VendorID        string            // Код номенклатура
	VendorCode      string            // Артикул
	Title           string            // Наименование
	Description     string            // Описание номенклатуры
	CreatedAt       time.Time         // Дата создания
	UpdatedAt       time.Time         // Дана обновления
	Brand           Brand             // Бренд
	Characteristics []*Characteristic // Характеристики номенклатуры
	Categories      []Category        // Категории номенклатуры
	Sizes           []Size            // Размеры
}

type Characteristic struct {
	ID    int64
	Title string   // Наименование характеристики
	Value []string // Значение характеристики

	CardID int64 // Номенклатура владелец
}

type Brand struct {
	ID   int64
	Name string // Наименование  бренда

	SellerID int64 // Продавец
}

type Category struct {
	ID    int64
	Title string

	CardID   int64
	SellerID int64
}

type Size struct {
	ID       int64
	TechSize string // Технический размер (64-127)
	Title    string // Произвольное написание

	CardID  int64
	PriceID int64 // Цена размера
}

type Barcode struct {
	Barcode string // Штрихкод

	SizeID   int64
	SellerID int64
}
