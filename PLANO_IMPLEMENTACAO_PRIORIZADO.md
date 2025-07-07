# 🚀 PLANO DE IMPLEMENTAÇÃO PRIORIZADO - DIREITO LUX

## 📅 Data: 05/01/2025
## 🎯 Meta: MVP Funcional em 3-4 semanas

---

## 📊 STATUS ATUAL REAL

**Descoberta Crítica**: Sistema tem apenas 30% de funcionalidade real (não 80-90% como documentado)

**Serviços Funcionais**:
- ✅ Auth Service: 100% completo
- ⚠️ Tenant Service: 10% completo (só GET por ID)
- ❌ Process Service: 0% funcional (só templates inúteis)
- ❌ Outros 7 serviços: Não implementados ou quebrados

---

## 🎯 ESTRATÉGIA: VERTICAL SLICE FUNCIONAL

### Princípio: Menos features, mais qualidade
- **Ao invés de**: 10 serviços 50% quebrados
- **Focar em**: 4 serviços 100% funcionais
- **Resultado**: Um fluxo completo que funciona

---

## 🔴 SEMANA 1: CORREÇÕES CRÍTICAS (7 dias)

### **Dia 1-2: Process Service - Endpoint Stats**
**Prioridade**: CRÍTICA (Dashboard quebrado)

**Tarefas**:
1. **Criar schema PostgreSQL**:
```sql
CREATE TABLE processes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(id),
    number VARCHAR(255) NOT NULL UNIQUE,
    court VARCHAR(255),
    subject TEXT,
    status VARCHAR(50) DEFAULT 'active',
    monitoring BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Índices para performance
CREATE INDEX idx_processes_tenant_id ON processes(tenant_id);
CREATE INDEX idx_processes_status ON processes(status);
CREATE INDEX idx_processes_number ON processes(number);
```

2. **Implementar endpoint stats**:
```go
// GET /api/v1/processes/stats
func (h *ProcessHandler) GetStats(c *gin.Context) {
    tenantID := c.GetHeader("X-Tenant-ID")
    
    var stats ProcessStats
    err := h.db.Get(&stats, `
        SELECT 
            COUNT(*) as total,
            COUNT(*) FILTER (WHERE status = 'active') as active,
            COUNT(*) FILTER (WHERE status = 'paused') as paused,
            COUNT(*) FILTER (WHERE status = 'archived') as archived,
            COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '30 days') as this_month
        FROM processes 
        WHERE tenant_id = $1
    `, tenantID)
    
    if err != nil {
        c.JSON(500, gin.H{"error": "Erro ao buscar estatísticas"})
        return
    }
    
    c.JSON(200, stats)
}
```

3. **Adicionar dados de teste**:
```sql
-- Inserir processos de teste para cada tenant
INSERT INTO processes (tenant_id, number, court, subject, status) VALUES
('11111111-1111-1111-1111-111111111111', '5001234-12.2024.8.26.0100', 'TJSP', 'Ação de Cobrança', 'active'),
('11111111-1111-1111-1111-111111111111', '5001235-12.2024.8.26.0100', 'TJSP', 'Ação Trabalhista', 'active'),
-- ... mais dados para todos os tenants
```

**Resultado**: Dashboard volta a funcionar com dados reais

### **Dia 3-4: Tenant Service Completo**
**Prioridade**: ALTA (Billing quebrado)

**Endpoints a implementar**:
```go
// GET /api/v1/tenants/current
func (s *Server) getCurrentTenant(c *gin.Context) {
    tenantID := c.GetHeader("X-Tenant-ID")
    
    var tenant TenantDB
    err := s.db.Get(&tenant, `
        SELECT id, legal_name, COALESCE(name, legal_name) as name, 
               email, plan_type, status, created_at, updated_at 
        FROM tenants WHERE id = $1
    `, tenantID)
    
    if err != nil {
        c.JSON(404, gin.H{"error": "Tenant não encontrado"})
        return
    }
    
    c.JSON(200, tenant)
}

// GET /api/v1/tenants/subscription
func (s *Server) getSubscription(c *gin.Context) {
    tenantID := c.GetHeader("X-Tenant-ID")
    
    // Buscar dados de assinatura
    subscription := SubscriptionInfo{
        PlanType: "professional",
        Status: "active",
        ExpiresAt: "2025-12-31T23:59:59Z",
        Features: []string{"whatsapp", "email", "unlimited_search"},
        Limits: map[string]int{
            "processes": 200,
            "clients": 100,
            "daily_queries": 500,
        },
    }
    
    c.JSON(200, subscription)
}

// GET /api/v1/tenants/quotas
func (s *Server) getQuotas(c *gin.Context) {
    tenantID := c.GetHeader("X-Tenant-ID")
    
    // Calcular quotas atuais vs. limites
    quotas := QuotaInfo{
        Processes: QuotaDetail{Used: 45, Limit: 200, Percentage: 22.5},
        Clients: QuotaDetail{Used: 12, Limit: 100, Percentage: 12.0},
        DailyQueries: QuotaDetail{Used: 127, Limit: 500, Percentage: 25.4},
        MonthlyQueries: QuotaDetail{Used: 2456, Limit: 15000, Percentage: 16.4},
    }
    
    c.JSON(200, quotas)
}
```

