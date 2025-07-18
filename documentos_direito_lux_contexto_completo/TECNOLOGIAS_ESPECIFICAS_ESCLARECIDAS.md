# 🔧 TECNOLOGIAS ESPECÍFICAS ESCLARECIDAS - DIREITO LUX

## 📋 **ESCLARECIMENTOS TÉCNICOS SOLICITADOS**

### **🎯 OBJETIVO**
Esclarecer as tecnologias que estavam "às escuras" para o usuário:
- **MCP Service** (Model Context Protocol)
- **DataJud Service** (CNJ API Integration)
- **Search Service** (Elasticsearch + Vector Search)

---

## 🤖 **MCP SERVICE - MODEL CONTEXT PROTOCOL**

### **📋 O que é MCP?**
O Model Context Protocol (MCP) é um protocolo desenvolvido pela Anthropic para padronizar como assistentes de IA interagem com ferramentas e dados externos.

### **🔧 Implementação no Direito Lux**

#### **Arquitetura MCP**
```go
// MCP Service - Porta 8088
type MCPService struct {
    toolRegistry   *ToolRegistry
    sessionManager *SessionManager
    claudeClient   *AnthropicClient
    eventBus      *EventBus
}

// Ferramentas jurídicas disponíveis
type ToolRegistry struct {
    processTools     []ProcessTool
    searchTools      []SearchTool
    notificationTools []NotificationTool
    aiTools          []AITool
    reportTools      []ReportTool
}
```

#### **17+ Ferramentas Jurídicas**
```yaml
PROCESS_TOOLS:
├── process_query - Consultar processo específico
├── process_list - Listar processos do cliente
├── process_movements - Obter movimentações
├── process_status - Verificar status atual
└── process_summary - Resumo inteligente

SEARCH_TOOLS:
├── jurisprudence_search - Buscar jurisprudência
├── similar_cases - Casos similares
├── legal_precedents - Precedentes legais
└── semantic_search - Busca semântica

NOTIFICATION_TOOLS:
├── send_whatsapp - Enviar WhatsApp
├── send_email - Enviar email
├── send_telegram - Enviar Telegram
└── notification_template - Templates

AI_TOOLS:
├── document_analysis - Análise de documentos
├── legal_summary - Resumo jurídico
├── argument_generator - Gerador de argumentos
└── risk_assessment - Avaliação de riscos

REPORT_TOOLS:
├── process_report - Relatório de processo
├── client_dashboard - Dashboard cliente
└── performance_metrics - Métricas
```

#### **Integração com Claude 3.5 Sonnet**
```go
func (s *MCPService) ProcessConversation(ctx context.Context, message string, userID string) (*ConversationResponse, error) {
    // 1. Identificar ferramentas necessárias
    tools := s.toolRegistry.IdentifyTools(message)
    
    // 2. Preparar contexto MCP
    mcpContext := &MCPContext{
        UserID:        userID,
        TenantID:      s.getTenantID(userID),
        AvailableTools: tools,
        ConversationHistory: s.getHistory(userID),
    }
    
    // 3. Enviar para Claude com MCP
    response, err := s.claudeClient.ProcessWithMCP(ctx, &ClaudeRequest{
        Message:    message,
        Context:    mcpContext,
        Tools:      tools,
        MaxTokens:  4096,
    })
    
    // 4. Executar ferramentas se necessário
    if len(response.ToolCalls) > 0 {
        for _, toolCall := range response.ToolCalls {
            result, err := s.executeTool(ctx, toolCall)
            if err != nil {
                return nil, err
            }
            response.ToolResults = append(response.ToolResults, result)
        }
    }
    
    return response, nil
}
```

### **🔌 Integração com WhatsApp/Telegram**
```go
// Luxia Bot - Assistente inteligente
type LuxiaBot struct {
    mcpService     *MCPService
    whatsappClient *WhatsAppClient
    telegramClient *TelegramClient
}

func (b *LuxiaBot) HandleWhatsAppMessage(ctx context.Context, msg *WhatsAppMessage) error {
    // 1. Identificar usuário
    user := b.identifyUser(msg.From)
    
    // 2. Processar via MCP
    response, err := b.mcpService.ProcessConversation(ctx, msg.Text, user.ID)
    if err != nil {
        return err
    }
    
    // 3. Enviar resposta formatada
    return b.whatsappClient.SendMessage(msg.From, response.Text)
}
```

