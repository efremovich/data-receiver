package entity

import "time"

type Wb2Card struct{
  NMID int64
  KTID int
  NMUUID string
  CreatedAt time.Time
  UpdatedAt time.Time

  CardID int64
}

