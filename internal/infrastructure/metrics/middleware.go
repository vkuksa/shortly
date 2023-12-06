package metrics

import (
	"net/http"
	"time"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		next.ServeHTTP(w, r)

		requestCount.WithLabelValues(r.Method, r.URL.Path).Inc()
		requestSeconds.WithLabelValues(r.Method, r.URL.Path).Add(float64(time.Since(t).Seconds()))
	})
}