---

## 🏛️ **DATAJUD SERVICE - CNJ API INTEGRATION**

### **📋 O que é DataJud?**
DataJud é a API oficial do CNJ (Conselho Nacional de Justiça) que centraliza dados processuais de todos os tribunais brasileiros.

### **🔧 Implementação no Direito Lux**

#### **Arquitetura DataJud**
```go
// DataJud Service - Porta 8084
type DataJudService struct {
    httpClient     *DataJudHTTPClient
    cnpjPool       *CNPJPoolManager
    rateLimiter    *RateLimiter
    circuitBreaker *CircuitBreaker
    cache          *CacheManager
}

// Client HTTP Real (substitui mock)
type DataJudHTTPClient struct {
    baseURL        string
    apiKey         string
    certificate    *tls.Certificate
    timeout        time.Duration
    retryAttempts  int
}
```

#### **Pool de CNPJs**
```go
// Sistema de rotação de CNPJs
type CNPJPoolManager struct {
    providers []CNPJProvider
    current   int
    mutex     sync.RWMutex
}

type CNPJProvider struct {
    CNPJ        string
    Name        string
    Active      bool
    DailyLimit  int
    UsedToday   int
    LastUsed    time.Time
}

func (p *CNPJPoolManager) GetNextCNPJ() (*CNPJProvider, error) {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    
    // Rotacionar CNPJs para distribuir carga
    for i := 0; i < len(p.providers); i++ {
        provider := &p.providers[p.current]
        p.current = (p.current + 1) % len(p.providers)
        
        if provider.Active && provider.UsedToday < provider.DailyLimit {
            provider.UsedToday++
            provider.LastUsed = time.Now()
            return provider, nil
        }
    }
    
    return nil, errors.New("no available CNPJ providers")
}
```

#### **Rate Limiting e Circuit Breaker**
```go
// Rate Limiter - 120 requests/minuto
type RateLimiter struct {
    requests    int
    windowStart time.Time
    limit       int
    window      time.Duration
}

func (r *RateLimiter) Allow() bool {
    now := time.Now()
    
    // Reset window if expired
    if now.Sub(r.windowStart) > r.window {
        r.requests = 0
        r.windowStart = now
    }
    
    if r.requests >= r.limit {
        return false
    }
    
    r.requests++
    return true
}

// Circuit Breaker - Proteção contra falhas
type CircuitBreaker struct {
    state         CircuitState
    failures      int
    lastFailTime  time.Time
    timeout       time.Duration
    maxFailures   int
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    if cb.state == CircuitOpen {
        if time.Since(cb.lastFailTime) > cb.timeout {
            cb.state = CircuitHalfOpen
        } else {
            return errors.New("circuit breaker is open")
        }
    }
    
    err := fn()
    if err != nil {
        cb.failures++
        cb.lastFailTime = time.Now()
        
        if cb.failures >= cb.maxFailures {
            cb.state = CircuitOpen
        }
        
        return err
    }
    
    cb.failures = 0
    cb.state = CircuitClosed
    return nil
}
```

#### **Query para Elasticsearch**
```go
// Construção de queries para CNJ
type ElasticsearchQueryBuilder struct {
    client *elasticsearch.Client
}

func (b *ElasticsearchQueryBuilder) BuildProcessQuery(processNumber string) map[string]interface{} {
    return map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []map[string]interface{}{
                    {
                        "term": map[string]interface{}{
                            "numeroProcesso": processNumber,
                        },
                    },
                },
            },
        },
        "sort": []map[string]interface{}{
            {
                "dataUltimaAtualizacao": map[string]interface{}{
                    "order": "desc",
                },
            },
        },
        "size": 1,
    }
}
```

### **📊 Mapeamento de Tribunais**
```go
// Mapeamento de códigos de tribunal
var TribunalMapping = map[string]string{
    "8.26": "TJSP", // São Paulo
    "8.19": "TJRJ", // Rio de Janeiro
    "8.12": "TJMG", // Minas Gerais
    "8.25": "TJSP", // São Paulo (2ª instância)
    "4.03": "TRF3", // Tribunal Regional Federal 3ª Região
    "5.01": "TST",  // Tribunal Superior do Trabalho
}

func (s *DataJudService) GetTribunalByCode(code string) string {
    if tribunal, exists := TribunalMapping[code]; exists {
        return tribunal
    }
    return "UNKNOWN"
}
```

