package opcua_monitor

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	exporter "github.com/doteich/OPC-UA-Logger/exporters"
	"github.com/doteich/OPC-UA-Logger/setup"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/monitor"
)

func errorHandler(err error) {
	panic(err)
}

func MonitorItems(ctx context.Context, nodeMonitor *monitor.NodeMonitor, interval time.Duration, lag time.Duration, wg *sync.WaitGroup, nodes []setup.NodeObject) {
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

				go exporter.PublishData(id, msg.Value.Value(), msg.SourceTimestamp)
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
