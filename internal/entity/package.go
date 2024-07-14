package entity

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Описание пакета по которому создается пакет.
type PackageDescription struct {
	PackageName string          // Наименование пакета.
  Cursor      int             // Курсор пакета.
	SendURL     string          // URL для отправки.
	PackageType PackageType     // Тип пакета.
	Description json.RawMessage // Описание пакета.
}

// Тип пакета.
type PackageType string

const (
	PackageTypeCard  = PackageType("CARD")  // Пакет с товарным карточками.
	PackageTypeOrder = PackageType("ORDER") // Пакет с заказами.
	PackageTypeSele  = PackageType("SALE")  // Пакет с продажами.
)

func StringToPackageType(s string) (PackageType, error) {
	s = strings.ToUpper(s)

	switch s {
	case "CARD":
		return PackageTypeCard, nil
	case "ORDER":
		return PackageTypeOrder, nil
	case "SALE":
		return PackageTypeSele, nil
	default:
		return "", fmt.Errorf("неизвестный тип пакета: %s", s)
	}
}

// Структура пакета.
type Package struct {
	ID        int64         // Идентификатор в БД.
	Type      PackageType   // Тип пакета.
	SendURL   string        // URL для отправки пакета.
  Cursor    int           // Курсор пакета.
	CreatedAt time.Time     // Дата создания пакета.
	Status    PackageStatus // Статус пакета.
	ErrorText string        // Текст ошибки.
	ErrorCode string        // Код ошибки.
}

// Статус пакета.
type PackageStatus string

const (
	PackageStatusCreated = PackageStatus("CREATED") // Пакет создан.
	PackageStatusSuccess = PackageStatus("SUCCESS") // Пакет успешно обработан.
	PackageStatusFailed  = PackageStatus("FAILED")  // Обработка пакета завершена с ошибкой.
)
