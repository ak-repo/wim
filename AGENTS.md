# AGENTS.md - OpenCode Instructions for WIM

This file contains non-obvious repo-specific actions and architecture notes for WIM (Warehouse Inventory Management).

## Quick Start (Exact Commands)

```bash
# Start all infrastructure (postgres, redis, kafka)
make docker-up

# Database setup (MUST run after docker-up)
make migrate-up

# Start backend services (separate terminals)
make run-api      # API server
make run-worker   # Kafka worker

# Start frontend
cd frontend && npm install && npm run dev
```

## Critical Environment Configuration

**Backend requires specific env vars from `.env`:**
- `DATABASE_URL=postgres://wim_user:wim_pass@localhost:5432/warehouse_inventory?sslmode=disable`
- `REDIS_PORT=6380` (not default 6379)
- `REDIS_PASSWORD=wim_redis_pass`
- `KAFKA_BROKERS=localhost:9092`
- `KAFKA_TOPIC=warehouse-events`
- `AUTH_JWT_SECRET` - must be set for auth to work

**Frontend uses `.env` variables at build time**

**Configuration priority**: `.env` > `config/config.yaml` > `config/config.example.yaml`

## Architecture & Code Generation

- **Not a monorepo**: Backend (Go) and frontend (React/TypeScript) are separate with no shared codegen
- **No code generation**: Types must be kept manually in sync between backend and frontend
- **Backend router**: Uses `chi` (not Gin as README suggests) - see `internal/handler/*`
- **Frontend state**: TanStack Query for server state, Zustand for client state
- **Validation**: Zod schemas in frontend, manual validation in backend

## Development Workflow

```bash
# Run single test (example)
go test ./internal/service/product -run TestCreateProduct

# Check migration status
make migrate-status

# Create new migration
make migrate-create name=add_user_table

# View logs
make logs

# Database shell access
make db-shell
```

## Testing Quirks

- **No test database**: Tests use same DATABASE_URL as dev (be careful)
- **No seed data**: Database starts empty after migration
- **Frontend has no tests configured**

## Build & Deploy

```bash
# Build binaries
make build  # Creates bin/api and bin/worker

# Full reset (destroys all data)
make docker-reset
```
