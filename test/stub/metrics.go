package stub

type MetricsCollector struct{}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{}
}

func (c *MetricsCollector) CollectHTTPError(_, _ string, _ ...string) error {
	return nil
}
