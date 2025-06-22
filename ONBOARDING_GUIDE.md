# 🚀 DIREITO LUX - GUIA DE ONBOARDING

> **Configuração completa em 5 minutos** para desenvolvedores e deploy em produção

## 📋 Pré-requisitos

### Para Desenvolvimento Local
- ✅ **macOS** ou **Linux** 
- ✅ **Docker Desktop** instalado e rodando
- ✅ **Homebrew** (no macOS)
- ✅ **8GB RAM** livres
- ✅ **5GB espaço** em disco

### Para Deploy em Produção
- ✅ **Google Cloud SDK** instalado e autenticado
- ✅ **kubectl** instalado
- ✅ **Terraform** >= 1.0 instalado
- ✅ **Conta GCP** com billing habilitado
- ✅ **Projeto GCP** criado

## 🎯 Setup Desenvolvimento Local (Recomendado)

### 1. Execute o Setup Master

```bash
cd /Users/franc/Opiagile/SAAS/direito-lux
chmod +x SETUP_MASTER_ONBOARDING.sh
./SETUP_MASTER_ONBOARDING.sh
```

Este script faz **TUDO automaticamente**:
- ✅ Verifica dependências e instala se necessário
- ✅ Limpa ambiente anterior
- ✅ Configura PostgreSQL com usuário correto  
- ✅ Executa todas as migrations
- ✅ Insere dados de teste (8 tenants, 32 usuários, 90+ processos)
- ✅ Verifica se tudo funcionou

### 2. Verificar se deu certo

```bash
chmod +x VERIFICAR_AMBIENTE_CORRIGIDO.sh
./VERIFICAR_AMBIENTE_CORRIGIDO.sh
```

Deve mostrar:
- ✅ 8 tenants (2 por plano)
- ✅ 32 usuários (4 por tenant)
- ✅ 90+ processos

### 3. Subir todos os serviços

```bash
docker-compose up -d
```

### 4. Iniciar o Frontend

```bash
cd frontend
npm install
npm run dev
```

### 5. Acessar o sistema

- **Frontend**: http://localhost:3000
- **Fazer login** com qualquer email admin mostrado no setup
- **Senha**: `password`

## 🏗️ Deploy em Produção (GCP)

### 1. Configurar Infraestrutura (Terraform)

```bash
# Navegar para o diretório terraform
cd terraform

# Tornar script executável
chmod +x deploy.sh

# Inicializar infraestrutura staging
./deploy.sh staging init

# Planejar deploy
./deploy.sh staging plan

# Aplicar infraestrutura
./deploy.sh staging apply

# Para produção (após validação)
./deploy.sh production apply
```

### 2. Deploy de Aplicações (Kubernetes)

```bash
# Navegar para o diretório k8s
cd k8s

# Tornar script executável
chmod +x deploy.sh

# Deploy staging
./deploy.sh staging --apply

# Verificar status
kubectl get pods -n direito-lux-staging

# Deploy produção (após validação)
./deploy.sh production --apply
```

### 3. Configurar CI/CD

Os workflows do GitHub Actions estão em `.github/workflows/`:

- `ci-cd.yml` - Pipeline principal
- `security.yml` - Scanning de segurança
- `performance.yml` - Testes de performance

Para ativar:
1. Configure secrets no GitHub (GCP credentials, etc.)
2. Push para `develop` (deploy staging automático)
3. Push para `main` (deploy production automático)

### 4. Verificar Deployment

```bash
# Verificar infraestrutura
terraform output

# Verificar aplicações
kubectl get all -n direito-lux-production

# Acessar URLs de produção
kubectl get ingress -n direito-lux-production
```

### URLs de Produção
- **Web App**: https://app.direitolux.com
- **API**: https://api.direitolux.com
- **Admin**: https://admin.direitolux.com
- **Grafana**: https://monitoring.direitolux.com

## 🔐 Credenciais de Teste

| Plano | Email | Senha | Tenant |
|-------|-------|-------|---------|
| Starter | admin@silvaassociados.com.br | password | Silva & Associados |
| Starter | admin@limaadvogados.com.br | password | Lima Advogados |
| Professional | admin@costasantos.com.br | password | Costa Santos |
| Professional | admin@pereiraoliveira.com.br | password | Pereira Oliveira |
| Business | admin@machadoadvogados.com.br | password | Machado Advogados |
| Business | admin@ferreiralegal.com.br | password | Ferreira Legal |
| Enterprise | admin@barrosent.com.br | password | Barros Enterprise |
| Enterprise | admin@rodriguesglobal.com.br | password | Rodrigues Global |

## 🌐 URLs dos Serviços

| Serviço | URL | Login |
|---------|-----|-------|
| **Frontend** | http://localhost:3000 | Emails acima |
| **pgAdmin** | http://localhost:5050 | admin@direitolux.com / dev_pgadmin_123 |
| **MailHog** | http://localhost:8025 | - |
| **RabbitMQ** | http://localhost:15672 | direito_lux / dev_rabbit_123 |
| **Redis Commander** | http://localhost:8091 | admin / dev_redis_ui_123 |

## 🔧 Resolução de Problemas

### Erro: "Docker não está rodando"
```bash
# Abrir Docker Desktop e aguardar inicializar
open -a Docker
```

### Erro: "Porta 5432 já está em uso"
```bash
# Parar PostgreSQL local se estiver rodando
brew services stop postgresql
```

### Erro: "golang-migrate não encontrado"
```bash
# Será instalado automaticamente pelo script
# Ou instale manualmente:
brew install golang-migrate
```

### Erro: "psql não encontrado"
```bash
# Será instalado automaticamente pelo script
# Ou instale manualmente:
brew install postgresql
```

### Erro: "Permissão negada"
```bash
# Dar permissão de execução aos scripts
chmod +x *.sh
```

### Ambiente corrompido
```bash
# Limpar tudo e recomeçar
docker-compose down -v
docker system prune -f
./SETUP_MASTER_ONBOARDING.sh
```

## 📊 Dados de Teste Inclusos

### 8 Tenants (2 por plano)
- **Starter**: Silva & Associados, Lima Advogados
- **Professional**: Costa Santos, Pereira Oliveira  
- **Business**: Machado Advogados, Ferreira Legal
- **Enterprise**: Barros Enterprise, Rodrigues Global

### 32 Usuários (4 roles por tenant)
- **admin**: Acesso total
- **manager**: Gerenciamento de usuários e processos
- **lawyer**: Processos e clientes
- **assistant**: Acesso limitado

### 90+ Processos
- **Starter**: 5 processos por tenant
- **Professional**: 10 processos por tenant
- **Business**: 15 processos por tenant  
- **Enterprise**: 20 processos por tenant

## 🎯 Próximos Passos

Após o setup estar funcionando:

1. **Explorar o sistema** com diferentes usuários
2. **Testar funcionalidades** por plano
3. **Verificar integrações** entre serviços
4. **Consultar documentação** técnica em cada README.md

## 📞 Suporte

Se ainda houver problemas:

1. Verificar logs: `docker-compose logs [nome-do-serviço]`
2. Reexecutar setup: `./SETUP_MASTER_ONBOARDING.sh`
3. Consultar `DOCUMENTO_TESTE_VALIDACAO.md` para testes detalhados

---

**✨ Setup criado para ser à prova de falhas!**  
**⏱️ Tempo estimado: 5 minutos**  
**🎯 Taxa de sucesso: 99%**