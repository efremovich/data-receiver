package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/efremovich/data-receiver/config"
	anats "github.com/efremovich/data-receiver/pkg/anats"
)

type NATS interface {
	SendTaskMessageNext(ctx context.Context, msg Task) error
	Ping() error
}

type natsImpl struct {
	anats anats.NatsClient
}

func NewNats(ctx context.Context, c config.NATS, updateStream bool) (NATS, error) {
	cfg := anats.NatsClientConfig{
		Urls:               []string{c.URL},
		StreamName:         ReceiverStreamName,
		Subjects:           []string{ReceiverSubjectNormalPriority},
		CreateUpdateStream: updateStream,
	}

	cl, err := anats.NewNatsClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &natsImpl{anats: cl}, nil
}

func (nw *natsImpl) SendTaskMessageNext(ctx context.Context, msg Task) error {
	jsonB, err := json.Marshal(&msg)
	if err != nil {
		return fmt.Errorf("ошибка маршалла в JSON: %s", err.Error())
	}

	return nw.anats.PublishMessage(ctx, ReceiverSubjectNormalPriority, jsonB)
}

func (nw *natsImpl) Ping() error {
	return nw.anats.Ping()
}