---

## 🔍 **SEARCH SERVICE - ELASTICSEARCH + VECTOR SEARCH**

### **📋 O que é Search Service?**
Serviço responsável por busca avançada, jurisprudência e busca semântica usando Elasticsearch e PostgreSQL com pgvector.

### **🔧 Implementação no Direito Lux**

#### **Arquitetura Search**
```go
// Search Service - Porta 8086
type SearchService struct {
    elasticClient *elasticsearch.Client
    vectorStore   *VectorStore
    aiService     *AIServiceClient
    cache         *RedisCache
}

// Vector Store com pgvector
type VectorStore struct {
    db *sql.DB
}

// Elasticsearch Configuration
type ElasticsearchConfig struct {
    Addresses    []string
    Username     string
    Password     string
    IndexPrefix  string
    Shards       int
    Replicas     int
}
```

#### **Busca Híbrida: Texto + Vetores**
```go
func (s *SearchService) HybridSearch(ctx context.Context, query string, tenantID string) (*SearchResults, error) {
    // 1. Busca tradicional no Elasticsearch
    textResults, err := s.elasticTextSearch(ctx, query, tenantID)
    if err != nil {
        return nil, err
    }
    
    // 2. Busca vetorial semântica
    vectorResults, err := s.vectorSemanticSearch(ctx, query, tenantID)
    if err != nil {
        return nil, err
    }
    
    // 3. Combinar e rankear resultados
    combinedResults := s.combineResults(textResults, vectorResults)
    
    return &SearchResults{
        Results:     combinedResults,
        Total:       len(combinedResults),
        QueryTime:   time.Since(start),
        Method:      "hybrid",
    }, nil
}
```

#### **Busca Semântica com pgvector**
```go
func (s *SearchService) vectorSemanticSearch(ctx context.Context, query string, tenantID string) ([]SearchResult, error) {
    // 1. Gerar embedding da query
    embedding, err := s.aiService.GenerateEmbedding(ctx, query)
    if err != nil {
        return nil, err
    }
    
    // 2. Busca por similaridade
    sql := `
        SELECT 
            document_id,
            content,
            metadata,
            1 - (embedding <=> $1) as similarity_score
        FROM document_embeddings 
        WHERE tenant_id = $2 
        ORDER BY embedding <=> $1 
        LIMIT 20
    `
    
    rows, err := s.vectorStore.db.QueryContext(ctx, sql, 
        pq.Array(embedding), tenantID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var results []SearchResult
    for rows.Next() {
        var result SearchResult
        err := rows.Scan(
            &result.DocumentID,
            &result.Content,
            &result.Metadata,
            &result.SimilarityScore,
        )
        if err != nil {
            continue
        }
        results = append(results, result)
    }
    
    return results, nil
}
```

#### **Indexação de Documentos**
```go
func (s *SearchService) IndexDocument(ctx context.Context, doc *Document) error {
    // 1. Indexar no Elasticsearch (busca textual)
    elasticDoc := map[string]interface{}{
        "content":     doc.Content,
        "title":       doc.Title,
        "type":        doc.Type,
        "tenant_id":   doc.TenantID,
        "created_at":  doc.CreatedAt,
        "metadata":    doc.Metadata,
    }
    
    _, err := s.elasticClient.Index(
        fmt.Sprintf("%s_documents", doc.TenantID),
        strings.NewReader(toJSON(elasticDoc)),
        s.elasticClient.Index.WithDocumentID(doc.ID),
    )
    if err != nil {
        return err
    }
    
    // 2. Gerar embedding e indexar no pgvector
    embedding, err := s.aiService.GenerateEmbedding(ctx, doc.Content)
    if err != nil {
        return err
    }
    
    _, err = s.vectorStore.db.ExecContext(ctx, `
        INSERT INTO document_embeddings 
        (id, tenant_id, document_type, document_id, content, embedding, metadata)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) DO UPDATE SET
        content = $5,
        embedding = $6,
        metadata = $7
    `, doc.ID, doc.TenantID, doc.Type, doc.ID, doc.Content, 
       pq.Array(embedding), toJSON(doc.Metadata))
    
    return err
}
```

