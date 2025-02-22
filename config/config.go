package config

import (
	"os"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Config struct {
	ServiceName           string                 `env:"SERVICE_NAME, default=data-receiver"`
	PGWriterConn          string                 `env:"POSTGRES_WRITER_CONN"`
	PGReaderConn          string                 `env:"POSTGRES_READER_CONN"`
	LogLevel              int                    `env:"LOG_LEVEL, default=-4"` // debug = -4, info = 0, warn = 4
	BrokerConsumerURL     []string               `env:"BROKER_CONSUMER_URL" validate:"required"`
	BrokerPublisherURL    []string               `env:"BROKER_PUBLISHER_URL" validate:"required"`
	Gateway               Gateway                `env:", prefix=GATEWAY_"`
	MarketPlaces          map[string]MarketPlace `env:"MARKETPLACE_"`
	Queue                 Queue                  `env:", prefix=QUEUE_"`
	ProcessTimeoutSeconds int                    `env:"TIMEOUT, default=600"`
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

type MarketPlace struct {
	ID   string `env:"ID"`   // Для вставки в базу
	Name string `env:"NAME"` // Наименование маркетплейса
	// В случае если подключение к сервису требует пару логин:пароль, ключ:токен, то записываем через запятую.
	Token string `env:"TOKEN"` // API ключ или id:token или логин:пароль
	Type  string `env:"TYPE"`
}

func (c *Config) FillMarketPlaceMap() {
	prefix := "MARKETPLACE_"

	c.MarketPlaces = make(map[string]MarketPlace)

	for k, v := range getEnvMap(prefix) {
		parts := strings.SplitN(k[len(prefix):], "_", 2)
		if len(parts) != 2 {
			continue
		}

		marketPlaceID := parts[0]
		fieldName := parts[1]

		if _, exists := c.MarketPlaces[marketPlaceID]; !exists {
			c.MarketPlaces[marketPlaceID] = MarketPlace{}
		}

		storage := c.MarketPlaces[marketPlaceID]
		setField(&storage, fieldName, v)
		c.MarketPlaces[marketPlaceID] = storage
	}
}

func getEnvMap(prefix string) map[string]string {
	envMap := make(map[string]string)

	for _, e := range os.Environ() {
		if strings.HasPrefix(e, prefix) {
			pair := strings.SplitN(e, "=", 2)
			envMap[pair[0]] = pair[1]
		}
	}
	return envMap
}

func setField(config *MarketPlace, fieldName, value string) {
	t := reflect.TypeOf(config).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("env")

		if strings.TrimSpace(tag) == strings.ToUpper(fieldName) {
			f := reflect.ValueOf(config).Elem()
			fieldName = fieldToVarName(fieldName)
			ff := f.FieldByName(fieldName)

			if ff.IsValid() && ff.CanSet() {
				switch ff.Kind() {
				case reflect.String:
					ff.SetString(value)
				case reflect.Int:
					if i, err := strconv.Atoi(value); err == nil {
						ff.SetInt(int64(i))
					}
				case reflect.Array:
				case reflect.Bool:
				case reflect.Chan:
				case reflect.Complex128:
				case reflect.Complex64:
				case reflect.Float32:
				case reflect.Float64:
				case reflect.Func:
				case reflect.Int16:
				case reflect.Int32:
				case reflect.Int64:
				case reflect.Int8:
				case reflect.Interface:
				case reflect.Invalid:
				case reflect.Map:
				case reflect.Pointer:
				case reflect.Slice:
				case reflect.Struct:
				case reflect.Uint:
				case reflect.Uint16:
				case reflect.Uint32:
				case reflect.Uint64:
				case reflect.Uint8:
				case reflect.Uintptr:
				case reflect.UnsafePointer:
				default:
					panic("unexpected reflect.Kind")
				}
			}
		}
	}
}

func fieldToVarName(fieldName string) string {
	parts := strings.Split(fieldName, "_")
	for i, part := range parts {
		if part != "ID" && part != "URL" {
			parts[i] = cases.Title(language.English).String(strings.ToLower(part))
		}
	}
	return strings.Join(parts, "")
}
