//go:build go1.19 || go1.20

// благодаря тегу build код будет собираться только актуальными версиями компилятора, так как
// система opentelemetry требует генерироков.

// про теги - см
// https://stackoverflow.com/questions/38439275/golang-build-tags-for-a-particular-go-version-possible
// https://www.programming-books.io/essential/go/conditional-compilation-with-build-tags-d1980344374d45c082c914c2aafa50cf

package generic_client

import (
	"context"
	"crypto/tls"
	"strings"

	"github.com/rs/zerolog/log"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// DialWithContext устанавливает соединение с удалённым GRPC сервером
func DialWithContext(ctx context.Context, cfg Config) (*grpc.ClientConn, error) {
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultTimeout
	}
	cwt, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()
	log.Debug().Msgf("Соединяемся с %s через GRPC", cfg.Addr)
	var dialOpts []grpc.DialOption

	if cfg.TLS {
		// включаем шифрование
		if cfg.InsecureSkipVerify {
			log.Warn().Msgf("Шифруем соединение, но НЕ проверяем сертификат у %s", cfg.Addr)
		} else {
			log.Debug().Msgf("Шифруем соединение и проверяем сертификат у %s", cfg.Addr)
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			ServerName:         strings.Split(cfg.Addr, ":")[0], // нужно для SNI
			InsecureSkipVerify: cfg.InsecureSkipVerify,          // если у удалённого сервера невалидный сертификат
		})))
	} else {
		log.Warn().Msgf("Внимание, используем соединение без шифрования с %s", cfg.Addr)
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if cfg.Token != "" {
		// добавляем авторизацию по JWT токену
		log.Debug().Msgf("Используем токен для соединения с %s", cfg.Addr)
		dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(tokenAuth{
			Token:  cfg.Token,
			Secure: cfg.TLS,
		}))
	} else {
		log.Debug().Msgf("Соединяемся с %s без авторизации.", cfg.Addr)
	}
	if cfg.Tracing {
		log.Debug().Msgf("Включён трейсинг для GRPC соединения с %s", cfg.Addr)
		// OpenTelemetry требует поддержки генериков и Go версии 1.19
		dialOpts = append(dialOpts, grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))
		dialOpts = append(dialOpts, grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
	}
	if cfg.Block {
		log.Debug().Msgf("Соединяемся с %s используя блокирующее соединение.", cfg.Addr)
		dialOpts = append(dialOpts, grpc.WithBlock())
	}
	// добавляем дополнительные опции соединения
	dialOpts = append(dialOpts, cfg.ExtraDialOptions...)
	return grpc.DialContext(cwt, cfg.Addr, dialOpts...)
}
