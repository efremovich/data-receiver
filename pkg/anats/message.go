package astral_nats

import (
	"time"
)

type Message interface {
	GetID() string
	GetData() []byte
	GetCreated() time.Time
	GetRetryCount() int
	IsLastAttempt() bool
}

var _ Message = MessageImpl{}

type MessageImpl struct {
	Id           string
	data         []byte
	created      time.Time
	retryCounter int
	maxRetry     int
}

func (m MessageImpl) GetID() string {
	return m.Id
}

func (m MessageImpl) GetData() []byte {
	return m.data
}

func (m MessageImpl) GetCreated() time.Time {
	return m.created
}

func (m MessageImpl) GetRetryCount() int {
	return m.retryCounter
}

func (m MessageImpl) IsLastAttempt() bool {
	return m.retryCounter == m.maxRetry
}
