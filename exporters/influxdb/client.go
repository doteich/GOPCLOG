package influxdb

import (
	"context"
	"time"

	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

var ctx context.Context
var client influxdb2.Client
var writeApi api.WriteAPIBlocking

func CreateConnection(namespace, org, bucket, connectionString, token string) {
	ctx = context.Background()

	client = influxdb2.NewClient(connectionString, token)
	writeApi = client.WriteAPIBlocking(org, bucket)
}

func WriteData(nodeId string, nodeName string, value interface{}, timestamp time.Time, logName string, server string, datatype string, namespace string) {
	tags := map[string]string{
		"nodeId":    nodeId,
		"nodeName":  nodeName,
		"server":    server,
		"datatype":  datatype,
		"namespace": namespace,
	}
	fields := map[string]interface{}{
		nodeName: value,
	}
	point := write.NewPoint(logName, tags, fields, timestamp)
	if err := writeApi.WritePoint(context.Background(), point); err != nil {
		logging.LogError(err, "Failed to write point", "influxdb")
	}
}
