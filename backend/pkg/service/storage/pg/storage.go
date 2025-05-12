package pg

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"time"

	"github.com/ONLYOFFICE/onlyoffice-miro/backend/internal/pkg/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func generateStatementName(query string) string {
	h := md5.Sum([]byte(query))
	return fmt.Sprintf("stmt_%x", h)
}

type postgresStorage[ID comparable, T any] struct {
	pool      *pgxpool.Pool
	processor service.StorageProcessor[ID, T, pgx.Row]
}

func NewPostgresStorage[ID comparable, T any](
	pool *pgxpool.Pool,
	processor service.StorageProcessor[ID, T, pgx.Row],
) (service.Storage[ID, T], error) {
	if pool == nil {
		return nil, ErrNilPool
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &postgresStorage[ID, T]{
		pool:      pool,
		processor: processor,
	}, nil
}

func (s *postgresStorage[ID, T]) Find(ctx context.Context, id ID) (T, error) {
	var result T

	query, args, scanner := s.processor.BuildSelectQuery(id)
	opts := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	}

	sName := generateStatementName(query)
	err := executeTransactionally(ctx, s.pool, opts, sName, query, func(tx pgx.Tx) error {
		row := tx.QueryRow(ctx, sName, args...)
		var err error
		result, err = scanner(row)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrNoRowsAffected
			}

			return err
		}

		return nil
	})

	return result, err
}

func (s *postgresStorage[ID, T]) Insert(ctx context.Context, id ID, value T) (T, error) {
	query, args := s.processor.BuildInsertQuery(id, value)
	opts := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}

	sName := generateStatementName(query)
	return value, executeTransactionally(ctx, s.pool, opts, sName, query, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, sName, args...)
		return err
	})
}

func (s *postgresStorage[ID, T]) Update(ctx context.Context, id ID, value T) (T, error) {
	query, args := s.processor.BuildUpdateQuery(id, value)
	opts := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}

	sName := generateStatementName(query)
	return value, executeTransactionally(ctx, s.pool, opts, sName, query, func(tx pgx.Tx) error {
		res, err := tx.Exec(ctx, sName, args...)
		if err != nil {
			return err
		}

		if res.RowsAffected() == 0 {
			return ErrNoRowsAffected
		}

		return nil
	})
}

func (s *postgresStorage[ID, T]) Delete(ctx context.Context, id ID) error {
	query, args := s.processor.BuildDeleteQuery(id)
	opts := pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	}

	sName := generateStatementName(query)
	return executeTransactionally(ctx, s.pool, opts, sName, query, func(tx pgx.Tx) error {
		res, err := tx.Exec(ctx, sName, args...)
		if err != nil {
			return err
		}

		if res.RowsAffected() == 0 {
			return ErrNoRowsAffected
		}

		return nil
	})
}