**Resultado**: Billing pages voltam a funcionar

### **Dia 5-6: Process Service CRUD Básico**
**Prioridade**: ALTA (CRUD quebrado)

**Substituir endpoints `/templates` por `/processes`**:
```go
func (s *Server) setupRoutes() {
    api := s.router.Group("/api/v1")
    
    // REMOVER: templates (inúteis)
    // ADICIONAR: processes (funcionais)
    processes := api.Group("/processes")
    {
        processes.GET("", handlers.ListProcesses())       // Listar
        processes.POST("", handlers.CreateProcess())      // Criar
        processes.GET("/:id", handlers.GetProcess())      // Buscar
        processes.PUT("/:id", handlers.UpdateProcess())   // Atualizar
        processes.DELETE("/:id", handlers.DeleteProcess()) // Excluir
        processes.GET("/stats", handlers.GetStats())      // Estatísticas
        processes.GET("/:id/movements", handlers.GetMovements()) // Movimentações
        processes.POST("/:id/monitor", handlers.MonitorProcess()) // Monitorar
    }
}
```

**Handlers básicos**:
```go
func (h *ProcessHandler) ListProcesses(c *gin.Context) {
    tenantID := c.GetHeader("X-Tenant-ID")
    
    var processes []Process
    err := h.db.Select(&processes, `
        SELECT id, number, court, subject, status, monitoring, created_at, updated_at
        FROM processes 
        WHERE tenant_id = $1 
        ORDER BY created_at DESC
    `, tenantID)
    
    if err != nil {
        c.JSON(500, gin.H{"error": "Erro ao listar processos"})
        return
    }
    
    c.JSON(200, gin.H{"data": processes, "total": len(processes)})
}

func (h *ProcessHandler) CreateProcess(c *gin.Context) {
    tenantID := c.GetHeader("X-Tenant-ID")
    
    var req CreateProcessRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Dados inválidos"})
        return
    }
    
    processID := uuid.New()
    _, err := h.db.Exec(`
        INSERT INTO processes (id, tenant_id, number, court, subject, status)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, processID, tenantID, req.Number, req.Court, req.Subject, "active")
    
    if err != nil {
        c.JSON(500, gin.H{"error": "Erro ao criar processo"})
        return
    }
    
    c.JSON(201, gin.H{"id": processID, "message": "Processo criado com sucesso"})
}
```

**Resultado**: CRUD de processos funciona completamente

### **Dia 7: Testes e Integração**
**Prioridade**: CRÍTICA (Validação)

**Testes E2E**:
1. Login → Dashboard (estatísticas aparecem)
2. Login → Processos → Criar/Editar/Listar
3. Login → Billing (dados de assinatura aparecem)
4. Todos os 8 tenants funcionando

**Resultado**: MVP básico 100% funcional

---

## 🟡 SEMANA 2: SERVIÇOS AUXILIARES (7 dias)

### **Dia 8-10: Corrigir Serviços Quebrados**

#### Notification Service
**Problema**: `.air.toml` faltando
**Solução**:
```bash
# Remover Air e usar go run direto
cd services/notification-service
# Criar dockerfile simples
docker build -t notification-service .
# Atualizar docker-compose.yml
```

#### Search Service  
**Problema**: Dependência Fx quebrada
**Solução**:
```go
// Remover Fx, usar conexão direta como tenant-service
func main() {
    // Setup simples como auth-service
    router := gin.New()
    router.GET("/health", healthCheck)
    router.POST("/api/v1/search", basicSearch)
    router.Listen(":8086")
}
```

#### AI Service
**Problema**: Python não responde
**Solução**:
```python
# Verificar se FastAPI está rodando corretamente
# Debug logs, health check
from fastapi import FastAPI
app = FastAPI()

@app.get("/health")
def health():
    return {"status": "healthy"}

@app.post("/api/v1/analysis/document")
def analyze(data: dict):
    return {"summary": "Análise em desenvolvimento"}
