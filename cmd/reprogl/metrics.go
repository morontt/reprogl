package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getPrometheusMetricsHandler(cs ...prometheus.Collector) http.Handler {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	reg.MustRegister(cs...)

	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})
}
