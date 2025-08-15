# USL Server Makefile

.PHONY: test test-unit test-integration test-coverage build run clean lint fmt help

# Default target
all: test build

# Build the server
build:
	@echo "Building USL server..."
	go build -o bin/server ./cmd/server

# Run the server
run:
	@echo "Starting USL server..."
	go run ./cmd/server

# Run all tests
test: test-unit test-integration

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	go test -v ./internal/handlers/... ./internal/models/... ./internal/repositories/...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v ./test/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Test the HTMX templates specifically
test-templates:
	@echo "Testing HTMX templates..."
	go test -v ./internal/handlers/ -run TestUserHandler

# Lint the code
lint:
	@echo "Running linter..."
	golangci-lint run

# Format the code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet the code
vet:
	@echo "Vetting code..."
	go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Build CSS (Tailwind)
build-css:
	@echo "Building Tailwind CSS..."
	npx tailwindcss -i ./static/src/input.css -o ./static/dist/output.css --watch

# Build CSS for production
build-css-prod:
	@echo "Building Tailwind CSS for production..."
	npx tailwindcss -i ./static/src/input.css -o ./static/dist/output.css --minify

# Install dependencies
deps:
	@echo "Installing Go dependencies..."
	go mod download
	@echo "Installing Node.js dependencies..."
	npm install

# Run development server with auto-reload
dev: build-css
	@echo "Starting development server..."
	air

# Database migrations (placeholder)
migrate-up:
	@echo "Running database migrations..."
	# Add migration commands here

migrate-down:
	@echo "Rolling back database migrations..."
	# Add rollback commands here

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t usl-server .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 usl-server

# Help target
help:
	@echo "Available targets:"
	@echo "  build          - Build the server binary"
	@echo "  run            - Run the server"
	@echo "  test           - Run all tests"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-templates - Test HTMX templates specifically"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  vet            - Vet code"
	@echo "  clean          - Clean build artifacts"
	@echo "  build-css      - Build Tailwind CSS (watch mode)"
	@echo "  build-css-prod - Build Tailwind CSS for production"
	@echo "  deps           - Install dependencies"
	@echo "  dev            - Run development server with auto-reload"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  help           - Show this help message"