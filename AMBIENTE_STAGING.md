# 🎯 Ambiente STAGING - Direito Lux

**Data de Criação**: 08/01/2025  
**Status**: ⚠️ **CRÍTICO** - Implementação obrigatória antes da produção

## 🚨 **CONTEXTO CRÍTICO**

**⚠️ DESCOBERTA IMPORTANTE**: Auditoria externa (08/01/2025) identificou que o ambiente DEV atual usa implementações **MOCK** e tokens **DEMO**, não garantindo funcionamento em produção.

### ❌ **Problemas Identificados:**
- **DataJud Service**: Implementação completamente mock (não faz calls reais para CNJ)
- **APIs Externas**: Todas configuradas com tokens demo/mock  
- **Webhooks**: Não configurados (URLs públicas necessárias)
- **Certificados**: Autenticação CNJ digital não implementada

### ✅ **Objetivo do STAGING:**
Criar ambiente intermediário com **APIs reais** e **quotas limitadas** para validação completa antes da produção.

## 🏗️ **ARQUITETURA STAGING**

### **Infraestrutura**
```
┌─────────────────────────────────────────────────────────────┐
│                    AMBIENTE STAGING                         │
├─────────────────────────────────────────────────────────────┤
│  💻 Frontend (staging.direitolux.com.br)                   │
│  └── Next.js 14 + SSL/TLS                                  │
├─────────────────────────────────────────────────────────────┤
│  🌐 API Gateway (api.staging.direitolux.com.br)            │
│  └── Kong + Rate Limiting + HTTPS                          │
├─────────────────────────────────────────────────────────────┤
│  🔧 Microserviços com APIs REAIS (quotas limitadas)        │
│  ├── Auth Service (JWT)                                    │
│  ├── Process Service (CRUD real)                           │
│  ├── DataJud Service ⚠️ HTTP CLIENT REAL                   │
│  ├── AI Service (OpenAI real, quota limitada)              │
│  ├── Notification (WhatsApp/Telegram real, webhooks)       │
│  └── Search/MCP/Report Services                            │
├─────────────────────────────────────────────────────────────┤
│  📊 Infraestrutura                                         │
│  ├── PostgreSQL (Cloud SQL ou similar)                     │
│  ├── Redis (cache distribuído)                             │
│  ├── RabbitMQ (mensageria)                                 │
│  └── Elasticsearch (busca)                                 │
└─────────────────────────────────────────────────────────────┘
```

## 🔑 **CONFIGURAÇÕES STAGING (APIs REAIS)**

### **1. DataJud Service - IMPLEMENTAÇÃO REAL** 
```bash
# Configurações obrigatórias (ambiente staging)
DATAJUD_API_KEY=staging_cnj_api_key_real
DATAJUD_API_URL=https://api-publica.datajud.cnj.jus.br
DATAJUD_CERTIFICATE_PATH=/certs/staging_cnpj.p12
DATAJUD_CERTIFICATE_PASSWORD=staging_cert_password
DATAJUD_RATE_LIMIT_DAILY=1000                    # Limitado para staging
DATAJUD_TIMEOUT=30s
DATAJUD_RETRY_ATTEMPTS=3
```

### **2. AI Service - APIs Reais**
```bash
# OpenAI (quota limitada)
OPENAI_API_KEY=sk-staging-real-key-limited-quota
OPENAI_MODEL=gpt-3.5-turbo
OPENAI_MAX_TOKENS=2000
OPENAI_RATE_LIMIT_RPM=100                        # Limitado para staging

# HuggingFace (opcional)
HUGGINGFACE_TOKEN=staging_hf_token_limited

# Anthropic/Claude (MCP Service)
ANTHROPIC_API_KEY=sk-ant-staging-real-key
ANTHROPIC_MODEL=claude-3-5-sonnet-20241022
ANTHROPIC_MAX_TOKENS=4096
```

### **3. Notification Service - Webhooks Reais**
```bash
# WhatsApp Business API (Meta)
WHATSAPP_ACCESS_TOKEN=staging_meta_business_token
WHATSAPP_PHONE_NUMBER_ID=staging_phone_number_id
WHATSAPP_VERIFY_TOKEN=staging_verify_token_secure
WHATSAPP_WEBHOOK_URL=https://api.staging.direitolux.com.br/webhook/whatsapp

# Telegram Bot
TELEGRAM_BOT_TOKEN=staging_telegram_bot_token_real
TELEGRAM_WEBHOOK_URL=https://api.staging.direitolux.com.br/webhook/telegram
TELEGRAM_WEBHOOK_SECRET=staging_telegram_secret

# Email SMTP (real)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=staging@direitolux.com.br
SMTP_PASSWORD=staging_email_app_password
SMTP_FROM_EMAIL=staging@direitolux.com.br
SMTP_USE_TLS=true
```

