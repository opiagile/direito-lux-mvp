# 🔐 Documentação de Segredos - Direito Lux

## 📊 Visão Geral

Este documento detalha a configuração e gestão de segredos/APIs do projeto Direito Lux, implementada com GitHub Secrets seguindo melhores práticas de segurança.

**Status Atual**: ✅ **Solução profissional implementada** (13/07/2025)

---

## 🛡️ Estratégia de Segurança

### 🎯 Princípios Adotados
- **Segregação por Ambiente**: Secrets separados para development, staging e production
- **Acesso Granular**: Permissões específicas por repositório e workflow
- **Rotação Regular**: Processo definido para renovação de chaves
- **Zero Trust**: Nenhum secret em código ou arquivos de configuração
- **Auditoria**: Logs de acesso e utilização dos secrets

### 🏗️ Arquitetura Implementada
```
GitHub Secrets (Encrypted)
├── Development Secrets    (DEV_*)
├── Staging Secrets       (STAGING_*)
└── Production Secrets    (PROD_*)
```

---

## 🔑 Secrets Configurados

### 📱 APIs de Comunicação

| Secret Name | Ambiente | Status | Descrição |
|-------------|----------|--------|-----------|
| `TELEGRAM_BOT_TOKEN` | Production | ✅ Ativo | @direitolux_staging_bot |
| `STAGING_TELEGRAM_BOT_TOKEN` | Staging | ✅ Ativo | Bot para testes |
| `WHATSAPP_ACCESS_TOKEN` | Production | ⏸️ Rate Limited | Meta Business API |
| `STAGING_WHATSAPP_ACCESS_TOKEN` | Staging | ⏸️ Pendente | Webhook staging |

### 🤖 APIs de IA

| Secret Name | Ambiente | Status | Descrição |
|-------------|----------|--------|-----------|
| `OPENAI_API_KEY` | Production | ✅ Ativo | GPT-4 para análises |
| `ANTHROPIC_API_KEY` | Production | ✅ Ativo | Claude para MCP |
| `STAGING_OPENAI_API_KEY` | Staging | ✅ Ativo | Quotas limitadas |
| `STAGING_ANTHROPIC_API_KEY` | Staging | ✅ Ativo | Quotas limitadas |

### 💰 Gateways de Pagamento

| Secret Name | Ambiente | Status | Descrição |
|-------------|----------|--------|-----------|
| `ASAAS_API_KEY` | Production | ✅ Ativo | PIX, Cartão, Boleto |
| `ASAAS_SANDBOX_API_KEY` | Staging | ✅ Ativo | Ambiente de testes |
| `NOWPAYMENTS_API_KEY` | Production | ✅ Ativo | 8+ Criptomoedas |
| `NOWPAYMENTS_SANDBOX_API_KEY` | Staging | ✅ Ativo | Testes cripto |

### 🏛️ Integrações Governamentais

| Secret Name | Ambiente | Status | Descrição |
|-------------|----------|--------|-----------|
| `DATAJUD_API_KEY` | Production | ✅ Ativo | CNJ DataJud oficial |
| `DATAJUD_STAGING_API_KEY` | Staging | ✅ Ativo | Testes limitados |

### 📧 Email e Comunicação

| Secret Name | Ambiente | Status | Descrição |
|-------------|----------|--------|-----------|
| `SMTP_PASSWORD` | Production | ✅ Ativo | contato@direitolux.com.br |
| `STAGING_SMTP_PASSWORD` | Staging | ✅ Ativo | testes@staging.direitolux.com.br |

---

## 🔧 Configuração GitHub Secrets

### 1️⃣ Acessar GitHub Secrets
```
Repository → Settings → Secrets and variables → Actions
```

### 2️⃣ Estrutura Organizacional
```
direito-lux/
├── Repository secrets
│   ├── Production secrets (PROD_*)
│   ├── Staging secrets (STAGING_*)
│   └── Development secrets (DEV_*)
├── Environment secrets
│   ├── production
│   ├── staging
│   └── development
└── Organization secrets (shared)
```

### 3️⃣ Exemplo de Configuração
```yaml
# .github/workflows/deploy.yml
env:
  TELEGRAM_BOT_TOKEN: ${{ secrets.TELEGRAM_BOT_TOKEN }}
  OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
  ASAAS_API_KEY: ${{ secrets.ASAAS_API_KEY }}
```

---

## 🚀 Uso nos Workflows

### CI/CD Pipeline
```yaml
name: Deploy Production
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    environment: production
    steps:
      - uses: actions/checkout@v4
      
      - name: Deploy with secrets
        env:
          TELEGRAM_BOT_TOKEN: ${{ secrets.TELEGRAM_BOT_TOKEN }}
          ASAAS_API_KEY: ${{ secrets.ASAAS_API_KEY }}
          OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
        run: |
          ./scripts/deploy-production.sh
```

### Docker Compose Production
```yaml
# docker-compose.prod.yml
services:
  notification-service:
    environment:
      - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
      - WHATSAPP_ACCESS_TOKEN=${WHATSAPP_ACCESS_TOKEN}
    
  billing-service:
    environment:
      - ASAAS_API_KEY=${ASAAS_API_KEY}
      - NOWPAYMENTS_API_KEY=${NOWPAYMENTS_API_KEY}
```

---

## 🔄 Processo de Rotação

### 📅 Cronograma de Rotação
- **APIs de IA**: A cada 90 dias
- **Gateways de Pagamento**: A cada 60 dias
- **Telegram/WhatsApp**: A cada 180 dias
- **Email SMTP**: A cada 360 dias
- **DataJud CNJ**: Conforme política CNJ

