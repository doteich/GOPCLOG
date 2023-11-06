package main

import (
	"net/http"

	exporter "github.com/doteich/OPC-UA-Logger/exporters"
	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/doteich/OPC-UA-Logger/machine/opcua_monitor"
	"github.com/doteich/OPC-UA-Logger/setup"
)

func main() {

	config := setup.SetConfig()
	logging.InitLogger()

	if config.ClientConfig.GenerateCert {
		setup.GeneratePEMFiles()
	}

	exporter.InitExporters(config)

	go func() {
		r := server()
		http.ListenAndServe(":6000", r)
	}()

	opcua_monitor.CreateOPCUAMonitor(config)
}
