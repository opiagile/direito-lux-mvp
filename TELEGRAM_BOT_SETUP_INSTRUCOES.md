# 🤖 INSTRUÇÕES: Criar Bot Telegram para Direito Lux

## 📋 Passo a Passo - Criar Bot Real

### **PASSO 1: Acessar BotFather**
1. Abra o Telegram (mobile ou desktop)
2. Busque por: `@BotFather`
3. Clique em "START" para iniciar conversa

### **PASSO 2: Criar Novo Bot**
1. Envie o comando: `/newbot`
2. **Nome do bot**: `Direito Lux Staging`
3. **Username**: `direitolux_staging_bot` (deve terminar com `_bot`)

### **PASSO 3: Salvar Token**
O BotFather retornará uma mensagem como:
```
Done! Congratulations on your new bot. You will find it at 
t.me/direitolux_staging_bot. You can now add a description...

Use this token to access the HTTP API:
1234567890:AAAA-BBBBccccDDDDeeeeFFFGGGhhhhIII

Keep your token secure and store it safely, it can be used by 
anyone to control your bot.
```

**🔑 COPIE E SALVE ESTE TOKEN COM SEGURANÇA!**

### **PASSO 4: Configurar Comandos**
1. Envie: `/setcommands`
2. Selecione seu bot
3. Copie e cole:
```
start - 🚀 Iniciar conversa com o bot
help - 🆘 Ajuda e comandos disponíveis
status - 📊 Status dos seus processos
agenda - 📅 Agenda de prazos importantes
busca - 🔍 Buscar processos jurídicos
relatorio - 📈 Relatório rápido dos processos
configurar - ⚙️ Configurações do bot
```

### **PASSO 5: Configurar Descrição**
1. Envie: `/setdescription`
2. Selecione seu bot
3. Copie e cole:
```
🏛️ Direito Lux - Assistente Jurídico Inteligente

Bot oficial para monitoramento de processos, notificações automáticas e análises jurídicas.

✅ Ambiente STAGING - Para testes e validação
```

### **PASSO 6: Testar Bot**
1. Acesse: `https://t.me/direitolux_staging_bot`
2. Clique em "START"
3. Envie `/help`
4. **Ainda não haverá resposta** (normal - webhook não configurado)

---

## 🔧 Configuração no Código

### **Atualizar Token Real**
Depois de obter o token, atualize o arquivo:

```bash
# Editar arquivo
nano services/notification-service/.env

# Substituir linha:
TELEGRAM_BOT_TOKEN=SEU_TOKEN_REAL_AQUI
```

### **Testar Configuração**
```bash
# Executar teste
./test_telegram_bot.sh
```

---

## 📝 Exemplo de Token Real
```
# Formato correto:
TELEGRAM_BOT_TOKEN=1234567890:AAAA-BBBBccccDDDDeeeeFFFGGGhhhhIII

# ❌ Formato incorreto:
TELEGRAM_BOT_TOKEN=mock_telegram_token
```

---

## 🚀 Próximos Passos

1. ✅ **Criar bot** - Via BotFather
2. ✅ **Configurar token** - No arquivo .env
3. ✅ **Testar conexão** - Via script de teste
4. ⏳ **Configurar webhook** - HTTPS necessário
5. ⏳ **Testar integração** - Mensagens reais

---

## 💡 Dicas Importantes

- **Token é secreto** - Nunca compartilhe publicamente
- **Bot username** - Deve terminar com `_bot`
- **Webhook** - Precisa de HTTPS para funcionar
- **Rate limiting** - Telegram tem limite de 30 mensagens/minuto

---

**Status**: Aguardando criação do bot real via BotFather