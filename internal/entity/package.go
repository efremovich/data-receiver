package entity

import (
	"fmt"
	"strings"
	"time"
)

// Описание пакета по которому создается пакет.
type PackageDescription struct {
	Cursor      int               // Курсор пакета.
	Limit       int               // Количество записей в запросе
	UpdatedAt   time.Time        // Дата обновления.
	PackageType PackageType       // Тип пакета.
	Seller      string            // Код продавца (wb, ozon, yandex, 1с)
	Query       map[string]string // Параметры запроса
}

// Тип пакета.
type PackageType string

const (
	PackageTypeCard  = PackageType("CARD")  // Пакет с товарным карточками.
	PackageTypeOrder = PackageType("ORDER") // Пакет с заказами.
	PackageTypeSale  = PackageType("SALE")  // Пакет с продажами.
	PackageTypeStock = PackageType("STOCK") // Пакет с остатками.
)

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
