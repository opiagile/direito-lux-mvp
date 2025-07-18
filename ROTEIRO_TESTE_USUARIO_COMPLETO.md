# 🧪 ROTEIRO DE TESTE - JORNADA COMPLETA DO USUÁRIO

## 📋 VISÃO GERAL

Este roteiro simula a experiência completa de um novo escritório de advocacia usando o Direito Lux pela primeira vez, com monitoramento em tempo real dos logs do GCP.

**Objetivo:** Validar todo o fluxo desde o cadastro até o uso diário, capturando erros em tempo de execução.

---

## 🚀 PREPARAÇÃO DO AMBIENTE

### 1. **Iniciar Sistema**
```bash
# Verificar estado atual
./scripts/gcp-cost-optimizer.sh costs

# Iniciar sistema se necessário
./scripts/gcp-cost-optimizer.sh start

# Aguardar 2-3 minutos e verificar
curl -k https://35.188.198.87/api/health
```

### 2. **Abrir Monitoramento de Logs (4 terminais)**

#### **Terminal 1 - Logs do Frontend**
```bash
kubectl logs -n direito-lux-staging -l app=frontend -f --tail=10
```

#### **Terminal 2 - Logs do Auth Service**
```bash
kubectl logs -n direito-lux-staging -l app=auth-service -f --tail=10
```

#### **Terminal 3 - Logs do Tenant Service**
```bash
kubectl logs -n direito-lux-staging -l app=tenant-service -f --tail=10
```

#### **Terminal 4 - Logs de Erros (Todos os Serviços)**
```bash
kubectl logs -n direito-lux-staging --all-containers=true -f | grep -E "(ERROR|ERRO|error|Error|failed|Failed|FAILED|panic|PANIC)"
```

### 3. **Dashboard de Status**
```bash
# Em outro terminal, executar a cada 10 segundos
watch -n 10 'kubectl get pods -n direito-lux-staging | grep -v Running'
```

---

## 📝 CENÁRIO DE TESTE: ESCRITÓRIO "COSTA ADVOGADOS"

### **Dados do Novo Escritório:**
- **Nome:** Costa Advogados
- **CNPJ:** 12.345.678/0001-90
- **Email:** contato@costaadvogados.com.br
- **Telefone:** (11) 98765-4321
- **Plano:** Professional (R$299/mês)

### **Usuário Administrador:**
- **Nome:** Dr. João Costa
- **Email:** joao@costaadvogados.com.br
- **Senha:** S3nh@F0rt3!2025

---

## 🎯 JORNADA DE TESTE PASSO A PASSO

### **FASE 1: CADASTRO DO ESCRITÓRIO (Sign-up)**

#### 1.1 **Acessar Página de Cadastro**
```bash
# Verificar se frontend está respondendo
curl -k https://35.188.198.87/ -I

# Abrir no browser
open https://35.188.198.87/signup
```

**✅ Verificar nos logs:**
- [ ] Frontend: Request para `/signup`
- [ ] Nenhum erro 404 ou 500

#### 1.2 **Preencher Formulário de Cadastro**
- [ ] Nome do Escritório: Costa Advogados
- [ ] CNPJ: 12.345.678/0001-90
- [ ] Email: contato@costaadvogados.com.br
- [ ] Telefone: (11) 98765-4321
- [ ] Plano: Professional

**✅ Verificar nos logs:**
- [ ] Tenant Service: `POST /api/v1/tenants`
- [ ] Auth Service: `POST /api/v1/auth/register`
- [ ] Nenhum erro de validação

#### 1.3 **Criar Usuário Admin**
- [ ] Nome: Dr. João Costa
- [ ] Email: joao@costaadvogados.com.br  
- [ ] Senha: S3nh@F0rt3!2025
- [ ] Confirmar Senha: S3nh@F0rt3!2025

**✅ Verificar nos logs:**
- [ ] Auth Service: Hash de senha sendo criado
- [ ] Database: Insert em `users` table
- [ ] Email Service: Tentativa de envio (pode falhar se não configurado)

