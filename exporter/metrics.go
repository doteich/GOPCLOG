package exporter

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	gauge = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: "newNamespace", Name: "custom_gauge_metric", Help: "New Gauge for collecting gopclog gauge metrics"}, []string{"NodeId", "NodeName"})
)

var (
	counter = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: "newNamespace", Name: "custom_counter_metric", Help: "New Gauge for collecting gopclog gauge metrics"}, []string{"NodeId", "NodeName"})
)

func ExposeMetrics() {
	go registerCustomMetrics()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":4444", nil)

}

func registerCustomMetrics() {
	setMetricsValue("gauge", "test", "tagtest", 10.00)
}

func setMetricsValue(MetricsType string, Name string, Tag string, value float64) {
	gauge.WithLabelValues(Tag, Name).Set(value)
}
