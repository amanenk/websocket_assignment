package clients_storage

import (
	"github.com/gorilla/websocket"
	"sync"
)

// SafeCounter is safe to use concurrently.
type ClientsStorage struct {
	Mutex   sync.RWMutex
	Clients map[*websocket.Conn]*Client
}

func (c *ClientsStorage) Subscribe(ws *websocket.Conn) {
	c.Mutex.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.Clients[ws].Sub = true
	c.Mutex.Unlock()

}

func (c *ClientsStorage) Unsubscribe(ws *websocket.Conn) {
	c.Mutex.Lock()
	c.Clients[ws].Sub = false
	c.Mutex.Unlock()
}

func (c *ClientsStorage) Delete(conn *websocket.Conn) {
	c.Mutex.Lock()
	delete(c.Clients, conn)
	c.Mutex.Unlock()
}

func (c *ClientsStorage) Add(client *Client) {
	c.Mutex.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.Clients[client.socket] = client
	c.Mutex.Unlock()
}

func (c *ClientsStorage) GetClientsCount() int {
	c.Mutex.RLock()
	clientsCount := len(c.Clients)
	c.Mutex.RUnlock()
	return clientsCount
}

func (c *ClientsStorage) GetClientId() int {
	c.Mutex.RLock()
	clientId := len(c.Clients)
	c.Mutex.RUnlock()
	return clientId
}

var ClientStorage = ClientsStorage{Clients: make(map[*websocket.Conn]*Client)}
