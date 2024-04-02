package websockets

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"github.com/gorilla/websocket"
)

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

var wsManager Manager

type Manager struct {
	clients ClientList
	sync.RWMutex
}

func NewManager() {

}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {

	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		return true // Remove in Production
	}

	connection, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	newClient := SetNewClient(connection, m)

	m.addClient(newClient)

	logging.LogGeneric("debug", "new websocket client added", "websocket")

	go newClient.writeMessages()
	go newClient.readMessages()

}

func (m *Manager) addClient(c *Client) {
	m.Lock()

	defer m.Unlock()

	m.clients[c] = true
}

func (m *Manager) removeClient(c *Client) {
	m.Lock()
	defer m.Unlock()

	if _, found := m.clients[c]; found {
		c.connection.Close()
		delete(m.clients, c)
	}
}

func (m *Manager) BroadcastMessage(payload []byte) {
	for client := range m.clients {
		client.OutChan <- payload
	}
}
