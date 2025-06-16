# =============================================================================
# Direito Lux - Makefile
# =============================================================================

.PHONY: help setup up down logs build test clean lint docker-build docker-push

# Default target
.DEFAULT_GOAL := help

# Colors for output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

# Variables
DOCKER_REGISTRY ?= gcr.io/direito-lux
VERSION ?= $(shell git describe --tags --always --dirty)
ENV ?= development

## Display help information
help:
	@echo "$(CYAN)Direito Lux - Development Commands$(RESET)"
	@echo ""
	@echo "$(GREEN)Setup Commands:$(RESET)"
	@awk '/^##.*setup/ { gsub(/##/, "", $$0); print "  " $$0 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Development Commands:$(RESET)"
	@awk '/^##.*dev/ { gsub(/##/, "", $$0); print "  " $$0 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Testing Commands:$(RESET)"
	@awk '/^##.*test/ { gsub(/##/, "", $$0); print "  " $$0 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Docker Commands:$(RESET)"
	@awk '/^##.*docker/ { gsub(/##/, "", $$0); print "  " $$0 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(GREEN)Deployment Commands:$(RESET)"
	@awk '/^##.*deploy/ { gsub(/##/, "", $$0); print "  " $$0 }' $(MAKEFILE_LIST)
	@echo ""

# =============================================================================
# SETUP COMMANDS
# =============================================================================

## setup: Complete local environment setup
setup:
	@echo "$(CYAN)Setting up Direito Lux local environment...$(RESET)"
	@./scripts/setup-local.sh

## setup-env: Create .env file from template
setup-env:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "$(GREEN)Created .env file. Please edit with your configurations.$(RESET)"; \
	else \
		echo "$(YELLOW).env file already exists.$(RESET)"; \
	fi

## setup-data: Create sample data for development
setup-data:
	@echo "$(CYAN)Creating sample data...$(RESET)"
	@./scripts/create-sample-data.sh

## setup-tools: Install development tools
setup-tools:
	@echo "$(CYAN)Installing development tools...$(RESET)"
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@pip install black flake8 mypy

# =============================================================================
# DEVELOPMENT COMMANDS
# =============================================================================

## dev: Start development environment
up:
	@echo "$(CYAN)Starting Direito Lux services...$(RESET)"
	@docker-compose up -d

## dev: Stop development environment
down:
	@echo "$(CYAN)Stopping Direito Lux services...$(RESET)"
	@docker-compose down

## dev: View logs of all services
logs:
	@docker-compose logs -f

## dev: View logs of specific service
logs-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "$(RED)Usage: make logs-service SERVICE=auth-service$(RESET)"; \
		exit 1; \
	fi
	@docker-compose logs -f $(SERVICE)

## dev: Restart specific service
restart:
	@if [ -z "$(SERVICE)" ]; then \
		echo "$(RED)Usage: make restart SERVICE=auth-service$(RESET)"; \
		exit 1; \
	fi
	@docker-compose restart $(SERVICE)

## dev: Rebuild and restart all services
rebuild:
	@echo "$(CYAN)Rebuilding all services...$(RESET)"
	@docker-compose up --build -d

## dev: Show status of all services
status:
	@docker-compose ps

## dev: Clean up development environment
clean:
	@echo "$(CYAN)Cleaning up development environment...$(RESET)"
	@docker-compose down -v
	@docker system prune -f
	@docker volume prune -f

# =============================================================================
# DATABASE COMMANDS
# =============================================================================

## dev: Connect to PostgreSQL
db-connect:
	@docker-compose exec postgres psql -U direito_lux -d direito_lux_dev

## dev: Create database backup
db-backup:
	@echo "$(CYAN)Creating database backup...$(RESET)"
	@mkdir -p backups
	@docker-compose exec -T postgres pg_dump -U direito_lux direito_lux_dev > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql

## dev: Restore database from backup
db-restore:
	@if [ -z "$(BACKUP)" ]; then \
		echo "$(RED)Usage: make db-restore BACKUP=backups/backup_20240115_143022.sql$(RESET)"; \
		exit 1; \
	fi
	@echo "$(CYAN)Restoring database from $(BACKUP)...$(RESET)"
	@docker-compose exec -T postgres psql -U direito_lux -d direito_lux_dev < $(BACKUP)

## dev: Reset database with sample data
db-reset:
	@echo "$(CYAN)Resetting database...$(RESET)"
	@docker-compose exec postgres psql -U direito_lux -d direito_lux_dev -c "DROP SCHEMA IF EXISTS auth CASCADE;"
	@docker-compose exec postgres psql -U direito_lux -d direito_lux_dev -c "DROP SCHEMA IF EXISTS tenant CASCADE;"
	@docker-compose exec postgres psql -U direito_lux -d direito_lux_dev -c "DROP SCHEMA IF EXISTS process CASCADE;"
	@make setup-data

