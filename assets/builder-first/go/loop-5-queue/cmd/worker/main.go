// Worker process for Loop 5.
//
// Run:
//   go run ./cmd/worker
//
// Reads from the `notifications` stream, processes each message at-least-once
// with idempotent dedup, routes poison messages to the DLQ.

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	// "loop5/internal/queue"   — adjust to your module path
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// Graceful shutdown on SIGTERM/SIGINT.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigCh
		log.Println("shutdown signal — draining")
		cancel()
	}()

	// TODO 1: build the queue client (queue.New from QUEUE_ADDR env)
	// TODO 2: q.EnsureGroup(ctx) at startup
	// TODO 3: identify this consumer (hostname + pid is fine)
	// TODO 4: the main loop — see structure below

	host, _ := os.Hostname()
	if host == "" {
		host = "worker"
	}
	consumerName := host + "-" + strconv.Itoa(os.Getpid())

	for {
		select {
		case <-ctx.Done():
			log.Println("worker stopped")
			return
		default:
		}

		// msgs, err := q.Read(ctx, consumerName, 5*time.Second)
		// if err != nil { log.Printf("read: %v", err); continue }
		//
		// for _, m := range msgs {
		//     if err := process(ctx, q, m); err != nil {
		//         log.Printf("process %s: %v", m.ID, err)
		//         if queue.MaxDeliveriesReached(m.DeliveryCount) {
		//             _ = q.SendToDLQ(ctx, m, err.Error())
		//         }
		//         // do NOT Ack — let it redeliver via XPENDING
		//         continue
		//     }
		//     _ = q.Ack(ctx, m.ID)
		// }

		_ = consumerName
		_ = time.Second
		log.Println("TODO: implement worker loop body")
		time.Sleep(time.Second)
	}
}

// process handles one message. Returns nil on success.
//
// The pattern:
//   1. Compute a dedup key from the message (e.g., "type:link_id")
//   2. SeenBefore? → ack and skip (this is the at-least-once → idempotent step)
//   3. Do the work (look up link, format the email, send/log it)
//   4. MarkSeen
//
// Ordering matters: MarkSeen must happen AFTER the work, not before. If you
// mark-seen before the work and crash mid-work, the message is treated as done
// even though the side effect didn't occur.
//
// Conversely, if work is itself idempotent (e.g., the email service has its
// own dedup), the dedup set is belt-and-suspenders — fine, just be honest about
// what's load-bearing.
func process(ctx context.Context /*, q *queue.Client, m queue.Message */) error {
	// TODO: implement.
	return nil
}
