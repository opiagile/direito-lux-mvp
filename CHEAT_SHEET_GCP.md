# ⚡ CHEAT SHEET - CONTROLE RÁPIDO DO AMBIENTE GCP

## 🎯 COMANDOS MAIS USADOS

### **🚀 LIGAR SISTEMA**
```bash
./scripts/gcp-cost-optimizer.sh start
```
**Tempo:** 2-3 minutos | **Custo:** R$0,60/hora

### **🛑 DESLIGAR SISTEMA**
```bash
./scripts/gcp-cost-optimizer.sh stop
```
**Tempo:** 1 minuto | **Custo:** R$0/hora

### **📊 VER CUSTOS**
```bash
./scripts/gcp-cost-optimizer.sh costs
```

### **🚨 PARAR TUDO (EMERGÊNCIA)**
```bash
./scripts/migrate-to-cloud-run.sh emergency
```

---

## 🔄 FLUXO DIÁRIO RECOMENDADO

### **🌅 MANHÃ (Iniciar trabalho)**
```bash
# 1. Verificar status
./scripts/gcp-cost-optimizer.sh costs

# 2. Iniciar sistema
./scripts/gcp-cost-optimizer.sh start

# 3. Aguardar 2-3 minutos

# 4. Testar se funciona
curl -k https://35.188.198.87/api/health
```

### **🌙 NOITE (Finalizar trabalho)**
```bash
# 1. Parar sistema
./scripts/gcp-cost-optimizer.sh stop

# 2. Confirmar que parou
./scripts/gcp-cost-optimizer.sh costs
```

---

## 💰 TABELA DE CUSTOS

| Configuração | Custo/Hora | Custo/Dia | Custo/Mês |
|-------------|------------|-----------|-----------|
| **6 nodes (original)** | R$3,60 | R$87,00 | R$2.610 |
| **1 node (otimizado)** | R$0,60 | R$14,40 | R$450 |
| **0 nodes (parado)** | R$0,00 | R$0,00 | R$0 |
| **Cloud Run** | R$0,10 | R$2,40 | R$30 |

---

## 🎛️ CONFIGURAÇÕES DISPONÍVEIS

### **Opção 1: Manual (Recomendado para desenvolvimento)**
- Liga/desliga quando quiser
- Custo controlado por você
- Economia: 83% se usar 8h/dia

### **Opção 2: Auto-shutdown (Recomendado para equipe)**
```bash
./scripts/setup-auto-shutdown.sh setup
```
- Para às 23h automaticamente
- Liga via página web
- Economia: 89%

### **Opção 3: Cloud Run (Recomendado para staging)**
```bash
./scripts/migrate-to-cloud-run.sh setup-cloudrun
```
- Escala automaticamente
- Só paga quando usa
- Economia: 98%

---

## 📱 MONITORAMENTO RÁPIDO

### **Verificar se está rodando:**
```bash
# Ver quantos nodes
gcloud container clusters describe direito-lux-gke-staging \
  --region=us-central1 --project=direito-lux-staging-2025 \
  --format="value(currentNodeCount)"
```

### **Testar sistema:**
```bash
# Frontend
curl -k https://35.188.198.87/api/health

# Auth Service
curl -k https://35.188.198.87/api/v1/auth/health
```

---

## 🆘 PROBLEMAS COMUNS

### **Sistema não responde:**
```bash
# Reiniciar
./scripts/gcp-cost-optimizer.sh stop
sleep 30
./scripts/gcp-cost-optimizer.sh start
```

### **Custo alto:**
```bash
# Verificar recursos
./scripts/gcp-cost-optimizer.sh costs

# Parar imediatamente
./scripts/migrate-to-cloud-run.sh emergency
```

### **Esqueci de desligar:**
```bash
# Parar agora
./scripts/gcp-cost-optimizer.sh stop

# Configurar auto-shutdown
./scripts/setup-auto-shutdown.sh setup
```

---

## 📞 SUPORTE TÉCNICO

### **Logs importantes:**
```bash
# Pods status
kubectl get pods -n direito-lux-staging

# Frontend logs
kubectl logs -n direito-lux-staging -l app=frontend

# Auth service logs
kubectl logs -n direito-lux-staging -l app=auth-service
```

### **Recursos GCP:**
```bash
# Instâncias rodando
gcloud compute instances list

# Clusters
gcloud container clusters list
```

---

## 💡 DICAS DE ECONOMIA

### **Máxima economia (98%):**
- Migre para Cloud Run
- Só paga quando usa
- Zero manutenção

### **Economia intermediária (89%):**
- Configure auto-shutdown
- Para automaticamente à noite
- Liga sob demanda

### **Economia básica (83%):**
- Controle manual
- Ligue só quando trabalhar
- Desligue sempre após uso

---

## 🎯 METAS DE CUSTO

### **Custo zero:** Sistema parado 24/7
### **Custo baixo:** R$15/mês (Cloud Run)
### **Custo médio:** R$450/mês (GKE otimizado)
### **Custo alto:** R$2.610/mês (configuração original)

**🎯 OBJETIVO:** Manter sempre abaixo de R$100/mês