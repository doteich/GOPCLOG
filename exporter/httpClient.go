package exporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var path string

type Payload struct {
	NodeId    string      `json:"nodeid"`
	NodeName  string      `json:"nodeName"`
	Value     interface{} `json:"value"` // Data type could be either uint32,string, float32, int16
	Timestamp time.Time   `json:"timestamp"`
	LogName   string      `json:"logName"`
	Server    string      `json:"server"`
}

func InitRoutes(p string) {
	path = p
}

func PostLoggedData(nodeId string, nodeName string, value interface{}, timestamp time.Time, logName string, server string) {

	client := http.Client{}

	newPayload := Payload{NodeId: nodeId, NodeName: nodeName, Value: value, Timestamp: timestamp, LogName: logName, Server: server}

	jsonPayload, err := json.Marshal(newPayload)

	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequest("POST", path, bytes.NewBuffer(jsonPayload))

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {

		Buffer.Error("Failed to reach "+path, zap.Any("payload", &newPayload))
		Logs.Error("Connection refused for " + path)

	} else {
		statusCode := resp.StatusCode
		if statusCode > 399 {
			Buffer.Error("Target returned a bad response: "+fmt.Sprint(statusCode), zap.Any("payload", &newPayload))
			Logs.Error("Target returned a bad response: " + fmt.Sprint(statusCode))

		}
		defer resp.Body.Close()
	}

}
