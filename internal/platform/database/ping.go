package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Narrow enough that *pgxpool.Pool satisfies it without wrapping
type Querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Ping proves the connection is alive; read-only, no schema required
func Ping(ctx context.Context, q Querier) error {
	var one int
	return q.QueryRow(ctx, "SELECT 1").Scan(&one)
}
