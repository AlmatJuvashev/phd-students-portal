DB_URL?=postgres://postgres:postgres@localhost:5432/phd?sslmode=disable
AIR?=air

.PHONY: run migrate-up migrate-down seed dev

run:
	go run ./cmd/server

migrate-up:
	migrate -database $(DB_URL) -path db/migrations up

migrate-down:
	migrate -database $(DB_URL) -path db/migrations down

seed:
	go run ./cmd/server --seed

dev:
	$(AIR)
