# Arquitetura Direito Lux - Full Cycle Approach

## Visão Geral - Event-Driven Microservices

### Bounded Contexts (DDD)
```
┌─────────────────────────────────────────────────────────────┐
│                        API Gateway                          │
│                    (Kong ou Traefik)                       │
└────────┬────────────────────────────────────┬──────────────┘
         │                                    │
    ┌────▼─────┐     ┌─────────────┐     ┌────▼─────┐
    │   MCP    │     │     BFF     │     │  BFF     │
    │ Service  │ 🤖  │   WhatsApp  │     │   Web    │
    │(Bot Hub) │     │   Telegram  │     │          │
    └────┬─────┘     └─────┬───────┘     └────┬─────┘
         │                 │                   │
┌────────▼─────────────────▼───────────────────▼──────────────┐
│                    Service Mesh (Istio)                     │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│  │   Auth      │  │  Process    │  │ Notification│       │
│  │  Service    │  │  Service    │  │   Service   │       │
│  └─────────────┘  └─────────────┘  └─────────────┘       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│  │  DataJud    │  │     AI      │  │   Search    │       │
│  │  Service    │  │  Service    │  │   Service   │       │
│  └─────────────┘  └─────────────┘  └─────────────┘       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│  │   Tenant    │  │  Analytics  │  │   Report    │       │
│  │  Service    │  │  Service    │  │   Service   │       │
│  └─────────────┘  └─────────────┘  └─────────────┘       │
└─────────────────────────────────────────────────────────────┘
                              │
                    ┌─────────▼──────────┐
                    │   Message Broker   │
                    │  RabbitMQ/Kafka    │
                    └────────────────────┘
```

## Microserviços Detalhados

### 1. Auth Service (Golang + Keycloak)
```go
// Hexagonal Architecture
├── domain/
│   ├── entities/
│   │   ├── user.go
│   │   └── tenant.go
│   └── ports/
│       ├── auth_repository.go
│       └── token_service.go
├── application/
│   ├── usecases/
│   │   ├── authenticate_user.go
│   │   └── validate_tenant.go
│   └── dto/
├── infrastructure/
│   ├── keycloak/
│   ├── database/
│   └── messaging/
└── interfaces/
    ├── grpc/
    └── http/
```

**Responsabilidades:**
- Autenticação multi-tenant
- Autorização baseada em roles (RBAC)
- JWT token management
- Session management
- Integração com Keycloak

### 2. Process Service (Golang)
```go
// CQRS Pattern
├── domain/
│   ├── aggregates/
│   │   └── process_aggregate.go
│   ├── events/
│   │   ├── process_created.go
│   │   └── process_updated.go
│   └── commands/
│       └── monitor_process.go
├── application/
│   ├── commands/
│   ├── queries/
│   └── projections/
└── infrastructure/
    ├── eventstore/
    └── read_model/
```

**Responsabilidades:**
- Gestão do ciclo de vida dos processos
- Event sourcing para auditoria completa
- CQRS para separar leitura/escrita
- Cache distribuído (Redis)

### 3. DataJud Service (Golang)
```go
// Circuit Breaker + Rate Limiting
├── domain/
│   ├── services/
│   │   └── datajud_client.go
│   └── value_objects/
│       └── process_data.go
├── application/
│   ├── cache_strategy/
│   ├── rate_limiter/
│   └── circuit_breaker/
└── infrastructure/
    ├── http_client/
    ├── redis_cache/
    └── metrics/
```

**Responsabilidades:**
- Integração com API DataJud
- Rate limiting (10k/dia)
- Cache inteligente por tenant
- Circuit breaker para resiliência
- Retry com backoff exponencial

