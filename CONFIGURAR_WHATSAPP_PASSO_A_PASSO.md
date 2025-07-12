# 📱 CONFIGURAR WHATSAPP BUSINESS API - GUIA PASSO A PASSO

## 🎯 OBJETIVO
Configurar WhatsApp Business API para receber e enviar mensagens no ambiente de staging com 100 mensagens gratuitas por dia.

---

## 📋 PASSO 1: Criar Conta Meta for Developers

### 1.1 Acessar Portal
1. **Acesse**: https://developers.facebook.com/
2. **Clique em "Começar"**
3. **Faça login** com sua conta pessoal do Facebook
4. **Aceite os termos** de desenvolvedor

### 1.2 Verificar Conta
1. **Confirme seu número de telefone** se solicitado
2. **Verificação por email** se necessário

---

## 📋 PASSO 2: Criar App WhatsApp Business

### 2.1 Criar Novo App
1. **Clique em "Meus Apps"**
2. **Criar App** → **Business**
3. **Nome do App**: `Direito Lux Staging`
4. **Email de contato**: Seu email
5. **Finalidade**: Business

### 2.2 Adicionar WhatsApp Business
1. **Na dashboard do app**, encontre **WhatsApp Business API**
2. **Clique em "Configurar"**
3. **Aceite os termos** do WhatsApp Business

---

## 📋 PASSO 3: Configurar Número de Teste

### 3.1 Número Gratuito de Teste
1. **Vá para "API Setup"**
2. **Na seção "Send and receive messages"**
3. **Você verá um número de teste** (ex: +1 555-0199)
4. **Este número é gratuito** por 90 dias
5. **Limite**: 100 mensagens/dia

### 3.2 Adicionar Número de Destino
1. **Na seção "To"**, clique em **"Manage"**
2. **Adicione seu número pessoal** (com código do país)
3. **Exemplo**: +5511999999999
4. **Clique em "Add Number"**
5. **Verifique no WhatsApp** - você receberá um código
6. **Digite o código** para confirmar

---

## 📋 PASSO 4: Obter Credenciais

### 4.1 Access Token Temporário
1. **Na seção "API Setup"**
2. **Copie o "Temporary Access Token"**
3. **⚠️ ATENÇÃO**: Este token expira em 24h!

### 4.2 Phone Number ID
1. **Na mesma seção "API Setup"**
2. **Copie o "Phone number ID"** (longo número)

### 4.3 Business Account ID
1. **Vá para "WhatsApp Business Management"**
2. **Copie o "Business account ID"**

---

## 📋 PASSO 5: Configurar Webhook

### 5.1 Configurar URL
1. **Vá para "Configuration" → "Webhooks"**
2. **Clique em "Configure webhooks"**
3. **Callback URL**: `https://direito-lux-staging.loca.lt/webhook/whatsapp`
4. **Verify token**: `direito_lux_staging_2025`
5. **Clique em "Verify and save"**

### 5.2 Selecionar Eventos
Marque as opções:
- ✅ **messages**
- ✅ **message_deliveries**
- ✅ **message_reads**
- ✅ **messaging_postbacks**

---

## 📋 PASSO 6: Testar Configuração

### 6.1 Teste de Verificação
```bash
# O Meta fará uma requisição GET automaticamente
curl "https://direito-lux-staging.loca.lt/webhook/whatsapp?hub.mode=subscribe&hub.challenge=teste&hub.verify_token=direito_lux_staging_2025"
```

### 6.2 Teste de Envio
1. **Na seção "API Setup"**
2. **Na caixa "To"**: Selecione seu número
3. **Message**: Digite uma mensagem de teste
4. **Clique em "Send message"**
5. **Você deve receber no WhatsApp**

---

## 📋 PASSO 7: Configurar Token Permanente (Opcional)

### 7.1 System User (Para produção)
1. **Vá para "Business Settings"**
2. **System Users** → **Add**
3. **Nome**: `Direito Lux System`
4. **Role**: Admin

### 7.2 Gerar Token Permanente
1. **Selecione o System User criado**
2. **Clique em "Generate New Token"**
3. **Selecione o app**: Direito Lux Staging
4. **Permissões**: 
   - `whatsapp_business_messaging`
   - `whatsapp_business_management`
5. **Copie e salve o token**

---

## 🔧 CONFIGURAÇÃO NO SISTEMA

### Atualizar Variáveis de Ambiente
```bash
# Editar arquivo .env
nano services/notification-service/.env

# Substituir as linhas:
WHATSAPP_ACCESS_TOKEN=EAAxxxxxxxxxxxxxxxx
WHATSAPP_PHONE_NUMBER_ID=123456789012345
WHATSAPP_BUSINESS_ACCOUNT_ID=123456789012345
WHATSAPP_VERIFY_TOKEN=direito_lux_staging_2025
WHATSAPP_WEBHOOK_URL=https://direito-lux-staging.loca.lt/webhook/whatsapp
```

### Reiniciar Serviço
```bash
docker-compose restart notification-service
```

---

## ✅ VALIDAÇÃO FINAL

### Checklist de Sucesso
- [ ] App criado no Meta for Developers
- [ ] WhatsApp Business API adicionado
- [ ] Número de teste ativo
- [ ] Seu número adicionado à lista
- [ ] Access Token copiado
- [ ] Phone Number ID copiado
- [ ] Webhook configurado e verificado
- [ ] Teste de envio funcionando
- [ ] Variáveis atualizadas no sistema

### Teste E2E
1. **Envie mensagem** via API do Meta
2. **Verifique recebimento** no seu WhatsApp
3. **Responda a mensagem** no WhatsApp
4. **Verifique logs** do notification service

---

## 🚨 LIMITAÇÕES DO TESTE

### Número de Teste Gratuito
- ✅ **Válido**: 90 dias
- ✅ **Mensagens**: 100/dia
- ✅ **Destinatários**: 5 números
- ❌ **Templates**: Apenas "hello_world"

### Para Produção
- 🔄 **Número verificado** necessário
- 💰 **~R$ 0,15/mensagem**
- 📝 **Templates customizados** (aprovação Meta)

---

## 📚 LINKS ÚTEIS

- [WhatsApp Business API Docs](https://developers.facebook.com/docs/whatsapp/cloud-api)
- [Getting Started Guide](https://developers.facebook.com/docs/whatsapp/cloud-api/get-started)
- [Webhook Setup](https://developers.facebook.com/docs/whatsapp/cloud-api/webhooks)
- [Message Templates](https://developers.facebook.com/docs/whatsapp/business-management-api/message-templates)

---

**⏱️ Tempo estimado**: 30 minutos  
**💰 Custo**: R$ 0 (gratuito para teste)  
**📱 Resultado**: WhatsApp funcional para staging