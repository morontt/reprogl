package metrics

import "github.com/prometheus/client_golang/prometheus"

const namespace = "reprogl"

type Metrics struct {
	Generic *GenericMetrics
}

func New() *Metrics {
	return &Metrics{
		Generic: NewGenericMetrics(),
	}
}

func (m *Metrics) Collectors() []prometheus.Collector {
	cs := make([]prometheus.Collector, 0)
	cs = append(cs, m.Generic.Collectors()...)

	return cs
}