# =============================================================================
# TESTING COMMANDS
# =============================================================================

## test: Run all tests
test:
	@echo "$(CYAN)Running all tests...$(RESET)"
	@$(MAKE) test-go
	@$(MAKE) test-python

## test: Run Go tests
test-go:
	@echo "$(CYAN)Running Go tests...$(RESET)"
	@for dir in services/*/; do \
		if [ -f "$$dir/go.mod" ]; then \
			echo "Testing $$dir"; \
			cd $$dir && go test -v -race -coverprofile=coverage.out ./...; \
			cd ../..; \
		fi \
	done

## test: Run Python tests
test-python:
	@echo "$(CYAN)Running Python tests...$(RESET)"
	@if [ -d "services/ai-service" ]; then \
		cd services/ai-service && python -m pytest tests/ -v; \
	fi

## test: Run integration tests
test-integration:
	@echo "$(CYAN)Running integration tests...$(RESET)"
	@go test -v -tags=integration ./tests/integration/...

## test: Generate test coverage report
test-coverage:
	@echo "$(CYAN)Generating coverage report...$(RESET)"
	@mkdir -p coverage
	@go test -coverprofile=coverage/coverage.out ./...
	@go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "$(GREEN)Coverage report: coverage/coverage.html$(RESET)"

# =============================================================================
# LINTING COMMANDS
# =============================================================================

## dev: Run all linters
lint:
	@$(MAKE) lint-go
	@$(MAKE) lint-python

## dev: Run Go linter
lint-go:
	@echo "$(CYAN)Running Go linter...$(RESET)"
	@golangci-lint run ./...

## dev: Run Python linter
lint-python:
	@echo "$(CYAN)Running Python linter...$(RESET)"
	@if [ -d "services/ai-service" ]; then \
		cd services/ai-service && black --check . && flake8 . && mypy .; \
	fi

## dev: Fix Go formatting
fmt-go:
	@echo "$(CYAN)Formatting Go code...$(RESET)"
	@gofmt -w .
	@goimports -w .

## dev: Fix Python formatting
fmt-python:
	@echo "$(CYAN)Formatting Python code...$(RESET)"
	@if [ -d "services/ai-service" ]; then \
		cd services/ai-service && black .; \
	fi

# =============================================================================
# PROTOBUF GENERATION
# =============================================================================

## dev: Generate protobuf files
proto-gen:
	@echo "$(CYAN)Generating protobuf files...$(RESET)"
	@mkdir -p shared/proto/gen
	@protoc --go_out=shared/proto/gen --go_opt=paths=source_relative \
		--go-grpc_out=shared/proto/gen --go-grpc_opt=paths=source_relative \
		shared/proto/*.proto

# =============================================================================
# DOCKER COMMANDS
# =============================================================================

## docker: Build all Docker images
docker-build:
	@echo "$(CYAN)Building Docker images...$(RESET)"
	@for service in auth-service tenant-service process-service datajud-service notification-service ai-service; do \
		if [ -d "services/$$service" ]; then \
			echo "Building $$service..."; \
			docker build -t $(DOCKER_REGISTRY)/$$service:$(VERSION) services/$$service/; \
		fi \
	done

## docker: Build specific service Docker image
docker-build-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "$(RED)Usage: make docker-build-service SERVICE=auth-service$(RESET)"; \
		exit 1; \
	fi
	@echo "$(CYAN)Building $(SERVICE) Docker image...$(RESET)"
	@docker build -t $(DOCKER_REGISTRY)/$(SERVICE):$(VERSION) services/$(SERVICE)/

## docker: Push all Docker images
docker-push:
	@echo "$(CYAN)Pushing Docker images...$(RESET)"
	@for service in auth-service tenant-service process-service datajud-service notification-service ai-service; do \
		if [ -d "services/$$service" ]; then \
			echo "Pushing $$service..."; \
			docker push $(DOCKER_REGISTRY)/$$service:$(VERSION); \
		fi \
	done

## docker: Build and push all images
docker-release: docker-build docker-push

# =============================================================================
# KUBERNETES COMMANDS
# =============================================================================

## deploy: Deploy to local Kubernetes
k8s-deploy-local:
	@echo "$(CYAN)Deploying to local Kubernetes...$(RESET)"
	@kubectl apply -k infrastructure/kubernetes/overlays/local

## deploy: Deploy to staging
k8s-deploy-staging:
	@echo "$(CYAN)Deploying to staging...$(RESET)"
	@kubectl apply -k infrastructure/kubernetes/overlays/staging

## deploy: Deploy to production
k8s-deploy-production:
	@echo "$(CYAN)Deploying to production...$(RESET)"
	@kubectl apply -k infrastructure/kubernetes/overlays/production

## deploy: Delete local deployment
k8s-delete-local:
	@kubectl delete -k infrastructure/kubernetes/overlays/local

# =============================================================================
# TERRAFORM COMMANDS
# =============================================================================

## deploy: Plan Terraform changes for development
tf-plan-dev:
	@echo "$(CYAN)Planning Terraform for development...$(RESET)"
	@cd infrastructure/terraform/environments/development && terraform plan

## deploy: Apply Terraform for development
tf-apply-dev:
	@echo "$(CYAN)Applying Terraform for development...$(RESET)"
	@cd infrastructure/terraform/environments/development && terraform apply

## deploy: Plan Terraform changes for production
tf-plan-prod:
	@echo "$(CYAN)Planning Terraform for production...$(RESET)"
	@cd infrastructure/terraform/environments/production && terraform plan

## deploy: Apply Terraform for production
tf-apply-prod:
	@echo "$(CYAN)Applying Terraform for production...$(RESET)"
	@cd infrastructure/terraform/environments/production && terraform apply

# =============================================================================
# MONITORING COMMANDS
# =============================================================================

## dev: Open Grafana dashboard
open-grafana:
	@open http://localhost:3000

## dev: Open Prometheus
open-prometheus:
	@open http://localhost:9090

## dev: Open Jaeger
open-jaeger:
	@open http://localhost:16686

## dev: Open pgAdmin
open-pgadmin:
	@open http://localhost:5050

## dev: Open RabbitMQ Management
open-rabbitmq:
	@open http://localhost:15672

# =============================================================================
# UTILITY COMMANDS
# =============================================================================

## dev: Show environment info
info:
	@echo "$(CYAN)Direito Lux Environment Info$(RESET)"
	@echo "Version: $(VERSION)"
	@echo "Environment: $(ENV)"
	@echo "Docker Registry: $(DOCKER_REGISTRY)"
	@echo ""
	@echo "$(GREEN)Service Status:$(RESET)"
	@docker-compose ps

## dev: Generate API documentation
docs-api:
	@echo "$(CYAN)Generating API documentation...$(RESET)"
	@mkdir -p docs/api
	@swag init -g main.go -o docs/api

## dev: Start development with file watching
dev-watch:
	@echo "$(CYAN)Starting development with file watching...$(RESET)"
	@air

## dev: Check system health
health-check:
	@echo "$(CYAN)Checking system health...$(RESET)"
	@curl -f http://localhost:8000/health || echo "$(RED)API Gateway not responding$(RESET)"
	@curl -f http://localhost:3000/api/health || echo "$(RED)Grafana not responding$(RESET)"
	@curl -f http://localhost:9090/-/healthy || echo "$(RED)Prometheus not responding$(RESET)"

# =============================================================================
# GIT HOOKS
# =============================================================================

## dev: Install Git hooks
git-hooks:
	@echo "$(CYAN)Installing Git hooks...$(RESET)"
	@cp scripts/git-hooks/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "$(GREEN)Git hooks installed$(RESET)"

# =============================================================================
# PERFORMANCE
# =============================================================================

## test: Run performance tests
perf-test:
	@echo "$(CYAN)Running performance tests...$(RESET)"
	@if command -v k6 > /dev/null; then \
		k6 run tests/performance/load-test.js; \
	else \
		echo "$(RED)k6 not installed. Install from https://k6.io$(RESET)"; \
	fi

## test: Run benchmark tests
benchmark:
	@echo "$(CYAN)Running benchmark tests...$(RESET)"
	@go test -bench=. -benchmem ./...

# =============================================================================
# SECURITY
# =============================================================================

## test: Run security scan
security-scan:
	@echo "$(CYAN)Running security scan...$(RESET)"
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "$(RED)gosec not installed. Run: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest$(RESET)"; \
	fi

## test: Check for vulnerabilities
vuln-check:
	@echo "$(CYAN)Checking for vulnerabilities...$(RESET)"
	@go list -json -m all | nancy sleuth

# =============================================================================
# EXAMPLES
# =============================================================================

## Examples of common commands:
## make setup                    # Complete setup
## make up                       # Start services
## make logs-service SERVICE=auth-service
## make test                     # Run all tests
## make docker-build-service SERVICE=auth-service
## make k8s-deploy-local         # Deploy to local k8s