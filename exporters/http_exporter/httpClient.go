package http_exporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/doteich/OPC-UA-Logger/exporters/metrics_exporter"
)

var path string

const retryThreshold uint8 = 10

var numberSuccessfulMessages uint8 = 0

var bufferSize int = 0

var retryInProgress bool = false

type Payload struct {
	NodeId    string      `json:"nodeId"`
	NodeName  string      `json:"nodeName"`
	Value     interface{} `json:"value"` // Data type could be either uint32,string, float32, int16
	Timestamp time.Time   `json:"timestamp"`
	LogName   string      `json:"logName"`
	Server    string      `json:"server"`
	DataType  string      `json:"dataType"`
}

func InitRoutes(p string) {
	path = p
}

func PostLoggedData(nodeId string, nodeName string, value interface{}, timestamp time.Time, logName string, server string, datatype string) {

	client := http.Client{}

	newPayload := Payload{NodeId: nodeId, NodeName: nodeName, Value: value, Timestamp: timestamp, LogName: logName, Server: server, DataType: datatype}

	jsonPayload, err := json.Marshal(newPayload)

	if err != nil {
		logging.LogError(err, "Failed to marshal Payload", "http_exporter")
	}

	req, err := http.NewRequest("POST", path, bytes.NewBuffer(jsonPayload))

	if err != nil {
		logging.LogError(err, "Failed to construct http Request", "http_exporter")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {

		metrics_exporter.Failed_requests.WithLabelValues(path).Inc()
		logging.LogError(err, "Connection refused for "+path, "http_exporter")

	} else {
		defer resp.Body.Close()
		statusCode := resp.StatusCode
		if statusCode > 399 {
			logging.LogError(nil, "Target returned a bad response: "+fmt.Sprint(statusCode), "http_exporter")
			metrics_exporter.Failed_requests.WithLabelValues(path).Inc()
		}

	}

}
