package entity

import "time"

type Wb2Card struct {
	ID        int64
	NMID      int64  // Артикул WB
	KTID      int    // Идентификатор карточки товара
	NMUUID    string // Внутренний идентификатор
	CreatedAt time.Time
	UpdatedAt time.Time

	CardID int64
}
