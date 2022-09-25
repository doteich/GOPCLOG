package exporter

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsObject struct {
	Namespace   string
	MetricsType string
	Name        string
}

func ExposeMetrics() {
	go registerCustomMetrics()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":4444", nil)

}

func registerCustomMetrics() {
	gauge := promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: "newNamespace", Name: "custom_gauge_metric", Help: "New Gauge for collecting gopclog gauge metrics"}, []string{"tag", "name"})
	counter := promauto.NewCounterVec(prometheus.CounterOpts{Namespace: "newNamespace", Name: "custom_gauge_metric", Help: "New Gauge for collecting gopclog gauge metrics"}, []string{"tag", "name"})
	prometheus.Register(gauge)
	prometheus.Register(counter)

}
