package config

import "testing"

func TestDatabaseConfigDSN_UsesURLWhenProvided(t *testing.T) {
	cfg := DatabaseConfig{URL: "postgres://u:p@localhost:5432/db?sslmode=disable"}

	if got := cfg.DSN(); got != cfg.URL {
		t.Fatalf("expected URL dsn, got %q", got)
	}
}

func TestDatabaseConfigDSN_BuildsFromFields(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "wim_user",
		Password: "wim_pass",
		Database: "warehouse_inventory",
		SSLMode:  "disable",
	}

	want := "postgres://wim_user:wim_pass@localhost:5432/warehouse_inventory?sslmode=disable"
	if got := cfg.DSN(); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
