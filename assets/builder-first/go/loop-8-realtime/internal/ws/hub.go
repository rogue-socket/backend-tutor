// Package ws — in-memory hub for WebSocket clients on a single instance.
//
// Phase A: hub broadcasts among locally-connected clients only.
// Phase B (bridge.go): hub also receives from Redis pub/sub, broadcasts to locals.

package ws

import (
	"encoding/json"
	"log"
	"sync"
)

// Event is the broadcast envelope. Keep it small — every connected client
// receives a copy.
type Event struct {
	ID   string          `json:"id"`     // for dedup across the Redis bridge
	Type string          `json:"type"`   // e.g., "link.created"
	Data json.RawMessage `json:"data"`
}

// Client represents a connected WebSocket. The send channel is buffered;
// the writer goroutine reads from it and writes to the underlying conn.
//
// If `send` fills (slow client), the Hub drops events for this client rather
// than blocking — see Hub.Broadcast.
type Client struct {
	UserID int64
	send   chan []byte // serialized Event JSON
	close  chan struct{}
}

// Hub maintains the set of clients and broadcasts events.
type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]struct{}

	// Phase B hook: when set, Hub.Broadcast also publishes to Redis.
	publishToBridge func(Event)
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*Client]struct{})}
}

// SetBridge wires up Phase B publish-to-Redis.
func (h *Hub) SetBridge(publish func(Event)) {
	h.publishToBridge = publish
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	delete(h.clients, c)
	h.mu.Unlock()
	close(c.close)
}

// Broadcast sends an event to all connected clients (and to the Redis bridge
// if wired).
//
// TODO: fan out the event payload via each client's send channel. Use a
// non-blocking send (select with default) so a slow client drops the event
// instead of blocking the broadcast loop:
//
//   select {
//   case c.send <- payload:
//   default:
//       // slow consumer — log and drop, or close the connection
//   }
//
// This is backpressure handling. The alternative — blocking the broadcast on
// a slow client — turns one slow client into a service-wide latency spike.
func (h *Hub) Broadcast(e Event) {
	payload, err := json.Marshal(e)
	if err != nil {
		log.Printf("hub: marshal event: %v", err)
		return
	}

	// Local fanout.
	h.mu.RLock()
	clients := make([]*Client, 0, len(h.clients))
	for c := range h.clients {
		clients = append(clients, c)
	}
	h.mu.RUnlock()

	for _, c := range clients {
		_ = c
		_ = payload
		// TODO: non-blocking send.
	}

	// Bridge fanout.
	if h.publishToBridge != nil {
		h.publishToBridge(e)
	}
}

// BroadcastFromBridge is called by the Redis bridge on incoming events.
// Same as Broadcast but does NOT republish (would loop).
func (h *Hub) BroadcastFromBridge(e Event) {
	old := h.publishToBridge
	h.publishToBridge = nil
	defer func() { h.publishToBridge = old }()
	h.Broadcast(e)
}
