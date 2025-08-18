# USL Server Makefile

.PHONY: test test-unit test-integration test-coverage test-templates-full test-template-contracts test-template-startup test-template-regression test-template-security test-template-performance test-pre-commit build run clean lint fmt help

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
test: test-unit test-integration test-smoke test-template-contracts

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

# Template validation test suite

# Run full template validation suite
test-templates-full: test-template-contracts test-template-startup test-template-regression test-template-security
	@echo "✅ All template validation tests completed"

# Test template-handler contracts (catches Issue #35 type bugs)
test-template-contracts:
	@echo "Testing template-handler contracts..."
	go test -v ./internal/usl/handlers/ -run TestTemplateHandlerContracts

# Test startup-time template validation
test-template-startup:
	@echo "Testing startup template validation..."
	go test -v ./internal/usl/handlers/ -run TestStartupTemplateValidation

# Test for specific regressions like Issue #35
test-template-regression:
	@echo "Testing template regressions..."
	go test -v ./internal/usl/handlers/ -run TestIssue35

# Test template security (XSS, injection)
test-template-security:
	@echo "Testing template security..."
	go test -v ./internal/usl/handlers/ -run TestTemplateSecurityEdgeCases

# Test template performance under load
test-template-performance:
	@echo "Testing template performance..."
	go test -v ./internal/usl/handlers/ -run TestTemplatePerformanceUnderLoad

# Property-based template testing
test-template-property:
	@echo "Running property-based template tests..."
	go test -v ./internal/usl/handlers/ -run TestPropertyBasedTemplateValidation

# Fast pre-commit template validation
test-pre-commit:
	@echo "Running pre-commit template validation..."
	go test -v ./test/ -run TestCIPreCommitValidation

# Integration template validation for CI/CD
test-template-ci:
	@echo "Running CI template validation..."
	go test -v ./test/ -run TestFullTemplateValidationSuite

# Legacy template test (keeping for compatibility)
test-templates:
	@echo "Testing HTMX templates..."
	go test -v ./internal/handlers/ -run TestUserHandler

# Run smoke tests for critical functionality
test-smoke:
	@echo "Running smoke tests..."
	go test -v ./test/ -run TestSmoke

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
	@echo "  build                  - Build the server binary"
	@echo "  run                    - Run the server"
	@echo "  test                   - Run all tests"
	@echo "  test-unit              - Run unit tests only"
	@echo "  test-integration       - Run integration tests only"
	@echo "  test-coverage          - Run tests with coverage report"
	@echo ""
	@echo "Template Validation:"
	@echo "  test-templates-full    - Run complete template validation suite"
	@echo "  test-template-contracts - Test template-handler contracts (prevents Issue #35)"
	@echo "  test-template-startup  - Test startup-time template validation"
	@echo "  test-template-regression - Test for specific template regressions"
	@echo "  test-template-security - Test template security (XSS, injection)"
	@echo "  test-template-performance - Test template performance under load"
	@echo "  test-template-property - Property-based template testing"
	@echo "  test-pre-commit        - Fast pre-commit template validation"
	@echo "  test-template-ci       - CI/CD template validation suite"
	@echo ""
	@echo "Development:"
	@echo "  lint                   - Run linter"
	@echo "  fmt                    - Format code"
	@echo "  vet                    - Vet code"
	@echo "  clean                  - Clean build artifacts"
	@echo "  build-css              - Build Tailwind CSS (watch mode)"
	@echo "  build-css-prod         - Build Tailwind CSS for production"
	@echo "  deps                   - Install dependencies"
	@echo "  dev                    - Run development server with auto-reload"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build           - Build Docker image"
	@echo "  docker-run             - Run Docker container"
	@echo ""
	@echo "  help                   - Show this help message"