```

### **Dia 11-14: Implementações Básicas**

#### DataJud Service
- Endpoints básicos com dados mockados
- Rate limiting simulado
- Cache em memória

#### Search Service
- Busca simples no PostgreSQL
- Sem Elasticsearch (por enquanto)
- Busca por número, tribunal, assunto

#### Notification Service
- Templates fixos
- Email via SMTP
- WhatsApp mockado

---

## 🟢 SEMANA 3: FEATURES AVANÇADAS (7 dias)

### **Dia 15-17: Report Service**

**Criar configuração**:
```yaml
# docker-compose.yml
report-service:
  build: ./services/report-service
  ports:
    - "8088:8088"
  environment:
    - PORT=8088
    - DB_HOST=postgres
```

**Endpoints básicos**:
- GET `/api/v1/reports` - Lista relatórios
- GET `/api/v1/dashboards` - Lista dashboards
- GET `/api/v1/kpis/calculate` - KPIs básicos

### **Dia 18-19: AI Service Funcional**

**Integração OpenAI**:
```python
import openai

@app.post("/api/v1/analysis/document")
def analyze_document(data: dict):
    response = openai.ChatCompletion.create(
        model="gpt-3.5-turbo",
        messages=[{
            "role": "user", 
            "content": f"Analise este documento jurídico: {data['text']}"
        }]
    )
    return {"analysis": response.choices[0].message.content}
```

### **Dia 20-21: MCP Service**

**Interface conversacional básica**:
- Comandos predefinidos
- Integração com outros serviços
- Respostas estruturadas

---

## 🔧 SEMANA 4: POLIMENTO (7 dias)

### **Dia 22-24: Performance e Otimização**
- Índices de banco otimizados
- Cache Redis para queries frequentes
- Otimização de consultas

### **Dia 25-26: Documentação e Testes**
- Documentação APIs com Swagger
- Testes automatizados básicos
- Guias de uso

### **Dia 27-28: Deploy e Monitoramento**
- Scripts de deploy
- Logs estruturados
- Métricas básicas

---

## 📊 MARCOS DE ENTREGA

### **MVP 1 (Semana 1): Sistema Básico**
- ✅ Login funcional
- ✅ Dashboard com dados reais
- ✅ CRUD de processos
- ✅ Billing informativo

### **MVP 2 (Semana 2): Serviços Conectados**
- ✅ Busca básica
- ✅ Notificações email
- ✅ DataJud mockado
- ✅ AI básica

### **MVP 3 (Semana 3): Plataforma Completa**
- ✅ Relatórios avançados
- ✅ AI funcional
- ✅ Interface conversacional
- ✅ Features premium

### **MVP 4 (Semana 4): Produção**
- ✅ Performance otimizada
- ✅ Documentação completa
- ✅ Deploy automatizado
- ✅ Monitoramento ativo

---

## 🎯 DEFINIÇÃO DE "PRONTO"

### Para cada Endpoint:
- [ ] Responde corretamente
- [ ] Trata erros adequadamente
- [ ] Tem dados de teste
- [ ] Frontend conectado
- [ ] Documentado no Swagger

### Para cada Serviço:
- [ ] Health check funciona
- [ ] Logs estruturados
- [ ] Métricas básicas
- [ ] Dockerfile funcional
- [ ] Variáveis de ambiente

### Para cada Feature:
- [ ] Testada manualmente
- [ ] Integrada com outros serviços
- [ ] UX adequada
- [ ] Performance aceitável
- [ ] Documentação atualizada

---

## 🚨 RISCOS E MITIGAÇÕES

### Risco 1: Scope Creep
**Mitigação**: Foco rigoroso no MVP, features extras para depois

### Risco 2: Bloqueios técnicos
**Mitigação**: Implementações simples primeiro, refinamento depois

### Risco 3: Integração complexa
**Mitigação**: Serviços independentes, acoplamento mínimo

### Risco 4: Performance
**Mitigação**: Otimização apenas após funcionalidade

---

## 💡 PRINCÍPIOS DE IMPLEMENTAÇÃO

### 1. **Simplicity First**
- Implementações diretas
- Menos abstrações
- Código legível

### 2. **Functionality Over Perfection**
- Funciona > Perfeito
- Iteração rápida
- Feedback constante

### 3. **Vertical Over Horizontal**
- Um fluxo completo
- Features end-to-end
- Valor real entregue

### 4. **Real Data Only**
- Sem mocks em produção
- Dados de teste realistas
- Integração real

---

**Criado em**: 05/01/2025  
**Revisão**: Após cada marco  
**Status**: Plano aprovado para execução