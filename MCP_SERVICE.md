# MCP Service - Model Context Protocol para Direito Lux

## 📋 Visão Geral

O **MCP Service** é um módulo inovador que implementa o Model Context Protocol da Anthropic, permitindo que usuários interajam com todas as funcionalidades do Direito Lux através de interfaces conversacionais (bots). Esta é uma funcionalidade diferencial única no mercado jurídico brasileiro.

## 🎯 Proposta de Valor

### **Diferencial Competitivo**
- **Primeiro SaaS jurídico** no Brasil com integração MCP
- **Interface conversacional** natural para advogados
- **Democratização** do acesso à tecnologia jurídica
- **Eficiência operacional** através de comandos naturais

### **Casos de Uso Transformadores**

```
👤 "Mostre todos os processos em andamento do escritório Silva & Associados"
🤖 [MCP] → Auth + Tenant + Process Services → Lista processos filtrados

👤 "Crie um alerta para o processo 0001234-56.2024.8.26.0100 quando houver movimentação"
🤖 [MCP] → Process + Notification Services → Configura monitoramento automático

👤 "Busque jurisprudências sobre responsabilidade civil em acidentes de trânsito"
🤖 [MCP] → AI + Search Services → Análise semântica + resultados relevantes

👤 "Gere um relatório de produtividade do último mês"
🤖 [MCP] → Process + Report Services → Dashboard personalizado

👤 "Qual o status dos pagamentos dos clientes em atraso?"
🤖 [MCP] → Tenant + Process Services → Relatório financeiro
```

## 🏗️ Arquitetura Técnica

### **Stack Tecnológica**
- **Framework**: Go 1.21+ com Arquitetura Hexagonal
- **Protocol**: MCP (Model Context Protocol) v1.0
- **AI Integration**: Claude 3.5 Sonnet via Anthropic API
- **Authentication**: JWT + Multi-tenant integration
- **Communication**: gRPC + REST APIs
- **Caching**: Redis para performance
- **Monitoring**: Prometheus + Jaeger

### **Arquitetura de Integração**

```
┌─────────────────────────────────────────────────────────────────┐
│                    INTERFACE LAYER                             │
├─────────────────┬─────────────────┬─────────────────────────────┤
│   Claude Chat   │   WhatsApp      │   Telegram    │   Slack     │
│   Interface     │   Business API  │   Bot API     │   Bot API   │
└─────────┬───────┴─────────┬───────┴─────────┬─────┴─────────────┘
          │                 │                 │
          └─────────────────┼─────────────────┘
                            │
                ┌───────────▼───────────┐
                │     MCP SERVICE       │
                │  (Protocol Handler)   │
                │                       │
                │ ┌─────────────────────┐ │
                │ │   Tool Registry     │ │
                │ │ - process_tools     │ │
                │ │ - search_tools      │ │
                │ │ - ai_tools          │ │
                │ │ - notification_tools│ │
                │ │ - report_tools      │ │
                │ │ - admin_tools       │ │
                │ └─────────────────────┘ │
                │                       │
                │ ┌─────────────────────┐ │
                │ │   Context Manager   │ │
                │ │ - Session state     │ │
                │ │ - User context      │ │
                │ │ - Tenant isolation  │ │
                │ └─────────────────────┘ │
                └───────────┬───────────┘
                            │
                ┌───────────▼───────────┐
                │     API GATEWAY       │
                │   (Kong/Traefik)      │
                └───────────┬───────────┘
                            │
    ┌───────────────────────┼───────────────────────┐
    │                       │                       │
┌───▼───┐  ┌─────▼─────┐  ┌─▼───┐  ┌─────▼─────┐  ┌─▼────┐
│ Auth  │  │  Process  │  │ AI  │  │  Search   │  │Report│
│Service│  │  Service  │  │ Svc │  │  Service  │  │ Svc  │
└───────┘  └───────────┘  └─────┘  └───────────┘  └──────┘
```

## 🛠️ Ferramentas MCP Implementadas

### **1. Process Management Tools**

#### **process_search**
```json
{
  "name": "process_search",
  "description": "Buscar processos por diversos critérios",
  "parameters": {
    "client_name": "string (opcional)",
    "process_number": "string (opcional)", 
    "status": "enum [active, archived, suspended]",
    "date_range": "object {start, end}",
    "court": "string (opcional)",
    "lawyer": "string (opcional)"
  }
}
```

