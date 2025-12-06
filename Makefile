.PHONY: help dev dev-services dev-backend dev-frontend dev-stop \
        migrate migrate-new migrate-down \
        test test-backend test-frontend \
        build build-backend build-frontend clean \
        install fmt lint

# Configuration
DEV_COMPOSE := docker compose -f docker-compose.dev.yml
DATABASE_URL ?= postgres://trackable:secret@localhost:5432/trackable?sslmode=disable

# Default target
help:
	@echo "Trackable Development Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev              - Start services + backend + frontend"
	@echo "  make dev-services     - Start dev services (db, mailcatcher)"
	@echo "  make dev-backend      - Start backend with hot-reload"
	@echo "  make dev-frontend     - Start frontend with HMR"
	@echo "  make dev-stop         - Stop dev services"
	@echo ""
	@echo "Database:"
	@echo "  make migrate          - Run database migrations"
	@echo "  make migrate-new      - Create new migration (name=xxx)"
	@echo "  make migrate-down     - Rollback last migration"
	@echo ""
	@echo "Testing:"
	@echo "  make test             - Run all tests"
	@echo "  make test-backend     - Run backend tests"
	@echo "  make test-frontend    - Run frontend tests"
	@echo ""
	@echo "Build:"
	@echo "  make build            - Build production images"
	@echo "  make build-backend    - Build backend binary"
	@echo "  make build-frontend   - Build frontend assets"
	@echo "  make clean            - Clean build artifacts"
	@echo ""
	@echo "Other:"
	@echo "  make install          - Install dependencies"
	@echo "  make fmt              - Format code"
	@echo "  make lint             - Lint code"

# =============================================================================
# Development
# =============================================================================

# Start everything for development
dev: dev-services
	@echo "Waiting for database..."
	@sleep 2
	@$(MAKE) -j2 dev-backend dev-frontend

# Start dev services (database + mailcatcher)
dev-services:
	$(DEV_COMPOSE) up -d

# Start backend with hot-reload
dev-backend:
	cd backend && DATABASE_URL="$(DATABASE_URL)" air

# Start frontend with HMR
dev-frontend:
	cd frontend && npm run dev

# Stop dev services
dev-stop:
	$(DEV_COMPOSE) down

# =============================================================================
# Database
# =============================================================================

# Run migrations
migrate:
	cd backend && go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
		-path ./migrations \
		-database "$(DATABASE_URL)" \
		up

# Create new migration
migrate-new:
ifndef name
	$(error Usage: make migrate-new name=migration_name)
endif
	cd backend && go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
		create -ext sql -dir ./migrations -seq $(name)

# Rollback last migration
migrate-down:
	cd backend && go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
		-path ./migrations \
		-database "$(DATABASE_URL)" \
		down 1

# =============================================================================
# Testing
# =============================================================================

test: test-backend test-frontend

test-backend:
	cd backend && go test -v ./...

test-frontend:
	cd frontend && npm test

# =============================================================================
# Build
# =============================================================================

# Build production images
build:
	docker compose build

build-backend:
	cd backend && go build -ldflags="-w -s" -o bin/server ./cmd/server

build-frontend:
	cd frontend && npm run build

clean:
	rm -rf backend/bin backend/tmp
	rm -rf frontend/dist frontend/node_modules/.vite
	$(DEV_COMPOSE) down -v --remove-orphans

# =============================================================================
# Other
# =============================================================================

install:
	cd backend && go mod download
	cd frontend && npm install

fmt:
	cd backend && go fmt ./...
	cd frontend && npm run format 2>/dev/null || true

lint:
	cd backend && go vet ./...
	cd frontend && npm run lint 2>/dev/null || true
