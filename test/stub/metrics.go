package stub

type MetricsCollector struct{}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{}
}

func (c *MetricsCollector) CollectHttpError(method, path string, labels ...string) error {
	return nil
}