### **COMANDO DE VERIFICAÇÃO:**
```bash
# Verificar se tenant foi criado
kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -c "SELECT id, name, email, plan_type FROM tenants WHERE email='contato@costaadvogados.com.br';"

# Verificar se usuário foi criado
kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -c "SELECT id, email, first_name, last_name, role FROM users WHERE email='joao@costaadvogados.com.br';"
```

---

### **FASE 2: PRIMEIRO LOGIN**

#### 2.1 **Fazer Login**
```bash
# Testar via API primeiro
curl -k -X POST https://35.188.198.87/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: $(kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -t -c "SELECT id FROM tenants WHERE email='contato@costaadvogados.com.br';" | tr -d ' ')" \
  -d '{"email":"joao@costaadvogados.com.br","password":"S3nh@F0rt3!2025"}'
```

**✅ Verificar nos logs:**
- [ ] Auth Service: `POST /api/v1/auth/login`
- [ ] JWT Token gerado
- [ ] Session criada
- [ ] Nenhum erro de autenticação

#### 2.2 **Acessar Dashboard**
- [ ] Login via browser
- [ ] Dashboard carrega corretamente
- [ ] Nome do usuário aparece
- [ ] Menu lateral visível

**✅ Verificar nos logs:**
- [ ] Frontend: Requests para API
- [ ] Tenant Service: Busca dados do escritório
- [ ] Nenhum erro de CORS

---

### **FASE 3: CADASTRO DE PROCESSO**

#### 3.1 **Acessar Área de Processos**
- [ ] Clicar em "Processos" no menu
- [ ] Clicar em "Novo Processo"

**✅ Verificar nos logs:**
- [ ] Process Service: `GET /api/v1/processes`
- [ ] Paginação funcionando

#### 3.2 **Cadastrar Novo Processo**
**Dados do Processo:**
- **Número:** 1234567-89.2025.8.26.0100
- **Tribunal:** TJSP
- **Comarca:** São Paulo
- **Vara:** 1ª Vara Cível
- **Classe:** Ação de Cobrança
- **Assunto:** Cobrança - Prestação de Serviços
- **Valor da Causa:** R$ 50.000,00

**Partes:**
- **Autor:** Costa Advogados (representando Cliente XYZ)
- **Réu:** Empresa ABC Ltda

**✅ Verificar nos logs:**
- [ ] Process Service: `POST /api/v1/processes`
- [ ] Validação do número do processo
- [ ] DataJud Service: Tentativa de busca (pode falhar se não configurado)
- [ ] Notification Service: Preparação de notificações

### **COMANDO DE VERIFICAÇÃO:**
```bash
# Verificar se processo foi criado
kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -c "SELECT id, case_number, court, status FROM processes WHERE case_number='1234567-89.2025.8.26.0100';"
```

---

### **FASE 4: TESTE DE NOTIFICAÇÕES**

#### 4.1 **Configurar Notificações**
- [ ] Acessar "Configurações" → "Notificações"
- [ ] Ativar notificações por Email
- [ ] Ativar notificações por WhatsApp (se disponível)

**✅ Verificar nos logs:**
- [ ] Notification Service: Configurações salvas
- [ ] Fila RabbitMQ: Mensagens enfileiradas

#### 4.2 **Simular Movimentação**
```bash
# Criar movimentação via API
PROCESS_ID=$(kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -t -c "SELECT id FROM processes WHERE case_number='1234567-89.2025.8.26.0100';" | tr -d ' ')

curl -k -X POST https://35.188.198.87/api/v1/processes/$PROCESS_ID/movements \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "description": "Juntada de Petição",
    "date": "2025-07-15T10:00:00Z",
    "type": "PETICAO"
  }'
```

**✅ Verificar nos logs:**
- [ ] Process Service: Movimentação criada
- [ ] Notification Service: Notificação processada
- [ ] Email/WhatsApp: Tentativa de envio

