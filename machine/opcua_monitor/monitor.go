package opcua_monitor

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	exporter "github.com/doteich/OPC-UA-Logger/exporters"
	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/doteich/OPC-UA-Logger/setup"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/monitor"
)

func MonitorItems(ctx context.Context, nodeMonitor *monitor.NodeMonitor, interval int, lag time.Duration, wg *sync.WaitGroup, nodes []setup.NodeObject) {
	nodeArr := make([]string, 0)

	for _, node := range nodes {
		nodeArr = append(nodeArr, node.NodeId)
	}

	sub, err := nodeMonitor.Subscribe(
		ctx,
		&opcua.SubscriptionParameters{
			Interval: time.Duration(interval) * time.Second,
		},
		func(s *monitor.Subscription, msg *monitor.DataChangeMessage) {
			if msg.Error != nil {
				logging.LogError(msg.Error, "Error with received subscription message", "opcua")
			} else {
				go exporter.PublishData(msg.NodeID.String(), msg.Value.Value(), msg.SourceTimestamp)
				// log.Printf("[callback] sub=%d ts=%s node=%s value=%v", s.SubscriptionID(), msg.SourceTimestamp.UTC().Format(time.RFC3339), msg.NodeID, msg.Value.Value())
			}
			time.Sleep(lag)
		},
		nodeArr...)

	if err != nil {
		log.Fatal(err)
	}

	defer cleanup(ctx, sub, wg)

	<-ctx.Done()

}

func cleanup(ctx context.Context, sub *monitor.Subscription, wg *sync.WaitGroup) {

	fmt.Printf("stats: sub=%d delivered=%d dropped=%d", sub.SubscriptionID(), sub.Delivered(), sub.Dropped())
	sub.Unsubscribe(ctx)
	wg.Done()
}
