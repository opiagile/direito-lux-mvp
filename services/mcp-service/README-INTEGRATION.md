# MCP Service - Integração Real

## 🚀 Ambiente de Desenvolvimento

### Pré-requisitos

- Docker e Docker Compose
- Go 1.21+
- Git

### Configuração Rápida

```bash
# 1. Clonar e navegar para o diretório
cd services/mcp-service

# 2. Executar script de desenvolvimento
chmod +x scripts/start-dev.sh
./scripts/start-dev.sh

# Ou com limpeza de volumes
./scripts/start-dev.sh --clean

# Ou apenas infraestrutura (sem iniciar o serviço)
./scripts/start-dev.sh --clean --migrate --no-start
```

### Serviços Disponíveis

| Serviço | Endereço | Credenciais |
|---------|----------|-------------|
| **PostgreSQL** | localhost:5434 | mcp_user / mcp_pass_dev |
| **Redis** | localhost:6380 | redis_pass_dev |
| **RabbitMQ** | localhost:5673 | mcp_user / rabbit_pass_dev |
| **RabbitMQ Management** | http://localhost:15673 | mcp_user / rabbit_pass_dev |
| **Jaeger UI** | http://localhost:16687 | - |
| **MCP Service** | http://localhost:8084 | - |
| **Metrics (Prometheus)** | http://localhost:9084/metrics | - |

## 🗄️ Banco de Dados

### Schema

O banco `direito_lux_mcp` contém as seguintes tabelas:

- `tenants` - Informações dos tenants
- `users` - Usuários do sistema
- `mcp_sessions` - Sessões de conversação MCP
- `mcp_messages` - Histórico de mensagens
- `mcp_tool_executions` - Log de execução de ferramentas
- `mcp_events` - Eventos de domínio
- `mcp_quotas` - Controle de cotas por tenant
- `mcp_bot_configs` - Configurações dos bots

### Dados de Teste

O banco é inicializado com:
- 2 tenants de exemplo (premium e basic)
- 2 usuários de teste
- Configurações de quotas iniciais

### Conexão Manual

```bash
# PostgreSQL
psql -h localhost -p 5434 -U mcp_user -d direito_lux_mcp

# Redis
redis-cli -h localhost -p 6380 -a redis_pass_dev
```

## 🔧 Configuração

### Variáveis de Ambiente

O arquivo `.env.development` contém todas as configurações necessárias:

```bash
# Application
APP_NAME=mcp-service
APP_ENV=development
PORT=8084
LOG_LEVEL=debug

# Database
DB_HOST=localhost
DB_PORT=5434
DB_NAME=direito_lux_mcp
DB_USER=mcp_user
DB_PASSWORD=mcp_pass_dev

# Redis
REDIS_HOST=localhost
REDIS_PORT=6380
REDIS_PASSWORD=redis_pass_dev

# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5673
RABBITMQ_USER=mcp_user
RABBITMQ_PASSWORD=rabbit_pass_dev
RABBITMQ_VHOST=mcp_vhost

# Claude API
CLAUDE_API_KEY=sk-ant-api03-your-key-here
```

### Personalização

Para usar suas próprias configurações:

1. Copie `.env.development` para `.env.local`
2. Modifique os valores necessários
3. Execute: `source .env.local && ./scripts/start-dev.sh`

## 🧪 Testes

### Teste de Compilação

```bash
# Compilar projeto
go build -o bin/mcp-service ./cmd/main.go

# Executar testes unitários
go test ./...

# Teste de integração
go run cmd/main.go
```

### Teste de Conexões

```bash
# Testar PostgreSQL
go run cmd/main.go 2>&1 | grep "PostgreSQL conectado"

# Testar APIs
curl http://localhost:8084/health
curl http://localhost:8084/api/v1/sessions
```

## 📊 Monitoramento

### Métricas

- **Prometheus**: http://localhost:9084/metrics
- **Health Check**: http://localhost:8084/health

### Logs

```bash
# Logs do MCP Service
docker-compose -f docker-compose.dev.yml logs -f

# Logs específicos
docker-compose -f docker-compose.dev.yml logs postgres
docker-compose -f docker-compose.dev.yml logs redis
docker-compose -f docker-compose.dev.yml logs rabbitmq
```

### Tracing

- **Jaeger UI**: http://localhost:16687
- **Service Name**: mcp-service

## 🔄 Workflow de Desenvolvimento

### 1. Inicialização

```bash
# Primeira vez
./scripts/start-dev.sh --clean --migrate

# Desenvolvimento contínuo
./scripts/start-dev.sh
```

### 2. Desenvolvimento

```bash
# Hot reload com air (opcional)
go install github.com/cosmtrek/air@latest
air

# Ou execução manual
go run cmd/main.go
```

### 3. Testes

```bash
# Testes unitários
go test ./internal/...

# Testes de integração
go test ./tests/integration/...

# Benchmark
go test -bench=. ./...
```

### 4. Debug

```bash
# Logs detalhados
LOG_LEVEL=debug go run cmd/main.go

# Profiling
go tool pprof http://localhost:8084/debug/pprof/profile
```

## 🚨 Troubleshooting

### Problemas Comuns

1. **Porta já em uso**
   ```bash
   # Parar containers
   docker-compose -f docker-compose.dev.yml down
   
   # Verificar processos
   lsof -i :5434  # PostgreSQL
   lsof -i :6380  # Redis
   lsof -i :5673  # RabbitMQ
   ```

2. **Banco não conecta**
   ```bash
   # Verificar se container está rodando
   docker-compose -f docker-compose.dev.yml ps
   
   # Reiniciar PostgreSQL
   docker-compose -f docker-compose.dev.yml restart postgres
   ```

3. **Volumes corrompidos**
   ```bash
   # Limpar tudo e reiniciar
   ./scripts/start-dev.sh --clean
   ```

### Logs de Debug

```bash
# Ver logs de inicialização
go run cmd/main.go 2>&1 | tee startup.log

# Monitorar em tempo real
tail -f startup.log | grep "ERROR\|WARN"
```

## 🔒 Segurança

### Desenvolvimento vs Produção

- **Desenvolvimento**: Credenciais fixas para facilitar setup
- **Produção**: Usar secrets/vault para credenciais sensíveis

### Variáveis Sensíveis

Nunca commitar no git:
- CLAUDE_API_KEY
- DB_PASSWORD (produção)
- JWT_SECRET (produção)
- Tokens de bots (produção)

## 📚 Próximos Passos

1. ✅ Integração com PostgreSQL, Redis, RabbitMQ
2. 🔄 Deploy para ambiente DEV
3. 🔄 Implementação completa dos handlers
4. 🔄 Testes de carga e performance
5. 🔄 Documentação da API (Swagger)
6. 🔄 CI/CD pipeline