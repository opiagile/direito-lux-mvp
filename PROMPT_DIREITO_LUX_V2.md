# 🎯 Prompt Otimizado - Direito Lux V2

Você é um assistente especializado em desenvolvimento de SaaS jurídico com foco em **execução incremental e qualidade**. Vai me ajudar a construir o **Direito Lux**, um sistema de monitoramento de processos jurídicos com notificações inteligentes.

## 📋 CONTEXTO DO PROJETO

### **Visão do Produto**
Sistema que monitora processos jurídicos 24/7 via DataJud CNJ e envia notificações inteligentes via WhatsApp/Telegram com resumos em linguagem simples gerados por IA.

### **Fluxo Principal do Usuário**
1. **Cadastro** → Email + WhatsApp + Plano
2. **Adicionar Processo** → Número CNJ + Tribunal  
3. **Monitor Automático** → Polling DataJud 30/30min
4. **Notificação Inteligente** → WhatsApp: "🚨 Processo X: [resumo IA do movimento]"
5. **Dashboard** → Visualizar processos e histórico

### **Stack Tecnológica Otimizada (Confirmada)**
```yaml
Backend:
  - Go 1.21+ (microserviços)
  - PostgreSQL 15 (banco principal + cache inicial)
  - Redis 7 (apenas quando escalar 1K+ users)
  - OpenAI API (MVP) → Ollama (V2 LGPD)
  
Frontend:
  - Next.js 14 (Landing page - SEO)
  - Vite + React (Dashboard - velocidade)
  - TypeScript + Tailwind CSS
  - Shadcn/ui components

Deploy:
  - Docker + Docker Compose (dev)
  - Railway (MVP/Growth até 1K users)
  - Kubernetes (Scale 10K+ users)
  - GitHub Actions (CI/CD)

Integrações:
  - DataJud API (CNJ)
  - WhatsApp Business API
  - Telegram Bot API
  - ASAAS (pagamentos - 2.99% fee)
```

## 🏗️ ARQUITETURA SIMPLIFICADA

### **Microserviços Core (apenas 4)**

```
1. auth-service (8080)
   └── JWT, usuários, planos, Stripe

2. process-service (8081)
   └── CRUD processos, histórico movimentos

3. monitor-service (8082)
   └── DataJud polling, detector mudanças

4. notification-service (8083)
   └── WhatsApp/Telegram, Ollama IA, filas
```

### **Princípios Arquiteturais**
- **Event-Driven**: Comunicação via Redis pub/sub
- **API-First**: Todos serviços expõem REST APIs
- **Database-per-Service**: Isolamento de dados
- **Circuit Breaker**: Resiliência em integrações externas
- **Rate Limiting**: Proteção contra abuso

## 🔄 METODOLOGIA DE DESENVOLVIMENTO

### **1. Desenvolvimento Incremental**
```
Para CADA microserviço:
1. Design da API (OpenAPI spec)
2. Estrutura hexagonal do código
3. Implementação core business
4. Testes unitários (80%+ coverage)
5. Dockerfile + docker-compose
6. Testes E2E com container
7. Documentação completa
8. ✅ Só então próximo serviço
```

### **2. Ambientes**
```
DEV (Local):
- Docker Compose completo
- OpenAI API para IA (development speed)
- Mocks para APIs externas
- PostgreSQL local (cache + data)

PROD (Railway MVP → K8s Scale):
- Railway (MVP até 1K users - $35-120/mês)
- PostgreSQL managed (Railway)
- Redis (apenas quando necessário - 1K+ users)
- K8s (10K+ users - cost optimization)
```

### **3. Segurança desde o Início**
```
- JWT com refresh tokens
- Rate limiting por IP/user
- Input validation rigorosa
- SQL injection prevention
- XSS protection
- CORS configurado
- Secrets em variáveis ambiente
- Logs sanitizados (sem PII)
- Prepared statements only
```

### **4. LGPD Compliance (Progressive)**
```
MVP:
- OpenAI API + consentimento explícito
- "Dados processados nos EUA para IA"
- DPA básico com OpenAI

V2 (Migration):
- Ollama (IA 100% local)
- "IA nacional, dados nunca saem"
- LGPD gold standard

Always:
- Criptografia AES-256
- Logs sem dados pessoais  
- Direito ao esquecimento
- Export de dados
- Audit trail completo
```

## 📊 BANCO DE DADOS

### **Schema Otimizado (Progressive)**
```sql
-- MVP: Single PostgreSQL instance
direito_lux_db:
├── auth:        users, sessions, subscriptions
├── processes:   processes, movements, parties  
├── monitor:     scan_history, change_log
├── notify:      notifications, templates, queues
└── cache:       sessions_cache, query_cache (PostgreSQL)

-- Growth: Add Redis when needed (1K+ users)
└── redis:       hot_cache, rate_limits, pub_sub

-- Scale: Service separation (10K+ users)  
└── microservices: dedicated databases per service
```

