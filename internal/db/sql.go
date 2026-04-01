package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ak-repo/wim/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func OpenSQLConnection(ctx context.Context, cfg config.DatabaseConfig) (*sql.DB, error) {
	sqlDB, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("open sql connection: %w", err)
	}

	sqlDB.SetMaxOpenConns(int(cfg.MaxConns))

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

	for attempt := 1; attempt <= retries; attempt++ {
		err = sqlDB.PingContext(ctx)
		if err == nil {
			return sqlDB, nil
		}

		if attempt == retries {
			sqlDB.Close()
			return nil, fmt.Errorf("ping sql database after %d attempts: %w", retries, err)
		}

		select {
		case <-ctx.Done():
			sqlDB.Close()
			return nil, fmt.Errorf("ping sql database canceled: %w", ctx.Err())
		case <-time.After(delay):
		}

		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
		}
	}

	sqlDB.Close()
	return nil, fmt.Errorf("ping sql database failed")
}
