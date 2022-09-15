package setup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const path string = "http://localhost:3001/log"          // Main Path for setting http calls
const backupPath string = "http://localhost:3001/backup" // Route if main path is not reachable

type Payload struct {
	NodeId    string      `json:"nodeid"`
	NodeName  string      `json:"nodeName"`
	Value     interface{} `json:"value"` // Data type could be either uint32,string, float32, int16
	Timestamp time.Time   `json:"timestamp"`
	LogName   string      `json:"logName"`
	Server    string      `json:"server"`
}

func PostLoggedData(nodeId string, value interface{}, timestamp time.Time) {
	client := http.Client{}

	c := SetConfig()

	var nodeName string

	for _, obj := range c.Nodes {
		if obj.NodeId == nodeId {
			nodeName = obj.NodeName
		}
	}

	newPayload := Payload{NodeId: nodeId, NodeName: nodeName, Value: value, Timestamp: timestamp, LogName: c.LoggerConfig.Name, Server: c.ClientConfig.Url}

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
		fmt.Println(err)
	}

	statusCode := resp.StatusCode

	if statusCode > 399 {
		req, err := http.NewRequest("POST", backupPath, bytes.NewBuffer(jsonPayload))
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "application/json")
		client.Do(req)
	}
	defer resp.Body.Close()

	fmt.Println(resp)
}
