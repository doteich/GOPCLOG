package main

import (
	"context"
	"time"

	exporter "github.com/doteich/OPC-UA-Logger/exporters"
	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/doteich/OPC-UA-Logger/machine/opcua_monitor"
	"github.com/doteich/OPC-UA-Logger/setup"
)

func main() {

	config := setup.SetConfig()
	logging.InitLogger()

	if config.ClientConfig.GenerateCert {
		if err := setup.CreateKeyPair(); err != nil {
			panic(err)
		}
	}

	exporter.InitExporters(config)

	ctx := context.Background()

	t := time.NewTicker(20 * time.Second)

	opcua_monitor.CreateConnectionWatcher(ctx, t, config)

}
