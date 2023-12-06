package main

import (
	"context"
	"errors"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vkuksa/shortly/internal/infrastructure/config"
	"github.com/vkuksa/shortly/internal/infrastructure/http"
	"github.com/vkuksa/shortly/internal/infrastructure/metrics"
	"github.com/vkuksa/shortly/internal/infrastructure/storage/inmem"
	"github.com/vkuksa/shortly/internal/interface/controller"
	"github.com/vkuksa/shortly/internal/interface/repository"
	"github.com/vkuksa/shortly/internal/usecase"
	"golang.org/x/sync/errgroup"
)

type App struct {
	HTTPServer    *http.Server
	MetricsServer *http.Server
}

func NewApp(conf *config.AppConfig) (*App, error) {
	storage := inmem.NewStorage()
	linkRepository := repository.New(storage)
	linkUsecase := usecase.NewLinkUseCase(linkRepository)
	linkController := controller.NewLinkController(linkUsecase)

	return &App{
		HTTPServer:    makeServer(conf.HTTPServer, linkController),
		MetricsServer: makeMetricsServer(conf.MetricsServer),
	}, nil
}

func makeServer(conf *http.Config, linkController *controller.LinkController) *http.Server {
	router := chi.NewRouter()
	router.Use(middleware.Timeout(10 * time.Second))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(metrics.Middleware)

	linkController.Register(router)

	return http.NewServer(conf.BuildAddr(), router)
}

func makeMetricsServer(conf *metrics.Config) *http.Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Handle("/metrics", promhttp.Handler())

	return http.NewServer(conf.BuildAddr(), router)
}

func (app *App) Stop(ctx context.Context) error {
	return errors.Join(app.HTTPServer.Close(ctx), app.MetricsServer.Close(ctx))
}

func (app *App) Run(ctx context.Context) (err error) {
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return app.HTTPServer.Run(ctx)
	})

	g.Go(func() error {
		return app.MetricsServer.Run(ctx)
	})

	g.Go(func() error {
		<-gCtx.Done()
		return app.Stop(ctx)
	})

	return g.Wait()
}
