package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/app"
)

func main() {
	var envPath string

	flag.StringVar(&envPath, "envPath", "", "путь к local.env")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	cfg, err := config.NewConfig(envPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
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
