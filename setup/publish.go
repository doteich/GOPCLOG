// package for dispatching monitored data to specified routes (metrics, http endpoint etc.)
package setup

import (
	"time"

	"github.com/doteich/OPC-UA-Logger/exporter"
)

func PublishData(nodeId string, iface interface{}, timestamp time.Time) {

	config := SetConfig()

	for _, node := range config.Nodes {
		if node.NodeId == nodeId {

			exporter.PostLoggedData(node.NodeId, node.NodeName, iface, timestamp, config.LoggerConfig.Name, config.ClientConfig.Url)

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
		exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case int8:
		value := float64(v)
		exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case int16:
		value := float64(v)
		exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case int32:
		value := float64(v)
		exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case uint:
		value := float64(v)
		exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case uint8:
		value := float64(v)
		exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case uint16:
		value := float64(v)
		exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case uint32:
		value := float64(v)
		exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case float32:
		value := float64(v)
		exporter.SetMetricsValue(metricsType, nodeId, name, value)
	case float64:
		exporter.SetMetricsValue(metricsType, nodeId, name, v)
	}

}
