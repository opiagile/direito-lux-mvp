# 🤖 COMANDOS DE MONITORAMENTO PARA IA

## 📋 VISÃO GERAL

Este documento contém comandos específicos para a IA monitorar o sistema em tempo real durante os testes, capturando erros e verificando processos.

---

## 🚀 COMANDOS DE INICIALIZAÇÃO

### **1. Verificar Estado Inicial**
```bash
# Status do cluster
gcloud container clusters describe direito-lux-gke-staging --region=us-central1 --format="value(currentNodeCount)"

# Pods rodando
kubectl get pods -n direito-lux-staging --no-headers | grep Running | wc -l

# Serviços críticos
kubectl get pods -n direito-lux-staging | grep -E "(auth|tenant|process|frontend|postgres)" | grep -v Running
```

### **2. Iniciar Sistema (se necessário)**
```bash
# Verificar custo atual
./scripts/gcp-cost-optimizer.sh costs

# Iniciar
./scripts/gcp-cost-optimizer.sh start

# Aguardar pods ficarem prontos
kubectl wait --for=condition=ready pod -l app=frontend -n direito-lux-staging --timeout=180s
kubectl wait --for=condition=ready pod -l app=auth-service -n direito-lux-staging --timeout=180s
```

---

## 📊 COMANDOS DE MONITORAMENTO CONTÍNUO

### **1. Dashboard de Status em Tempo Real**
```bash
# Executar em terminal dedicado
watch -n 5 'echo "=== STATUS GERAL ==="; \
kubectl get pods -n direito-lux-staging | grep -v Running; \
echo ""; \
echo "=== ERROS RECENTES (2min) ==="; \
kubectl logs -n direito-lux-staging --all-containers=true --since=2m 2>/dev/null | grep -i error | tail -3; \
echo ""; \
echo "=== REQUESTS/MIN ==="; \
kubectl logs -n direito-lux-staging -l app=frontend --since=1m 2>/dev/null | grep -c "GET\|POST"'
```

### **2. Captura de Logs por Serviço**

#### **Auth Service**
```bash
# Logs em tempo real com filtro
kubectl logs -n direito-lux-staging -l app=auth-service -f | grep -E "(POST|GET|ERROR|login|register|401|403)"
```

#### **Tenant Service**
```bash
# Logs com destaque de erros
kubectl logs -n direito-lux-staging -l app=tenant-service -f | grep --color=always -E "(ERROR|error|failed|POST /api/v1/tenants|GET /api/v1/tenants)"
```

#### **Frontend**
```bash
# Logs de requisições e erros
kubectl logs -n direito-lux-staging -l app=frontend -f | grep -E "(Error|error|failed|fetch|api/|status: [4-5][0-9][0-9])"
```

### **3. Agregador de Erros Global**
```bash
# Todos os erros de todos os serviços
kubectl logs -n direito-lux-staging --all-containers=true -f | \
grep --color=always -E "(ERROR|ERRO|error|Error|failed|Failed|FAILED|panic|PANIC|Exception|exception|401|403|404|500|502|503)" | \
while read line; do echo "[$(date '+%H:%M:%S')] $line"; done
```

---

## 🔍 COMANDOS DE VERIFICAÇÃO ESPECÍFICA

### **1. Verificar Processo por Nome**
```bash
# Função para verificar se processo está funcionando
check_service_health() {
    local service=$1
    echo "🔍 Verificando $service..."
    
    # Verificar pod
    echo "📦 Pod status:"
    kubectl get pods -n direito-lux-staging -l app=$service
    
    # Verificar logs recentes
    echo "📝 Últimos erros (5min):"
    kubectl logs -n direito-lux-staging -l app=$service --since=5m 2>/dev/null | grep -i error | tail -3
    
    # Verificar endpoint de saúde
    echo "🏥 Health check:"
    case $service in
        "auth-service")
            kubectl exec -n direito-lux-staging deploy/$service -- wget -qO- http://localhost:8080/health 2>/dev/null || echo "❌ Health check falhou"
            ;;
        "tenant-service")
            kubectl exec -n direito-lux-staging deploy/$service -- wget -qO- http://localhost:8080/health 2>/dev/null || echo "❌ Health check falhou"
            ;;
        "frontend")
            kubectl exec -n direito-lux-staging deploy/$service -- wget -qO- http://localhost:3000/api/health 2>/dev/null || echo "❌ Health check falhou"
            ;;
    esac
    echo "---"
}

# Usar: check_service_health auth-service
```

