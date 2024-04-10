// Package storage_client предоставляет набор функций для подключения и использования единого файлового хранилища
// Более подробное описание: https://astraltrack.atlassian.net/wiki/spaces/EDO/pages/3163259010
package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.astralnalog.ru/utils/protos/generic_client"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	pb "git.astralnalog.ru/utils/protos/storage/proto"
)

const TIMEOUT = time.Second * 3 // Максимальное время выполнения каждого запроса по умолчанию

const (
	methodGetFile      = "GetFile"
	methodSaveFile     = "SaveFile"
	methodRemoveFile   = "RemoveFile"
	methodSetFileAttrs = "SetFileAttrs"
	methodGetFileAttrs = "GetFileAttrs"
)

var (
	ErrStorageClientConn = errors.New("ошибка подключения к файловому хранилищу")
	ErrStorageSave       = errors.New("ошибка сохранения файла")
)

// Conn - Подключение к единому файловому хранилищу
type Conn struct {
	serviceName string
	timeout     time.Duration
	cfg         generic_client.Config
	client      pb.RpcStorageClient
	conn        *grpc.ClientConn
	ctx         context.Context
	ctxCancel   context.CancelFunc
}

// NewConn - Создание нового подключение к единому файловому хранилищу.
// По окончанию работы требуется закрыть соединение, вызвав Conn.Close.
// При этом все активные запросы завершат свою работу.
func NewConn(serviceName string, cfg generic_client.Config) (c *Conn, err error) {
	defer func() {
		if err != nil {
			c = nil
			err = errWrapper(err)
		}
	}()

	// Подготовка и валидация конфига
	if cfg.Timeout == 0 {
		cfg.Timeout = TIMEOUT
	}
	if err = validator.New().Struct(cfg); err != nil {
		return
	}

	// Создание нового gRPC подключения к серверу
	c = new(Conn)
	c.timeout = cfg.Timeout
	c.serviceName = serviceName
	c.ctx, c.ctxCancel = context.WithCancel(context.Background())

	conn, err := generic_client.DialWithContext(c.ctx, cfg)
	if err != nil {
		return
	}
	c.conn = conn
	c.client = pb.NewRpcStorageClient(c.conn)
	return
}

// Ping позволяет проверить состояние GRPC соединения со Storage
func (s *Conn) Ping(ctx context.Context) (err error) {
	state := s.conn.GetState()
	switch state {
	case connectivity.Ready, connectivity.Idle:
		return nil
	default:
		return fmt.Errorf("%s: некорректное состояние grpc соединения", state)
	}
}

// Close - Отключение от сервера хранилища и завершение всех активных запросов.
// Если вызвать несколько раз подряд, то выполнится только 1 раз, а в остальных случаях вернет nil
func (s *Conn) Close() (err error) {
	defer func() {
		if err != nil {
			err = errWrapper(err)
		}
	}()

	if s.conn == nil {
		return
	}

	cs := s.conn.GetState()
	if cs == connectivity.Shutdown {
		return
	}

	s.ctxCancel()
	err = s.conn.Close()
	return
}

func errWrapper(err error) error {
	return fmt.Errorf("storage: %w", err)
}

func errTimeout(method string, input interface{}) error {
	return fmt.Errorf("response timeout reached: func=%s args=%+v", method, input)
}

// ctxDead - проверяет жив контекст или нет.
// Вызывается в теле функций, имеющих context в аргументах:
// - в начале больших циклов чтобы не переходить к следующей итерации в данном цикле
// Если контекст оказывается мертв, то выполнение функции сразу следует прекратить
func ctxDead(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
