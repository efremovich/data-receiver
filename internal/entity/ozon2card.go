package entity

import "time"

type Ozon2Card struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time

	Card Card
}
