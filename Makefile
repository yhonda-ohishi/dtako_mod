.PHONY: help test test-contract test-integration setup-db run build clean swag swag-init

# Default target
help:
	@echo "Available targets:"
	@echo "  make test          - Run all tests"
	@echo "  make test-contract - Run contract tests only"
	@echo "  make test-integration - Run integration tests only"
	@echo "  make setup-db      - Setup test databases"
	@echo "  make run           - Run the application"
	@echo "  make build         - Build the application"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make swag          - Generate Swagger documentation"
	@echo "  make swag-init     - Initialize and generate Swagger documentation"

# Run all tests
test:
	go test ./...

# Run contract tests
test-contract:
	go test ./tests/contract/...

# Run integration tests
test-integration:
	go test ./tests/integration/...

# Setup test databases
setup-db:
	go run cmd/setup-test-db/main.go

# Run the application
run:
	go run main.go

# Build the application
build:
	go build -o bin/dtako_mod main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Generate test coverage report
coverage-report:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Generate Swagger documentation
swag:
	swag init -g main.go --output docs
	@echo "Swagger documentation generated in docs/"

# Initialize and generate Swagger documentation (first time setup)
swag-init:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g main.go --output docs
	@echo "Swagger documentation initialized and generated in docs/"

# Run application with swagger docs regeneration
run-with-docs: swag run