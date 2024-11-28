package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	conf "github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/app"
	"github.com/efremovich/data-receiver/pkg/aconf/v3"
)

func main() {
	var envPath string

	flag.StringVar(&envPath, "envPath", "./config/vars.env", "путь к local.env")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfg := conf.Config{}

	if envPath != "" {
		err := aconf.PreloadEnvsFile(envPath)
		if err != nil {
			log.Fatalf("ошибка загрузки конфигурационного файла: %s", err.Error())
		}
	}

	if err := aconf.Load(&cfg); err != nil {
		log.Fatalf("ошибка инициализации конфигурации: %s", err.Error())
	}

	app, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("app.New failed: %s", err.Error())
	}

	err = app.Start(ctx)
	if err != nil {
		log.Fatalf("app.Start failed: %s", err.Error())
	}
}
