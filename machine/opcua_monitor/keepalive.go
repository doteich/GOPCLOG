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

func StartKeepAlive(ctx context.Context, nodeMonitor *monitor.NodeMonitor, lag time.Duration, tChan chan bool) {

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
	Subs[id] = s_struct{sub: sub, tChan: tChan}

	defer cleanup(sub)

	for {
		select {
		case <-tChan:
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

	time.Sleep(10 * time.Second)

	for {
		select {
		case <-ctx.Done():
			for id, s_struct := range Subs {
				s_struct.tChan <- true
				if err := s_struct.sub.Unsubscribe(context.Background()); err != nil {
					logging.LogError(err, "error unsubscribing at shutdown", "opcua")
				}

				delete(Subs, id)
			}

			return
		default:

			if time.Since(last_keepalive) > 5*time.Second {

				logging.LogGeneric("info", "subscriptions timed out - reinit sub", "opcua")

				for id, s_struct := range Subs {
					if err := s_struct.sub.Unsubscribe(ctx); err != nil {
						logging.LogError(err, "error unsubscribing", "opcua")
					}
					delete(Subs, id)
				}
				go StartKeepAlive(ctx, NodeMonitor, 1*time.Second, make(chan bool))
				go MonitorItems(ctx, NodeMonitor, iv, 1000, nodes, make(chan bool))

			}

		}
	}
}