#### **process_monitor**
```json
{
  "name": "process_monitor",
  "description": "Configurar monitoramento automático de processos",
  "parameters": {
    "process_id": "string (obrigatório)",
    "notification_types": "array [movement, deadline, court_decision]",
    "channels": "array [whatsapp, email, telegram]",
    "frequency": "enum [immediate, daily, weekly]"
  }
}
```

#### **process_create**
```json
{
  "name": "process_create",
  "description": "Adicionar novo processo ao monitoramento",
  "parameters": {
    "cnj_number": "string (obrigatório)",
    "client_name": "string (obrigatório)",
    "process_type": "string",
    "responsible_lawyer": "string",
    "tags": "array de strings"
  }
}
```

### **2. AI & Analysis Tools**

#### **jurisprudence_search**
```json
{
  "name": "jurisprudence_search", 
  "description": "Busca semântica em decisões judiciais",
  "parameters": {
    "query": "string (obrigatório)",
    "court_types": "array [STF, STJ, TJ, TRT, TST]",
    "similarity_threshold": "float (0.0-1.0)",
    "max_results": "integer (default: 10)",
    "date_range": "object {start, end}"
  }
}
```

#### **case_similarity_analysis**
```json
{
  "name": "case_similarity_analysis",
  "description": "Análise de similaridade entre casos",
  "parameters": {
    "base_case": "string - descrição do caso",
    "comparison_cases": "array de IDs ou descrições", 
    "analysis_dimensions": "array [semantic, legal, factual, procedural]"
  }
}
```

#### **document_analysis**
```json
{
  "name": "document_analysis",
  "description": "Análise completa de documentos legais",
  "parameters": {
    "document_text": "string (obrigatório)",
    "analysis_type": "enum [risk, compliance, entities, classification]",
    "legal_area": "string (opcional)"
  }
}
```

#### **legal_document_generation**
```json
{
  "name": "legal_document_generation",
  "description": "Geração de documentos jurídicos",
  "parameters": {
    "document_type": "enum [contract, petition, opinion, motion]",
    "template_variables": "object - variáveis dinâmicas",
    "quality_level": "enum [draft, standard, professional, premium]",
    "customizations": "object - personalizações específicas"
  }
}
```

### **3. Search & Discovery Tools**

#### **advanced_search**
```json
{
  "name": "advanced_search",
  "description": "Busca avançada com filtros complexos",
  "parameters": {
    "query": "string (obrigatório)",
    "indices": "array de índices",
    "filters": "object - filtros complexos",
    "aggregations": "array - métricas desejadas",
    "sort_by": "string",
    "page_size": "integer"
  }
}
```

#### **search_suggestions**
```json
{
  "name": "search_suggestions",
  "description": "Sugestões automáticas de busca",
  "parameters": {
    "partial_query": "string (obrigatório)",
    "context": "enum [processes, jurisprudence, documents]",
    "max_suggestions": "integer (default: 5)"
  }
}
```

### **4. Notification & Communication Tools**

#### **notification_setup**
```json
{
  "name": "notification_setup",
  "description": "Configurar notificações personalizadas",
  "parameters": {
    "trigger_type": "enum [process_update, deadline, keyword_match]",
    "conditions": "object - condições específicas",
    "channels": "array [whatsapp, email, telegram, push]",
    "template_id": "string (opcional)",
    "priority": "enum [low, normal, high, critical]"
  }
}
```

#### **bulk_notification**
```json
{
  "name": "bulk_notification",
  "description": "Envio em massa de notificações",
  "parameters": {
    "recipient_filter": "object - critérios de destinatários",
    "message_template": "string",
    "variables": "object - variáveis do template",
    "schedule": "datetime (opcional)",
    "channels": "array"
  }
}
```

### **5. Reporting & Analytics Tools**

#### **generate_report**
```json
{
  "name": "generate_report",
  "description": "Geração de relatórios personalizados",
  "parameters": {
    "report_type": "enum [productivity, financial, processes, performance]",
    "date_range": "object {start, end}",
    "filters": "object - filtros específicos",
    "format": "enum [pdf, excel, json]",
    "include_charts": "boolean"
  }
}
```

