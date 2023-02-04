package main

import (
	"fmt"

	exporter "github.com/doteich/OPC-UA-Logger/exporters"
	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/doteich/OPC-UA-Logger/machine/opcua_monitor"
	"github.com/doteich/OPC-UA-Logger/setup"
)

func main() {

	config := setup.SetConfig()

	if config.ClientConfig.GenerateCert {
		setup.GeneratePEMFiles()
	}

	exporter.InitExporters(&config)
	fmt.Println("INIT01")
	logging.InitLogs()

	fmt.Println("INIT")
	opcua_monitor.CreateOPCUAMonitor(config)

}