### **4. URLs Públicas Obrigatórias**
```bash
# Frontend
FRONTEND_URL=https://staging.direitolux.com.br

# API Gateway  
API_GATEWAY_URL=https://api.staging.direitolux.com.br

# Webhooks (HTTPS obrigatório)
WEBHOOK_BASE_URL=https://api.staging.direitolux.com.br

# CORS Origins
CORS_ALLOWED_ORIGINS=https://staging.direitolux.com.br,https://api.staging.direitolux.com.br
```

## 🔐 **SEGURANÇA STAGING**

### **Certificados SSL/TLS**
```bash
# Certificados Let's Encrypt ou similar
SSL_CERT_PATH=/etc/ssl/certs/staging.direitolux.com.br.crt
SSL_KEY_PATH=/etc/ssl/private/staging.direitolux.com.br.key

# Certificado CNJ (obrigatório para DataJud)
CNJ_CERT_PATH=/certs/staging_cnpj.p12
CNJ_CERT_PASSWORD=staging_cert_password
CNJ_CERT_TYPE=A1  # ou A3
```

### **Autenticação e Autorização**
```bash
# JWT Secrets (únicos para staging)
JWT_SECRET=staging_jwt_secret_unique_key_2025
JWT_ALGORITHM=HS256
JWT_EXPIRATION_MINUTES=60

# API Keys rotativas
API_KEY_ROTATION_DAYS=30
API_KEY_CURRENT=staging_api_key_v1
API_KEY_PREVIOUS=staging_api_key_v0
```

## 🧪 **IMPLEMENTAÇÕES OBRIGATÓRIAS**

### **1. DataJud HTTP Client Real**

**Arquivo**: `services/datajud-service/internal/infrastructure/http/datajud_client.go`

```go
// IMPLEMENTAÇÃO REAL (substitui mock)
package http

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "net/http"
    "time"
)

type DataJudClient struct {
    httpClient *http.Client
    baseURL    string
    apiKey     string
    certPath   string
    certPass   string
}

func NewDataJudClient(config *DataJudConfig) (*DataJudClient, error) {
    // Carregar certificado P12
    cert, err := loadP12Certificate(config.CertPath, config.CertPass)
    if err != nil {
        return nil, fmt.Errorf("erro ao carregar certificado: %w", err)
    }

    // Configurar TLS com certificado cliente
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        MinVersion:   tls.VersionTLS12,
    }

    // Cliente HTTP com timeout e TLS
    httpClient := &http.Client{
        Timeout: 30 * time.Second,
        Transport: &http.Transport{
            TLSClientConfig: tlsConfig,
        },
    }

    return &DataJudClient{
        httpClient: httpClient,
        baseURL:    config.APIURL,
        apiKey:     config.APIKey,
        certPath:   config.CertPath,
        certPass:   config.CertPass,
    }, nil
}

func (c *DataJudClient) QueryProcess(processNumber, courtID string) (*ProcessResponse, error) {
    // Implementação real da consulta CNJ
    url := fmt.Sprintf("%s/processo/%s/%s", c.baseURL, courtID, processNumber)
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    // Headers obrigatórios CNJ
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("User-Agent", "DireitoLux/1.0")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("erro na requisição: %w", err)
    }
    defer resp.Body.Close()

    // Parse response real
    var result ProcessResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("erro ao decodificar response: %w", err)
    }

    return &result, nil
}
```

### **2. Webhook Handlers**

**WhatsApp Webhook**: `services/notification-service/internal/infrastructure/http/handlers/whatsapp_webhook.go`

```go
func (h *WebhookHandler) HandleWhatsApp(c *gin.Context) {
    // Verificar token de verificação
    mode := c.Query("hub.mode")
    token := c.Query("hub.verify_token")
    challenge := c.Query("hub.challenge")

    if mode == "subscribe" && token == h.config.WhatsApp.VerifyToken {
        c.String(200, challenge)
        return
    }

    // Processar webhook de mensagem
    var webhook WhatsAppWebhook
    if err := c.ShouldBindJSON(&webhook); err != nil {
        c.JSON(400, gin.H{"error": "invalid payload"})
        return
    }

    // Processar mensagens recebidas
    for _, entry := range webhook.Entry {
        for _, change := range entry.Changes {
            if change.Value.Messages != nil {
                for _, message := range change.Value.Messages {
                    h.processIncomingMessage(message)
                }
            }
        }
    }

    c.JSON(200, gin.H{"status": "ok"})
}
```

### **3. Configuração Docker Compose Staging**