### 4. Notification Service (Golang)
```go
// Strategy Pattern para multi-canal
├── domain/
│   ├── interfaces/
│   │   └── notification_channel.go
│   └── strategies/
│       ├── whatsapp_strategy.go
│       ├── email_strategy.go
│       └── telegram_strategy.go
├── application/
│   ├── notification_orchestrator.go
│   └── template_engine/
└── infrastructure/
    ├── whatsapp_business/
    ├── smtp/
    └── telegram_bot/
```

**Responsabilidades:**
- Orquestração de notificações multicanal
- Template engine para mensagens
- Retry mechanism
- Delivery tracking
- Webhooks para status

### 5. AI Service (Python)
```python
# Clean Architecture
├── domain/
│   ├── entities/
│   │   └── legal_document.py
│   └── use_cases/
│       ├── summarize_process.py
│       └── explain_terms.py
├── application/
│   ├── nlp_pipeline/
│   └── ml_models/
├── infrastructure/
│   ├── openai_adapter/
│   ├── huggingface_adapter/
│   └── model_cache/
└── interfaces/
    ├── grpc/
    └── http/
```

**Responsabilidades:**
- Processamento de linguagem natural
- Resumos inteligentes
- Explicação de termos jurídicos
- Análise de sentimento
- Classificação de documentos

### 6. Document Service (Golang)
```go
// Repository Pattern
├── domain/
│   ├── entities/
│   │   └── document.go
│   └── repositories/
│       └── document_repository.go
├── application/
│   ├── document_generator/
│   └── template_manager/
└── infrastructure/
    ├── s3_storage/
    ├── pdf_generator/
    └── ocr_service/
```

**Responsabilidades:**
- Geração de documentos jurídicos
- Armazenamento seguro (S3)
- OCR para digitalização
- Versionamento de documentos
- Templates customizáveis

### 7. Tenant Service (Golang)
```go
// Multi-tenant isolation
├── domain/
│   ├── entities/
│   │   └── tenant.go
│   └── policies/
│       └── data_isolation.go
├── application/
│   ├── tenant_provisioning/
│   └── quota_management/
└── infrastructure/
    ├── database_sharding/
    └── tenant_middleware/
```

**Responsabilidades:**
- Provisionamento de tenants
- Isolamento de dados
- Gestão de quotas
- Billing integration
- Feature flags por tenant

### 8. Analytics Service (Golang)
```go
// Event-driven analytics
├── domain/
│   ├── aggregations/
│   └── metrics/
├── application/
│   ├── event_processors/
│   └── report_generators/
└── infrastructure/
    ├── clickhouse/
    ├── prometheus/
    └── grafana/
```

**Responsabilidades:**
- Métricas em tempo real
- Dashboards customizados
- Jurimetria (ML predictions)
- Reports por tenant
- Data warehouse

### 9. MCP Service (Golang + Claude API) 🤖 DIFERENCIAL EXCLUSIVO
```go
// Model Context Protocol Architecture
├── domain/
│   ├── entities/
│   │   ├── mcp_session.go
│   │   ├── conversation_context.go
│   │   └── tool_registry.go
│   └── tools/
│       ├── process_tools.go      // 5 ferramentas
│       ├── ai_tools.go           // 4 ferramentas
│       ├── search_tools.go       // 2 ferramentas
│       ├── notification_tools.go // 2 ferramentas
│       ├── report_tools.go       // 2 ferramentas
│       └── admin_tools.go        // 2 ferramentas
├── application/
│   ├── mcp_orchestrator/
│   ├── context_manager/
│   ├── tool_executor/
│   └── bot_interfaces/
│       ├── whatsapp_handler.go
│       ├── telegram_handler.go
│       ├── claude_chat_handler.go
│       └── slack_handler.go
└── infrastructure/
    ├── claude_api/
    ├── webhook_adapters/
    ├── session_cache/
    └── quota_manager/
```

