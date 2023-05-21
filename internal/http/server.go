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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	shortly "github.com/vkuksa/shortly/internal/domain"
)

const (
	// ShutdownTimeout is the time given for outstanding requests to finish before shutdown.
	DefaultShutdownTimeout = 10 * time.Second

	promNamespace = "shortly"
)

var (
	// Generic HTTP metrics.
	requestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: promNamespace,
		Name:      "http_request_count",
		Help:      "Total number of requests by route",
	}, []string{"method", "path"})

	requestSeconds = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: promNamespace,
		Name:      "http_request_seconds",
		Help:      "Total amount of request time by route, in seconds",
	}, []string{"method", "path"})

	errorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: promNamespace,
		Name:      "error_count",
		Help:      "Total number of errors",
	}, []string{"method", "path", "code", "message"})
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

	// Add middlewares for logging, timeout, panic recovery and metrics tracking
	router.Use(middleware.Timeout(DefaultShutdownTimeout))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(trackMetrics)

	// Register domain handlers
	s.registerLinkRoutes(router)

	s.srv.Handler = router

	// Enable internal administrative endpoints.
	go func() {
		listenAndServeAdmin()
	}()

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
func (s *Server) port() int {
	if s.ln == nil {
		return 0
	}
	return s.ln.Addr().(*net.TCPAddr).Port
}

// url returns the local base URL of the running server.
func (s *Server) url() string {
	// Use localhost unless a domain is specified.
	domain := "localhost"
	if s.Domain != "" {
		domain = s.Domain
	}
	scheme := "http"
	if s.Scheme != "" {
		scheme = s.Scheme
	}

	return fmt.Sprintf("%s://%s:%d", scheme, domain, s.port())
}

// Handles errors gracefully
func (s *Server) handleError(w http.ResponseWriter, r *http.Request, err error) {
	code, message := shortly.ErrorCode(err), shortly.ErrorMessage(err)

	errorCount.WithLabelValues(r.Method, r.URL.Path, code, message).Inc()
	// TODO: report internal errors

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

// ListenAndServeAdmin runs an HTTP server with administrative infromation
func listenAndServeAdmin() error {
	h := http.NewServeMux()
	h.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":6060", h)
}

// trackMetrics is middleware for tracking the request count and timing per route.
func trackMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtain path template & start time of request.
		t := time.Now()

		// Delegate to next handler in middleware chain.
		next.ServeHTTP(w, r)

		requestCount.WithLabelValues(r.Method, r.URL.Path).Inc()
		requestSeconds.WithLabelValues(r.Method, r.URL.Path).Add(float64(time.Since(t).Seconds()))
	})
}
