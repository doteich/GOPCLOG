package opcua_monitor

import (
	"context"
	"sync"
	"time"

	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/doteich/OPC-UA-Logger/setup"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/monitor"
)

var (
	last_keepalive time.Time
)

func StartKeepAlive(ctx context.Context, nodeMonitor *monitor.NodeMonitor, lag time.Duration) {

	last_keepalive = time.Now()

	ch := make(chan *monitor.DataChangeMessage, 1)

	node := make([]string, 0)
	node = append(node, "i=2258")

	logging.LogGeneric("info", "Starting Keepalive", "opcua")

	sub, err := nodeMonitor.ChanSubscribe(ctx, &opcua.SubscriptionParameters{Interval: 10 * time.Second, Priority: 1}, ch, node...)

	if err != nil {
		logging.LogError(err, "Error starting the subscription for keepalive", "opcua")
		return
	}

	id := sub.SubscriptionID()
	Subs[id] = sub

	defer cleanup(ctx, sub)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			if msg.Error != nil {
				logging.LogError(msg.Error, "Error with received keepalive message", "opcua")
			} else {
				last_keepalive = time.Now()
			}
			time.Sleep(lag)
		}
	}

}

func MonitorSubscriptions(ctx context.Context, wg *sync.WaitGroup, iv int, nodes []setup.NodeObject) {

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:

			if time.Since(last_keepalive) > 60*time.Second {

				logging.LogGeneric("info", "subscriptions timed out - reinit sub", "opcua")

				for id, sub := range Subs {
					if err := sub.Unsubscribe(ctx); err != nil {
						logging.LogError(err, "error unsubscribing", "opcua")
					}
					delete(Subs, id)
				}
				go StartKeepAlive(ctx, NodeMonitor, 1*time.Second)
				go MonitorItems(ctx, NodeMonitor, iv, 1000, nodes)

			}

			time.Sleep(60 * time.Second)
		}
	}
}
