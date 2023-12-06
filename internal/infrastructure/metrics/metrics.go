package metrics

import (
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

	ErrorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: promNamespace,
		Name:      "error_count",
		Help:      "Total number of errors",
	}, []string{"method", "path", "code", "message"})
)

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestSeconds)
	prometheus.MustRegister(ErrorCount)
}
