package opcua_monitor

import (
	"context"
	"fmt"
	"sync"
	"time"

	exporter "github.com/doteich/OPC-UA-Logger/exporters"
	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/doteich/OPC-UA-Logger/setup"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/monitor"
	"github.com/gopcua/opcua/ua"
)

func MonitorItems(ctx context.Context, nodeMonitor *monitor.NodeMonitor, interval int, lag time.Duration, wg *sync.WaitGroup, nodes []setup.NodeObject) {

	sub, err := nodeMonitor.Subscribe(
		ctx,
		&opcua.SubscriptionParameters{
			Interval: time.Duration(interval) * time.Second,
			Priority: 10,
		},
		func(s *monitor.Subscription, msg *monitor.DataChangeMessage) {
			if msg.Error != nil {
				logging.LogError(msg.Error, "Error with received subscription message", "opcua")
			} else {
				go exporter.PublishData(msg.NodeID.String(), msg.Value.Value(), msg.SourceTimestamp)
			}
			time.Sleep(lag)
		},
	)

	for _, node := range nodes {
		_, err := sub.AddMonitorItemsWithContext(ctx, monitor.Request{NodeID: ua.MustParseNodeID(node.NodeId), MonitoringParameters: &ua.MonitoringParameters{QueueSize: 1, SamplingInterval: 1000}, MonitoringMode: ua.MonitoringModeReporting})
		if err != nil {
			logging.LogError(err, "Error while adding node to subscription- node:"+node.NodeId, "opcua")
		}
	}

	if err != nil {
		logging.LogError(err, "Error with subscription", "opcua")
		return
	}

	defer cleanup(ctx, sub, wg)

	<-ctx.Done()

}

func cleanup(ctx context.Context, sub *monitor.Subscription, wg *sync.WaitGroup) {

	fmt.Printf("stats: sub=%d delivered=%d dropped=%d", sub.SubscriptionID(), sub.Delivered(), sub.Dropped())
	sub.Unsubscribe(ctx)
	wg.Done()
}
