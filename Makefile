.PHONY: build run test test-integration docker-build docker-up docker-down clean mocks \
       install-migrate migrate-create migrate-up migrate-down migrate-version migrate-force

# Generate mocks
mocks:
	go generate ./...

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run the application locally
run:
	go run ./cmd/server

# Run unit tests
test:
	go test -v ./...

# Run integration tests (requires running database)
test-integration:
	go test -v -tags=integration ./...

# Build Docker image
docker-build:
	docker-compose build

# Start all services
docker-up:
	docker-compose up -d

# Stop all services
docker-down:
	docker-compose down

# Stop and remove volumes
docker-clean:
	docker-compose down -v

# View logs
logs:
	docker-compose logs -f

# Clean build artifacts
clean:
	rm -rf bin/

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Download dependencies
deps:
	go mod download
	go mod tidy

# golang-migrate
# Set DATABASE_URL env var: export DATABASE_URL="postgres://user:pass@localhost:5432/db?sslmode=disable"
MIGRATIONS_DIR=internal/repository/migrations

# Install golang-migrate CLI
install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create new migration (usage: make migrate-create name=create_users)
migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=migration_name"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

# Run all up migrations
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$$DATABASE_URL" up

# Rollback one migration
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$$DATABASE_URL" down 1

# Show current version
migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$$DATABASE_URL" version

# Force set version (usage: make migrate-force version=1)
migrate-force:
	@if [ -z "$(version)" ]; then echo "Usage: make migrate-force version=N"; exit 1; fi
	migrate -path $(MIGRATIONS_DIR) -database "$$DATABASE_URL" force $(version)
