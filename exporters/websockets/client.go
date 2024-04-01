package websockets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/doteich/OPC-UA-Logger/setup"
	"github.com/gopcua/opcua/ua"
	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	OutChan    chan []byte
}

type InboundMessageEvent struct {
	Operation string `json:"operation"`
}

var (
	pongWait     = 10 * time.Second
	pingInterval = 9 * time.Second
)

func SetNewClient(connection *websocket.Conn, manager *Manager) *Client {
	return &Client{connection: connection, manager: manager, OutChan: make(chan []byte)}
}

func (c *Client) writeMessages() {

	ticker := time.NewTicker(pingInterval)

	defer func() {
		ticker.Stop()
		c.manager.removeClient(c)
	}()
	for {

		select {

		case message, ok := <-c.OutChan:
			if !ok {
				c.connection.WriteMessage(websocket.CloseMessage, nil)
				fmt.Println("Error with received message")
				return
			}

			c.connection.WriteMessage(websocket.TextMessage, message)

		case <-ticker.C:
			err := c.connection.WriteMessage(websocket.PingMessage, []byte{})

			if err != nil {
				return
			}
		}

	}
}

func (c *Client) readMessages() {

	defer func() {
		c.manager.removeClient(c)
	}()

	c.connection.SetReadDeadline(time.Now().Add(pongWait))

	c.connection.SetPongHandler(c.pongHandler)
	for {

		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error reading message: %v", err)
			}

			break
		}

		var inboundMsg InboundMessageEvent

		err = json.Unmarshal(payload, &inboundMsg)

		if err != nil {
			fmt.Println(err)

		}

		if inboundMsg.Operation == "bulk_read" {
			allValues := TriggerBulkRead()
			var newPayload Payload

			for key, val := range allValues {
				metaData, dType, _ := findNodeDetails(key, val)

				newPayload = Payload{NodeId: key, Value: val, Timestamp: time.Now(), NodeName: metaData.NodeName, LogName: setup.PubConfig.LoggerConfig.Name, Server: setup.PubConfig.ClientConfig.Url, DataType: dType}

				bytesArr, _ := json.Marshal(newPayload)
				c.OutChan <- bytesArr
				delete(allValues, key)
			}

		}

	}
}

func (c *Client) pongHandler(PongMsg string) error {
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

func ReadNodes(nodeId string) (interface{}, error) {

	id, _ := ua.ParseNodeID(nodeId)

	obj := &ua.ReadRequest{
		NodesToRead: []*ua.ReadValueID{
			{NodeID: id},
		},
	}
	resp, err := ws_opcclient.Read(context.Background(), obj)

	if err != nil {
		fmt.Printf("Error while reading %s", nodeId)
		return nil, err
	}

	if resp.Results[0].Status == ua.StatusBad || resp.Results[0].Value == nil {
		return nil, errors.New("received status code bad while reading")
	}

	return resp.Results[0].Value.Value(), nil

}

func TriggerBulkRead() map[string]interface{} {

	idMap := make(map[string]interface{})

	for _, node := range setup.PubConfig.Nodes {

		val, err := ReadNodes(node.NodeId)

		if err != nil {
			continue
		}

		idMap[node.NodeId] = val

	}

	return idMap
}

func findNodeDetails(nodeId string, iface interface{}) (setup.NodeObject, string, error) {

	var dataType string

	switch iface.(type) {
	case int:
		dataType = "int"
	case int8:
		dataType = "int8"
	case int16:
		dataType = "int16"
	case int32:
		dataType = "int32"
	case uint8:
		dataType = "uint8"
	case uint16:
		dataType = "uint16"
	case uint32:
		dataType = "uint32"
	case float32:
		dataType = "float32"
	case float64:
		dataType = "float64"
	case string:
		dataType = "string"
	case bool:
		dataType = "bool"
	}

	for _, node := range setup.PubConfig.Nodes {
		if nodeId == node.NodeId {
			return node, dataType, nil
		}
	}

	return setup.NodeObject{}, "", errors.New("node not found")
}
