# 🤖 TESTE DO ZERO COM ACOMPANHAMENTO DA IA

## 📋 RESUMO EXECUTIVO

**Objetivo:** Testar jornada completa de novo usuário com IA monitorando logs do GCP em tempo real

**Cenário:** Escritório "Costa Advogados" se cadastrando e usando o sistema pela primeira vez

**Duração:** 45-60 minutos

**Monitoramento:** IA captura erros em tempo real e verifica processos

---

## 🎯 PREPARAÇÃO (10 minutos)

### **1. Iniciar Sistema e Monitoramento**
```bash
# Iniciar sistema
./scripts/monitor-teste-usuario.sh start

# Verificar se está pronto
./scripts/monitor-teste-usuario.sh status
```

### **2. Abrir 3 Terminais para IA**

#### **Terminal 1: Dashboard Geral**
```bash
./scripts/monitor-teste-usuario.sh dashboard
```

#### **Terminal 2: Captura de Erros**
```bash
./scripts/monitor-teste-usuario.sh errors
```

#### **Terminal 3: Comandos Sob Demanda**
```bash
# Usar conforme necessário durante teste
./scripts/monitor-teste-usuario.sh db
./scripts/monitor-teste-usuario.sh performance
```

---

## 🧪 EXECUÇÃO DO TESTE (35 minutos)

### **FASE 1: Cadastro (10 min)**

#### **👤 Ação do Usuário:**
1. Acessar https://35.188.198.87/signup
2. Preencher dados do escritório:
   - Nome: Costa Advogados
   - CNPJ: 12.345.678/0001-90
   - Email: contato@costaadvogados.com.br
   - Telefone: (11) 98765-4321
   - Plano: Professional
3. Criar usuário admin:
   - Nome: Dr. João Costa
   - Email: joao@costaadvogados.com.br
   - Senha: S3nh@F0rt3!2025

#### **🤖 IA Executa:**
```bash
# Iniciar monitoramento específico
./scripts/monitor-teste-usuario.sh signup
```

#### **🔍 IA Verifica:**
```bash
# Verificar se dados foram salvos
./scripts/monitor-teste-usuario.sh db

# Verificar tenant criado
kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -c "SELECT id, name, email FROM tenants WHERE email='contato@costaadvogados.com.br';"

# Verificar usuário criado
kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -c "SELECT id, email, first_name FROM users WHERE email='joao@costaadvogados.com.br';"
```

**✅ Checkpoint:** IA confirma se cadastro foi bem-sucedido

---

### **FASE 2: Primeiro Login (5 min)**

#### **👤 Ação do Usuário:**
1. Acessar https://35.188.198.87/login
2. Fazer login com joao@costaadvogados.com.br / S3nh@F0rt3!2025
3. Verificar se dashboard carrega

#### **🤖 IA Executa:**
```bash
# Monitorar tentativas de login
./scripts/monitor-teste-usuario.sh login joao@costaadvogados.com.br
```

#### **🔍 IA Verifica:**
```bash
# Testar login via API
curl -k -X POST https://35.188.198.87/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"joao@costaadvogados.com.br","password":"S3nh@F0rt3!2025"}'

# Verificar session criada
kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -c "SELECT user_id, created_at FROM sessions WHERE user_id=(SELECT id FROM users WHERE email='joao@costaadvogados.com.br');"
```

**✅ Checkpoint:** IA confirma se login foi bem-sucedido

---

### **FASE 3: Cadastro de Processo (10 min)**

#### **👤 Ação do Usuário:**
1. Clicar em "Processos" → "Novo Processo"
2. Preencher dados:
   - Número: 1234567-89.2025.8.26.0100
   - Tribunal: TJSP
   - Comarca: São Paulo
   - Classe: Ação de Cobrança
   - Valor: R$ 50.000,00
3. Salvar processo

#### **🤖 IA Executa:**
```bash
# Monitorar logs do process-service
kubectl logs -n direito-lux-staging -l app=process-service -f | grep -E "(POST|processes|created|error)"
```

#### **🔍 IA Verifica:**
```bash
# Verificar se processo foi criado
kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -c "SELECT id, case_number, court, status FROM processes WHERE case_number='1234567-89.2025.8.26.0100';"

# Verificar estatísticas
./scripts/monitor-teste-usuario.sh db
```

**✅ Checkpoint:** IA confirma se processo foi cadastrado

---

### **FASE 4: Teste de Busca (5 min)**

#### **👤 Ação do Usuário:**
1. Usar barra de busca
2. Buscar por "1234567"
3. Verificar se resultado aparece

#### **🤖 IA Executa:**
```bash
# Monitorar search-service
kubectl logs -n direito-lux-staging -l app=search-service -f | grep -E "(search|query|results)"
```

#### **🔍 IA Verifica:**
```bash
# Testar busca via API
curl -k "https://35.188.198.87/api/v1/processes/search?q=1234567"

# Verificar performance
./scripts/monitor-teste-usuario.sh performance
```

**✅ Checkpoint:** IA confirma se busca funcionou

---

### **FASE 5: Teste de Notificações (5 min)**

#### **👤 Ação do Usuário:**
1. Acessar Configurações → Notificações
2. Ativar notificações por email
3. Simular movimentação no processo

#### **🤖 IA Executa:**
```bash
# Monitorar notification-service
kubectl logs -n direito-lux-staging -l app=notification-service -f | grep -E "(notification|email|queue|sent)"
```

