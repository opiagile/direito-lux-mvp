# 💰 CONFIGURAR GATEWAYS DE PAGAMENTO - DIREITO LUX

## 🎯 OBJETIVO
Configurar ASAAS (PIX/Cartão) + NOWPayments (Crypto) para monetização completa do SaaS.

---

## 🏦 ASAAS (PIX + CARTÃO DE CRÉDITO)

### 📋 **Passo a Passo ASAAS**

#### 1️⃣ **Criar Conta**
1. **Acesse**: https://www.asaas.com/
2. **Cadastro gratuito** → Pessoa Jurídica
3. **Dados da empresa**: Direito Lux Tecnologia
4. **CNPJ**: (se tiver) ou Pessoa Física inicial
5. **Ativação**: 24-48h para aprovação

#### 2️⃣ **Configurar API**
1. **Dashboard** → **Integrações** → **API**
2. **Gerar API Key**: Ambiente de Sandbox primeiro
3. **Copiar chaves**:
   - Sandbox: `$aact_YTU5YTE0M2Jxxxxxxxxxxxxxxxx`
   - Produção: `$aact_YTU5YTE0M2Jxxxxxxxxxxxxxxxx`

#### 3️⃣ **Webhook Configuration**
```
Webhook URL: https://direito-lux-staging.loca.lt/webhook/asaas
Events: payment_created, payment_confirmed, payment_overdue
```

### 💳 **Produtos ASAAS Disponíveis**
- ✅ **PIX**: Instantâneo, taxa 0,99%
- ✅ **Cartão de Crédito**: Taxa 4,99% + R$0,39
- ✅ **Boleto**: Taxa R$3,49
- ✅ **Transferência**: Taxa R$3,49

---

## ₿ NOWPAYMENTS (CRYPTO)

### 📋 **Passo a Passo NOWPayments**

#### 1️⃣ **Criar Conta**
1. **Acesse**: https://nowpayments.io/
2. **Sign Up** → Business Account
3. **KYC**: Upload de documentos
4. **Aprovação**: 2-5 dias úteis

#### 2️⃣ **Configurar API**
1. **Dashboard** → **Settings** → **API**
2. **Generate API Key**: 
   - Sandbox: `NP-xxxxxxxxxxxxxxxx`
   - Production: `NP-xxxxxxxxxxxxxxxx`

#### 3️⃣ **Webhook Configuration**
```
Webhook URL: https://direito-lux-staging.loca.lt/webhook/nowpayments
Events: payment_waiting, payment_confirmed, payment_failed
```

### ₿ **Cryptos Suportadas**
- ✅ **Bitcoin (BTC)**
- ✅ **Ethereum (ETH)**  
- ✅ **USDT/USDC** (Stablecoins)
- ✅ **PIX via Binance Pay**
- ✅ **+150 altcoins**

---

## 💰 PLANOS DE ASSINATURA

### 💵 **Preços Definidos**
```json
{
  "starter": {
    "price": 99.00,
    "currency": "BRL",
    "interval": "monthly",
    "features": ["50 processos", "20 clientes", "100 consultas/dia"]
  },
  "professional": {
    "price": 299.00,
    "currency": "BRL", 
    "interval": "monthly",
    "features": ["200 processos", "100 clientes", "500 consultas/dia"]
  },
  "business": {
    "price": 699.00,
    "currency": "BRL",
    "interval": "monthly", 
    "features": ["500 processos", "500 clientes", "2000 consultas/dia"]
  },
  "enterprise": {
    "price": 1999.00,
    "currency": "BRL",
    "interval": "monthly",
    "features": ["Ilimitado", "Ilimitado", "10k consultas/dia", "White-label"]
  }
}
```

---

## 🔧 IMPLEMENTAÇÃO TÉCNICA

### 📱 **Frontend - Checkout**
```javascript
// components/checkout/PaymentMethod.tsx
const PaymentMethods = {
  PIX: 'asaas_pix',
  CREDIT_CARD: 'asaas_card', 
  CRYPTO: 'nowpayments_crypto',
  BOLETO: 'asaas_boleto'
}
```

### 🔙 **Backend - Billing Service**
```go
// services/billing-service/internal/handlers/payment_handler.go
type PaymentGateway interface {
    CreatePayment(amount float64, method string) (*Payment, error)
    ProcessWebhook(payload []byte) error
}

type AsaasGateway struct {
    apiKey string
    baseURL string
}

type NOWPaymentsGateway struct {
    apiKey string
    baseURL string
}
```

