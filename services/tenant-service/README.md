# Template Service - Direito Lux

Template de microserviço Go seguindo os padrões Full Cycle para o projeto Direito Lux.

## 🎯 Visão Geral

Este template implementa as melhores práticas para microserviços Go, incluindo:

- **Hexagonal Architecture** (Ports & Adapters)
- **Event-Driven Architecture** com RabbitMQ
- **Observabilidade** completa (Logs, Metrics, Tracing)
- **Multi-tenancy** com isolamento por tenant
- **Clean Code** e SOLID principles

## 🏗️ Arquitetura

```
template-service/
├── cmd/
│   └── server/           # Entry point da aplicação
├── internal/
│   ├── domain/           # Regras de negócio
│   ├── application/      # Casos de uso
│   └── infrastructure/   # Implementações técnicas
│       ├── config/       # Configurações
│       ├── database/     # Conexão com banco
│       ├── events/       # Event bus
│       ├── http/         # Servidor HTTP
│       ├── logging/      # Sistema de logs
│       ├── metrics/      # Métricas Prometheus
│       └── tracing/      # Distributed tracing
├── migrations/           # Migrações do banco
├── tests/               # Testes
└── docs/                # Documentação
```

## 🚀 Quick Start

### Desenvolvimento Local

```bash
# 1. Copiar template para novo serviço
cp -r template-service auth-service
cd auth-service

# 2. Atualizar go.mod
go mod edit -module github.com/direito-lux/auth-service

# 3. Atualizar imports no código
find . -name "*.go" -exec sed -i 's|github.com/direito-lux/template-service|github.com/direito-lux/auth-service|g' {} +

# 4. Configurar variáveis de ambiente
cp .env.example .env
# Editar .env com suas configurações

# 5. Executar localmente
go run cmd/server/main.go

# Ou usar live reload
air
```

### Docker Development

```bash
# Build da imagem de desenvolvimento
docker build -f Dockerfile.dev -t direito-lux/auth-service:dev .

# Executar com Docker Compose (do diretório raiz)
docker-compose up auth-service
```

## 📋 Funcionalidades Implementadas

### ✅ Infraestrutura Base
- [x] Configuração via environment variables
- [x] Logging estruturado com Zap
- [x] Metrics com Prometheus
- [x] Distributed tracing com Jaeger
- [x] Health checks (liveness/readiness)
- [x] Graceful shutdown

### ✅ HTTP Server
- [x] Gin framework com middlewares
- [x] CORS configurável
- [x] Rate limiting
- [x] Request ID tracking
- [x] Tenant isolation middleware

### ✅ Database
- [x] PostgreSQL com SQLx
- [x] Connection pooling
- [x] Database migrations
- [x] Health checks
- [x] Metrics de performance

### ✅ Event System
- [x] Event-driven architecture
- [x] RabbitMQ integration
- [x] Event sourcing support
- [x] Retry mechanism
- [x] Dead letter queue

### ✅ Observabilidade
- [x] Structured logging
- [x] Context propagation
- [x] Performance metrics
- [x] Business metrics
- [x] Error tracking

## 🔧 Configuração

### Variáveis de Ambiente

| Variável | Descrição | Padrão |
|----------|-----------|---------|
| `SERVICE_NAME` | Nome do serviço | `template-service` |
| `PORT` | Porta HTTP | `8080` |
| `LOG_LEVEL` | Nível de log | `info` |
| `ENVIRONMENT` | Ambiente | `development` |
| `DB_HOST` | Host PostgreSQL | `localhost` |
| `DB_PORT` | Porta PostgreSQL | `5432` |
| `REDIS_HOST` | Host Redis | `localhost` |
| `RABBITMQ_URL` | URL RabbitMQ | `amqp://...` |
| `JAEGER_ENDPOINT` | Endpoint Jaeger | `http://localhost:14268` |

### Configuração de Desenvolvimento

```bash
# .env para desenvolvimento local
SERVICE_NAME=auth-service
PORT=8081
LOG_LEVEL=debug
ENVIRONMENT=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=direito_lux_dev
DB_USER=direito_lux
DB_PASSWORD=dev_password_123

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=dev_redis_123

# RabbitMQ
RABBITMQ_URL=amqp://direito_lux:dev_rabbit_123@localhost:5672/direito_lux

# Observabilidade
JAEGER_ENDPOINT=http://localhost:14268/api/traces
METRICS_ENABLED=true
METRICS_PORT=9090
```

