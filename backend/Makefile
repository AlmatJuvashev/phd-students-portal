DB_URL?=postgres://postgres:postgres@localhost:5435/phd?sslmode=disable
AIR?=air

# Automatically determine APP_PORT for `make dev` when not provided.
ifeq ($(filter dev,$(MAKECMDGOALS)),dev)
ifndef APP_PORT
ifdef PORT
APP_PORT := $(PORT)
else
APP_PORT := $(strip $(shell python3 -c 'import socket; s=socket.socket(); s.bind(("", 0)); print(s.getsockname()[1]); s.close()'))
endif
endif
endif
# VERSION is used by `make migrate-force VERSION=<n>` to set the forced migration version.
VERSION?=1
PORT_TO_KILL?=8280

.PHONY: help run migrate-up migrate-down migrate-force seed seed-demo dev kill-port

help:
	@echo "Available targets:"
	@echo "  make run               # start the API server"
	@echo "  make dev               # run with Air live reload (AIR env configurable)"
	@echo "  make kill-port         # kill processes bound to port $(PORT_TO_KILL)"
	@echo "  make migrate-up        # apply all pending migrations"
	@echo "  make migrate-down      # roll back the most recent migration"
	@echo "  make migrate-force VERSION=N  # force-set schema version to N (use after fixing dirty state)"
	@echo "  make seed              # seed initial data"
	@echo "  make seed-demo         # seed comprehensive demo data for demo.university"
	@echo "Environment variables:"
	@echo "  DB_URL=$(DB_URL)"
	@echo "  VERSION=$(VERSION)"
	@echo "  AIR=$(AIR)"

run:
	go run ./cmd/server

migrate-up:
	migrate -database '$(DB_URL)' -path db/migrations up

migrate-down:
	migrate -database '$(DB_URL)' -path db/migrations down

migrate-force:
	migrate -database '$(DB_URL)' -path db/migrations force $(VERSION)

seed:
	go run ./cmd/server --seed

seed-demo:
	@echo "Seeding comprehensive demo data for demo.university tenant..."
	go run ./cmd/mock

dev:
	@echo "Starting dev server on APP_PORT=$(APP_PORT)"
	@if command -v $(AIR) >/dev/null 2>&1; then \
		APP_PORT=$(APP_PORT) $(AIR); \
	else \
		echo "Warning: '$(AIR)' command not found. Falling back to 'go run ./cmd/server'"; \
		APP_PORT=$(APP_PORT) go run ./cmd/server; \
	fi

kill-port:
	@echo "Killing processes listening on port $(PORT_TO_KILL)"
	@lsof -ti:$(PORT_TO_KILL) | xargs kill -9 2>/dev/null || true

# ============= Test Targets =============
TEST_DB_URL?=postgres://postgres:postgres@localhost:5435/phd_test?sslmode=disable
GOLANGCI_LINT_VERSION=v1.55.2

create-test-db:
	@echo "Creating phd_test database if it doesn't exist..."
	@psql "$(DB_URL)" -tc "SELECT 1 FROM pg_database WHERE datname = 'phd_test'" | grep -q 1 || \
		psql "$(DB_URL)" -c "CREATE DATABASE phd_test"
	@echo "Test database ready."

test: create-test-db
	@echo "Running all tests..."
	TEST_DATABASE_URL=$(TEST_DB_URL) go test ./... -count=1

test-unit:
	@echo "Running UNIT tests only (fast, no DB)..."
	go test -v -run Unit ./... -count=1

test-integration: create-test-db
	@echo "Running INTEGRATION tests only..."
	TEST_DATABASE_URL=$(TEST_DB_URL) go test -v -skip Unit ./... -count=1

test-coverage: create-test-db
	@echo "Running tests with coverage analysis..."
	TEST_DATABASE_URL=$(TEST_DB_URL) go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	@echo "Generating HTML report..."
	go tool cover -html=coverage.out -o coverage.html

lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Installing $(GOLANGCI_LINT_VERSION)..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION); \
		$$(go env GOPATH)/bin/golangci-lint run; \
	fi

.PHONY: lint test test-unit test-integration test-coverage
test-info:
	@echo "=== High-Performance Testing Infrastructure ==="
	@echo "The backend now supports true package-level parallelism using database cloning."
	@echo ""
	@echo "How it works:"
	@echo "1. One time setup: A 'phd_test_base' database is created and migrated."
	@echo "2. Per-package: Each Go package (handlers, services, etc.) gets its own clone (e.g., 'phd_test_handlers')."
	@echo "3. Isolation: This allows running tests with '-p 8' or more without data interference."
	@echo ""
	@echo "Commands:"
	@echo "  make test           - Standard sequential run (slowest, safest)"
	@echo "  make test-parallel  - High-speed run with 8 parallel packages (recommended)"
	@echo "  make test-fast      - Ultra-high speed with 10 parallel packages"
	@echo "  make test-coverage  - Parallel run with coverage output"
	@echo "  make test-info      - Show this helper"
	@echo "==============================================="


