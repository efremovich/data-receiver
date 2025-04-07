package entity

import "time"

type Cost struct {
	ID         int64
	ExternalID string
	CardID     int64
	Amount     float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
