package opcua_monitor

import (
	"context"
	"sync"
	"time"

	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/monitor"
)

func StartKeepAlive(ctx context.Context, nodeMonitor *monitor.NodeMonitor, lag time.Duration, wg *sync.WaitGroup) {

	ch := make(chan *monitor.DataChangeMessage, 1)

	node := make([]string, 0)
	node = append(node, "i=2258")

	logging.LogGeneric("info", "Starting Keepalive", "opcua")

	sub, err := nodeMonitor.ChanSubscribe(ctx, &opcua.SubscriptionParameters{Interval: 10 * time.Second, Priority: 1}, ch, node...)

	if err != nil {
		logging.LogError(err, "Error starting the subscription for keepalive", "opcua")
	}

	defer cleanup(ctx, sub, wg)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			if msg.Error != nil {
				logging.LogError(msg.Error, "Error with received keepalive message", "opcua")
			}
			time.Sleep(lag)
		}
	}

}
