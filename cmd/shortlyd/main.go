package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sethvargo/go-envconfig"
	"github.com/vkuksa/shortly/internal/infrastructure/config"
	"github.com/vkuksa/shortly/internal/infrastructure/http"
	"github.com/vkuksa/shortly/internal/infrastructure/metrics"
	"github.com/vkuksa/shortly/internal/infrastructure/storage/inmem"
	"github.com/vkuksa/shortly/internal/interface/controller/rest"
	"github.com/vkuksa/shortly/internal/interface/repository"
	"github.com/vkuksa/shortly/internal/link"
)

type ChiRegisterer interface {
	Register(router chi.Router)
}

func makeHttpServer(conf *http.Config, controllers ...ChiRegisterer) *http.Server {
	router := chi.NewRouter()
	router.Use(middleware.Timeout(10 * time.Second))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(metrics.Middleware)

	for _, c := range controllers {
		c.Register(router)
	}

	return http.NewServer(conf.BuildAddr(), router)
}

func makeMetricsServer(conf *metrics.Config) *http.Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Handle("/metrics", promhttp.Handler())

	return http.NewServer(conf.BuildAddr(), router)
}

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

	metrics := metrics.Collector

	storage := inmem.NewStorage()
	linkRepository := repository.New(storage)
	linkUsecase := link.NewUseCase(linkRepository)
	linkController := rest.NewLinkController(linkUsecase, metrics)

	httpServer := makeHttpServer(cfg.HTTPServerConfig, linkController)
	metricsServer := makeMetricsServer(cfg.MetricsServerConfig)

	app := NewApp(httpServer, metricsServer)
	if err := app.Run(ctx); err != nil {
		log.Fatal("run: ", err.Error())
	}

	log.Print("Gracefull shutdown")
}
