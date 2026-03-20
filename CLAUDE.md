# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build & Run
```bash
# Build the API server
go build -o bin/api cmd/api/main.go

# Build the worker
go build -o bin/worker cmd/worker/main.go

# Run the API server
./bin/api

# Run the worker
./bin/worker
```

### Docker
```bash
# Build Docker images
docker-compose build

# Start all services
docker-compose up -d

# Stop all services
docker-compose down
```

### Testing
```bash
# Run unit tests
go test ./... -v

# Run integration tests
docker-compose run --rm api go test ./... -v
```

### Linting
```bash
# Run Go lint
golangci-lint run

# Run Dockerfile lint
dockerfile-lint .
```

## Architecture Overview

The system follows a microservices architecture with the following components:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           WIM System                                    │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌──────────────┐      ┌──────────────┐      ┌──────────────┐         │
│  │   Clients    │──────▶│   REST API   │──────▶│  Services   │         │
│  │  (Frontend)  │      │    (Gin)     │      │  (Business) │         │
│  └──────────────┘      └──────────────┘      └──────┬───────┘         │
│                                                      │                  │
│                       ┌───────────────────────────────┼────────────┐   │
│                       │                               │            │   │
│                       ▼                               ▼            ▼   │
│              ┌────────────────┐           ┌────────────────┐          │
│              │  PostgreSQL     │           │     Redis      │          │
│              │  (Persistence)  │           │    (Cache)     │          │
│              └────────────────┘           └────────────────┘          │
│                                                      │                  │
│                                                      ▼                  │
│              ┌────────────────┐      ┌─────────────────────┐          │
│              │    Kafka       │─────▶│     Worker          │          │
│              │    Queue       │      │  (Async Processing)│          │
│              └────────────────┘      └─────────────────────┘          │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Configuration

Create a `config.yaml` file with the following structure:

```yaml
server:
  port: 8080
  host: 0.0.0.0

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  database: warehouse_inventory
  ssl_mode: disable
  max_conns: 25

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

kafka:
  brokers:
    - localhost:9092
  topic: warehouse-events
  group_id: warehouse-worker

worker:
  pool_size: 5
  queue_size: 100
  retry_count: 3
  retry_delay: 1s
  batch_size: 10

log_level: info
```

## Hooks

Include important Git hooks from `.git/hooks/`:

- `pre-commit.sample`: Run linters and tests before commits
- `pre-push.sample`: Validate changes before pushing
- `prepare-commit-msg.sample`: Auto-generate commit messages

## References

- README.md: Detailed documentation and API examples
- go.mod: Dependency management
- docker-compose.yml: Container orchestration
- config.yaml: Application configuration

## Notes

1. The system uses Gin for the REST API and Kafka for async processing
2. PostgreSQL handles persistent storage while Redis caches frequently accessed data
3. Worker processes handle background tasks like inventory recalculations
4. All services are containerized for consistent environments