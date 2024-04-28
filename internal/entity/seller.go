package entity

import "time"

type Seller struct {
	ID        int64
	Title     string // Наименование продавца
	IsEnable  bool   // Признак активности
	ExtID     int    // Внешний ID
	CreatedAt time.Time
	UpdatedAt time.Time
}
