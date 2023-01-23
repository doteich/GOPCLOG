package exporter

import (
	"time"

	"github.com/doteich/OPC-UA-Logger/exporters/http_exporter"
	"github.com/doteich/OPC-UA-Logger/exporters/metrics_exporter"
	"github.com/doteich/OPC-UA-Logger/setup"
)

func PublishData(nodeId string, iface interface{}, timestamp time.Time) {

	config := setup.SetConfig()

	for _, node := range config.Nodes {
		if node.NodeId == nodeId {

			var dataType string
			var metricsValue float64

			switch v := iface.(type) {
			case int:
				dataType = "int"
				metricsValue = float64(v)
			case int8:
				dataType = "int8"
				metricsValue = float64(v)
			case int16:
				dataType = "int16"
				metricsValue = float64(v)
			case int32:
				dataType = "int32"
				metricsValue = float64(v)
			case uint8:
				dataType = "uint8"
				metricsValue = float64(v)
			case uint16:
				dataType = "uint16"
				metricsValue = float64(v)
			case uint32:
				dataType = "uint32"
				metricsValue = float64(v)
			case float32:
				dataType = "float32"
				metricsValue = float64(v)
			case float64:
				dataType = "float64"
				metricsValue = v
			case string:
				dataType = "string"
			case bool:
				dataType = "bool"
			}

			http_exporter.PostLoggedData(node.NodeId, node.NodeName, iface, timestamp, config.LoggerConfig.Name, config.ClientConfig.Url, dataType)

			if config.LoggerConfig.MetricsEnabled && (dataType != "bool" || dataType != "string") {
				ExportMetric(node.MetricsType, node.NodeId, node.NodeName, metricsValue)
			}
		}
	}

}

func ExportMetric(metricsType string, nodeId string, name string, metricsValue float64) {

	metrics_exporter.SetMetricsValue(metricsType, nodeId, name, metricsValue)

}
