.PHONY: build run test test-integration docker-build docker-up docker-down clean mocks

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
