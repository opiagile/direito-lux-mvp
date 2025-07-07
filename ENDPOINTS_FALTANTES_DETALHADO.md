# 🔍 ANÁLISE DETALHADA - ENDPOINTS FALTANTES

## 📅 Data: 05/01/2025
## 🎯 Objetivo: Mapear exatamente quais APIs faltam implementar

---

## 📊 RESUMO CRÍTICO

**Descoberta Principal**: Process Service só tem endpoints de `/templates`, não de processos reais!

**Status Real**:
- ✅ **Auth Service**: 100% funcional
- ⚠️ **Tenant Service**: 10% funcional (só GET por ID)
- ❌ **Process Service**: 0% funcional para processos (só templates)
- ❌ **Todos os outros**: Não implementados

---

## 🔍 ANÁLISE POR SERVIÇO

### ✅ Auth Service (Porta 8081) - COMPLETO

**Endpoints Implementados**:
- ✅ POST `/api/v1/auth/login`
- ✅ POST `/api/v1/auth/logout`
- ✅ POST `/api/v1/auth/refresh`
- ✅ GET `/api/v1/auth/validate`
- ✅ CRUD completo `/api/v1/users/*`

**Status**: Pronto para produção

---

### ⚠️ Tenant Service (Porta 8082) - 10% COMPLETO

**Endpoints Implementados**:
- ✅ GET `/api/v1/tenants/:id`
- ✅ GET `/health`

**Endpoints Faltantes Críticos**:
```
❌ GET /api/v1/tenants/current
❌ GET /api/v1/tenants/subscription  
❌ GET /api/v1/tenants/quotas
❌ GET /api/v1/tenants (listar)
❌ POST /api/v1/tenants (criar)
❌ PUT /api/v1/tenants/:id (atualizar)
```

**Impacto**: Frontend chama esses endpoints e recebe 404

---

### ❌ Process Service (Porta 8083) - 0% FUNCIONAL

**PROBLEMA CRÍTICO**: Só implementa templates, não processos!

**Endpoints Implementados (Inúteis)**:
- ✅ GET `/api/v1/templates`
- ✅ POST `/api/v1/templates` 
- ✅ GET `/api/v1/templates/:id`
- ✅ PUT `/api/v1/templates/:id`
- ✅ DELETE `/api/v1/templates/:id`

**Endpoints Esperados pelo Frontend (TODOS FALTANDO)**:
```
❌ GET /api/v1/processes
❌ POST /api/v1/processes
❌ GET /api/v1/processes/:id
❌ PUT /api/v1/processes/:id
❌ DELETE /api/v1/processes/:id
❌ GET /api/v1/processes/:id/movements
❌ POST /api/v1/processes/:id/monitor
❌ DELETE /api/v1/processes/:id/unmonitor
❌ GET /api/v1/processes/stats  ⚠️ CRÍTICO: Dashboard espera
```

**Impacto**: Dashboard quebrado, CRUD de processos não funciona

---

### ❌ DataJud Service (Porta 8084) - TEMPLATE

**Status**: Container roda mas só tem handlers template

**Endpoints Esperados (TODOS FALTANDO)**:
```
❌ POST /api/v1/datajud/search
❌ GET /api/v1/datajud/process/:number
❌ GET /api/v1/datajud/process/:number/movements
❌ GET /api/v1/datajud/stats
```

**Impacto**: Integração CNJ não existe

---

### ❌ Notification Service (Porta 8085) - QUEBRADO

**Status**: Crash loop (`.air.toml` não encontrado)

**Endpoints Esperados (TODOS FALTANDO)**:
```
❌ GET /api/v1/notifications
❌ POST /api/v1/notifications
❌ GET /api/v1/notifications/:id
❌ PUT /api/v1/notifications/:id/read
❌ GET /api/v1/notifications/preferences
❌ GET /api/v1/notifications/templates
❌ GET /api/v1/notifications/stats
```

**Impacto**: WhatsApp, email, Telegram não funcionam

---

### ❌ Search Service (Porta 8086) - QUEBRADO

**Status**: Crash loop (dependência Fx quebrada)

**Endpoints Esperados (TODOS FALTANDO)**:
```
❌ POST /api/v1/search
❌ POST /api/v1/search/advanced
❌ GET /api/v1/search/suggestions
❌ POST /api/v1/search/aggregate
❌ POST /api/v1/index
```

**Impacto**: Busca manual ilimitada (vendida) não funciona

---

### ❌ AI Service (Porta 8087) - MUDO

**Status**: Container roda mas não responde

**Endpoints Esperados (TODOS FALTANDO)**:
```
❌ POST /api/v1/analysis/document
❌ POST /api/v1/jurisprudence/search
❌ POST /api/v1/jurisprudence/similarity
❌ POST /api/v1/generation/document
❌ GET /api/v1/analysis/history
❌ GET /api/v1/analysis/types
```

**Impacto**: Diferencial IA não existe

---

### ❌ Report Service - NÃO EXISTE NO DOCKER-COMPOSE

**Status**: Configuração ausente

**Endpoints Esperados (TODOS FALTANDO)**:
```
❌ GET /api/v1/reports
❌ POST /api/v1/reports
❌ GET /api/v1/reports/:id
❌ GET /api/v1/reports/:id/download
❌ DELETE /api/v1/reports/:id
❌ GET /api/v1/reports/stats

❌ GET /api/v1/dashboards
❌ POST /api/v1/dashboards
❌ GET /api/v1/dashboards/:id
❌ PUT /api/v1/dashboards/:id
❌ DELETE /api/v1/dashboards/:id
❌ GET /api/v1/dashboards/:id/data
❌ POST /api/v1/dashboards/:id/widgets

❌ GET /api/v1/schedules
❌ POST /api/v1/schedules
❌ GET /api/v1/schedules/:id
❌ PUT /api/v1/schedules/:id
❌ DELETE /api/v1/schedules/:id

❌ GET /api/v1/kpis
❌ POST /api/v1/kpis/calculate
```

