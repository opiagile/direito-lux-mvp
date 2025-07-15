# 🎯 GUIA OPERACIONAL - GERENCIAMENTO DIÁRIO DO AMBIENTE GCP

## 📋 RESUMO EXECUTIVO

**Problema:** Custo de R$115 em 2 dias (R$1.725/mês) com 6 nodes rodando 24/7  
**Solução:** Scripts para controlar ambiente sob demanda  
**Economia:** Até 98% (R$20.340/ano)

---

## 🚀 CENÁRIOS DE USO DIÁRIO

### 1. 🌅 **INICIAR AMBIENTE PARA TRABALHAR**

#### Situação: Preciso acessar o sistema staging

```bash
# Opção A: Cluster mínimo (R$15/dia)
./scripts/gcp-cost-optimizer.sh start

# Opção B: Cluster otimizado (R$15/dia com auto-shutdown)
./scripts/gcp-cost-optimizer.sh optimize
```

**Tempo:** 2-3 minutos para ficar pronto  
**Custo:** R$0,60/hora (R$14,40/dia se deixar ligado)

#### Verificar se está funcionando:
```bash
# Verificar status
./scripts/gcp-cost-optimizer.sh costs

# Testar sistema
curl -k https://35.188.198.87/api/health
```

### 2. 🌙 **PARAR AMBIENTE APÓS TRABALHAR**

#### Situação: Terminei o trabalho, quero economizar

```bash
# Parar cluster imediatamente
./scripts/gcp-cost-optimizer.sh stop
```

**Resultado:** Custo vai para R$0/hora  
**Economia:** R$14,40/dia

### 3. ⚡ **MODO EMERGÊNCIA - PARAR TUDO AGORA**

#### Situação: Custos muito altos, preciso parar AGORA

```bash
# Parar tudo imediatamente
./scripts/migrate-to-cloud-run.sh emergency
```

**Resultado:** Todo o cluster para, custo R$0/hora

---

## 🤖 AUTOMATIZAÇÃO COMPLETA

### 1. **CONFIGURAR AUTO-SHUTDOWN**

#### Situação: Quero que pare automaticamente à noite

```bash
# Configurar para parar às 23h e iniciar sob demanda
./scripts/setup-auto-shutdown.sh setup
```

**O que é criado:**
- ✅ Cloud Function para gerenciar cluster
- ✅ Cloud Scheduler para parar às 23:00
- ✅ Página web para iniciar sistema

**Fluxo:**
1. Sistema para às 23:00 automaticamente
2. Quando você acessa https://35.188.198.87
3. Página aparece com botão "Iniciar Sistema"
4. Clica e aguarda 2-3 minutos
5. Sistema fica disponível normalmente

### 2. **MIGRAR PARA CLOUD RUN (RECOMENDADO)**

#### Situação: Quero máxima economia (98%) e zero manutenção

```bash
# Migrar completamente para Cloud Run
./scripts/migrate-to-cloud-run.sh setup-cloudrun
```

**Benefícios:**
- ✅ Escala para zero automaticamente
- ✅ Só paga quando há tráfego
- ✅ Sem gerenciamento de nodes
- ✅ Custo: R$30/mês (98% economia)

---

## 📊 MONITORAMENTO E CONTROLE

### **Verificar Custos Atuais:**
```bash
# Ver recursos ativos e custos
./scripts/gcp-cost-optimizer.sh costs
```

### **Verificar Status do Cluster:**
```bash
# Quantos nodes estão rodando
gcloud container clusters describe direito-lux-gke-staging \
  --region=us-central1 \
  --project=direito-lux-staging-2025 \
  --format="value(currentNodeCount)"
```

### **Verificar Status dos Pods:**
```bash
# Se cluster estiver rodando
kubectl get pods -n direito-lux-staging
```

---

## 🎯 CENÁRIOS COMUNS

### **Cenário 1: Trabalho das 9h às 18h**
```bash
# 9h - Iniciar
./scripts/gcp-cost-optimizer.sh start

# 18h - Parar
./scripts/gcp-cost-optimizer.sh stop
```
**Custo:** R$5,40/dia (9h × R$0,60/h)

