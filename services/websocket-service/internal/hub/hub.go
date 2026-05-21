package hub

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	UserID   string
	TenantID string
}

type Hub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]*Client
}

func New() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]*Client),
	}
}

func (h *Hub) Add(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c.Conn] = c
}

func (h *Hub) Remove(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, conn)
}

func (h *Hub) Broadcast(messageType int, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for conn := range h.clients {
		_ = conn.WriteMessage(messageType, data)
	}
}

func (h *Hub) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
