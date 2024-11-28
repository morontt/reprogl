package metrics

import "github.com/prometheus/client_golang/prometheus"

type HttpRequestsMetrics struct {
	counter *prometheus.CounterVec
}

func NewHttpRequestsMetrics() *HttpRequestsMetrics {
	return &HttpRequestsMetrics{
		counter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "http_requests",
				Name:      "count",
				Help:      "Number of HTTP requests",
			},
			[]string{"code", "method"},
		),
	}
}

func (m *HttpRequestsMetrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.counter,
	}
}
