package database

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
)

type fakeQuerier struct {
	err error
}

type fakeRow struct{ err error }

func (r fakeRow) Scan(dest ...any) error { return r.err }

func (q fakeQuerier) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return fakeRow{err: q.err}
}

func TestPing_OK(t *testing.T) {
	if err := Ping(context.Background(), fakeQuerier{}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPing_PropagatesQueryError(t *testing.T) {
	want := errors.New("connection refused")
	err := Ping(context.Background(), fakeQuerier{err: want})
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}
