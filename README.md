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

- ✅ **WhatsApp em todos os planos** - Receba notificações diretamente no WhatsApp
- ✅ **Busca manual ilimitada** - Consulte processos sem restrições
- ✅ **IA integrada** - Resumos automáticos e explicação de termos jurídicos
- ✅ **Multi-tenant** - Isolamento completo entre escritórios
- ✅ **Alta disponibilidade** - Arquitetura cloud-native no GCP

## 🚀 Funcionalidades

### Core Features
- 📊 **Monitoramento Automático** - Acompanhe mudanças em processos 24/7
- 📱 **Notificações Multicanal** - WhatsApp, Email, Telegram e Push
- 🤖 **Assistente Virtual** - IA para análise e sumarização
- 📈 **Dashboard Analytics** - Visualize métricas e tendências
- 🔍 **Busca Avançada** - Encontre processos rapidamente
- 📄 **Geração de Documentos** - Templates personalizáveis
- 🔮 **Predição de Resultados** - ML para análise preditiva

### Planos Disponíveis

| Funcionalidade | Starter | Professional | Business | Enterprise |
|----------------|---------|--------------|----------|------------|
| Processos | 50 | 200 | 500 | Ilimitado |
| Usuários | 2 | 5 | 15 | Ilimitado |
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
         │                       │                         │
         └───────────────────────┴─────────────────────────┘
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
│  Service   │  │  Service   │  │  Service   │  │   Service   │
└────────────┘  └────────────┘  └────────────┘  └─────────────┘
```

## 🚀 Começando

### Pré-requisitos

- Docker Desktop 4.0+
- Go 1.21+
- Node.js 18+
- Python 3.11+
- Make

### Quick Start

```bash
# 1. Clone o repositório
git clone https://github.com/direito-lux/direito-lux.git
cd direito-lux

# 2. Configure o ambiente
cp .env.example .env
# Edite .env com suas configurações

# 3. Inicie os serviços de infraestrutura
docker-compose up -d postgres redis rabbitmq

# 4. Execute as migrações
./scripts/setup-postgres.sh

# 5. Inicie o Auth Service
cd services/auth-service
make dev

# 6. Teste a API
curl http://localhost:8081/health
```

### Desenvolvimento com Docker Compose

```bash
# Iniciar todos os serviços
docker-compose up -d

# Ver logs
docker-compose logs -f

# Parar tudo
docker-compose down
```

## 📚 Documentação

### 📋 Documentação Principal
- [**Status da Implementação**](./STATUS_IMPLEMENTACAO.md) - ✅ O que está pronto e ❌ o que falta
- [**Setup do Ambiente**](./SETUP_AMBIENTE.md) - 🔧 Guia completo de instalação
- [**Visão Geral**](./VISAO_GERAL_DIREITO_LUX.md) - 🎯 Detalhes do produto e planos
- [**Arquitetura Full Cycle**](./ARQUITETURA_FULLCYCLE.md) - 🏗️ Arquitetura técnica detalhada
- [**Event Storming**](./EVENT_STORMING_DIREITO_LUX.md) - 📊 Domain modeling
- [**Roadmap**](./ROADMAP_IMPLEMENTACAO.md) - 🗓️ Plano de implementação
- [**Processo de Documentação**](./PROCESSO_DOCUMENTACAO.md) - 📝 Como manter docs atualizadas

### 🔗 URLs de Desenvolvimento

| Serviço | URL | Credenciais |
|---------|-----|-------------|
| **API Gateway** | http://localhost:8000 | - |
| **Auth Service** | http://localhost:8081 | - |
| **PostgreSQL** | localhost:5432 | direito_lux/dev_password_123 |
| **Redis** | localhost:6379 | dev_redis_123 |
| **RabbitMQ** | http://localhost:15672 | direito_lux/dev_rabbit_123 |
| **Keycloak** | http://localhost:8080 | admin/admin123 |
| **Jaeger** | http://localhost:16686 | - |
| **Prometheus** | http://localhost:9090 | - |
| **Grafana** | http://localhost:3000 | admin/admin123 |
| **Kibana** | http://localhost:5601 | - |

## 📊 Status do Projeto

### ✅ Implementado
- ✅ Documentação completa e planejamento
- ✅ Ambiente Docker com 15+ serviços
- ✅ Template de microserviço Go (Hexagonal Architecture)
- ✅ Auth Service completo com JWT + Multi-tenant
- ✅ Migrações de banco de dados
- ✅ Event-driven architecture base

### 🚧 Em Desenvolvimento
- 🔄 Tenant Service
- 🔄 CI/CD Pipeline

### ⏳ Próximos Passos
1. Tenant Service (gerenciamento de organizações)
2. Process Service (core business)
3. DataJud Service (integração CNJ)
4. Notification Service (WhatsApp/Email)
5. AI Service (Python/FastAPI)

**Progresso Total**: ~25% completo

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
./scripts/create-service.sh nome-do-servico

# Ver status dos containers
docker-compose ps

# Conectar ao PostgreSQL
docker-compose exec postgres psql -U direito_lux -d direito_lux_dev

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