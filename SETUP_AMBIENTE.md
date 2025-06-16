# Setup do Ambiente - Direito Lux

## 📋 Pré-requisitos

### Sistema Operacional
- macOS, Linux ou Windows com WSL2
- Mínimo 16GB RAM (recomendado 32GB)
- 50GB de espaço em disco livre

### Software Necessário
- **Docker Desktop** 4.0+ com Docker Compose
- **Go** 1.21+
- **Node.js** 18+ e npm
- **Python** 3.11+
- **Git** 2.30+
- **Make** (geralmente já instalado)
- **kubectl** (para Kubernetes)
- **Terraform** 1.5+ (para IaC)
- **Google Cloud SDK** (gcloud)

### Ferramentas de Desenvolvimento
```bash
# Go tools
go install github.com/cosmtrek/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Node tools
npm install -g @angular/cli
npm install -g expo-cli

# Python tools
pip install poetry
pip install pre-commit
```

## 🚀 Setup Inicial

### 1. Clonar o Repositório
```bash
git clone https://github.com/direito-lux/direito-lux.git
cd direito-lux
```

### 2. Configurar Variáveis de Ambiente
```bash
# Copiar arquivo de exemplo
cp .env.example .env

# Editar com suas configurações
# IMPORTANTE: Gerar secrets seguros para produção
nano .env
```

### 3. Configurar Git Hooks
```bash
# Instalar pre-commit hooks
pre-commit install

# Configurar git
git config user.name "Seu Nome"
git config user.email "seu@email.com"
```

## 🐳 Docker Environment

### 1. Build das Imagens Base
```bash
# Build de todas as imagens
docker-compose build

# Ou build específico
docker-compose build auth-service
```

### 2. Iniciar Infraestrutura Base
```bash
# Iniciar serviços de infraestrutura primeiro
docker-compose up -d postgres redis rabbitmq

# Aguardar health checks
docker-compose ps
# Todos devem estar "healthy"
```

### 3. Executar Migrações
```bash
# Criar databases
./scripts/setup-postgres.sh

# Rodar migrações do Auth Service
cd services/auth-service
migrate -path migrations -database "postgres://direito_lux:dev_password_123@localhost:5432/direito_lux_dev?sslmode=disable" up
```

### 4. Popular Dados de Desenvolvimento
```bash
# Executar script de seed
./scripts/seed-data.sh

# Isso cria:
# - Tenants de exemplo
# - Usuários de teste
# - Processos mock
```

### 5. Iniciar Todos os Serviços
```bash
# Iniciar tudo
docker-compose up -d

# Verificar logs
docker-compose logs -f

# Verificar status
docker-compose ps
```

## 🔧 Desenvolvimento Local

### Auth Service
```bash
cd services/auth-service

# Copiar configurações
cp .env.example .env

# Instalar dependências
go mod download

# Executar com hot reload
air

# Ou executar diretamente
go run cmd/server/main.go

# Executar testes
go test ./...

# Coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Criar Novo Serviço
```bash
# Usar script de criação
./scripts/create-service.sh tenant-service

# Navegar para o serviço
cd services/tenant-service

# Configurar e executar
cp .env.example .env
go mod tidy
make dev
```

## 🌐 URLs e Acessos

### Aplicação
- **API Gateway**: http://localhost:8000
- **Auth Service**: http://localhost:8081
- **Tenant Service**: http://localhost:8082
- **Process Service**: http://localhost:8083

### Infraestrutura
- **PostgreSQL**: localhost:5432
  - User: `direito_lux`
  - Password: `dev_password_123`
  - Database: `direito_lux_dev`

- **Redis**: localhost:6379
  - Password: `dev_redis_123`

- **RabbitMQ**: http://localhost:15672
  - User: `direito_lux`
  - Password: `dev_rabbit_123`

### Observabilidade
- **Jaeger UI**: http://localhost:16686
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000
  - User: `admin`
  - Password: `admin123`
- **Kibana**: http://localhost:5601

### Ferramentas
- **Keycloak**: http://localhost:8080
  - Admin: `admin`
  - Password: `admin123`
- **MinIO**: http://localhost:9001
  - Access Key: `minioadmin`
  - Secret Key: `minioadmin`
- **Mailhog**: http://localhost:8025

## 🧪 Testando a Aplicação

### 1. Health Checks
```bash
# Auth Service
curl http://localhost:8081/health
curl http://localhost:8081/ready

