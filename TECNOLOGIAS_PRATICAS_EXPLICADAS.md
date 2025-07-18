# 🛠️ TECNOLOGIAS NA PRÁTICA - EXEMPLOS REAIS

## 🎯 **COMO FUNCIONAM NA REALIDADE**

Para esclarecer completamente as tecnologias específicas dos serviços MCP, DataJud e Search, vou mostrar **exemplos práticos** de como funcionam:

---

## 🤖 **MCP SERVICE - EXEMPLO PRÁTICO**

### **🔄 Fluxo Real de Uso**

**1. Advogado manda mensagem no WhatsApp:**
```
👤 Advogado: "Mostre processos pendentes do cliente João Silva"
```

**2. MCP Service processa via Claude API:**
```go
// Código real do MCP Service
func (s *MCPService) ProcessWhatsAppMessage(ctx context.Context, msg *WhatsAppMessage) error {
    // 1. Parse da mensagem em linguagem natural
    tool, params := s.parseNaturalLanguage(msg.Text)
    
    // 2. Executa ferramenta específica
    switch tool {
    case "process_search":
        return s.executeProcessSearch(ctx, params)
    case "client_search":
        return s.executeClientSearch(ctx, params)
    }
}

// Chama Claude API
claudeResponse := s.claudeClient.Complete(ctx, &anthropic.CompletionRequest{
    Model: "claude-3-5-sonnet-20241022",
    Messages: []anthropic.Message{
        {
            Role: "user",
            Content: "Extraia parâmetros: 'Mostre processos pendentes do cliente João Silva'",
        },
    },
    Tools: s.getAvailableTools(), // 17+ ferramentas
})
```

**3. Sistema responde no WhatsApp:**
```
🤖 Bot: Encontrei 3 processos pendentes para João Silva:

📋 1001234-56.2024.8.26.0100 (TJSP)
Status: Aguardando decisão
Última atualização: 2 dias atrás

📋 2002345-67.2024.8.26.0200 (TJSP)  
Status: Audiência marcada
Próxima data: 20/01/2025

📋 3003456-78.2024.8.26.0300 (TJSP)
Status: Recurso pendente
Prazo: 5 dias restantes

Deseja detalhes de algum processo específico?
```

---

## 🏛️ **DATAJUD SERVICE - EXEMPLO PRÁTICO**

### **🔄 Fluxo Real de Consulta CNJ**

**1. Sistema faz consulta automática (polling 30/30min):**
```go
// Código real do DataJud Service
func (s *DataJudService) QueryProcessUpdates(ctx context.Context, processNumber string) (*ProcessData, error) {
    // 1. Seleciona CNPJ do pool (round-robin)
    cnpj := s.cnpjPool.GetNextAvailable()
    
    // 2. Monta query Elasticsearch (padrão CNJ)
    query := &ElasticsearchQuery{
        Query: ElasticsearchQueryClause{
            Bool: &ElasticsearchBoolQuery{
                Must: []interface{}{
                    map[string]interface{}{
                        "term": map[string]interface{}{
                            "numeroProcesso": processNumber,
                        },
                    },
                    map[string]interface{}{
                        "term": map[string]interface{}{
                            "codigoTribunal": "TJSP",
                        },
                    },
                },
            },
        },
        Size: 10,
    }
    
    // 3. HTTP request para API CNJ
    resp, err := s.httpClient.Post(
        "https://api-publica.datajud.cnj.jus.br/api_publica_tjsp/_search",
        "application/json",
        queryJSON,
    )
    
    // 4. Parse response Elasticsearch
    return s.parseElasticsearchResponse(resp)
}
```

**2. Response real da API CNJ:**
```json
{
  "took": 5,
  "timed_out": false,
  "hits": {
    "total": {"value": 1, "relation": "eq"},
    "hits": [
      {
        "_index": "tjsp-processos",
        "_id": "1001234-56.2024.8.26.0100",
        "_score": 1.0,
        "_source": {
          "numeroProcesso": "1001234-56.2024.8.26.0100",
          "codigoTribunal": "TJSP",
          "dataUltimaAtualizacao": "2024-01-15T10:30:00Z",
          "movimentos": [
            {
              "dataMovimento": "2024-01-15T10:30:00Z",
              "codigoMovimento": "123",
              "descricaoMovimento": "Juntada de petição da parte autora",
              "tipoMovimento": "JUNTADA"
            }
          ]
        }
      }
    ]
  }
}
```

**3. Sistema detecta mudança e dispara notificação:**
```go
// Detector de mudanças
if lastUpdate.After(process.LastKnownUpdate) {
    // Novo movimento detectado!
    event := &ProcessUpdateEvent{
        ProcessNumber: processNumber,
        NewMovement: movement,
        Timestamp: time.Now(),
    }
    s.eventBus.Publish(ctx, event)
}
```

---

## 🔍 **SEARCH SERVICE - EXEMPLO PRÁTICO**

### **🔄 Fluxo Real de Busca**

**1. Advogado busca no dashboard:**
```
🔍 Busca: "responsabilidade civil médica danos morais"
```