#### **dashboard_metrics**
```json
{
  "name": "dashboard_metrics",
  "description": "Métricas para dashboard em tempo real",
  "parameters": {
    "metric_types": "array [processes_count, notifications_sent, search_volume]",
    "period": "enum [today, week, month, quarter, year]",
    "tenant_filter": "string (opcional)",
    "user_filter": "string (opcional)"
  }
}
```

### **6. Administrative Tools**

#### **user_management**
```json
{
  "name": "user_management",
  "description": "Gerenciamento de usuários e permissões",
  "parameters": {
    "action": "enum [list, create, update, deactivate]",
    "user_data": "object (para create/update)",
    "filters": "object (para list)",
    "permissions": "array (para create/update)"
  }
}
```

#### **tenant_analytics**
```json
{
  "name": "tenant_analytics",
  "description": "Análise de uso por escritório/tenant",
  "parameters": {
    "tenant_id": "string (opcional - admin only)",
    "metrics": "array [usage, quotas, performance, costs]",
    "period": "enum [day, week, month, quarter]",
    "export_format": "enum [json, csv, pdf]"
  }
}
```

### **7. Integration Tools**

#### **datajud_sync**
```json
{
  "name": "datajud_sync",
  "description": "Sincronização manual com DataJud CNJ",
  "parameters": {
    "process_numbers": "array de números CNJ",
    "force_update": "boolean",
    "priority": "enum [low, normal, high]",
    "callback_webhook": "string (opcional)"
  }
}
```

#### **external_api_call**
```json
{
  "name": "external_api_call",
  "description": "Chamadas para APIs externas autorizadas",
  "parameters": {
    "api_endpoint": "string",
    "method": "enum [GET, POST, PUT, DELETE]",
    "headers": "object",
    "body": "object (para POST/PUT)",
    "auth_required": "boolean"
  }
}
```

## 💬 Interfaces de Comunicação

### **1. WhatsApp Business API**
```
Usuário: "Oi, quero ver meus processos ativos"
Bot: "Olá! Encontrei 15 processos ativos. Gostaria de ver:
1. 📊 Resumo geral
2. 📋 Lista detalhada  
3. 🔍 Filtrar por cliente
4. ⚠️ Apenas urgentes"

Usuário: "Opção 2"
Bot: "📋 *Processos Ativos (15):*

🏛️ *TJ-SP*
• 1001234-56.2024.8.26.0100
  Cliente: João Silva
  Status: Aguardando decisão
  Última atualização: 2 dias atrás

• 2001234-78.2024.8.26.0200
  Cliente: Maria Santos
  Status: Em fase de instrução
  Última atualização: 5 dias atrás

[+ 13 outros processos]

Digite 'mais' para ver todos ou 'filtrar' para buscar específico."
```

### **2. Telegram Bot**
```
/start - Iniciar sessão
/processos - Listar processos
/buscar <termo> - Buscar jurisprudência  
/relatorio <tipo> - Gerar relatório
/notificar <processo> - Configurar alerta
/ajuda - Comandos disponíveis

Exemplo:
Usuário: /buscar responsabilidade civil médica
Bot: 🔍 Buscando jurisprudências sobre "responsabilidade civil médica"...

📚 Encontradas 47 decisões relevantes:

🏛️ STJ - REsp 1.234.567/SP (95% similar)
"A responsabilidade civil do médico..."
📅 Julgado em: 15/03/2024

🏛️ TJ-SP - Apelação 1001234-56.2023.8.26.0100 (92% similar) 
"Em casos de responsabilidade médica..."
📅 Julgado em: 08/02/2024

[Ver mais resultados] [Análise detalhada] [Salvar pesquisa]
```

### **3. Claude Chat Interface**
```
Usuário: "Preciso de um relatório completo sobre a produtividade do escritório no último trimestre, incluindo número de processos novos, conclusões, tempo médio de resolução e faturamento por advogado."

Claude: Vou gerar um relatório completo de produtividade trimestral para você. Deixe-me coletar os dados necessários:

🔍 Coletando dados de processos...
📊 Calculando métricas de produtividade...
💰 Compilando dados financeiros...

📈 **RELATÓRIO DE PRODUTIVIDADE - Q1 2024**

**VISÃO GERAL**
• Processos novos: 127 (+15% vs Q4 2023)
• Processos concluídos: 98 
• Tempo médio de resolução: 8.3 meses (-12% vs trimestre anterior)
• Taxa de sucesso: 87%

**POR ADVOGADO**
👨‍💼 Dr. João Silva
• Novos processos: 43
• Concluídos: 31
• Faturamento: R$ 145.000
• Tempo médio: 7.2 meses

👩‍💼 Dra. Maria Santos  
• Novos processos: 38
• Concluídos: 29
• Faturamento: R$ 132.000
• Tempo médio: 8.1 meses

[Continua com dados detalhados...]

Gostaria que eu gere algum gráfico específico ou analise algum aspecto em particular?
```

