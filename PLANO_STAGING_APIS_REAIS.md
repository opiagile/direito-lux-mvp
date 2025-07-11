# 🚀 PLANO STAGING - APIS REAIS COM QUOTAS LIMITADAS

## 📊 Resumo Executivo

**Objetivo**: Configurar ambiente STAGING com APIs reais mas quotas limitadas para validação E2E completa antes da produção.

**Timeline**: 1-2 dias de trabalho

**Prioridade**: 🔥 **MÁXIMA** - Próximo marco crítico após DataJud real

---

## 🎯 CONTEXTO E JUSTIFICATIVA

### Por que STAGING com APIs Reais é Crítico?

1. **Validação de Integração Real** - Testar fluxos completos com APIs reais
2. **Detecção de Problemas** - Identificar issues antes da produção  
3. **Webhooks Funcionais** - Testar notificações bidirecionais reais
4. **Performance Real** - Medir latência com APIs externas reais
5. **Demonstrações Para Clientes** - Funcionalidades reais para prospects

### Estado Atual vs Objetivo

```
ANTES (Demo):                    DEPOIS (Staging Real):
┌─────────────────┐             ┌─────────────────┐
│   Serviços      │             │   Serviços      │
│                 │             │                 │
│ 📱 WhatsApp Demo │   ═══►     │ 📱 WhatsApp Real│
│ 🤖 Telegram Demo│             │ 🤖 Telegram Real│
│ 🧠 OpenAI Demo  │             │ 🧠 OpenAI Real  │
│ 🤖 Claude Demo  │             │ 🤖 Claude Real  │
└─────────────────┘             └─────────────────┘
                                        │
                                ┌─────────────────┐
                                │ Quotas Limitadas│
                                │ URLs HTTPS      │
                                │ Webhooks Ativos │
                                └─────────────────┘
```

---

## 🔧 FASE 1: WHATSAPP BUSINESS API

### 1.1 Configuração WhatsApp Meta Business (2-3 horas)

#### Passos de Setup:
1. **Criar Meta Business Account**
   - Acesso: https://business.facebook.com/
   - Criar conta business (se não existir)
   - Verificar conta business

2. **Configurar WhatsApp Business API**
   - Acessar Meta Developers: https://developers.facebook.com/
   - Criar novo app tipo "Business"
   - Adicionar produto "WhatsApp Business API"
   - Configurar número de telefone de teste

3. **Obter Tokens de Acesso**
   ```bash
   # Tokens necessários:
   WHATSAPP_ACCESS_TOKEN=EAAxxxxxx (token real de desenvolvimento)
   WHATSAPP_PHONE_NUMBER_ID=123456789012345 (número real)
   WHATSAPP_VERIFY_TOKEN=staging_webhook_verify_token
   WHATSAPP_BUSINESS_ACCOUNT_ID=123456789012345
   ```

4. **Configurar Webhook HTTPS**
   - URL necessária: `https://staging.direitolux.com.br/webhook/whatsapp`
   - Método: ngrok ou domínio staging real
   - Verificação: webhook verification token

#### Limitações Staging:
- **25 números de teste** grátis para desenvolvimento
- **1000 mensagens/mês** no tier gratuito
- **Rate limit**: 20 msg/segundo

---

## 🤖 FASE 2: TELEGRAM BOT API

### 2.1 Configuração Telegram Bot (1 hora)

#### Passos de Setup:
1. **Criar Bot via BotFather**
   - Telegram: https://t.me/BotFather
   - Comando: `/newbot`
   - Nome: "Direito Lux Staging Bot"
   - Username: `@direitolux_staging_bot`

2. **Obter Token**
   ```bash
   # Token real do bot:
   TELEGRAM_BOT_TOKEN=1234567890:AAAA-BBBBccccDDDDeeeeFFFGGGhhhhIII
   ```

3. **Configurar Webhook**
   - URL: `https://staging.direitolux.com.br/webhook/telegram`
   - Método: `setWebhook` API call

#### Limitações Staging:
- **30 mensagens/segundo** por bot
- **Ilimitado** para desenvolvimento
- **Grátis** sempre

---

## 🧠 FASE 3: OPENAI API

### 3.1 Configuração OpenAI (30 minutos)

