# Direito Lux - Development Deployment

## 🚀 Deploy Completo do Ambiente de Desenvolvimento

Este deployment inclui todos os serviços principais do Direito Lux:

- **AI Service** (Python/FastAPI) - Análise de documentos e IA
- **Search Service** (Go) - Busca avançada com Elasticsearch  
- **MCP Service** (Go) - Model Context Protocol para bots
- **Infraestrutura completa** - PostgreSQL, Redis, RabbitMQ, Elasticsearch, Jaeger

## 📋 Pré-requisitos

- Docker Desktop ou Docker Engine + Docker Compose
- 8GB+ RAM disponível
- 10GB+ espaço em disco

## 🎯 Quick Start

```bash
# 1. Navegar para o diretório
cd services/

# 2. Dar permissão de execução ao script
chmod +x scripts/deploy-dev.sh

# 3. Deploy completo (primeira vez)
./scripts/deploy-dev.sh --clean --build

# 4. Deploy normal (dias seguintes)
./scripts/deploy-dev.sh
```

## 🛠️ Comandos Disponíveis

### Comandos Principais

```bash
# Iniciar todos os serviços
./scripts/deploy-dev.sh start

# Parar todos os serviços
./scripts/deploy-dev.sh stop

# Reiniciar todos os serviços
./scripts/deploy-dev.sh restart

# Ver status dos serviços
./scripts/deploy-dev.sh status

# Ver endpoints disponíveis
./scripts/deploy-dev.sh endpoints

# Executar testes de conectividade
./scripts/deploy-dev.sh test
```

### Logs e Monitoramento

```bash
# Ver logs de todos os serviços
./scripts/deploy-dev.sh logs

# Ver logs de um serviço específico
./scripts/deploy-dev.sh logs ai-service
./scripts/deploy-dev.sh logs search-service

# Seguir logs em tempo real
docker-compose -f docker-compose.dev.yml logs -f ai-service
```

### Opções Avançadas

```bash
# Limpar volumes e reiniciar do zero
./scripts/deploy-dev.sh --clean start

# Reconstruir imagens Docker
./scripts/deploy-dev.sh --build start

# Baixar imagens mais recentes
./scripts/deploy-dev.sh --pull start

# Combinar opções
./scripts/deploy-dev.sh --clean --build --pull start
```

## 🌐 Endpoints Disponíveis

### Serviços Principais

| Serviço | URL | Descrição |
|---------|-----|-----------|
| **AI Service** | http://localhost:8000 | Análise de documentos e IA |
| **Search Service** | http://localhost:8086 | Busca avançada |
| **AI Health** | http://localhost:8000/health | Health check AI |
| **Search Health** | http://localhost:8086/health | Health check Search |
| **AI Docs** | http://localhost:8000/docs | Swagger UI - AI Service |
| **Search Docs** | http://localhost:8086/docs | Swagger UI - Search Service |

### Infraestrutura

| Serviço | Endpoint | Credenciais |
|---------|----------|-------------|
| **PostgreSQL (Main)** | localhost:5432 | direito_lux / direito_lux_pass_dev |
| **PostgreSQL (MCP)** | localhost:5434 | mcp_user / mcp_pass_dev |
| **Redis (Main)** | localhost:6379 | redis_pass_dev |
| **Redis (MCP)** | localhost:6380 | redis_pass_dev |
| **RabbitMQ (Main)** | localhost:5672 | direito_lux / rabbit_pass_dev |
| **RabbitMQ (MCP)** | localhost:5673 | mcp_user / rabbit_pass_dev |
| **Elasticsearch** | http://localhost:9200 | - |

### Monitoramento

| Serviço | URL | Credenciais |
|---------|-----|-------------|
| **RabbitMQ Management** | http://localhost:15672 | direito_lux / rabbit_pass_dev |
| **RabbitMQ Mgmt (MCP)** | http://localhost:15673 | mcp_user / rabbit_pass_dev |
| **Jaeger Tracing** | http://localhost:16686 | - |

## 🗄️ Banco de Dados

### Esquemas Criados

- `public` - Tabelas comuns (tenants, users)
- `ai_service` - Análises e embeddings
- `search_service` - Índices e consultas
- `auth_service` - Autenticação
- `process_service` - Processos jurídicos
- `tenant_service` - Gestão de tenants
- `notification_service` - Notificações

### Dados de Desenvolvimento

```sql
-- Tenants de teste
'11111111-1111-1111-1111-111111111111' -> 'Tenant Dev' (premium)
'22222222-2222-2222-2222-222222222222' -> 'Tenant Test' (basic)

-- Usuários de teste
dev@direito-lux.com (admin)
test@direito-lux.com (user)
```

### Conectar Manualmente

```bash
# PostgreSQL Principal
docker exec -it services_postgres_1 psql -U direito_lux -d direito_lux_dev

# PostgreSQL MCP
docker exec -it services_mcp-postgres_1 psql -U mcp_user -d direito_lux_mcp

# Redis Principal
docker exec -it services_redis_1 redis-cli -a redis_pass_dev

# Redis MCP  
docker exec -it services_mcp-redis_1 redis-cli -a redis_pass_dev
```

