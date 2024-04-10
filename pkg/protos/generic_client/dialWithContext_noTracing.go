package generic_client

import (
	"context"
	"crypto/tls"
	"strings"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// DialWithContextNoTracing устанавливает соединение с удалённым GRPC сервером без использования трейсинга
func DialWithContextNoTracing(ctx context.Context, cfg Config) (*grpc.ClientConn, error) {
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
		log.Debug().Msgf("Включёна имитация трейсинга для GRPC соединения с %s", cfg.Addr)
	}
	if cfg.Block {
		log.Debug().Msgf("Соединяемся с %s используя блокирующее соединение.", cfg.Addr)
		dialOpts = append(dialOpts, grpc.WithBlock())
	}
	// добавляем дополнительные опции соединения
	dialOpts = append(dialOpts, cfg.ExtraDialOptions...)
	return grpc.DialContext(cwt, cfg.Addr, dialOpts...)
}
