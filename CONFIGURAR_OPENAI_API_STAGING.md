# 🤖 CONFIGURAR OPENAI API - STAGING COM BUDGET LIMITADO

## 🎯 CONFIGURAÇÃO STAGING - ORÇAMENTO: $10/MÊS

### 1. CRIAR CONTA OPENAI

#### 1.1 Acesso e Cadastro
1. Acesse: https://platform.openai.com/
2. Crie conta (gratuita)
3. Confirme email

#### 1.2 Configurar Billing
1. Acesse: https://platform.openai.com/account/billing/overview
2. Adicione método de pagamento
3. **IMPORTANTE**: Configurar limite de gastos

### 2. CONFIGURAR BUDGET LIMITS

#### 2.1 Usage Limits (Crítico para Staging)
```
Acesse: https://platform.openai.com/account/billing/limits

Configurar:
- Monthly Budget: $10.00
- Email Notifications: ✅ Enabled
- Hard Limit: ✅ Stop requests when limit reached
- Soft Limit: $8.00 (80% alert)
```

#### 2.2 Rate Limits (Tier 1 - Gratuito)
```
Requests per minute: 3
Tokens per minute: 150,000
Tokens per day: 10,000,000
```

### 3. GERAR API KEY

#### 3.1 Criar API Key Staging
```
Acesse: https://platform.openai.com/api-keys
Create new secret key:
- Name: "Direito Lux Staging"
- Permissions: All
- Expiration: 90 days
```

#### 3.2 Configurar API Key
```bash
# Exemplo de API Key (sempre sk-...)
OPENAI_API_KEY=sk-proj-abcd1234567890abcdef1234567890abcdef1234567890
```

### 4. CONFIGURAR MODELOS PARA STAGING

#### 4.1 Modelos Recomendados (Custo/Benefício)
```
Resumos de Processos:
- Modelo: gpt-4o-mini
- Custo: $0.150/1M input tokens
- Custo: $0.600/1M output tokens
- Ideal para: Resumos, análises

Embeddings (Busca):
- Modelo: text-embedding-3-small
- Custo: $0.020/1M tokens
- Ideal para: Busca semântica
```

#### 4.2 Configuração de Tokens
```bash
# Limites por requisição para economizar
MAX_TOKENS_RESUMO=500      # ~$0.30 por resumo
MAX_TOKENS_ANALISE=300     # ~$0.18 por análise
MAX_TOKENS_BUSCA=100       # ~$0.06 por busca
```

### 5. VARIÁVEIS DE AMBIENTE

```bash
# OpenAI API Configuration
OPENAI_API_KEY=sk-proj-abcd1234567890abcdef1234567890abcdef1234567890
OPENAI_MODEL_RESUMO=gpt-4o-mini
OPENAI_MODEL_ANALISE=gpt-4o-mini
OPENAI_MODEL_EMBEDDING=text-embedding-3-small
OPENAI_MAX_TOKENS_RESUMO=500
OPENAI_MAX_TOKENS_ANALISE=300
OPENAI_MAX_TOKENS_BUSCA=100
OPENAI_TEMPERATURE=0.7
```

### 6. CONFIGURAR RATE LIMITING

#### 6.1 Implementar Client Rate Limiting
```go
// Em ai-service: internal/infrastructure/openai/client.go
const (
    MaxRequestsPerMinute = 2  // Abaixo do limite de 3
    MaxTokensPerMinute   = 100000  // Abaixo do limite de 150k
)
```

#### 6.2 Configurar Retry Policy
```go
// Retry com backoff exponencial
retryPolicy := &RetryPolicy{
    MaxRetries:      3,
    InitialDelay:    1 * time.Second,
    MaxDelay:        30 * time.Second,
    BackoffFactor:   2.0,
}
```

### 7. MONITORAMENTO DE CUSTOS

#### 7.1 Implementar Cost Tracking
```go
// Tracking de custos por operação
type CostTracker struct {
    InputTokens  int64
    OutputTokens int64
    TotalCost    float64
}

func (c *CostTracker) CalculateCost(model string, inputTokens, outputTokens int64) {
    // Cálculo baseado no modelo
}
```

#### 7.2 Alertas de Custo
```bash
# Alertas quando atingir limites
COST_ALERT_THRESHOLD=8.00  # $8 (80% do budget)
COST_HARD_LIMIT=10.00      # $10 (100% do budget)
```

### 8. TESTAR CONFIGURAÇÃO

#### 8.1 Teste Básico
```bash
curl -X POST "https://api.openai.com/v1/chat/completions" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [{"role": "user", "content": "Teste OpenAI API - Direito Lux Staging"}],
    "max_tokens": 50
  }'
```

#### 8.2 Teste Embedding
```bash
curl -X POST "https://api.openai.com/v1/embeddings" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "text-embedding-3-small",
    "input": "Teste embedding - processo jurídico"
  }'
```

### 9. LIMITAÇÕES STAGING

#### 9.1 Budget Limits
- ✅ **Budget total**: $10/mês
- ✅ **Soft limit**: $8/mês (alerta)
- ✅ **Hard limit**: $10/mês (parar)

#### 9.2 Usage Estimates
```
Com gpt-4o-mini:
- Resumos: ~33 resumos/mês ($0.30 cada)
- Análises: ~55 análises/mês ($0.18 cada)
- Buscas: ~166 buscas/mês ($0.06 cada)
```

#### 9.3 Rate Limits
- ✅ **Requests**: 3/min → usar 2/min
- ✅ **Tokens**: 150k/min → usar 100k/min
- ✅ **Tokens diários**: 10M/dia

### 10. PRÓXIMOS PASSOS

1. ✅ Configurar OpenAI API com budget limitado
2. ⏳ Configurar Claude API com budget limitado
3. ⏳ Testar integração completa
4. ⏳ Validar custos e performance

### 11. CUSTO STAGING

- **OpenAI API**: $10/mês (budget limitado)
- **Estimativa uso**: ~$5-8/mês (desenvolvimento)
- **Monitoramento**: Alertas automáticos

**Total OpenAI**: $10/mês máximo para staging 💰

### 12. MIGRAÇÃO PARA PRODUÇÃO

Quando sair do staging:
1. Aumentar budget para $100-500/mês
2. Usar modelos mais potentes se necessário
3. Implementar caching para reduzir custos
4. Monitoramento avançado de custos

---

## 🔗 LINKS ÚTEIS

- [OpenAI Platform](https://platform.openai.com/)
- [Pricing](https://openai.com/pricing)
- [API Reference](https://platform.openai.com/docs/api-reference)
- [Usage Limits](https://platform.openai.com/account/billing/limits)
- [Rate Limits](https://platform.openai.com/docs/guides/rate-limits)