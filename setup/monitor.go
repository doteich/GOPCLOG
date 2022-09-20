package setup

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/monitor"
	"github.com/gopcua/opcua/ua"
)

func errorHandler(err error) {
	panic(err)
}

func ValidateEndpoint(ctx context.Context, endpoint string, policy string, mode string) *ua.EndpointDescription {
	endpoints, err := opcua.GetEndpoints(ctx, endpoint)

	if err != nil {
		fmt.Println(endpoint + policy + mode)
		errorHandler(err)
	}

	ep := opcua.SelectEndpoint(endpoints, policy, ua.MessageSecurityModeFromString(mode))

	if ep == nil {
		panic("No Matching Endpoint Found - Check Configuration")
	}

	return ep

}

func SetClientOptions(config *Config, ep *ua.EndpointDescription) []opcua.Option {

	// basic params
	connectionParams := []opcua.Option{
		opcua.SecurityPolicy(config.ClientConfig.SecurityPolicy),
		opcua.SecurityModeString(config.ClientConfig.SecurityMode),
		opcua.CertificateFile("./certs/cert.pem"),
		opcua.PrivateKeyFile("./certs/private_key.pem"),
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

func CreateClientConnection(ep string, options []opcua.Option) *opcua.Client {

	return opcua.NewClient(ep, options...)

}

func MonitorItems(ctx context.Context, nodeMonitor *monitor.NodeMonitor, interval time.Duration, lag time.Duration, wg *sync.WaitGroup, nodes []NodeObject) {
	ch := make(chan *monitor.DataChangeMessage, 16)
	fmt.Println(nodes)
	nodeArr := make([]string, 0)

	for _, node := range nodes {
		nodeArr = append(nodeArr, node.NodeId)
	}

	sub, err := nodeMonitor.ChanSubscribe(ctx, &opcua.SubscriptionParameters{Interval: interval}, ch, nodeArr...)

	if err != nil {
		errorHandler(err)
	}

	defer cleanup(ctx, sub, wg)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			if msg.Error != nil {
				log.Printf("[channel ] sub=%d error=%s", sub.SubscriptionID(), msg.Error)
			} else {

				id := msg.NodeID.String()
				fmt.Println(id)
				//PostLoggedData(id, msg.Value.Value(), msg.SourceTimestamp)
				log.Printf("[channel ] sub=%d ts=%s node=%s value=%v", sub.SubscriptionID(), msg.SourceTimestamp.UTC().Format(time.RFC3339), msg.NodeID, msg.Value.Value())
			}
			time.Sleep(lag)
		}
	}
}

func cleanup(ctx context.Context, sub *monitor.Subscription, wg *sync.WaitGroup) {
	fmt.Printf("stats: sub=%d delivered=%d dropped=%d", sub.SubscriptionID(), sub.Delivered(), sub.Dropped())
	sub.Unsubscribe(ctx)
	wg.Done()
}
