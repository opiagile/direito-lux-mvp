#!/bin/bash

# =============================================================================
# Script para criar novo microserviço a partir do template
# =============================================================================

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verificar parâmetros
if [ $# -ne 1 ]; then
    echo "Uso: $0 <nome-do-servico>"
    echo ""
    echo "Exemplos:"
    echo "  $0 auth-service"
    echo "  $0 tenant-service"
    echo "  $0 process-service"
    exit 1
fi

SERVICE_NAME=$1

# Validar nome do serviço
if [[ ! $SERVICE_NAME =~ ^[a-z][a-z0-9\-]*[a-z0-9]$ ]]; then
    log_error "Nome do serviço deve estar em lowercase e usar apenas letras, números e hífens"
    log_error "Exemplo: auth-service, tenant-service"
    exit 1
fi

# Verificar se template existe
if [ ! -d "template-service" ]; then
    log_error "Diretório template-service não encontrado"
    log_error "Execute este script a partir do diretório raiz do projeto"
    exit 1
fi

# Verificar se serviço já existe
if [ -d "services/$SERVICE_NAME" ]; then
    log_error "Serviço $SERVICE_NAME já existe em services/$SERVICE_NAME"
    exit 1
fi

# Banner
echo -e "${BLUE}"
cat << "EOF"
╔══════════════════════════════════════════════════════════════╗
║                     DIREITO LUX                             ║
║                 Criar Novo Microserviço                     ║
╚══════════════════════════════════════════════════════════════╝
EOF
echo -e "${NC}"

log_info "Criando serviço: $SERVICE_NAME"

# Criar diretório do serviço
mkdir -p "services/$SERVICE_NAME"

# Copiar template
log_info "Copiando template..."
cp -r template-service/* "services/$SERVICE_NAME/"
cp template-service/.air.toml "services/$SERVICE_NAME/"

# Criar diretórios adicionais se não existirem
mkdir -p "services/$SERVICE_NAME/migrations"
mkdir -p "services/$SERVICE_NAME/tests/unit"
mkdir -p "services/$SERVICE_NAME/tests/integration"
mkdir -p "services/$SERVICE_NAME/docs"

# Converter nome para diferentes formatos
SERVICE_NAME_UNDERSCORE=$(echo $SERVICE_NAME | tr '-' '_')
SERVICE_NAME_CAMEL=$(echo $SERVICE_NAME | sed -r 's/(^|-)([a-z])/\U\2/g' | sed 's/-//g')
SERVICE_NAME_TITLE=$(echo $SERVICE_NAME | sed 's/-/ /g' | sed 's/\b\w/\U&/g')

log_info "Atualizando arquivos..."

# Atualizar go.mod
sed -i "s|github.com/direito-lux/template-service|github.com/direito-lux/$SERVICE_NAME|g" "services/$SERVICE_NAME/go.mod"

# Atualizar imports em todos os arquivos Go
find "services/$SERVICE_NAME" -name "*.go" -exec sed -i "s|github.com/direito-lux/template-service|github.com/direito-lux/$SERVICE_NAME|g" {} +

# Atualizar README
sed -i "s/Template Service/$SERVICE_NAME_TITLE Service/g" "services/$SERVICE_NAME/README.md"
sed -i "s/template-service/$SERVICE_NAME/g" "services/$SERVICE_NAME/README.md"
sed -i "s/auth-service/$SERVICE_NAME/g" "services/$SERVICE_NAME/README.md"

# Atualizar configurações padrão
sed -i "s/template-service/$SERVICE_NAME/g" "services/$SERVICE_NAME/internal/infrastructure/config/config.go"

# Atualizar Air configuration
sed -i "s/template-service/$SERVICE_NAME/g" "services/$SERVICE_NAME/.air.toml"

# Atualizar Dockerfiles
sed -i "s/Template Service/$SERVICE_NAME_TITLE Service/g" "services/$SERVICE_NAME/Dockerfile"
sed -i "s/template-service/$SERVICE_NAME/g" "services/$SERVICE_NAME/Dockerfile"
sed -i "s/template-service/$SERVICE_NAME/g" "services/$SERVICE_NAME/Dockerfile.dev"

# Criar arquivos específicos do serviço

# .env.example específico
cat > "services/$SERVICE_NAME/.env.example" << EOF
# =============================================================================
# $SERVICE_NAME_TITLE Service - Environment Variables
# =============================================================================

# Application
SERVICE_NAME=$SERVICE_NAME
PORT=8080
LOG_LEVEL=debug
ENVIRONMENT=development
VERSION=1.0.0

# Database
DB_HOST=postgres
DB_PORT=5432
DB_NAME=direito_lux_dev
DB_USER=direito_lux
DB_PASSWORD=dev_password_123
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=300s

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=dev_redis_123
REDIS_DATABASE=0
REDIS_POOL_SIZE=10

# RabbitMQ
RABBITMQ_URL=amqp://direito_lux:dev_rabbit_123@rabbitmq:5672/direito_lux
RABBITMQ_EXCHANGE=direito_lux.events
RABBITMQ_ROUTING_KEY=$SERVICE_NAME_UNDERSCORE
RABBITMQ_QUEUE=$SERVICE_NAME_UNDERSCORE.events

# Observabilidade
JAEGER_ENDPOINT=http://jaeger:14268/api/traces
JAEGER_SERVICE_NAME=$SERVICE_NAME
METRICS_ENABLED=true
METRICS_PORT=9090

# External Services
AUTH_SERVICE_URL=http://auth-service:8080
TENANT_SERVICE_URL=http://tenant-service:8080
PROCESS_SERVICE_URL=http://process-service:8080
DATAJUD_SERVICE_URL=http://datajud-service:8080
NOTIFICATION_SERVICE_URL=http://notification-service:8080
AI_SERVICE_URL=http://ai-service:8000
EOF

# Makefile específico
cat > "services/$SERVICE_NAME/Makefile" << EOF
# =============================================================================
# $SERVICE_NAME_TITLE Service - Makefile
# =============================================================================

.PHONY: help build run test lint clean docker-build docker-run

# Variables
SERVICE_NAME := $SERVICE_NAME
DOCKER_IMAGE := direito-lux/\$(SERVICE_NAME)
VERSION := \$(shell git describe --tags --always --dirty)

## Display help
help:
	@echo "$SERVICE_NAME_TITLE Service Commands:"
	@echo ""
	@grep -E '^##' \$(MAKEFILE_LIST) | sed 's/##//g'

## Build the service
build:
	@echo "Building $SERVICE_NAME..."
	@go build -o bin/\$(SERVICE_NAME) cmd/server/main.go

## Run the service locally
run:
	@echo "Running $SERVICE_NAME..."
	@go run cmd/server/main.go

## Run with live reload
dev:
	@echo "Running $SERVICE_NAME with live reload..."
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
	@docker build -t \$(DOCKER_IMAGE):\$(VERSION) .
	@docker tag \$(DOCKER_IMAGE):\$(VERSION) \$(DOCKER_IMAGE):latest

## Build development Docker image
docker-build-dev:
	@echo "Building development Docker image..."
	@docker build -f Dockerfile.dev -t \$(DOCKER_IMAGE):dev .

## Run with Docker
docker-run:
	@echo "Running with Docker..."
	@docker run --rm -p 8080:8080 \$(DOCKER_IMAGE):latest

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
	@if [ -z "\$(NAME)" ]; then echo "Usage: make migrate-create NAME=migration_name"; exit 1; fi
	@migrate create -ext sql -dir migrations \$(NAME)
EOF

# .gitignore específico
cat > "services/$SERVICE_NAME/.gitignore" << EOF
# Binários
bin/
tmp/

# Logs
*.log

# Coverage
coverage.out
coverage.html

# Environment
.env

# IDE
.vscode/
.idea/

# OS
.DS_Store
Thumbs.db

# Go
vendor/
EOF

# Atualizar docker-compose.yml para incluir o novo serviço
log_info "Atualizando docker-compose.yml..."

# Definir porta baseada no nome do serviço
case $SERVICE_NAME in
    "auth-service")
        PORT="8081"
        ;;
    "tenant-service")
        PORT="8082"
        ;;
    "process-service")
        PORT="8083"
        ;;
    "datajud-service")
        PORT="8084"
        ;;
    "notification-service")
        PORT="8085"
        ;;
    "ai-service")
        PORT="8086"
        ;;
    *)
        PORT="8087"
        ;;
esac

# Adicionar serviço ao docker-compose.yml (se não existir)
if ! grep -q "  $SERVICE_NAME:" docker-compose.yml; then
    # Encontrar onde inserir o novo serviço (antes das ferramentas de desenvolvimento)
    sed -i "/# =============================================================================\n# DESENVOLVIMENTO E UTILIDADES/i\\
  # $SERVICE_NAME_TITLE\\
  $SERVICE_NAME:\\
    build:\\
      context: ./services/$SERVICE_NAME\\
      dockerfile: Dockerfile.dev\\
    container_name: direito-lux-$SERVICE_NAME\\
    environment:\\
      - PORT=8080\\
      - DB_HOST=postgres\\
      - DB_PORT=5432\\
      - DB_NAME=direito_lux_dev\\
      - DB_USER=direito_lux\\
      - DB_PASSWORD=dev_password_123\\
      - REDIS_HOST=redis\\
      - REDIS_PORT=6379\\
      - REDIS_PASSWORD=dev_redis_123\\
      - RABBITMQ_URL=amqp://direito_lux:dev_rabbit_123@rabbitmq:5672/direito_lux\\
      - JAEGER_ENDPOINT=http://jaeger:14268/api/traces\\
      - LOG_LEVEL=debug\\
    ports:\\
      - \"$PORT:8080\"\\
    depends_on:\\
      postgres:\\
        condition: service_healthy\\
      redis:\\
        condition: service_healthy\\
      rabbitmq:\\
        condition: service_healthy\\
    volumes:\\
      - ./services/$SERVICE_NAME:/app\\
      - /app/vendor\\
    networks:\\
      - direito-lux-network\\
    restart: unless-stopped\\
\\
" docker-compose.yml
fi

# Inicializar go mod
log_info "Inicializando Go module..."
cd "services/$SERVICE_NAME"
go mod tidy
cd ../..

log_success "Serviço $SERVICE_NAME criado com sucesso!"

echo ""
echo -e "${GREEN}📁 Serviço criado em: ${BLUE}services/$SERVICE_NAME${NC}"
echo ""
echo -e "${YELLOW}🔧 Próximos passos:${NC}"
echo ""
echo -e "${BLUE}1. Editar configurações:${NC}"
echo "   cd services/$SERVICE_NAME"
echo "   cp .env.example .env"
echo "   # Editar .env com configurações específicas"
echo ""
echo -e "${BLUE}2. Implementar domínio:${NC}"
echo "   # Criar entidades em internal/domain/"
echo "   # Implementar casos de uso em internal/application/"
echo "   # Implementar handlers HTTP em internal/infrastructure/http/"
echo ""
echo -e "${BLUE}3. Executar localmente:${NC}"
echo "   cd services/$SERVICE_NAME"
echo "   make run"
echo "   # ou com live reload:"
echo "   make dev"
echo ""
echo -e "${BLUE}4. Executar com Docker:${NC}"
echo "   docker-compose up $SERVICE_NAME"
echo ""
echo -e "${BLUE}5. Executar testes:${NC}"
echo "   cd services/$SERVICE_NAME"
echo "   make test"
echo ""
echo -e "${GREEN}✅ URLs do serviço:${NC}"
echo "   • API: http://localhost:$PORT"
echo "   • Health: http://localhost:$PORT/health"
echo "   • Metrics: http://localhost:9090/metrics"
echo ""