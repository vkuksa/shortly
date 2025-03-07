package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/vkuksa/shortly/internal/app"
	"github.com/vkuksa/shortly/internal/infrastructure/config"
)

func main() {
	a := app.New(
		config.NewMongo(),
		config.NewServer(),
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		_ = a.Stop(ctx)
	}()

	a.Run()
}
