package astral_nats

import (
	"context"
)

type headerName string

const (
	traceIdHeaderName headerName = "trace_id"
	headerIdName      headerName = "Nats-Msg-Id"
)

type MessageResultEnum int8

const (
	MessageResultEnumSuccess    MessageResultEnum = 1
	MessageResultEnumTempError  MessageResultEnum = 2
	MessageResultEnumFatalError MessageResultEnum = 3
)

type HandlerFunc func(ctx context.Context, m Message) MessageResultEnum

type SubscribeOptions struct {
	Workers           int
	MaxDeliver        int
	NakTimeoutSeconds int
	AckWaitSeconds    int
	MaxAckPending     int
}

type NatsClientConfig struct {
	Urls               []string
	StreamName         string
	Subjects           []string
	CreateUpdateStream bool // если false - не будет обновлять стрим, а просто попробует его получить.
}

type ConsumerInfo struct {
	Name        string   `json:"name"`
	Subject     string   `json:"subject,omitempty"`
	Subjects    []string `json:"subjects,omitempty"`
	CountQueue  int      `json:"count_queue"`
	CountInWork int      `json:"count_in_work"`
}

type publishMessage struct {
	Data []byte
	ID   string
}
