package metrics_exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var custom_gauge *prometheus.GaugeVec
var custom_counter *prometheus.GaugeVec
var custom_technical_counter *prometheus.CounterVec
var Failed_requests *prometheus.CounterVec
var opcua_connects *prometheus.CounterVec

func ExposeMetrics(namespace string) {

	custom_gauge = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Name: "custom_gauge_metric", Help: "Metric for collecting gopclog gauge type tag values"}, []string{"NodeId", "NodeName"})
	custom_counter = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: namespace, Name: "custom_counter_metric", Help: "Metric for collecting gopclog counter type tag values"}, []string{"NodeId", "NodeName"})
	custom_technical_counter = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Name: "custom_technical_counter_metric", Help: "Metric for collecting gopclog technical counter type tag values"}, []string{"NodeId", "NodeName"})
	Failed_requests = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Name: "number_failed_request", Help: "Metric for collecting gopclog number of failed request to the specified target URL"}, []string{"url"})
	opcua_connects = promauto.NewCounterVec(prometheus.CounterOpts{Namespace: namespace, Name: "number_opcua_connection_attempts", Help: "Metric for collecting gopclog number of failed request to the specified target URL"}, []string{"server"})
}

func SetMetricsValue(metricsType string, nodeId string, tagname string, tagValue float64) {

	if metricsType == "Gauge" {
		custom_gauge.WithLabelValues(nodeId, tagname).Set(tagValue)
	}

	if metricsType == "Counter" {
		custom_counter.WithLabelValues(nodeId, tagname).Set(tagValue)
	}

	if metricsType == "Technical Counter" {
		custom_technical_counter.WithLabelValues(nodeId, tagname).Add(tagValue)
	}

}

func LogReconnects(s string) {
	opcua_connects.WithLabelValues(s).Add(1)
}