### **Cenário 2: Trabalho esporádico**
```bash
# Configurar auto-shutdown
./scripts/setup-auto-shutdown.sh setup
```
**Custo:** R$0 quando não usa, R$0,60/h quando usa

### **Cenário 3: Demo para cliente**
```bash
# Migrar para Cloud Run
./scripts/migrate-to-cloud-run.sh setup-cloudrun
```
**Custo:** R$1/mês + R$0,10 por demo

---

## 🚨 PROCEDIMENTOS DE EMERGÊNCIA

### **Custo Disparando:**
```bash
# IMEDIATO
./scripts/migrate-to-cloud-run.sh emergency

# VERIFICAR
./scripts/gcp-cost-optimizer.sh costs
```

### **Sistema Não Responde:**
```bash
# Reiniciar cluster
./scripts/gcp-cost-optimizer.sh stop
sleep 30
./scripts/gcp-cost-optimizer.sh start
```

### **Preciso Deletar Tudo:**
```bash
# Deletar cluster completamente
gcloud container clusters delete direito-lux-gke-staging \
  --region=us-central1 \
  --project=direito-lux-staging-2025
```

---

## 📋 CHECKLIST DIÁRIO

### **Antes de Sair do Trabalho:**
- [ ] Sistema funcionando como esperado?
- [ ] Preciso deixar rodando à noite?
- [ ] Se não, executar: `./scripts/gcp-cost-optimizer.sh stop`

### **Início do Trabalho:**
- [ ] Verificar custo atual: `./scripts/gcp-cost-optimizer.sh costs`
- [ ] Iniciar sistema: `./scripts/gcp-cost-optimizer.sh start`
- [ ] Aguardar 2-3 minutos
- [ ] Testar: `curl -k https://35.188.198.87/api/health`

### **Weekly Review:**
- [ ] Verificar custos no console GCP
- [ ] Considerar migração para Cloud Run se uso for baixo
- [ ] Avaliar se auto-shutdown está funcionando

---

## 🔧 TROUBLESHOOTING

### **Problema: Sistema não inicia**
```bash
# Verificar se há problemas
kubectl get pods -n direito-lux-staging
kubectl logs -n direito-lux-staging -l app=frontend
```

### **Problema: Custo alto inesperado**
```bash
# Ver recursos rodando
./scripts/gcp-cost-optimizer.sh costs
gcloud compute instances list
```

### **Problema: Auto-shutdown não funciona**
```bash
# Verificar Cloud Scheduler
gcloud scheduler jobs list
gcloud scheduler jobs run shutdown-cluster
```

---

## 💡 RECOMENDAÇÕES FINAIS

### **Para Desenvolvimento Individual:**
- Use GKE com auto-shutdown
- Custo: R$450/mês
- Comando: `./scripts/setup-auto-shutdown.sh setup`

### **Para Staging/Demo:**
- Use Cloud Run
- Custo: R$30/mês
- Comando: `./scripts/migrate-to-cloud-run.sh setup-cloudrun`

### **Para Produção:**
- Use GKE otimizado sem auto-shutdown
- Custo: R$450/mês
- Comando: `./scripts/gcp-cost-optimizer.sh optimize`

---

## 📞 COMANDOS RÁPIDOS

| Ação | Comando | Tempo | Custo |
|------|---------|-------|-------|
| **Iniciar** | `./scripts/gcp-cost-optimizer.sh start` | 2-3 min | R$0,60/h |
| **Parar** | `./scripts/gcp-cost-optimizer.sh stop` | 1 min | R$0/h |
| **Otimizar** | `./scripts/gcp-cost-optimizer.sh optimize` | 5 min | R$0,60/h |
| **Emergência** | `./scripts/migrate-to-cloud-run.sh emergency` | 2 min | R$0/h |
| **Auto-shutdown** | `./scripts/setup-auto-shutdown.sh setup` | 10 min | R$0,60/h |
| **Cloud Run** | `./scripts/migrate-to-cloud-run.sh setup-cloudrun` | 20 min | R$0,10/h |

**🎯 RESULTADO:** Controle total sobre custos GCP com economia de até 98%