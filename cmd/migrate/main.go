package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ak-repo/wim/config"
	dbutil "github.com/ak-repo/wim/internal/db"
	"github.com/pressly/goose/v3"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	db, err := dbutil.OpenSQLConnection(ctx, cfg.Database)
	if err != nil {
		panic(fmt.Errorf("open db: %w", err))
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		panic(fmt.Errorf("set dialect: %w", err))
	}

	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "up":
		err = goose.Up(db, "migrations")
	case "down":
		err = goose.Down(db, "migrations")
	case "status":
		err = goose.Status(db, "migrations")
	default:
		panic(fmt.Errorf("unknown command %q: use up|down|status", cmd))
	}

	if err != nil {
		panic(fmt.Errorf("goose %s: %w", cmd, err))
	}
}