## 📊 APIs Expostas

### Health Checks
```http
GET /health          # Liveness probe
GET /ready           # Readiness probe
```

### Example Endpoints
```http
GET    /api/v1/ping              # Teste simples
GET    /api/v1/templates         # Listar templates
POST   /api/v1/templates         # Criar template
GET    /api/v1/templates/:id     # Buscar template
PUT    /api/v1/templates/:id     # Atualizar template
DELETE /api/v1/templates/:id     # Remover template
```

### Métricas
```http
GET /metrics         # Métricas Prometheus (porta 9090)
```

## 🧪 Testes

```bash
# Testes unitários
go test ./...

# Testes com coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Testes de integração
go test -tags=integration ./...

# Benchmark
go test -bench=. ./...
```

## 📈 Métricas Coletadas

### HTTP Metrics
- `http_requests_total` - Total de requisições HTTP
- `http_request_duration_seconds` - Duração das requisições
- `http_requests_in_flight` - Requisições em andamento

### Database Metrics
- `database_queries_total` - Total de queries
- `database_query_duration_seconds` - Duração das queries
- `database_connections` - Status das conexões

### Business Metrics
- `tenant_operations_total` - Operações por tenant
- `user_operations_total` - Operações de usuário

## 🔍 Logging

### Structured Logging
```go
// Exemplo de log estruturado
logging.LogInfo(ctx, logger, "Usuário criado",
    zap.String("user_id", userID),
    zap.String("tenant_id", tenantID),
    zap.String("operation", "create_user"),
)
```

### Context Propagation
```go
// Propagar informações via contexto
ctx = logging.WithTenantID(ctx, tenantID)
ctx = logging.WithUserID(ctx, userID)
```

## 🎯 Distributed Tracing

```go
// Exemplo de tracing
span, ctx := tracing.StartSpanFromContext(ctx, "create_user")
defer span.Finish()

// Adicionar tags
span.SetTag("user.id", userID)
span.SetTag("tenant.id", tenantID)

// Tracing automático para operações
err := tracing.TracedOperation(ctx, "database_save", func(ctx context.Context) error {
    return repository.Save(ctx, user)
})
```

## 🔒 Multi-tenancy

### Tenant Isolation
```go
// Middleware extrai tenant ID do header
func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantID := c.GetHeader("X-Tenant-ID")
        c.Set("tenant_id", tenantID)
        
        ctx := logging.WithTenantID(c.Request.Context(), tenantID)
        c.Request = c.Request.WithContext(ctx)
        
        c.Next()
    }
}
```

### Database Isolation
```go
// Schema por tenant
err := db.SetTenantSchema(ctx, tenantID)
```

## 🚀 Deploy

### Build Production
```bash
# Build da imagem
docker build -t direito-lux/auth-service:v1.0.0 .

# Push para registry
docker push direito-lux/auth-service:v1.0.0
```

### Kubernetes
```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: auth-service
  template:
    spec:
      containers:
      - name: auth-service
        image: direito-lux/auth-service:v1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
```

## 📚 Padrões Implementados

### Hexagonal Architecture
- **Domain**: Entidades e regras de negócio
- **Application**: Casos de uso e orchestração
- **Infrastructure**: Implementações técnicas (DB, HTTP, etc.)

### Event-Driven
- **Domain Events**: Eventos de negócio
- **Event Handlers**: Processamento assíncrono
- **Event Store**: Armazenamento de eventos (opcional)

### CQRS (Command Query Responsibility Segregation)
- **Commands**: Operações de escrita
- **Queries**: Operações de leitura
- **Projections**: Views otimizadas

## 🤝 Contribuição

### Code Style
- Seguir padrões Go (gofmt, golint)
- Comentários em português
- Testes obrigatórios (min 80% coverage)

### Git Workflow
1. Feature branch
2. Pull Request
3. Code Review
4. Merge to main

## 📄 Licença

MIT License - ver [LICENSE](LICENSE) para detalhes.