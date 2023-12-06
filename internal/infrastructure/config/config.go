package config

import (
	"github.com/vkuksa/shortly/internal/infrastructure/http"
	"github.com/vkuksa/shortly/internal/infrastructure/metrics"
)

type AppConfig struct {
	HTTPServer    *http.Config    `env:",prefix=HTTP_SERVER_"`
	MetricsServer *metrics.Config `env:",prefix=METRICS_HTTP_SERVER_"`
}
