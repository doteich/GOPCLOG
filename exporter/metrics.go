package exporter

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	custom_gauge = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: "newNamespace", Name: "custom_gauge_metric", Help: "New Gauge for collecting gopclog gauge metrics"}, []string{"NodeId", "NodeName"})
)

var (
	custom_counter = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: "newNamespace", Name: "custom_counter_metric", Help: "New Gauge for collecting gopclog gauge metrics"}, []string{"NodeId", "NodeName"})
)

func ExposeMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":4444", nil)

}

func SetMetricsValue(metricsType string, nodeId string, tagname string, tagValue float64) {

	if metricsType == "gauge" {
		custom_gauge.WithLabelValues(nodeId, tagname).Set(tagValue)
	}

	if metricsType == "counter" {
		custom_counter.WithLabelValues(nodeId, tagname).Add(tagValue)
	}

}
