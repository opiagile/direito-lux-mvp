# MCP Service - Model Context Protocol

## 🎯 Visão Geral

O **MCP Service** é o módulo de interface conversacional do Direito Lux que implementa o Model Context Protocol da Anthropic. Este serviço permite que usuários interajam com todas as funcionalidades da plataforma através de comandos naturais via bots inteligentes.

### 🏆 Diferencial Estratégico

- **Primeiro SaaS jurídico brasileiro** com interface conversacional completa
- **17+ ferramentas MCP especializadas** para o domínio jurídico
- **Multi-plataforma**: WhatsApp, Telegram, Claude Chat, Slack
- **Democratização da tecnologia jurídica** via linguagem natural

## 🛠️ Ferramentas MCP Implementadas (17+)

### 📋 Process Management (5 ferramentas)
- `process_search` - Buscar processos por critérios
- `process_monitor` - Configurar monitoramento automático
- `process_create` - Adicionar novo processo
- `process_details` - Obter detalhes completos
- `process_update_status` - Atualizar status

### 🧠 AI & Analysis (4 ferramentas)
- `jurisprudence_search` - Busca semântica em decisões
- `case_similarity_analysis` - Análise de similaridade
- `document_analysis` - Análise de documentos legais
- `legal_document_generation` - Geração de documentos

### 🔍 Search & Discovery (2 ferramentas)
- `advanced_search` - Busca avançada com filtros
- `search_suggestions` - Sugestões automáticas

### 📱 Notification & Communication (2 ferramentas)
- `notification_setup` - Configurar notificações
- `bulk_notification` - Envio em massa

### 📊 Reporting & Analytics (2 ferramentas)
- `generate_report` - Relatórios personalizados
- `dashboard_metrics` - Métricas em tempo real

### 👥 Administrative (2 ferramentas)
- `user_management` - Gerenciamento de usuários
- `tenant_analytics` - Análise por escritório

## 🏗️ Arquitetura

### Domain Layer
```
internal/domain/
├── mcp_session.go          # Sessões de conversação
├── conversation_context.go # Contexto e histórico
├── tool_registry.go        # Registro de ferramentas
├── quota.go               # Controle de quotas
├── events.go              # Eventos de domínio
└── tools/
    ├── process_tools.go    # Ferramentas de processos
    ├── ai_tools.go         # Ferramentas de IA
    ├── search_tools.go     # Ferramentas de busca
    ├── notification_tools.go # Ferramentas de notificação
    ├── report_tools.go     # Ferramentas de relatórios
    └── admin_tools.go      # Ferramentas administrativas
```

### Infrastructure Layer
```
internal/infrastructure/
├── config/
│   └── config.go          # Configurações específicas MCP
├── claude/
│   └── client.go          # Cliente Claude API
├── bots/
│   ├── whatsapp/          # Bot WhatsApp Business
│   ├── telegram/          # Bot Telegram
│   ├── claude_chat/       # Interface Claude Chat
│   └── slack/             # Bot Slack
├── quota/
│   └── manager.go         # Gerenciamento de quotas
└── session/
    └── manager.go         # Gerenciamento de sessões
```

## 📊 Sistema de Quotas por Plano

| Plano | Comandos/Mês | Interfaces Disponíveis |
|-------|--------------|------------------------|
| Starter | ❌ 0 | Nenhuma |
| Professional | ✅ 200 | WhatsApp + Telegram |
| Business | ✅ 1.000 | WhatsApp + Telegram + Claude Chat |
| Enterprise | ✅ Ilimitado | Todas + Slack + Ferramentas Customizadas |

## 🤖 Exemplos de Uso

### WhatsApp Business
```
👤: "Mostre meus processos ativos"
🤖: Encontrei 15 processos ativos:

🏛️ TJ-SP
• 1001234-56.2024.8.26.0100
  Cliente: João Silva
  Status: Aguardando decisão
  Última atualização: 2 dias atrás

Deseja ver mais detalhes ou filtrar por cliente?
```

### Telegram Bot
```
👤: /buscar responsabilidade civil médica
🤖: 🔍 Buscando jurisprudências...

📚 Encontradas 47 decisões relevantes:

🏛️ STJ - REsp 1.234.567/SP (95% similar)
"A responsabilidade civil do médico..."
📅 Julgado em: 15/03/2024

[Ver mais resultados] [Análise detalhada]
```