### 🌐 **Webhooks**
```go
// Webhook handlers for payment confirmations
func (h *BillingHandler) HandleAsaasWebhook(c *gin.Context) {
    // Process ASAAS payment confirmation
}

func (h *BillingHandler) HandleNOWPaymentsWebhook(c *gin.Context) {
    // Process crypto payment confirmation
}
```

---

## 🚀 CONFIGURAÇÃO RÁPIDA

### **Script Automático**
```bash
# Executar configuração completa
./setup_payment_gateways.sh

# Ou manual:
# 1. Configure ASAAS
# 2. Configure NOWPayments  
# 3. Teste webhooks
# 4. Deploy billing service
```

### **GitHub Secrets Necessários**
```bash
ASAAS_API_KEY_SANDBOX=seu_token_sandbox
ASAAS_API_KEY_PRODUCTION=seu_token_producao
NOWPAYMENTS_API_KEY_SANDBOX=seu_token_sandbox
NOWPAYMENTS_API_KEY_PRODUCTION=seu_token_producao
```

---

## 📊 FLUXO DE PAGAMENTO

### **1. Cliente escolhe plano**
```
Frontend → Billing Service → Gateway escolhido
```

### **2. Processamento**
```
ASAAS (PIX/Cartão) OU NOWPayments (Crypto)
```

### **3. Confirmação via webhook**
```
Gateway → Webhook → Billing Service → Ativação
```

### **4. Ativação automática**
```
Billing → Tenant Service → Process Service → Ativo
```

---

## 🧪 TESTES DE INTEGRAÇÃO

### **ASAAS Sandbox**
```bash
# Testar criação de cobrança PIX
curl -X POST https://sandbox.asaas.com/api/v3/payments \
  -H "access_token: $ASAAS_SANDBOX_KEY" \
  -d '{
    "customer": "cus_G7Dvo8w8IsP4",
    "billingType": "PIX", 
    "dueDate": "2025-07-15",
    "value": 99.00,
    "description": "Direito Lux - Plano Starter"
  }'
```

### **NOWPayments Sandbox**
```bash
# Testar pagamento crypto
curl -X POST https://api-sandbox.nowpayments.io/v1/payment \
  -H "x-api-key: $NOWPAYMENTS_SANDBOX_KEY" \
  -d '{
    "price_amount": 99.00,
    "price_currency": "BRL",
    "pay_currency": "btc",
    "order_id": "direito-lux-001",
    "order_description": "Direito Lux - Plano Starter"
  }'
```

---

## 📈 MÉTRICAS E ANALYTICS

### **Dashboard de Pagamentos**
- 💰 Receita mensal por gateway
- 📊 Conversão por método de pagamento
- 🔄 Taxa de renovação por plano
- 💳 Chargebacks e disputas

### **Relatórios Financeiros**
- 📈 MRR (Monthly Recurring Revenue)
- 📊 Churn rate por plano
- 💰 LTV (Customer Lifetime Value)
- 🎯 CAC (Customer Acquisition Cost)

---

## 🔒 SEGURANÇA E COMPLIANCE

### **LGPD**
- ✅ Dados de pagamento não armazenados
- ✅ Tokenização via gateways
- ✅ Logs auditáveis
- ✅ Consentimento explícito

### **PCI DSS**
- ✅ ASAAS e NOWPayments são certificados
- ✅ Não processamos cartões diretamente
- ✅ HTTPS obrigatório
- ✅ Webhooks autenticados

---

## 📞 SUPORTE

### **ASAAS**
- 📧 Email: suporte@asaas.com
- 📱 WhatsApp: (16) 3003-4200
- 🌐 Docs: https://docs.asaas.com/

### **NOWPayments**  
- 📧 Email: support@nowpayments.io
- 💬 Telegram: @nowpaymentsio
- 🌐 Docs: https://documenter.getpostman.com/view/7907941/

---

## 🎯 PRÓXIMOS PASSOS

1. ✅ **HOJE**: Criar contas ASAAS + NOWPayments
2. 🔧 **2 DIAS**: Implementar billing service completo
3. 🧪 **3 DIAS**: Testes de integração E2E
4. 🚀 **1 SEMANA**: Deploy em produção com pagamentos reais

---

**💰 RESULTADO**: Sistema de pagamento completo com PIX, cartão, boleto e crypto para maximizar conversão!