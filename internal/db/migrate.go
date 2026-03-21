package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/pressly/goose/v3"
)

const migrationsDir = "migrations"

func RunMigrations(db *sql.DB) error {
	log := slog.Default().With("component", "db-migrations")

	if err := goose.SetDialect("postgres"); err != nil {
		log.Error("failed to configure migration dialect", "error", err)
		return fmt.Errorf("set goose dialect: %w", err)
	}

	fromVersion, err := goose.GetDBVersion(db)
	if err != nil {
		log.Error("failed to fetch current migration version", "error", err)
		return fmt.Errorf("get current migration version: %w", err)
	}

	log.Info("running database migrations", "from_version", fromVersion, "dir", migrationsDir)

	if err := goose.Up(db, migrationsDir); err != nil {
		log.Error("database migration failed", "error", err)
		return fmt.Errorf("run migrations up: %w", err)
	}

	toVersion, err := goose.GetDBVersion(db)
	if err != nil {
		log.Error("failed to fetch migration version after run", "error", err)
		return fmt.Errorf("get migration version after run: %w", err)
	}

	log.Info("database migrations complete", "to_version", toVersion)

	return nil
}
