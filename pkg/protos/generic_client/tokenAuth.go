package generic_client

import (
	"context"
)

// tokenAuth реализует интерфейс https://pkg.go.dev/google.golang.org/grpc/credentials#PerRPCCredentials
type tokenAuth struct {
	Token  string
	Secure bool
}

// GetRequestMetadata вызывается при каждом созданном запросе, чтобы добавить в него JWT токен из tokenAuth
func (t tokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": "Bearer " + t.Token,
	}, nil
}

// RequireTransportSecurity вызывается при каждом созданном GPRC запросе, чтобы проверить, является ли соединение зашифрованным
func (t tokenAuth) RequireTransportSecurity() bool {
	return t.Secure
}