**Arquivo**: `docker-compose.staging.yml`

```yaml
version: '3.8'

services:
  # DataJud Service com configurações reais
  datajud-service:
    build:
      context: ./services/datajud-service
      dockerfile: Dockerfile.staging
    environment:
      - ENVIRONMENT=staging
      - DATAJUD_API_KEY=${DATAJUD_API_KEY}
      - DATAJUD_API_URL=https://api-publica.datajud.cnj.jus.br
      - DATAJUD_CERTIFICATE_PATH=/certs/staging.p12
      - DATAJUD_CERTIFICATE_PASSWORD=${CNJ_CERT_PASSWORD}
    volumes:
      - ./certs:/certs:ro
    ports:
      - "8084:8080"

  # AI Service com APIs reais
  ai-service:
    build:
      context: ./services/ai-service
      dockerfile: Dockerfile.staging
    environment:
      - ENVIRONMENT=staging
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY}
      - HUGGINGFACE_TOKEN=${HUGGINGFACE_TOKEN}
    ports:
      - "8000:8000"

  # Notification Service com webhooks
  notification-service:
    build:
      context: ./services/notification-service
      dockerfile: Dockerfile.staging
    environment:
      - ENVIRONMENT=staging
      - WHATSAPP_ACCESS_TOKEN=${WHATSAPP_ACCESS_TOKEN}
      - WHATSAPP_WEBHOOK_URL=https://api.staging.direitolux.com.br/webhook/whatsapp
      - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
      - TELEGRAM_WEBHOOK_URL=https://api.staging.direitolux.com.br/webhook/telegram
    ports:
      - "8085:8080"

  # Nginx SSL proxy
  nginx-ssl:
    image: nginx:alpine
    ports:
      - "443:443"
      - "80:80"
    volumes:
      - ./nginx/staging.conf:/etc/nginx/nginx.conf
      - ./certs/ssl:/etc/ssl/certs
    depends_on:
      - datajud-service
      - ai-service
      - notification-service
```

## 🧪 **PLANO DE TESTES STAGING**

### **1. Testes de Integração CNJ**
```bash
# Teste DataJud real
curl -X GET "https://api.staging.direitolux.com.br/datajud/v1/process/0000001-12.2023.8.02.0001/TJAL" \
  -H "Authorization: Bearer ${STAGING_JWT}" \
  -H "X-Tenant-ID: silva-associados"

# Resposta esperada (real):
{
  "status": "success",
  "data": {
    "numero": "0000001-12.2023.8.02.0001",
    "tribunal": "TJAL",
    "situacao": "Em Andamento",
    "movimentos": [...]
  },
  "from_cache": false,
  "processing_mode": "real_cnj_api"
}
```

### **2. Testes WhatsApp Real**
```bash
# Configurar webhook Meta Business
curl -X POST "https://graph.facebook.com/v17.0/${PHONE_NUMBER_ID}/webhooks" \
  -H "Authorization: Bearer ${WHATSAPP_ACCESS_TOKEN}" \
  -d '{
    "webhook_url": "https://api.staging.direitolux.com.br/webhook/whatsapp",
    "verify_token": "'${WHATSAPP_VERIFY_TOKEN}'"
  }'

# Enviar mensagem teste
curl -X POST "https://api.staging.direitolux.com.br/notification/v1/send" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: silva-associados" \
  -d '{
    "type": "whatsapp",
    "to": "+5511999999999",
    "message": "Teste staging - notificação real"
  }'
```

### **3. Testes AI/OpenAI Real**
```bash
# Análise de documento real
curl -X POST "https://api.staging.direitolux.com.br/ai/v1/analysis" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: silva-associados" \
  -d '{
    "text": "Processo de execução fiscal contra devedor...",
    "analysis_type": "summary"
  }'

# Resposta esperada (real):
{
  "result": "Resumo: Trata-se de execução fiscal...",
  "confidence": 0.95,
  "analysis_type": "summary",
  "processing_mode": "openai_real",
  "tokens_used": 150
}
```

## 🚀 **SCRIPT DE DEPLOY STAGING**

**Arquivo**: `scripts/deploy-staging.sh`

