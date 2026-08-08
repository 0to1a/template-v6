// Placeholder background service; a template for future periodic jobs
package dbping

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// Minimal surface needed for the ping; *pgxpool.Pool satisfies this
type Execer interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// Runs until ctx is canceled. First ping waits for the first tick, so
// nothing fires immediately at startup. Failures are logged, not fatal.
func Run(ctx context.Context, exec Execer, ticks <-chan time.Time) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticks:
			if _, err := exec.Exec(ctx, "SELECT 1"); err != nil {
				log.Printf("dbping: select 1 failed: %v", err)
			}
		}
	}
}
