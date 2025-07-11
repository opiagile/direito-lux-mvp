# 🧠 CONFIGURAR CLAUDE API - STAGING COM BUDGET LIMITADO

## 🎯 CONFIGURAÇÃO STAGING - ORÇAMENTO: $10/MÊS

### 1. CRIAR CONTA ANTHROPIC

#### 1.1 Acesso e Cadastro
1. Acesse: https://console.anthropic.com/
2. Crie conta (gratuita)
3. Confirme email

#### 1.2 Configurar Billing
1. Acesse: https://console.anthropic.com/account/billing
2. Adicione método de pagamento
3. **IMPORTANTE**: Configurar limite de gastos

### 2. CONFIGURAR BUDGET LIMITS

#### 2.1 Usage Limits (Crítico para Staging)
```
Acesse: https://console.anthropic.com/account/billing

Configurar:
- Monthly Budget: $10.00
- Usage Alerts: ✅ Enabled at $8.00
- Hard Stop: ✅ Enabled at $10.00
- Email Notifications: ✅ Enabled
```

#### 2.2 Rate Limits (Tier 1)
```
Requests per minute: 10
Tokens per minute: 50,000
Tokens per day: 1,000,000
```

### 3. GERAR API KEY

#### 3.1 Criar API Key Staging
```
Acesse: https://console.anthropic.com/account/keys
Create Key:
- Name: "Direito Lux Staging"
- Permissions: All
- Expiration: 90 days
```

#### 3.2 Configurar API Key
```bash
# Exemplo de API Key (sempre sk-ant-...)
ANTHROPIC_API_KEY=sk-ant-api03-abcd1234567890abcdef1234567890abcdef1234567890
```

### 4. CONFIGURAR MODELOS PARA STAGING

#### 4.1 Modelos Recomendados (Custo/Benefício)
```
Análises Jurídicas:
- Modelo: claude-3-haiku-20240307
- Input: $0.25/1M tokens
- Output: $1.25/1M tokens
- Ideal para: Análises rápidas

Resumos Complexos:
- Modelo: claude-3-sonnet-20240229
- Input: $3.00/1M tokens
- Output: $15.00/1M tokens
- Ideal para: Resumos detalhados
```

#### 4.2 Configuração de Tokens
```bash
# Limites por requisição para economizar
MAX_TOKENS_ANALISE_HAIKU=300   # ~$0.38 por análise
MAX_TOKENS_RESUMO_SONNET=500   # ~$7.50 por resumo
MAX_TOKENS_CONSULTA=200        # ~$0.25 por consulta
```

### 5. VARIÁVEIS DE AMBIENTE

```bash
# Anthropic Claude API Configuration
ANTHROPIC_API_KEY=sk-ant-api03-abcd1234567890abcdef1234567890abcdef1234567890
CLAUDE_MODEL_ANALISE=claude-3-haiku-20240307
CLAUDE_MODEL_RESUMO=claude-3-sonnet-20240229
CLAUDE_MAX_TOKENS_ANALISE=300
CLAUDE_MAX_TOKENS_RESUMO=500
CLAUDE_MAX_TOKENS_CONSULTA=200
CLAUDE_TEMPERATURE=0.3
```

### 6. CONFIGURAR RATE LIMITING

#### 6.1 Implementar Client Rate Limiting
```go
// Em ai-service: internal/infrastructure/claude/client.go
const (
    MaxRequestsPerMinute = 8   // Abaixo do limite de 10
    MaxTokensPerMinute   = 40000  // Abaixo do limite de 50k
)
```

#### 6.2 Configurar Retry Policy
```go
// Retry com backoff exponencial
retryPolicy := &RetryPolicy{
    MaxRetries:      3,
    InitialDelay:    2 * time.Second,
    MaxDelay:        60 * time.Second,
    BackoffFactor:   2.0,
}
```

### 7. MONITORAMENTO DE CUSTOS