# Metrics
curl http://localhost:9090/metrics
```

### 2. Autenticação
```bash
# Login
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: 00000000-0000-0000-0000-000000000001" \
  -d '{
    "email": "admin@example.com",
    "password": "Admin@123"
  }'

# Usar o token retornado para próximas requisições
export TOKEN="seu-jwt-token-aqui"

# Validar token
curl http://localhost:8081/api/v1/auth/validate \
  -H "Authorization: Bearer $TOKEN"
```

### 3. Criar Usuário
```bash
curl -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: 00000000-0000-0000-0000-000000000001" \
  -d '{
    "email": "novo@example.com",
    "password": "Senha@123",
    "first_name": "Novo",
    "last_name": "Usuário",
    "role": "client"
  }'
```

## 🔍 Troubleshooting

### Problemas Comuns

#### 1. Portas em Uso
```bash
# Verificar portas
lsof -i :8081
lsof -i :5432

# Matar processo
kill -9 <PID>
```

#### 2. Docker sem Espaço
```bash
# Limpar Docker
docker system prune -a
docker volume prune
```

#### 3. Problemas de Permissão
```bash
# macOS/Linux
sudo chown -R $(whoami) .

# Permissões do Docker
sudo usermod -aG docker $USER
```

#### 4. Health Check Falhando
```bash
# Verificar logs específicos
docker-compose logs postgres
docker-compose logs auth-service

# Reiniciar serviço
docker-compose restart postgres
```

### Logs e Debug

#### Verificar Logs
```bash
# Todos os serviços
docker-compose logs -f

# Serviço específico
docker-compose logs -f auth-service

# Últimas 100 linhas
docker-compose logs --tail=100 auth-service
```

#### Debug com Delve
```bash
# Instalar Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug do Auth Service
cd services/auth-service
dlv debug cmd/server/main.go
```

#### Jaeger Tracing
1. Abrir http://localhost:16686
2. Selecionar serviço "auth-service"
3. Buscar traces por operação

## 🛠️ Comandos Úteis

### Docker Compose
```bash
# Start/Stop
docker-compose up -d
docker-compose down

# Rebuild
docker-compose build --no-cache auth-service
docker-compose up -d --force-recreate auth-service

# Executar comando em container
docker-compose exec auth-service sh

# Ver recursos
docker stats
```

### Makefile Commands
```bash
# No diretório do serviço
make help       # Ver comandos disponíveis
make build      # Build do binário
make test       # Executar testes
make lint       # Executar linter
make docker-build  # Build Docker
```

### Database
```bash
# Conectar ao PostgreSQL
docker-compose exec postgres psql -U direito_lux -d direito_lux_dev

# Backup
docker-compose exec postgres pg_dump -U direito_lux direito_lux_dev > backup.sql

# Restore
docker-compose exec -T postgres psql -U direito_lux direito_lux_dev < backup.sql
```

## 🚀 Próximos Passos

1. **Configurar IDE** (VSCode/GoLand)
   - Instalar extensões Go
   - Configurar debugger
   - Setup linters

2. **Familiarizar com o código**
   - Explorar template-service
   - Entender estrutura hexagonal
   - Revisar auth-service

3. **Começar desenvolvimento**
   - Criar branch feature
   - Implementar Tenant Service
   - Escrever testes

## 📚 Documentação Adicional

- [Arquitetura Full Cycle](./ARQUITETURA_FULLCYCLE.md)
- [Event Storming](./EVENT_STORMING_DIREITO_LUX.md)
- [Roadmap](./ROADMAP_IMPLEMENTACAO.md)
- [Status da Implementação](./STATUS_IMPLEMENTACAO.md)

## 🆘 Suporte

Em caso de problemas:
1. Verificar logs detalhados
2. Consultar documentação
3. Abrir issue no GitHub
4. Contatar equipe de desenvolvimento