### Claude Chat Interface
```
👤: "Preciso de um relatório de produtividade trimestral completo"

🤖: Vou gerar um relatório completo para você:

📈 RELATÓRIO DE PRODUTIVIDADE - Q1 2024

VISÃO GERAL
• Processos novos: 127 (+15% vs Q4 2023)
• Processos concluídos: 98
• Tempo médio: 8.3 meses (-12%)
• Taxa de sucesso: 87%

POR ADVOGADO
👨‍💼 Dr. João Silva: 43 novos, 31 concluídos, R$ 145.000
👩‍💼 Dra. Maria Santos: 38 novos, 29 concluídos, R$ 132.000

Gostaria de algum gráfico específico ou análise detalhada?
```

## 🚀 Configuração e Setup

### 1. Variáveis de Ambiente

```bash
# MCP Service
MCP_SERVICE_NAME=mcp-service
MCP_VERSION=1.0.0
PORT=8088

# Claude API
ANTHROPIC_API_KEY=your_claude_api_key
ANTHROPIC_MODEL=claude-3-5-sonnet-20241022
ANTHROPIC_MAX_TOKENS=4096

# WhatsApp Business API
WHATSAPP_ACCESS_TOKEN=your_whatsapp_token
WHATSAPP_PHONE_NUMBER_ID=your_phone_number_id
WHATSAPP_VERIFY_TOKEN=your_verify_token

# Telegram Bot
TELEGRAM_BOT_TOKEN=your_telegram_token
TELEGRAM_WEBHOOK_URL=https://your-domain.com/webhook/telegram

# Slack Bot (Enterprise)
SLACK_BOT_TOKEN=your_slack_token
SLACK_SIGNING_SECRET=your_slack_secret

# Performance & Security
MCP_CACHE_TTL=300s
MCP_MAX_CONCURRENT_REQUESTS=100
MCP_RATE_LIMIT_PER_USER=60
MCP_SESSION_TIMEOUT=30m
MCP_JWT_SECRET=your_jwt_secret
```

### 2. Executar Localmente

```bash
# Instalar dependências
cd services/mcp-service
go mod tidy

# Executar
make run

# Com live reload
make dev

# Executar testes
make test
```

### 3. Docker

```bash
# Build
docker build -t direito-lux-mcp-service .

# Run
docker run -p 8088:8088 --env-file .env direito-lux-mcp-service

# Com docker-compose
docker-compose up mcp-service
```

## 📋 APIs Disponíveis

### Health Check
```bash
GET /health
GET /ready
```

### MCP Management
```bash
POST /api/v1/mcp/sessions        # Criar sessão
GET  /api/v1/mcp/sessions/:id    # Obter sessão
POST /api/v1/mcp/execute         # Executar ferramenta
GET  /api/v1/mcp/tools           # Listar ferramentas
```

### Bot Webhooks
```bash
POST /webhook/whatsapp           # WhatsApp webhook
POST /webhook/telegram           # Telegram webhook
POST /webhook/slack              # Slack webhook
```

### Quota Management
```bash
GET /api/v1/quotas/:tenant_id    # Status da quota
GET /api/v1/quotas/stats         # Estatísticas de uso
```

## 🔒 Segurança

- **JWT Authentication** com multi-tenant isolation
- **Rate Limiting** por usuário e tenant
- **Input Validation** em todas as ferramentas
- **Quota Enforcement** por plano
- **Audit Logging** completo
- **LGPD Compliance** para dados conversacionais

## 📈 Monitoramento

- **Prometheus Metrics** para todas as operações
- **Jaeger Tracing** para requests distribuídos
- **Health Checks** com dependências
- **Custom Dashboards** para métricas MCP
- **Alertas** para quotas e performance

## 🎯 Próximos Passos

1. **Integração Real** com serviços existentes (Process, AI, Search, etc.)
2. **Interface Voice** para comandos por voz
3. **Custom Tools** para Enterprise (ferramentas específicas por escritório)
4. **Multi-idiomas** (inglês, espanhol)
5. **Advanced Analytics** para uso de MCP

## 📞 Endpoints de Teste

```bash
# Testar saúde
curl http://localhost:8088/health

# Listar ferramentas disponíveis
curl -H "Authorization: Bearer <jwt>" \
     http://localhost:8088/api/v1/mcp/tools

# Executar ferramenta (exemplo)
curl -X POST \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <jwt>" \
     -d '{"tool_name":"process_search","parameters":{"status":"active"}}' \
     http://localhost:8088/api/v1/mcp/execute
```

---

**🤖 MCP Service - Democratizando acesso à tecnologia jurídica via linguagem natural**