package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "reprogl"

type Metrics struct {
	Generic  *GenericMetrics
	Requests *HttpRequestsMetrics
}

func New() *Metrics {
	return &Metrics{
		Generic:  NewGenericMetrics(),
		Requests: NewHttpRequestsMetrics(),
	}
}

func (m *Metrics) Collectors() []prometheus.Collector {
	cs := make([]prometheus.Collector, 0)
	cs = append(cs, m.Generic.Collectors()...)
	cs = append(cs, m.Requests.Collectors()...)

	return cs
}

func (m *Metrics) SetInfo(version, buildTime string) {
	m.Generic.Info.With(prometheus.Labels{"version": version, "build_time": buildTime}).Set(1.0)
}

func (m *Metrics) IncrementRequestCount(statusCode int, method string) {
	m.Requests.Counter.With(
		prometheus.Labels{
			"code":   strconv.Itoa(statusCode),
			"method": method,
		}).Inc()
}
