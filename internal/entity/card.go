package entity

import "time"

type Card struct {
	ID              int64             // id в бд
	ExternalID      int64             // id в магазине
	VendorID        string            // Код номенклатура
	VendorCode      string            // Артикул
	Title           string            // Наименование
	Description     string            // Описание номенклатуры
	CreatedAt       time.Time         // Дата создания
	UpdatedAt       time.Time         // Дана обновления
	Brand           Brand             // Бренд
	Dimensions      Dimensions        // Размеры упаковки
	Characteristics []*CardCharacteristic // Характеристики номенклатуры
	Categories      []*Category       // Категории номенклатуры
	Sizes           []*Size           // Размеры
	Barcodes        []*Barcode        // Штрихкоды
	MediaFile       []*MediaFile      // Фоточки
}

type Characteristic struct {
	ID    int64
	Title string // Наименование характеристики
}

type CardCharacteristic struct {
	ID               int64
	Value            []string // Значение характеристики
	Title            string   // Текстовое значение характеристики
	CharacteristicID int64    // Ссылка на справочник характеристики
	CardID           int64    // Номенклатура владелец
}

type Brand struct {
	ID         int64
	ExternalID int64  // id в магазине
	Title      string // Наименование  бренда

	SellerID int64 // Продавец
}

type Category struct {
	ID         int64
	ExternalID int64 // id в магазине
	Title      string

	SellerID int64
}

type Size struct {
	ID         int64
	ExternalID int64  // id в магазине
	TechSize   string // Технический размер (64-127)
	Title      string // Произвольное написание

	CardID  int64
	PriceID int64 // Цена размера
}

type Barcode struct {
	Barcode    string // Штрихкод
	ExternalID int64  // id в магазине

	SizeID   int64
	SellerID int64
}

type Dimensions struct {
	Width  int
	Height int
	Length int

	CardID int64
}
