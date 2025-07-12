# 🔐 CONFIGURAR GITHUB SECRETS - DIREITO LUX

## 🎯 INSTRUÇÕES RÁPIDAS (5 MINUTOS)

### 1️⃣ **Acessar Repository Settings**
1. Vá para: https://github.com/SEU_USUARIO/direito-lux
2. **Settings** → **Secrets and variables** → **Actions**
3. Clique em **"New repository secret"**

### 2️⃣ **Configurar Secrets Obrigatórios**

#### 🤖 **TELEGRAM BOT**
```
Name: TELEGRAM_BOT_TOKEN
Value: 7927061803:AAGC5GMerAe9CVegcl85o6BTFj2hqkcjO04
```

#### 📱 **WHATSAPP API** (Quando obtiver)
```
Name: WHATSAPP_ACCESS_TOKEN
Value: EAAxxxxxxxxxxxxxxxx

Name: WHATSAPP_PHONE_NUMBER_ID  
Value: 123456789012345

Name: WHATSAPP_BUSINESS_ACCOUNT_ID
Value: 123456789012345
```

#### 🤖 **AI SERVICES**
```
Name: OPENAI_API_KEY
Value: sk-xxxxxxxxxxxxxxxx

Name: ANTHROPIC_API_KEY
Value: sk-ant-xxxxxxxxxxxxxxxx
```

#### 💰 **PAYMENT GATEWAYS**
```
Name: ASAAS_API_KEY
Value: $aact_YTU5YTE0M2Jxxxxxxxxxxxxxxxx

Name: NOWPAYMENTS_API_KEY
Value: NP-xxxxxxxxxxxxxxxx
```

#### 📧 **EMAIL**
```
Name: SMTP_PASSWORD
Value: sua_senha_email_aqui
```

#### 🏛️ **DATAJUD CNJ**
```
Name: DATAJUD_API_KEY
Value: sua_chave_cnj_aqui

Name: DATAJUD_CERTIFICATE_PASSWORD
Value: senha_certificado_digital
```

#### 🗄️ **DATABASE**
```
Name: DB_PASSWORD
Value: direito_lux_pass_production

Name: RABBITMQ_PASSWORD
Value: direito_lux_rabbit_pass
```

---

## 🔄 **WORKFLOW AUTOMÁTICO CRIADO**

✅ **Arquivo**: `.github/workflows/deploy-with-secrets.yml`

**Funcionalidades**:
- 🔐 Validação de secrets obrigatórios
- 🧪 Testes automatizados dos microserviços
- 🐳 Build com Docker usando secrets
- 📱 Teste de integração Telegram/WhatsApp
- 🔍 Auditoria de segurança automática
- 🏥 Health checks do ambiente staging

---

## 🚀 **COMO USAR**

### **Automático (Recomendado)**
1. Configure os secrets acima no GitHub
2. Push para `main` → Deploy automático
3. Pull Request → Testes automáticos

### **Manual**
```bash
# Testar localmente com secrets
export TELEGRAM_BOT_TOKEN="seu_token"
export WHATSAPP_ACCESS_TOKEN="seu_token"

# Executar serviços
docker-compose up -d
```

---

## ✅ **VANTAGENS GITHUB SECRETS**

### **🔒 Segurança**
- Criptografia AES-256-GCM
- Mascaramento automático nos logs
- Acesso controlado por permissões

### **🚀 Produção Ready**
- Zero configuração adicional
- Integração nativa CI/CD
- Auditoria automática

### **💰 Custo**
- Totalmente gratuito
- Incluído no GitHub

---

## 🔍 **VALIDAÇÃO**

### **Verificar Secrets Configurados**
1. Repository → Settings → Secrets
2. Deve ver todos os secrets listados
3. ✅ = Configurado | ❌ = Faltando

### **Testar Deploy**
1. Push qualquer mudança para `main`
2. Actions → Ver execução do workflow
3. Verificar se todos os jobs passaram

### **Testar Staging**
1. Acessar: https://direito-lux-staging.loca.lt
2. Testar webhook Telegram: `@direitolux_staging_bot`
3. Verificar logs: `docker-compose logs -f notification-service`

---

## 🆘 **TROUBLESHOOTING**

### **Secret não funciona**
```bash
# Verificar se está mascarado nos logs
echo "Token: $TELEGRAM_BOT_TOKEN"  # Deve aparecer ***
```

### **Workflow falhando**
1. Actions → Ver logs detalhados
2. Procurar por "❌" nos steps
3. Configurar secret faltante

### **Webhook não responde**
1. Verificar se TELEGRAM_BOT_TOKEN está correto
2. Testar localmente primeiro
3. Verificar logs do container

---

## 🎯 **PRÓXIMOS PASSOS**

1. ✅ **Configure secrets críticos**: TELEGRAM_BOT_TOKEN
2. 🔄 **Push para main**: Ativar workflow automático  
3. 📱 **Teste WhatsApp**: Adicionar tokens quando obtiver
4. 💰 **Configurar pagamento**: ASAAS + NOWPayments
5. 🏛️ **DataJud real**: Certificado digital + API key

---

## 📞 **SUPORTE**

- **Email**: contato@direitolux.com.br
- **Bot**: @direitolux_staging_bot  
- **Docs**: SECRETS_MANAGEMENT_OPTIONS.md

---

**🔐 LEMBRE-SE**: 
- Nunca commitar secrets no código
- Usar sempre GitHub Secrets em produção
- Rotacionar tokens periodicamente