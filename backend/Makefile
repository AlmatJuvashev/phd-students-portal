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

.PHONY: help run migrate-up migrate-down migrate-force seed dev kill-port

help:
	@echo "Available targets:"
	@echo "  make run               # start the API server"
	@echo "  make dev               # run with Air live reload (AIR env configurable)"
	@echo "  make kill-port         # kill processes bound to port $(PORT_TO_KILL)"
	@echo "  make migrate-up        # apply all pending migrations"
	@echo "  make migrate-down      # roll back the most recent migration"
	@echo "  make migrate-force VERSION=N  # force-set schema version to N (use after fixing dirty state)"
	@echo "  make seed              # seed demo data"
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
