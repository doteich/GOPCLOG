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
var backupPath string

/*
const path string = "http://localhost:3001/log"          // Main Path for setting http calls
const backupPath string = "http://localhost:3001/backup" // Route if main path is not reachable
*/

type Payload struct {
	NodeId    string      `json:"nodeid"`
	NodeName  string      `json:"nodeName"`
	Value     interface{} `json:"value"` // Data type could be either uint32,string, float32, int16
	Timestamp time.Time   `json:"timestamp"`
	LogName   string      `json:"logName"`
	Server    string      `json:"server"`
}

func InitRoutes(p string, b string) {
	path = p
	backupPath = b
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
		Logs.Error("Unable to reach "+path, zap.ByteString("payload", jsonPayload))
	} else {
		statusCode := resp.StatusCode

		if statusCode > 399 {
			req, err := http.NewRequest("POST", backupPath, bytes.NewBuffer(jsonPayload))
			if err != nil {
				fmt.Println(err)
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(resp)
			}

		}
		defer resp.Body.Close()
	}

}