```bash
#!/bin/bash

set -e

echo "🎯 Deploy Ambiente STAGING - Direito Lux"
echo "========================================"

# 1. Verificar pré-requisitos
echo "🔍 Verificando pré-requisitos..."
check_prerequisites() {
    command -v docker >/dev/null 2>&1 || { echo "❌ Docker não encontrado"; exit 1; }
    command -v docker-compose >/dev/null 2>&1 || { echo "❌ Docker Compose não encontrado"; exit 1; }
    [ -f ".env.staging" ] || { echo "❌ Arquivo .env.staging não encontrado"; exit 1; }
    [ -f "certs/staging.p12" ] || { echo "❌ Certificado CNJ não encontrado"; exit 1; }
    echo "✅ Pré-requisitos ok"
}

# 2. Carregar variáveis staging
echo "📝 Carregando configurações staging..."
source .env.staging

# 3. Build das imagens staging
echo "🔨 Build das imagens staging..."
docker-compose -f docker-compose.staging.yml build

# 4. Implementar DataJud client real
echo "⚠️ Verificando implementação DataJud real..."
if grep -q "placeholder" services/datajud-service/internal/application/datajud_service.go; then
    echo "❌ DataJud Service ainda usa implementação mock!"
    echo "   Implementar HTTP client real antes do deploy staging"
    exit 1
fi

# 5. Deploy staging
echo "🚀 Iniciando deploy staging..."
docker-compose -f docker-compose.staging.yml up -d

# 6. Aguardar serviços
echo "⏳ Aguardando serviços ficarem prontos..."
sleep 30

# 7. Testes de saúde
echo "🧪 Executando testes de saúde..."
./scripts/test-staging-health.sh

echo "✅ Deploy staging concluído!"
echo "🌐 URLs:"
echo "   Frontend: https://staging.direitolux.com.br"
echo "   API: https://api.staging.direitolux.com.br"
echo "   Grafana: https://grafana.staging.direitolux.com.br"
```

## 📋 **CHECKLIST STAGING**

### **⚠️ Pré-Deploy (Obrigatório)**
- [ ] **Implementar DataJud HTTP Client real** (substitui mock)
- [ ] **Obter certificado digital CNJ** (A1/A3)
- [ ] **Configurar chaves API reais** (OpenAI, WhatsApp, Telegram, Anthropic)
- [ ] **Configurar domínios staging** (staging.direitolux.com.br)
- [ ] **Certificados SSL/TLS** (Let's Encrypt ou similar)

### **🔧 Configuração**
- [ ] **Arquivo .env.staging** com todas as chaves reais
- [ ] **docker-compose.staging.yml** configurado
- [ ] **Webhooks URLs** públicas configuradas
- [ ] **Rate limiting** configurado para staging
- [ ] **Quotas limitadas** para APIs pagas

### **🧪 Validação**
- [ ] **Teste DataJud real** (consulta processo CNJ)
- [ ] **Teste WhatsApp real** (envio + webhook)
- [ ] **Teste Telegram real** (bot + webhook) 
- [ ] **Teste AI real** (OpenAI/Claude)
- [ ] **Teste E2E completo** (fluxo usuário)

### **🚀 Deploy**
- [ ] **Deploy staging executado** sem erros
- [ ] **Todos os serviços healthy** 
- [ ] **Testes de integração** passando
- [ ] **Monitoramento** configurado
- [ ] **Documentação** atualizada

## 🎯 **CRONOGRAMA ESTIMADO**

| Fase | Duração | Responsabilidades |
|------|---------|------------------|
| **Implementação DataJud** | 1 dia | HTTP client real, certificado CNJ |
| **Configuração APIs** | 0.5 dia | Chaves reais, webhooks, domínios |
| **Deploy e Testes** | 0.5 dia | Docker staging, testes E2E |
| **Validação Final** | 0.5 dia | Testes completos, documentação |
| **Total** | **2-3 dias** | **Staging pronto para produção** |

## ⚠️ **OBSERVAÇÕES IMPORTANTES**

### **Custos Estimados (Staging)**
- **OpenAI**: ~$20-50/mês (quota limitada)
- **WhatsApp Business**: Grátis até 1000 mensagens/mês
- **Telegram Bot**: Grátis
- **CNJ DataJud**: Grátis até limite diário
- **Hospedagem**: $50-100/mês (VPS ou cloud básico)

### **Segurança**
- **Chaves em ambiente**: Nunca committar chaves reais no Git
- **Rotação**: Rotacionar chaves staging mensalmente  
- **Monitoramento**: Logs de acesso e rate limiting
- **Backup**: Backup diário do PostgreSQL staging

### **Limitações Staging**
- **Quotas reduzidas**: Para evitar custos altos
- **Dados limitados**: Apenas dados de teste
- **Usuários limitados**: Apenas equipe de desenvolvimento
- **Uptime**: Pode ter janelas de manutenção

---

**📋 PRÓXIMOS PASSOS**: Após criação do staging, validar que ambiente funciona 100% com APIs reais antes do deploy em produção.

**🎯 OBJETIVO**: Garantir que PRODUÇÃO funcionará exatamente como STAGING, eliminando surpresas no go-live.