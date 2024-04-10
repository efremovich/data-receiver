package entity

import "time"

type TpEvent struct {
	TpID        int64           `db:"tp_id"`
	CreatedAt   time.Time       `db:"created_at"`
	EventType   TpEventTypeEnum `db:"event_type"`
	Description string          `db:"description"`
}

type TpEventTypeEnum string

var (
	CreatedTpEventType TpEventTypeEnum = "CREATED"
	SuccessEventType   TpEventTypeEnum = "SUCCESS"
	GotAgainEventType  TpEventTypeEnum = "GOT_AGAIN" // Прислали повторно.
	ReprocessEventType TpEventTypeEnum = "REPROCESS" // Поступил ручной запрос на переобработку.
	ErrorEventType     TpEventTypeEnum = "ERROR"
	SendTaskNext       TpEventTypeEnum = "SEND_TASK_NEXT" // Отправили задачу на обработку дальше.
)
