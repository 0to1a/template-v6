package background

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TC: each tick invokes the job, and a failing tick doesn't stop the next one
func TestRun_TicksInvokeJobAndSurviveErrors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ticks := make(chan time.Time)

	var calls int
	done := make(chan struct{})
	go func() {
		Run(ctx, ticks, func(context.Context) error {
			calls++
			if calls == 1 {
				return errors.New("boom")
			}
			return nil
		})
		close(done)
	}()

	ticks <- time.Now() // fails, must not stop the loop
	ticks <- time.Now() // succeeds

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run did not return after cancel")
	}

	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

// TC: cancelling ctx stops the loop even with no pending tick
func TestRun_StopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ticks := make(chan time.Time)

	done := make(chan struct{})
	go func() {
		Run(ctx, ticks, func(context.Context) error { return nil })
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run did not return after cancel")
	}
}
