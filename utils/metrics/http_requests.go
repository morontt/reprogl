package metrics

import "github.com/prometheus/client_golang/prometheus"

type HttpRequestsMetrics struct {
	counter         *prometheus.CounterVec
	handlerDuration *prometheus.SummaryVec
}

func NewHttpRequestsMetrics() *HttpRequestsMetrics {
	return &HttpRequestsMetrics{
		counter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: requestSubsystem,
				Name:      "count",
				Help:      "Number of HTTP requests",
			},
			[]string{"code", "method"},
		),
		handlerDuration: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  namespace,
				Subsystem:  requestSubsystem,
				Name:       "handler_duration",
				Help:       "Duration of HTTP handler processing in milliseconds",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.005},
			},
			[]string{"handler"},
		),
	}
}

func (m *HttpRequestsMetrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		m.counter,
		m.handlerDuration,
	}
}
