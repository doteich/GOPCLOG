package opcua_monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/monitor"
	"github.com/gopcua/opcua/ua"
)

var (
	last_keepalive time.Time
)

func StartKeepAlive(pctx context.Context, ctx context.Context, m *monitor.NodeMonitor) {

	last_keepalive = time.Now()

	sub, err := m.Subscribe(pctx, &opcua.SubscriptionParameters{Interval: 10 * time.Second}, func(s *monitor.Subscription, dcm *monitor.DataChangeMessage) {
		if dcm.Error != nil {
			logging.LogError(fmt.Errorf("error with received keepalive message: %s - nodeid %s", dcm.Error.Error(), dcm.NodeID), "", "keepalive")

		} else if dcm.Status != ua.StatusOK {
			logging.LogError(fmt.Errorf("received bad status for keepalive message: %s - nodeid %s", dcm.Value.StatusCode(), dcm.NodeID), "", "keepalive")

		} else {
			last_keepalive = time.Now()
		}
	})

	if err != nil {
		logging.LogError(fmt.Errorf("error while creating subscription: %s", err.Error()), "", "keepalive")
		return
	}

	sub.AddMonitorItems(pctx, monitor.Request{NodeID: ua.MustParseNodeID("i=2258"), MonitoringMode: ua.MonitoringModeReporting, MonitoringParameters: &ua.MonitoringParameters{DiscardOldest: true, QueueSize: 1}})

	id := sub.SubscriptionID()
	Subs[id] = sub
	defer TerminateSub(ctx, sub, id)

	logging.LogGeneric("info", "Starting keepalive with id: "+fmt.Sprint(id), "opcua")

	<-ctx.Done()
}
