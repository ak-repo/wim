SHELL := /bin/bash

GOOSE_VERSION ?= v3.24.1
GOOSE_CMD := go run github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
MIGRATIONS_DIR := migrations
DATABASE_URL ?= postgres://wim_user:wim_pass@localhost:5432/warehouse_inventory?sslmode=disable

.PHONY: run-api run-worker build docker-up docker-down docker-reset migrate-up migrate-down migrate-status migrate-create logs ps db-shell

run-api:
	go run ./cmd/api

run-worker:
	go run ./cmd/worker

build:
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

docker-reset:
	docker compose down -v --remove-orphans
	docker compose up --build -d

migrate-up:
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

migrate-down:
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

migrate-status:
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" status

migrate-create:
	@if [ -z "$(name)" ]; then echo "usage: make migrate-create name=<migration_name>"; exit 1; fi
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) create $(name) sql

logs:
	docker compose logs -f --tail=200

ps:
	docker compose ps

db-shell:
	docker compose exec postgres psql -U wim_user -d warehouse_inventory
