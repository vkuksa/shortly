package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	// ShutdownTimeout is the time given for outstanding requests to finish before shutdown.
	DefaultShutdownTimeout = 10 * time.Second
)

type Server struct {
	srv  *http.Server
	once sync.Once
}

func NewServer(addr string, router chi.Router) *Server {
	s := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	return &Server{srv: s}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		_ = s.shutdown(ctx)
	}()

	slog.Info("Starting a server", slog.Any("addr", s.srv.Addr))
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Serve(rw http.ResponseWriter, req *http.Request) {
	s.srv.Handler.ServeHTTP(rw, req)
}

func (s *Server) Close(ctx context.Context) error {
	return s.shutdown(ctx)
}

func (s *Server) shutdown(ctx context.Context) (err error) {
	s.once.Do(func() {
		ctx, cancel := context.WithTimeout(ctx, DefaultShutdownTimeout)
		defer cancel()

		err = s.srv.Shutdown(ctx)
	})
	return
}
