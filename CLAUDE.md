# Contexto para Claude - Projeto Direito Lux

## 🎯 Sobre o Projeto

O Direito Lux é uma plataforma SaaS para monitoramento automatizado de processos jurídicos, integrada com a API DataJud do CNJ, oferecendo notificações multicanal e análise com IA.

## 🏗️ Arquitetura

- **Microserviços** em Go (Hexagonal Architecture)
- **Event-Driven** com RabbitMQ
- **Multi-tenant** com isolamento por schema PostgreSQL
- **Cloud-native** para GCP com Kubernetes
- **AI Service** em Python para análises

## 📋 Processo de Desenvolvimento

### 🔄 Ao Finalizar Cada Módulo/Serviço

**IMPORTANTE**: Sempre atualizar a documentação após implementar qualquer componente!

1. **STATUS_IMPLEMENTACAO.md**
   - Mover item de "O que Falta" para "O que está Implementado"
   - Atualizar percentual de progresso
   - Adicionar detalhes do que foi implementado

2. **README.md**
   - Atualizar seção "Status do Projeto"
   - Adicionar URLs de desenvolvimento
   - Atualizar comandos úteis

3. **SETUP_AMBIENTE.md**
   - Adicionar instruções de setup do novo módulo
   - Incluir novas variáveis de ambiente
   - Documentar troubleshooting

4. **Documentação do Módulo**
   - Criar README.md específico no diretório do serviço
   - Documentar APIs e eventos
   - Incluir exemplos de uso

### 📝 Padrões de Código

- **Go**: Seguir template em `template-service/`
- **Comentários**: Sempre em português
- **Commits**: Conventional Commits
- **Testes**: Mínimo 80% coverage
- **APIs**: Documentar com Swagger/OpenAPI
- **Snippets**: Máximo 40 linhas por vez

### 🚀 Comandos Importantes

```bash
# Criar novo serviço
./scripts/create-service.sh nome-service

# Rodar migrações
cd services/[nome-service]
migrate -path migrations -database "postgres://..." up

# Executar testes
make test
make test-coverage
```

## 📊 Status Atual (Atualizado 07/01/2025)

- ✅ **Implementado (85% do projeto)**: 
  - Documentação completa (visão, arquitetura, roadmap)
  - Event Storming e Domain Modeling
  - Docker Compose com 15+ serviços
  - Template de microserviço Go
  - **10 Microserviços Core 100% funcionais**: Auth, Tenant, Process, DataJud, Notification, AI, Search, MCP, Report
  - **Frontend Next.js 14 completo** - CRUD processos, busca, billing, dashboard
  - **Infrastructure completa**: K8s, Terraform, CI/CD GitHub Actions
  
- ⚠️ **Auditoria Externa Concluída (07/01/2025)**: 
  - ✅ Todas configurações de APIs externas verificadas
  - ❌ **DataJud Service identificado como MOCK** - precisa implementação real
  - ⚠️ Todas as chaves configuradas para DEV (demo/mock tokens)
  - ✅ Ambiente funcional para desenvolvimento e testes de arquitetura
  
- 🎯 **Próximo Marco Crítico: AMBIENTE STAGING**
  - ❌ **DataJud HTTP Client real** - substituir mock por implementação CNJ
  - ⚠️ **APIs reais com quotas limitadas** - OpenAI, WhatsApp, Telegram, CNJ  
  - ⚠️ **Certificado digital A1/A3** para autenticação CNJ obrigatória
  - ⚠️ **Webhooks HTTPS** para WhatsApp e Telegram
  - ✅ **Validação E2E com dados reais** antes da produção

**Progresso Total**: ~85% completo (desenvolvimento), próximo: STAGING (2-3 dias)

## 🔗 Documentação Principal

Consultar sempre:
- [PROCESSO_DOCUMENTACAO.md](./PROCESSO_DOCUMENTACAO.md) - Como manter docs atualizadas
- [STATUS_IMPLEMENTACAO.md](./STATUS_IMPLEMENTACAO.md) - Status detalhado
- [ARQUITETURA_FULLCYCLE.md](./ARQUITETURA_FULLCYCLE.md) - Arquitetura técnica

## ⚠️ Lembretes Importantes

1. **Sempre atualizar documentação ao finalizar implementações**
2. **Usar Event-Driven Architecture para comunicação entre serviços**
3. **Implementar health checks e métricas em todos os serviços**
4. **Seguir padrão de multi-tenancy com header X-Tenant-ID**
5. **Todos os serviços devem ter Dockerfile e docker-compose entry**

## 🚨 LIÇÕES APRENDIDAS - AUDITORIA EXTERNA (07/01/2025)

### ⚠️ **CONFIGURAÇÕES DEV ≠ PROD**

