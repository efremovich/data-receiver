package entity

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Описание пакета по которому создается пакет.
type PackageDescription struct {
	Cursor      string            `json:"cursor"` // Курсор пакета.
	LastID      string            `json:"last_id"`
	Limit       int               `json:"limit"`        // Количество записей в запросе
	UpdatedAt   time.Time         `json:"updated_at"`   // Дата обновления.
	PackageType PackageType       `json:"package_type"` // Тип пакета.
	Seller      MarketplaceType   `json:"seller"`       // Код продавца (wb, ozon, yandex, 1с)
	Query       map[string]string `json:"query"`        // Параметры запроса
	Delay       int               `json:"delay"`        // Задержка перед следующей загрузкой
}

// Тип пакета.
type PackageType string

const (
	PackageTypeCard        = PackageType("CARD")       // Пакет с товарным карточками.
	PackageTypeOrder       = PackageType("ORDER")      // Пакет с заказами.
	PackageTypeSale        = PackageType("SALE")       // Пакет с продажами.
	PackageTypeStock       = PackageType("STOCK")      // Пакет с остатками.
	PackageTypeSaleReports = PackageType("SALEREPORT") // Пакет с отчетом по продажами.
)

func (p PackageDescription) GetCursor() string {
	if p.Cursor == "" {
		return "0"
	}
	return p.Cursor
}
func StringToPackageType(s string) (PackageType, error) {
	s = strings.ToUpper(s)

	switch s {
	case "CARD":
		return PackageTypeCard, nil
	case "ORDER":
		return PackageTypeOrder, nil
	case "SALE":
		return PackageTypeSale, nil
	case "STOCK":
		return PackageTypeStock, nil
	case "SALEREPORT":
		return PackageTypeStock, nil
	default:
		return "", fmt.Errorf("неизвестный тип пакета: %s", s)
	}
}

// Статус пакета.
type PackageStatus string

const (
	PackageStatusCreated = PackageStatus("CREATED") // Пакет создан.
	PackageStatusSuccess = PackageStatus("SUCCESS") // Пакет успешно обработан.
	PackageStatusFailed  = PackageStatus("FAILED")  // Обработка пакета завершена с ошибкой.
)

func (p *PackageDescription) ConvertCursorToInt() int {
	v, err := strconv.Atoi(p.Cursor)
	if err != nil {
		return 0
	}

	return v
}
