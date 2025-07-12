# 📱 INSTRUÇÕES RÁPIDAS - WHATSAPP BUSINESS API

## ⏱️ 30 MINUTOS PARA CONFIGURAR

### 1️⃣ CRIAR APP META (10 min)

1. **Acesse**: https://developers.facebook.com/
2. **Login** com Facebook pessoal
3. **Criar App** → Business → "Direito Lux Staging"
4. **Adicionar** WhatsApp Business API

### 2️⃣ OBTER CREDENCIAIS (10 min)

1. **API Setup** → Copie:
   - **Access Token** (temporário 24h)
   - **Phone Number ID**
2. **Adicione seu número** na lista de teste
3. **Verifique no WhatsApp** (código de confirmação)

### 3️⃣ CONFIGURAR WEBHOOK (5 min)

1. **Configuration** → **Webhooks**
2. **URL**: `https://direito-lux-staging.loca.lt/webhook/whatsapp`
3. **Token**: `direito_lux_staging_2025`
4. **Eventos**: messages, deliveries, reads

### 4️⃣ CONFIGURAR SISTEMA (5 min)

```bash
# Execute o script de configuração
./configure_whatsapp_api.sh

# Cole as credenciais quando solicitado
# O script fará todo o resto automaticamente!
```

---

## 📱 TESTAR WHATSAPP

1. **No Meta console**: Envie mensagem de teste
2. **Receba no WhatsApp** 
3. **Responda no WhatsApp**
4. **Verifique logs**: `docker-compose logs -f notification-service`

---

## ✅ RESULTADO

- ✅ **100 mensagens/dia gratuitas**
- ✅ **Webhook funcionando**
- ✅ **Envio e recebimento** 
- ✅ **Staging completo**

---

## 🆘 PROBLEMAS?

- **Webhook failed?** Verifique se túnel está ativo
- **Token inválido?** Gere novo token temporário
- **Número não autorizado?** Adicione à lista de teste

---

**📖 Guia completo**: [CONFIGURAR_WHATSAPP_PASSO_A_PASSO.md](./CONFIGURAR_WHATSAPP_PASSO_A_PASSO.md)