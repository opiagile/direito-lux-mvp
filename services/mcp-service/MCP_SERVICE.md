# MCP Service - Model Context Protocol

> **📋 Este arquivo foi movido**: A documentação completa do MCP Service está disponível em:
> 
> - [README.md](./README.md) - Documentação principal do MCP Service
> - [README-INTEGRATION.md](./README-INTEGRATION.md) - Guia de integração e deploy
> 
> Por favor, consulte esses arquivos para informações atualizadas sobre o MCP Service.

## 🤖 Resumo Executivo

O **MCP Service** é o diferencial estratégico do Direito Lux - o primeiro SaaS jurídico brasileiro com interface conversacional completa via bots inteligentes.

### 🏆 Características Principais

- **17+ ferramentas MCP** especializadas para advogados
- **Multi-bot**: WhatsApp, Telegram, Claude Chat, Slack
- **Claude 3.5 Sonnet** como modelo de IA principal
- **Sistema de quotas** por plano (200/1000/ilimitado comandos/mês)
- **Arquitetura hexagonal** em Go 1.21+
- **Infraestrutura completa** com PostgreSQL, Redis, RabbitMQ

### 📊 Status de Implementação

- ✅ **Domain layer**: 17+ ferramentas especificadas
- ✅ **Infrastructure layer**: Completa (config, database, events, HTTP, messaging)
- ✅ **Application layer**: Handlers para sessões, ferramentas e bots
- ✅ **Compilação**: Testada e funcionando
- ✅ **Deploy DEV**: Configurado com infraestrutura separada
- ✅ **Documentação**: Completa e atualizada

### 🚀 Como Usar

Consulte [README-INTEGRATION.md](./README-INTEGRATION.md) para instruções detalhadas de deploy e integração.

```bash
# Deploy rápido
cd services/mcp-service
./scripts/start-dev.sh --clean --build
```

### 📚 Documentação Completa

- [README.md](./README.md) - Documentação técnica completa
- [README-INTEGRATION.md](./README-INTEGRATION.md) - Integração e deploy