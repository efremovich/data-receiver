package generic_client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// Полезная статья
// https://grpc.github.io/grpc/core/md_doc_keepalive.html

// Ping проверяет работоспособность соединения на основе получения состояния пула соединений
// https://track.astral.ru/soft/browse/EDO-9452
func Ping(ctx context.Context, conn *grpc.ClientConn) (err error) {
	state := conn.GetState()
	switch state {
	case connectivity.Ready, connectivity.Idle:
		return nil
	default:
		return fmt.Errorf("%s: некорректное состояние grpc соединения", state)
	}
}
