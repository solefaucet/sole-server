package list

import (
	"container/list"
	"sync"

	"github.com/solefaucet/sole-server/services/hub"
)

// Hub stores connections in doubly linked list
type Hub struct {
	conns  *list.List
	rwlock sync.RWMutex
}

var _ hub.Hub = &Hub{}

// New creates a new Hub
func New() *Hub {
	return &Hub{
		conns: list.New(),
	}
}

// PutConn puts a new connection into list
func (h *Hub) PutConn(c hub.Conn) {
	h.rwlock.Lock()
	defer h.rwlock.Unlock()
	h.conns.PushFront(c)
}

// Broadcast broadcast message to all connections
func (h *Hub) Broadcast(m []byte) {
	go h.broadcast(m)
}

// Len returns number of active connections
func (h *Hub) Len() int {
	h.rwlock.RLock()
	defer h.rwlock.RUnlock()
	return h.conns.Len()
}

// broadcast and remove errored connections
func (h *Hub) broadcast(raw []byte) {
	h.rwlock.Lock()
	defer h.rwlock.Unlock()

	var next *list.Element
	for e := h.conns.Front(); e != nil; e = next {
		next = e.Next()
		conn := e.Value.(hub.Conn)
		if err := conn.Write(raw); err != nil {
			conn.Close()
			h.conns.Remove(e)
		}
	}
}
