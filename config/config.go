package config

type Config struct {
	ServiceName        string   `env:"SERVICE_NAME, default=receiver"`
	PGWriterConn       string   `env:"POSTGRES_WRITER_CONN"`
	PGReaderConn       string   `env:"POSTGRES_READER_CONN"`
	LogLevel           int      `env:"LOG_LEVEL, default=-4"` // debug = -4, info = 0, warn = 4
	BrokerConsumerURL  []string `env:"BROKER_CONSUMER_URL" validate:"required"`
	BrokerPublisherURL []string `env:"BROKER_PUBLISHER_URL" validate:"required"`
	Gateway            Gateway  `env:", prefix=GATEWAY_"`
	Nats               NATS     `env:", prefix=NATS_"`
	Seller             Seller   `env:", prefix=SELLER_"`
	Queue              Queue    `env:", prefix=QUEUE_"`
}

type Gateway struct {
	AuthToken        string `env:"AUTH_TOKEN"`
	PathToSwaggerDir string `env:"SWAGGER_PATH, default=docs/swagger"`
	HTTP             Adr    `env:", prefix=HTTP_"`
	GRPC             Adr    `env:", prefix=GRPC_"`
}

type Adr struct {
	Host string `env:"HOST"`
	Port string `env:"PORT"`
}

type NATS struct {
	URLS string `env:"URLS" validate:"required"`
}

// Настройки очереди.
type Queue struct {
	Workers           int `env:"WORKERS, default=1"`             // Количество потоков получения сообщений из очереди.
	MaxDeliver        int `env:"MAX_DELIVER, default=1"`         // Максимальное количество попыток получить сообщение.
	NakTimeoutSeconds int `env:"NAK_TIMEOUT_SECONDS, default=2"` // Время через которое будет повторяться попытка получить сообщение.
	AckWaitSeconds    int `env:"ACK_WAIT_SECONDS, default=3"`    // Время на обработку полученного сообщения.
	MaxAckPending     int `env:"MAX_ACK_PENDING, default=10000"` // Максимальное количество сообщений, которые могут быть ожидающими подтверждения.
}

type Seller struct {
	URL                   string `env:"URL"`
	Token                 string `env:"TOKEN"`
	ProcessTimeoutSeconds int    `env:"TIMEOUT, default=10"`
}
