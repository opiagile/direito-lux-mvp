# 🎯 ESCLARECIMENTOS ARQUITETURA FINAL

## ❓ **SUAS DÚVIDAS ESPECÍFICAS RESPONDIDAS**

### **1. MCP SERVICE (Bot Luxia) - ACESSO CLIENTE**

#### **🤖 Quem acessa o bot Luxia?**
```yaml
Usuários do Bot:
✅ ADVOGADOS: Acesso total (admin do escritório)
✅ CLIENTES: Acesso limitado (apenas seus processos)
✅ FUNCIONÁRIOS: Acesso por permissão (assistentes, estagiários)
```

#### **🔐 Controle de Acesso por Usuário:**
```go
// Código real do controle de acesso
func (s *MCPService) ProcessBotMessage(ctx context.Context, msg *BotMessage) error {
    user := s.authenticateUser(ctx, msg.UserID)
    
    switch user.Role {
    case "ADVOGADO":
        // Acesso total: todos os processos do escritório
        return s.executeToolWithFullAccess(ctx, msg.Tool, msg.Params)
        
    case "CLIENTE":
        // Acesso limitado: apenas processos onde é parte
        processes := s.getClientProcesses(ctx, user.ID)
        return s.executeToolWithLimitedAccess(ctx, msg.Tool, processes)
        
    case "FUNCIONARIO":
        // Acesso por permissão: processos atribuídos
        permissions := s.getUserPermissions(ctx, user.ID)
        return s.executeToolWithPermissions(ctx, msg.Tool, permissions)
    }
}
```

#### **📱 Exemplos de Uso por Tipo de Usuário:**

**Cliente (João Silva):**
```
👤 João: "Luxia, como está meu processo de danos morais?"
🤖 Luxia: Olá João! Seu processo 1001234-56.2024.8.26.0100:
         📋 Status: Aguardando decisão
         📅 Última atualização: 2 dias atrás
         ⏰ Próxima audiência: 25/01/2025
         
         Alguma dúvida sobre o andamento?
```

**Advogado (Dr. Silva):**
```
👤 Dr. Silva: "Luxia, relatório de todos os processos ativos"
🤖 Luxia: Dr. Silva, encontrei 47 processos ativos:
         📊 15 TJSP, 8 TJRJ, 12 TRT2, 7 TRF3, 5 STJ
         🔥 3 processos com prazo vencendo (< 5 dias)
         💰 Valor total em discussão: R$ 2.3M
         
         Deseja detalhes dos processos urgentes?
```

---

### **2. DATAJUD SERVICE - CNPJ DA CONSULTA**

#### **🏛️ CNPJ utilizado nas consultas:**
```yaml
RESPOSTA: CNPJ DA SUA EMPRESA (não do cliente)

Motivo:
- DataJud CNJ exige CNPJ do REQUERENTE da consulta
- Sua empresa é o "usuário" da API CNJ
- Clientes não têm acesso direto ao CNJ
- Você atua como intermediário autorizado
```

#### **🔑 Configuração Real:**
```go
// Configuração do DataJud Service
type DataJudConfig struct {
    // CNPJ da sua empresa (Direito Lux)
    CompanyCNPJ: "12.345.678/0001-90",  // SEU CNPJ
    
    // Pool de CNPJs para escalar (se necessário)
    CNPJPool: []string{
        "12.345.678/0001-90",  // CNPJ principal
        "12.345.678/0002-71",  // CNPJ filial 1 (se houver)
        "12.345.678/0003-52",  // CNPJ filial 2 (se houver)
    },
    
    // Quota por CNPJ (10K/dia por CNPJ)
    DailyQuota: 10000,
}
```

#### **⚖️ Aspectos Legais:**
```yaml
Legitimidade:
✅ Empresa de software pode consultar CNJ
✅ Advocacia terceirizada é permitida
✅ Cliente autoriza consulta via contrato
✅ Dados são públicos (processos judiciais)

Auditoria:
✅ Todo acesso é logado com CNPJ da empresa
✅ CNJ rastreia quem consultou o que
✅ Relatórios mensais são obrigatórios
```

---

### **3. SEARCH SERVICE - LOCAL vs ONLINE**

#### **🔍 Tipos de Consulta:**
```yaml
RESPOSTA: HÍBRIDO (Local + Online)

Base Local (Search Service):
✅ Processos já monitorados
✅ Documentos internos
✅ Jurisprudência indexada
✅ Histórico de movimentos

Consulta Online (DataJud Service):
✅ Processos novos (primeira consulta)
✅ Atualizações em tempo real
✅ Verificação de mudanças
✅ Dados que não temos localmente
```

#### **🔄 Fluxo Híbrido Real:**
```go
// Código real do fluxo híbrido
func (s *SearchService) SearchProcess(ctx context.Context, processNumber string) (*ProcessData, error) {
    // 1. PRIMEIRO: Busca local (rápida)
    localResult := s.searchLocal(ctx, processNumber)
    if localResult != nil && localResult.IsRecent() {
        return localResult, nil // Retorna dados locais se recentes
    }
    
    // 2. SEGUNDO: Consulta online (DataJud)
    onlineResult := s.dataJudService.QueryProcess(ctx, processNumber)
    if onlineResult != nil {
        // 3. TERCEIRO: Atualiza base local
        s.updateLocalIndex(ctx, onlineResult)
        return onlineResult, nil
    }
    
    // 4. FALLBACK: Retorna dados locais (mesmo se antigos)
    return localResult, nil
}
```

