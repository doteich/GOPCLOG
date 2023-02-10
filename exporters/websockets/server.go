package websockets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Payload struct {
	NodeId    string      `json:"nodeid"`
	NodeName  string      `json:"nodeName"`
	Value     interface{} `json:"value"` // Data type could be either uint32,string, float32, int16
	Timestamp time.Time   `json:"timestamp"`
	LogName   string      `json:"logName"`
	Server    string      `json:"server"`
	DataType  string      `json:"dataType"`
}

func InitWebsockets() {

	RouteHandler()
	http.ListenAndServe(":8080", nil)

}

func RouteHandler() {
	wsManager.clients = make(ClientList)
	http.HandleFunc("/ws", wsManager.ServeWS)
}

func BroadcastToWebsocket(nodeId string, nodeName string, value interface{}, timestamp time.Time, logName string, server string, datatype string) {
	newPayload := Payload{NodeId: nodeId, NodeName: nodeName, Value: value, Timestamp: timestamp, LogName: logName, Server: server, DataType: datatype}
	byteArr, err := json.Marshal(newPayload)

	if err != nil {
		fmt.Println(err)
		return
	}

	wsManager.BroadcastMessage(byteArr)
}
