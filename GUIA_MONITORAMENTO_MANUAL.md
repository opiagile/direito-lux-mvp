# 🚨 GUIA DE MONITORAMENTO MANUAL - TESTES FRONTEND

## 📊 **COMANDOS PARA MONITORAR LOGS**

### **1. Monitorar Auth Service (CRÍTICO)**
```bash
# Em um terminal separado:
docker-compose logs -f auth-service

# Procurar por:
# ✅ "login successful" - Login OK
# ❌ "authentication failed" - Erro de login
# ⚠️ "invalid token" - Problema de JWT
# 🔐 "unauthorized" - Acesso negado
```

### **2. Monitorar PostgreSQL (CRÍTICO)**
```bash
# Em outro terminal:
docker-compose logs -f postgres

# Procurar por:
# ✅ "connection received" - Conexões OK
# ❌ "FATAL" - Erro crítico
# ⚠️ "ERROR" - Erro de query
# 📊 "SELECT/INSERT/UPDATE" - Operações DB
```

### **3. Status Geral dos Serviços**
```bash
# Verificar periodicamente:
docker-compose ps

# Deve mostrar:
# auth-service: Up
# postgres: Up  
# rabbitmq: Up
# (outros podem estar offline)
```

---

## 🧪 **ROTEIRO DE TESTES FRONTEND**

### **FASE 1: Testes de Login (5 min)**

#### **Teste 1.1: Login Starter**
- URL: http://localhost:3000/login
- Email: `admin@silvaassociados.com.br`
- Senha: `password`
- **Esperado:** Redirecionamento para dashboard
- **Monitorar:** Logs de auth-service para "login successful"

#### **Teste 1.2: Login Professional**
- Email: `admin@costasantos.com.br` 
- Senha: `password`
- **Esperado:** Dashboard com features Professional
- **Diferencial:** MCP Bot deve estar disponível

#### **Teste 1.3: Login Business**
- Email: `admin@machadoadvogados.com.br`
- Senha: `password`
- **Esperado:** Dashboard Business com mais recursos

#### **Teste 1.4: Login Enterprise**
- Email: `admin@barrosent.com.br`
- Senha: `password`
- **Esperado:** Dashboard completo sem limitações

### **FASE 2: Testes de Interface por Role (10 min)**

#### **Como Admin (Acesso Total)**
- ✅ **Menu Gestão de Usuários** deve estar visível
- ✅ **Configurações de Billing** deve estar acessível
- ✅ **Todos os relatórios** disponíveis
- ✅ **Configurações do tenant** acessíveis

#### **Como Manager (testar com gerente@silvaassociados.com.br)**
- ✅ **Relatórios e Analytics** visíveis
- ❌ **Billing** NÃO deve aparecer
- ✅ **Dashboard executivo** disponível
- ❌ **Gestão de usuários** limitada

#### **Como Lawyer (testar com advogado@silvaassociados.com.br)**
- ✅ **Seção Processos** com CRUD completo
- ✅ **AI Assistant** disponível
- ❌ **Configurações** não disponíveis
- ❌ **Billing** não visível

#### **Como Assistant (testar com cliente@silvaassociados.com.br)**
- ✅ **Visualização de processos** (somente leitura)
- ❌ **Criação/edição** bloqueada
- ❌ **Relatórios** não disponíveis
- ❌ **Configurações** bloqueadas

### **FASE 3: Testes de Funcionalidades (15 min)**

#### **Teste 3.1: Navegação entre Páginas**
- Dashboard → Processos → AI → Configurações
- **Monitorar:** Requests de API nos logs
- **Esperado:** Navegação suave sem erros

#### **Teste 3.2: Tentativa de Acesso Negado**
- Como Assistant, tentar acessar /admin
- **Esperado:** Redirect ou página de erro
- **Monitorar:** Logs "unauthorized" ou "403"

#### **Teste 3.3: Logout e Re-login**
- Fazer logout
- Tentar acessar página protegida
- Fazer login novamente
- **Esperado:** Fluxo completo funcionando

---

## 🔍 **ERROS COMUNS E SOLUÇÕES**

### **❌ "Network Error" no Frontend**
```bash
# Verificar se Auth Service está respondendo:
curl http://localhost:8081/health

# Se não responder:
docker-compose restart auth-service
```

### **❌ "Database Connection Error"**
```bash
# Verificar PostgreSQL:
docker-compose exec postgres pg_isready -U direito_lux

# Se falhar:
docker-compose restart postgres
```

### **❌ "Unauthorized" após login**
```bash
# Verificar logs do Auth Service:
docker-compose logs auth-service | grep -i "token\|jwt"

# Token pode estar expirado ou inválido
```

### **❌ Página em branco ou erro JS**
```bash
# Verificar console do navegador (F12)
# Procurar por:
# - Erros de CORS
# - Falhas de API
# - Problemas de JavaScript
```

---

## 📊 **MÉTRICAS PARA ACOMPANHAR**

### **Durante os Testes:**
- **Response Time:** < 500ms para login
- **Success Rate:** 100% dos logins
- **Errors:** 0 erros críticos
- **Navigation:** Todas as páginas carregando

### **Logs de Sucesso Esperados:**
```
[auth-service] login successful for user: admin@silvaassociados.com.br
[auth-service] token generated successfully
[postgres] connection received: host=172.x.x.x
[postgres] statement: SELECT * FROM users WHERE email = $1
```

### **Logs de Erro para Alertar:**
```
[auth-service] ERRO: authentication failed
[auth-service] ERRO: invalid credentials
[postgres] FATAL: database connection failed
[postgres] ERROR: relation "users" does not exist
```

---

## 🎯 **CHECKLIST DE VALIDAÇÃO**

### **✅ Login e Autenticação**
- [ ] Login Starter funcionando
- [ ] Login Professional funcionando  
- [ ] Login Business funcionando
- [ ] Login Enterprise funcionando
- [ ] Logout funcionando
- [ ] Redirecionamento pós-login correto

### **✅ Interface e Permissões**
- [ ] Menu Admin completo
- [ ] Menu Manager sem billing
- [ ] Menu Lawyer com processos
- [ ] Menu Assistant limitado
- [ ] Tentativas de acesso negado bloqueadas

### **✅ Navegação e Performance**
- [ ] Todas as páginas carregando
- [ ] Navegação entre seções fluida
- [ ] Response times adequados
- [ ] Sem erros no console do browser

### **✅ Multi-tenancy**
- [ ] Dados isolados por tenant
- [ ] Cross-tenant access bloqueado
- [ ] Headers X-Tenant-ID funcionando

---

## 🚀 **APÓS OS TESTES**

### **Se Tudo Funcionou:**
```bash
# Executar bateria completa:
./EXECUTAR_TODOS_TESTES.sh

# Ou específicos:
./TESTAR_AUTENTICACAO.sh
./TESTAR_PLANOS.sh
```

### **Se Houveram Erros:**
1. **Documentar** todos os erros encontrados
2. **Capturar** screenshots dos problemas
3. **Coletar** logs específicos dos momentos de erro
4. **Listar** funcionalidades que não funcionaram

---

**🎯 Objetivo: Validar que o frontend está 100% funcional para autenticação e navegação básica antes de prosseguir com testes mais complexos!**