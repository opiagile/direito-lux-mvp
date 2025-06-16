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

## 📊 Status Atual

- ✅ **Implementado**: 
  - Documentação completa (visão, arquitetura, roadmap)
  - Event Storming e Domain Modeling
  - Docker Compose com 15+ serviços
  - Template de microserviço Go
  - Auth Service completo (JWT, multi-tenant, CRUD)
  
- 🚧 **Em Desenvolvimento**: Tenant Service
- ⏳ **Próximos**: Process Service, DataJud Service, Notification Service

**Progresso Total**: ~25% completo

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

## 📁 Estrutura do Projeto

```
direito-lux/
├── services/               # Microserviços
│   ├── auth-service/      ✅ Implementado
│   ├── tenant-service/    🚧 Em desenvolvimento
│   ├── process-service/   ⏳ Próximo
│   └── ...
├── template-service/      ✅ Template base
├── infrastructure/        # IaC e K8s
├── scripts/              # Scripts úteis
└── docs/                 # Documentação
```

## 🛠️ Ferramentas de Desenvolvimento

- Air (Go hot reload)
- golangci-lint (Go linter)
- migrate (database migrations)
- swag (Swagger generator)
- pre-commit hooks