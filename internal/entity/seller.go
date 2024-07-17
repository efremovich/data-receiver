package entity

import "time"

type Seller struct {
	ID         int64
	Title      string // Наименование продавца
	IsEnabled  bool   // Признак активности
	ExternalID string // Внешний ID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
