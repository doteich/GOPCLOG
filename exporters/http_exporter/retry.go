package http_exporter

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/doteich/OPC-UA-Logger/exporters/metrics_exporter"
)

type BufferEntry struct {
	Message string  `json:"msg"`
	Payload Payload `json:"payload"`
}

func Resend() {
	retryInProgress = true

	messages := ReadLogFile()

	for _, obj := range messages {
		payload := obj.Payload
		PostLoggedData(payload.NodeId, payload.NodeName, payload.Value, payload.Timestamp, payload.LogName, payload.Server, payload.DataType)
		metrics_exporter.Failed_requests.WithLabelValues(path).Add(-1)
	}
	bufferSize = 0

	retryInProgress = false

}

func ReadLogFile() []BufferEntry {

	dat, err := os.ReadFile("tmp/logs/buffer.json")

	if err != nil {
		logging.Logs.Error(fmt.Sprint(err))
	}

	jsonString := string(dat)
	jsonString = jsonString[0 : len(jsonString)-1]

	jsonString = "[" + jsonString + "]"

	var jsonArr []BufferEntry

	if err := json.Unmarshal([]byte(jsonString), &jsonArr); err != nil {
		logging.Logs.Error(fmt.Sprint(err))
	}

	os.Truncate("tmp/logs/buffer.json", 0)

	return jsonArr
}
