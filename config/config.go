package config

import (
	"fmt"

	"github.com/efremovich/data-receiver/pkg/aconf/v3"
)

type Config struct {
	ServiceName    string `env:"SERVICE_NAME, default=receiver"`
	SelfOperatorID string `env:"SELFID"`
	NATS           `env:", prefix=NATS_"`
	PGWriterConn   string      `env:"POSTGRES_WRITER_CONN"`
	PGReaderConn   string      `env:"POSTGRES_READER_CONN"`
	LogLevel       int         `env:"LOG_LEVEL, default=-4"` // debug = -4, info = 0, warn = 4
	OperatorAPI    OperatorAPI `env:", prefix=OPERATOR_"`
	Storage        Storage     `env:", prefix=STORAGE_"`
	Gateway        Gateway     `env:", prefix=GATEWAY_"`
	Packer         Packer      `env:", prefix=PACKER_"`
}

type Storage struct {
	URL                string `env:"URL"`
	Token              string `env:"TOKEN"`
	UseTLS             bool   `env:"USE_TLS, default=false"`
	TimeoutSeconds     int    `env:"TIMEOUT_SECONDS, default=5"`
	InsecureSkipVerify bool   `env:"INSECURE_SKIP_VERIFY, default=true"`

	// Для юнит-тестов и локального запуска.
	UseMockStorage bool `env:"USE_MOCK, default=false"`
}

type OperatorAPI struct {
	BaseURL        string `env:"URL"`
	TimeoutSeconds int    `env:"TIMEOUT, default=10"`
	Login          string `env:"LOGIN"`
	Password       string `env:"PASSWORD"`
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

type Packer struct {
	CertData []byte `env:"CERT"`
	KeyData  []byte `env:"KEY"`
}

type NATS struct {
	URL             string            `env:"URL" validate:"required"`
	ServiceSubjects map[string]string `env:"SERVICE_SUBJECTS"`

	QueueReceiver Queue `env:", prefix=QUEUE_RECEIVER_"`
	QueueSender   Queue `env:", prefix=QUEUE_SENDER_"`
	QueueMarking  Queue `env:", prefix=QUEUE_MARKING_"`
}

type Queue struct {
	Workers               int `env:"WORKERS"`
	Repeats               int `env:"REPEATS"`
	NakTimeoutSeconds     int `env:"NAK_TIMEOUT_SECONDS"`
	ProcessTimeoutSeconds int `env:"PROCESS_TIMEOUT_SECONDS"`
	MaxAckPending         int `env:"MAX_ACK_PENDING"`
}

func NewConfig(envPath string) (*Config, error) {
	cfg := &Config{}

	if envPath != "" {
		err := aconf.PreloadEnvsFile(envPath)
		if err != nil {
			return nil, fmt.Errorf("ошибка при обработке файла env: %s", err.Error())
		}
	}

	if err := aconf.Load(cfg); err != nil {
		return nil, fmt.Errorf("aconf.Load failed: %s", err.Error())
	}

	return cfg, nil
}
