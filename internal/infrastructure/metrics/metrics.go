package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	promNamespace = "shortly"
)

// Singleton for all metrics collection
var Collector = newCollector()

type collector struct {
	requestCount   *prometheus.CounterVec
	requestSeconds *prometheus.CounterVec
	errorCount     *prometheus.CounterVec
}

func newCollector() *collector {
	rc := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: promNamespace,
		Name:      "http_request_count",
		Help:      "Total number of requests by route",
	}, []string{"method", "path"})

	rs := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: promNamespace,
		Name:      "http_request_seconds",
		Help:      "Total amount of request time by route, in seconds",
	}, []string{"method", "path"})

	ec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: promNamespace,
		Name:      "error_count",
		Help:      "Total number of errors",
	}, []string{"method", "path", "code", "message"})

	prometheus.MustRegister(rc)
	prometheus.MustRegister(rs)
	prometheus.MustRegister(ec)

	return &collector{errorCount: ec, requestCount: rc, requestSeconds: rs}
}

func (c *collector) CollectHTTPError(method, path string, labels ...string) error {
	labels = append(labels, method, path)
	counter, err := c.errorCount.GetMetricWithLabelValues(labels...)
	if err != nil {
		return fmt.Errorf("get metric: %w", err)
	}

	counter.Inc()
	return nil
}
