# 🎯 PLANO TÉCNICO COMPLETO - ProcessAlert WhatsApp MicroSaaS

## 📋 VISÃO GERAL DO PRODUTO

**Nome**: ProcessAlert WhatsApp  
**Tagline**: "Nunca mais perca um prazo - Receba movimentos processuais no WhatsApp"  
**Mercado**: Advogados brasileiros (1M+ potenciais clientes)  
**Prazo**: 14 dias para MVP funcional  

---

## 🎯 CORE VALUE PROPOSITION

### **Problema Resolvido:**
- Advogados perdem prazos por não monitorar processos constantemente
- Emails são ignorados, sistemas jurídicos são complexos
- Consulta manual DataJud é trabalhosa e inconsistente

### **Solução Única:**
- **Monitor automático 24/7** de processos via DataJud CNJ
- **Notificação INSTANTÂNEA no WhatsApp** (canal preferido dos advogados)
- **Resumo IA** dos movimentos em linguagem simples
- **Setup em 2 minutos**: WhatsApp + Processo = Alertas automáticos

---

## 💰 MODELO DE NEGÓCIO

### **Planos de Assinatura:**
| Plano | Preço/mês | Processos | Consultas/dia | Target |
|-------|-----------|-----------|---------------|---------|
| **Starter** | R$ 29 | 5 processos | 50 consultas | Advogado solo |
| **Professional** | R$ 99 | 25 processos | 250 consultas | Escritório pequeno |
| **Business** | R$ 299 | 100 processos | 1000 consultas | Escritório médio |

### **Revenue Projetado:**
- **Mês 1**: 50 clientes × R$99 = R$4.950
- **Mês 3**: 200 clientes × R$99 = R$19.800  
- **Mês 6**: 500 clientes × R$99 = R$49.500
- **Ano 1**: 1000 clientes × R$99 = R$99.000/mês

---

## 🏗️ ARQUITETURA TÉCNICA SIMPLIFICADA

### **Stack Tecnológica:**
```
Frontend: Next.js 14 + Tailwind (Landing + Dashboard)
Backend: Go (3 microserviços mínimos)
Database: PostgreSQL (single instance)
Queue: Redis (simple pub/sub)
Deploy: Railway/Render (simplicidade máxima)
Payments: Stripe (direto, sem intermediários)
```

### **Microserviços Core (3 apenas):**

#### **1. Auth Service** (Port 8080)
```
Responsabilidades:
- Cadastro/Login advogados
- JWT tokens
- Planos/Billing via Stripe
- CRUD básico usuários

APIs:
POST /register → Cria conta + Stripe customer
POST /login → JWT token
GET /profile → Dados do usuário
POST /subscribe → Criar assinatura Stripe
```

#### **2. Monitor Service** (Port 8081)  
```
Responsabilidades:
- CRUD processos monitorados
- Polling DataJud (30/30min)
- Detect movimentos novos
- Enviar para fila notificação

APIs:
POST /processes → Adiciona processo ao monitor
GET /processes → Lista processos do usuário
DELETE /processes/:id → Remove monitoramento
GET /movements/:process → Histórico movimentos

Background Job:
- Cron 30min → Consulta DataJud todos processos ativos
- Compara com último movimento salvo
- Se novo movimento → Redis queue
```

#### **3. Notification Service** (Port 8082)
```
Responsabilidades:
- WhatsApp Business API integration
- IA resumo movimentos (OpenAI GPT-4)
- Fila de notificações (Redis)
- Rate limiting WhatsApp

APIs:
POST /notify → Enviar notificação manual
GET /notifications → Histórico notificações

Background Job:
- Redis consumer → Processa fila notificações
- OpenAI → Resume movimento jurídico
- WhatsApp → Envia: "🚨 Processo 123: [resumo IA]"
```

---

## 🗄️ BANCO DE DADOS (PostgreSQL)

### **Schema Simplificado (5 tabelas):**