**Responsabilidades:**
- **Interface Conversacional**: WhatsApp, Telegram, Claude Chat, Slack
- **17+ Ferramentas MCP**: Comandos naturais para todas as funcionalidades
- **Context Management**: Sessões conversacionais com memória
- **Multi-tenant Security**: Isolamento completo entre escritórios
- **Quota Management**: 200/1000/ilimitado comandos por plano
- **Tool Registry**: Sistema dinâmico de registro de ferramentas
- **Bot Orchestration**: Coordenação entre múltiplas interfaces

**Diferencial Estratégico:**
- Primeiro SaaS jurídico brasileiro com interface conversacional completa
- Democratização do acesso via linguagem natural
- Redução da curva de aprendizado para advogados
- Automação de tarefas via comandos de voz/texto

## Infraestrutura e DevOps - IMPLEMENTADO COMPLETO!

### 🏗️ Infrastructure as Code (Terraform)
```hcl
# terraform/main.tf
module "networking" {
  source = "./modules/networking"
  
  project_id    = var.project_id
  environment   = var.environment
  vpc_cidr      = var.vpc_cidr
  subnet_cidrs  = var.subnet_cidrs
}

module "gke" {
  source = "./modules/gke"
  
  project_id     = var.project_id
  environment    = var.environment
  network_id     = module.networking.network_id
  subnet_id      = module.networking.private_subnet_id
  node_pools     = var.gke_node_pools
}

module "database" {
  source = "./modules/database"
  
  project_id       = var.project_id
  environment      = var.environment
  network_id       = module.networking.network_id
  database_config  = var.database_config
  redis_config     = var.redis_config
}
```

**Recursos Provisionados:**
- VPC com subnets segmentadas (public, private, database, GKE)
- GKE cluster regional com private nodes e multiple node pools
- Cloud SQL PostgreSQL com HA e read replicas
- Redis com Standard HA tier
- Load Balancer global com SSL termination
- Cloud DNS com certificados gerenciados
- IAM roles e service accounts
- Network policies e firewall rules

### ☸️ Kubernetes Production (Implementado)
```yaml
# k8s/production/services/process-service.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: process-service
  namespace: direito-lux-production
  labels:
    app: process-service
    tier: backend
spec:
  replicas: 5
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 2
      maxUnavailable: 1
  selector:
    matchLabels:
      app: process-service
  template:
    metadata:
      labels:
        app: process-service
        tier: backend
    spec:
      serviceAccountName: process-service-sa
      containers:
      - name: process-service
        image: gcr.io/direito-lux-production/process-service:latest
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: postgres-secret
              key: url
        - name: REDIS_URL
          valueFrom:
            secretKeyRef:
              name: redis-secret
              key: url
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: process-service-hpa
  namespace: direito-lux-production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: process-service
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

**Estrutura Kubernetes Completa:**
```
k8s/
├── deploy.sh                    # Script de deploy automatizado
├── staging/                     # Ambiente staging
│   ├── databases/              # PostgreSQL, Redis
│   ├── services/               # Todos os microserviços
│   ├── ingress/                # Ingress controllers e SSL
│   └── monitoring/             # Prometheus, Grafana, Jaeger
├── production/                  # Ambiente production
│   ├── databases/              # HA databases
│   ├── services/               # Scaled microservices
│   ├── ingress/                # Production ingress
│   └── monitoring/             # Production monitoring
└── shared/                      # Resources compartilhados
    ├── namespaces.yaml
    ├── network-policies.yaml
    └── rbac.yaml
