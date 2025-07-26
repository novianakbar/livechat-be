# Variables
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= root
DB_PASSWORD ?= postgre123
DB_NAME ?= livechat_db
MIGRATE_PATH ?= ./migrations
BINARY_NAME ?= livechat-be

# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application"
	@echo "  dev             - Run the application in development mode"
	@echo "  test            - Run tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  deps            - Install dependencies"
	@echo "  migrate-up      - Run database migrations"
	@echo "  migrate-down    - Rollback database migrations"
	@echo "  migrate-create  - Create new migration file"
	@echo "  db-create       - Create database"
	@echo "  db-drop         - Drop database"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  oss-setup       - Setup OSS-specific tables and data"
	@echo "  oss-seed        - Seed OSS test data"
	@echo "  oss-test        - Run OSS-specific tests"

# Build the application
.PHONY: build
build:
	@echo "Building application..."
	go build -o dist/$(BINARY_NAME) cmd/main.go

# Run the application
.PHONY: run
run: build
	@echo "Running application..."
	./dist/$(BINARY_NAME)

# Run in development mode
.PHONY: dev
dev:
	@echo "Running in development mode..."
	go run cmd/main.go

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -rf dist/
	go clean

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Database migrations
.PHONY: migrate-up
migrate-up:
	@echo "Running database migrations..."
	migrate -path $(MIGRATE_PATH) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	@echo "Rolling back database migrations..."
	migrate -path $(MIGRATE_PATH) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

.PHONY: migrate-create
migrate-create:
	@echo "Creating new migration file..."
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $$name

# Database operations
.PHONY: db-create
db-create:
	@echo "Creating database..."
	createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)

.PHONY: db-drop
db-drop:
	@echo "Dropping database..."
	dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)

# Docker commands
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(BINARY_NAME)

# Setup development environment
.PHONY: setup
setup: deps
	@echo "Setting up development environment..."
	cp .env.example .env
	@echo "Please edit .env file with your configuration"

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Generate swagger docs
.PHONY: docs
docs:
	@echo "Generating swagger documentation..."
	swag init -g cmd/main.go -o docs/

# Run all checks
.PHONY: check
check: fmt lint test
	@echo "All checks passed!"

# OSS-specific commands
.PHONY: oss-setup
oss-setup:
	@echo "Setting up OSS-specific tables and data..."
	# Add commands to setup OSS-specific tables and data

.PHONY: oss-seed
oss-seed:
	@echo "Seeding OSS test data..."
	# Add commands to seed OSS test data

.PHONY: oss-test
oss-test:
	@echo "Running OSS-specific tests..."
	# Add commands to run OSS-specific tests
