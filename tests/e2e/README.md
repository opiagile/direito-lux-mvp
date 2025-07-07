# Testes E2E - Direito Lux

## 🎯 Visão Geral

Suite completa de testes End-to-End para validar o funcionamento integrado do sistema Direito Lux. Os testes cobrem desde autenticação multi-tenant até fluxos completos de negócio.

## 📋 Testes Implementados

### 🔑 Autenticação (`auth.test.js`)
- ✅ Login com todos os 4 tenants (Starter, Professional, Business, Enterprise)
- ✅ Validação de tokens JWT
- ✅ Isolamento multi-tenant
- ✅ Casos de erro (credenciais inválidas, headers ausentes)
- ✅ Refresh tokens
- ✅ Performance de autenticação

### 📋 CRUD de Processos (`processes.test.js`)
- ✅ Dashboard e estatísticas por tenant
- ✅ Listar processos com paginação e filtros
- ✅ Criar processos com validação CNJ
- ✅ Obter detalhes de processo específico
- ✅ Atualizar processos
- ✅ Deletar processos
- ✅ Isolamento cross-tenant
- ✅ Performance de operações

### 📊 Dashboard e Relatórios (`dashboard.test.js`)
- ✅ KPIs do Process Service
- ✅ Atividades recentes do Report Service
- ✅ Métricas adicionais do dashboard
- ✅ Health checks dos serviços
- ✅ Integração com frontend
- ✅ Performance do dashboard
- ✅ Graceful degradation

### 🔄 Fluxo Completo (`full-flow.test.js`)
- ✅ Fluxo principal: Registro → Monitoramento → Dashboard
- ✅ Operações multi-tenant simultâneas
- ✅ Performance do fluxo completo
- ✅ Testes de resiliência
- ✅ Consistência após operações simultâneas

## 🚀 Como Executar

### Pré-requisitos

1. **Serviços rodando**:
   ```bash
   cd services/
   ./scripts/deploy-dev.sh start
   ```

2. **Node.js** (versão 18+)

### Instalação

```bash
cd tests/e2e/
npm install
```

### Execução

#### Executar todos os testes
```bash
./run-tests.sh
```

#### Executar testes específicos
```bash
./run-tests.sh --auth           # Apenas autenticação
./run-tests.sh --processes      # Apenas processos
./run-tests.sh --dashboard      # Apenas dashboard
./run-tests.sh --full-flow      # Apenas fluxo completo
```

#### Opções úteis
```bash
./run-tests.sh --setup          # Verificar ambiente
./run-tests.sh --install        # Instalar dependências
./run-tests.sh --verbose        # Output detalhado
./run-tests.sh --help           # Ajuda
```

#### Comandos npm diretos
```bash
npm test                        # Todos os testes
npm run test:auth              # Testes de autenticação
npm run test:processes         # Testes de processos
npm run test:dashboard         # Testes de dashboard
npm run test:full-flow         # Testes de fluxo completo
npm run test:watch             # Modo watch
```

## 📊 Configuração

### Tenants de Teste

Os testes utilizam os seguintes tenants (conforme STATUS_IMPLEMENTACAO.md):

| Tenant | Email | Plano | Senha |
|--------|-------|-------|-------|
| Silva & Associados | admin@silvaassociados.com.br | Starter | password |
| Costa & Santos | admin@costasantos.com.br | Professional | password |
| Machado Advogados | admin@machadoadvogados.com.br | Business | password |
| Barros Enterprise | admin@barrosent.com.br | Enterprise | password |

### URLs dos Serviços

| Serviço | URL | Porta |
|---------|-----|-------|
| Auth Service | http://localhost:8081 | 8081 |
| Tenant Service | http://localhost:8082 | 8082 |
| Process Service | http://localhost:8083 | 8083 |
| DataJud Service | http://localhost:8084 | 8084 |
| Notification Service | http://localhost:8085 | 8085 |
| AI Service | http://localhost:8000 | 8000 |
| Search Service | http://localhost:8086 | 8086 |
| Report Service | http://localhost:8087 | 8087 |
| MCP Service | http://localhost:8088 | 8088 |

## 🧪 Estrutura dos Testes

