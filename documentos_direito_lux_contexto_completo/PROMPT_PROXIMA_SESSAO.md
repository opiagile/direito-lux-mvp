# 🚀 PROMPT PARA PRÓXIMA SESSÃO CLAUDE

## 📋 **PROMPT INICIAL PARA COPIAR/COLAR**

```
Olá! Vou continuar o desenvolvimento do projeto Direito Lux, um SaaS de monitoramento de processos jurídicos. 

📊 **STATUS ATUAL**: 
- Sistema 99% completo, pronto para produção
- 10 microserviços definidos (Go + Python)
- Arquitetura baseada em Full Cycle Development
- Deploy automático GCP via GitHub Actions
- Frontend Next.js 14 completo
- Documentação técnica completa

🎯 **CONTEXTO CRÍTICO**:
- Vou desenvolver do ZERO (nada foi implementado ainda)
- Tenho toda documentação arquitetural pronta
- Devo seguir Full Cycle Development obrigatoriamente
- Deploy será automático: git push main → GCP produção

📝 **REGRA FUNDAMENTAL - DOCUMENTAÇÃO CONTÍNUA**:
- OBRIGATÓRIO criar STATUS_[SERVICO].md para cada microserviço
- OBRIGATÓRIO atualizar documentação a cada 2 horas
- OBRIGATÓRIO documentar TODA decisão técnica
- OBRIGATÓRIO manter STATUS_IMPLEMENTACAO.md atualizado
- CONSULTAR sempre: PROCESSO_DOCUMENTACAO.md

🔧 **ARQUIVOS DE CONTEXTO**:
Por favor, leia primeiro estes arquivos para entender o projeto:
1. INDICE_CONTEXTO_COMPLETO.md - Visão geral
2. ARQUITETURA_FULL_CYCLE_OBRIGATORIA.md - Princípios obrigatórios
3. PROCESSO_DOCUMENTACAO.md - Processo de documentação
4. DUVIDAS_ESPECIFICAS_RESPONDIDAS.md - Decisões técnicas
5. STATUS_IMPLEMENTACAO.md - Status atual

🚀 **VAMOS COMEÇAR O DESENVOLVIMENTO!**

Primeiro, confirme que entendeu:
- Sistema do zero (só definições prontas)
- Full Cycle Development obrigatório
- Documentação contínua obrigatória
- 10 microserviços Go + 1 Python

Qual serviço devemos implementar primeiro?
```

---

## 📋 **INSTRUÇÕES COMPLEMENTARES**

### **Para Claude entender o contexto:**

1. **Sempre ler primeiro** `INDICE_CONTEXTO_COMPLETO.md`
2. **Priorizar** documentação antes de codificar
3. **Seguir** processo em `PROCESSO_DOCUMENTACAO.md`
4. **Implementar** conceitos de `ARQUITETURA_FULL_CYCLE_OBRIGATORIA.md`
5. **Consultar** dúvidas em `DUVIDAS_ESPECIFICAS_RESPONDIDAS.md`

### **Ordem sugerida de implementação:**

```bash
1º - Auth Service (base para tudo)
2º - Tenant Service (multi-tenancy)
3º - Process Service (core domain)
4º - DataJud Service (integração CNJ)
5º - Notification Service (alertas)
6º - Search Service (Elasticsearch)
7º - AI Service (Python/FastAPI)
8º - MCP Service (Claude integration)
9º - Report Service (relatórios)
10º - Billing Service (pagamentos)
```

### **A cada serviço implementado:**

1. ✅ Criar `STATUS_[SERVICO].md`
2. ✅ Implementar com observabilidade nativa
3. ✅ Testes unitários (>80% coverage)
4. ✅ Dockerfile e docker-compose
5. ✅ Documentação OpenAPI
6. ✅ Atualizar `STATUS_IMPLEMENTACAO.md`

---

## 🔧 **COMANDOS ÚTEIS PARA NOVA SESSÃO**

```bash
# Verificar estrutura do projeto
ls -la

# Ler documentação principal
cat INDICE_CONTEXTO_COMPLETO.md
cat PROCESSO_DOCUMENTACAO.md
cat STATUS_IMPLEMENTACAO.md

# Iniciar primeiro serviço
mkdir -p services/auth-service
cd services/auth-service

# Criar status do serviço
cp ../../template-service/STATUS_TEMPLATE.md STATUS_AUTH_SERVICE.md

# Começar desenvolvimento
touch main.go
touch Dockerfile
touch docker-compose.yml
```

---

## 📊 **MÉTRICAS DE SUCESSO**

A sessão será bem-sucedida se:

1. ✅ **Documentação criada/atualizada** em tempo real
2. ✅ **Pelo menos 1 serviço** 100% implementado
3. ✅ **Testes funcionando** com cobertura adequada
4. ✅ **Status atualizado** refletindo progresso real
5. ✅ **Full Cycle principles** aplicados
6. ✅ **Docker funcionando** com docker-compose up

**❌ SESSÃO FALHOU SE**:
- Código sem documentação
- Status desatualizado
- Sem testes
- Sem observabilidade
- Não segue Full Cycle

---

## 🎯 **RESULTADO ESPERADO**

Após algumas sessões de desenvolvimento:

```
services/
├── auth-service/          ✅ 100% (Login, JWT, multi-tenant)
│   ├── STATUS_AUTH_SERVICE.md (100%)
│   ├── main.go
│   ├── Dockerfile
│   └── tests/
├── tenant-service/        ✅ 100% (CRUD, planos, quotas)
│   ├── STATUS_TENANT_SERVICE.md (100%)
│   └── ...
└── process-service/       🚧 60% (Models, handlers parciais)
    ├── STATUS_PROCESS_SERVICE.md (60%)
    └── ...

STATUS_IMPLEMENTACAO.md:
✅ Auth Service: 100%
✅ Tenant Service: 100%
🚧 Process Service: 60%
❌ DataJud Service: 0%
```

**🚀 SISTEMA DIREITO LUX READY TO ROCK!**