```

### 🔄 CI/CD Pipeline (GitHub Actions - Implementado)
```yaml
# .github/workflows/ci-cd.yml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  security-scan:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Run Trivy vulnerability scanner
      uses: aquasecurity/trivy-action@master
    - name: SAST with CodeQL
      uses: github/codeql-action/init@v2

  build-and-test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        service: [auth-service, process-service, tenant-service, 
                 datajud-service, notification-service, ai-service,
                 search-service, mcp-service, report-service]
    steps:
    - uses: actions/checkout@v4
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Build ${{ matrix.service }}
      run: |
        cd services/${{ matrix.service }}
        go build -v ./...
        go test -v ./...
        go test -race -coverprofile=coverage.out ./...
    
    - name: Build Docker image
      run: |
        cd services/${{ matrix.service }}
        docker build -t gcr.io/${{ secrets.GCP_PROJECT }}/${{ matrix.service }}:${{ github.sha }} .
    
    - name: Push to GCR
      if: github.ref == 'refs/heads/main' || github.ref == 'refs/heads/develop'
      run: |
        echo ${{ secrets.GCP_SA_KEY }} | base64 -d | docker login -u _json_key --password-stdin https://gcr.io
        docker push gcr.io/${{ secrets.GCP_PROJECT }}/${{ matrix.service }}:${{ github.sha }}

  deploy-staging:
    if: github.ref == 'refs/heads/develop'
    needs: [build-and-test, security-scan]
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Setup kubectl
      uses: azure/k8s-set-context@v1
      with:
        method: service-account
        k8s-url: ${{ secrets.GKE_CLUSTER_URL }}
        k8s-secret: ${{ secrets.GKE_SA_KEY }}
    
    - name: Deploy to staging
      run: |
        cd k8s
        ./deploy.sh staging --apply --image-tag=${{ github.sha }}

  deploy-production:
    if: github.ref == 'refs/heads/main'
    needs: [build-and-test, security-scan]
    runs-on: ubuntu-latest
    environment: production
    steps:
    - uses: actions/checkout@v4
    - name: Deploy to production
      run: |
        cd k8s
        ./deploy.sh production --apply --image-tag=${{ github.sha }}

  performance-tests:
    needs: deploy-staging
    runs-on: ubuntu-latest
    steps:
    - name: Run k6 performance tests
      run: |
        docker run --rm -v $PWD/tests:/tests grafana/k6 run /tests/load-test.js
