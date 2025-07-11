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

## 🚀 Setup Local (Desenvolvimento)

### 1. Clonar o Repositório
```bash
git clone https://github.com/direito-lux/direito-lux.git
cd direito-lux
```

### 1.1. Setup Automatizado (Recomendado - 100% FUNCIONAL! ✨)
```bash
# Setup definitivo com auth service 100% funcional
chmod +x SETUP_DATABASE_DEFINITIVO.sh
./SETUP_DATABASE_DEFINITIVO.sh

# ✅ Isso irá:
# - Configurar database com schema corrigido
# - Criar 32 usuários de teste funcionais
# - Inicializar auth service na porta 8081
# - Validar login JWT funcionando

# Verificar se funcionou
./scripts/utilities/CHECK_SERVICES_STATUS.sh
```

### 1.2. Scripts Essenciais (Ambiente Limpo - Redução de 75%)

Após a **grande limpeza**, mantemos apenas os scripts essenciais:

```bash
# ⭐ CONFIGURAÇÃO INICIAL
./SETUP_DATABASE_DEFINITIVO.sh               # Setup definitivo com auth 100% funcional
./CLEAN_ENVIRONMENT_TOTAL.sh                 # Limpeza total quando necessário

# 🛠️ DESENVOLVIMENTO DIÁRIO  
./START_LOCAL_DEV.sh                         # Iniciar ambiente de desenvolvimento
./scripts/utilities/CHECK_SERVICES_STATUS.sh # Verificar status dos serviços
./test-local.sh                              # Testar funcionalidades
./stop-services.sh                           # Parar serviços

# 📦 BUILD E DEPLOY
./build-all.sh                               # Compilar todos os microserviços
./start-services.sh                          # Iniciar serviços localmente
./create-service.sh                          # Criar novo microserviço
```

📋 **Consulte** [`SCRIPTS_ESSENCIAIS.md`](./SCRIPTS_ESSENCIAIS.md) **para documentação completa dos 17 scripts organizados**

### 1.3. Setup Frontend
```bash
cd frontend
npm install
npm run dev
# Frontend: http://localhost:3000
# Grafana: http://localhost:3002

# ✅ NOVO: Páginas de Autenticação Completas (08/01/2025)
# - /register - Registro público em 3 etapas (tenant → admin → plano)
# - /forgot-password - Recuperação de senha
# - /reset-password - Reset de senha com validação
# - /login - Login existente atualizado
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

### 3. Executar Migrações (Automatizado)
```bash
# As migrações são executadas automaticamente pelo SETUP_COMPLETE_FIXED.sh
# Para execução manual:
./scripts/utilities/execute_migrations.sh

# Ou manualmente:
cd services/auth-service
migrate -path migrations -database "postgres://direito_lux:direito_lux_pass_dev@localhost:5432/direito_lux_dev?sslmode=disable" up

# ✅ NOVA: Migração de Password Reset Tokens (08/01/2025)
# 004_create_password_reset_tokens_table.sql - Para recuperação de senha
# Criada automaticamente com o sistema completo de autenticação
```

### 4. Popular Dados de Desenvolvimento (Automatizado)
```bash
# Os dados são inseridos automaticamente pelo SETUP_COMPLETE_FIXED.sh
# Dados incluídos: SEED_DATABASE_COMPLETE.sql
# - 8 tenants (2 por plano)
# - 32 usuários (4 por tenant)
# - 90+ processos de exemplo
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
- **DataJud Service**: http://localhost:8084
- **Report Service**: http://localhost:8087
- **🆕 Billing Service**: http://localhost:8089
- **Frontend Web App**: http://localhost:3000

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
- **Grafana**: http://localhost:3002 (admin / dev_grafana_123)
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

## 🏗️ Setup Produção (GCP + Kubernetes)

### Pré-requisitos para Produção
```bash
# Instalar ferramentas necessárias
brew install google-cloud-sdk
brew install kubernetes-cli
brew install terraform

# Autenticar com GCP
gcloud auth login
gcloud auth application-default login

# Configurar projeto
gcloud config set project direito-lux-production
```

### 1. Infraestrutura (Terraform)

#### 1.1. Deploy Infraestrutura Staging
```bash
cd terraform

# Tornar script executável
chmod +x deploy.sh

# Inicializar
./deploy.sh staging init

# Planejar
./deploy.sh staging plan

# Aplicar
./deploy.sh staging apply
```