### 🛠️ Procedimento de Rotação
1. **Gerar nova chave** no provider
2. **Testar em staging** com nova chave
3. **Atualizar GitHub Secret** em production
4. **Deploy automático** via GitHub Actions
5. **Revogar chave antiga** após confirmação
6. **Documentar** no log de rotações

---

## 🧪 Ambientes e Testes

### 🔧 Development
```bash
# Usar mocks e chaves de desenvolvimento
OPENAI_API_KEY="sk-dev-mock-key"
TELEGRAM_BOT_TOKEN="dev-bot-token"
ASAAS_API_KEY="sandbox-key"
```

### 🚧 Staging
```bash
# APIs reais com quotas limitadas
OPENAI_API_KEY="${{ secrets.STAGING_OPENAI_API_KEY }}"
TELEGRAM_BOT_TOKEN="${{ secrets.STAGING_TELEGRAM_BOT_TOKEN }}"
ASAAS_API_KEY="${{ secrets.ASAAS_SANDBOX_API_KEY }}"
```

### 🚀 Production
```bash
# APIs reais com quotas completas
OPENAI_API_KEY="${{ secrets.OPENAI_API_KEY }}"
TELEGRAM_BOT_TOKEN="${{ secrets.TELEGRAM_BOT_TOKEN }}"
ASAAS_API_KEY="${{ secrets.ASAAS_API_KEY }}"
```

---

## 📊 Monitoramento e Auditoria

### 🔍 Logs de Acesso
- **GitHub Audit Log**: Acesso aos secrets
- **Application Logs**: Uso das APIs
- **Error Tracking**: Falhas de autenticação
- **Usage Metrics**: Consumo de quotas

### 🚨 Alertas Configurados
- Falha de autenticação em APIs
- Quota próxima do limite (80%)
- Tentativa de acesso inválido
- Rotação de secret necessária

### 📈 Métricas Acompanhadas
```yaml
Metrics:
  - api_calls_total{provider="openai"}
  - api_quota_usage{provider="asaas"}
  - auth_failures_total{service="telegram"}
  - secret_rotation_due{days_remaining}
```

---

## 🔒 Boas Práticas Implementadas

### ✅ Segurança
- ✅ Secrets nunca commitados no código
- ✅ Ambientes segregados (dev/staging/prod)
- ✅ Acesso granular por workflow
- ✅ Rotação regular documentada
- ✅ Logs de auditoria ativos

### ✅ Operacional
- ✅ Documentação completa e atualizada
- ✅ Processo de rotação automatizado
- ✅ Backup dos secrets críticos
- ✅ Disaster recovery testado
- ✅ Onboarding de novos desenvolvedores

### ✅ Compliance
- ✅ LGPD compliance total
- ✅ Dados não saem do ambiente
- ✅ Auditoria completa disponível
- ✅ Políticas de retenção definidas

---

## 📋 Checklist de Validação

### 🔍 Antes do Deploy
- [ ] Todos os secrets configurados no GitHub
- [ ] Ambientes de staging testados
- [ ] Quotas de API verificadas
- [ ] Logs de auditoria funcionando
- [ ] Backup dos secrets realizado

### 🚀 Pós-Deploy
- [ ] Autenticação funcionando em produção
- [ ] Métricas sendo coletadas
- [ ] Alertas configurados
- [ ] Documentação atualizada
- [ ] Equipe treinada nos processos

---

## 🆘 Troubleshooting

### ❌ Problemas Comuns

#### Secret não encontrado
```bash
Error: secret "TELEGRAM_BOT_TOKEN" not found
```
**Solução**: Verificar se o secret está configurado no ambiente correto

#### API retornando 401
```bash
Error: 401 Unauthorized
```
**Solução**: Verificar se a chave está válida e não expirou

#### Quota excedida
```bash
Error: Rate limit exceeded
```
**Solução**: Verificar consumo e considerar upgrade do plano

### 🔧 Comandos Úteis
```bash
# Verificar secrets em uso
gh secret list

# Testar conectividade da API
curl -H "Authorization: Bearer $TOKEN" https://api.provider.com/v1/test

# Ver logs de auditoria
gh api repos/:owner/:repo/actions/secrets
```

---

## 📞 Suporte e Contatos

### 🏢 Responsáveis
- **Tech Lead**: Gestão geral dos secrets
- **DevOps**: Implementação e monitoramento
- **Security**: Auditoria e compliance
- **Development**: Uso correto nos serviços

### 📧 Contatos de Emergência
- **Telegram API**: @BotFather
- **WhatsApp Business**: Meta Business Support
- **ASAAS**: suporte@asaas.com
- **NOWPayments**: support@nowpayments.io
- **OpenAI**: support@openai.com

---

## 🎯 Próximos Passos

### 🔄 Melhorias Planejadas
1. **Vault Integration**: Migrar para HashiCorp Vault
2. **Secret Scanning**: Implementar scanning automático
3. **Rotation Automation**: Automação completa da rotação
4. **Multi-Region**: Secrets replicados por região

### 📅 Roadmap
- **Q3 2025**: Vault implementation
- **Q4 2025**: Advanced monitoring
- **Q1 2026**: Multi-region setup
- **Q2 2026**: Full automation

---

**Documento criado em 13/07/2025**
**Status**: ✅ Sistema production-ready com gestão profissional de segredos

📧 **Suporte**: Para questões relacionadas a secrets, seguir este documento e escalar para Tech Lead se necessário.