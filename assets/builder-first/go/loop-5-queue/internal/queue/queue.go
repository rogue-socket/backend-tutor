// Package queue wraps Redis Streams for the notification pipeline.
//
// Two streams: `notifications` (work) and `notifications.dlq` (failures).
// One consumer group: `workers`.

package queue

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	StreamMain    = "notifications"
	StreamDLQ     = "notifications.dlq"
	GroupWorkers  = "workers"
	maxDeliveries = 5
	dedupTTL      = 24 * time.Hour
)

type Client struct {
	rdb *redis.Client
}

func New(addr string) *Client {
	return &Client{
		rdb: redis.NewClient(&redis.Options{Addr: addr}),
	}
}

// EnsureGroup creates the consumer group if it doesn't exist.
// Idempotent — call once at worker startup.
func (c *Client) EnsureGroup(ctx context.Context) error {
	// XGROUP CREATE notifications workers $ MKSTREAM
	// On "BUSYGROUP" error, treat as success.
	err := c.rdb.XGroupCreateMkStream(ctx, StreamMain, GroupWorkers, "$").Err()
	if err != nil && err.Error() == "BUSYGROUP Consumer Group name already exists" {
		return nil
	}
	return err
}

// Enqueue adds a message to the main stream.
// TODO: implement with XAdd; the message body is a map[string]any.
func (c *Client) Enqueue(ctx context.Context, fields map[string]any) error {
	_ = c.rdb.XAdd
	return errors.New("TODO: implement Enqueue")
}

// Message holds a redelivered or fresh message.
type Message struct {
	ID            string
	Fields        map[string]any
	DeliveryCount int
}

// Read blocks for up to `block` waiting for messages.
//
// TODO:
//   - XReadGroup with Group=workers, Consumer=consumerName, Streams=[notifications, ">"], Count=10, Block=block
//   - For each returned message, populate DeliveryCount via XPending lookup (or pass the count from elsewhere)
//   - Return nil, nil on no messages (timeout) — caller loops
func (c *Client) Read(ctx context.Context, consumerName string, block time.Duration) ([]Message, error) {
	return nil, errors.New("TODO: implement Read")
}

// Ack acknowledges a message — caller MUST call after successful processing.
func (c *Client) Ack(ctx context.Context, id string) error {
	return c.rdb.XAck(ctx, StreamMain, GroupWorkers, id).Err()
}

// SeenBefore checks the dedup set; SADD must be called after successful processing.
func (c *Client) SeenBefore(ctx context.Context, key string) (bool, error) {
	return c.rdb.SIsMember(ctx, "notifications:seen", key).Result()
}

func (c *Client) MarkSeen(ctx context.Context, key string) error {
	pipe := c.rdb.TxPipeline()
	pipe.SAdd(ctx, "notifications:seen", key)
	pipe.Expire(ctx, "notifications:seen", dedupTTL)
	_, err := pipe.Exec(ctx)
	return err
}

// SendToDLQ routes a poisoned message to the DLQ stream and acks it on the main.
//
// TODO:
//   - XAdd to StreamDLQ with the original fields plus a `_dlq_reason` field
//   - XAck the original on the main stream
//   - Log + emit a metric (Loop 6 will wire up real metrics)
func (c *Client) SendToDLQ(ctx context.Context, m Message, reason string) error {
	return errors.New("TODO: implement SendToDLQ")
}

// MaxDeliveriesReached is a convenience for the worker loop.
func MaxDeliveriesReached(deliveryCount int) bool {
	return deliveryCount >= maxDeliveries
}