**Impacto**: Relatórios e dashboards não funcionam

---

### ❌ MCP Service - NÃO EXISTE NO DOCKER-COMPOSE

**Status**: Configuração ausente

**Endpoints Esperados (TODOS FALTANDO)**:
```
❌ GET /api/v1/mcp/sessions
❌ POST /api/v1/mcp/sessions/:id/messages
❌ GET /api/v1/mcp/tools
❌ POST /api/v1/mcp/execute
❌ GET /api/v1/mcp/stats
```

**Impacto**: Interface conversacional não existe

---

## 🎯 ENDPOINTS CRÍTICOS POR PRIORIDADE

### 🔴 PRIORIDADE 1 - DASHBOARD FUNCIONAL (Semana 1)

#### Process Service - `/api/v1/processes/stats`
**Problema**: Dashboard quebra sem esse endpoint
**Implementação**:
```sql
-- Schema necessário
CREATE TABLE processes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    number VARCHAR(255) NOT NULL,
    court VARCHAR(255),
    subject TEXT,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

**Response esperado**:
```json
{
  "total": 150,
  "active": 120,
  "paused": 20,
  "archived": 10,
  "this_month": 25,
  "last_update": "2025-01-05T10:30:00Z"
}
```

#### Tenant Service - Endpoints faltantes
```go
// GET /api/v1/tenants/current
func (s *Server) getCurrentTenant(c *gin.Context) {
    tenantID := c.GetHeader("X-Tenant-ID")
    // Buscar tenant atual
}

// GET /api/v1/tenants/subscription  
func (s *Server) getSubscription(c *gin.Context) {
    // Retornar dados de assinatura
}

// GET /api/v1/tenants/quotas
func (s *Server) getQuotas(c *gin.Context) {
    // Retornar quotas de uso
}
```

### 🟡 PRIORIDADE 2 - CRUD BÁSICO (Semana 1-2)

#### Process Service - CRUD Completo
```go
// SUBSTITUIR endpoints /templates por /processes
func (s *Server) setupRoutes() {
    api := s.router.Group("/api/v1")
    
    // Processos (não templates!)
    processes := api.Group("/processes")
    {
        processes.GET("", handlers.ListProcesses())
        processes.POST("", handlers.CreateProcess()) 
        processes.GET("/:id", handlers.GetProcess())
        processes.PUT("/:id", handlers.UpdateProcess())
        processes.DELETE("/:id", handlers.DeleteProcess())
        processes.GET("/stats", handlers.GetStats()) // CRÍTICO
        processes.GET("/:id/movements", handlers.GetMovements())
        processes.POST("/:id/monitor", handlers.MonitorProcess())
    }
}
```

### 🟢 PRIORIDADE 3 - FUNCIONALIDADES AVANÇADAS (Semana 2-3)

#### DataJud Service
- Integração real com API CNJ
- Rate limiting (10k consultas/dia)
- Cache de resultados

#### Notification Service
- Providers WhatsApp/Email/Telegram
- Templates configuráveis
- Filas com RabbitMQ

#### Search Service
- Integração Elasticsearch
- Indexação automática
- Busca fulltext

#### AI Service
- Integração OpenAI/Claude
- Análise de documentos
- Geração de relatórios

---

## 📋 PLANO DE IMPLEMENTAÇÃO DETALHADO

### Dia 1-2: Process Service Stats
1. Criar schema `processes` no PostgreSQL
2. Implementar endpoint `GET /api/v1/processes/stats`
3. Dashboard volta a funcionar

### Dia 3-4: Tenant Service Completo
1. Implementar `/current`, `/subscription`, `/quotas`
2. Billing pages voltam a funcionar

### Dia 5-7: Process Service CRUD
1. Substituir todos endpoints `/templates` por `/processes`
2. Implementar CRUD completo
3. Páginas de processos funcionam

### Semana 2: Serviços Auxiliares
1. Corrigir Notification, Search, AI Services
2. Implementar funcionalidades básicas
3. Integração entre serviços

### Semana 3: Report e MCP Services
1. Criar configurações docker-compose
2. Implementar endpoints básicos
3. Features avançadas

---

## 🚨 ALERTAS CRÍTICOS

### 1. **Process Service é FALSO**
- Documentação afirma "completo com CQRS"
- Realidade: só tem templates inúteis
- Precisa reescrita completa

### 2. **Frontend está 70% quebrado**
- Todas as chamadas retornam 404
- Dashboard não funciona
- CRUD não funciona

### 3. **Documentação enganosa**
- Status real: 30% (não 80-90%)
- Muitos serviços são só templates
- Expectativas desalinhadas

---

## 💡 RECOMENDAÇÃO FINAL

**Foco**: Implementar 1 fluxo completo funcional ao invés de 10 serviços quebrados.

**MVP Sugerido**:
1. Auth ✅ (já funciona)
2. Tenant completo (2 dias)
3. Process básico (3 dias) 
4. Dashboard funcional (resultado)

**Resultado**: Sistema mínimo mas 100% funcional em 1 semana.

---

**Criado em**: 05/01/2025  
**Status**: Análise completa  
**Próximo**: Implementação focada