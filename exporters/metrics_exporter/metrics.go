package metrics_exporter

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var custom_gauge *prometheus.GaugeVec
var custom_counter *prometheus.CounterVec
var Failed_requests *prometheus.CounterVec
var last_Count float64
var diff float64

func ExposeMetrics(namespace string) {

	custom_gauge = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Name: "custom_gauge_metric", Help: "Metric for collecting gopclog gauge type tag values"}, []string{"NodeId", "NodeName"})
	custom_counter = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Name: "custom_counter_metric", Help: "Metric for collecting gopclog counter type tag values"}, []string{"NodeId", "NodeName"})
	Failed_requests = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Name: "number_failed_request", Help: "Metric for collecting gopclog number of failed request to the specified target URL"}, []string{"url"})
	last_Count = 0
	diff = 0

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":4444", nil)

}

func SetMetricsValue(metricsType string, nodeId string, tagname string, tagValue float64) {

	if metricsType == "Gauge" {
		custom_gauge.WithLabelValues(nodeId, tagname).Set(tagValue)
	}

	if metricsType == "Counter" {
		diff = tagValue - last_Count
		custom_counter.WithLabelValues(nodeId, tagname).Add(diff)
		last_Count = tagValue
	}

}
