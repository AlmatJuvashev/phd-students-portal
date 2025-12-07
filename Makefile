.PHONY: help up down dev migrate-up migrate-down kill-ports install logs

# Configuration
BACKEND_DIR := backend
FRONTEND_DIR := frontend
MAILSERVER_COMPOSE := mailserver/docker-compose.yml
ROOT_COMPOSE := docker-compose.yml

# Ports to clean up
PORTS := 8080 5173 5435 9090 9091 8025 1025

help:
	@echo "KazNMU PhD Portal - Root Makefile"
	@echo "---------------------------------"
	@echo "Usage:"
	@echo "  make up            Start infrastructure (Postgres, Minio, Mailpit) in background"
	@echo "  make down          Stop all infrastructure containers"
	@echo "  make dev           Start infrastructure + Backend + Frontend (interactive)"
	@echo "  make install       Install dependencies (Go mod + NPM)"
	@echo "  make migrate-up    Run database migrations"
	@echo "  make migrate-down  Rollback last migration"
	@echo "  make kill-ports    Kill processes on ports $(PORTS)"
	@echo "  make logs          Follow infrastructure logs"

up:
	@echo "Starting infrastructure..."
	docker compose -f $(ROOT_COMPOSE) -f $(MAILSERVER_COMPOSE) up -d

down:
	@echo "Stopping infrastructure..."
	docker compose -f $(ROOT_COMPOSE) -f $(MAILSERVER_COMPOSE) down

dev: up
	@echo "Starting applications..."
	@trap 'kill 0' SIGINT; \
	(cd $(BACKEND_DIR) && make run) & \
	(cd $(FRONTEND_DIR) && npm run dev) & \
	wait

install:
	@echo "Installing Backend dependencies..."
	cd $(BACKEND_DIR) && go mod download
	@echo "Installing Frontend dependencies..."
	cd $(FRONTEND_DIR) && npm install

migrate-up:
	@echo "Running migrations..."
	cd $(BACKEND_DIR) && make migrate-up

migrate-down:
	@echo "Rolling back migration..."
	cd $(BACKEND_DIR) && make migrate-down

kill-ports:
	@echo "Killing processes on ports: $(PORTS)"
	@lsof -ti:$(shell echo $(PORTS) | tr ' ' ',') | xargs kill -9 2>/dev/null || true
	@echo "Ports cleared."

logs:
	docker compose -f $(ROOT_COMPOSE) -f $(MAILSERVER_COMPOSE) logs -f