---

### **FASE 5: TESTE DE BUSCA E FILTROS**

#### 5.1 **Buscar Processo**
- [ ] Usar barra de busca
- [ ] Buscar por número: "1234567"
- [ ] Resultado aparece corretamente

**✅ Verificar nos logs:**
- [ ] Search Service: Query executada
- [ ] ElasticSearch: Busca realizada (se configurado)
- [ ] Tempo de resposta < 1 segundo

#### 5.2 **Filtrar Processos**
- [ ] Filtrar por status: "Ativo"
- [ ] Filtrar por tribunal: "TJSP"
- [ ] Filtrar por data

**✅ Verificar nos logs:**
- [ ] Process Service: Filtros aplicados
- [ ] Query SQL otimizada
- [ ] Nenhum erro de timeout

---

### **FASE 6: TESTE DE RELATÓRIOS**

#### 6.1 **Gerar Relatório**
- [ ] Acessar "Relatórios"
- [ ] Selecionar "Relatório de Processos"
- [ ] Período: Últimos 30 dias
- [ ] Gerar PDF

**✅ Verificar nos logs:**
- [ ] Report Service: Geração iniciada
- [ ] PDF criado com sucesso
- [ ] Download funcionando

---

## 🔍 MONITORAMENTO CONTÍNUO

### **Comandos de Verificação em Tempo Real**

#### **1. Status Geral do Sistema**
```bash
# Criar script de monitoramento
cat > monitor.sh << 'EOF'
#!/bin/bash
while true; do
  clear
  echo "=== DIREITO LUX - MONITOR DE SISTEMA ==="
  echo "Horário: $(date)"
  echo ""
  echo "PODS STATUS:"
  kubectl get pods -n direito-lux-staging | grep -E "(NAME|Running|Error|Crash)"
  echo ""
  echo "ÚLTIMOS ERROS (5 min):"
  kubectl logs -n direito-lux-staging --all-containers=true --since=5m 2>/dev/null | grep -E "(ERROR|ERRO|failed|Failed)" | tail -5
  echo ""
  echo "REQUESTS/SEGUNDO:"
  kubectl logs -n direito-lux-staging -l app=frontend --since=1m 2>/dev/null | grep "GET\|POST" | wc -l
  echo ""
  echo "USO DE MEMÓRIA:"
  kubectl top pods -n direito-lux-staging 2>/dev/null | head -5
  sleep 10
done
EOF

chmod +x monitor.sh
./monitor.sh
```

#### **2. Captura de Erros Específicos**
```bash
# Erros de autenticação
kubectl logs -n direito-lux-staging -l app=auth-service -f | grep -E "(401|403|authentication|unauthorized)"

# Erros de database
kubectl logs -n direito-lux-staging --all-containers=true -f | grep -E "(database|postgres|connection refused|timeout)"

# Erros de API
kubectl logs -n direito-lux-staging -l app=frontend -f | grep -E "(fetch failed|network error|CORS)"
```

#### **3. Verificação de Processos**
```bash
# Verificar se processo específico está rodando
check_process() {
  local process=$1
  echo "Verificando $process..."
  
  case $process in
    "auth")
      curl -sk https://35.188.198.87/api/v1/auth/health | jq '.'
      ;;
    "tenant")
      curl -sk https://35.188.198.87/api/v1/tenants/health | jq '.'
      ;;
    "process")
      curl -sk https://35.188.198.87/api/v1/processes/health | jq '.'
      ;;
    "frontend")
      curl -sk https://35.188.198.87/api/health | jq '.'
      ;;
    *)
      echo "Processo desconhecido: $process"
      ;;
  esac
}

# Usar: check_process auth
```

---

## 📊 MÉTRICAS DE SUCESSO

### **Performance**
- [ ] Tempo de login < 2 segundos
- [ ] Tempo de carregamento dashboard < 3 segundos
- [ ] Busca de processos < 1 segundo
- [ ] Geração de relatório < 5 segundos

