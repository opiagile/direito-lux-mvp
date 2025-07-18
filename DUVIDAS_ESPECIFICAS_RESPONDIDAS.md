# 🎯 DÚVIDAS ESPECÍFICAS RESPONDIDAS - DIREITO LUX

## 📋 **SUAS QUESTÕES TÉCNICAS IMPORTANTES**

---

## 1. 🏛️ **DATAJUD CNJ - INFORMAÇÕES RECEBIDAS E ARMAZENAMENTO**

### **📊 Dados Recebidos da API CNJ DataJud**

```json
{
  "hits": {
    "hits": [
      {
        "_source": {
          "numeroProcesso": "1001234-56.2024.8.26.0100",
          "classe": {
            "codigo": 436,
            "nome": "Ação de Cobrança"
          },
          "assunto": [
            {
              "codigo": 1396,
              "nome": "Responsabilidade Civil"
            }
          ],
          "tribunal": "TJSP",
          "grau": "1",
          "dataAjuizamento": "2024-01-15T10:30:00Z",
          "dataUltimaAtualizacao": "2024-07-15T14:22:00Z",
          "orgaoJulgador": {
            "codigo": 1234,
            "nome": "1ª Vara Cível Central"
          },
          "movimento": [
            {
              "codigo": 123,
              "nome": "Juntada de Documento",
              "dataHora": "2024-07-15T14:22:00Z",
              "complemento": "Certidão de intimação"
            }
          ],
          "partes": [
            {
              "tipo": "Autor",
              "nome": "JOÃO DA SILVA",
              "documento": "123.456.789-00"
            },
            {
              "tipo": "Réu", 
              "nome": "EMPRESA XYZ LTDA",
              "documento": "12.345.678/0001-90"
            }
          ],
          "advogados": [
            {
              "nome": "Dr. Pedro Advogado",
              "oab": "SP123456",
              "parte": "Autor"
            }
          ]
        }
      }
    ]
  }
}
```

### **🗄️ Estrutura de Armazenamento Local**

```sql
-- Tabela principal do processo
CREATE TABLE processes (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    numero_processo VARCHAR(25) UNIQUE NOT NULL,
    classe_codigo INTEGER,
    classe_nome VARCHAR(200),
    tribunal VARCHAR(10),
    grau INTEGER,
    data_ajuizamento TIMESTAMP,
    data_ultima_atualizacao TIMESTAMP,
    orgao_julgador_codigo INTEGER,
    orgao_julgador_nome VARCHAR(200),
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Tabela de movimentações (histórico)
CREATE TABLE process_movements (
    id UUID PRIMARY KEY,
    process_id UUID REFERENCES processes(id),
    codigo INTEGER,
    nome VARCHAR(500),
    data_hora TIMESTAMP,
    complemento TEXT,
    sequence_number INTEGER, -- Para ordenação
    created_at TIMESTAMP DEFAULT NOW()
);

-- Tabela de partes
CREATE TABLE process_parties (
    id UUID PRIMARY KEY,
    process_id UUID REFERENCES processes(id),
    tipo VARCHAR(50), -- Autor, Réu, etc.
    nome VARCHAR(200),
    documento VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Tabela de advogados
CREATE TABLE process_lawyers (
    id UUID PRIMARY KEY,
    process_id UUID REFERENCES processes(id),
    nome VARCHAR(200),
    oab VARCHAR(20),
    parte VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Tabela de consultas (auditoria)
CREATE TABLE datajud_queries (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    process_number VARCHAR(25),
    query_type VARCHAR(50), -- 'initial', 'update', 'movement'
    success BOOLEAN,
    response_size INTEGER,
    query_time_ms INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### **📝 Gerenciamento de Histórico**

**ESTRATÉGIA: HÍBRIDA (Update + Insert)**

```go
// Algoritmo de sincronização
func (s *DataJudService) SynchronizeProcess(processNumber string) error {
    // 1. Consultar API CNJ
    cnjData := s.queryDataJud(processNumber)
    
    // 2. Buscar processo local
    localProcess := s.getLocalProcess(processNumber)
    
    if localProcess == nil {
        // PRIMEIRA CONSULTA: Insert completo
        return s.insertNewProcess(cnjData)
    }
    
    // ATUALIZAÇÃO: Update processo + Insert novos movimentos
    if cnjData.DataUltimaAtualizacao.After(localProcess.UpdatedAt) {
        // Update dados principais
        s.updateProcessInfo(localProcess.ID, cnjData)
        
        // Insert apenas movimentos novos
        s.insertNewMovements(localProcess.ID, cnjData.Movements)
        
        // Log da consulta
        s.logQuery(processNumber, "update", true)
    }
    
    return nil
}
```

**VANTAGENS:**
- ✅ **Performance**: Só atualiza o que mudou
- ✅ **Histórico completo**: Mantém todos os movimentos
- ✅ **Auditoria**: Log de todas as consultas
- ✅ **Economia**: Evita consultas desnecessárias

---

## 2. 🔍 **VECTOR SEARCH - BASE DE DADOS**

### **📊 Estratégia de Vector Database**

**SOLUÇÃO: PostgreSQL + pgvector (NÃO Pinecone)**

```sql
-- Extensão pgvector no PostgreSQL
CREATE EXTENSION IF NOT EXISTS vector;

