# =============================================================================
# The Asheville Setlist - Makefile
# =============================================================================

.PHONY: help dev stop db migrate migrate-down seed test lint build clean

# Default target
help:
	@echo "The Asheville Setlist - Development Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Development:"
	@echo "  dev          Start local development environment (DB + services)"
	@echo "  stop         Stop all containers"
	@echo "  logs         Show container logs"
	@echo ""
	@echo "Database:"
	@echo "  db           Start PostgreSQL container only"
	@echo "  db-shell     Open psql shell"
	@echo "  migrate      Run database migrations"
	@echo "  migrate-down Rollback last migration"
	@echo "  migrate-new  Create new migration (usage: make migrate-new name=add_users)"
	@echo "  seed         Load seed data"
	@echo "  db-reset     Drop and recreate database with seeds"
	@echo ""
	@echo "Backend:"
	@echo "  api          Run Go API server"
	@echo "  scraper      Run scraper once"
	@echo "  test         Run all tests"
	@echo "  lint         Run linters"
	@echo "  build        Build binaries"
	@echo ""
	@echo "Frontend:"
	@echo "  web          Run Next.js dev server"
	@echo "  web-build    Build frontend for production"
	@echo ""
	@echo "Utilities:"
	@echo "  clean        Clean build artifacts"
	@echo "  install      Install dependencies"

# =============================================================================
# ENVIRONMENT
# =============================================================================

# Load .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Default database URL for local development
DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/asheville_setlist?sslmode=disable

# =============================================================================
# DEVELOPMENT
# =============================================================================

# Start everything
dev: db
	@echo "✓ Database is running on localhost:5432"
	@echo ""
	@echo "Next steps:"
	@echo "  make migrate  - Run database migrations"
	@echo "  make seed     - Load seed data"
	@echo "  make api      - Start Go API server"
	@echo "  make web      - Start Next.js frontend"

# Stop all containers
stop:
	docker compose down

# Show logs
logs:
	docker compose logs -f

# =============================================================================
# DATABASE
# =============================================================================

# Start PostgreSQL only
db:
	docker compose up -d db
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 2
	@until docker compose exec -T db pg_isready -U postgres > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "✓ PostgreSQL is ready"

# Open psql shell
db-shell:
	docker compose exec db psql -U postgres -d asheville_setlist

# Run migrations
migrate:
	@echo "Running migrations..."
	migrate -path backend/migrations -database "$(DATABASE_URL)" up
	@echo "✓ Migrations complete"

# Rollback last migration
migrate-down:
	@echo "Rolling back last migration..."
	migrate -path backend/migrations -database "$(DATABASE_URL)" down 1

# Create new migration
migrate-new:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-new name=migration_name"; \
		exit 1; \
	fi
	migrate create -ext sql -dir backend/migrations -seq $(name)

# Load seed data
seed:
	@echo "Loading seed data..."
	@for f in seeds/*.sql; do \
		echo "  Running $$f..."; \
		docker compose exec -T db psql -U postgres -d asheville_setlist < $$f; \
	done
	@echo "✓ Seed data loaded"

# Reset database (drop, recreate, migrate, seed)
db-reset: db
	@echo "⚠️  Resetting database..."
	docker compose exec -T db psql -U postgres -c "DROP DATABASE IF EXISTS asheville_setlist"
	docker compose exec -T db psql -U postgres -c "CREATE DATABASE asheville_setlist"
	@$(MAKE) migrate
	@$(MAKE) seed
	@echo "✓ Database reset complete"

# =============================================================================
# BACKEND
# =============================================================================

# Run API server
api:
	cd backend && go run ./cmd/api

# Run scraper
scraper:
	cd backend && go run ./cmd/scraper

# Run tests
test:
	cd backend && go test -v ./...

# Run linter
lint:
	cd backend && golangci-lint run

# Build binaries
build:
	cd backend && go build -o bin/api ./cmd/api
	cd backend && go build -o bin/scraper ./cmd/scraper

# =============================================================================
# FRONTEND
# =============================================================================

# Run Next.js dev server
web:
	cd frontend && npm run dev

# Build frontend
web-build:
	cd frontend && npm run build

# =============================================================================
# UTILITIES
# =============================================================================

# Install dependencies
install:
	cd backend && go mod download
	cd frontend && npm install

# Clean build artifacts
clean:
	rm -rf backend/bin
	rm -rf frontend/.next
	rm -rf frontend/node_modules/.cache

# =============================================================================
# DOCKER COMMANDS
# =============================================================================

# Build Docker images
docker-build:
	docker build -t asheville-api ./backend/cmd/api
	docker build -t asheville-scraper ./backend/cmd/scraper

# Start with pgAdmin
dev-full:
	docker compose --profile tools up -d
