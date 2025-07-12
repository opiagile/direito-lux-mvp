# 🚀 INSTRUÇÕES RÁPIDAS - TELEGRAM BOT

## ⏱️ 5 MINUTOS PARA CONFIGURAR

### 1️⃣ CRIAR BOT (2 min)

1. **Abra Telegram** → Busque `@BotFather`
2. Envie `/newbot`
3. Nome: `Direito Lux Staging Bot`
4. Username: `direitolux_staging_bot`
5. **COPIE O TOKEN** (tipo: `7458394857:AAHKz9XjB8vK_2QxYz0-fG8kNvM_xQz7890`)

### 2️⃣ CONFIGURAR BOT (3 min)

```bash
# Execute o script de configuração
./configure_telegram_bot.sh

# Cole o token quando solicitado
# O script fará todo o resto automaticamente!
```

### 3️⃣ PRONTO! ✅

O script irá:
- ✅ Validar o token
- ✅ Configurar comandos
- ✅ Atualizar .env
- ✅ Reiniciar serviço
- ✅ Configurar webhook
- ✅ Testar integração

---

## 📱 TESTAR BOT

1. Acesse: https://t.me/direitolux_staging_bot
2. Clique em "START"
3. Envie `/help`
4. Bot deve responder!

---

## 🆘 PROBLEMAS?

- **Username já existe?** Tente: `direitolux_staging_2025_bot`
- **Token inválido?** Verifique se copiou corretamente
- **Bot não responde?** Verifique logs: `docker-compose logs -f notification-service`

---

**📖 Guia completo**: [CRIAR_TELEGRAM_BOT_PASSO_A_PASSO.md](./CRIAR_TELEGRAM_BOT_PASSO_A_PASSO.md)