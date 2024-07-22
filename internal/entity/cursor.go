package entity

import "time"

type Cursor struct {
	Position  int
	UpdatedAt time.Time
}
