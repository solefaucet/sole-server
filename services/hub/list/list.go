package list

import (
	"container/list"
	"sync"

	"github.com/freeusd/solebtc/services/hub"
)

// Hub stores connections in doubly linked list
type Hub struct {
	conns  *list.List
	rwlock sync.RWMutex

	onSendError func(error)
}

var _ hub.Hub = &Hub{}

// New creates a new Hub
func New(onSendError func(error)) *Hub {
	return &Hub{
		conns:       list.New(),
		onSendError: onSendError,
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
	removeErrConns := func(errConns []*list.Element) {
		if len(errConns) == 0 {
			return
		}

		h.rwlock.Lock()
		defer h.rwlock.Unlock()

		for i := range errConns {
			h.conns.Remove(errConns[i])
		}
	}

	broadcast := func(msg []byte, onError func(error, *list.Element)) {
		errConns := []*list.Element{}
		defer removeErrConns(errConns)

		h.rwlock.RLock()
		defer h.rwlock.RUnlock()
		for e := h.conns.Front(); e != nil; e = e.Next() {
			if err := e.Value.(hub.Conn).Write(msg); err != nil {
				// NOTE:
				// DO NOT remove element from list here
				// it can cause race condition because it's a read lock
				onError(err, e)
			}
		}
	}

	go func(msg []byte) {
		// collect errored connections
		errConns := []*list.Element{}
		onError := func(err error, e *list.Element) {
			errConns = append(errConns, e)
			h.onSendError(err)
		}

		broadcast(msg, onError)
		removeErrConns(errConns)
	}(m)
}

// Len returns number of active connections
func (h *Hub) Len() int {
	h.rwlock.RLock()
	defer h.rwlock.RUnlock()
	return h.conns.Len()
}
