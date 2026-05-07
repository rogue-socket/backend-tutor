// Loop 6 — cache-aside on the "popular links" endpoint.
//
// The single-flight group is the load-bearing primitive: ensures that for any
// given cache key, only ONE in-flight loader runs, even if 1000 concurrent
// requests arrive on a cache miss.

package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

const (
	popularKey = "popular:v1"        // bump v1→v2 to invalidate everything
	popularTTL = 60 * time.Second
)

// Link mirrors the application's Link type. In your real code, import it.
type Link struct {
	ID    int64    `json:"id"`
	URL   string   `json:"url"`
	Title string   `json:"title"`
	Tags  []string `json:"tags,omitempty"`
}

type Popular struct {
	rdb *redis.Client
	sfg singleflight.Group
	// loader runs the actual DB query when there's a cache miss.
	// Inject it so tests can swap.
	loader func(ctx context.Context) ([]Link, error)
}

func NewPopular(rdb *redis.Client, loader func(ctx context.Context) ([]Link, error)) *Popular {
	return &Popular{rdb: rdb, loader: loader}
}

// GetOrCompute returns the popular links list, using the cache when warm and
// loading from the underlying source on miss. Concurrent callers on a miss
// share a single loader execution.
//
// TODO 1: try to read from Redis (rdb.Get).
// TODO 2: on hit, json.Unmarshal and return.
// TODO 3: on miss, sfg.Do(popularKey, ...) — inside the closure:
//           a. call p.loader(ctx)
//           b. on success, json.Marshal and rdb.Set with popularTTL
//           c. return the slice
// TODO 4: assert the result type and return.
func (p *Popular) GetOrCompute(ctx context.Context) ([]Link, error) {
	_ = json.Marshal
	_ = json.Unmarshal
	return nil, errors.New("TODO: implement GetOrCompute")
}

// Invalidate bumps or deletes the cache. Use the version-bump strategy
// (set a new key) when you can't atomically delete; use direct delete when
// you can.
//
// Simplest correct implementation: rdb.Del(ctx, popularKey).
func (p *Popular) Invalidate(ctx context.Context) error {
	return p.rdb.Del(ctx, popularKey).Err()
}
