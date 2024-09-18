package brokerpublisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	anats "github.com/efremovich/data-receiver/pkg/anats"
)

// TODO возможно такое стоит вынести в конфигурации.
const (
	OutgoingStreamName            = "package-sender-stream"
	OutgoingSubjectNormalPriority = "package-sender.inbox"
)

// Брокер издатель.
type BrokerPublisher interface {
	// Проверка работы брокера.
	Ping() error
	// Отправка пакета в сервис package-sender.
	SendPackage(ctx context.Context, p *entity.PackageDescription) error
}

// Имплементация брокера издателя.
type brokerPublisherImpl struct {
	anats anats.NatsClient
}

// Инициализация брокера издателя.
func NewBrokerPublisher(ctx context.Context, urls []string, updateStream bool) (BrokerPublisher, error) {
	// TODO заменить реализацию на реализацию пакета package-sender.
	cfg := anats.NatsClientConfig{
		Urls:               urls,
		StreamName:         OutgoingStreamName,
		Subjects:           []string{OutgoingSubjectNormalPriority},
		CreateUpdateStream: updateStream,
	}

	cl, err := anats.NewNatsClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &brokerPublisherImpl{anats: cl}, nil
}

// Проверка подключение к брокеру.
func (b brokerPublisherImpl) Ping() error {
	return b.anats.Ping()
}

// Отправка события.
func (b brokerPublisherImpl) SendPackage(ctx context.Context, p *entity.PackageDescription) error {
	// Преобразование пакета приложения в сообщение.
	msg := tmpPackageSenderMsg{
		Cursor:      p.Cursor,
		UpdatedAt:   &p.UpdatedAt,
		Limit:       p.Limit,
		Seller:      p.Seller,
		PackageType: string(p.PackageType),
		Delay:       p.Delay,
	}

	// Сериализация пакета.
	msgB, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации сообщения для package-sender: %s", err.Error())
	}

	// Публикация события.
	return b.anats.PublishMessage(ctx, OutgoingSubjectNormalPriority, msgB)
}

type tmpPackageSenderMsg struct {
	Cursor      int        `json:"cursor"` // Указатель на последнюю полученную запись из внешнего источника
	UpdatedAt   *time.Time `json:"updated_at"`
	Limit       int        `json:"limit"`
	Seller      string     `json:"seller"`
	PackageType string     `json:"package_type"`
	Delay       int        `json:"delay"`
}
