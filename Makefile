.PHONY: help build run test clean docker-build docker-run frontend-install frontend-start frontend-build all

# Default target
help:
	@echo "Available commands:"
	@echo "  build           - Build the Go binary"
	@echo "  run            - Run the Go service"
	@echo "  test           - Run Go tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker Compose"
	@echo "  frontend-install - Install frontend dependencies"
	@echo "  frontend-start - Start frontend development server"
	@echo "  frontend-build - Build frontend for production"
	@echo "  all            - Build and run everything"

# Build the Go binary
build:
	@echo "Building Go binary..."
	go build -o bin/crypto-insights cmd/main.go

# Run the Go service
run: build
	@echo "Running Go service..."
	./bin/crypto-insights

# Run Go tests
test:
	@echo "Running Go tests..."
	go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -f Dockerfile.go -t crypto-insights:latest .

# Run with Docker Compose
docker-run:
	@echo "Starting services with Docker Compose..."
	docker-compose -f docker-compose.go.yml up -d

# Stop Docker services
docker-stop:
	@echo "Stopping Docker services..."
	docker-compose -f docker-compose.go.yml down

# Install frontend dependencies
frontend-install:
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

# Start frontend development server
frontend-start: frontend-install
	@echo "Starting frontend development server..."
	cd frontend && npm start

# Build frontend for production
frontend-build: frontend-install
	@echo "Building frontend for production..."
	cd frontend && npm run build

# Install Go dependencies
deps:
	@echo "Installing Go dependencies..."
	go mod tidy
	go mod download

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Lint Go code
lint:
	@echo "Linting Go code..."
	golangci-lint run

# Run everything
all: deps build docker-run frontend-install
	@echo "All services started successfully!"
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"

# Development mode
dev: deps
	@echo "Starting development mode..."
	@echo "Backend will be available at http://localhost:8080"
	@echo "Frontend will be available at http://localhost:3000"
	@echo "Press Ctrl+C to stop"
	@make -j2 run frontend-start
