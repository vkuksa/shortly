package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sethvargo/go-envconfig"
	"github.com/vkuksa/shortly/internal/infrastructure/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()

	cfg := &config.AppConfig{}
	err := envconfig.Process(ctx, cfg)
	if err != nil {
		slog.Error("failed to parse configuration", slog.Any("err", err.Error()))
		os.Exit(1)
	}

	app, err := NewApp(cfg)
	if err != nil {
		log.Fatal("NewApp: ", err.Error())
	}

	if err := app.Run(ctx); err != nil {
		log.Fatal("Run: ", err.Error())
	}

	log.Print("Gracefull shutdown")
}
