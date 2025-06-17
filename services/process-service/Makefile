# =============================================================================
# Process Service - Makefile
# =============================================================================

.PHONY: help build run test lint clean docker-build docker-run

# Variables
SERVICE_NAME := process-service
DOCKER_IMAGE := direito-lux/$(SERVICE_NAME)
VERSION := $(shell git describe --tags --always --dirty)

## Display help
help:
	@echo "Process Service Commands:"
	@echo ""
	@grep -E '^##' $(MAKEFILE_LIST) | sed 's/##//g'

## Build the service
build:
	@echo "Building process-service..."
	@go build -o bin/$(SERVICE_NAME) cmd/server/main.go

## Run the service locally
run:
	@echo "Running process-service..."
	@go run cmd/server/main.go

## Run with live reload
dev:
	@echo "Running process-service with live reload..."
	@air

## Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

## Generate test coverage report
test-coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

## Format code
fmt:
	@echo "Formatting code..."
	@gofmt -w .
	@goimports -w .

## Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/ tmp/ coverage.out coverage.html

## Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(VERSION) .
	@docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest

## Build development Docker image
docker-build-dev:
	@echo "Building development Docker image..."
	@docker build -f Dockerfile.dev -t $(DOCKER_IMAGE):dev .

## Run with Docker
docker-run:
	@echo "Running with Docker..."
	@docker run --rm -p 8083:8080 $(DOCKER_IMAGE):latest

## Generate API documentation
docs:
	@echo "Generating API docs..."
	@swag init -g cmd/server/main.go -o docs/

## Database migration up
migrate-up:
	@echo "Running database migrations..."
	@migrate -path migrations -database "postgres://direito_lux:dev_password_123@localhost:5432/direito_lux_dev?sslmode=disable" up

## Database migration down
migrate-down:
	@echo "Rolling back database migrations..."
	@migrate -path migrations -database "postgres://direito_lux:dev_password_123@localhost:5432/direito_lux_dev?sslmode=disable" down

## Create new migration
migrate-create:
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=migration_name"; exit 1; fi
	@migrate create -ext sql -dir migrations $(NAME)

## CQRS specific commands

## Rebuild read models
rebuild-read-models:
	@echo "Rebuilding CQRS read models..."
	@go run cmd/tools/rebuild-read-models.go

## Sync event store
sync-event-store:
	@echo "Syncing event store..."
	@go run cmd/tools/sync-event-store.go

## Process monitoring commands

## Start monitoring
start-monitoring:
	@echo "Starting process monitoring..."
	@go run cmd/tools/start-monitoring.go

## Sync all processes
sync-processes:
	@echo "Syncing all processes with DataJud..."
	@go run cmd/tools/sync-processes.go