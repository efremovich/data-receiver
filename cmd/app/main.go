package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/app"
	"github.com/efremovich/data-receiver/pkg/aconf/v3"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	c := config.Config{}

	if err := aconf.Load(&c); err != nil {
		log.Fatalf("aconf.Load failed: %s", err.Error())
	}
	app, err := app.New(ctx, c)
	if err != nil {
		log.Fatalf("app.New failed: %s", err.Error())
	}

	err = app.Start(ctx)
	if err != nil {
		log.Fatalf("app.Start failed: %s", err.Error())
	}
}
