package opcua_monitor

import (
	"context"
	"fmt"
	"time"

	exporter "github.com/doteich/OPC-UA-Logger/exporters"
	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/doteich/OPC-UA-Logger/exporters/metrics_exporter"
	"github.com/doteich/OPC-UA-Logger/exporters/websockets"
	"github.com/doteich/OPC-UA-Logger/setup"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/monitor"
	"github.com/gopcua/opcua/ua"
)

var (
	opcclient *opcua.Client
	Subs      map[uint32]*monitor.Subscription
)

func CreateOPCUAMonitor(ctx context.Context, config *setup.Config) error {

	ep, err := ValidateEndpoint(ctx, config.ClientConfig.Url, config.ClientConfig.SecurityPolicy, config.ClientConfig.SecurityMode)

	if err != nil {
		return fmt.Errorf("no matching endpoint found - check configuration: %s", err.Error())
	}

	connectionParams := SetClientOptions(config, ep)

	opcclient, err = CreateClientConnection(config.ClientConfig.Url, connectionParams)

	if err != nil {
		return fmt.Errorf("error while creating opc client: %s", err.Error())
	}

	err = opcclient.Connect(ctx)

	if err != nil {
		return fmt.Errorf("error connecting to opcua server: %s", err.Error())
	}

	websockets.InitOPCUARead(opcclient)
	exporter.SetOPCUAClient(opcclient)

	return nil

}

func ValidateEndpoint(ctx context.Context, endpoint string, policy string, mode string) (*ua.EndpointDescription, error) {
	endpoints, err := opcua.GetEndpoints(ctx, endpoint)

	if err != nil {
		return nil, err
	}

	ep := opcua.SelectEndpoint(endpoints, policy, ua.MessageSecurityModeFromString(mode))

	if ep == nil {
		return nil, err
	}

	return ep, nil

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

func InitSubs(pctx context.Context, ctx context.Context, conf *setup.Config) error {
	m, err := monitor.NewNodeMonitor(opcclient)

	if err != nil {
		return err
	}

	go StartKeepAlive(pctx, ctx, m)
	go MonitorItems(pctx, ctx, m, conf.LoggerConfig.Interval, 1000, conf.Nodes)

	return nil

}

func CreateConnectionWatcher(ctx context.Context, t *time.Ticker, conf *setup.Config) {
	Subs = make(map[uint32]*monitor.Subscription)
	var sub_ctx context.Context
	var cancel func()

	sub_ctx, cancel = context.WithCancel(ctx)

	if err := CreateOPCUAMonitor(ctx, conf); err != nil {
		logging.LogError(err, "error initializing watcher", "opcua")
		return
	}
	if err := InitSubs(ctx, sub_ctx, conf); err != nil {
		logging.LogError(err, "error initializing watcher monitor", "opcua")
		return
	}
	last_keepalive = time.Now()
	metrics_exporter.LogReconnects(conf.ClientConfig.Url)

	for {

		select {
		case <-t.C:
			diff := time.Since(last_keepalive).Seconds()

			if diff > 60 {
				logging.LogGeneric("warning", "received last keepalive message more than 60 s ago - reinit subs", "submonitor")
				metrics_exporter.LogReconnects(conf.ClientConfig.Url)
				cancel()
				if err := opcclient.Close(ctx); err != nil {
					logging.LogError(err, "error closing opc ua client connection on reconnect", "opcua")
				}
				time.Sleep(10 * time.Second)
				sub_ctx, cancel = context.WithCancel(ctx)

				if err := CreateOPCUAMonitor(ctx, conf); err != nil {
					logging.LogError(err, "error while reestablishing opc ua client connection", "opcua")
					continue
				}

				InitSubs(ctx, sub_ctx, conf)
			}
		case <-ctx.Done():
			logging.LogGeneric("warning", "shutting down due to context cancel", "opcua")
			cancel()
			opcclient.Close(ctx)
		}

	}

}

// func ConnectionCheck(ctx context.Context, t *time.Ticker, wg *sync.WaitGroup, conf *setup.Config) {
// 	var sub_ctx context.Context
// 	var cancel func()

// 	sub_ctx, cancel = context.WithCancel(ctx)

// 	InitSubs(ctx, sub_ctx, conf)

// 	for {
// 		select {
// 		case <-t.C:
// 			diff := time.Since(last_keepalive).Seconds()

// 			if diff > 60 {
// 				logging.LogGeneric("warning", "received last keepalive message more than 60 s ago - reinit subs", "submonitor")
// 				cancel()
// 				time.Sleep(10 * time.Second)
// 				sub_ctx, cancel = context.WithCancel(ctx)
// 				InitSubs(ctx, sub_ctx, conf)
// 			}

// 		case <-ctx.Done():
// 			cancel()
// 			wg.Done()
// 			return
// 		}
// 	}

// }
