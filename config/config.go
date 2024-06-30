package config

import (
	"fmt"

	"github.com/efremovich/data-receiver/pkg/aconf/v3"
)

type Config struct {
	ServiceName  string  `env:"SERVICE_NAME, default=receiver"`
	PGWriterConn string  `env:"POSTGRES_WRITER_CONN"`
	PGReaderConn string  `env:"POSTGRES_READER_CONN"`
	LogLevel     int     `env:"LOG_LEVEL, default=-4"` // debug = -4, info = 0, warn = 4
	Gateway      Gateway `env:", prefix=GATEWAY_"`
	Nats         NATS    `env:", prefix=NATS_"`
	Seller       Seller  `env:", prefix=SELLER_"`
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

type Queue struct {
	Workers               int `env:"WORKERS"`
	Repeats               int `env:"REPEATS"`
	NakTimeoutSeconds     int `env:"NAK_TIMEOUT_SECONDS"`
	ProcessTimeoutSeconds int `env:"PROCESS_TIMEOUT_SECONDS"`
	MaxAckPending         int `env:"MAX_ACK_PENDING"`
}

type Seller struct {
	URL                   string `env:"URL"`
	Token                 string `env:"TOKEN"`
	ProcessTimeoutSeconds int    `env:"TIMEOUT, default=10"`
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