#### Passos de Setup:
1. **Criar Conta OpenAI**
   - Acesso: https://platform.openai.com/
   - Criar account (se não existir)
   - Adicionar método de pagamento

2. **Configurar API Key**
   ```bash
   # API Key real com quota limitada:
   OPENAI_API_KEY=sk-proj-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   OPENAI_MODEL=gpt-3.5-turbo  # Modelo mais barato
   OPENAI_MAX_TOKENS=1000      # Limitação para staging
   ```

3. **Configurar Quota Limitada**
   - **Budget Limit**: $10/mês
   - **Hard Limit**: $15/mês (segurança)
   - **Usage Tracking**: Alertas em $5, $8

#### Limitações Staging:
- **$10/mês** de budget
- **gpt-3.5-turbo** apenas (mais barato)
- **1000 tokens máximo** por request
- **Monitoring ativo** de custos

---

## 🤖 FASE 4: ANTHROPIC CLAUDE API

### 4.1 Configuração Claude (30 minutos)

#### Passos de Setup:
1. **Criar Conta Anthropic**
   - Acesso: https://console.anthropic.com/
   - Criar account
   - Adicionar método de pagamento

2. **Configurar API Key**
   ```bash
   # API Key real com quota limitada:
   ANTHROPIC_API_KEY=sk-ant-api03-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   ANTHROPIC_MODEL=claude-3-haiku-20240307  # Modelo mais barato
   ANTHROPIC_MAX_TOKENS=1000
   ```

3. **Configurar Quota Limitada**
   - **Budget Limit**: $10/mês
   - **Usage Alerts**: $5, $8

---

## 🌐 FASE 5: URLS HTTPS PÚBLICAS

### 5.1 Opção A: ngrok (Rápido - 30 minutos)

#### Setup ngrok:
```bash
# Instalar ngrok
brew install ngrok  # macOS
# ou baixar de https://ngrok.com/

# Configurar túnel
ngrok http 8085  # notification-service port

# URLs resultantes:
https://abcd1234.ngrok.io/webhook/whatsapp
https://abcd1234.ngrok.io/webhook/telegram
```

#### Vantagens:
- ✅ Setup em 30 minutos
- ✅ HTTPS automático
- ✅ Ideal para staging

#### Limitações:
- ❌ URL muda a cada restart
- ❌ Rate limiting no plano gratuito

### 5.2 Opção B: Domínio Staging Real (2-3 horas)

#### Setup domínio próprio:
```bash
# Registrar subdomínio
staging.direitolux.com.br

# Configurar DNS
A record: staging.direitolux.com.br → IP_SERVIDOR
CNAME: *.staging.direitolux.com.br → staging.direitolux.com.br

# URLs resultantes:
https://staging.direitolux.com.br/webhook/whatsapp
https://staging.direitolux.com.br/webhook/telegram
```

#### Vantagens:
- ✅ URL fixa
- ✅ Profissional
- ✅ SSL configurável

---

## 🔧 FASE 6: CONFIGURAÇÃO TÉCNICA

### 6.1 Atualizar docker-compose.yml

```yaml
# Adicionar environment variables reais
notification-service:
  environment:
    # WhatsApp Real
    - WHATSAPP_ACCESS_TOKEN=${WHATSAPP_ACCESS_TOKEN_STAGING}
    - WHATSAPP_PHONE_NUMBER_ID=${WHATSAPP_PHONE_NUMBER_ID_STAGING}
    - WHATSAPP_VERIFY_TOKEN=${WHATSAPP_VERIFY_TOKEN_STAGING}
    - WHATSAPP_WEBHOOK_URL=https://staging.direitolux.com.br/webhook/whatsapp
    
    # Telegram Real
    - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN_STAGING}
    - TELEGRAM_WEBHOOK_URL=https://staging.direitolux.com.br/webhook/telegram

ai-service:
  environment:
    # OpenAI Real com limitações
    - OPENAI_API_KEY=${OPENAI_API_KEY_STAGING}
    - OPENAI_MODEL=gpt-3.5-turbo
    - OPENAI_MAX_TOKENS=1000
    - OPENAI_BUDGET_LIMIT=10  # $10/mês

mcp-service:
  environment:
    # Anthropic Real com limitações
    - ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY_STAGING}
    - ANTHROPIC_MODEL=claude-3-haiku-20240307
    - ANTHROPIC_MAX_TOKENS=1000
    
    # APIs compartilhadas
    - WHATSAPP_ACCESS_TOKEN=${WHATSAPP_ACCESS_TOKEN_STAGING}
    - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN_STAGING}
```

