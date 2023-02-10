package websockets

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	OutChan    chan []byte
}

func SetNewClient(connection *websocket.Conn, manager *Manager) *Client {
	return &Client{connection: connection, manager: manager, OutChan: make(chan []byte)}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()
	for {
		message, ok := <-c.OutChan
		if !ok {
			c.connection.WriteMessage(websocket.CloseMessage, nil)
			fmt.Println("Error with received message")
			return
		}
		c.connection.WriteMessage(websocket.TextMessage, message)
	}
}
