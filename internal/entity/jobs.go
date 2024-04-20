package entity

import "time"

type Jobs struct {
  ID int64
  Description string
  Status string
  CreatedAt time.Time
}

type EventEnum struct{
  ID int64
  EventDescription string
}
