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
			log.Printf("ошибка загрузки конфигурационного файла: %s", err.Error())

			return
		}
	}

	if err := aconf.Load(&cfg); err != nil {
		log.Printf("ошибка инициализации конфигурации: %s", err.Error())

		return
	}
	// Подгрузим настройки маркетплейса в конфигурацию
	cfg.FillMarketPlaceMap()

	app, err := app.New(ctx, cfg)
	if err != nil {
		log.Printf("app.New failed: %s", err.Error())

		return
	}

	err = app.Start(ctx)
	if err != nil {
		log.Printf("app.Start failed: %s", err.Error())

		return
	}
}
