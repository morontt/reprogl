package metrics

import "github.com/prometheus/client_golang/prometheus"

type GenericMetrics struct {
	Info *prometheus.GaugeVec
}

func NewGenericMetrics() *GenericMetrics {
	m := &GenericMetrics{
		Info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "info",
			Help:      "Info about blog engine",
		}, []string{"version", "build_time"}),
	}

	return m
}

func (m *GenericMetrics) SetInfo(version, buildTime string) {
	m.Info.With(prometheus.Labels{"version": version, "build_time": buildTime}).Set(1.0)
}

func (m *GenericMetrics) Collectors() []prometheus.Collector {
	return []prometheus.Collector{m.Info}
}