#### 7.1 Implementar Cost Tracking
```go
// Tracking de custos por operação
type ClaudeCostTracker struct {
    InputTokens  int64
    OutputTokens int64
    Model        string
    TotalCost    float64
}

func (c *ClaudeCostTracker) CalculateCost(model string, inputTokens, outputTokens int64) {
    // Cálculo baseado no modelo Claude
    switch model {
    case "claude-3-haiku-20240307":
        inputCost := float64(inputTokens) * 0.25 / 1000000
        outputCost := float64(outputTokens) * 1.25 / 1000000
        c.TotalCost = inputCost + outputCost
    case "claude-3-sonnet-20240229":
        inputCost := float64(inputTokens) * 3.00 / 1000000
        outputCost := float64(outputTokens) * 15.00 / 1000000
        c.TotalCost = inputCost + outputCost
    }
}
```

#### 7.2 Alertas de Custo
```bash
# Alertas quando atingir limites
CLAUDE_COST_ALERT_THRESHOLD=8.00  # $8 (80% do budget)
CLAUDE_COST_HARD_LIMIT=10.00      # $10 (100% do budget)
```

### 8. TESTAR CONFIGURAÇÃO

#### 8.1 Teste Básico com Haiku
```bash
curl -X POST "https://api.anthropic.com/v1/messages" \
  -H "x-api-key: $ANTHROPIC_API_KEY" \
  -H "Content-Type: application/json" \
  -H "anthropic-version: 2023-06-01" \
  -d '{
    "model": "claude-3-haiku-20240307",
    "max_tokens": 100,
    "messages": [{"role": "user", "content": "Teste Claude API - Direito Lux Staging"}]
  }'
```

#### 8.2 Teste Análise Jurídica
```bash
curl -X POST "https://api.anthropic.com/v1/messages" \
  -H "x-api-key: $ANTHROPIC_API_KEY" \
  -H "Content-Type: application/json" \
  -H "anthropic-version: 2023-06-01" \
  -d '{
    "model": "claude-3-haiku-20240307",
    "max_tokens": 300,
    "messages": [{"role": "user", "content": "Analise este processo: Ação de cobrança em andamento"}]
  }'
```

### 9. LIMITAÇÕES STAGING

#### 9.1 Budget Limits
- ✅ **Budget total**: $10/mês
- ✅ **Soft limit**: $8/mês (alerta)
- ✅ **Hard limit**: $10/mês (parar)

#### 9.2 Usage Estimates
```
Com claude-3-haiku (análises):
- Análises: ~26 análises/mês ($0.38 cada)

Com claude-3-sonnet (resumos):
- Resumos: ~1.3 resumos/mês ($7.50 cada)

Estratégia: 80% Haiku + 20% Sonnet
```

#### 9.3 Rate Limits
- ✅ **Requests**: 10/min → usar 8/min
- ✅ **Tokens**: 50k/min → usar 40k/min
- ✅ **Tokens diários**: 1M/dia

### 10. ESTRATÉGIA DE USO OTIMIZADA

#### 10.1 Distribuição de Modelos
```bash
# Usar Haiku para análises rápidas (80% do uso)
CLAUDE_MODEL_ANALISE=claude-3-haiku-20240307

# Usar Sonnet apenas para resumos críticos (20% do uso)
CLAUDE_MODEL_RESUMO_PREMIUM=claude-3-sonnet-20240229
```

#### 10.2 Caching Inteligente
```go
// Cache análises similares para economizar
type AnalysisCache struct {
    ProcessType string
    Analysis    string
    Timestamp   time.Time
}
```

### 11. PRÓXIMOS PASSOS

1. ✅ Configurar Claude API com budget limitado
2. ⏳ Testar integração completa
3. ⏳ Validar custos e performance
4. ⏳ Implementar monitoramento

### 12. CUSTO STAGING

- **Claude API**: $10/mês (budget limitado)
- **Estimativa uso**: ~$6-9/mês (desenvolvimento)
- **Monitoramento**: Alertas automáticos

**Total Claude**: $10/mês máximo para staging 💰

### 13. MIGRAÇÃO PARA PRODUÇÃO

Quando sair do staging:
1. Aumentar budget para $100-500/mês
2. Usar mais Claude Sonnet para resumos premium
3. Implementar caching avançado
4. Monitoramento detalhado de custos

---

## 🔗 LINKS ÚTEIS

- [Anthropic Console](https://console.anthropic.com/)
- [API Documentation](https://docs.anthropic.com/claude/reference/)
- [Pricing](https://www.anthropic.com/pricing)
- [Rate Limits](https://docs.anthropic.com/claude/reference/rate-limits)
- [Best Practices](https://docs.anthropic.com/claude/docs/best-practices)