# 🚀 SETUP INICIAL - AMBIENTE GCP OTIMIZADO

## 📋 QUANDO USAR ESTE GUIA

- ✅ Primeira vez configurando o ambiente GCP
- ✅ Após deletar cluster para economizar
- ✅ Configurando novo projeto GCP
- ✅ Migrando de desenvolvimento para staging

---

## 🎯 ESCOLHA SUA ESTRATÉGIA

### **1. 🟢 CLOUD RUN (RECOMENDADO)**
**Custo:** R$30/mês | **Economia:** 98% | **Manutenção:** Zero

```bash
# Setup completo Cloud Run
./scripts/migrate-to-cloud-run.sh setup-cloudrun
```

### **2. 🟡 GKE COM AUTO-SHUTDOWN**
**Custo:** R$450/mês | **Economia:** 83% | **Manutenção:** Baixa

```bash
# Setup GKE otimizado
./scripts/gcp-cost-optimizer.sh optimize
./scripts/setup-auto-shutdown.sh setup
```

### **3. 🔴 GKE MANUAL**
**Custo:** Variável | **Economia:** 50-90% | **Manutenção:** Alta

```bash
# Setup básico - você controla
./scripts/gcp-cost-optimizer.sh start
```

---

## 🛠️ SETUP PASSO A PASSO

### **PASSO 1: Verificar Ambiente**
```bash
# Verificar se tem cluster
gcloud container clusters list --project=direito-lux-staging-2025

# Verificar custos atuais
./scripts/gcp-cost-optimizer.sh costs
```

### **PASSO 2: Decidir Estratégia**

#### **Se custos estão altos (>R$50/dia):**
```bash
# EMERGÊNCIA - Parar tudo
./scripts/migrate-to-cloud-run.sh emergency
```

#### **Se não tem cluster:**
```bash
# Criar cluster mínimo
gcloud container clusters create direito-lux-gke-staging \
  --region=us-central1 \
  --project=direito-lux-staging-2025 \
  --num-nodes=1 \
  --machine-type=e2-small \
  --enable-autoscaling \
  --min-nodes=0 \
  --max-nodes=3
```

### **PASSO 3: Configurar Database**
```bash
# Configurar port-forward para PostgreSQL
kubectl port-forward -n direito-lux-staging svc/postgres-service 5432:5432 &

# Executar migrations
./scripts/setup-staging-database.sh
```

### **PASSO 4: Configurar Economia**

#### **Para Cloud Run:**
```bash
./scripts/migrate-to-cloud-run.sh setup-cloudrun
```

#### **Para GKE com auto-shutdown:**
```bash
./scripts/setup-auto-shutdown.sh setup
```

#### **Para GKE manual:**
```bash
./scripts/gcp-cost-optimizer.sh optimize
```

---

## 🧪 TESTES APÓS SETUP

### **Teste 1: Sistema funcionando**
```bash
# Health check
curl -k https://35.188.198.87/api/health

# Resposta esperada: {"status":"healthy","timestamp":"..."}
```

### **Teste 2: Auth Service**
```bash
# Login
curl -k -X POST https://35.188.198.87/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: 550e8400-e29b-41d4-a716-446655440001" \
  -d '{"email":"admin@silvaassociados.com.br","password":"password"}'

# Resposta esperada: {"access_token":"...","user":{...}}
```

### **Teste 3: Tenant Service**
```bash
# Dados do tenant
curl -k https://35.188.198.87/api/v1/tenants/550e8400-e29b-41d4-a716-446655440001 \
  -H "X-Tenant-ID: 550e8400-e29b-41d4-a716-446655440001"

# Resposta esperada: {"data":{"name":"Silva & Associados",...}}
```

### **Teste 4: Frontend**
```bash
# Acessar no browser
open https://35.188.198.87

# Fazer login com:
# Email: admin@silvaassociados.com.br
# Senha: password
```

---

## 📊 CONFIGURAÇÕES ESPECÍFICAS

### **🟢 CLOUD RUN SETUP**

**Benefícios:**
- Escala para zero automaticamente
- Só paga quando há tráfego
- Sem gerenciamento de nodes

**Limitações:**
- Tempo de cold start (2-3 segundos)
- Melhor para baixo tráfego

**Quando usar:**
- Staging/Demo
- Desenvolvimento individual
- Tráfego esporádico

### **🟡 GKE AUTO-SHUTDOWN SETUP**

**Benefícios:**
- Performance constante
- Controle total
- Para automaticamente

**Limitações:**
- Precisa gerenciar nodes
- Tempo de startup (2-3 minutos)

**Quando usar:**
- Desenvolvimento em equipe
- Tráfego previsível
- Precisa de performance

### **🔴 GKE MANUAL SETUP**

**Benefícios:**
- Controle total
- Performance máxima
- Configuração customizada

**Limitações:**
- Precisa ligar/desligar manualmente
- Custo alto se esquecer ligado

**Quando usar:**
- Produção
- Desenvolvimento intensivo
- Tráfego alto

---

## 🔧 CONFIGURAÇÕES AVANÇADAS

### **Configurar Domínio Custom:**
```bash
# Apontar DNS para IP do load balancer
# staging.direitolux.com.br → 35.188.198.87
```

### **Configurar SSL/TLS:**
```bash
# Já configurado automaticamente com cert-manager
# Certificados Let's Encrypt renovam automaticamente
```

### **Configurar Monitoramento:**
```bash
# Prometheus já deployado
# Acesso: kubectl port-forward -n direito-lux-staging svc/prometheus 9090:9090
```

---

## 🚨 TROUBLESHOOTING INICIAL

### **Problema: Cluster não existe**
```bash
# Criar cluster
gcloud container clusters create direito-lux-gke-staging \
  --region=us-central1 \
  --project=direito-lux-staging-2025 \
  --num-nodes=1 \
  --machine-type=e2-small
```

### **Problema: Database sem dados**
```bash
# Executar setup completo
./scripts/setup-staging-database.sh
```

### **Problema: Ingress não funciona**
```bash
# Verificar ingress
kubectl get ingress -n direito-lux-staging

# Aplicar ingress correto
kubectl apply -f ingress-simples-apis.yaml
```

### **Problema: Custos altos**
```bash
# Verificar recursos
./scripts/gcp-cost-optimizer.sh costs

# Otimizar imediatamente
./scripts/gcp-cost-optimizer.sh optimize
```

---

## 📋 CHECKLIST FINAL

### **Após Setup Completo:**
- [ ] Sistema responde em https://35.188.198.87
- [ ] Login funciona com admin@silvaassociados.com.br
- [ ] Custos estão controlados (verificar com `costs`)
- [ ] Estratégia de economia configurada
- [ ] Testes passando

### **Configurações Opcionais:**
- [ ] Auto-shutdown configurado
- [ ] Cloud Run migrado
- [ ] Domínio DNS configurado
- [ ] Monitoramento ativo

### **Documentação:**
- [ ] Leu GUIA_OPERACIONAL_GCP.md
- [ ] Tem CHEAT_SHEET_GCP.md como referência
- [ ] Conhece comandos de emergência

---

## 🎯 RESULTADO ESPERADO

**Sistema funcionando:**
- ✅ Frontend acessível
- ✅ APIs funcionando
- ✅ Database com dados
- ✅ Custos controlados

**Economia configurada:**
- ✅ Scripts funcionando
- ✅ Auto-shutdown (opcional)
- ✅ Cloud Run (opcional)
- ✅ Monitoramento ativo

**Próximos passos:**
1. Usar CHEAT_SHEET_GCP.md para operação diária
2. Configurar desenvolvimento local se necessário
3. Planejar migração para produção