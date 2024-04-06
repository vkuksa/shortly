package config

import (
	"github.com/vkuksa/shortly/internal/infrastructure/http"
	"github.com/vkuksa/shortly/internal/infrastructure/metrics"
)

type AppConfig struct {
	HTTPServerConfig        *http.Config    `env:",prefix=HTTP_SERVER_"`
	MetricsServerConfig     *metrics.Config `env:",prefix=METRICS_HTTP_SERVER_"`
	MongodbConnectionString string          `env:"MONGODB_CONNECTION_STRING"`
}