## 📊 Monetização e Planos

### **Estrutura de Acesso por Plano**

| Funcionalidade | Starter | Professional | Business | Enterprise |
|---------------|---------|--------------|----------|------------|
| **MCP Access** | ❌ | ✅ | ✅ | ✅ |
| **Comandos/mês** | - | 200 | 1.000 | Ilimitado |
| **WhatsApp Bot** | ❌ | ✅ | ✅ | ✅ |
| **Telegram Bot** | ❌ | ✅ | ✅ | ✅ |
| **Claude Interface** | ❌ | ❌ | ✅ | ✅ |
| **Slack Integration** | ❌ | ❌ | ❌ | ✅ |
| **Custom Tools** | ❌ | ❌ | ❌ | ✅ |
| **Advanced AI** | ❌ | Básico | Avançado | Premium |
| **Report Generation** | ❌ | ✅ | ✅ | ✅ |
| **Bulk Operations** | ❌ | ❌ | ✅ | ✅ |

### **Precificação de Features Premium**

```yaml
MCP_FEATURES:
  basic_commands: Incluído nos planos
  ai_analysis: +R$ 0.50 por análise
  document_generation: +R$ 2.00 por documento
  advanced_reports: +R$ 5.00 por relatório
  bulk_operations: +R$ 0.10 por item processado
  external_integrations: +R$ 10.00/mês por integração
```

## 🚀 Roadmap de Implementação

### **Fase 1: Foundation (Semana 1)**
- ✅ Criar estrutura do MCP Service
- ✅ Implementar autenticação e multi-tenancy
- ✅ Configurar conexão com Claude API
- ✅ Desenvolver sistema de ferramentas base
- ✅ Criar middleware de logging e métricas

### **Fase 2: Core Tools (Semana 2)**
- ✅ Implementar process_search e process_monitor
- ✅ Desenvolver notification_setup e bulk_notification  
- ✅ Criar advanced_search e search_suggestions
- ✅ Implementar generate_report básico
- ✅ Testes unitários das ferramentas core

### **Fase 3: AI Integration (Semana 3)**
- ✅ Integrar jurisprudence_search com AI Service
- ✅ Implementar case_similarity_analysis
- ✅ Desenvolver document_analysis
- ✅ Criar legal_document_generation
- ✅ Otimizar performance com cache Redis

### **Fase 4: Bot Interfaces (Semana 4)**
- ✅ WhatsApp Business API integration
- ✅ Telegram Bot implementation
- ✅ Claude Chat interface setup
- ✅ Context management e session handling
- ✅ Testes de integração completos

### **Fase 5: Advanced Features (Semana 5)**
- 🔄 Slack Bot integration
- 🔄 Custom tools framework
- 🔄 Advanced analytics e reporting
- 🔄 Bulk operations optimization
- 🔄 External API integrations

### **Fase 6: Production Ready (Semana 6)**
- 🔄 Security hardening
- 🔄 Performance optimization
- 🔄 Monitoring e alerting
- 🔄 Documentation completa
- 🔄 Deploy em produção

## 🔒 Segurança e Compliance

### **Autenticação e Autorização**
- **JWT Integration**: Reutilização do sistema existente
- **Multi-tenant Isolation**: Garantia de isolamento de dados
- **Role-based Access**: Controle granular por função
- **Rate Limiting**: Prevenção de abuso por usuário/tenant
- **Audit Logging**: Rastreamento completo de ações

### **Proteção de Dados**
- **Encryption in Transit**: TLS 1.3 para todas as comunicações
- **Encryption at Rest**: Dados sensíveis criptografados
- **Data Anonymization**: Logs sem informações pessoais
- **LGPD Compliance**: Conformidade com legislação brasileira
- **Retention Policies**: Políticas de retenção de dados