-- Tabela de embeddings
CREATE TABLE document_embeddings (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    document_type VARCHAR(50), -- 'process', 'jurisprudence', 'document'
    document_id UUID,
    content TEXT,
    embedding vector(1536), -- OpenAI embedding dimension
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Índice para busca por similaridade
CREATE INDEX ON document_embeddings USING ivfflat (embedding vector_cosine_ops);
```

### **🤖 Geração de Embeddings**

```go
// AI Service - Geração de embeddings
func (s *AIService) GenerateEmbedding(text string) ([]float32, error) {
    // Usar OpenAI text-embedding-ada-002
    response, err := s.openaiClient.CreateEmbedding(ctx, openai.EmbeddingRequest{
        Model: "text-embedding-ada-002",
        Input: text,
    })
    
    return response.Data[0].Embedding, nil
}

// Search Service - Busca por similaridade
func (s *SearchService) SemanticSearch(query string, limit int) ([]SearchResult, error) {
    queryEmbedding := s.aiService.GenerateEmbedding(query)
    
    sql := `
        SELECT document_id, content, metadata, 
               1 - (embedding <=> $1) as similarity
        FROM document_embeddings 
        WHERE tenant_id = $2
        ORDER BY embedding <=> $1 
        LIMIT $3
    `
    
    return s.db.Query(sql, queryEmbedding, tenantID, limit)
}
```

### **🎯 Por que PostgreSQL + pgvector?**

**VANTAGENS:**
- ✅ **Custo zero**: Sem APIs externas pagas
- ✅ **Controle total**: Dados ficam no nosso ambiente
- ✅ **Performance**: Busca local ultra-rápida
- ✅ **Integração**: Mesmo banco dos dados relacionais
- ✅ **Segurança**: Dados jurídicos não saem do ambiente

**VS Pinecone:**
- ❌ **Custo alto**: $0.096/1M queries
- ❌ **Vendor lock-in**: Dependência externa
- ❌ **Latência**: Requisições HTTP
- ❌ **Compliance**: Dados jurídicos em servidor externo

---

## 3. 🤖 **LUXIA - SUPORTE AO SISTEMA**

### **✅ SIM! Luxia pode tirar dúvidas sobre o sistema**

```go
// MCP Service - Ferramenta de suporte
type SystemHelpTool struct {
    knowledgeBase map[string]string
    faqs          []FAQ
}

func (t *SystemHelpTool) Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error) {
    question := params["question"].(string)
    userRole := params["user_role"].(string) // ADVOGADO, FUNCIONARIO, CLIENTE
    
    // Buscar resposta na base de conhecimento
    answer := t.searchKnowledgeBase(question, userRole)
    
    if answer == "" {
        // Gerar resposta usando IA
        answer = t.generateContextualAnswer(question, userRole)
    }
    
    return &ToolResult{
        Success: true,
        Data: map[string]interface{}{
            "answer": answer,
            "helpful_links": t.getHelpfulLinks(userRole),
            "next_steps": t.suggestNextSteps(question, userRole),
        },
    }, nil
}
```

### **📚 Base de Conhecimento por Papel**

```yaml
ADVOGADO:
  - "Como adicionar novo processo?"
  - "Como interpretar relatórios?"
  - "Como configurar notificações?"
  - "Como gerenciar usuários?"

FUNCIONARIO:
  - "Como buscar processos?"
  - "Como atualizar dados de cliente?"
  - "Como usar a busca de jurisprudência?"

CLIENTE:
  - "Como ver meus processos?"
  - "O que significa cada status?"
  - "Como receber notificações?"
  - "Como entrar em contato com meu advogado?"
```

### **💬 Exemplos de Interação**

```
👤 Cliente: "Luxia, como eu vejo o andamento do meu processo?"
🤖 Luxia: "Olá! Para ver seus processos:
          
          1. Acesse o menu 'Meus Processos'
          2. Clique no processo desejado
          3. Veja a timeline de movimentações
          
          📱 Você também recebe notificações automáticas 
          no WhatsApp sempre que houver atualizações!"

👤 Funcionário: "Como faço busca de jurisprudência?"
🤖 Luxia: "Para buscar jurisprudência:
          
          1. Acesse 'Pesquisa Jurisprudencial'
          2. Digite termos relacionados ao caso
          3. Use filtros (tribunal, data, área)
          
          💡 Dica: Use termos específicos como 'dano moral' 
          para resultados mais precisos!"
