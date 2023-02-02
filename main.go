package main

import (
	"strings"

	"github.com/doteich/OPC-UA-Logger/exporters/db_exporter"
	"github.com/doteich/OPC-UA-Logger/exporters/http_exporter"
	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/doteich/OPC-UA-Logger/exporters/metrics_exporter"
	"github.com/doteich/OPC-UA-Logger/machine/opcua_monitor"
	"github.com/doteich/OPC-UA-Logger/setup"
)

func main() {

	config := setup.SetConfig()

	if config.ClientConfig.GenerateCert {
		setup.GeneratePEMFiles()
	}

	namespace := strings.Replace(config.LoggerConfig.Name, " ", "", -1)

	db_exporter.SetupDBConnection(namespace)

	http_exporter.InitRoutes(config.LoggerConfig.TargetURL)
	logging.InitLogs()

	if config.LoggerConfig.MetricsEnabled {

		go metrics_exporter.ExposeMetrics(namespace)
	}

	opcua_monitor.CreateOPCUAMonitor(config)

	//exporter.ReadLogFile()

}
