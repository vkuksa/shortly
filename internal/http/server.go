package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	shortly "github.com/vkuksa/shortly/internal/domain"
)

const (
	// ShutdownTimeout is the time given for outstanding requests to finish before shutdown.
	DefaultShutdownTimeout = 10 * time.Second
)

type Server struct {
	ln  net.Listener
	srv *http.Server

	Addr   string
	Scheme string
	Domain string

	LinkService shortly.LinkService
}

func NewServer() *Server {
	s := &Server{
		srv: &http.Server{},
	}

	router := chi.NewRouter()

	// Add middlewares for logging, timeout and panic recovery
	router.Use(middleware.Timeout(DefaultShutdownTimeout))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	// Register handlers
	s.registerLinkRoutes(router)

	s.srv.Handler = router

	return s
}

// Open validates the server options and begins listening on the bind address.
func (s *Server) Open() (err error) {
	if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
		return err
	}

	go func() {
		err := s.srv.Serve(s.ln)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("[http] error: %s", err)
		}
	}()

	return nil
}

// Close gracefully shuts down the server.
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	defer cancel()
	return s.srv.Shutdown(ctx)
}

// Port returns the TCP port for the running server.
// This is useful in tests where we allocate a random port by using ":0".
func (s *Server) Port() int {
	if s.ln == nil {
		return 0
	}
	return s.ln.Addr().(*net.TCPAddr).Port
}

// URL returns the local base URL of the running server.
func (s *Server) URL() string {
	// Use localhost unless a domain is specified.
	domain := "localhost"
	if s.Domain != "" {
		domain = s.Domain
	}
	scheme := "http"
	if s.Scheme != "" {
		scheme = s.Scheme
	}

	return fmt.Sprintf("%s://%s:%d", scheme, domain, s.Port())
}

// Handles errors gracefully
func (s *Server) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	code, message := shortly.ErrorCode(err), shortly.ErrorMessage(err)

	// TODO: track metrics
	// TODO: report internal errors

	s.LogError(r, err)
	http.Error(w, message, ErrorStatusCode(code))
}

// LogError logs an error with the HTTP route information.
func (s *Server) LogError(r *http.Request, err error) {
	log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
}

// lookup of application error codes to HTTP status codes.
var codes = map[string]int{
	shortly.ErrConflict: http.StatusConflict,
	shortly.ErrInvalid:  http.StatusBadRequest,
	shortly.ErrNotFound: http.StatusNotFound,
	shortly.ErrInternal: http.StatusInternalServerError,
}

// ErrorStatusCode returns the associated HTTP status code for a WTF error code.
func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}
