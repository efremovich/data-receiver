package config

import "time"

type Config struct {
	ServiceName        string   `env:"SERVICE_NAME, default=receiver"`
	PGWriterConn       string   `env:"POSTGRES_WRITER_CONN"`
	PGReaderConn       string   `env:"POSTGRES_READER_CONN"`
	LogLevel           int      `env:"LOG_LEVEL, default=-4"` // debug = -4, info = 0, warn = 4
	BrokerConsumerURL  []string `env:"BROKER_CONSUMER_URL" validate:"required"`
	BrokerPublisherURL []string `env:"BROKER_PUBLISHER_URL" validate:"required"`
	Gateway            Gateway  `env:", prefix=GATEWAY_"`
	Seller             Sellers  `env:", prefix=SELLER_"`
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

// Настройки очереди.
type Queue struct {
	Workers           int `env:"WORKERS, default=1"`             // Количество потоков получения сообщений из очереди.
	MaxDeliver        int `env:"MAX_DELIVER, default=1"`         // Максимальное количество попыток получить сообщение.
	NakTimeoutSeconds int `env:"NAK_TIMEOUT_SECONDS, default=2"` // Время через которое будет повторяться попытка получить сообщение.
	AckWaitSeconds    int `env:"ACK_WAIT_SECONDS, default=60"`   // Время на обработку полученного сообщения.
	MaxAckPending     int `env:"MAX_ACK_PENDING, default=10000"` // Максимальное количество сообщений, которые могут быть ожидающими подтверждения.
}

// Конфигурация для создания api клиентов для получения данных.
type Sellers struct {
	WB    SellerWB    `env:", prefix=WB_"`
	OZON  SellerOZON  `env:", prefix=OZON_"`
	OdinC SellerOdinC `env:", prefix=1C_"`
}

type SellerWB struct {
	URLMarketPlace        string   `env:"URL_MP"`
	URLContent            string   `env:"URL_CONTENT"`
	URL                   string   `env:"URL"`
	URLStat               string   `env:"URL_STAT"`
	Token                 []string `env:"TOKEN"`
	TokenStat             []string `env:"TOKEN_STAT"`
	ProcessTimeoutSeconds int      `env:"TIMEOUT, default=15"`
	Schedule              Schedule `env:", prefix=SCHEDULE_"`
}

type SellerOZON struct {
	URL                   string   `env:"URL"`
	APIKey                []string `env:"APIKEY"`
	ClientID              []string `env:"CLIENTID"`
	ProcessTimeoutSeconds int      `env:"TIMEOUT, default=15"`
	Schedule              Schedule `env:", prefix=SCHEDULE_"`
}

type SellerOdinC struct {
	URL                   string   `env:"URL"`
	Login                 string   `env:"LOGIN"`
	Password              string   `env:"PASSWORD"`
	ProcessTimeoutSeconds int      `env:"TIMEOUT, default=15"`
	Schedule              Schedule `env:", prefix=SCHEDULE_"`
}

type Schedule struct {
	StartTime string        `env:"STARTTIME"` // Время начала первого запуска
	Interval  time.Duration `env:"INTERVAL"`  // Интервал в секундах между запусками
}
