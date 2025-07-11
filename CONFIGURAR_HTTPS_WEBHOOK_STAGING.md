# 🌐 CONFIGURAR HTTPS WEBHOOK STAGING - OPÇÕES

## ⚡ OPÇÃO A: LOCALTUNNEL (RÁPIDA - SEM AUTH)

### Instalar e usar localtunnel:
```bash
# Instalar localtunnel (alternativa ao ngrok)
npm install -g localtunnel

# Criar túnel para notification service (porta 8085)
lt --port 8085 --subdomain direito-lux-staging

# Resultado esperado:
# your url is: https://direito-lux-staging.loca.lt
```

### Testar túnel:
```bash
# Testar se está funcionando
curl https://direito-lux-staging.loca.lt/health
```

---

## 🔧 OPÇÃO B: NGROK (PRECISA AUTH)

### 1. Criar conta ngrok:
1. Acesse: https://dashboard.ngrok.com/signup
2. Crie conta gratuita
3. Confirme email

### 2. Obter authtoken:
1. Acesse: https://dashboard.ngrok.com/get-started/your-authtoken
2. Copie o authtoken
3. Configure: `ngrok config add-authtoken SEU_TOKEN`

### 3. Criar túnel:
```bash
ngrok http 8085
```

---

## 🚀 OPÇÃO C: CLOUDFLARE TUNNEL (AVANÇADA)

### 1. Instalar cloudflared:
```bash
brew install cloudflare/cloudflare/cloudflared
```

### 2. Criar túnel:
```bash
cloudflared tunnel --url http://localhost:8085
```

---

## ✅ RESULTADO ESPERADO

Qualquer opção deve retornar URL HTTPS tipo:
- `https://direito-lux-staging.loca.lt` (localtunnel)
- `https://abc123.ngrok.io` (ngrok)
- `https://xyz.trycloudflare.com` (cloudflare)

## 🔗 PRÓXIMOS PASSOS

1. ✅ Obter URL HTTPS pública
2. ⏳ Configurar webhook Telegram
3. ⏳ Configurar webhook WhatsApp
4. ⏳ Testar webhooks funcionando

**Recomendação**: Use localtunnel por ser mais rápido para staging!