#### **Agregações e Facetas**
```go
func (s *SearchService) SearchWithFacets(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
    query := map[string]interface{}{
        "query": s.buildQuery(req),
        "aggs": map[string]interface{}{
            "by_tribunal": map[string]interface{}{
                "terms": map[string]interface{}{
                    "field": "tribunal.keyword",
                    "size":  10,
                },
            },
            "by_class": map[string]interface{}{
                "terms": map[string]interface{}{
                    "field": "classe.keyword",
                    "size":  20,
                },
            },
            "by_year": map[string]interface{}{
                "date_histogram": map[string]interface{}{
                    "field":    "dataAjuizamento",
                    "interval": "year",
                },
            },
        },
        "size": req.Size,
        "from": req.From,
    }
    
    res, err := s.elasticClient.Search(
        s.elasticClient.Search.WithContext(ctx),
        s.elasticClient.Search.WithIndex(req.Index),
        s.elasticClient.Search.WithBody(strings.NewReader(toJSON(query))),
    )
    
    // Processar resposta e extrair facetas
    return s.parseSearchResponse(res)
}
```

### **🔄 Cache Strategy**
```go
type CacheStrategy struct {
    redis  *redis.Client
    ttl    time.Duration
}

func (s *SearchService) CachedSearch(ctx context.Context, query string, tenantID string) (*SearchResults, error) {
    // 1. Verificar cache
    cacheKey := fmt.Sprintf("search:%s:%s", tenantID, hashQuery(query))
    cached, err := s.cache.Get(ctx, cacheKey)
    if err == nil {
        var results SearchResults
        json.Unmarshal([]byte(cached), &results)
        return &results, nil
    }
    
    // 2. Executar busca
    results, err := s.HybridSearch(ctx, query, tenantID)
    if err != nil {
        return nil, err
    }
    
    // 3. Salvar no cache
    serialized, _ := json.Marshal(results)
    s.cache.Set(ctx, cacheKey, serialized, 5*time.Minute)
    
    return results, nil
}
```

---

## 🔄 **INTEGRAÇÃO ENTRE SERVIÇOS**

### **Event-Driven Communication**
```go
// Eventos entre serviços
type ProcessIndexedEvent struct {
    ProcessID    string    `json:"process_id"`
    TenantID     string    `json:"tenant_id"`
    Content      string    `json:"content"`
    Timestamp    time.Time `json:"timestamp"`
}

// Search Service escuta eventos do Process Service
func (s *SearchService) HandleProcessIndexedEvent(event *ProcessIndexedEvent) error {
    doc := &Document{
        ID:        event.ProcessID,
        TenantID:  event.TenantID,
        Content:   event.Content,
        Type:      "process",
        CreatedAt: event.Timestamp,
    }
    
    return s.IndexDocument(context.Background(), doc)
}
```

### **API Gateway Pattern**
```go
// Roteamento entre serviços
type ServiceRouter struct {
    mcpService    *MCPServiceClient
    searchService *SearchServiceClient
    datajudService *DataJudServiceClient
}

func (r *ServiceRouter) RouteRequest(ctx context.Context, req *APIRequest) (*APIResponse, error) {
    switch req.Service {
    case "mcp":
        return r.mcpService.Process(ctx, req)
    case "search":
        return r.searchService.Search(ctx, req)
    case "datajud":
        return r.datajudService.Query(ctx, req)
    default:
        return nil, errors.New("unknown service")
    }
}
```

---

## 🎯 **RESULTADO FINAL**

### **✅ Tecnologias Esclarecidas**

1. **MCP Service**: Protocolo Claude + 17 ferramentas jurídicas + integração WhatsApp/Telegram
2. **DataJud Service**: API CNJ + pool CNPJs + rate limiting + circuit breaker
3. **Search Service**: Elasticsearch + pgvector + busca híbrida + cache inteligente

### **🔧 Integração Completa**
- Event-driven architecture
- API Gateway pattern
- Observabilidade nativa
- Cache multi-layer
- Resilience patterns

**🚀 TODAS AS TECNOLOGIAS ESPECÍFICAS FORAM COMPLETAMENTE ESCLARECIDAS!**