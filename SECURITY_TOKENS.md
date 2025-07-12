# 🔒 GERENCIAMENTO DE TOKENS - SEGURANÇA

## ⚠️ IMPORTANTE: Resolução GitGuardian

**GitGuardian detectou token exposto em commits anteriores.**

### 🚨 AÇÃO TOMADA

1. **Tokens removidos** dos arquivos commitados
2. **Placeholders seguros** implementados
3. **Gitignore atualizado** para prevenir futuras exposições
4. **Variáveis de ambiente** (.env) protegidas

---

## 📋 ONDE FICAM OS TOKENS AGORA

### ✅ SEGURO (.env - NÃO COMMITADO)
```bash
# services/notification-service/.env
TELEGRAM_BOT_TOKEN=SEU_TOKEN_REAL_AQUI
WHATSAPP_ACCESS_TOKEN=SEU_TOKEN_WHATSAPP_AQUI
```

### ❌ REMOVIDO DOS COMMITADOS
- `test_telegram_bot_manual.sh` → Usa `${TELEGRAM_BOT_TOKEN}`
- `telegram_bot_info.txt` → `[REMOVIDO_POR_SEGURANCA]`

---

## 🔧 CONFIGURAÇÃO SEGURA

### 1. Para Desenvolvimento Local
```bash
# Copie o arquivo de exemplo
cp services/notification-service/.env.example services/notification-service/.env

# Edite com seus tokens reais
nano services/notification-service/.env
```

### 2. Para Produção/Staging
```bash
# Use variáveis de ambiente do sistema
export TELEGRAM_BOT_TOKEN="seu_token_real"
export WHATSAPP_ACCESS_TOKEN="seu_token_real"
```

---

## 🛡️ BOAS PRÁTICAS IMPLEMENTADAS

1. **Never commit tokens** - Sempre use .env
2. **Gitignore robusto** - Todos os arquivos sensíveis protegidos
3. **Placeholders** - Scripts usam variáveis de ambiente
4. **Documentação** - Este arquivo explica o processo

---

## 📝 ARQUIVO .env.example

```bash
# Telegram Bot (STAGING)
TELEGRAM_BOT_TOKEN=YOUR_TELEGRAM_BOT_TOKEN_HERE
TELEGRAM_WEBHOOK_URL=https://your-domain.com/webhook/telegram

# WhatsApp Business API (STAGING)
WHATSAPP_ACCESS_TOKEN=YOUR_WHATSAPP_TOKEN_HERE
WHATSAPP_PHONE_NUMBER_ID=YOUR_PHONE_NUMBER_ID_HERE
WHATSAPP_WEBHOOK_URL=https://your-domain.com/webhook/whatsapp
WHATSAPP_VERIFY_TOKEN=your_verify_token_here
```

---

## 🚀 PARA NOVOS TOKENS

1. **Crie o token** na plataforma (Telegram, WhatsApp, etc.)
2. **Adicione no .env** (nunca no código)
3. **Teste localmente**
4. **Configure no servidor** via variáveis de ambiente

---

## 📧 CONTATO PARA SUPORTE

Email: contato@direitolux.com.br  
Bot Telegram: @direitolux_staging_bot

---

**⚠️ LEMBRE-SE**: Tokens são como senhas. Nunca compartilhe ou comite!