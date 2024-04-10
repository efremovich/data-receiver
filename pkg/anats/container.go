package astral_nats

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// поднимает временный контейнер с nats для тестов. возвращает url для подключения к натс.
func CreateNatsTempContainer(ctx context.Context) (string, error) {
	req := testcontainers.ContainerRequest{
		Image:        "nats:latest",
		ExposedPorts: []string{"4222/tcp"},
		AutoRemove:   true,
		WaitingFor:   wait.ForListeningPort("4222/tcp"),
	}
	req.Cmd = []string{"-js"}

	natsContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", err
	}

	natsPort, err := natsContainer.MappedPort(ctx, "4222")
	if err != nil {
		return "", err
	}

	host, err := natsContainer.Host(ctx)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("nats://%s:%s", host, natsPort.Port())

	return url, nil
}
