# 🤖 CONFIGURAR TELEGRAM BOT STAGING - PASSO A PASSO

## 📋 Instruções para Criar Bot Real

### **PASSO 1: Criar Bot via BotFather**

1. **Abrir Telegram** (mobile ou desktop)
2. **Buscar**: `@BotFather`
3. **Iniciar conversa**: `/start`
4. **Criar novo bot**: `/newbot`
5. **Nome do bot**: `Direito Lux Staging Bot`
6. **Username do bot**: `direitolux_staging_bot` (deve terminar com _bot)

### **PASSO 2: Obter Token**

O BotFather retornará algo como:
```
Done! Congratulations on your new bot. You will find it at t.me/direitolux_staging_bot.

Use this token to access the HTTP API:
1234567890:AAAA-BBBBccccDDDDeeeeFFFGGGhhhhIII

Keep your token secure and store it safely, it can be used by anyone to control your bot.
```

**IMPORTANTE**: Salvar este token com segurança!

### **PASSO 3: Configurar Comandos do Bot**

Enviar para BotFather:
```
/setcommands
```

Selecionar o bot criado e enviar:
```
help - 🆘 Ajuda e comandos disponíveis
status - 📊 Status dos seus processos
agenda - 📅 Agenda de prazos
busca - 🔍 Buscar processos
relatorio - 📈 Relatório rápido
configurar - ⚙️ Configurações
```

### **PASSO 4: Configurar Descrição**

```
/setdescription
```

Texto:
```
🏛️ Direito Lux - Assistente Jurídico Inteligente

Bot oficial para monitoramento de processos, notificações automáticas e análises jurídicas.

✅ Ambiente STAGING - Para testes e validação
```

### **PASSO 5: Testar Bot Manualmente**

1. Acesse: `https://t.me/direitolux_staging_bot`
2. Clique em **START**
3. Envie: `/help`
4. **Não haverá resposta ainda** (normal - webhook não configurado)

---

## 🔧 CONFIGURAÇÃO TÉCNICA

### **Token Obtido** (exemplo):
```
TELEGRAM_BOT_TOKEN_STAGING=1234567890:AAAA-BBBBccccDDDDeeeeFFFGGGhhhhIII
```

### **URLs Necessárias** (quando tivermos HTTPS):
```
TELEGRAM_WEBHOOK_URL=https://staging.direitolux.com.br/webhook/telegram
```

### **Configurar Webhook** (depois do HTTPS):
```bash
curl -X POST "https://api.telegram.org/bot${BOT_TOKEN}/setWebhook" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://staging.direitolux.com.br/webhook/telegram",
    "allowed_updates": ["message", "callback_query"]
  }'
```

---

## ✅ PRÓXIMOS PASSOS

1. ✅ **Bot criado** - Token obtido
2. ⏳ **Configurar HTTPS** - ngrok ou domínio 
3. ⏳ **Configurar webhook** - Conectar ao serviço
4. ⏳ **Testar integração** - Enviar mensagens reais

**Status**: Pronto para integração quando tivermos URLs HTTPS!