# Acoustics Calculator Makefile

.PHONY: run build test fmt clean deps

# Default target
all: fmt test build

# Run the application
run:
	@echo "Starting Acoustics Calculator server..."
	go run cmd/server/main.go

# Build the application
build:
	@echo "Building Acoustics Calculator..."
	@mkdir -p bin
	go build -o bin/acoustics-calculator cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f acoustics.db*

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Development setup
setup: deps
	@echo "Setting up development environment..."
	@echo "Run 'make run' to start the server"

# Lint code (requires golangci-lint)
lint:
	@echo "Running linter..."
	golangci-lint run

# Database reset
reset-db:
	@echo "Resetting database..."
	rm -f acoustics.db
	@echo "Database reset. Run 'make run' to recreate it."

# Development with hot reload (requires air)
dev:
	@echo "Starting development server with hot reload..."
	air

# Help
help:
	@echo "Available commands:"
	@echo "  run      - Start the application server"
	@echo "  build    - Build the application binary"
	@echo "  test     - Run all tests"
	@echo "  fmt      - Format Go code"
	@echo "  clean    - Remove build artifacts"
	@echo "  deps     - Install/update dependencies"
	@echo "  setup    - Initial development setup"
	@echo "  lint     - Run code linter"
	@echo "  reset-db - Reset the database"
	@echo "  dev      - Start with hot reload (requires air)"
	@echo "  help     - Show this help message"
