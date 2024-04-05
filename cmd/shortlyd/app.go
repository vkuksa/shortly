package main

import (
	"context"
	"errors"

	"github.com/vkuksa/shortly/internal/infrastructure/http"
	"golang.org/x/sync/errgroup"
)

type App struct {
	restServer    *http.Server
	metricsServer *http.Server
}

func NewApp(restServer, metricsServer *http.Server) *App {
	return &App{
		restServer:    restServer,
		metricsServer: metricsServer,
	}
}

func (app *App) Stop(ctx context.Context) error {
	return errors.Join(app.restServer.Close(ctx), app.metricsServer.Close(ctx))
}

func (app *App) Run(ctx context.Context) (err error) {
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return app.restServer.Run(ctx)
	})

	g.Go(func() error {
		return app.metricsServer.Run(ctx)
	})

	g.Go(func() error {
		<-gCtx.Done()
		return app.Stop(ctx)
	})

	return g.Wait()
}