```sql
-- Usuários/Advogados
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR UNIQUE NOT NULL,
    password_hash VARCHAR NOT NULL,
    whatsapp VARCHAR NOT NULL,
    plan VARCHAR DEFAULT 'starter',
    stripe_customer_id VARCHAR,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Processos Monitorados
CREATE TABLE monitored_processes (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    process_number VARCHAR NOT NULL,
    tribunal VARCHAR NOT NULL,
    last_movement_date TIMESTAMP,
    last_movement_text TEXT,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Movimentos Detectados
CREATE TABLE movements (
    id UUID PRIMARY KEY,
    process_id UUID REFERENCES monitored_processes(id),
    movement_date TIMESTAMP NOT NULL,
    movement_text TEXT NOT NULL,
    summary_ai TEXT, -- Resumo gerado pela IA
    detected_at TIMESTAMP DEFAULT NOW()
);

-- Notificações Enviadas
CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    process_id UUID REFERENCES monitored_processes(id),
    movement_id UUID REFERENCES movements(id),
    message TEXT NOT NULL,
    whatsapp_id VARCHAR, -- ID da mensagem WhatsApp
    status VARCHAR DEFAULT 'sent',
    sent_at TIMESTAMP DEFAULT NOW()
);

-- Quotas/Limites
CREATE TABLE quota_usage (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    date DATE NOT NULL,
    datajud_queries INTEGER DEFAULT 0,
    whatsapp_messages INTEGER DEFAULT 0,
    PRIMARY KEY (user_id, date)
);
```

---

## 🔄 FLUXO DE FUNCIONAMENTO

### **1. Onboarding (2 minutos):**
```
1. Advogado acessa landing page
2. Clica "Começar Grátis"
3. Preenche: Email, WhatsApp, Senha
4. Seleciona plano (14 dias trial gratuito)
5. Adiciona primeiro processo para monitorar
6. Recebe WhatsApp: "Olá! Processo 123 está sendo monitorado 24/7"
```

### **2. Monitor Automático:**
```
A cada 30 minutos:
1. Cron job consulta DataJud para todos processos ativos
2. Compara último movimento salvo vs movimento atual DataJud
3. Se detectou movimento novo:
   a) Salva movement na tabela movements
   b) Envia para Redis queue: {user_id, process_id, movement_id}
```

### **3. Notificação WhatsApp:**
```
Redis consumer processa fila:
1. Busca dados: user.whatsapp, process.number, movement.text
2. OpenAI resume movimento: "Nova petição juntada aos autos"
3. WhatsApp send: "🚨 Processo 12345: Nova petição juntada aos autos. Prazo: 15 dias para manifestação."
4. Salva notification enviada
```

### **4. Dashboard Simples:**
```
- Lista processos monitorados (ativo/pausado)
- Últimos 10 movimentos detectados
- Botão: "Adicionar Processo"
- Configurações: WhatsApp, plano, billing
```

---

## 🛠️ IMPLEMENTAÇÃO - CRONOGRAMA 14 DIAS

### **SEMANA 1: CORE BACKEND + MVP**

#### **Dias 1-2: Setup Projeto + Auth Service**
```bash
Estrutura:
direito-lux/
├── services/
│   ├── auth-service/     # Go service - JWT, Stripe, CRUD users
│   ├── monitor-service/  # Go service - DataJud polling, CRUD processes
│   └── notification-service/ # Go service - WhatsApp, OpenAI, Redis
├── frontend/            # Next.js 14 - Landing + Dashboard
├── database/           # PostgreSQL schema + migrations  
├── docker-compose.yml  # Local development
└── deploy/            # Railway/Render configs

Tarefas Dia 1-2:
✅ Setup Go modules (3 services)
✅ PostgreSQL schema + migrations
✅ Auth Service: JWT + Stripe integration
✅ Docker compose local
✅ Testes básicos auth
```

#### **Dias 3-4: Monitor Service + DataJud**
```bash
Tarefas Dia 3-4:
✅ CRUD processos monitorados
✅ DataJud HTTP client (real API CNJ)
✅ Background job polling (cron 30min)
✅ Detect novos movimentos
✅ Redis queue integration
✅ Testes DataJud + polling
```

#### **Dias 5-6: Notification Service + WhatsApp**
```bash
Tarefas Dia 5-6:
✅ WhatsApp Business API setup
✅ OpenAI GPT-4 integration (resumos)
✅ Redis consumer (fila notificações)
✅ Rate limiting WhatsApp
✅ Histórico notificações
✅ Testes notificação end-to-end
```

#### **Dia 7: Integration Testing**
```bash
Tarefas Dia 7:
✅ Teste fluxo completo: Registro → Monitor → Notificação
✅ Deploy staging (Railway)
✅ Testes com DataJud real
✅ Webhook WhatsApp funcionando
```