### 6.2 Criar .env.staging

```bash
# WHATSAPP BUSINESS API (Real com quotas limitadas)
WHATSAPP_ACCESS_TOKEN_STAGING=EAAxxxxxx...
WHATSAPP_PHONE_NUMBER_ID_STAGING=123456789012345
WHATSAPP_BUSINESS_ACCOUNT_ID_STAGING=123456789012345
WHATSAPP_VERIFY_TOKEN_STAGING=staging_webhook_verify_2025

# TELEGRAM BOT API (Real)
TELEGRAM_BOT_TOKEN_STAGING=1234567890:AAAA-BBBBccccDDDDeeeeFFFGGGhhhhIII

# OPENAI API (Real com budget limitado)
OPENAI_API_KEY_STAGING=sk-proj-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
OPENAI_MODEL_STAGING=gpt-3.5-turbo
OPENAI_MAX_TOKENS_STAGING=1000
OPENAI_BUDGET_LIMIT_STAGING=10

# ANTHROPIC CLAUDE API (Real com budget limitado)
ANTHROPIC_API_KEY_STAGING=sk-ant-api03-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
ANTHROPIC_MODEL_STAGING=claude-3-haiku-20240307
ANTHROPIC_MAX_TOKENS_STAGING=1000

# URLS PUBLICAS STAGING
WEBHOOK_BASE_URL_STAGING=https://staging.direitolux.com.br
WHATSAPP_WEBHOOK_URL_STAGING=https://staging.direitolux.com.br/webhook/whatsapp
TELEGRAM_WEBHOOK_URL_STAGING=https://staging.direitolux.com.br/webhook/telegram

# LIMITAÇÕES DE SEGURANÇA
STAGING_MODE=true
RATE_LIMIT_STAGING=true
BUDGET_ALERTS_ENABLED=true
USAGE_MONITORING=true
```

---

## 🧪 FASE 7: TESTES DE VALIDAÇÃO E2E

### 7.1 Teste WhatsApp Real (15 minutos)
```bash
# 1. Enviar mensagem via API
curl -X POST "https://graph.facebook.com/v18.0/${PHONE_NUMBER_ID}/messages" \
  -H "Authorization: Bearer ${WHATSAPP_ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "messaging_product": "whatsapp",
    "to": "+5511999999999",
    "type": "text",
    "text": {
      "body": "🚀 Teste STAGING - Direito Lux funcionando com WhatsApp real!"
    }
  }'

# 2. Verificar webhook recebido
# 3. Testar resposta automática via bot
```

### 7.2 Teste Telegram Real (15 minutos)
```bash
# 1. Configurar webhook
curl -X POST "https://api.telegram.org/bot${BOT_TOKEN}/setWebhook" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://staging.direitolux.com.br/webhook/telegram"
  }'

# 2. Enviar mensagem para bot
# 3. Verificar resposta automática
```

### 7.3 Teste OpenAI Real (10 minutos)
```bash
# 1. Testar análise de documento via AI Service
curl -X POST "http://localhost:8000/api/v1/ai/analyze" \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Analise esta petição inicial sobre responsabilidade civil.",
    "max_tokens": 500
  }'

# 2. Verificar resposta e uso de quota
```

### 7.4 Teste Claude Real (10 minutos)
```bash
# 1. Testar MCP Bot com Claude
curl -X POST "http://localhost:8088/api/v1/mcp/execute" \
  -H "Content-Type: application/json" \
  -d '{
    "tool": "document_analyze",
    "parameters": {
      "text": "Análise jurídica de contrato"
    }
  }'

# 2. Verificar resposta e uso de quota
```

---

## 📊 FASE 8: MONITORAMENTO E ALERTAS

### 8.1 Dashboard de Quotas
```javascript
// Implementar dashboard simples para monitorar:
- WhatsApp: mensagens enviadas/1000 limite
- Telegram: mensagens enviadas (ilimitado)
- OpenAI: custos/$10 budget
- Claude: custos/$10 budget
- Webhooks: requests recebidos
```

