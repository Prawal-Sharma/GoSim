.PHONY: build run clean test deps

# Build the server binary
build:
	go build -o gosim cmd/server/main.go

# Run the server
run:
	go run cmd/server/main.go

# Download dependencies
deps:
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	rm -f gosim
	go clean

# Run tests
test:
	go test ./...

# Run with hot reload (requires air)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Running without hot reload..."; \
		go run cmd/server/main.go; \
	fi

# Install development tools
install-tools:
	go install github.com/cosmtrek/air@latest

# Build for production
prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gosim cmd/server/main.go

# Format code
fmt:
	go fmt ./...

# Check for issues
lint:
	@if command -v golint > /dev/null; then \
		golint ./...; \
	else \
		echo "golint not installed. Skipping..."; \
	fi

# Run server with race detection
race:
	go run -race cmd/server/main.go

help:
	@echo "Available commands:"
	@echo "  make build   - Build the server binary"
	@echo "  make run     - Run the server"
	@echo "  make deps    - Download dependencies"
	@echo "  make clean   - Clean build artifacts"
	@echo "  make test    - Run tests"
	@echo "  make dev     - Run with hot reload"
	@echo "  make prod    - Build for production"
	@echo "  make fmt     - Format code"
	@echo "  make lint    - Check for issues"
	@echo "  make race    - Run with race detection"