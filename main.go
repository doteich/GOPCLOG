package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/doteich/OPC-UA-Logger/exporters/db_exporter"
	"github.com/doteich/OPC-UA-Logger/exporters/http_exporter"
	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/doteich/OPC-UA-Logger/exporters/metrics_exporter"
	"github.com/doteich/OPC-UA-Logger/machine/opcua_monitor"
	"github.com/doteich/OPC-UA-Logger/setup"
	"github.com/gopcua/opcua/monitor"
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

	//exporter.ReadLogFile()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-signalCh
		println()
		cancel()
	}()

	ep := opcua_monitor.ValidateEndpoint(ctx, config.ClientConfig.Url, config.ClientConfig.SecurityPolicy, config.ClientConfig.SecurityMode)

	connectionParams := opcua_monitor.SetClientOptions(&config, ep)

	client := opcua_monitor.CreateClientConnection(config.ClientConfig.Url, connectionParams)
	err := client.Connect(ctx)

	if err != nil {
		fmt.Println(err)
	}

	defer client.CloseSessionWithContext(ctx)

	nodeMonitor, err := monitor.NewNodeMonitor(client)

	if err != nil {
		panic("Failed to setup monitor")
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go opcua_monitor.MonitorItems(ctx, nodeMonitor, time.Duration(config.LoggerConfig.Interval*1000000000), 1000, wg, config.Nodes)

	<-ctx.Done()
	defer ShowDone()

}

func ShowDone() {
	fmt.Println("ABORTING")

}