### **2. Verificar Fluxo de Cadastro**
```bash
# Monitorar cadastro de novo tenant em tempo real
monitor_signup() {
    echo "👀 Monitorando fluxo de cadastro..."
    
    # Terminal 1: Tenant service
    kubectl logs -n direito-lux-staging -l app=tenant-service -f | grep --color=always -E "(POST.*tenants|tenant.*created|error)" &
    PID1=$!
    
    # Terminal 2: Auth service  
    kubectl logs -n direito-lux-staging -l app=auth-service -f | grep --color=always -E "(POST.*register|user.*created|password|error)" &
    PID2=$!
    
    # Terminal 3: Database
    kubectl logs -n direito-lux-staging -l app=postgres -f | grep --color=always -E "(INSERT|ERROR)" &
    PID3=$!
    
    echo "Monitorando... (Ctrl+C para parar)"
    wait
}
```

### **3. Verificar Fluxo de Login**
```bash
# Monitorar tentativas de login
monitor_login() {
    local email=$1
    echo "🔐 Monitorando login para: $email"
    
    kubectl logs -n direito-lux-staging -l app=auth-service -f | \
    grep --color=always -E "(login.*$email|JWT|token|401|403|success|failed)"
}

# Usar: monitor_login "joao@costaadvogados.com.br"
```

---

## 💾 COMANDOS DE VALIDAÇÃO DE DADOS

### **1. Verificar Dados no PostgreSQL**
```bash
# Função para queries rápidas
db_query() {
    local query=$1
    kubectl exec -n direito-lux-staging deploy/postgres -- \
    psql -U direito_lux -d direito_lux_staging -c "$query"
}

# Verificar tenant criado
check_tenant() {
    local email=$1
    db_query "SELECT id, name, email, plan_type, status FROM tenants WHERE email='$email';"
}

# Verificar usuário criado
check_user() {
    local email=$1
    db_query "SELECT id, email, first_name, last_name, role, status FROM users WHERE email='$email';"
}

# Verificar processo criado
check_process() {
    local case_number=$1
    db_query "SELECT id, case_number, court, status FROM processes WHERE case_number='$case_number';"
}

# Estatísticas gerais
db_stats() {
    echo "📊 ESTATÍSTICAS DO BANCO:"
    db_query "SELECT 'Tenants' as tipo, COUNT(*) as total FROM tenants UNION ALL SELECT 'Users', COUNT(*) FROM users UNION ALL SELECT 'Processes', COUNT(*) FROM processes;"
}
```

### **2. Verificar Filas RabbitMQ**
```bash
# Ver filas e mensagens
check_queues() {
    kubectl exec -n direito-lux-staging deploy/rabbitmq -- rabbitmqctl list_queues name messages consumers
}

# Monitorar mensagens em tempo real
monitor_messages() {
    watch -n 2 'kubectl exec -n direito-lux-staging deploy/rabbitmq -- rabbitmqctl list_queues name messages | grep -v "^Listing"'
}
```

---

## 🚨 COMANDOS DE DIAGNÓSTICO DE PROBLEMAS

### **1. Diagnóstico Rápido**
```bash
# Função all-in-one para diagnóstico
quick_diagnosis() {
    echo "🏥 DIAGNÓSTICO RÁPIDO DO SISTEMA"
    echo "================================"
    
    # Pods com problema
    echo "❌ Pods com problemas:"
    kubectl get pods -n direito-lux-staging | grep -v Running | grep -v NAME
    
    # Últimos eventos de erro
    echo -e "\n⚠️  Eventos recentes:"
    kubectl get events -n direito-lux-staging --field-selector type=Warning --sort-by='.lastTimestamp' | tail -5
    
    # Erros nos logs (últimos 5 min)
    echo -e "\n🔴 Erros recentes (5 min):"
    kubectl logs -n direito-lux-staging --all-containers=true --since=5m 2>/dev/null | grep -i error | tail -10
    
    # Uso de recursos
    echo -e "\n📊 Top 5 pods por CPU:"
    kubectl top pods -n direito-lux-staging --sort-by=cpu | head -6
    
    echo -e "\n💾 Top 5 pods por Memória:"
    kubectl top pods -n direito-lux-staging --sort-by=memory | head -6
}
```

