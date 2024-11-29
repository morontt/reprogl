package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "reprogl"

type Metrics struct {
	generic  *GenericMetrics
	requests *HttpRequestsMetrics

	uptime prometheus.Counter
}

func New() *Metrics {
	return &Metrics{
		generic:  NewGenericMetrics(),
		requests: NewHttpRequestsMetrics(),
		uptime:   NewUptimeMetrics(),
	}
}

func (m *Metrics) Collectors() []prometheus.Collector {
	cs := make([]prometheus.Collector, 0)
	cs = append(cs, m.generic.Collectors()...)
	cs = append(cs, m.requests.Collectors()...)
	cs = append(cs, prometheus.Collector(m.uptime))

	return cs
}

func (m *Metrics) SetInfo(version, buildTime string) {
	m.generic.info.With(prometheus.Labels{"version": version, "build_time": buildTime}).Set(1.0)
}

func (m *Metrics) IncrementRequestCount(statusCode int, method string) {
	m.requests.counter.With(
		prometheus.Labels{
			"code":   strconv.Itoa(statusCode),
			"method": method,
		}).Inc()
}
