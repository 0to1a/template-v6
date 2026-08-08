package dbping

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type fakeExecer struct {
	calls int32
	err   error
}

func (f *fakeExecer) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	atomic.AddInt32(&f.calls, 1)
	return pgconn.CommandTag{}, f.err
}

func TestRun_NoPingBeforeFirstTick(t *testing.T) {
	exec := &fakeExecer{}
	ticks := make(chan time.Time)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go Run(ctx, exec, ticks)

	if got := atomic.LoadInt32(&exec.calls); got != 0 {
		t.Fatalf("calls before any tick = %d, want 0", got)
	}
}

func TestRun_TickInvokesExec(t *testing.T) {
	exec := &fakeExecer{}
	ticks := make(chan time.Time)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		Run(ctx, exec, ticks)
		close(done)
	}()

	ticks <- time.Now()
	ticks <- time.Now()
	ticks <- time.Now()

	cancel()
	<-done

	if got := atomic.LoadInt32(&exec.calls); got != 3 {
		t.Fatalf("calls = %d, want 3", got)
	}
}

func TestRun_ExecFailureIsLoggedAndLoopContinues(t *testing.T) {
	execErr := errors.New("connection reset")
	exec := &fakeExecer{err: execErr}
	ticks := make(chan time.Time)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var logs bytes.Buffer
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	done := make(chan struct{})
	go func() {
		Run(ctx, exec, ticks)
		close(done)
	}()

	ticks <- time.Now()
	ticks <- time.Now()

	cancel()
	<-done

	if got := atomic.LoadInt32(&exec.calls); got != 2 {
		t.Fatalf("calls = %d, want 2 (loop must continue after failure)", got)
	}
	if !strings.Contains(logs.String(), "select 1 failed") {
		t.Fatalf("log output missing failure message: %q", logs.String())
	}
}

func TestRun_CancelStopsPromptly(t *testing.T) {
	exec := &fakeExecer{}
	ticks := make(chan time.Time)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		Run(ctx, exec, ticks)
		close(done)
	}()

	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run did not stop after cancel")
	}
}
