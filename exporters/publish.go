package exporter

import (
	"strings"
	"time"

	"github.com/doteich/OPC-UA-Logger/exporters/db_exporter"
	"github.com/doteich/OPC-UA-Logger/exporters/metrics_exporter"
	"github.com/doteich/OPC-UA-Logger/setup"
)

func PublishData(nodeId string, iface interface{}, timestamp time.Time) {

	config := setup.SetConfig()

	for _, node := range config.Nodes {
		if node.NodeId == nodeId {

			//http_exporter.PostLoggedData(node.NodeId, node.NodeName, iface, timestamp, config.LoggerConfig.Name, config.ClientConfig.Url)
			namespace := strings.Replace(config.LoggerConfig.Name, " ", "", -1)
			db_exporter.InsertValues(namespace, node.NodeId, node.NodeName, iface, timestamp, config.LoggerConfig.Name, config.ClientConfig.Url)

			if config.LoggerConfig.MetricsEnabled {
				ExportMetric(node.MetricsType, node.NodeId, node.NodeName, iface)
			}
		}
	}

}

func ExportMetric(metricsType string, nodeId string, name string, iface interface{}) {

	switch v := iface.(type) {
	case int:
		value := float64(v)
		metrics_exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case int8:
		value := float64(v)
		metrics_exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case int16:
		value := float64(v)
		metrics_exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case int32:
		value := float64(v)
		metrics_exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case uint:
		value := float64(v)
		metrics_exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case uint8:
		value := float64(v)
		metrics_exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case uint16:
		value := float64(v)
		metrics_exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case uint32:
		value := float64(v)
		metrics_exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case float32:
		value := float64(v)
		metrics_exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case float64:
		metrics_exporter.SetMetricsValue(metricsType, nodeId, name, v)
	}

}
