# Task Scheduler Makefile

# Build configuration
BINARY_NAME=task-scheduler
GO_FILES=$(shell find . -name "*.go" -not -path "./vendor/*")

# Build the application
.PHONY: all
all: clean build test

.PHONY: build
build:
	@echo "Building application..."
	@mkdir -p bin
	@go build -o bin/$(BINARY_NAME) cmd/api/main.go
	@echo "Build completed: bin/$(BINARY_NAME)"

# Build for Windows
.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p bin
	@GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME).exe cmd/api/main.go

# Build for Linux
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p bin
	@GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux cmd/api/main.go

# Build for macOS
.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p bin
	@GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin cmd/api/main.go

# Production build (optimized)
.PHONY: prod-build
prod-build:
	@echo "Building for production..."
	@mkdir -p bin
	@CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/$(BINARY_NAME) cmd/api/main.go
	@echo "Production build completed: bin/$(BINARY_NAME)"

# Run the application
.PHONY: run
run:
	@echo "Starting Task Scheduler..."
	@go run cmd/api/main.go

# Run with hot reload (Air)
.PHONY: watch
watch:
	@echo "Starting Task Scheduler with hot reload..."
	@if command -v air > /dev/null 2>&1; then \
		air; \
	else \
		echo "Installing Air..."; \
		go install github.com/air-verse/air@latest; \
		air; \
	fi

# Run with hot reload (PowerShell alternative)
.PHONY: watch-ps
watch-ps:
	@echo "Starting Task Scheduler with PowerShell hot reload..."
	@powershell -ExecutionPolicy Bypass -File scripts/hot-reload.ps1

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./... -v

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -v -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Docker operations
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	@docker build -t task-scheduler:latest .

.PHONY: docker-run
docker-run:
	@echo "Starting application with Docker Compose..."
	@docker-compose up --build -d
	@echo "Application started at http://localhost:8080"
	@echo "Database available at localhost:5432"

.PHONY: docker-down
docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down

.PHONY: docker-logs
docker-logs:
	@docker-compose logs -f app

.PHONY: docker-clean
docker-clean: docker-down
	@echo "Cleaning up Docker resources..."
	@docker-compose down -v
	@docker system prune -f

# Database operations (using docker-compose)
.PHONY: db-up
db-up:
	@echo "Starting database..."
	@docker-compose up -d psql_bp

.PHONY: db-down
db-down:
	@echo "Stopping database..."
	@docker-compose stop psql_bp

# Linting and formatting
.PHONY: lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@gofmt -s -w $(GO_FILES)
	@go mod tidy

# Security scanning
.PHONY: security
security:
	@echo "Running security scan..."
	@if command -v gosec > /dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Installing..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi

# Clean up
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@rm -f $(BINARY_NAME) $(BINARY_NAME).exe *.log

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod verify

# Development setup
.PHONY: dev-setup
dev-setup: deps
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp .env.example .env; fi
	@echo "Please update .env file with your configuration"
	@mkdir -p bin logs
	@echo "Development setup complete!"

# Quality assurance
.PHONY: qa
qa: fmt lint security test
	@echo "Quality assurance checks completed!"

# Help
.PHONY: help
help:
	@echo "Task Scheduler - Available Commands:"
	@echo ""
	@echo "Build Commands:"
	@echo "  build           - Build the application"
	@echo "  build-windows   - Build for Windows"
	@echo "  build-linux     - Build for Linux" 
	@echo "  build-darwin    - Build for macOS"
	@echo "  prod-build      - Build optimized for production"
	@echo ""
	@echo "Development:"
	@echo "  run             - Run the application locally"
	@echo "  watch           - Run with hot reload (auto-restart on changes)"
	@echo "  dev-setup       - Set up development environment"
	@echo "  deps            - Download dependencies"
	@echo ""
	@echo "Testing:"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Start with Docker Compose"
	@echo "  docker-down     - Stop Docker containers"
	@echo "  docker-logs     - View application logs"
	@echo "  docker-clean    - Clean up Docker resources"
	@echo ""
	@echo "Database:"
	@echo "  db-up           - Start database only"
	@echo "  db-down         - Stop database"
	@echo ""
	@echo "Code Quality:"
	@echo "  fmt             - Format code"
	@echo "  lint            - Run linter"
	@echo "  security        - Run security scan"
	@echo "  qa              - Run all quality checks"
	@echo ""
	@echo "Maintenance:"
	@echo "  clean           - Clean build artifacts"
	@echo "  help            - Show this help message"