**2. Search Service processa:**
```go
// Código real do Search Service
func (s *SearchService) Search(ctx context.Context, query string) (*SearchResponse, error) {
    // 1. Verifica cache Redis primeiro
    cacheKey := fmt.Sprintf("search:%s:%s", query, tenantID)
    if cached := s.redis.Get(ctx, cacheKey); cached != nil {
        return cached, nil
    }
    
    // 2. Monta query Elasticsearch multi-match
    esQuery := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []interface{}{
                    map[string]interface{}{
                        "multi_match": map[string]interface{}{
                            "query": query,
                            "fields": []string{
                                "titulo^2",    // Peso 2x
                                "conteudo",    // Peso 1x
                                "resumo^1.5", // Peso 1.5x
                            },
                            "fuzziness": "AUTO",
                            "type": "best_fields",
                        },
                    },
                },
                "filter": []interface{}{
                    map[string]interface{}{
                        "term": map[string]interface{}{
                            "tenant_id": tenantID,
                        },
                    },
                },
            },
        },
        "highlight": map[string]interface{}{
            "fields": map[string]interface{}{
                "titulo": map[string]interface{}{},
                "conteudo": map[string]interface{}{
                    "fragment_size": 150,
                    "number_of_fragments": 3,
                },
            },
        },
        "size": 20,
    }
    
    // 3. Executa busca no Elasticsearch
    resp, err := s.esClient.Search(
        s.esClient.Search.WithContext(ctx),
        s.esClient.Search.WithIndex("direito_lux_processos"),
        s.esClient.Search.WithBody(esQuery),
    )
    
    // 4. Salva no cache Redis (5 minutos)
    s.redis.Set(ctx, cacheKey, result, 5*time.Minute)
    
    return result, nil
}
```

**3. Response do Elasticsearch:**
```json
{
  "took": 12,
  "hits": {
    "total": {"value": 847, "relation": "eq"},
    "max_score": 2.45,
    "hits": [
      {
        "_index": "direito_lux_processos",
        "_score": 2.45,
        "_source": {
          "numero_processo": "1001234-56.2024.8.26.0100",
          "titulo": "Ação de Indenização por Danos Morais",
          "tribunal": "TJSP",
          "data_distribuicao": "2024-01-10",
          "status": "Em andamento",
          "cliente": "João Silva",
          "area_juridica": "Direito Civil"
        },
        "highlight": {
          "titulo": ["Ação de Indenização por <em>Danos Morais</em>"],
          "conteudo": [
            "...configurada a <em>responsabilidade civil</em> do réu...",
            "...os <em>danos morais</em> sofridos pelo autor...",
            "...tratamento <em>médico</em> inadequado..."
          ]
        }
      }
    ]
  }
}
```

**4. Frontend mostra resultado:**
```typescript
// Resultado renderizado no dashboard
const SearchResults = () => (
  <div className="search-results">
    <h3>847 resultados encontrados</h3>
    
    <div className="result-item">
      <h4>Ação de Indenização por <mark>Danos Morais</mark></h4>
      <p>Processo: 1001234-56.2024.8.26.0100 • TJSP</p>
      <p>Cliente: João Silva • Status: Em andamento</p>
      
      <div className="highlights">
        <p>...configurada a <mark>responsabilidade civil</mark> do réu...</p>
        <p>...os <mark>danos morais</mark> sofridos pelo autor...</p>
        <p>...tratamento <mark>médico</mark> inadequado...</p>
      </div>
      
      <span className="score">Relevância: 98%</span>
    </div>
  </div>
);
```

---

## 🎯 **INTEGRAÇÃO DOS 3 SERVIÇOS**

### **🔄 Fluxo Completo Real**

**1. Advogado pergunta no WhatsApp:**
```
👤 "Encontre processos similares ao 1001234-56.2024.8.26.0100"
```

**2. MCP Service coordena:**
```go
// MCP chama DataJud para obter dados do processo
processData := s.dataJudService.GetProcessDetails(ctx, "1001234-56.2024.8.26.0100")

// MCP chama Search para encontrar similares
searchQuery := fmt.Sprintf("area:%s tribunal:%s", processData.Area, processData.Tribunal)
similarProcesses := s.searchService.Search(ctx, searchQuery)

// MCP usa Claude para gerar resposta inteligente
response := s.claudeClient.GenerateResponse(ctx, processData, similarProcesses)
```

**3. Bot responde no WhatsApp:**
```
🤖 Encontrei 5 processos similares ao 1001234-56.2024.8.26.0100:

📋 Processo: 2002345-67.2024.8.26.0200
Similaridade: 94%
Cliente: Maria Santos  
Status: Julgado procedente
💰 Valor: R$ 15.000 (danos morais)

📋 Processo: 3003456-78.2024.8.26.0300
Similaridade: 89%
Cliente: Pedro Oliveira
Status: Em recurso
⏰ Tempo tramitação: 8 meses

Deseja análise detalhada ou jurisprudência relacionada?
```

---

## ✅ **CONCLUSÃO PRÁTICA**

### **🤖 MCP Service:**
- **Real**: Bot WhatsApp que entende linguagem natural
- **Tecnologia**: Claude API + 17 ferramentas específicas
- **Resultado**: Advogado controla sistema via conversa

### **🏛️ DataJud Service:**
- **Real**: Consulta oficial CNJ com Elasticsearch
- **Tecnologia**: HTTP client + query builder + circuit breaker
- **Resultado**: Dados processuais atualizados automaticamente

### **🔍 Search Service:**
- **Real**: Busca interna ultra-rápida com relevância
- **Tecnologia**: Elasticsearch 8.11 + Redis cache + vector search
- **Resultado**: Encontra processos similares em <100ms

**Agora ficou claro como essas tecnologias funcionam na prática?** 🎯