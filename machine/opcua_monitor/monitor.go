package opcua_monitor

import (
	"context"
	"errors"
	"fmt"
	"time"

	exporter "github.com/doteich/OPC-UA-Logger/exporters"
	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/doteich/OPC-UA-Logger/setup"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/monitor"
	"github.com/gopcua/opcua/ua"
)

func MonitorItems(pctx context.Context, ctx context.Context, m *monitor.NodeMonitor, interval int, lag time.Duration, nodes []setup.NodeObject) {

	sub, err := m.Subscribe(
		ctx,
		&opcua.SubscriptionParameters{
			Interval: time.Duration(interval) * time.Second,
			Priority: 10,
		},
		func(s *monitor.Subscription, msg *monitor.DataChangeMessage) {
			if msg.Error != nil {
				logging.LogError(msg.Error, "Error with received subscription message", "opcua")
			} else {
				if msg.Status == ua.StatusBad {
					logging.LogError(errors.New("received status Code bad"), "error with received subscription message", "opcua")
				} else {
					go exporter.PublishData(msg.NodeID.String(), msg.Value.Value(), msg.SourceTimestamp)
				}

			}
			time.Sleep(lag)
		},
	)

	if err != nil {
		logging.LogError(err, "Error with subscription", "opcua")
		return
	}

	for _, node := range nodes {
		_, err := sub.AddMonitorItems(ctx, monitor.Request{NodeID: ua.MustParseNodeID(node.NodeId), MonitoringParameters: &ua.MonitoringParameters{QueueSize: 1, SamplingInterval: 1000}, MonitoringMode: ua.MonitoringModeReporting})
		if err != nil {
			logging.LogError(err, "Error while adding node to subscription- node:"+node.NodeId, "opcua")
			continue
		}
	}

	id := sub.SubscriptionID()

	Subs[id] = sub

	logging.LogGeneric("info", "Starting Subscription with id: "+fmt.Sprint(id), "opcua")

	defer TerminateSub(pctx, sub, id)

	<-ctx.Done()

}

func TerminateSub(ctx context.Context, s *monitor.Subscription, id uint32) {
	logging.LogGeneric("info", fmt.Sprintf("terminating subscription with id: %d - delivered: %d - dropped: %d", id, s.Delivered(), s.Dropped()), "terminate_sub")
	delete(Subs, id)
	s.Unsubscribe(ctx)
}