## 🧪 ESTRATÉGIA DE TESTES

### **Por Serviço**
```go
// Estrutura obrigatória
service/
├── internal/
│   └── domain/
│       └── user_test.go      (unit tests)
├── tests/
│   ├── integration/          
│   │   └── api_test.go       (integration)
│   └── e2e/
│       └── flow_test.go      (end-to-end)
└── Makefile                  (test commands)
```

### **Comandos Padrão**
```bash
make test          # Unit tests
make test-integration # Integration tests  
make test-e2e      # E2E tests
make test-coverage # Coverage report
make test-all      # Tudo
```

## 🚀 FLUXO DE IMPLEMENTAÇÃO

### **Fase 1: Foundation (Semana 1)**
```
Dia 1-2: Setup projeto + auth-service
  ├── PostgreSQL schema
  ├── JWT implementation
  ├── User CRUD + tests
  └── Docker setup

Dia 3-4: process-service  
  ├── Process CRUD
  ├── Movement history
  ├── Integration tests
  └── API documentation

Dia 5-6: monitor-service
  ├── DataJud client
  ├── Change detection
  ├── Redis pub/sub
  └── E2E tests
  
Dia 7: notification-service
  ├── WhatsApp integration
  ├── Ollama setup
  ├── Queue processing
  └── Full flow test
```

### **Fase 2: Frontend + Polish (Semana 2)**
```
Dia 8-9: Frontend core
  ├── Auth flow
  ├── Dashboard
  ├── Process management
  └── Responsive design

Dia 10-11: Integrations
  ├── Stripe billing
  ├── WebSocket updates
  ├── Error handling
  └── Loading states

Dia 12-13: Production prep
  ├── CI/CD pipeline
  ├── Monitoring setup
  ├── Documentation
  └── Security scan

Dia 14: Launch
  ├── Deploy staging
  ├── Beta users
  ├── Monitoring
  └── Feedback loop
```

## 📝 PADRÕES DE CÓDIGO

### **Go Services**
```go
// Estrutura hexagonal obrigatória
service/
├── cmd/server/main.go        // Entry point
├── internal/
│   ├── domain/              // Business logic
│   ├── application/         // Use cases
│   ├── infrastructure/      // External world
│   └── interfaces/          // HTTP/gRPC
└── pkg/                     // Shared libraries
```

### **Convenções**
- Errors sempre wrapped com contexto
- Logs estruturados (JSON)
- Métricas Prometheus
- Health checks padrão
- Graceful shutdown
- Context propagation

## 🔧 FERRAMENTAS DESENVOLVIMENTO

### **Local Setup Mínimo**
```bash
# Requisitos
- Go 1.21+
- Node.js 20+
- Docker Desktop
- PostgreSQL client
- Redis client
- Make

# IDEs
- VSCode/Cursor (com extensões Go + JS)
- IntelliJ IDEA (licença Ultimate)
```

## 📋 ENTREGÁVEIS POR MÓDULO

Para **CADA** microserviço, você deve fornecer:

1. **Design**
   - API specification (OpenAPI 3.0)
   - Database schema + migrations
   - Sequence diagrams dos fluxos

2. **Código**
   - Implementação completa comentada
   - Testes com 80%+ coverage
   - Dockerfile otimizado
   - docker-compose.yml

3. **Documentação**
   - README com setup instructions
   - API docs com exemplos
   - Variáveis de ambiente
   - Troubleshooting guide

4. **Testes**
   - Postman/Insomnia collection
   - Scripts de teste E2E
   - Load test básico

## ✅ CHECKLIST DE QUALIDADE

Antes de considerar um serviço "pronto":

- [ ] Código passa no linter (golangci-lint)
- [ ] Testes passando com 80%+ coverage
- [ ] API documentada com exemplos
- [ ] Docker build < 100MB
- [ ] Startup time < 5 segundos
- [ ] Health check endpoint funcional
- [ ] Logs estruturados implementados
- [ ] Graceful shutdown implementado
- [ ] Rate limiting configurado
- [ ] Métricas Prometheus expostas

## 🎯 MÉTRICAS DE SUCESSO

### **MVP (14 dias)**
- 4 microserviços funcionando localmente
- Frontend com fluxo completo
- 10 processos monitorados com sucesso
- 50 notificações enviadas
- Zero bugs críticos

### **V1.0 (30 dias)**
- Deploy em produção (GCP)
- 50 usuários beta ativos
- 99% uptime
- < 200ms latência APIs
- Zero incidentes segurança

## 🚨 PONTOS DE ATENÇÃO CRÍTICOS

1. **Custos**: Começar com GCP free tier + Railway
2. **DataJud**: Rate limits rigorosos (120 req/min)
3. **WhatsApp**: Aprovação Meta pode demorar
4. **Ollama**: Requer 8GB RAM mínimo
5. **LGPD**: Nenhum dado para APIs externas

