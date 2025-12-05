.PHONY: help dev dev-backend dev-frontend db migrate migrate-new test build clean

# Default target
help:
	@echo "Trackable Development Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev            - Start all services (db, backend, frontend)"
	@echo "  make dev-backend    - Start only backend with hot-reload"
	@echo "  make dev-frontend   - Start only frontend with HMR"
	@echo "  make db             - Start only database"
	@echo ""
	@echo "Database:"
	@echo "  make migrate        - Run database migrations"
	@echo "  make migrate-new    - Create new migration (name=migration_name)"
	@echo "  make migrate-down   - Rollback last migration"
	@echo ""
	@echo "Testing:"
	@echo "  make test           - Run all tests"
	@echo "  make test-backend   - Run backend tests"
	@echo "  make test-frontend  - Run frontend tests"
	@echo ""
	@echo "Build:"
	@echo "  make build          - Build production images"
	@echo "  make clean          - Clean build artifacts"

# Start all services
dev:
	docker-compose up db -d
	@echo "Waiting for database to be ready..."
	@sleep 3
	@make -j2 dev-backend dev-frontend

# Start backend with hot-reload
dev-backend:
	cd backend && air -c .air.toml

# Start frontend with HMR
dev-frontend:
	cd frontend && npm run dev

# Start only database
db:
	docker-compose up db -d

# Run migrations
migrate:
	cd backend && go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
		-path ./migrations \
		-database "$${DATABASE_URL:-postgres://trackable:secret@localhost:5432/trackable?sslmode=disable}" \
		up

# Create new migration
migrate-new:
ifndef name
	$(error name is required. Usage: make migrate-new name=migration_name)
endif
	cd backend && go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
		create -ext sql -dir ./migrations -seq $(name)

# Rollback last migration
migrate-down:
	cd backend && go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
		-path ./migrations \
		-database "$${DATABASE_URL:-postgres://trackable:secret@localhost:5432/trackable?sslmode=disable}" \
		down 1

# Run all tests
test: test-backend test-frontend

# Run backend tests
test-backend:
	cd backend && go test -v ./...

# Run frontend tests
test-frontend:
	cd frontend && npm test

# Build production images
build:
	docker-compose build

# Build backend only
build-backend:
	cd backend && go build -ldflags="-w -s" -o bin/server ./cmd/server

# Build frontend only
build-frontend:
	cd frontend && npm run build

# Clean build artifacts
clean:
	rm -rf backend/bin backend/tmp
	rm -rf frontend/dist frontend/node_modules/.vite
	docker-compose down -v --remove-orphans

# Install dependencies
install:
	cd backend && go mod download
	cd frontend && npm install

# Format code
fmt:
	cd backend && go fmt ./...
	cd frontend && npm run format

# Lint code
lint:
	cd backend && go vet ./...
	cd frontend && npm run lint
