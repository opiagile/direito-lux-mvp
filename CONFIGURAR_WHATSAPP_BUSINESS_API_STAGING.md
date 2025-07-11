# 📱 CONFIGURAR WHATSAPP BUSINESS API - STAGING

## 🎯 CONFIGURAÇÃO STAGING COM QUOTAS LIMITADAS

### 1. CRIAR CONTA FACEBOOK BUSINESS
1. Acesse: https://business.facebook.com/
2. Crie conta Business (gratuita)
3. Confirme empresa e dados

### 2. CONFIGURAR WHATSAPP BUSINESS API

#### 2.1 Acesse Meta for Developers
1. Acesse: https://developers.facebook.com/
2. Faça login com conta Facebook Business
3. Crie novo app → "Business" → "WhatsApp Business API"

#### 2.2 Configurar App WhatsApp
```
Nome do App: Direito Lux Staging
Categoria: Business
Propósito: Notificações jurídicas
```

#### 2.3 Configurar Webhook
```
Webhook URL: https://locking-model-sports-anti.trycloudflare.com/webhook/whatsapp
Verify Token: direito_lux_staging_2025
```

#### 2.4 Configurar Número de Teste
- WhatsApp Business API fornece número de teste gratuito
- Limite: 100 mensagens/dia (perfeito para staging)
- Válido para 90 dias

### 3. OBTER CREDENCIAIS

#### 3.1 Access Token (Temporário - 24h)
```
App → WhatsApp → API Setup → Temporary Access Token
```

#### 3.2 Configurar Permanent Access Token
```
App → Settings → Advanced → System Users
Create System User: "Direito Lux Staging"
Assign to App: Direito Lux Staging
Generate Token: whatsapp_business_messaging, whatsapp_business_management
```

### 4. VARIÁVEIS DE AMBIENTE

```bash
# WhatsApp Business API
WHATSAPP_ACCESS_TOKEN=EAAxxxxxxxxxxxxxxxx
WHATSAPP_PHONE_NUMBER_ID=123456789012345
WHATSAPP_BUSINESS_ACCOUNT_ID=123456789012345
WHATSAPP_VERIFY_TOKEN=direito_lux_staging_2025
WHATSAPP_WEBHOOK_URL=https://locking-model-sports-anti.trycloudflare.com/webhook/whatsapp
```

### 5. CONFIGURAR WEBHOOK NO FACEBOOK

#### 5.1 Webhook Configuration
```
Callback URL: https://locking-model-sports-anti.trycloudflare.com/webhook/whatsapp
Verify Token: direito_lux_staging_2025
```

#### 5.2 Webhook Fields
- ✅ messages
- ✅ messaging_postbacks
- ✅ message_deliveries
- ✅ message_reads

### 6. TESTAR CONFIGURAÇÃO

#### 6.1 Verificar Webhook
```bash
# Meta vai fazer GET request para verificar webhook
curl -X GET "https://locking-model-sports-anti.trycloudflare.com/webhook/whatsapp?hub.mode=subscribe&hub.challenge=teste&hub.verify_token=direito_lux_staging_2025"
```

#### 6.2 Testar Envio de Mensagem
```bash
curl -X POST "https://graph.facebook.com/v18.0/PHONE_NUMBER_ID/messages" \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "messaging_product": "whatsapp",
    "to": "SEU_TELEFONE_TESTE",
    "type": "text",
    "text": {
      "body": "Teste WhatsApp API - Direito Lux Staging"
    }
  }'
```

### 7. LIMITAÇÕES STAGING

#### 7.1 Número de Teste
- ✅ Gratuito por 90 dias
- ✅ 100 mensagens/dia
- ✅ Perfeito para staging

#### 7.2 Destinatários
- ✅ 5 números de telefone para teste
- ✅ Números devem ser verificados no console

#### 7.3 Templates
- ✅ Templates pré-aprovados disponíveis
- ✅ "hello_world" template gratuito
- ✅ Templates customizados precisam aprovação

### 8. PRÓXIMOS PASSOS

1. ✅ Configurar WhatsApp Business API
2. ⏳ Configurar OpenAI API com budget limitado
3. ⏳ Configurar Claude API com budget limitado
4. ⏳ Testar integração completa

### 9. CUSTO STAGING

- **WhatsApp Business API**: R$ 0/mês (teste gratuito)
- **Mensagens**: R$ 0 (até 100/dia)
- **Templates**: R$ 0 (hello_world gratuito)

**Total WhatsApp**: R$ 0/mês para staging 🎉

### 10. MIGRAÇÃO PARA PRODUÇÃO

Quando sair do staging:
1. Configurar número WhatsApp Business verificado
2. Configurar templates personalizados
3. Aumentar quotas de mensagens
4. Custos: ~R$ 0,15/mensagem

---

## 🔗 LINKS ÚTEIS

- [WhatsApp Business API Docs](https://developers.facebook.com/docs/whatsapp)
- [Meta for Developers](https://developers.facebook.com/)
- [Webhook Setup Guide](https://developers.facebook.com/docs/whatsapp/cloud-api/webhooks)
- [Message Templates](https://developers.facebook.com/docs/whatsapp/business-management-api/message-templates)