```

---

## 4. 📅 **AGENDAMENTO E NOTIFICAÇÕES**

### **🚨 IMPORTANTE: Funcionalidades NÃO implementadas**

As funcionalidades mencionadas **"Agendar reuniões"** e **"notificar sobre atualizações"** foram citadas como exemplos de capacidades potenciais, mas **NÃO estão implementadas no sistema atual**.

### **📋 Status Real das Funcionalidades**

```yaml
IMPLEMENTADO:
  ✅ Notificações de movimentações processuais
  ✅ Alertas de prazos
  ✅ Notificações multicanal (WhatsApp, Email, Telegram)
  ✅ Templates de notificação personalizáveis

NÃO IMPLEMENTADO:
  ❌ Agendamento de reuniões
  ❌ Calendário integrado
  ❌ Lembretes de compromissos
  ❌ Sincronização com Google Calendar
```

### **🔮 Implementação Futura (se necessário)**

```go
// Estrutura para agendamento (FUTURO)
type Appointment struct {
    ID          string    `json:"id"`
    TenantID    string    `json:"tenant_id"`
    ClientID    string    `json:"client_id"`
    LawyerID    string    `json:"lawyer_id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    StartTime   time.Time `json:"start_time"`
    EndTime     time.Time `json:"end_time"`
    Status      string    `json:"status"` // scheduled, confirmed, cancelled
    CreatedAt   time.Time `json:"created_at"`
}

// Notificação de agendamento (FUTURO)
func (s *NotificationService) SendAppointmentReminder(appointment *Appointment) error {
    message := fmt.Sprintf(`
🗓️ *Lembrete de Reunião*

📋 *Assunto:* %s
👤 *Cliente:* %s
🕐 *Horário:* %s
📍 *Local:* %s

Confirme sua presença respondendo SIM.
    `, appointment.Title, appointment.ClientName, 
       appointment.StartTime.Format("15:04"), appointment.Location)
    
    return s.whatsappProvider.SendMessage(appointment.ClientPhone, message)
}
```

---

## 5. 🚫 **STAGING REMOVIDO - CORREÇÃO DE DOCUMENTAÇÃO**

### **✅ Ambientes Definitivos**

```yaml
AMBIENTES:
  DEV (Local): 
    - Docker Compose
    - Dados reais (sem mocks)
    - Desenvolvimento e testes
    - URL: http://localhost:3000

  PRODUCTION (GCP):
    - Kubernetes (GKE)
    - Cloud SQL PostgreSQL
    - Dados reais
    - URL: https://app.direitolux.com.br
```

### **📝 Correção na Documentação**

**ANTES (incorreto):**
```
INFRASTRUCTURA DE STAGING
Endpoint de Teste: https://35.188.198.87
Staging deployment confirmado
```

**DEPOIS (correto):**
```
INFRASTRUCTURA DE PRODUÇÃO
Endpoint de Produção: https://app.direitolux.com.br
Deploy direto: DEV → PRODUCTION
```

---

## 6. 🚫 **SEM MOCKS - DADOS REAIS**

### **✅ Confirmação: Dados Reais em DEV**

```go
// DataJud Service - SEM MOCK
func (s *DataJudService) QueryProcess(processNumber string) (*ProcessData, error) {
    // SEMPRE usar HTTP client real
    return s.httpClient.QueryProcess(processNumber)
}

// AI Service - SEM MOCK
func (s *AIService) AnalyzeDocument(document string) (*Analysis, error) {
    // SEMPRE usar OpenAI API real
    return s.openaiClient.Analyze(document)
}

// Notification Service - SEM MOCK
func (s *NotificationService) SendWhatsApp(phone, message string) error {
    // SEMPRE usar WhatsApp Business API real
    return s.whatsappClient.SendMessage(phone, message)
}
```

### **🎯 Estratégia de Validação**

```yaml
DESENVOLVIMENTO:
  - Dados reais CNJ
  - APIs reais (com quotas de dev)
  - Testes com dados reais
  - Validação completa

PRODUÇÃO:
  - Mesma implementação
  - Sem surpresas
  - Zero refatoração
  - Deploy direto
```

---

## 📊 **RESUMO DAS RESPOSTAS**

| Dúvida | Resposta | Status |
|--------|----------|---------|
| **DataJud dados** | JSON completo + estrutura SQL detalhada | ✅ Esclarecido |
| **Histórico** | Híbrido: Update processo + Insert movimentos | ✅ Esclarecido |
| **Vector search** | PostgreSQL + pgvector (não Pinecone) | ✅ Esclarecido |
| **Luxia suporte** | SIM - base conhecimento por papel | ✅ Esclarecido |
| **Agendamento** | NÃO implementado (exemplo teórico) | ✅ Esclarecido |
| **Staging** | Removido - só DEV + PROD | ✅ Corrigido |
| **Mocks** | Removidos - dados reais sempre | ✅ Confirmado |

---

## 🎯 **PRÓXIMOS PASSOS**

1. **Corrigir** VALIDACAO_REGISTRO_COMPLETA.md (remover staging)
2. **Implementar** estrutura SQL detalhada do DataJud
3. **Configurar** pgvector no PostgreSQL
4. **Expandir** base de conhecimento da Luxia
5. **Documentar** fluxo sem mocks

**🔥 Todas as dúvidas específicas foram respondidas com implementação detalhada!**