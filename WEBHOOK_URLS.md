# 🌐 URLs dos Webhooks - Direito Lux Staging

## 📡 Túnel HTTPS Ativo

**URL Base**: https://direito-lux-staging.loca.lt

### 🔗 Endpoints Disponíveis

| Serviço | Endpoint | URL Completa |
|---------|----------|--------------|
| **Health Check** | `/health` | https://direito-lux-staging.loca.lt/health |
| **Telegram Webhook** | `/webhook/telegram` | https://direito-lux-staging.loca.lt/webhook/telegram |
| **WhatsApp Webhook** | `/webhook/whatsapp` | https://direito-lux-staging.loca.lt/webhook/whatsapp |
| **Billing ASAAS** | `/billing/webhooks/asaas` | https://direito-lux-staging.loca.lt/billing/webhooks/asaas |
| **Billing Crypto** | `/billing/webhooks/crypto` | https://direito-lux-staging.loca.lt/billing/webhooks/crypto |

## ✅ Status de Configuração

### 🤖 Telegram Bot
- **Status**: ✅ Túnel configurado
- **Webhook URL**: https://direito-lux-staging.loca.lt/webhook/telegram
- **Próximo passo**: Configurar webhook via BotFather com token real

### 📱 WhatsApp Business API
- **Status**: ✅ Túnel configurado
- **Webhook URL**: https://direito-lux-staging.loca.lt/webhook/whatsapp
- **Próximo passo**: Configurar no Meta Business

### 💰 Payment Gateways
- **ASAAS**: https://direito-lux-staging.loca.lt/billing/webhooks/asaas
- **NOWPayments**: https://direito-lux-staging.loca.lt/billing/webhooks/crypto

## 🛠️ Comandos Úteis

### Verificar Status do Túnel
```bash
curl -s https://direito-lux-staging.loca.lt/health
```

### Configurar Webhook Telegram (quando tiver token)
```bash
curl -X POST "https://api.telegram.org/bot<TOKEN>/setWebhook" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://direito-lux-staging.loca.lt/webhook/telegram",
    "allowed_updates": ["message", "callback_query"]
  }'
```

### Reiniciar Túnel
```bash
pkill -f "npx localtunnel"
npx localtunnel --port 8085 --subdomain direito-lux-staging
```

## 🚨 Notas Importantes

1. **Túnel temporário**: Esta URL é válida apenas enquanto o processo estiver rodando
2. **Subdomain**: Usando `direito-lux-staging` para consistência
3. **SSL**: HTTPS automático via localtunnel
4. **Rate limiting**: Sem limitações para desenvolvimento

## 🎯 Próximos Passos

1. ✅ **Túnel HTTPS configurado**
2. ⏳ **Configurar webhook Telegram** (aguarda token real)
3. ⏳ **Configurar webhook WhatsApp** (aguarda API keys)
4. ⏳ **Testar webhooks** (envio/recebimento)

---

**Atualizado em**: 2025-07-11  
**Status**: ✅ Funcionando