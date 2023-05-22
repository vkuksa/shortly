package http

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	shortly "github.com/vkuksa/shortly/internal/domain"
)

const (
	// ShutdownTimeout is the time given for outstanding requests to finish before shutdown.
	DefaultShutdownTimeout = 10 * time.Second
)

type Server struct {
	ln  net.Listener
	srv *http.Server

	conf *Config

	LinkService shortly.LinkService

	Prom struct {
		ln  net.Listener
		srv *http.Server
	}
}

func NewServer(c Config) *Server {
	s := &Server{
		srv:  &http.Server{},
		conf: &c,
	}

	router := chi.NewRouter()

	// Add middlewares for logging, timeout, panic recovery and metrics tracking
	router.Use(middleware.Timeout(DefaultShutdownTimeout))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	if c.Prometheus.Enabled {
		initMetrics()
		router.Use(trackMetrics)
		s.Prom.srv = &http.Server{}
	}

	// Register domain handlers
	s.registerLinkRoutes(router)

	s.srv.Handler = router

	return s
}

// Open validates the server options and begins listening on the bind address.
func (s *Server) Open() (err error) {
	// Start a separate http server for prometheus
	if s.conf.Prometheus.Enabled {
		if err = s.listenAndServePrometheus(); err != nil {
			return err
		}
	}

	if s.ln, err = net.Listen("tcp", ":"+s.conf.Port); err != nil {
		return err
	}

	go func() {
		_ = s.srv.Serve(s.ln)
	}()

	return nil
}

// Close gracefully shuts down the server.
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	defer cancel()

	if s.conf.Prometheus.Enabled {
		_ = s.Prom.srv.Shutdown(ctx)
	}

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

// url returns the local base URL of the running server.
func (s *Server) URL() string {
	// Use localhost unless a domain is specified.
	host := "localhost"
	if s.conf.Host != "" {
		host = s.conf.Host
	}
	scheme := "http"
	if s.conf.Scheme != "" {
		scheme = s.conf.Scheme
	}

	return fmt.Sprintf("%s://%s:%d", scheme, host, s.Port())
}

// Handles errors gracefully
func (s *Server) handleError(w http.ResponseWriter, r *http.Request, err error) {
	code, message := shortly.ErrorCode(err), shortly.ErrorMessage(err)

	errorCount.WithLabelValues(r.Method, r.URL.Path, code, message).Inc()

	s.logError(r, err)
	http.Error(w, message, errorStatusCode(code))
}

// LogError logs an error with the HTTP route information.
func (s *Server) logError(r *http.Request, err error) {
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
func errorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}

// listenAndServePrometheus runs an HTTP server with prometheus metrics
func (s *Server) listenAndServePrometheus() (err error) {
	h := http.NewServeMux()
	h.Handle("/metrics", promhttp.Handler())
	s.Prom.srv.Handler = h

	if s.Prom.ln, err = net.Listen("tcp", ":"+s.conf.Prometheus.Port); err != nil {
		return err
	}

	go func() {
		_ = s.Prom.srv.Serve(s.Prom.ln)
	}()

	return nil
}
