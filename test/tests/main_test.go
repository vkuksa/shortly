package test

import (
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/vkuksa/shortly/internal/infrastructure/config"
	"github.com/vkuksa/shortly/internal/infrastructure/http"
	"github.com/vkuksa/shortly/internal/infrastructure/storage/inmem"
	"github.com/vkuksa/shortly/internal/interface/controller/rest"
	"github.com/vkuksa/shortly/internal/interface/repository"
	"github.com/vkuksa/shortly/internal/link"
	"github.com/vkuksa/shortly/test/stub"
)

const (
	TestingHost = "localhost"
	TestingPort = 8081
)

var httpServer *http.Server

func makeHTTPServer(conf *http.Config, linkController *rest.LinkController) *http.Server {
	router := chi.NewRouter()
	linkController.Register(router)
	return http.NewServer(conf.BuildAddr(), router)
}

func newTestAppConfig() *config.AppConfig {
	return &config.AppConfig{
		HTTPServerConfig: &http.Config{
			Host: TestingHost,
			Port: TestingPort,
		},
	}
}

func TestMain(m *testing.M) {
	cfg := newTestAppConfig()

	errorHandler := stub.NewErrorHandler()

	storage := inmem.NewStorage()
	linkRepository := repository.New(storage)
	linkUsecase := link.NewUseCase(linkRepository)
	linkController := rest.NewLinkController(linkUsecase, errorHandler)

	httpServer = makeHTTPServer(cfg.HTTPServerConfig, linkController)

	os.Exit(m.Run())
}
