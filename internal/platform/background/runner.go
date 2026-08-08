// One goroutine per job; no scheduler framework needed for a single ticker
package background

import (
	"context"
	"log"
	"time"
)

// Run calls job on every tick until ctx is done. A job failure is logged
// and never stops the loop, so the next tick still runs.
func Run(ctx context.Context, ticks <-chan time.Time, job func(context.Context) error) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticks:
			if err := job(ctx); err != nil {
				log.Printf("background: job failed: %v", err)
			}
		}
	}
}