## 🧪 Testes e Validação

### Teste Rápido dos Serviços

```bash
# AI Service
curl http://localhost:8000/health

# Search Service
curl http://localhost:8086/health

# Elasticsearch
curl http://localhost:9200/_health

# PostgreSQL (através do AI Service)
curl http://localhost:8000/api/v1/analysis/
```

### Teste de Funcionalidades

```bash
# 1. Análise de texto (AI Service)
curl -X POST "http://localhost:8000/api/v1/analysis/document" \
  -H "Content-Type: application/json" \
  -d '{"text": "Este é um contrato de prestação de serviços", "analysis_type": "classification"}'

# 2. Busca (Search Service)
curl -X POST "http://localhost:8086/api/v1/search" \
  -H "Content-Type: application/json" \
  -d '{"query": "contrato", "tenant_id": "11111111-1111-1111-1111-111111111111"}'

# 3. Health checks automatizados
./scripts/deploy-dev.sh test
```

## 🔧 Desenvolvimento

### Hot Reload

Os serviços estão configurados para hot reload durante desenvolvimento:

- **AI Service**: Mudanças no código Python são detectadas automaticamente
- **Search Service**: Precisa ser recompilado (`docker-compose restart search-service`)

### Logs de Debug

```bash
# Logs detalhados de um serviço
docker-compose -f docker-compose.dev.yml logs -f ai-service

# Logs com filtro
docker-compose -f docker-compose.dev.yml logs ai-service | grep ERROR

# Logs de múltiplos serviços
docker-compose -f docker-compose.dev.yml logs -f ai-service search-service
```

### Modificar Configurações

1. Editar arquivo `.env.development` do serviço
2. Reiniciar o serviço:
   ```bash
   docker-compose -f docker-compose.dev.yml restart ai-service
   ```

## 🚨 Troubleshooting

### Problemas Comuns

#### 1. Portas em Uso
```bash
# Verificar o que está usando a porta
lsof -i :8000  # AI Service
lsof -i :8086  # Search Service
lsof -i :5432  # PostgreSQL

# Parar todos os containers
./scripts/deploy-dev.sh stop
```

#### 2. Problemas de Memória
```bash
# Verificar uso de recursos
docker stats

# Limpar containers não utilizados
docker system prune -a

# Reiniciar com limpeza
./scripts/deploy-dev.sh --clean start
```

#### 3. Elasticsearch Não Inicia
```bash
# Aumentar vm.max_map_count (Linux/Mac)
sudo sysctl -w vm.max_map_count=262144

# Windows WSL
wsl -d docker-desktop sysctl -w vm.max_map_count=262144
```

#### 4. Serviços Não Respondem
```bash
# Verificar logs de erro
./scripts/deploy-dev.sh logs | grep -i error

# Verificar health checks
docker-compose -f docker-compose.dev.yml ps

# Reiniciar serviço específico
docker-compose -f docker-compose.dev.yml restart ai-service
```

### Reset Completo

```bash
# 1. Parar tudo
./scripts/deploy-dev.sh stop

# 2. Limpar volumes e imagens
docker system prune -a --volumes

# 3. Reiniciar do zero
./scripts/deploy-dev.sh --clean --build start
```

## 📊 Monitoramento

### Métricas Disponíveis

- **Jaeger**: Tracing distribuído em http://localhost:16686
- **RabbitMQ**: Filas e mensagens em http://localhost:15672
- **Elasticsearch**: Status do cluster em http://localhost:9200/_cluster/health

### Logs Estruturados

Todos os serviços produzem logs estruturados em JSON:

```bash
# Ver logs estruturados
docker-compose -f docker-compose.dev.yml logs ai-service | jq .

# Filtrar por nível
docker-compose -f docker-compose.dev.yml logs ai-service | jq 'select(.level == "ERROR")'
```

## 🔒 Segurança (Desenvolvimento)

> ⚠️ **ATENÇÃO**: Este setup é apenas para desenvolvimento. Não usar em produção!

- Senhas fixas e simples
- Sem TLS/SSL
- Sem autenticação em alguns serviços
- Dados não persistentes entre resets

## 📚 Próximos Passos

1. ✅ Deploy de desenvolvimento funcionando
2. 🔄 Implementar testes de integração automatizados
3. 🔄 Configurar CI/CD pipeline
4. 🔄 Setup de ambiente de staging
5. 🔄 Configurações de produção
6. 🔄 Monitoramento avançado (Prometheus + Grafana)

## 🆘 Suporte

Em caso de problemas:

1. Verificar logs: `./scripts/deploy-dev.sh logs`
2. Verificar status: `./scripts/deploy-dev.sh status`  
3. Tentar reset: `./scripts/deploy-dev.sh --clean start`
4. Verificar documentação do serviço específico