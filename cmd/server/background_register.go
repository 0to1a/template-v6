package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"project/internal/platform/background"
	"project/internal/platform/database"
)

// All background services live in this one file, one function per job.
// Add a new job by writing its func(context.Context) error and appending
// another registerBackgroundJob(ctx, job, interval) call below.
func registerBackground(ctx context.Context, pool *pgxpool.Pool) (stop func()) {
	stops := []func(){
		registerBackgroundJob(ctx, pingDatabase(pool), time.Minute),
	}
	return func() {
		for _, stop := range stops {
			stop()
		}
	}
}

// registerBackgroundJob runs job every interval until stopped. A failure is
// logged by background.Run and never stops the loop, so the next tick still runs.
func registerBackgroundJob(ctx context.Context, job func(context.Context) error, interval time.Duration) (stop func()) {
	jobCtx, cancel := context.WithCancel(ctx)
	ticker := time.NewTicker(interval)
	go background.Run(jobCtx, ticker.C, job)
	return func() {
		cancel()
		ticker.Stop()
	}
}

// pingDatabase checks connectivity with SELECT 1.
func pingDatabase(pool *pgxpool.Pool) func(context.Context) error {
	return func(ctx context.Context) error {
		return database.Ping(ctx, pool)
	}
}