### 8.2 Alertas Automáticos
```bash
# Configurar alertas para:
- 80% da quota WhatsApp
- $8 de gasto OpenAI ($10 limite)
- $8 de gasto Claude ($10 limite)
- Webhook failures
- Rate limiting ativo
```

---

## ⚠️ RISCOS E MITIGAÇÕES

### 🚨 Riscos Técnicos

#### 1. Custos Descontrolados
**Risco**: APIs cobrarem mais que esperado
**Mitigação**: 
- Hard limits em $15/mês por API
- Monitoring automático
- Alertas em múltiplos níveis

#### 2. Webhooks Instáveis  
**Risco**: URLs HTTPS instáveis
**Mitigação**:
- Backup com ngrok + domínio real
- Health checks automáticos
- Retry automático com exponential backoff

#### 3. Rate Limiting
**Risco**: Exceder limites das APIs
**Mitigação**:
- Circuit breakers configurados
- Queue de mensagens
- Fallback para modo demo quando necessário

### 💰 Riscos de Negócio

#### 1. Tempo de Setup
**Risco**: Aprovações demoradas (Meta Business)
**Mitigação**:
- Processo paralelo com todas as APIs
- Documentação step-by-step
- Fallback para ngrok se domínio atrasar

---

## 📅 TIMELINE DETALHADO

### DIA 1 (6-8 horas)
- **08:00-10:00**: Setup WhatsApp Business API
- **10:00-11:00**: Setup Telegram Bot
- **11:00-12:00**: Setup OpenAI + Claude APIs
- **13:00-15:00**: Configurar URLs HTTPS (ngrok ou domínio)
- **15:00-17:00**: Atualizar docker-compose + .env.staging
- **17:00-18:00**: Deploy inicial e testes básicos

### DIA 2 (4-6 horas)  
- **08:00-10:00**: Testes E2E WhatsApp + Telegram
- **10:00-12:00**: Testes AI (OpenAI + Claude)
- **13:00-15:00**: Monitoramento + alertas
- **15:00-16:00**: Documentação e handover
- **16:00-17:00**: Testes finais + validação

---

## 🎯 ENTREGÁVEIS

### Configuração
1. **APIs Reais Funcionais** - WhatsApp, Telegram, OpenAI, Claude
2. **URLs HTTPS Estáveis** - Webhooks funcionando
3. **Quotas Limitadas** - Segurança contra custos
4. **Monitoramento Ativo** - Dashboard + alertas

### Documentação  
1. **Guia de Setup APIs** - Step-by-step reproduzível
2. **Configuração de Webhooks** - URLs e verificação
3. **Monitoramento de Quotas** - Como acompanhar uso
4. **Troubleshooting Guide** - Problemas comuns

### Testes
1. **Validação E2E Completa** - Fluxos reais funcionando
2. **Performance com APIs Reais** - Latência medida
3. **Segurança Testada** - Rate limits e budget limits
4. **Monitoring Testado** - Alertas funcionando

---

## 🚀 PRÓXIMOS PASSOS PÓS-STAGING

### Otimizações (1-2 semanas)
1. **Performance Tuning** - Otimizar latência com APIs externas
2. **Caching Estratégico** - Reduzir custos de API
3. **Auto-scaling** - Ajustar recursos conforme demanda
4. **Advanced Monitoring** - Métricas detalhadas

### Preparação Produção (1 semana)
1. **Budget Real** - Aumentar para $100+/mês por API
2. **Domínio Produção** - direitolux.com.br
3. **SSL Certificates** - Let's Encrypt ou comprado
4. **Backup Systems** - Redundância de APIs

---

## 🎉 CONCLUSÃO

**🎯 STAGING é o unlock final para PRODUÇÃO**

Esta implementação remove a última barreira para lançar com confiança:

- ✅ **Validação Real** com APIs externas
- ✅ **Integração Completa** testada
- ✅ **Monitoramento Ativo** de custos/quotas  
- ✅ **Demonstrações Reais** para clientes
- ✅ **Confiança Total** para go-live

**Timeline**: 1-2 dias intensivos

**Custo**: ~$50/mês (quotas limitadas staging)

**ROI**: Unlock para produção + revenue

---

*Plano criado em 11/07/2025 - Ready for execution*

📧 **Próximo passo**: Iniciar configuração das APIs reais