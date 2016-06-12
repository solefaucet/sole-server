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
	go func() {
		// construct slice of connections
		elems := h.elems()

		// broadcast message
		// collect errored connections
		errElems := []*list.Element{}
		wg := &sync.WaitGroup{}
		for i := range elems {
			wg.Add(1)
			go func(e *list.Element) {
				defer wg.Done()
				if err := e.Value.(hub.Conn).Write(m); err != nil {
					errElems = append(errElems, e)
				}
			}(elems[i])
		}
		wg.Wait()

		// remove errored connections
		h.removeErrElems(errElems)
	}()
}

// Len returns number of active connections
func (h *Hub) Len() int {
	h.rwlock.RLock()
	defer h.rwlock.RUnlock()
	return h.conns.Len()
}

// iterate list, construct slice of connections
func (h *Hub) elems() []*list.Element {
	h.rwlock.RLock()
	defer h.rwlock.RUnlock()
	elems := make([]*list.Element, h.conns.Len())
	for e, i := h.conns.Front(), 0; e != nil; e, i = e.Next(), i+1 {
		elems[i] = e
	}

	return elems
}

// remove errored connections from list
func (h *Hub) removeErrElems(errElems []*list.Element) {
	h.rwlock.Lock()
	defer h.rwlock.Unlock()

	for i := range errElems {
		errElems[i].Value.(hub.Conn).Close()
		h.conns.Remove(errElems[i])
	}
}
