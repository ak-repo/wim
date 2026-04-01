package db

import (
	"context"
	"fmt"
	"time"

	"github.com/ak-repo/wim/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewConnection(ctx context.Context, cfg config.DatabaseConfig) (*DB, error) {
	connStr := cfg.DSN()

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	retries := cfg.ConnectRetries
	if retries < 1 {
		retries = 1
	}

	delay := cfg.ConnectRetryInitialDelay
	if delay <= 0 {
		delay = time.Second
	}

	maxDelay := cfg.ConnectRetryMaxDelay
	if maxDelay <= 0 {
		maxDelay = 30 * time.Second
	}

	var pingErr error
	for attempt := 1; attempt <= retries; attempt++ {
		pingErr = pool.Ping(ctx)
		if pingErr == nil {
			break
		}

		if attempt == retries {
			pool.Close()
			return nil, fmt.Errorf("failed to ping database after %d attempts: %w", retries, pingErr)
		}

		time.Sleep(delay)
		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
		}
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