#### **🔍 IA Verifica:**
```bash
# Verificar fila RabbitMQ
kubectl exec -n direito-lux-staging deploy/rabbitmq -- rabbitmqctl list_queues

# Verificar se notificação foi criada
kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -c "SELECT COUNT(*) FROM notifications WHERE user_id=(SELECT id FROM users WHERE email='joao@costaadvogados.com.br');"
```

**✅ Checkpoint:** IA confirma se notificações foram processadas

---

## 📊 RELATÓRIO FINAL (10 minutos)

### **🤖 IA Executa Diagnóstico Completo:**
```bash
# 1. Status final do sistema
./scripts/monitor-teste-usuario.sh status

# 2. Verificar dados finais
./scripts/monitor-teste-usuario.sh db

# 3. Verificar performance
./scripts/monitor-teste-usuario.sh performance

# 4. Coletar estatísticas finais
kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -c "
SELECT 
    'Tenants' as tipo, COUNT(*) as total FROM tenants 
UNION ALL SELECT 
    'Users', COUNT(*) FROM users 
UNION ALL SELECT 
    'Processes', COUNT(*) FROM processes 
UNION ALL SELECT 
    'Sessions', COUNT(*) FROM sessions 
UNION ALL SELECT 
    'Notifications', COUNT(*) FROM notifications;
"

# 5. Resumo de erros
grep -c "ERROR\|error\|failed" teste-usuario-*.log
```

---

## 📋 TEMPLATE DE RELATÓRIO PARA IA

### **RELATÓRIO DE TESTE - ESCRITÓRIO COSTA ADVOGADOS**

**Data/Hora:** [DATA_ATUAL]  
**Duração:** [TEMPO_TOTAL]  
**Sistema:** https://35.188.198.87

#### **🎯 RESUMO EXECUTIVO**
- ✅ Cadastro de escritório: [SUCESSO/FALHA]
- ✅ Primeiro login: [SUCESSO/FALHA]  
- ✅ Cadastro de processo: [SUCESSO/FALHA]
- ✅ Busca de processos: [SUCESSO/FALHA]
- ✅ Notificações: [SUCESSO/FALHA]

#### **📊 MÉTRICAS COLETADAS**
- Tempo médio de resposta: [X]ms
- Erros capturados: [X]
- Pods reiniciados: [X]
- Requests processados: [X]

#### **💾 DADOS CRIADOS**
- Tenants: [X] (antes) → [X] (depois)
- Users: [X] (antes) → [X] (depois)  
- Processes: [X] (antes) → [X] (depois)
- Sessions: [X] (antes) → [X] (depois)

#### **🚨 PROBLEMAS ENCONTRADOS**
1. **[PROBLEMA_1]**
   - Serviço: [NOME_SERVIÇO]
   - Erro: [MENSAGEM_ERRO]
   - Log: [LINHA_DO_LOG]
   - Impacto: [ALTO/MÉDIO/BAIXO]

#### **💡 RECOMENDAÇÕES**
1. [RECOMENDAÇÃO_1]
2. [RECOMENDAÇÃO_2]

#### **🔧 COMANDOS PARA REPRODUZIR PROBLEMAS**
```bash
# Problema 1
[COMANDO_REPRODUZIR_1]

# Problema 2  
[COMANDO_REPRODUZIR_2]
```

#### **📈 ANÁLISE DE PERFORMANCE**
- Sistema estável: [SIM/NÃO]
- Uso de memória: [NORMAL/ALTO]
- Uso de CPU: [NORMAL/ALTO]
- Latência aceitável: [SIM/NÃO]

---

## 🚀 COMANDOS RÁPIDOS PARA IA

### **Durante o Teste:**
```bash
# Verificar se processo mencionado está funcionando
check_process() {
    local process=$1
    case $process in
        "cadastro") kubectl logs -n direito-lux-staging -l app=tenant-service --tail=5 ;;
        "login") kubectl logs -n direito-lux-staging -l app=auth-service --tail=5 ;;
        "processo") kubectl logs -n direito-lux-staging -l app=process-service --tail=5 ;;
        "busca") kubectl logs -n direito-lux-staging -l app=search-service --tail=5 ;;
        "notificacao") kubectl logs -n direito-lux-staging -l app=notification-service --tail=5 ;;
        *) echo "Processo não reconhecido: $process" ;;
    esac
}

# Usar: check_process "login"
```

### **Capturar Evidências:**
```bash
# Salvar estado atual
mkdir -p evidencias/$(date +%Y%m%d-%H%M%S)
cd evidencias/$(date +%Y%m%d-%H%M%S)

# Coletar logs
kubectl logs -n direito-lux-staging --all-containers=true --since=1h > logs-completos.txt
kubectl get events -n direito-lux-staging > eventos.txt
kubectl top pods -n direito-lux-staging > recursos.txt

# Estatísticas do banco
kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -c "SELECT 'Tenants', COUNT(*) FROM tenants UNION ALL SELECT 'Users', COUNT(*) FROM users UNION ALL SELECT 'Processes', COUNT(*) FROM processes;" > estatisticas.txt
```

---

## 🎯 RESULTADO ESPERADO

**✅ Teste Bem-Sucedido:**
- Sistema responde em < 2 segundos
- Nenhum erro crítico nos logs
- Todos os dados salvos corretamente
- Funcionalidades básicas operacionais

**📊 Métricas Alvo:**
- Uptime: 100% durante teste
- Erros: 0 críticos
- Performance: < 1s para busca
- Dados: 100% persistidos

**🎉 Sistema validado para uso real por escritórios de advocacia!**