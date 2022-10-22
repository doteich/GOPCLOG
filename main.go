package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	certgen "github.com/doteich/OPC-UA-Logger/cert-gen"
	"github.com/doteich/OPC-UA-Logger/exporter"
	opcsetup "github.com/doteich/OPC-UA-Logger/setup"
	"github.com/gopcua/opcua/monitor"
)

func main() {

	config := opcsetup.SetConfig()

	if config.ClientConfig.GenerateCert {
		certgen.GeneratePEMFiles()
	}

	exporter.InitRoutes(config.LoggerConfig.TargetURL)
	exporter.InitLogs()

	if config.LoggerConfig.MetricsEnabled {

		namespace := strings.Replace(config.LoggerConfig.Name, " ", "", -1)

		go exporter.ExposeMetrics(namespace)
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

	ep := opcsetup.ValidateEndpoint(ctx, config.ClientConfig.Url, config.ClientConfig.SecurityPolicy, config.ClientConfig.SecurityMode)

	connectionParams := opcsetup.SetClientOptions(&config, ep)

	client := opcsetup.CreateClientConnection(config.ClientConfig.Url, connectionParams)
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

	go opcsetup.MonitorItems(ctx, nodeMonitor, time.Duration(config.LoggerConfig.Interval*1000000000), 1000, wg, config.Nodes)

	<-ctx.Done()
	defer ShowDone()

}

func ShowDone() {
	fmt.Println("ABORTING")

}
