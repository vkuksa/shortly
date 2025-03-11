package trace

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer("http-request-middleware")
		traceCtx, span := tracer.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		defer span.End()
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.target", r.URL.Path),
			attribute.String("http.host", r.Host),
			attribute.String("http.flavor", r.Proto),
			attribute.String("http.user_agent", r.UserAgent()),
		)

		cr := r.Clone(traceCtx)
		next.ServeHTTP(w, cr)
	})
}
