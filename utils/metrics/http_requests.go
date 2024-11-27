package metrics

import "github.com/prometheus/client_golang/prometheus"

type HttpRequestsMetrics struct {
	Counter *prometheus.CounterVec
}

func NewHttpRequestsMetrics() *HttpRequestsMetrics {
	return &HttpRequestsMetrics{
		Counter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "http_requests",
				Name:      "cnt",
				Help:      "Number of HTTP requests",
			},
			[]string{"code", "method"},
		),
	}
}

func (m *HttpRequestsMetrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.Counter,
	}
}
