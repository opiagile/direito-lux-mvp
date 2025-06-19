# Direito Lux - Plataforma de Monitoramento Jurídico

<p align="center">
  <strong>🏛️ Automatize o monitoramento de processos jurídicos com IA 🤖</strong>
</p>

<p align="center">
  <a href="#-sobre">Sobre</a> •
  <a href="#-funcionalidades">Funcionalidades</a> •
  <a href="#-arquitetura">Arquitetura</a> •
  <a href="#-começando">Começando</a> •
  <a href="#-documentação">Documentação</a> •
  <a href="#-status">Status</a>
</p>

## 🎯 Sobre

O **Direito Lux** é uma plataforma SaaS inovadora para monitoramento automatizado de processos jurídicos, integrada com a API DataJud do CNJ. Oferecemos notificações em tempo real via WhatsApp, análise inteligente com IA e uma experiência completa para escritórios de advocacia e departamentos jurídicos.

### 🏆 Diferenciais

- 🤖 **EXCLUSIVO: Interface Conversacional (MCP)** - Primeiro SaaS jurídico com bots inteligentes
- ✅ **WhatsApp em todos os planos** - Receba notificações diretamente no WhatsApp
- ✅ **Busca manual ilimitada** - Consulte processos sem restrições
- ✅ **IA integrada** - Resumos automáticos e explicação de termos jurídicos
- ✅ **Multi-tenant** - Isolamento completo entre escritórios
- ✅ **Alta disponibilidade** - Arquitetura cloud-native no GCP

## 🚀 Funcionalidades

### Core Features
- 🤖 **Bot Conversacional (MCP)** - Interaja via WhatsApp, Telegram e Claude Chat
- 📊 **Monitoramento Automático** - Acompanhe mudanças em processos 24/7
- 📱 **Notificações Multicanal** - WhatsApp, Email, Telegram e Push
- 🧠 **Assistente Virtual** - IA para análise e sumarização jurídica
- 📈 **Dashboard Analytics** - Visualize métricas e tendências
- 🔍 **Busca Avançada** - Encontre processos rapidamente
- 📄 **Geração de Documentos** - Templates personalizáveis
- 🔮 **Predição de Resultados** - ML para análise preditiva

### Planos Disponíveis

| Funcionalidade | Starter | Professional | Business | Enterprise |
|----------------|---------|--------------|----------|------------|
| Processos | 50 | 200 | 500 | Ilimitado |
| Usuários | 2 | 5 | 15 | Ilimitado |
| **Bot MCP** | ❌ | ✅ | ✅ | ✅ |
| **Comandos Bot/mês** | - | 200 | 1.000 | Ilimitado |
| WhatsApp | ✅ | ✅ | ✅ | ✅ |
| Busca Manual | Ilimitada | Ilimitada | Ilimitada | Ilimitada |
| IA Resumos | 10/mês | 50/mês | 200/mês | Ilimitado |
| Preço | R$ 99/mês | R$ 299/mês | R$ 699/mês | Sob consulta |

## 🏗️ Arquitetura

### Stack Tecnológica

- **Backend**: Go 1.21+ (microserviços)
- **AI/ML**: Python 3.11+ (FastAPI)
- **Frontend**: Next.js 14 + TypeScript
- **Mobile**: React Native + Expo
- **Database**: PostgreSQL 15 + Redis
- **Message Queue**: RabbitMQ
- **Cloud**: Google Cloud Platform
- **Orquestração**: Kubernetes (GKE)
- **Observabilidade**: Prometheus + Grafana + Jaeger

### Arquitetura de Microserviços

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Web App       │     │   Mobile App    │     │   WhatsApp Bot  │
└────────┬────────┘     └────────┬────────┘     └────────┬────────┘
         │                       │                       │
         └───────────────────────┴───────────────────────┴─────────┐
                                                                    │
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐│
│  Telegram Bot   │     │   Claude Chat   │     │    Slack Bot    ││
└────────┬────────┘     └────────┬────────┘     └────────┬────────┘│
         │                       │                       │         │
         └───────────────────────┴───────────────────────┴─────────┘
                                 │
                        ┌────────▼────────┐
                        │   MCP Service   │  🤖 NOVO!
                        │ (Bot Interface) │
                        └────────┬────────┘
                                 │
                        ┌────────▼────────┐
                        │   API Gateway   │
                        │  (Kong/Traefik) │
                        └────────┬────────┘
                                 │
     ┌───────────────────────────┴───────────────────────────┐
     │                                                       │
┌────▼─────┐  ┌─────▼─────┐  ┌─────▼─────┐  ┌─────▼─────┐  │
│   Auth   │  │  Process  │  │  DataJud  │  │    AI     │  │
│ Service  │  │  Service  │  │  Service  │  │  Service  │  │
└──────────┘  └───────────┘  └───────────┘  └───────────┘  │
                                                             │