### **SEMANA 2: FRONTEND + LAUNCH**

#### **Dias 8-9: Frontend Next.js**
```bash
Páginas:
- Landing page (hero, features, pricing, testimonials)
- /register (onboarding 2 minutos)
- /login 
- /dashboard (lista processos, adicionar processo)
- /settings (whatsapp, plano, billing)

Tarefas Dia 8-9:
✅ Landing page conversão otimizada
✅ Dashboard funcional (CRUD processos)
✅ Stripe Checkout integration
✅ Design responsivo (mobile-first)
```

#### **Dias 10-11: Polish + Billing**
```bash
Tarefas Dia 10-11:
✅ Stripe webhooks (renovação, cancelamento)
✅ Trial 14 dias implementation
✅ Quotas enforcement (limites por plano)
✅ Email transacional (onboarding, billing)
✅ Analytics básico (Posthog)
```

#### **Dias 12-13: Launch Preparation**
```bash
Tarefas Dia 12-13:
✅ Landing page SEO otimizada
✅ Blog posts (3 artigos: Como funciona, Casos de uso, Comparativo)
✅ Deploy produção (Railway)
✅ Domínio + SSL (processalert.com.br)
✅ Testes finais produção
```

#### **Dia 14: LAUNCH + Marketing**
```bash
Tarefas Dia 14:
✅ Product Hunt launch
✅ LinkedIn posts (advocacia + tech)
✅ WhatsApp para rede pessoal
✅ Email para beta users
✅ Monitor primeiros signups
```

---

## 🔧 STACK TÉCNICO DETALHADO

### **Backend (Go):**
```go
// Principais packages:
- gin-gonic/gin          // HTTP router
- golang-jwt/jwt         // JWT tokens
- stripe/stripe-go       // Payments
- lib/pq                 // PostgreSQL
- go-redis/redis         // Redis
- go-resty/resty         // HTTP client DataJud
- sashabaranov/go-openai // OpenAI GPT-4
```

### **Frontend (Next.js 14):**
```javascript
// Stack:
- Next.js 14 (App Router)
- TypeScript
- Tailwind CSS
- Shadcn/ui (componentes)
- Zustand (estado global)
- React Hook Form (formulários)
- Stripe Elements (checkout)
```

### **Infrastructure:**
```yaml
# Deploy Railway (simplicidade máxima):
- PostgreSQL (Railway)
- Redis (Railway)  
- 3 Go services (Railway)
- Next.js frontend (Vercel)
- Domínio: processalert.com.br
```

---

## 🚀 APIS DETALHADAS

### **Auth Service (8080):**
```go
// Principais endpoints:
POST   /api/auth/register          // Registro + Stripe customer
POST   /api/auth/login             // Login + JWT
GET    /api/auth/profile           // Dados usuário
PUT    /api/auth/profile           // Update perfil
POST   /api/auth/subscribe         // Criar assinatura
POST   /api/auth/cancel            // Cancelar assinatura
GET    /api/auth/billing           // Status billing
POST   /api/auth/webhook/stripe    // Stripe webhooks
```

### **Monitor Service (8081):**
```go
// Principais endpoints:
GET    /api/processes              // Lista processos usuário
POST   /api/processes              // Adiciona processo
PUT    /api/processes/:id          // Update processo (pause/resume)
DELETE /api/processes/:id          // Remove processo
GET    /api/processes/:id/movements // Histórico movimentos
POST   /api/processes/validate     // Valida número processo
GET    /api/quota/usage            // Uso atual quotas
```

### **Notification Service (8082):**
```go
// Principais endpoints:
GET    /api/notifications          // Histórico notificações
POST   /api/notifications/test     // Teste notificação manual
PUT    /api/notifications/settings // Config WhatsApp
GET    /api/notifications/status   // Status WhatsApp API
POST   /api/webhook/whatsapp       // WhatsApp webhooks
```

---

## 📱 INTEGRAÇÕES EXTERNAS