---

### **4. ATUALIZAÇÃO DE PROCESSOS - ESTRATÉGIA**

#### **🕐 Atualização Automática (Background):**
```yaml
RESPOSTA: BACKGROUND AUTOMÁTICO + SOB DEMANDA

Polling Automático:
✅ Processos monitorados: A cada 30 minutos
✅ Processos urgentes: A cada 15 minutos
✅ Processos arquivados: A cada 24 horas
✅ Horário: 6h às 22h (quando tribunais atualizam)

Sob Demanda:
✅ Advogado solicita atualização específica
✅ Cliente pergunta sobre processo
✅ Bot detecta prazo vencendo
✅ Relatório é solicitado
```

#### **⚡ Processo de Atualização Automática:**
```go
// Código real do processo de atualização
func (s *MonitorService) StartBackgroundMonitoring(ctx context.Context) {
    // Polling a cada 30 minutos
    ticker := time.NewTicker(30 * time.Minute)
    
    for {
        select {
        case <-ticker.C:
            s.updateAllMonitoredProcesses(ctx)
        case <-ctx.Done():
            return
        }
    }
}

func (s *MonitorService) updateAllMonitoredProcesses(ctx context.Context) {
    // 1. Busca processos monitorados
    processes := s.getMonitoredProcesses(ctx)
    
    // 2. Atualiza em lotes (respeitando rate limit)
    for _, batch := range s.createBatches(processes, 50) {
        s.updateProcessBatch(ctx, batch)
        time.Sleep(2 * time.Second) // Rate limiting
    }
}

func (s *MonitorService) updateProcessBatch(ctx context.Context, processes []Process) {
    for _, process := range processes {
        // 3. Consulta DataJud
        newData := s.dataJudService.QueryProcess(ctx, process.Number)
        
        // 4. Detecta mudanças
        if s.detectChanges(process, newData) {
            // 5. Atualiza base local
            s.updateLocalProcess(ctx, newData)
            
            // 6. Dispara notificação
            s.notifyProcessUpdate(ctx, process, newData)
        }
    }
}
```

---

### **5. ARQUITETURA COMPLETA ESCLARECIDA**

#### **🏗️ Fluxo Completo de Dados:**
```yaml
1. CLIENTE pergunta no WhatsApp:
   "Luxia, como está meu processo X?"

2. MCP SERVICE (Bot Luxia):
   - Autentica cliente
   - Verifica permissões (apenas processos do cliente)
   - Executa ferramenta process_search

3. SEARCH SERVICE (Busca Híbrida):
   - Busca local primeiro (PostgreSQL + Elasticsearch)
   - Se dados recentes: retorna local
   - Se dados antigos: consulta DataJud

4. DATAJUD SERVICE (Consulta CNJ):
   - Usa CNPJ da sua empresa
   - Consulta API oficial CNJ
   - Atualiza base local

5. NOTIFICATION SERVICE:
   - Responde no WhatsApp
   - Salva histórico da conversa
   - Agenda próximas verificações

6. MONITOR SERVICE (Background):
   - Polling automático a cada 30min
   - Detecta mudanças
   - Notifica automaticamente
```

#### **📊 Dados por Fonte:**
```yaml
Base Local (PostgreSQL + Elasticsearch):
├── Processos monitorados (completos)
├── Movimentos históricos
├── Documentos internos
├── Jurisprudência indexada
├── Relatórios gerados
└── Conversas do bot

Consulta Online (DataJud CNJ):
├── Processos novos
├── Atualizações recentes
├── Verificação de mudanças
└── Dados que não temos

APIs Externas:
├── WhatsApp (mensagens)
├── Telegram (mensagens)
├── Email (notificações)
└── OpenAI/Ollama (IA)
```

---

## ✅ **RESPOSTAS FINAIS**

### **🤖 MCP (Bot Luxia):**
- **Quem acessa**: Advogados (total) + Clientes (limitado) + Funcionários (permissão)
- **Dados**: Base local + APIs internas (não consulta CNJ diretamente)
- **Controle**: Autenticação + autorização por usuário

### **🏛️ DataJud:**
- **CNPJ**: DA SUA EMPRESA (não do cliente)
- **Legitimidade**: Empresa de software pode consultar CNJ
- **Auditoria**: Todo acesso é rastreado pelo CNJ

### **🔍 Search:**
- **Tipo**: Híbrido (Local primeiro, Online se necessário)
- **Local**: Processos já monitorados + documentos internos
- **Online**: Processos novos + atualizações

### **🔄 Atualização:**
- **Background**: Polling automático a cada 30min (6h-22h)
- **Sob demanda**: Quando solicitado por advogado/cliente
- **Estratégia**: Prioriza dados locais, consulta CNJ quando necessário

---

**🎯 Essas explicações esclareceram suas dúvidas sobre a arquitetura?**