### **2. Debug de Requisição Específica**
```bash
# Rastrear requisição através dos serviços
trace_request() {
    local request_id=$1
    echo "🔍 Rastreando request: $request_id"
    
    # Buscar em todos os logs
    kubectl logs -n direito-lux-staging --all-containers=true --since=10m | grep "$request_id" | sort -k1,2
}

# Capturar próxima requisição
capture_next_request() {
    echo "📸 Capturando próxima requisição de login..."
    kubectl logs -n direito-lux-staging -l app=frontend -f | grep -m1 "login" | awk '{print $NF}'
}
```

---

## 📈 COMANDOS DE PERFORMANCE

### **1. Métricas em Tempo Real**
```bash
# Dashboard de performance
performance_monitor() {
    watch -n 5 'echo "🚀 PERFORMANCE MONITOR"; \
    echo "====================="; \
    echo ""; \
    echo "⏱️  Response Times (último minuto):"; \
    kubectl logs -n direito-lux-staging -l app=frontend --since=1m 2>/dev/null | \
    grep "ms" | awk "{print \$NF}" | sort -n | tail -10; \
    echo ""; \
    echo "📊 Requests/segundo:"; \
    kubectl logs -n direito-lux-staging -l app=frontend --since=10s 2>/dev/null | grep -c "GET\|POST"; \
    echo ""; \
    echo "💾 Memory Usage:"; \
    kubectl top pods -n direito-lux-staging | grep -E "(auth|tenant|frontend)" | awk "{print \$1, \$3}"'
}
```

### **2. Teste de Carga Simples**
```bash
# Gerar carga para teste
load_test_login() {
    local email=$1
    local password=$2
    local count=${3:-10}
    
    echo "🔥 Testando $count logins simultâneos..."
    
    for i in $(seq 1 $count); do
        (
            time curl -sk -X POST https://35.188.198.87/api/v1/auth/login \
            -H "Content-Type: application/json" \
            -d "{\"email\":\"$email\",\"password\":\"$password\"}" \
            -w "\n$i: %{http_code} - %{time_total}s\n"
        ) &
    done
    wait
    echo "✅ Teste concluído"
}
```

---

## 🎯 COMANDOS PARA IA EXECUTAR DURANTE O TESTE

### **Sequência Recomendada:**

```bash
# 1. Preparar ambiente
./scripts/gcp-cost-optimizer.sh start
sleep 180  # Aguardar 3 minutos

# 2. Verificar sistema pronto
quick_diagnosis

# 3. Iniciar monitoramento (em terminais separados)
# Terminal 1:
watch -n 5 'kubectl get pods -n direito-lux-staging | grep -v Running'

# Terminal 2:
kubectl logs -n direito-lux-staging --all-containers=true -f | grep -i error

# Terminal 3:
monitor_signup  # Durante fase de cadastro

# Terminal 4:
db_stats  # Verificar dados periodicamente

# 4. Durante o teste, executar conforme necessário:
check_service_health auth-service
check_tenant "contato@costaadvogados.com.br"
check_user "joao@costaadvogados.com.br"
check_process "1234567-89.2025.8.26.0100"

# 5. Se houver problemas:
quick_diagnosis
trace_request "REQUEST_ID_AQUI"

# 6. Ao final:
performance_monitor  # Deixar rodando 2-3 minutos
db_stats  # Verificar totais finais
```

---

## 📝 TEMPLATE DE RELATÓRIO PARA IA

```markdown
## RELATÓRIO DE TESTE - [DATA/HORA]

### SUMÁRIO EXECUTIVO
- Sistema iniciado com sucesso: [SIM/NÃO]
- Cadastro funcionou: [SIM/NÃO]
- Login funcionou: [SIM/NÃO]
- Erros críticos encontrados: [QUANTIDADE]

### MÉTRICAS
- Tempo de inicialização: X minutos
- Tempo médio de login: X segundos
- Erros capturados: X
- Pods reiniciados: X

### PROBLEMAS ENCONTRADOS
1. [Descrição do problema]
   - Serviço afetado: 
   - Log do erro:
   - Comando para reproduzir:

### DADOS VALIDADOS
- Tenants criados: X
- Usuários criados: X
- Processos criados: X

### RECOMENDAÇÕES
1. [Ação recomendada]
2. [Ação recomendada]
```

**🎯 Com estes comandos, a IA pode monitorar efetivamente o sistema durante todo o teste!**