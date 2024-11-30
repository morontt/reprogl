package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace        = "reprogl"
	requestSubsystem = "http_requests"
)

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

func (m *Metrics) Duration(handlerName string, handlerFn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handlerFn(w, r)

		m.requests.handlerDuration.
			With(prometheus.Labels{"handler": handlerName}).
			Observe(0.001 * time.Since(start).Seconds())
	}
}
