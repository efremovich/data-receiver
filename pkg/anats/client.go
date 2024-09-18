package astral_nats

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsClient interface {
	PublishMessage(ctx context.Context, subject string, data []byte) error
	// публикация сообщения с дедубликацией по ID. окно дедубликации - 2 минуты
	PublishMessageDupe(ctx context.Context, subject string, data []byte, msgID string) error
	Subscribe(ctx context.Context, consumerName string, subject string, handler HandlerFunc, opt SubscribeOptions) error
	GetConsumersInfo(ctx context.Context) ([]ConsumerInfo, error)
	Ping() error
}

type natsClientImpl struct {
	conn   *nats.Conn
	js     jetstream.JetStream
	cons   []jetstream.Consumer
	stream jetstream.Stream
	cfg    NatsClientConfig
}

func NewNatsClient(ctx context.Context, c NatsClientConfig) (NatsClient, error) {
	opts := nats.Options{
		Servers:        c.Urls,
		MaxReconnect:   -1,
		ReconnectWait:  time.Second * 10,
		AllowReconnect: true,
	}

	conn, err := opts.Connect()
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	var stream jetstream.Stream
	if c.CreateUpdateStream {
		cfg := jetstream.StreamConfig{
			Name:       c.StreamName,
			Retention:  jetstream.WorkQueuePolicy,
			Subjects:   c.Subjects,
			Storage:    jetstream.MemoryStorage,
			Duplicates: time.Minute * 2,
		}

		stream, err = js.CreateOrUpdateStream(ctx, cfg)
		if err != nil {
			return nil, err
		}
	} else {
		stream, err = js.Stream(ctx, c.StreamName)
		if err != nil {
			return nil, err
		}
	}

	return &natsClientImpl{
		conn:   conn,
		stream: stream,
		js:     js,
		cfg:    c,
	}, nil
}

func (nw *natsClientImpl) PublishMessageDupe(ctx context.Context, subject string, data []byte, msgID string) error {
	id := fmt.Sprintf("%s_%s", subject, msgID)
	return nw.publishMessage(ctx, subject, publishMessage{Data: data, ID: id})
}

func (nw *natsClientImpl) PublishMessage(ctx context.Context, subject string, data []byte) error {
	return nw.publishMessage(ctx, subject, publishMessage{Data: data})
}

func (nw *natsClientImpl) publishMessage(ctx context.Context, subject string, pMsg publishMessage) error {
	if nw.conn == nil {
		return fmt.Errorf("нет активного подключения к NATS")
	}

	if !nw.steamHasSubject(subject) {
		return fmt.Errorf("subject %s не был указан при создании клиента", subject)
	}

	msg := &nats.Msg{
		Subject: subject,
		Data:    pMsg.Data,
		Header:  make(nats.Header),
	}
	if pMsg.ID != "" {
		msg.Header.Add(string(headerIdName), pMsg.ID)
	}

	traceID := ctx.Value(traceIdHeaderName)
	if traceID != nil {
		traceIdString, ok := traceID.(string)
		if ok {
			msg.Header.Add(string(traceIdHeaderName), traceIdString)
		}
	}

	_, err := nw.js.PublishMsg(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}

func (nw *natsClientImpl) Subscribe(ctx context.Context, consumerName string, subject string, handler HandlerFunc, opt SubscribeOptions) error {
	if nw.conn == nil {
		return fmt.Errorf("нет активного подключения к NATS")
	}

	if !nw.steamHasSubject(subject) {
		return fmt.Errorf("subject %s не был указан при создании клиента", subject)
	}

	consumerCfg := jetstream.ConsumerConfig{
		Durable:       consumerName,
		FilterSubject: subject,
		MaxDeliver:    opt.MaxDeliver,
		MaxAckPending: opt.MaxAckPending,
		AckWait:       time.Duration(opt.AckWaitSeconds) * time.Second,
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverAllPolicy,
	}

	cons, err := nw.stream.CreateOrUpdateConsumer(ctx, consumerCfg)
	if err != nil {
		return fmt.Errorf("ошибка создания консьюмера: %s", err.Error())
	}

	nw.cons = append(nw.cons, cons)

	for i := 0; i < opt.Workers; i++ {
		jscontext, err := cons.Consume(func(msg jetstream.Msg) {
			msgID := msg.Headers().Get(string(headerIdName))
			if msgID == "" {
				msgID = "without_id"
			}

			msgWrap := MessageImpl{
				data:     msg.Data(),
				maxRetry: opt.MaxDeliver,
				ID:       msgID,
			}

			metadata, err := msg.Metadata()
			if err == nil {
				msgWrap.retryCounter = int(metadata.NumDelivered)
				msgWrap.created = metadata.Timestamp
			}

			fmt.Printf("внутри либы натса: начинается обработка сообщения ID %s, time: %s. попытка %d\n", msgID, time.Now().UTC().Format(time.RFC3339Nano), msgWrap.retryCounter)

			traceID := msg.Headers().Get(string(traceIdHeaderName))
			if traceID == "" {
				traceID = strings.ReplaceAll(uuid.NewString(), "-", "")
			}

			newContext := context.WithValue(ctx, traceIdHeaderName, traceID)

			res := handler(newContext, msgWrap)
			fmt.Printf("внутри либы натса: завершается обработка сообщения ID %s, res: %d, time: %s\n", msgID, res, time.Now().UTC().Format(time.RFC3339Nano))

			switch res {
			case MessageResultEnumSuccess:
				err = msg.DoubleAck(ctx)
				if err != nil {
					fmt.Printf("внутри либы натса: ошибка DoubleAck. ID сообщения: %s, err: %s\n", msgID, err.Error())
				}
			case MessageResultEnumFatalError:
				err = msg.Term()
				if err != nil {
					fmt.Printf("внутри либы натса: ошибка Term. ID сообщения %s, err: %s\n", msgID, err.Error())
				}
			case MessageResultEnumTempError:
				// чтобы сообщение не висело в очереди
				if msgWrap.IsLastAttempt() {
					err = msg.Term()
					if err != nil {
						fmt.Printf("внутри либы натса: ошибка Term. ID сообщения %s, err: %s\n", msgID, err.Error())
					}
				} else {
					err = msg.NakWithDelay(time.Duration(opt.NakTimeoutSeconds) * time.Second)
					if err != nil {
						fmt.Printf("внутри либы натса: ошибка NakWithDelay. ID сообщения %s, err: %s\n", msgID, err.Error())
					}
				}
			}
		})
		if err != nil {
			return fmt.Errorf("ошибка создания консьюмера: %s", err.Error())
		}

		go func() {
			<-ctx.Done()
			jscontext.Stop()
		}()
	}

	return nil
}

func (nw *natsClientImpl) Ping() error {
	if nw.conn == nil {
		return fmt.Errorf("нет активного подключения к NATS")
	}

	err := nw.conn.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (nw *natsClientImpl) GetConsumersInfo(ctx context.Context) ([]ConsumerInfo, error) {
	var res []ConsumerInfo

	for _, c := range nw.cons {
		info, err := c.Info(ctx)
		if err != nil {
			return nil, err
		}

		res = append(res, ConsumerInfo{
			Name:        info.Name,
			Subject:     info.Config.FilterSubject,
			Subjects:    info.Config.FilterSubjects,
			CountQueue:  int(info.NumPending),
			CountInWork: info.NumAckPending,
		})
	}

	return res, nil
}

func (nw *natsClientImpl) steamHasSubject(s string) bool {
	subjects := nw.stream.CachedInfo().Config.Subjects

	for _, subject := range subjects {
		if subject == s {
			return true
		}
	}

	return false
}
