package generic_client

import (
	"time"

	"google.golang.org/grpc"
)

const DefaultTimeout = 10 * time.Second

// Config задаёт параметры соединения с удалённым GRPC сервером
type Config struct {
	// Addr - удалённый адрес, с которым соединяемся по GRPC протоколу.
	// Обычно все сервисы слушают GRPC на 8090 порту.
	// Пример - signer.astral-edo-2:8090. Протокол указывать не надо
	Addr string `yaml:"addr" conf:"addr" validate:"required,hostname_port"`
	// Timeout - длительность установки соединения с удалённым сервером, по умолчанию, если
	// сервер не ответил в течение 10 секунд, то возвращаем ошибку
	Timeout time.Duration `yaml:"timeout" conf:"timeout"`
	// Token - JWT токен, который используется для авторизации. Если поле пустое, то JWT
	// авторизация не производится
	Token string `yaml:"token" conf:"token"`
	// TLS - включать шифрование, по умолчанию выключено
	TLS bool `yaml:"tls" conf:"tls"`
	// InsecureSkipVerify при включённом шифровании, не проверять
	// валидность сертификата удалённого сервера. На проде так лучше не делать
	// По умолчанию отмена проверки сертификата выключена, и сертификат всегда проверяется, если
	// включен TLS
	InsecureSkipVerify bool `yaml:"insecureSkipVerify" conf:"insecureSkipVerify"`
	// Block - блокирующий режим работы унарных запросов. Также влияет на стрим запросы,
	// они обрываются по контексту, по умолчанию выключено, но иногда может потребоватсья включение
	// этого режима
	Block bool `yaml:"block" conf:"block"`

	// Tracing включает трассировку с помощью OpenTelemetry путём передачи заголовков с TraceID
	// Но, к солажению, данная функция требует Go версии от 1.18+ так как там используются генерики.
	// По умолчанию выключено
	Tracing bool `yaml:"tracing" conf:"tracing"`

	// ExtraDialOptions - дополнительные опции соединения, например, лимиты на размер запросов\ответов и т.д.
	ExtraDialOptions []grpc.DialOption `yaml:"-"`
}