**❌ Riscos Identificados:**
- **DataJud Service tem implementação MOCK** - não funciona em produção
- **APIs externas usam tokens demo** - WhatsApp, Telegram, OpenAI
- **Ambiente DEV não garante funcionamento em PROD**

### 🔧 **PREPARAÇÃO PARA STAGING**

**Configurações obrigatórias para ambiente staging:**

```bash
# Chaves reais (desenvolvimento limitado)
OPENAI_API_KEY=sk-real-but-limited-key
DATAJUD_API_KEY=real_cnj_staging_key
DATAJUD_CERTIFICATE_PATH=/certs/staging.p12
DATAJUD_CERTIFICATE_PASSWORD=staging_cert_password
WHATSAPP_ACCESS_TOKEN=staging_meta_token
TELEGRAM_BOT_TOKEN=staging_bot_token
ANTHROPIC_API_KEY=sk-ant-staging-key

# URLs públicas obrigatórias
WHATSAPP_WEBHOOK_URL=https://staging.direitolux.com.br/webhook/whatsapp
TELEGRAM_WEBHOOK_URL=https://staging.direitolux.com.br/webhook/telegram
```

### 📋 **PROCESSO STAGING**

1. **Implementar DataJud HTTP Client real** (substitui mock)
2. **Configurar certificado digital CNJ**
3. **Criar webhooks HTTPS públicos**
4. **Configurar APIs reais com quotas limitadas**
5. **Testes E2E com dados reais**
6. **Validação completa antes de produção**

### 🎯 **PRÓXIMAS SESSÕES**

- **Prioridade 1**: Implementar DataJud HTTP Client real
- **Prioridade 2**: Configurar ambiente staging com APIs reais
- **Prioridade 3**: Testes de integração E2E com dados reais

## 🎯 Diferenciais do Produto

- WhatsApp em TODOS os planos (diferencial competitivo)
- Busca manual ilimitada em todos os planos
- Integração com DataJud (limite 10k consultas/dia)
- IA para resumos adaptados (advogados e clientes)
- Multi-tenant com isolamento completo

## 💰 Planos de Assinatura

- **Starter**: R$99 (50 processos, 20 clientes, 100 consultas/dia)
- **Professional**: R$299 (200 processos, 100 clientes, 500 consultas/dia)
- **Business**: R$699 (500 processos, 500 clientes, 2000 consultas/dia)
- **Enterprise**: R$1999+ (ilimitado, 10k consultas/dia, white-label)

## 🏛️ Bounded Contexts

1. **Authentication & Identity** - Keycloak, JWT, RBAC
2. **Tenant Management** - Planos, quotas, billing
3. **Process Management** - Core domain, CQRS
4. **External Integration** - DataJud API, circuit breaker
5. **Notification System** - WhatsApp, Email, Telegram
6. **AI & Analytics** - Resumos, jurimetria
7. **Document Management** - Templates, assinaturas

## 🔧 Stack Tecnológica

- **Backend**: Go 1.21+ (microserviços)
- **AI/ML**: Python 3.11+ (FastAPI)
- **Frontend**: Next.js 14 + TypeScript
- **Mobile**: React Native + Expo
- **Database**: PostgreSQL 15 + Redis
- **Message Queue**: RabbitMQ
- **Cloud**: Google Cloud Platform
- **Container**: Docker + Kubernetes (GKE)
- **IaC**: Terraform
- **CI/CD**: GitHub Actions + ArgoCD
- **Observability**: Jaeger + Prometheus + Grafana

## 📁 Estrutura do Projeto (Atualizada)

```
direito-lux/
├── services/               # Microserviços (100% Implementados)
│   ├── auth-service/      ✅ Funcional (JWT, multi-tenant)
│   ├── tenant-service/    ✅ Funcional (planos, quotas)
│   ├── process-service/   ✅ Funcional (CQRS, CRUD)
│   ├── datajud-service/   ⚠️ Mock (precisa HTTP client real)
│   ├── notification-service/ ✅ Funcional (WhatsApp, email)
│   ├── ai-service/        ✅ Funcional (Python/FastAPI)
│   ├── search-service/    ✅ Funcional (Elasticsearch)
│   ├── mcp-service/       ✅ Funcional (Claude MCP)
│   └── report-service/    ✅ Funcional (dashboard, PDF)
├── template-service/      ✅ Template base Go
├── frontend/              ✅ Next.js 14 completo (CRUD, busca)
├── infrastructure/        ✅ K8s + Terraform completos
├── scripts/              ✅ Deploy e utilities
├── docs/                 ✅ Documentação completa
└── .github/workflows/    ✅ CI/CD GitHub Actions
```

## 🛠️ Ferramentas de Desenvolvimento

- Air (Go hot reload)
- golangci-lint (Go linter)
- migrate (database migrations)
- swag (Swagger generator)
- pre-commit hooks