### **Estabilidade**
- [ ] Nenhum pod reiniciando
- [ ] Nenhum erro 500 durante o teste
- [ ] Nenhum timeout de database
- [ ] Memória estável (não crescendo)

### **Funcionalidade**
- [ ] Cadastro completo funciona
- [ ] Login/logout funciona
- [ ] CRUD de processos funciona
- [ ] Notificações são enfileiradas
- [ ] Busca retorna resultados
- [ ] Relatórios são gerados

---

## 🚨 TROUBLESHOOTING DURANTE O TESTE

### **Se Login Falhar:**
```bash
# Verificar hash da senha
kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -c "SELECT email, password_hash FROM users WHERE email='joao@costaadvogados.com.br';"

# Verificar tenant_id
kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -c "SELECT id, email FROM tenants WHERE email='contato@costaadvogados.com.br';"
```

### **Se Frontend Não Carregar:**
```bash
# Verificar ingress
kubectl get ingress -n direito-lux-staging
kubectl describe ingress -n direito-lux-staging

# Testar direto no pod
kubectl port-forward -n direito-lux-staging deploy/frontend 3000:3000
# Acessar http://localhost:3000
```

### **Se API Retornar 404:**
```bash
# Verificar rotas do ingress
kubectl get ingress -n direito-lux-staging direito-lux-ingress-apis -o yaml

# Testar serviço diretamente
kubectl port-forward -n direito-lux-staging svc/auth-service 8081:8080
curl http://localhost:8081/api/v1/auth/health
```

---

## 📋 CHECKLIST FINAL

### **Preparação**
- [ ] Sistema iniciado
- [ ] 4 terminais com logs abertos
- [ ] Monitor.sh rodando
- [ ] Browser com DevTools aberto

### **Execução**
- [ ] Fase 1: Cadastro completo
- [ ] Fase 2: Login bem-sucedido
- [ ] Fase 3: Processo cadastrado
- [ ] Fase 4: Notificações testadas
- [ ] Fase 5: Busca funcionando
- [ ] Fase 6: Relatório gerado

### **Validação**
- [ ] Nenhum erro crítico nos logs
- [ ] Performance dentro do esperado
- [ ] Dados persistidos corretamente
- [ ] Sistema estável

### **Limpeza**
- [ ] Screenshots dos erros salvos
- [ ] Logs importantes coletados
- [ ] Dados de teste documentados
- [ ] Sistema parado (se necessário)

---

## 💾 COLETA DE EVIDÊNCIAS

### **Salvar Logs do Teste**
```bash
# Criar diretório
mkdir -p teste-$(date +%Y%m%d-%H%M%S)
cd teste-$(date +%Y%m%d-%H%M%S)

# Coletar logs
kubectl logs -n direito-lux-staging --all-containers=true --since=1h > all-logs.txt
kubectl get events -n direito-lux-staging > events.txt
kubectl top pods -n direito-lux-staging > resources.txt

# Queries importantes
kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -c "SELECT COUNT(*) as total_tenants FROM tenants;" > stats.txt
kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -c "SELECT COUNT(*) as total_users FROM users;" >> stats.txt
kubectl exec -n direito-lux-staging deploy/postgres -- psql -U direito_lux -d direito_lux_staging -c "SELECT COUNT(*) as total_processes FROM processes;" >> stats.txt

# Comprimir
cd ..
tar -czf teste-$(date +%Y%m%d-%H%M%S).tar.gz teste-$(date +%Y%m%d-%H%M%S)/
```

---

## 🎯 CONCLUSÃO

Este roteiro simula a jornada completa de um novo usuário, desde o cadastro até o uso diário do sistema, com monitoramento ativo de logs e captura de erros em tempo real.

**Tempo estimado:** 45-60 minutos

**Resultado esperado:** Sistema validado end-to-end com evidências coletadas