```

### Container Registry e CI/CD (Implementado)

### Message Broker (RabbitMQ/Kafka)
```yaml
# Topics/Exchanges
- process.created
- process.updated
- notification.send
- document.generated
- tenant.provisioned
- billing.processed
```

### Observabilidade Stack
- **Logs**: ELK Stack (Elasticsearch, Logstash, Kibana)
- **Metrics**: Prometheus + Grafana
- **Tracing**: Jaeger (OpenTelemetry)
- **APM**: New Relic ou DataDog

### Banco de Dados
- **PostgreSQL**: Sharding por tenant_id
- **Redis**: Cache distribuído
- **MongoDB**: Documentos e templates
- **ClickHouse**: Analytics e time-series

### Segurança
- **Vault**: Gestão de secrets
- **Cert-Manager**: SSL/TLS automático
- **Network Policies**: Isolamento entre services
- **RBAC**: Kubernetes native

## Padrões e Boas Práticas

### 1. Event-Driven Patterns
```go
// Exemplo de evento
type ProcessUpdatedEvent struct {
    EventID     string    `json:"event_id"`
    TenantID    string    `json:"tenant_id"`
    ProcessID   string    `json:"process_id"`
    UpdateType  string    `json:"update_type"`
    OccurredAt  time.Time `json:"occurred_at"`
    Payload     json.RawMessage `json:"payload"`
}
```

### 2. Saga Pattern para transações distribuídas
```go
// Orquestrador de Saga
type ProcessMonitoringSaga struct {
    // 1. Validate tenant quota
    // 2. Register process
    // 3. Schedule monitoring
    // 4. Send confirmation
    // Compensating transactions if needed
}
```

### 3. Idempotência
- Todos endpoints devem ser idempotentes
- Use Idempotency Keys
- Event deduplication

### 4. Rate Limiting por Tenant
```go
type RateLimiter struct {
    limiter *rate.Limiter
    tenantQuotas map[string]int
}
```

## Escalabilidade

### Horizontal Pod Autoscaling
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: process-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: process-service
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### Database Read Replicas
- Master para escrita
- Multiple read replicas
- Connection pooling
- Query routing

### CDN e Edge Computing
- CloudFlare para assets
- Edge functions para validações
- Regional deployments

## Disaster Recovery

### Backup Strategy
- Database: Daily snapshots
- Event Store: Continuous replication
- Documents: S3 cross-region replication

### Multi-region Setup
- Active-passive para começar
- Active-active para Enterprise

## Migração e Rollout

### Feature Flags (Unleash)
```go
if featureToggle.IsEnabled("new-ai-model", tenant.ID) {
    // Use new AI model
} else {
    // Use legacy model
}
```

### Blue-Green Deployment
- Zero downtime deployments
- Instant rollback capability
- A/B testing infrastructure

## Métricas de Sucesso

### SLIs (Service Level Indicators)
- Latência p99 < 200ms
- Disponibilidade > 99.9%
- Taxa de erro < 0.1%

### SLOs (Service Level Objectives)
- Process query time < 100ms
- Notification delivery < 5s
- AI response time < 2s

### SLAs (Service Level Agreements)
- Starter: 99.5% uptime
- Professional: 99.9% uptime
- Enterprise: 99.99% uptime

## Custo Estimado GCP (São Paulo - southamerica-east1)

### Ambiente Desenvolvimento
- GKE Autopilot: ~$74/mês
- Cloud SQL (1vCPU, 4GB): ~$70/mês
- Memorystore Redis (1GB): ~$45/mês
- Pub/Sub: ~$0 (free tier)
- Storage: ~$20/mês
- **Total Dev**: ~$209/mês

### Ambiente Produção
- GKE Autopilot: ~$74 + ~$200 (compute)
- Cloud SQL HA (2vCPU, 8GB): ~$140/mês
- Memorystore Redis HA (5GB): ~$90/mês
- Pub/Sub: ~$40/mês
- Load Balancer + CDN: ~$50/mês
- Cloud Storage: ~$40/mês
- Monitoring/Logging: ~$100/mês
- **Total Prod**: ~$734/mês inicial

### Scaling Costs
- Por 1000 usuários ativos: +$300/mês
- Storage por TB: +$20/mês
- CDN Bandwidth por TB: +$85/mês
- Pub/Sub por milhão msgs: +$40/mês

## Status Atual Full Cycle - IMPLEMENTADO!

### ✅ Concluído (100%)
1. **Domain Modeling** - Event Storming completo com 7 bounded contexts
2. **API First** - Contratos definidos para todos os serviços
3. **Kubernetes Production** - Manifests completos para staging e production
4. **Microserviços Core** - Todos os 10 serviços implementados
5. **CI/CD Pipeline** - GitHub Actions com build, test, security e deploy
6. **Observability** - Prometheus, Grafana, Jaeger implementados
7. **Infrastructure as Code** - Terraform completo para GCP
8. **Frontend Web App** - Next.js 14 com todas as funcionalidades

### ⏳ Próximos Passos (Finalizando)
1. **Testes de Integração E2E** - Validação de fluxos completos
2. **Mobile App** - React Native para iOS e Android
3. **Load Testing** - Performance testing com k6 em produção
4. **Documentação API** - OpenAPI/Swagger para todos os serviços

### 🏆 Arquitetura Final Alcançada
- ✅ **Cloud-Native**: Kubernetes + GCP + Terraform
- ✅ **Event-Driven**: RabbitMQ + Domain Events
- ✅ **Multi-tenant**: Isolamento completo por tenant
- ✅ **Observability**: Full stack monitoring
- ✅ **Security**: RBAC + Network Policies + SSL
- ✅ **Scalability**: HPA + Cluster Autoscaler
- ✅ **Disaster Recovery**: HA databases + backups
- ✅ **DevOps**: GitOps + Infrastructure as Code