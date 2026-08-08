package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"project/internal/platform/background"
	"project/internal/platform/database"
)

// One file per background service; never a second registry.
// Pings the DB every minute until stopped; failures are logged only.
func registerBackground(ctx context.Context, pool *pgxpool.Pool) (stop func()) {
	pingCtx, cancel := context.WithCancel(ctx)
	ticker := time.NewTicker(time.Minute)
	go background.Run(pingCtx, ticker.C, func(ctx context.Context) error {
		return database.Ping(ctx, pool)
	})
	return func() {
		cancel()
		ticker.Stop()
	}
}