### **API Security**
- **Input Validation**: Validação rigorosa de todos os inputs
- **Output Sanitization**: Sanitização de respostas
- **SQL Injection Prevention**: Queries parametrizadas
- **XSS Protection**: Headers de segurança
- **CORS Policy**: Política restritiva de CORS

## 📈 Métricas e KPIs

### **Métricas de Uso**
- **Daily Active Users (DAU)** via MCP
- **Commands per User** (média diária)
- **Tool Usage Distribution** (ferramentas mais usadas)
- **Response Time** (latência média)
- **Error Rate** (taxa de erro por ferramenta)

### **Métricas de Negócio**
- **User Engagement** (tempo de sessão médio)
- **Feature Adoption** (% usuários usando MCP)
- **Churn Reduction** (retenção com MCP vs sem MCP)
- **Revenue Impact** (receita adicional por features MCP)
- **Support Reduction** (redução em tickets de suporte)

### **Métricas Técnicas**
- **API Response Time** (99th percentile)
- **Claude API Usage** (tokens consumidos)
- **Cache Hit Rate** (eficiência do cache)
- **Error Distribution** (tipos de erro mais comuns)
- **Resource Utilization** (CPU, memória, rede)

## 🛡️ Disaster Recovery e Backup

### **Estratégia de Backup**
- **Context State**: Backup contínuo do estado das sessões
- **Tool Registry**: Versionamento das configurações de ferramentas
- **User Preferences**: Backup das preferências de usuário
- **Audit Logs**: Retenção de 5 anos para compliance

### **Failover Strategy**
- **Claude API Fallback**: Sistema de fallback para APIs Claude
- **Multi-region Deployment**: Deploy em múltiplas regiões
- **Graceful Degradation**: Funcionalidade limitada em caso de falha
- **Health Monitoring**: Monitoramento proativo da saúde do sistema

## 🔧 Configuração e Environment

### **Environment Variables**
```bash
# MCP Service Configuration
MCP_SERVICE_NAME=mcp-service
MCP_VERSION=1.0.0
MCP_PORT=8088

# Claude API Configuration  
ANTHROPIC_API_KEY=your_claude_api_key
ANTHROPIC_MODEL=claude-3-5-sonnet-20241022
ANTHROPIC_MAX_TOKENS=4096
ANTHROPIC_TIMEOUT=30s

# Bot Integrations
WHATSAPP_ACCESS_TOKEN=your_whatsapp_token
WHATSAPP_PHONE_NUMBER_ID=your_phone_number_id
WHATSAPP_VERIFY_TOKEN=your_verify_token

TELEGRAM_BOT_TOKEN=your_telegram_token
TELEGRAM_WEBHOOK_URL=https://your-domain.com/webhook/telegram

SLACK_BOT_TOKEN=your_slack_token
SLACK_SIGNING_SECRET=your_slack_secret

# Performance & Caching
MCP_CACHE_TTL=300s
MCP_MAX_CONCURRENT_REQUESTS=100
MCP_RATE_LIMIT_PER_USER=60
MCP_SESSION_TIMEOUT=30m

# Security
MCP_JWT_SECRET=your_jwt_secret
MCP_ENCRYPTION_KEY=your_encryption_key
MCP_ALLOWED_ORIGINS=https://app.direitolux.com

# Monitoring
MCP_METRICS_ENABLED=true
MCP_LOGGING_LEVEL=info
MCP_TRACING_ENABLED=true
```

## 📚 Recursos e Referências

### **Documentação Técnica**
- [Model Context Protocol Specification](https://spec.modelcontextprotocol.io/)
- [Anthropic Claude API Documentation](https://docs.anthropic.com/)
- [WhatsApp Business API Documentation](https://developers.facebook.com/docs/whatsapp)
- [Telegram Bot API Documentation](https://core.telegram.org/bots/api)

### **Best Practices**
- [MCP Tool Development Guide](https://modelcontextprotocol.io/docs/tools)
- [Conversational AI Design Principles](https://claude.ai/docs/conversational-ai)
- [Multi-tenant Architecture Patterns](https://docs.microsoft.com/en-us/azure/architecture/patterns/)

---

**🔄 Última Atualização**: 18/06/2025  
**📈 Status**: 📋 Planejamento Completo - Pronto para Implementação  
**👨‍💻 Responsável**: Full Cycle Developer  
**🎯 Próximo**: Implementação do MCP Service Foundation