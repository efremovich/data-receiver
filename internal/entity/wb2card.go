package entity

import "time"

type Seller2Card struct {
	ID         int64
	ExternalID int64  // Артикул WB
	KTID       int    // Идентификатор карточки товара
	NMUUID     string // Внутренний идентификатор
	CreatedAt  time.Time
	UpdatedAt  time.Time

	CardID   int64
	SellerID int64
}
