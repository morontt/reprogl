package metrics

import "github.com/prometheus/client_golang/prometheus"

type GenericMetrics struct {
	info *prometheus.GaugeVec
}

func NewGenericMetrics() *GenericMetrics {
	return &GenericMetrics{
		info: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "info",
				Help:      "Info about blog engine",
			},
			[]string{"version", "build_time"},
		),
	}
}

func (m *GenericMetrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{m.info}
}
