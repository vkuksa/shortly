package http

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vkuksa/shortly/internal/infrastructure/trace"
)

const (
	// ShutdownTimeout is the time given for outstanding requests to finish before shutdown.
	DefaultShutdownTimeout = 10 * time.Second
)

type ServerConfig interface {
	Port() int
}

type Server struct {
	*http.Server
}

func NewServer(cfg ServerConfig) Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Use(middleware.RequestID)
	router.Use(trace.Middleware)
	return Server{&http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port()),
		Handler: router,
	}}
}

func (s *Server) Run() error {
	slog.Info("Starting a server", slog.Any("addr", s.Addr))
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