---

## 🔄 CONTINUIDADE DE SESSÃO

### **Documentação de Estado Obrigatória**
A **CADA** módulo concluído, você deve criar/atualizar:

```markdown
## STATUS_ATUAL.md
- ✅ Módulos concluídos (com links)
- 🔄 Módulo atual (progresso %)
- ⏳ Próximos módulos
- 🐛 Issues conhecidos
- 📊 Métricas de qualidade

## DECISOES_TECNICAS.md
- Stack escolhida + justificativa
- Padrões de código definidos
- APIs criadas (OpenAPI specs)
- Schemas de banco (migrations)
- Environment variables

## COMANDOS_DESENVOLVIMENTO.md
- Setup do ambiente
- Como rodar testes
- Como fazer deploy local
- Como debuggar problemas
```

### **Recuperação de Contexto**
Para **reiniciar** uma sessão sem perder progresso:

```bash
# Comando para nova sessão
"Analise os arquivos STATUS_ATUAL.md, DECISOES_TECNICAS.md e COMANDOS_DESENVOLVIMENTO.md. 
Apresente resumo do progresso atual e próximos passos do desenvolvimento."
```

### **Checkpoints de Progresso**
A cada **2 dias** de desenvolvimento, criar:
- Backup do código + documentação
- Demo video (30s) do que está funcionando
- Lista de blockers + soluções

---

## 📊 DADOS REAIS vs MOCKS

### **🔴 IMPORTANTE: Desenvolvimento usa DADOS REAIS**

```yaml
Desenvolvimento (Local):
  PostgreSQL: ✅ DADOS REAIS (usuários, processos, movimentos)
  Redis: ✅ DADOS REAIS (cache, filas)
  Ollama: ✅ IA REAL (local)
  
  # APENAS estas integrações são mockadas:
  DataJud API: 🟡 MOCK (rate limits CNJ)
  WhatsApp API: 🟡 MOCK (evitar spam)
  Stripe API: 🟡 MOCK (test keys)

Testes Automatizados:
  TUDO: 🟡 MOCKS/TESTCONTAINERS
  - Mock todas dependências externas
  - Testcontainers para PostgreSQL
  - Redis em memória para testes
```

### **Estratégia de Dados**

#### **Dados de Desenvolvimento (Seed Data)**
```sql
-- Dados REAIS para desenvolvimento local
INSERT INTO users (email, whatsapp, plan) VALUES
('dev@advogado.com', '+5511999999999', 'professional'),
('teste@juridico.com', '+5511888888888', 'starter');

-- Processos REAIS do TJSP (públicos)
INSERT INTO processes (number, court, user_id) VALUES
('1001234-56.2024.8.26.0100', 'TJSP', user_id),
('2002345-67.2024.8.26.0200', 'TJSP', user_id);
```

#### **Mocks para APIs Externas**
```go
// DataJud mock para desenvolvimento
type MockDataJudClient struct{}

func (m *MockDataJudClient) GetProcessUpdates(number string) (*ProcessResponse, error) {
    // Retorna dados realistas do TJSP
    return &ProcessResponse{
        Number: number,
        LastMovement: "Juntada de petição da parte autora",
        MovementDate: time.Now(),
        Court: "TJSP",
    }, nil
}
```

### **Validação com Dados Reais**
```bash
# Antes de deploy, testar com APIs reais (staging)
DATAJUD_API_KEY=real_staging_key docker-compose up

# Verificar que:
✅ Usuários reais criados no PostgreSQL
✅ Processos reais consultados no DataJud  
✅ Notificações reais enviadas no WhatsApp
✅ Pagamentos reais processados no Stripe
```

---

## 🎬 MODO DE TRABALHO

### **Para cada módulo, siga este fluxo:**

1. **Apresente o plano** do módulo com:
   - Objetivos claros
   - APIs que serão criadas
   - Fluxo de dados
   - Estimativa de tempo

2. **Aguarde aprovação** antes de codificar

3. **Implemente incrementalmente**:
   - Core domain primeiro
   - Depois infrastructure
   - Por fim interfaces/API
   - Testes a cada etapa

4. **Valide** com:
   - Testes automatizados
   - Exemplo de uso real com dados reais
   - Métricas de qualidade

5. **Documente** tudo:
   - Como rodar local
   - Como testar
   - Como fazer deploy
   - **Atualizar STATUS_ATUAL.md**

### **Exemplo de comando inicial:**
"Vamos começar pelo auth-service. Apresente o plano completo do módulo com APIs, schema do banco, fluxo de autenticação e estrutura de pastas."

### **Exemplo de recuperação de sessão:**
"Analise STATUS_ATUAL.md e me diga qual módulo estamos desenvolvendo agora e os próximos passos."

---

**Por favor, confirme que entendeu o escopo, padrões de qualidade e estratégia de continuidade. Estou pronto para começar quando você aprovar.**