┌────────────┐  ┌────────────┐  ┌────────────┐  ┌──────────▼─┐
│   Tenant   │  │Notification│  │   Search   │  │   Report    │
│  Service   │  │  Service   │  │  Service   │  │   Service   │ ✅
└────────────┘  └────────────┘  └────────────┘  └─────────────┘
```

## 🚀 Começando

### Pré-requisitos

- Docker Desktop 4.0+
- Go 1.21+
- Node.js 18+
- Python 3.11+
- Make

### 🎯 Quick Start - Deploy DEV (NOVO)

```bash
# 1. Clone o repositório
git clone https://github.com/direito-lux/direito-lux.git
cd direito-lux/services

# 2. Deploy automatizado completo (primeira vez)
chmod +x scripts/deploy-dev.sh
./scripts/deploy-dev.sh --clean --build

# 3. Verificar serviços rodando
./scripts/deploy-dev.sh status

# 4. Testar conectividade
./scripts/deploy-dev.sh test

# 5. Ver logs em tempo real
./scripts/deploy-dev.sh logs
```

### 🔧 Comandos Úteis

```bash
# Deploy normal (dias seguintes)
./scripts/deploy-dev.sh

# Parar todos os serviços
./scripts/deploy-dev.sh stop

# Reiniciar serviços
./scripts/deploy-dev.sh restart

# Ver endpoints disponíveis
./scripts/deploy-dev.sh endpoints

# Logs de serviço específico
./scripts/deploy-dev.sh logs ai-service
./scripts/deploy-dev.sh logs search-service
```

### 🎛️ Desenvolvimento Manual

```bash
# Iniciar todos os serviços (método antigo)
docker-compose up -d

# Ver logs
docker-compose logs -f

# Parar tudo
docker-compose down
```

## 📚 Documentação

### 📋 Documentação Principal
- [**Status da Implementação**](./STATUS_IMPLEMENTACAO.md) - ✅ O que está pronto e ❌ o que falta
- [**Deploy DEV**](./services/README-DEPLOYMENT.md) - 🚀 Guia de deploy automatizado
- [**Diretrizes de Desenvolvimento**](./DIRETRIZES_DESENVOLVIMENTO.md) - 📐 Padrões e convenções obrigatórias
- [**Setup do Ambiente**](./SETUP_AMBIENTE.md) - 🔧 Guia completo de instalação
- [**Visão Geral**](./VISAO_GERAL_DIREITO_LUX.md) - 🎯 Detalhes do produto e planos
- [**Arquitetura Full Cycle**](./ARQUITETURA_FULLCYCLE.md) - 🏗️ Arquitetura técnica detalhada
- [**Event Storming**](./EVENT_STORMING_DIREITO_LUX.md) - 📊 Domain modeling
- [**Roadmap**](./ROADMAP_IMPLEMENTACAO.md) - 🗓️ Plano de implementação
- [**MCP Service**](./services/mcp-service/MCP_SERVICE.md) - 🤖 Model Context Protocol (diferencial)

### 🔗 URLs de Desenvolvimento (Deploy DEV)

| Serviço | URL | Credenciais |
|---------|-----|-------------|
| **AI Service** | http://localhost:8000 | - |
| **Search Service** | http://localhost:8086 | - |
| **Report Service** | http://localhost:8087 | - |
| **AI Service Docs** | http://localhost:8000/docs | - |
| **Search Service Health** | http://localhost:8086/health | - |
| **Report Service Health** | http://localhost:8087/health | - |
| **PostgreSQL (Main)** | localhost:5432 | direito_lux/direito_lux_pass_dev |
| **PostgreSQL (MCP)** | localhost:5434 | mcp_user/mcp_pass_dev |
| **Redis (Main)** | localhost:6379 | redis_pass_dev |
| **Redis (MCP)** | localhost:6380 | redis_pass_dev |
| **RabbitMQ (Main)** | http://localhost:15672 | direito_lux/rabbit_pass_dev |
| **RabbitMQ (MCP)** | http://localhost:15673 | mcp_user/rabbit_pass_dev |
| **Elasticsearch** | http://localhost:9200 | - |
| **Jaeger Tracing** | http://localhost:16686 | - |

## 📊 Status do Projeto

### ✅ Implementado (Completo)

#### 🎉 TODOS OS MICROSERVIÇOS CORE 100% IMPLEMENTADOS!
- ✅ Documentação completa e planejamento
- ✅ Ambiente Docker com 15+ serviços
- ✅ **Deploy DEV Environment** - Script automatizado com todos os serviços
- ✅ Template de microserviço Go (Hexagonal Architecture)
- ✅ **Auth Service** - JWT + Multi-tenant + PostgreSQL
- ✅ **Tenant Service** - Multi-tenancy e gestão de planos com quotas
- ✅ **Process Service** - CQRS + Event Sourcing + validação CNJ
- ✅ **DataJud Service** - Pool de CNPJs + circuit breaker + rate limiting
- ✅ **Notification Service** - Multicanal com WhatsApp/Email/Telegram providers
- ✅ **AI Service** - Python/FastAPI + ML para análise jurisprudencial (deploy ready)
- ✅ **Search Service** - Go + Elasticsearch para busca avançada (deploy ready)
- ✅ **MCP Service** - Model Context Protocol com 17+ ferramentas (diferencial único)
- ✅ **Report Service** - Dashboard executivo + geração PDF/Excel + agendamento cron
- ✅ **Frontend Web App** - Next.js 14 + TypeScript + Tailwind CSS (100% completo)
- ✅ Migrações de banco robustas com triggers e funções
- ✅ Event-driven architecture base
- ✅ Correções de qualidade e estabilidade aplicadas

### 🚀 Próximos Passos (NOVA FASE)
1. **Testar Ambiente Completo** - Frontend + Backend integrados localmente
2. **CI/CD Pipeline** - GitHub Actions para build/test/deploy automatizado
3. **Kubernetes Production** - Manifests e Helm charts para GCP
4. **Terraform IaC** - Infraestrutura versionada e reproduzível
5. **Testes de Integração** - End-to-end entre microserviços
6. **Mobile App** - React Native nativo

**Progresso Total**: 🎯 **100% dos microserviços core + Frontend Web App implementados** | ~75% do projeto total

## 💻 Frontend Web App

### Stack e Tecnologias
- **Framework**: Next.js 14 com App Router
- **Linguagem**: TypeScript 100%
- **Styling**: Tailwind CSS + Shadcn/ui
- **State Management**: Zustand (stores especializados)
- **Data Fetching**: React Query (@tanstack/react-query)
- **Forms**: React Hook Form + Zod validation
- **HTTP Client**: Axios com interceptors
- **Notifications**: Sonner toast system
- **Themes**: Next-themes (light/dark mode)

### Funcionalidades Implementadas
- 🔐 **Autenticação** - Login seguro com JWT
- 📊 **Dashboard** - KPIs e atividades em tempo real
- 📁 **Gestão de Processos** - CRUD, busca, filtros, visualizações
- 🤖 **AI Assistant** - Chat interface, análise docs, jurisprudência
- 🎨 **UI/UX** - Design system completo e responsivo
- 🔍 **Busca Global** - Header search integrada
- 🌙 **Dark Mode** - Sistema completo de temas
- 📱 **Mobile Responsive** - Otimizado para todos os dispositivos

### Como Executar
```bash
# Instalar dependências
cd frontend
npm install

