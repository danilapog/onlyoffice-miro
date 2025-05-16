package pg

import (
	"context"
	"runtime"
	"time"

	pgx "github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(connection string, opts ...Option) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connection)
	if err != nil {
		return nil, err
	}

	maxConns := runtime.NumCPU() * 4
	if maxConns < 5 {
		maxConns = 5
	}

	options := Options{
		MaxConnections:              maxConns,
		MinConnections:              2,
		ConnectionTimeout:           5 * time.Second,
		ConnectionIdleTimeout:       5 * time.Minute,
		ConnectionHealthcheckPeriod: 15 * time.Second,
		MaxConnLifetime:             1 * time.Hour,
		MaxRetries:                  3,
		RetryInterval:               500 * time.Millisecond,
	}

	for _, opt := range opts {
		opt(&options)
	}

	config.MaxConns = int32(options.MaxConnections)
	config.MinConns = int32(options.MinConnections)
	config.MaxConnIdleTime = options.ConnectionIdleTimeout
	config.HealthCheckPeriod = options.ConnectionHealthcheckPeriod
	config.MaxConnLifetime = options.MaxConnLifetime
	config.ConnConfig.ConnectTimeout = options.ConnectionTimeout
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, "SET statement_timeout = '5s'")
		return err
	}

	config.ConnConfig.RuntimeParams = map[string]string{
		"application_name": "onlyoffice-miro",
	}

	var pool *pgxpool.Pool
	var retryErr error

	for i := 0; i <= options.MaxRetries; i++ {
		if i > 0 {
			time.Sleep(options.RetryInterval)
		}

		pool, retryErr = pgxpool.NewWithConfig(context.Background(), config)
		if retryErr == nil {
			ctx, cancel := context.WithTimeout(context.Background(), options.ConnectionTimeout)
			err = pool.Ping(ctx)
			cancel()

			if err == nil {
				return pool, nil
			}

			pool.Close()
			retryErr = err
		}
	}

	if retryErr != nil {
		return nil, retryErr
	}

	return pool, nil
}
