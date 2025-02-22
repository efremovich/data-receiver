package entity

import "time"

type MarketPlace struct {
	ID         int64
	Title      string // Наименование продавца
	IsEnabled  bool   // Признак активности
	ExternalID string // Внешний ID
	Type       MarketplaceType
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// MarketplaceType представляет тип маркетплейса
type MarketplaceType string

const (
	Wildberries MarketplaceType = "wb"
	Ozon        MarketplaceType = "ozon"
	Yandex      MarketplaceType = "yandex"
	Sber        MarketplaceType = "sber"
	OdinAss     MarketplaceType = "1c"
)