# Desenvolvimento
npm run dev

# Build para produção
npm run build
npm start

# Type checking
npm run type-check

# Linting
npm run lint
```

### URLs do Frontend
- **Frontend Dev**: http://localhost:3000
- **Login**: http://localhost:3000/login
- **Dashboard**: http://localhost:3000/dashboard

## 🧪 Testes

```bash
# Auth Service
cd services/auth-service

# Testes unitários
make test

# Coverage
make test-coverage

# Testes de integração
make test-integration
```

## 🛠️ Comandos Úteis

```bash
# Criar novo microserviço
./create-service.sh nome-do-servico porta

# Compilar todos os serviços
./build-all.sh

# Iniciar todos os microserviços
./start-services.sh

# Parar todos os microserviços
./stop-services.sh

# Testar ambiente completo
./test-local.sh

# Ver status dos containers
docker-compose ps

# Conectar ao PostgreSQL
docker-compose exec postgres psql -U direito_lux -d direito_lux_dev

# Ver logs de um serviço
tail -f logs/auth-service.log

# Limpar ambiente
docker-compose down -v
```

## 🤝 Contribuindo

1. Fork o projeto
2. Crie sua feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

### Padrões de Código
- Go: `gofmt`, `golangci-lint`
- Commits: Conventional Commits
- Testes: Mínimo 80% coverage
- Comentários em português

## 📄 Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 👥 Time

- **Arquiteto de Software**: Full Cycle Developer
- **Stack**: Go + Python + React + GCP

## 📞 Suporte

- **Issues**: GitHub Issues
- **Email**: suporte@direitolux.com.br
- **Docs**: [Documentação completa](./docs/)

---

<p align="center">
  Feito com ❤️ para modernizar a advocacia brasileira 🇧🇷
</p>

<p align="center">
  <strong>Transformando a justiça através da tecnologia</strong>
</p>