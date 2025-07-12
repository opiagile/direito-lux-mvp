# 🤖 CRIAR BOT TELEGRAM - GUIA VISUAL PASSO A PASSO

## 📱 PASSO 1: Abrir Telegram e Encontrar BotFather

1. **Abra o Telegram** (mobile ou desktop)
2. **Busque por**: `@BotFather`
3. **Procure pelo bot oficial** com verificado ✅
4. **Clique em "START"** para iniciar conversa

---

## 🔨 PASSO 2: Criar Novo Bot

### 2.1 Envie o comando:
```
/newbot
```

### 2.2 BotFather responderá:
```
Alright, a new bot. How are we going to call it? 
Please choose a name for your bot.
```

### 2.3 Digite o nome do bot:
```
Direito Lux Staging Bot
```

### 2.4 BotFather perguntará:
```
Good. Now let's choose a username for your bot. 
It must end in `bot`. Like this, for example: TetrisBot or tetris_bot.
```

### 2.5 Digite o username:
```
direitolux_staging_bot
```

---

## 🎉 PASSO 3: Copiar Token

### 3.1 BotFather responderá com:
```
Done! Congratulations on your new bot. You will find it at 
t.me/direitolux_staging_bot. You can now add a description, 
about section and profile picture for your bot, see /help 
for a list of commands. By the way, when you've finished 
creating your cool bot, ping our Bot Support if you want a 
better username for it. Just make sure the bot is fully 
operational before you do this.

Use this token to access the HTTP API:
7458394857:AAHKz9XjB8vK_2QxYz0-fG8kNvM_xQz7890

Keep your token secure and store it safely, it can be used 
by anyone to control your bot.

For a description of the Bot API, see this page:
https://core.telegram.org/bots/api
```

### 3.2 **COPIE O TOKEN!** 
Exemplo: `7458394857:AAHKz9XjB8vK_2QxYz0-fG8kNvM_xQz7890`

---

## ⚙️ PASSO 4: Configurar Comandos do Bot

### 4.1 Envie:
```
/setcommands
```

### 4.2 Selecione seu bot:
```
@direitolux_staging_bot
```

### 4.3 Copie e cole EXATAMENTE este texto:
```
start - 🚀 Iniciar conversa com o bot
help - 🆘 Ajuda e comandos disponíveis
status - 📊 Status dos seus processos
agenda - 📅 Agenda de prazos importantes
busca - 🔍 Buscar processos jurídicos
relatorio - 📈 Relatório rápido dos processos
configurar - ⚙️ Configurações do bot
```

### 4.4 BotFather confirmará:
```
Success! Command list updated. /help
```

---

## 📝 PASSO 5: Configurar Descrição

### 5.1 Envie:
```
/setdescription
```

### 5.2 Selecione seu bot:
```
@direitolux_staging_bot
```

### 5.3 Copie e cole:
```
🏛️ Direito Lux - Assistente Jurídico Inteligente

Bot oficial para monitoramento de processos, notificações automáticas e análises jurídicas.

✅ Ambiente STAGING - Para testes e validação
```

### 5.4 BotFather confirmará:
```
Success! Description updated.
```

---

## 📸 PASSO 6: Adicionar Foto de Perfil (Opcional)

### 6.1 Envie:
```
/setuserpic
```

### 6.2 Selecione seu bot:
```
@direitolux_staging_bot
```

### 6.3 Envie uma imagem (logo do Direito Lux)

---

## ✅ PASSO 7: Testar o Bot

1. **Acesse**: https://t.me/direitolux_staging_bot
2. **Clique em "START"**
3. **Envie**: `/help`
4. **O bot ainda não responderá** (normal - webhook não configurado)

---

## 🔑 PASSO 8: Salvar Token no Sistema

### 8.1 Copie seu token real:
```
7458394857:AAHKz9XjB8vK_2QxYz0-fG8kNvM_xQz7890
```

### 8.2 Cole no campo abaixo quando executar o script:
```bash
./configure_telegram_bot.sh
```

---

## 📋 CHECKLIST FINAL

- [ ] Bot criado com sucesso
- [ ] Token copiado e salvo
- [ ] Comandos configurados
- [ ] Descrição adicionada
- [ ] Username: @direitolux_staging_bot
- [ ] Link: https://t.me/direitolux_staging_bot

---

## 🚨 DICAS IMPORTANTES

1. **NUNCA compartilhe o token publicamente**
2. **Salve o token em local seguro**
3. **Se o username estiver em uso**, tente:
   - `direitolux_staging_2025_bot`
   - `direitolux_dev_bot`
   - `direitolux_test_bot`

4. **Se errar algo**, use:
   - `/deletebot` - para deletar e começar de novo
   - `/mybots` - para ver seus bots

---

## 🎯 PRONTO!

Agora você tem:
- ✅ Bot criado
- ✅ Token em mãos
- ✅ Comandos configurados
- ✅ Pronto para integração

**Próximo passo**: Execute o script de configuração com seu token!