```
tests/e2e/
├── utils/
│   ├── config.js          # Configurações (URLs, tenants, dados de teste)
│   ├── api-helper.js      # Helper para chamadas API
│   ├── setup.js           # Setup inicial dos testes
│   └── cleanup.js         # Limpeza após testes
├── auth.test.js           # Testes de autenticação
├── processes.test.js      # Testes de CRUD de processos
├── dashboard.test.js      # Testes de dashboard e relatórios
├── full-flow.test.js      # Testes de fluxo completo
├── package.json           # Dependências e scripts
├── run-tests.sh           # Script principal
└── README.md              # Esta documentação
```

## 🔧 Utilitários

### ApiHelper

Classe helper que facilita chamadas para APIs:

```javascript
const apiHelper = require('./utils/api-helper');

// Login
await apiHelper.login('silva');

// Chamadas autenticadas
const response = await apiHelper.get('process', '/api/v1/processes', 'silva');
await apiHelper.post('process', '/api/v1/processes', data, 'silva');
await apiHelper.put('process', '/api/v1/processes/123', data, 'silva');
await apiHelper.delete('process', '/api/v1/processes/123', 'silva');

// Health checks
await apiHelper.checkHealth('process');
await apiHelper.waitForService('auth', 10, 2000);
```

### Configuração

Arquivo `utils/config.js` centraliza:
- URLs dos serviços
- Dados dos tenants
- Headers padrão
- Timeouts
- Dados de teste

## 🎯 Casos de Teste Cobertos

### Funcionais
- ✅ Autenticação multi-tenant
- ✅ CRUD completo de processos
- ✅ Dashboard com dados reais
- ✅ Relatórios e atividades
- ✅ Isolamento entre tenants
- ✅ Validação de dados
- ✅ Casos de erro

### Não-funcionais
- ✅ Performance (latência < 500ms para KPIs)
- ✅ Resiliência (graceful degradation)
- ✅ Concorrência (operações simultâneas)
- ✅ Consistência (dados íntegros)
- ✅ Escalabilidade (múltiplos tenants)

### Integração
- ✅ Process Service ↔ Dashboard
- ✅ Auth Service ↔ Todos os serviços
- ✅ Report Service ↔ KPIs
- ✅ Multi-tenant em todos os serviços

## 📈 Métricas de Performance

### Benchmarks Esperados
- Login: < 2 segundos
- Validação de token: < 500ms
- KPIs do dashboard: < 300ms
- Listar processos: < 1 segundo
- Criar processo: < 2 segundos
- Fluxo completo: < 7 segundos

### Monitoramento
Os testes incluem medição de performance e falham se os tempos estiverem acima dos limites definidos.

## 🚨 Troubleshooting

### Serviços não disponíveis
```bash
# Verificar status
./run-tests.sh --setup

# Iniciar serviços
cd ../../services/
./scripts/deploy-dev.sh start

# Verificar logs
./scripts/deploy-dev.sh logs
```

### Falhas de autenticação
```bash
# Verificar se Auth Service está rodando
curl http://localhost:8081/health

# Verificar credenciais no config.js
# Verificar se tenants estão no banco
```

### Timeouts
```bash
# Aumentar timeouts no config.js se necessário
# Verificar performance dos serviços
# Verificar resources do sistema (CPU, RAM)
```

### Falhas de isolamento
```bash
# Verificar se X-Tenant-ID está sendo enviado
# Verificar se middleware multi-tenant está ativo
# Verificar logs dos serviços
```

## 🧹 Limpeza

Os testes fazem limpeza automática:
- Tokens de autenticação são limpos após cada suite
- Processos criados durante testes são removidos
- Estado é restaurado após operações

## 📋 Checklist de Validação

Após executar os testes, verificar:

- [ ] Todos os testes passaram
- [ ] Performance dentro dos limites
- [ ] Isolamento multi-tenant funcionando
- [ ] Dados criados foram limpos
- [ ] Serviços ainda estão saudáveis
- [ ] Logs não mostram erros críticos

## 🚀 Próximos Passos

1. **Integração CI/CD**: Executar testes automaticamente
2. **Testes de Carga**: Validar comportamento sob stress
3. **Testes de Segurança**: Pentesting automatizado
4. **Testes de Mobile**: Validar app React Native
5. **Testes de Produção**: Smoke tests em produção

## 📞 Suporte

Em caso de problemas:

1. Verificar logs: `./run-tests.sh --verbose`
2. Verificar ambiente: `./run-tests.sh --setup`
3. Reiniciar serviços: `cd ../../services && ./scripts/deploy-dev.sh restart`
4. Consultar documentação dos serviços individuais

---

**🧪 Suite E2E Direito Lux - Validação completa do sistema integrado**