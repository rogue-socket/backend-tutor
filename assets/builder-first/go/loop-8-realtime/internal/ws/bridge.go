// Phase B: Redis pub/sub bridge between instances.
//
// One channel ("links.events") that every instance both publishes to and
// subscribes from. Local broadcasts → publish; remote messages → local
// broadcast (without republish).

package ws

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/redis/go-redis/v9"
)

const channel = "links.events"

type Bridge struct {
	rdb        *redis.Client
	hub        *Hub
	instanceID string // included in Event payload to detect "I sent this"
	seenIDs    map[string]struct{} // tiny dedup window; cap size
}

func NewBridge(rdb *redis.Client, hub *Hub, instanceID string) *Bridge {
	b := &Bridge{rdb: rdb, hub: hub, instanceID: instanceID, seenIDs: make(map[string]struct{})}
	hub.SetBridge(func(e Event) { b.publish(e) })
	return b
}

// Run subscribes to the channel and forwards to the local hub. Block until ctx
// is cancelled.
//
// TODO:
//   1. rdb.Subscribe(ctx, channel)
//   2. for each message:
//        json.Unmarshal payload to Event
//        if seen → skip (we sent it)
//        else → mark seen, BroadcastFromBridge
//   3. on ctx.Done → unsubscribe, close
func (b *Bridge) Run(ctx context.Context) error {
	_ = b.hub.BroadcastFromBridge
	_ = json.Unmarshal
	return errors.New("TODO: implement Run")
}

// publish is called by the Hub on every local broadcast.
//
// TODO: marshal Event to JSON, then b.rdb.Publish(ctx, channel, payload).
// Until you implement, this is a no-op so the rest of the code can compile.
func (b *Bridge) publish(e Event) {
	_ = b.rdb // remove once you wire up Publish
	log.Printf("bridge: TODO publish event %s (type=%s)", e.ID, e.Type)
}