### **1. DataJud CNJ API:**
```
Base URL: https://api-publica.datajud.cnj.jus.br
Rate Limit: 120 requests/minuto
Authentication: API Key + Certificado A1/A3

Endpoint principal:
GET /api_publica_v2/processo/{numero_processo}

Response:
{
  "hits": [{
    "numeroProcesso": "1234567-89.2023.8.26.0100",
    "dataHora": "2025-01-15T10:30:00Z",
    "movimento": {
      "descricao": "Juntada de petição"
    }
  }]
}
```

### **2. WhatsApp Business API:**
```
Provider: Meta Business (oficial)
Rate Limit: 1000 mensagens/dia (plano inicial)
Webhook URL: https://processalert.com.br/api/webhook/whatsapp

Envio mensagem:
POST https://graph.facebook.com/v18.0/{phone_number_id}/messages
{
  "to": "5511999999999",
  "type": "text",
  "text": {
    "body": "🚨 Processo 123: Nova petição juntada aos autos"
  }
}
```

### **3. OpenAI GPT-4:**
```
Model: gpt-4-turbo
Max tokens: 150 (resumos)
Cost: ~$0.03 por resumo

Prompt template:
"Resume este movimento processual em até 30 palavras, explicando de forma simples o que aconteceu e se há prazo para resposta:

Movimento: {movimento_texto}

Formato: [Ação] - [Prazo se houver]"

Exemplo output:
"Nova petição juntada aos autos - Prazo 15 dias para manifestação"
```

### **4. Stripe Payments:**
```
Products:
- starter_plan (R$ 29/mês)
- professional_plan (R$ 99/mês) 
- business_plan (R$ 299/mês)

Webhooks importantes:
- customer.subscription.created
- customer.subscription.updated  
- customer.subscription.deleted
- invoice.payment_succeeded
- invoice.payment_failed
```

---

## 🎨 UX/UI WIREFRAMES

### **Landing Page:**
```
Header: ProcessAlert | Nunca perca um prazo jurídico

Hero Section:
[Smartphone com WhatsApp] 
"Receba movimentos processuais no WhatsApp
em tempo real com resumo IA"

CTA: "Começar Grátis - 14 dias trial"

Features:
✅ Monitor 24/7 DataJud CNJ
✅ WhatsApp instantâneo  
✅ Resumo IA inteligente
✅ Setup em 2 minutos

Pricing:
[3 cards com planos]

Testimonials:
"Nunca mais perdi um prazo desde que uso ProcessAlert"
- Dr. João Silva, OAB/SP

Social Proof:
"500+ advogados monitoram 2000+ processos"
```

### **Dashboard:**
```
Sidebar:
- 📊 Dashboard
- ⚖️ Processos (5/25)
- 🔔 Notificações
- ⚙️ Configurações
- 💳 Plano

Main:
"Processos Monitorados"

[+ Adicionar Processo]

Table:
| Nº Processo | Tribunal | Último Movimento | Status |
|-------------|----------|------------------|---------|
| 1234567-89  | TJSP     | 12/01 - Petição  | 🟢 Ativo |
| 7654321-10  | TJRJ     | 10/01 - Despacho | 🟢 Ativo |

"Últimas Notificações"
📱 12/01 15:30 - Processo 1234567: Nova petição juntada
📱 10/01 09:15 - Processo 7654321: Despacho publicado
```

---

## 📊 MÉTRICAS & KPIs

### **Métricas de Produto:**
- **Signup Rate**: Landing → Registro (target: 3%)
- **Activation Rate**: Registro → Primeiro processo (target: 80%)
- **Retention D7**: Usuários ativos 7 dias (target: 70%)
- **Churn Rate**: Cancelamentos mensais (target: <5%)
- **NPS**: Net Promoter Score (target: 50+)

### **Métricas Técnicas:**
- **DataJud Success Rate**: Consultas bem-sucedidas (target: >95%)
- **WhatsApp Delivery Rate**: Mensagens entregues (target: >98%)  
- **Response Time**: APIs < 200ms (target: >95%)
- **Uptime**: Disponibilidade sistema (target: 99.9%)

### **Métricas de Negócio:**
- **MRR**: Monthly Recurring Revenue
- **CAC**: Customer Acquisition Cost (target: <R$50)
- **LTV**: Customer Lifetime Value (target: >R$1000)
- **LTV/CAC Ratio**: target >20x

---

## 🔒 SEGURANÇA & COMPLIANCE

### **Dados Jurídicos (LGPD):**
- Criptografia AES-256 dados sensíveis
- Logs auditoria todas ações usuário
- Backup daily com retenção 30 dias
- DPO contact: privacy@processalert.com.br

