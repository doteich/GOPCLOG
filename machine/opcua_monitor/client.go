package opcua_monitor

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	exporter "github.com/doteich/OPC-UA-Logger/exporters"
	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/doteich/OPC-UA-Logger/exporters/websockets"
	"github.com/doteich/OPC-UA-Logger/setup"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/monitor"
	"github.com/gopcua/opcua/ua"
)

type s_struct struct {
	sub   *monitor.Subscription
	tChan chan bool
}

var (
	opcclient   *opcua.Client
	Subs        map[uint32]s_struct
	NodeMonitor *monitor.NodeMonitor
)

func CreateOPCUAMonitor(config *setup.Config) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-signalCh
		println()
		cancel()
	}()

	ep := ValidateEndpoint(ctx, config.ClientConfig.Url, config.ClientConfig.SecurityPolicy, config.ClientConfig.SecurityMode)

	connectionParams := SetClientOptions(config, ep)

	var e error

	opcclient, e = CreateClientConnection(config.ClientConfig.Url, connectionParams)

	if e != nil {
		logging.LogError(e, "error while creating opc client", "opcua")
		return
	}

	err := opcclient.Connect(ctx)

	if err != nil {
		logging.LogError(err, "Error connecting to opcua server", "opcua")
		return
	}

	defer opcclient.CloseSession(ctx)

	NodeMonitor, err = monitor.NewNodeMonitor(opcclient)

	if err != nil {
		logging.LogError(err, "Error while setting up the node monitor", "opcua")
	}

	websockets.InitOPCUARead(opcclient)
	exporter.SetOPCUAClient(opcclient)

	wg := &sync.WaitGroup{}
	Subs = make(map[uint32]s_struct)

	go MonitorItems(ctx, NodeMonitor, config.LoggerConfig.Interval, 1000, config.Nodes, make(chan bool))
	go StartKeepAlive(ctx, NodeMonitor, 1000, make(chan bool))

	wg.Add(1)

	go MonitorSubscriptions(ctx, wg, config.LoggerConfig.Interval, config.Nodes)

	//<-ctx.Done()

	wg.Wait()

	defer func() {
		logging.LogGeneric("warning", "Shutting down opcua monitor", "opcua")
	}()

}

func ValidateEndpoint(ctx context.Context, endpoint string, policy string, mode string) *ua.EndpointDescription {
	endpoints, err := opcua.GetEndpoints(ctx, endpoint)

	if err != nil {
		logging.LogError(err, "No Matching Endpoint Found - Check Configuration", "opcua")
	}

	ep := opcua.SelectEndpoint(endpoints, policy, ua.MessageSecurityModeFromString(mode))

	if ep == nil {
		logging.LogError(nil, "No Matching Endpoint Found - Check Configuration", "opcua")
		panic("No Matching Endpoint Found - Check Configuration")
	}

	return ep

}

func SetClientOptions(config *setup.Config, ep *ua.EndpointDescription) []opcua.Option {

	// basic params
	connectionParams := []opcua.Option{
		opcua.SecurityPolicy(config.ClientConfig.SecurityPolicy),
		opcua.SecurityModeString(config.ClientConfig.SecurityMode),
		opcua.AutoReconnect(true),
		opcua.ReconnectInterval(time.Second * 20),
	}

	if config.ClientConfig.SecurityMode != "None" || config.ClientConfig.SecurityPolicy != "None" {
		connectionParams = append(connectionParams, opcua.CertificateFile("./certs/cert.pem"))
		connectionParams = append(connectionParams, opcua.PrivateKeyFile("./certs/key.pem"))
	}

	switch config.ClientConfig.AuthType {
	case "User & Password":
		connectionParams = append(connectionParams, opcua.AuthUsername(config.ClientConfig.Username, config.ClientConfig.Password))
		connectionParams = append(connectionParams, opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeUserName))
	case "Certificate":
		//connectionParams = append(connectionParams, opcua.AuthCertificate())
		connectionParams = append(connectionParams, opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeCertificate))
	default:
		connectionParams = append(connectionParams, opcua.AuthAnonymous())
		connectionParams = append(connectionParams, opcua.SecurityFromEndpoint(ep, ua.UserTokenTypeAnonymous))
	}

	return connectionParams
}

func CreateClientConnection(ep string, options []opcua.Option) (*opcua.Client, error) {

	return opcua.NewClient(ep, options...)

}
