package main

import (
	"context"
	"time"

	"project/internal/platform/dbping"
)

// One file per background job; starts dbping and returns its stop func
func registerBackground(ctx context.Context, exec dbping.Execer) (stop func()) {
	return startDBPing(ctx, exec, 60*time.Second)
}

// Runs the DB ping loop on its own ticker; stop cancels it and frees the ticker
func startDBPing(ctx context.Context, exec dbping.Execer, interval time.Duration) (stop func()) {
	pingCtx, cancel := context.WithCancel(ctx)
	ticker := time.NewTicker(interval)

	go func() {
		dbping.Run(pingCtx, exec, ticker.C)
		ticker.Stop()
	}()

	return cancel
}
