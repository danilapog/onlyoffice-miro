package pg

import (
	"context"
	"time"

	pgx "github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

func executeTransactionally(
	ctx context.Context,
	pool *pgxpool.Pool,
	opts pgx.TxOptions,
	statementName, statementQuery string,
	fn func(tx pgx.Tx) error,
) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	tx, err := conn.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	if _, err = tx.Prepare(ctx, statementName, statementQuery); err != nil {
		tctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		return tx.Rollback(tctx)
	}

	if err = fn(tx); err != nil {
		tctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		return tx.Rollback(tctx)
	}

	return tx.Commit(ctx)
}
