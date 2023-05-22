package http

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	promNamespace = "shortly"
)

var (
	// Generic HTTP metrics.
	requestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: promNamespace,
		Name:      "http_request_count",
		Help:      "Total number of requests by route",
	}, []string{"method", "path"})

	requestSeconds = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: promNamespace,
		Name:      "http_request_seconds",
		Help:      "Total amount of request time by route, in seconds",
	}, []string{"method", "path"})

	errorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: promNamespace,
		Name:      "error_count",
		Help:      "Total number of errors",
	}, []string{"method", "path", "code", "message"})
)

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

func initMetrics() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestSeconds)
	prometheus.MustRegister(errorCount)
}