#### 1.2. Deploy Infraestrutura Production
```bash
# Validar staging primeiro, depois
./deploy.sh production apply
```

#### 1.3. Verificar Infraestrutura
```bash
# Ver outputs
./deploy.sh production output

# Verificar no console GCP
gcloud compute instances list
gcloud container clusters list
gcloud sql instances list
```

### 2. Aplicações (Kubernetes)

#### 2.1. Deploy Staging
```bash
cd k8s

# Tornar script executável
chmod +x deploy.sh

# Deploy staging
./deploy.sh staging --apply

# Verificar status
kubectl get pods -n direito-lux-staging
kubectl get services -n direito-lux-staging
kubectl get ingress -n direito-lux-staging
```

#### 2.2. Deploy Production
```bash
# Deploy production (após validação staging)
./deploy.sh production --apply

# Verificar status
kubectl get all -n direito-lux-production

# Verificar URLs
kubectl get ingress -n direito-lux-production
```

### 3. CI/CD Pipeline

#### 3.1. Configurar Secrets GitHub
No GitHub, vá em Settings > Secrets and Variables > Actions:

```bash
# Secrets necessários
GCP_PROJECT=direito-lux-production
GCP_SA_KEY=<base64-encoded-service-account-key>
GKE_CLUSTER_URL=<cluster-endpoint>
GKE_SA_KEY=<base64-encoded-kubeconfig>
```

#### 3.2. Ativar Workflows
```bash
# Push para develop = deploy staging automático
git checkout develop
git push origin develop

# Push para main = deploy production automático
git checkout main
git merge develop
git push origin main
```

### 4. Monitoramento e Observabilidade

#### 4.1. Acessar Dashboards
```bash
# Grafana
kubectl port-forward -n monitoring svc/grafana 3000:80

# Prometheus
kubectl port-forward -n monitoring svc/prometheus 9090:9090

# Jaeger
kubectl port-forward -n monitoring svc/jaeger 16686:16686
```

#### 4.2. Logs Centralizados
```bash
# Ver logs de um pod
kubectl logs -f deployment/auth-service -n direito-lux-production

# Ver logs de todos os pods de um serviço
kubectl logs -f -l app=auth-service -n direito-lux-production

# Kibana (se configurado)
kubectl port-forward -n monitoring svc/kibana 5601:5601
```

### 5. Operações de Produção

#### 5.1. Scaling Manual
```bash
# Escalar deployment
kubectl scale deployment auth-service --replicas=10 -n direito-lux-production

# Ver status HPA
kubectl get hpa -n direito-lux-production
```

#### 5.2. Rolling Updates
```bash
# Atualizar imagem
kubectl set image deployment/auth-service \
  auth-service=gcr.io/direito-lux-production/auth-service:v2.0.0 \
  -n direito-lux-production

# Ver status do rollout
kubectl rollout status deployment/auth-service -n direito-lux-production

# Rollback se necessário
kubectl rollout undo deployment/auth-service -n direito-lux-production
```

#### 5.3. Backup e Restore
```bash
# Backup automático (Cloud SQL)
gcloud sql backups list --instance=direito-lux-db-production

# Restore se necessário
gcloud sql backups restore <backup-id> \
  --restore-instance=direito-lux-db-production
```

### 6. URLs de Produção

| Serviço | URL | Descrição |
|---------|-----|-----------|
| **Web App** | https://app.direitolux.com | Frontend principal |
| **API Gateway** | https://api.direitolux.com | APIs REST |
| **Admin** | https://admin.direitolux.com | Painel administrativo |
| **Monitoring** | https://monitoring.direitolux.com | Grafana dashboards |
| **Status** | https://status.direitolux.com | Status page |

### 7. Disaster Recovery

#### 7.1. Backup Strategy
- **Database**: Backups automáticos diários
- **Persistent Volumes**: Snapshots automáticos
- **Configuration**: Terraform state em GCS
- **Secrets**: Backup em Cloud Secret Manager

#### 7.2. Recovery Procedures
```bash
# Restore completo
./terraform/deploy.sh production apply --auto-approve
./k8s/deploy.sh production --apply

# Restore database specific
gcloud sql backups restore <backup-id> --restore-instance=<instance>
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