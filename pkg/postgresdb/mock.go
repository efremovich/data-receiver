package postgresdb

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// запускает посгрес в контейнере и выполняет миграции
// требует запущенный докер и установленный goose
// проверено только на винде
func GetMockConn(pathToMigrations string) (*DBConnection, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	req := testcontainers.ContainerRequest{
		Image:        "postgres",
		ExposedPorts: []string{"5432/tcp"},
		AutoRemove:   true,
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	pg, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	port, err := pg.MappedPort(ctx, "5432")
	if err != nil {
		return nil, "", err
	}

	connString := fmt.Sprintf("postgres://postgres:postgres@localhost:%d/postgres?sslmode=disable", port.Int())

	markingDB, err := New(context.Background(), connString, connString)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка при создании подключения к бд маркировки: %s", err.Error())
	}

	err = runMigrations(port.Int(), pathToMigrations)
	if err != nil {
		return nil, "", err
	}

	return markingDB, connString, nil
}

// проверено только на винде
func runMigrations(port int, pathToMigrations string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "goose", "postgres", fmt.Sprintf("user=postgres password=postgres dbname=postgres host=localhost port=%d sslmode=disable", port), "up")
	} else {
		cmd = exec.Command("bash", "-c", fmt.Sprintf("goose postgres \"user=postgres password=postgres dbname=postgres host=localhost port=%d sslmode=disable\" up", port))
	}

	cmd.Dir = pathToMigrations

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("oшибка выполнения команды: %s", output)
	}

	return nil
}
