DB_URL?=postgres://postgres:postgres@localhost:5435/phd?sslmode=disable
AIR?=air
# VERSION is used by `make migrate-force VERSION=<n>` to set the forced migration version.
VERSION?=1

.PHONY: help run migrate-up migrate-down migrate-force seed dev

help:
	@echo "Available targets:"
	@echo "  make run               # start the API server"
	@echo "  make dev               # run with Air live reload (AIR env configurable)"
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
	$(AIR)