### **Autenticação:**
- JWT tokens (15min lifetime)
- Refresh tokens (7 dias)
- Rate limiting: 100 requests/min per IP
- 2FA via WhatsApp (futuro)

### **API Security:**
- HTTPS obrigatório (TLS 1.3)
- CORS configurado
- Input validation todas APIs
- SQL injection protection

---

## 💰 PROJEÇÃO FINANCEIRA

### **Custos Mensais:**
```
Infrastructure (Railway):     $50
WhatsApp Business API:        $100  
OpenAI GPT-4:                $200
Stripe fees (3%):            $150
Domínio + SSL:               $10
Total:                       $510/mês
```

### **Break-even:**
```
Custo: R$510/mês = ~R$2.550 (USD 5.0 → BRL)
Receita por cliente: R$99/mês  
Break-even: 26 clientes pagantes

Meta mês 1: 50 clientes = R$4.950 - R$2.550 = R$2.400 lucro
Meta mês 3: 200 clientes = R$19.800 - R$2.550 = R$17.250 lucro
Meta mês 6: 500 clientes = R$49.500 - R$2.550 = R$46.950 lucro
```

### **Projeção 12 meses:**
| Mês | Clientes | MRR | Custos | Lucro |
|-----|----------|-----|--------|--------|
| 1   | 50       | R$4.950 | R$2.550 | R$2.400 |
| 3   | 200      | R$19.800 | R$2.550 | R$17.250 |
| 6   | 500      | R$49.500 | R$2.550 | R$46.950 |
| 12  | 1000     | R$99.000 | R$3.000 | R$96.000 |

---

## 🚀 GO-TO-MARKET STRATEGY

### **Pré-Launch (Dias 1-13):**
- Build in public (LinkedIn daily updates)
- Landing page + email capture
- Beta users (50 advogados da rede)

### **Launch Day (Dia 14):**
- Product Hunt launch
- LinkedIn viral post + demo video
- WhatsApp personal network
- Email beta users

### **Pós-Launch (Semanas 3-4):**
- Content marketing (blog SEO)
- Google Ads (keywords: "monitorar processo", "prazo processual")
- Parcerias (grupos WhatsApp advogados)
- Referral program (1 mês grátis)

### **Growth Channels:**
1. **SEO Content** (50% traffic target)
2. **LinkedIn** (25% traffic)  
3. **Google Ads** (15% traffic)
4. **Referrals** (10% traffic)

---

## ✅ DEFINITION OF DONE

### **MVP Ready Criteria:**
- [ ] 3 serviços Go deployados e funcionais
- [ ] DataJud integration com API real CNJ
- [ ] WhatsApp notificação end-to-end
- [ ] OpenAI resumos funcionando
- [ ] Frontend dashboard completo
- [ ] Stripe billing integration
- [ ] Landing page otimizada conversão
- [ ] Deploy produção estável
- [ ] 10 beta users testando com sucesso
- [ ] Documentação básica (API + onboarding)

### **Launch Ready Criteria:**
- [ ] 50 beta users feedback positivo
- [ ] Métricas básicas implementadas
- [ ] Support básico (email + WhatsApp)
- [ ] Blog com 3 artigos publicados
- [ ] SEO básico otimizado
- [ ] Product Hunt page preparada
- [ ] Demo video gravado
- [ ] Primeiros clientes pagantes (5+)

---

## 🎯 SUCCESS METRICS

### **14 dias (MVP):**
- ✅ Sistema funcionando end-to-end
- ✅ 10 beta users ativos
- ✅ 100 processos monitorados
- ✅ 50 notificações WhatsApp enviadas

### **30 dias:**
- 🎯 50 clientes pagantes
- 🎯 R$4.950 MRR
- 🎯 500 processos monitorados
- 🎯 95% uptime

### **90 dias:**
- 🎯 200 clientes pagantes
- 🎯 R$19.800 MRR
- 🎯 2000 processos monitorados  
- 🎯 Break-even alcançado

---

**🚀 ESTE PLANO É EXECUTÁVEL EM 14 DIAS COM FOCO LASER NO ESSENCIAL!**

Está pronto para começarmos a implementação? Podemos ajustar qualquer detalhe antes de partir para o código!