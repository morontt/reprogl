package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func NewUptimeMetrics() prometheus.Counter {
	c := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "uptime_total",
			Help:      "How long the blog engine has been running (in hours)",
		},
	)

	go func(cnt prometheus.Counter) {
		var last, current int
		start := time.Now()

		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				current = int(time.Since(start).Hours())
				if current > last {
					last = current
					cnt.Inc()
				}
			}
		}
	}(c)

	return c
}
