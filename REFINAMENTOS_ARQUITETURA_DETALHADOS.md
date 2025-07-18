# 🎯 REFINAMENTOS ARQUITETURA DETALHADOS

## 📋 **SUAS QUESTÕES ESPECÍFICAS RESPONDIDAS**

---

## 1. 🤖 **MCP SERVICE - DOCUMENTAÇÃO DE ACESSO**

### **📚 Documentação Completa do Controle de Acesso**

```go
// internal/domain/access_control.go
package domain

import (
    "context"
    "errors"
    "time"
)

// UserRole define os tipos de usuário
type UserRole string

const (
    RoleAdvogado    UserRole = "ADVOGADO"      // Acesso total
    RoleCliente     UserRole = "CLIENTE"       // Acesso limitado
    RoleFuncionario UserRole = "FUNCIONARIO"   // Acesso por permissão
    RoleEstagiario  UserRole = "ESTAGIARIO"    // Acesso supervisionado
)

// AccessLevel define níveis de acesso
type AccessLevel string

const (
    AccessFull       AccessLevel = "FULL"        // Todos os processos
    AccessLimited    AccessLevel = "LIMITED"     // Apenas processos próprios
    AccessPermission AccessLevel = "PERMISSION"  // Processos atribuídos
    AccessSupervised AccessLevel = "SUPERVISED"  // Acesso supervisionado
)

// UserAccess define permissões de acesso
type UserAccess struct {
    UserID          string        `json:"user_id"`
    TenantID        string        `json:"tenant_id"`
    Role            UserRole      `json:"role"`
    AccessLevel     AccessLevel   `json:"access_level"`
    AllowedTools    []string      `json:"allowed_tools"`
    AllowedProcesses []string     `json:"allowed_processes"`
    Supervisor      *string       `json:"supervisor,omitempty"`
    ExpiresAt       *time.Time    `json:"expires_at,omitempty"`
    CreatedAt       time.Time     `json:"created_at"`
    UpdatedAt       time.Time     `json:"updated_at"`
}

// ToolPermission define permissões por ferramenta
type ToolPermission struct {
    ToolName    string   `json:"tool_name"`
    AllowedRoles []UserRole `json:"allowed_roles"`
    RequiredPermissions []string `json:"required_permissions"`
}

// GetToolPermissions retorna permissões por ferramenta
func GetToolPermissions() map[string]ToolPermission {
    return map[string]ToolPermission{
        "process_search": {
            ToolName: "process_search",
            AllowedRoles: []UserRole{RoleAdvogado, RoleCliente, RoleFuncionario},
            RequiredPermissions: []string{"read_processes"},
        },
        "process_create": {
            ToolName: "process_create",
            AllowedRoles: []UserRole{RoleAdvogado, RoleFuncionario},
            RequiredPermissions: []string{"create_processes"},
        },
        "jurisprudence_search": {
            ToolName: "jurisprudence_search",
            AllowedRoles: []UserRole{RoleAdvogado, RoleFuncionario, RoleEstagiario},
            RequiredPermissions: []string{"read_jurisprudence"},
        },
        "generate_report": {
            ToolName: "generate_report",
            AllowedRoles: []UserRole{RoleAdvogado},
            RequiredPermissions: []string{"generate_reports"},
        },
        "bulk_notification": {
            ToolName: "bulk_notification",
            AllowedRoles: []UserRole{RoleAdvogado},
            RequiredPermissions: []string{"send_notifications"},
        },
    }
}

// AccessControlService serviço de controle de acesso
type AccessControlService struct {
    userRepository UserRepository
    auditLogger    AuditLogger
}

// ValidateToolAccess valida se usuário pode usar ferramenta
func (s *AccessControlService) ValidateToolAccess(ctx context.Context, userID, toolName string, params map[string]interface{}) error {
    // 1. Obter dados do usuário
    user, err := s.userRepository.GetUserAccess(ctx, userID)
    if err != nil {
        return err
    }
    
    // 2. Verificar permissões da ferramenta
    toolPermission, exists := GetToolPermissions()[toolName]
    if !exists {
        return errors.New("ferramenta não existe")
    }
    
    // 3. Verificar se role tem acesso
    if !s.hasRoleAccess(user.Role, toolPermission.AllowedRoles) {
        return errors.New("role não tem acesso a esta ferramenta")
    }
    
    // 4. Validar parâmetros específicos por role
    if err := s.validateRoleSpecificParams(ctx, user, toolName, params); err != nil {
        return err
    }
    
    // 5. Log de auditoria
    s.auditLogger.LogToolAccess(ctx, userID, toolName, params)
    
    return nil
}

// validateRoleSpecificParams valida parâmetros específicos por role
func (s *AccessControlService) validateRoleSpecificParams(ctx context.Context, user *UserAccess, toolName string, params map[string]interface{}) error {
    switch user.Role {
    case RoleCliente:
        // Cliente só pode acessar processos onde é parte
        if processID, ok := params["process_id"]; ok {
            if !s.isClientProcess(ctx, user.UserID, processID.(string)) {
                return errors.New("cliente não tem acesso a este processo")
            }
        }
        
    case RoleFuncionario:
        // Funcionário só pode acessar processos atribuídos
        if processID, ok := params["process_id"]; ok {
            if !s.isAssignedProcess(ctx, user.UserID, processID.(string)) {
                return errors.New("funcionário não tem acesso a este processo")
            }
        }
        
    case RoleEstagiario:
        // Estagiário precisa de supervisão
        if user.Supervisor == nil {
            return errors.New("estagiário precisa de supervisor")
        }
        
    case RoleAdvogado:
        // Advogado tem acesso total (sem validação adicional)
        break
    }
    
    return nil
}
```

### **🔐 Implementação de Controle por Mensagem**

```go
// internal/infrastructure/mcp/message_handler.go
func (h *MessageHandler) ProcessBotMessage(ctx context.Context, msg *BotMessage) error {
    // 1. Identificar usuário pela mensagem
    userID := h.identifyUser(ctx, msg)
    
    // 2. Validar acesso à ferramenta
    if err := h.accessControl.ValidateToolAccess(ctx, userID, msg.ToolName, msg.Params); err != nil {
        return h.sendAccessDeniedMessage(ctx, msg.ChatID, err.Error())
    }
    
    // 3. Executar ferramenta com contexto de usuário
    result, err := h.toolService.ExecuteWithUserContext(ctx, userID, msg.ToolName, msg.Params)
    if err != nil {
        return h.sendErrorMessage(ctx, msg.ChatID, err.Error())
    }
    
    // 4. Enviar resposta
    return h.sendSuccessMessage(ctx, msg.ChatID, result)
}
```

---

## 2. 🏛️ **DATAJUD SERVICE - GERENCIAMENTO MÚLTIPLOS CNPJs**

### **📋 Sistema de Pool de CNPJs Avançado**

```go
// internal/domain/cnpj_pool.go
package domain

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// CNPJProvider representa um CNPJ para consultas
type CNPJProvider struct {
    ID           string    `json:"id"`
    CNPJ         string    `json:"cnpj"`
    CompanyName  string    `json:"company_name"`
    DailyQuota   int       `json:"daily_quota"`
    UsedToday    int       `json:"used_today"`
    LastReset    time.Time `json:"last_reset"`
    IsActive     bool      `json:"is_active"`
    Priority     int       `json:"priority"`        // 1=alta, 2=média, 3=baixa
    RateLimit    int       `json:"rate_limit"`      // Requests por minuto
    LastRequest  time.Time `json:"last_request"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

// CNPJPoolManager gerencia pool de CNPJs
type CNPJPoolManager struct {
    providers    map[string]*CNPJProvider
    strategy     PoolStrategy
    mutex        sync.RWMutex
    repository   CNPJRepository
    metrics      MetricsCollector
}

// PoolStrategy define estratégias de seleção
type PoolStrategy string

const (
    StrategyRoundRobin  PoolStrategy = "round_robin"
    StrategyLeastUsed   PoolStrategy = "least_used"
    StrategyPriority    PoolStrategy = "priority"
    StrategyAvailable   PoolStrategy = "available"
)

// NewCNPJPoolManager cria novo gerenciador
func NewCNPJPoolManager(repository CNPJRepository, metrics MetricsCollector) *CNPJPoolManager {
    return &CNPJPoolManager{
        providers:  make(map[string]*CNPJProvider),
        strategy:   StrategyLeastUsed,
        repository: repository,
        metrics:    metrics,
    }
}

// InitializePool inicializa pool com CNPJs da empresa
func (m *CNPJPoolManager) InitializePool(ctx context.Context) error {
    // CNPJs da sua empresa (configuração inicial)
    initialCNPJs := []CNPJProvider{
        {
            ID:          "primary",
            CNPJ:        "12.345.678/0001-90",
            CompanyName: "Direito Lux LTDA",
            DailyQuota:  10000,
            IsActive:    true,
            Priority:    1, // Alta prioridade
            RateLimit:   120, // 120 req/min
        },
        {
            ID:          "secondary",
            CNPJ:        "12.345.678/0002-71",
            CompanyName: "Direito Lux Filial",
            DailyQuota:  10000,
            IsActive:    true,
            Priority:    2, // Média prioridade
            RateLimit:   120,
        },
    }
    
    // Adicionar CNPJs ao pool
    for _, cnpj := range initialCNPJs {
        if err := m.AddCNPJ(ctx, cnpj); err != nil {
            return fmt.Errorf("erro ao adicionar CNPJ %s: %w", cnpj.CNPJ, err)
        }
    }
    
    return nil
}

// AddCNPJ adiciona novo CNPJ ao pool
func (m *CNPJPoolManager) AddCNPJ(ctx context.Context, cnpj CNPJProvider) error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    // Validar CNPJ
    if !m.isValidCNPJ(cnpj.CNPJ) {
        return fmt.Errorf("CNPJ inválido: %s", cnpj.CNPJ)
    }
    
    // Verificar se já existe
    if _, exists := m.providers[cnpj.ID]; exists {
        return fmt.Errorf("CNPJ %s já existe no pool", cnpj.ID)
    }
    
    // Definir timestamps
    cnpj.CreatedAt = time.Now()
    cnpj.UpdatedAt = time.Now()
    cnpj.LastReset = time.Now()
    
    // Salvar no banco
    if err := m.repository.CreateCNPJ(ctx, &cnpj); err != nil {
        return err
    }
    
    // Adicionar ao pool
    m.providers[cnpj.ID] = &cnpj
    
    m.metrics.IncrementCounter("cnpj_pool_added", map[string]string{
        "cnpj_id": cnpj.ID,
    })
    
    return nil
}

// GetAvailableCNPJ retorna CNPJ disponível baseado na estratégia
func (m *CNPJPoolManager) GetAvailableCNPJ(ctx context.Context) (*CNPJProvider, error) {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    // Resetar quotas diárias se necessário
    m.resetDailyQuotasIfNeeded()
    
    // Filtrar CNPJs disponíveis
    available := m.getAvailableCNPJs()
    if len(available) == 0 {
        return nil, fmt.Errorf("nenhum CNPJ disponível no pool")
    }
    
    // Selecionar baseado na estratégia
    var selected *CNPJProvider
    switch m.strategy {
    case StrategyRoundRobin:
        selected = m.selectRoundRobin(available)
    case StrategyLeastUsed:
        selected = m.selectLeastUsed(available)
    case StrategyPriority:
        selected = m.selectByPriority(available)
    default:
        selected = available[0]
    }
    
    return selected, nil
}

// UseCNPJ marca CNPJ como usado
func (m *CNPJPoolManager) UseCNPJ(ctx context.Context, cnpjID string) error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    provider, exists := m.providers[cnpjID]
    if !exists {
        return fmt.Errorf("CNPJ %s não encontrado no pool", cnpjID)
    }
    
    // Incrementar uso
    provider.UsedToday++
    provider.LastRequest = time.Now()
    provider.UpdatedAt = time.Now()
    
    // Salvar no banco
    if err := m.repository.UpdateCNPJ(ctx, provider); err != nil {
        return err
    }
    
    // Métricas
    m.metrics.IncrementCounter("cnpj_usage", map[string]string{
        "cnpj_id": cnpjID,
    })
    
    return nil
}

// getAvailableCNPJs retorna CNPJs disponíveis
func (m *CNPJPoolManager) getAvailableCNPJs() []*CNPJProvider {
    var available []*CNPJProvider
    
    for _, provider := range m.providers {
        if provider.IsActive && provider.UsedToday < provider.DailyQuota {
            // Verificar rate limit
            if time.Since(provider.LastRequest) >= time.Minute/time.Duration(provider.RateLimit) {
                available = append(available, provider)
            }
        }
    }
    
    return available
}

// selectLeastUsed seleciona CNPJ menos usado
func (m *CNPJPoolManager) selectLeastUsed(available []*CNPJProvider) *CNPJProvider {
    if len(available) == 0 {
        return nil
    }
    
    least := available[0]
    for _, provider := range available[1:] {
        if provider.UsedToday < least.UsedToday {
            least = provider
        }
    }
    
    return least
}

// selectByPriority seleciona por prioridade
func (m *CNPJPoolManager) selectByPriority(available []*CNPJProvider) *CNPJProvider {
    if len(available) == 0 {
        return nil
    }
    
    highest := available[0]
    for _, provider := range available[1:] {
        if provider.Priority < highest.Priority { // Menor número = maior prioridade
            highest = provider
        }
    }
    
    return highest
}

// resetDailyQuotasIfNeeded reseta quotas diárias se mudou o dia
func (m *CNPJPoolManager) resetDailyQuotasIfNeeded() {
    now := time.Now()
    today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
    
    for _, provider := range m.providers {
        lastReset := time.Date(
            provider.LastReset.Year(),
            provider.LastReset.Month(),
            provider.LastReset.Day(),
            0, 0, 0, 0,
            provider.LastReset.Location(),
        )
        
        if today.After(lastReset) {
            provider.UsedToday = 0
            provider.LastReset = now
            provider.UpdatedAt = now
        }
    }
}

// GetPoolStatus retorna status do pool
func (m *CNPJPoolManager) GetPoolStatus(ctx context.Context) map[string]interface{} {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    status := map[string]interface{}{
        "total_cnpjs": len(m.providers),
        "strategy":    m.strategy,
        "providers":   make([]map[string]interface{}, 0),
    }
    
    totalQuota := 0
    totalUsed := 0
    activeCount := 0
    
    for _, provider := range m.providers {
        if provider.IsActive {
            activeCount++
        }
        totalQuota += provider.DailyQuota
        totalUsed += provider.UsedToday
        
        status["providers"] = append(status["providers"].([]map[string]interface{}), map[string]interface{}{
            "id":           provider.ID,
            "cnpj":         provider.CNPJ,
            "company_name": provider.CompanyName,
            "daily_quota":  provider.DailyQuota,
            "used_today":   provider.UsedToday,
            "available":    provider.DailyQuota - provider.UsedToday,
            "usage_percent": float64(provider.UsedToday) / float64(provider.DailyQuota) * 100,
            "is_active":    provider.IsActive,
            "priority":     provider.Priority,
            "last_request": provider.LastRequest,
        })
    }
    
    status["summary"] = map[string]interface{}{
        "active_cnpjs":     activeCount,
        "total_quota":      totalQuota,
        "total_used":       totalUsed,
        "total_available":  totalQuota - totalUsed,
        "usage_percent":    float64(totalUsed) / float64(totalQuota) * 100,
    }
    
    return status
}
```

### **🔧 Configuração Inicial**

```go
// internal/infrastructure/config/datajud_config.go
type DataJudConfig struct {
    // Pool de CNPJs
    CNPJPool struct {
        Strategy         PoolStrategy `envconfig:"CNPJ_POOL_STRATEGY" default:"least_used"`
        AutoAddCNPJs     bool         `envconfig:"AUTO_ADD_CNPJS" default:"true"`
        MaxCNPJs         int          `envconfig:"MAX_CNPJS" default:"10"`
        DefaultQuota     int          `envconfig:"DEFAULT_QUOTA" default:"10000"`
        DefaultRateLimit int          `envconfig:"DEFAULT_RATE_LIMIT" default:"120"`
    }
    
    // CNPJs da empresa (configuração inicial)
    CompanyCNPJs []struct {
        ID          string `envconfig:"CNPJ_ID"`
        CNPJ        string `envconfig:"CNPJ"`
        CompanyName string `envconfig:"COMPANY_NAME"`
        Priority    int    `envconfig:"PRIORITY" default:"1"`
        IsActive    bool   `envconfig:"IS_ACTIVE" default:"true"`
    }
}
```

---

## 3. 🔍 **SEARCH SERVICE - JURISPRUDÊNCIA**

### **📚 Busca de Jurisprudência para Embasamento**

```go
// internal/domain/jurisprudence_search.go
package domain

import (
    "context"
    "time"
)

// JurisprudenceType tipos de jurisprudência
type JurisprudenceType string

const (
    JurisprudenceSumula     JurisprudenceType = "SUMULA"
    JurisprudenceAcordao    JurisprudenceType = "ACORDAO"
    JurisprudenceDecisao    JurisprudenceType = "DECISAO"
    JurisprudenceRepercussao JurisprudenceType = "REPERCUSSAO"
    JurisprudenceInfome     JurisprudenceType = "INFORME"
)

// JurisprudenceSource fontes de jurisprudência
type JurisprudenceSource string

const (
    SourceSTF  JurisprudenceSource = "STF"
    SourceSTJ  JurisprudenceSource = "STJ"
    SourceTST  JurisprudenceSource = "TST"
    SourceTSE  JurisprudenceSource = "TSE"
    SourceTJSP JurisprudenceSource = "TJSP"
    SourceTJRJ JurisprudenceSource = "TJRJ"
    SourceTRF1 JurisprudenceSource = "TRF1"
    SourceTRT2 JurisprudenceSource = "TRT2"
)

// JurisprudenceEntry entrada de jurisprudência
type JurisprudenceEntry struct {
    ID              string              `json:"id"`
    Number          string              `json:"number"`
    Type            JurisprudenceType   `json:"type"`
    Source          JurisprudenceSource `json:"source"`
    Title           string              `json:"title"`
    Content         string              `json:"content"`
    Summary         string              `json:"summary"`
    Keywords        []string            `json:"keywords"`
    LegalArea       string              `json:"legal_area"`
    Rapporteur      string              `json:"rapporteur"`
    DecisionDate    time.Time           `json:"decision_date"`
    PublicationDate time.Time           `json:"publication_date"`
    SimilarCases    []string            `json:"similar_cases"`
    Precedents      []string            `json:"precedents"`
    Relevance       float64             `json:"relevance"`
    Vector          []float64           `json:"vector"`
    CreatedAt       time.Time           `json:"created_at"`
    UpdatedAt       time.Time           `json:"updated_at"`
}

// JurisprudenceSearchRequest requisição de busca
type JurisprudenceSearchRequest struct {
    Query           string                `json:"query"`
    CaseContext     string                `json:"case_context"`
    LegalArea       string                `json:"legal_area"`
    Sources         []JurisprudenceSource `json:"sources"`
    Types           []JurisprudenceType   `json:"types"`
    DateFrom        *time.Time            `json:"date_from"`
    DateTo          *time.Time            `json:"date_to"`
    SimilarToCase   string                `json:"similar_to_case"`
    MinRelevance    float64               `json:"min_relevance"`
    MaxResults      int                   `json:"max_results"`
    UseSemanticSearch bool                `json:"use_semantic_search"`
}

// JurisprudenceSearchResponse resposta de busca
type JurisprudenceSearchResponse struct {
    Query           string                `json:"query"`
    Results         []JurisprudenceEntry  `json:"results"`
    Total           int                   `json:"total"`
    ProcessTime     int                   `json:"process_time_ms"`
    SuggestedQueries []string             `json:"suggested_queries"`
    RelatedCases    []string              `json:"related_cases"`
}

// JurisprudenceSearchService serviço de busca de jurisprudência
type JurisprudenceSearchService struct {
    elasticClient   ElasticsearchClient
    vectorService   VectorService
    aiService       AIService
    cacheService    CacheService
    repository      JurisprudenceRepository
}

// SearchJurisprudence busca jurisprudência com IA
func (s *JurisprudenceSearchService) SearchJurisprudence(ctx context.Context, req *JurisprudenceSearchRequest) (*JurisprudenceSearchResponse, error) {
    // 1. Verificar cache
    cacheKey := s.buildCacheKey(req)
    if cached := s.cacheService.Get(ctx, cacheKey); cached != nil {
        return cached.(*JurisprudenceSearchResponse), nil
    }
    
    startTime := time.Now()
    
    // 2. Busca semântica com IA
    var results []JurisprudenceEntry
    var err error
    
    if req.UseSemanticSearch {
        results, err = s.searchSemantic(ctx, req)
    } else {
        results, err = s.searchTraditional(ctx, req)
    }
    
    if err != nil {
        return nil, err
    }
    
    // 3. Enriquecer com contexto do caso
    enrichedResults := s.enrichWithCaseContext(ctx, results, req.CaseContext)
    
    // 4. Gerar sugestões relacionadas
    suggestions := s.generateSuggestions(ctx, req.Query, results)
    
    // 5. Encontrar casos relacionados
    relatedCases := s.findRelatedCases(ctx, req.SimilarToCase, results)
    
    response := &JurisprudenceSearchResponse{
        Query:           req.Query,
        Results:         enrichedResults,
        Total:           len(enrichedResults),
        ProcessTime:     int(time.Since(startTime).Milliseconds()),
        SuggestedQueries: suggestions,
        RelatedCases:    relatedCases,
    }
    
    // 6. Salvar no cache
    s.cacheService.Set(ctx, cacheKey, response, 30*time.Minute)
    
    return response, nil
}

// searchSemantic busca semântica usando vectores
func (s *JurisprudenceSearchService) searchSemantic(ctx context.Context, req *JurisprudenceSearchRequest) ([]JurisprudenceEntry, error) {
    // 1. Gerar vetor da query
    queryVector, err := s.vectorService.GenerateVector(ctx, req.Query)
    if err != nil {
        return nil, err
    }
    
    // 2. Busca por similaridade vetorial
    query := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []interface{}{
                    map[string]interface{}{
                        "script_score": map[string]interface{}{
                            "query": map[string]interface{}{
                                "bool": map[string]interface{}{
                                    "filter": s.buildFilters(req),
                                },
                            },
                            "script": map[string]interface{}{
                                "source": "cosineSimilarity(params.query_vector, 'vector') + 1.0",
                                "params": map[string]interface{}{
                                    "query_vector": queryVector,
                                },
                            },
                        },
                    },
                },
            },
        },
        "size": req.MaxResults,
    }
    
    return s.elasticClient.SearchJurisprudence(ctx, query)
}

// FindSimilarCases encontra casos similares para embasamento
func (s *JurisprudenceSearchService) FindSimilarCases(ctx context.Context, caseNumber string, context string) ([]JurisprudenceEntry, error) {
    // 1. Obter dados do caso
    caseData, err := s.repository.GetProcessData(ctx, caseNumber)
    if err != nil {
        return nil, err
    }
    
    // 2. Extrair características do caso
    caseFeatures := s.extractCaseFeatures(caseData)
    
    // 3. Buscar jurisprudência similar
    searchReq := &JurisprudenceSearchRequest{
        Query:           fmt.Sprintf("%s %s", caseFeatures.LegalArea, caseFeatures.MainIssue),
        CaseContext:     context,
        LegalArea:       caseFeatures.LegalArea,
        UseSemanticSearch: true,
        MaxResults:      20,
        MinRelevance:    0.7,
    }
    
    response, err := s.SearchJurisprudence(ctx, searchReq)
    if err != nil {
        return nil, err
    }
    
    // 4. Filtrar por relevância para o caso específico
    relevantCases := s.filterByRelevance(response.Results, caseFeatures)
    
    return relevantCases, nil
}

// GenerateAppealBasis gera embasamento para recurso
func (s *JurisprudenceSearchService) GenerateAppealBasis(ctx context.Context, caseNumber string, appealType string) (*AppealBasis, error) {
    // 1. Analisar decisão atual
    caseData, err := s.repository.GetProcessData(ctx, caseNumber)
    if err != nil {
        return nil, err
    }
    
    // 2. Identificar pontos controvertidos
    controversialPoints := s.identifyControversialPoints(caseData)
    
    // 3. Buscar jurisprudência favorável
    var allPrecedents []JurisprudenceEntry
    
    for _, point := range controversialPoints {
        precedents, err := s.FindSimilarCases(ctx, caseNumber, point.Context)
        if err != nil {
            continue
        }
        
        // Filtrar apenas jurisprudência favorável
        favorablePrecedents := s.filterFavorablePrecedents(precedents, point.Position)
        allPrecedents = append(allPrecedents, favorablePrecedents...)
    }
    
    // 4. Gerar estrutura de embasamento
    basis := &AppealBasis{
        CaseNumber:          caseNumber,
        AppealType:          appealType,
        ControversialPoints: controversialPoints,
        SupportingPrecedents: allPrecedents,
        RecommendedStrategy: s.generateAppealStrategy(controversialPoints, allPrecedents),
        GeneratedAt:         time.Now(),
    }
    
    return basis, nil
}

// AppealBasis estrutura de embasamento para recurso
type AppealBasis struct {
    CaseNumber          string                    `json:"case_number"`
    AppealType          string                    `json:"appeal_type"`
    ControversialPoints []ControversialPoint      `json:"controversial_points"`
    SupportingPrecedents []JurisprudenceEntry     `json:"supporting_precedents"`
    RecommendedStrategy string                    `json:"recommended_strategy"`
    GeneratedAt         time.Time                 `json:"generated_at"`
}

// ControversialPoint ponto controvertido
type ControversialPoint struct {
    Issue       string  `json:"issue"`
    Context     string  `json:"context"`
    Position    string  `json:"position"`
    Strength    float64 `json:"strength"`
    Precedents  []string `json:"precedents"`
}
```

---

## 4. 🧠 **CONTROLE DE IA - ANTI-DELÍRIO**

### **🛡️ Sistema de Controle de Contexto**

```go
// internal/domain/ai_control.go
package domain

import (
    "context"
    "fmt"
    "regexp"
    "strings"
    "time"
)

// AIControlService serviço de controle de IA
type AIControlService struct {
    allowedTopics    []string
    blockedKeywords  []string
    contextValidator ContextValidator
    outputFilter     OutputFilter
    auditLogger      AuditLogger
}

// AIContext contexto permitido para IA
type AIContext struct {
    Domain          string            `json:"domain"`
    AllowedTopics   []string          `json:"allowed_topics"`
    UserContext     map[string]string `json:"user_context"`
    SystemBoundary  string            `json:"system_boundary"`
    MaxResponseSize int               `json:"max_response_size"`
}

// NewAIControlService cria serviço de controle
func NewAIControlService() *AIControlService {
    return &AIControlService{
        allowedTopics: []string{
            "direito civil",
            "direito penal",
            "direito trabalhista",
            "direito tributário",
            "direito empresarial",
            "direito constitucional",
            "direito administrativo",
            "processo civil",
            "processo penal",
            "jurisprudência",
            "legislação brasileira",
            "consulta processual",
            "direito lux",
        },
        blockedKeywords: []string{
            "médico",
            "diagnóstico",
            "tratamento",
            "receita",
            "medicina",
            "investimento",
            "dinheiro",
            "criptomoeda",
            "trading",
            "religião",
            "política",
            "relacionamento",
            "pessoal",
        },
    }
}

// ValidateAndProcessRequest valida e processa requisição de IA
func (s *AIControlService) ValidateAndProcessRequest(ctx context.Context, request *AIRequest) (*AIResponse, error) {
    // 1. Validar contexto
    if err := s.validateContext(ctx, request); err != nil {
        return nil, err
    }
    
    // 2. Preparar prompt com contexto limitado
    controlledPrompt := s.prepareControlledPrompt(request)
    
    // 3. Processar com IA
    rawResponse, err := s.processWithAI(ctx, controlledPrompt)
    if err != nil {
        return nil, err
    }
    
    // 4. Filtrar e validar resposta
    filteredResponse := s.filterResponse(rawResponse)
    
    // 5. Verificar se resposta está dentro do contexto
    if err := s.validateResponse(filteredResponse); err != nil {
        return nil, err
    }
    
    // 6. Log de auditoria
    s.auditLogger.LogAIInteraction(ctx, request, filteredResponse)
    
    return filteredResponse, nil
}

// validateContext valida se requisição está dentro do contexto permitido
func (s *AIControlService) validateContext(ctx context.Context, request *AIRequest) error {
    // 1. Verificar se contém palavras bloqueadas
    lowerQuery := strings.ToLower(request.Query)
    for _, blocked := range s.blockedKeywords {
        if strings.Contains(lowerQuery, blocked) {
            return fmt.Errorf("tópico não permitido: %s", blocked)
        }
    }
    
    // 2. Verificar se está relacionado a tópicos jurídicos
    isLegalTopic := false
    for _, allowed := range s.allowedTopics {
        if strings.Contains(lowerQuery, allowed) {
            isLegalTopic = true
            break
        }
    }
    
    if !isLegalTopic {
        return fmt.Errorf("consulta deve ser sobre tópicos jurídicos")
    }
    
    // 3. Validar contexto do usuário
    if request.UserContext.Role != "ADVOGADO" && request.UserContext.Role != "CLIENTE" {
        return fmt.Errorf("usuário não autorizado para consultas de IA")
    }
    
    return nil
}

// prepareControlledPrompt prepara prompt com contexto controlado
func (s *AIControlService) prepareControlledPrompt(request *AIRequest) string {
    systemPrompt := `
VOCÊ É LUXIA, ASSISTENTE DE IA JURÍDICA DO DIREITO LUX.

CONTEXTO RESTRITO:
- Você SOMENTE responde sobre assuntos jurídicos brasileiros
- Você NUNCA dá conselhos médicos, financeiros ou pessoais
- Você SEMPRE mantém foco no sistema Direito Lux
- Você NUNCA inventa informações sobre processos ou leis
- Você SEMPRE indica quando não tem certeza

DADOS DO USUÁRIO:
- Role: %s
- Tenant: %s
- Processos disponíveis: %d

INSTRUÇÕES:
1. Responda APENAS sobre direito brasileiro
2. Se a pergunta não for jurídica, redirecione educadamente
3. Se não souber algo, diga "Não tenho essa informação"
4. Mantenha respostas objetivas e profissionais
5. Cite fontes quando possível

PERGUNTA DO USUÁRIO:
%s`

    return fmt.Sprintf(systemPrompt,
        request.UserContext.Role,
        request.UserContext.TenantID,
        len(request.UserContext.AvailableProcesses),
        request.Query,
    )
}

// filterResponse filtra resposta para manter contexto
func (s *AIControlService) filterResponse(response *AIResponse) *AIResponse {
    // 1. Remover conteúdo não jurídico
    filtered := s.removeNonLegalContent(response.Content)
    
    // 2. Adicionar disclaimers se necessário
    if s.needsDisclaimer(filtered) {
        filtered = s.addLegalDisclaimer(filtered)
    }
    
    // 3. Limitar tamanho da resposta
    if len(filtered) > 2000 {
        filtered = filtered[:2000] + "...\n\nResposta limitada por segurança."
    }
    
    return &AIResponse{
        Content:     filtered,
        Confidence:  response.Confidence,
        Sources:     response.Sources,
        GeneratedAt: time.Now(),
    }
}

// validateResponse valida se resposta está dentro do contexto
func (s *AIControlService) validateResponse(response *AIResponse) error {
    content := strings.ToLower(response.Content)
    
    // 1. Verificar conteúdo bloqueado
    for _, blocked := range s.blockedKeywords {
        if strings.Contains(content, blocked) {
            return fmt.Errorf("resposta contém conteúdo bloqueado: %s", blocked)
        }
    }
    
    // 2. Verificar se menciona contexto jurídico
    legalIndicators := []string{
        "lei", "código", "jurisprudência", "tribunal", "processo",
        "direito", "legal", "judicial", "advocacia", "advogado",
    }
    
    hasLegalContent := false
    for _, indicator := range legalIndicators {
        if strings.Contains(content, indicator) {
            hasLegalContent = true
            break
        }
    }
    
    if !hasLegalContent {
        return fmt.Errorf("resposta não contém conteúdo jurídico relevante")
    }
    
    return nil
}

// removeNonLegalContent remove conteúdo não jurídico
func (s *AIControlService) removeNonLegalContent(content string) string {
    // Padrões para remover
    patterns := []string{
        `(?i)(não sou médico|consulte um médico|procure um médico)`,
        `(?i)(investir|investimento|trading|criptomoeda)`,
        `(?i)(relacionamento|pessoal|família|casamento)`,
        `(?i)(política|eleição|candidato|partido)`,
        `(?i)(religião|religioso|igreja|fé)`,
    }
    
    for _, pattern := range patterns {
        re := regexp.MustCompile(pattern)
        content = re.ReplaceAllString(content, "[conteúdo removido - fora do contexto jurídico]")
    }
    
    return content
}

// addLegalDisclaimer adiciona disclaimer legal
func (s *AIControlService) addLegalDisclaimer(content string) string {
    disclaimer := `

⚖️ IMPORTANTE: Esta resposta é apenas informativa e não substitui consulta jurídica profissional. Para questões específicas, consulte um advogado.`
    
    return content + disclaimer
}

// AIRequest requisição para IA
type AIRequest struct {
    Query       string            `json:"query"`
    Tool        string            `json:"tool"`
    UserContext AIUserContext     `json:"user_context"`
    SystemData  map[string]string `json:"system_data"`
}

// AIUserContext contexto do usuário
type AIUserContext struct {
    UserID             string   `json:"user_id"`
    Role               string   `json:"role"`
    TenantID           string   `json:"tenant_id"`
    AvailableProcesses []string `json:"available_processes"`
}

// AIResponse resposta da IA
type AIResponse struct {
    Content     string            `json:"content"`
    Confidence  float64           `json:"confidence"`
    Sources     []string          `json:"sources"`
    GeneratedAt time.Time         `json:"generated_at"`
}
```

---

## 5. 🔄 **POLLING AUTOMÁTICO - DOCUMENTAÇÃO CLARA**

### **📋 Sistema de Polling Documentado**

```go
// internal/domain/polling_system.go
package domain

import (
    "context"
    "fmt"
    "time"
)

// PollingScheduler agendador de polling
type PollingScheduler struct {
    processRepository ProcessRepository
    dataJudService    DataJudService
    notificationService NotificationService
    config            PollingConfig
    isRunning         bool
    stopChan          chan struct{}
}

// PollingConfig configuração do polling
type PollingConfig struct {
    // Intervalos de polling
    ActiveProcessInterval    time.Duration `json:"active_process_interval"`     // 30 minutos
    UrgentProcessInterval    time.Duration `json:"urgent_process_interval"`     // 15 minutos
    ArchivedProcessInterval  time.Duration `json:"archived_process_interval"`   // 24 horas
    
    // Horários de funcionamento
    StartHour    int `json:"start_hour"`     // 6h
    EndHour      int `json:"end_hour"`       // 22h
    
    // Controle de carga
    BatchSize         int           `json:"batch_size"`          // 50 processos por lote
    BatchInterval     time.Duration `json:"batch_interval"`      // 2 segundos entre lotes
    MaxConcurrency    int           `json:"max_concurrency"`     // 5 workers paralelos
    
    // Por plano
    PlansConfig map[string]PlanPollingConfig `json:"plans_config"`
}

// PlanPollingConfig configuração por plano
type PlanPollingConfig struct {
    PlanName         string        `json:"plan_name"`
    MaxProcesses     int           `json:"max_processes"`
    PollingInterval  time.Duration `json:"polling_interval"`
    PriorityBoost    bool          `json:"priority_boost"`
    RealTimeAlerts   bool          `json:"real_time_alerts"`
    CustomSchedule   bool          `json:"custom_schedule"`
}

// GetDefaultPollingConfig retorna configuração padrão
func GetDefaultPollingConfig() PollingConfig {
    return PollingConfig{
        ActiveProcessInterval:   30 * time.Minute,
        UrgentProcessInterval:   15 * time.Minute,
        ArchivedProcessInterval: 24 * time.Hour,
        StartHour:               6,
        EndHour:                 22,
        BatchSize:               50,
        BatchInterval:           2 * time.Second,
        MaxConcurrency:          5,
        PlansConfig: map[string]PlanPollingConfig{
            "starter": {
                PlanName:        "starter",
                MaxProcesses:    50,
                PollingInterval: 60 * time.Minute,   // 1 hora
                PriorityBoost:   false,
                RealTimeAlerts:  false,
                CustomSchedule:  false,
            },
            "professional": {
                PlanName:        "professional",
                MaxProcesses:    200,
                PollingInterval: 30 * time.Minute,   // 30 minutos
                PriorityBoost:   true,
                RealTimeAlerts:  true,
                CustomSchedule:  false,
            },
            "business": {
                PlanName:        "business",
                MaxProcesses:    500,
                PollingInterval: 15 * time.Minute,   // 15 minutos
                PriorityBoost:   true,
                RealTimeAlerts:  true,
                CustomSchedule:  true,
            },
            "enterprise": {
                PlanName:        "enterprise",
                MaxProcesses:    -1,                 // Ilimitado
                PollingInterval: 10 * time.Minute,   // 10 minutos
                PriorityBoost:   true,
                RealTimeAlerts:  true,
                CustomSchedule:  true,
            },
        },
    }
}

// StartPolling inicia sistema de polling
func (s *PollingScheduler) StartPolling(ctx context.Context) error {
    if s.isRunning {
        return fmt.Errorf("polling já está em execução")
    }
    
    s.isRunning = true
    s.stopChan = make(chan struct{})
    
    // Iniciar workers de polling
    for i := 0; i < s.config.MaxConcurrency; i++ {
        go s.pollingWorker(ctx, i)
    }
    
    // Iniciar agendador principal
    go s.schedulerLoop(ctx)
    
    return nil
}

// StopPolling para sistema de polling
func (s *PollingScheduler) StopPolling() {
    if !s.isRunning {
        return
    }
    
    s.isRunning = false
    close(s.stopChan)
}

// schedulerLoop loop principal do agendador
func (s *PollingScheduler) schedulerLoop(ctx context.Context) {
    // Ticker para processos ativos
    activeTicker := time.NewTicker(s.config.ActiveProcessInterval)
    defer activeTicker.Stop()
    
    // Ticker para processos urgentes
    urgentTicker := time.NewTicker(s.config.UrgentProcessInterval)
    defer urgentTicker.Stop()
    
    // Ticker para processos arquivados
    archivedTicker := time.NewTicker(s.config.ArchivedProcessInterval)
    defer archivedTicker.Stop()
    
    for {
        select {
        case <-activeTicker.C:
            if s.isWorkingHours() {
                s.scheduleActiveProcesses(ctx)
            }
            
        case <-urgentTicker.C:
            if s.isWorkingHours() {
                s.scheduleUrgentProcesses(ctx)
            }
            
        case <-archivedTicker.C:
            s.scheduleArchivedProcesses(ctx)
            
        case <-s.stopChan:
            return
            
        case <-ctx.Done():
            return
        }
    }
}

// scheduleActiveProcesses agenda processos ativos
func (s *PollingScheduler) scheduleActiveProcesses(ctx context.Context) {
    // Buscar processos ativos por plano
    for planName, planConfig := range s.config.PlansConfig {
        processes, err := s.processRepository.GetActiveProcessesByPlan(ctx, planName)
        if err != nil {
            continue
        }
        
        // Limitar por quota do plano
        if planConfig.MaxProcesses > 0 && len(processes) > planConfig.MaxProcesses {
            processes = processes[:planConfig.MaxProcesses]
        }
        
        // Processar em lotes
        s.processBatches(ctx, processes, planConfig)
    }
}

// scheduleUrgentProcesses agenda processos urgentes
func (s *PollingScheduler) scheduleUrgentProcesses(ctx context.Context) {
    urgentProcesses, err := s.processRepository.GetUrgentProcesses(ctx)
    if err != nil {
        return
    }
    
    // Processos urgentes têm prioridade máxima
    s.processBatches(ctx, urgentProcesses, PlanPollingConfig{
        PriorityBoost:  true,
        RealTimeAlerts: true,
    })
}

// processBatches processa processos em lotes
func (s *PollingScheduler) processBatches(ctx context.Context, processes []Process, config PlanPollingConfig) {
    // Dividir em lotes
    for i := 0; i < len(processes); i += s.config.BatchSize {
        end := i + s.config.BatchSize
        if end > len(processes) {
            end = len(processes)
        }
        
        batch := processes[i:end]
        
        // Enviar lote para processamento
        s.processBatch(ctx, batch, config)
        
        // Intervalo entre lotes
        time.Sleep(s.config.BatchInterval)
    }
}

// processBatch processa um lote de processos
func (s *PollingScheduler) processBatch(ctx context.Context, batch []Process, config PlanPollingConfig) {
    for _, process := range batch {
        // Verificar se deve processar baseado no plano
        if !s.shouldProcessNow(process, config) {
            continue
        }
        
        // Consultar atualizações
        s.checkProcessUpdates(ctx, process, config)
    }
}

// shouldProcessNow verifica se deve processar agora
func (s *PollingScheduler) shouldProcessNow(process Process, config PlanPollingConfig) bool {
    // Verificar última consulta
    timeSinceLastCheck := time.Since(process.LastChecked)
    
    // Respeitar intervalo do plano
    if timeSinceLastCheck < config.PollingInterval {
        return false
    }
    
    // Prioridade boost para planos premium
    if config.PriorityBoost && process.Priority == "HIGH" {
        return true
    }
    
    // Verificar se está dentro do horário permitido
    if !s.isWorkingHours() && !config.RealTimeAlerts {
        return false
    }
    
    return true
}

// checkProcessUpdates verifica atualizações de processo
func (s *PollingScheduler) checkProcessUpdates(ctx context.Context, process Process, config PlanPollingConfig) {
    // Consultar DataJud
    updates, err := s.dataJudService.GetProcessUpdates(ctx, process.Number)
    if err != nil {
        return
    }
    
    // Verificar se há mudanças
    if s.hasChanges(process, updates) {
        // Atualizar processo
        s.processRepository.UpdateProcess(ctx, process.ID, updates)
        
        // Notificar se necessário
        if config.RealTimeAlerts || s.isImportantChange(updates) {
            s.notificationService.NotifyProcessUpdate(ctx, process, updates)
        }
    }
    
    // Atualizar timestamp da última verificação
    s.processRepository.UpdateLastChecked(ctx, process.ID, time.Now())
}

// isWorkingHours verifica se está no horário de trabalho
func (s *PollingScheduler) isWorkingHours() bool {
    now := time.Now()
    hour := now.Hour()
    return hour >= s.config.StartHour && hour <= s.config.EndHour
}

// GetPollingStatus retorna status do polling
func (s *PollingScheduler) GetPollingStatus(ctx context.Context) map[string]interface{} {
    status := map[string]interface{}{
        "is_running":     s.isRunning,
        "working_hours":  s.isWorkingHours(),
        "config":         s.config,
        "next_runs":      s.getNextRuns(),
        "statistics":     s.getStatistics(ctx),
    }
    
    return status
}

// getNextRuns retorna próximas execuções
func (s *PollingScheduler) getNextRuns() map[string]time.Time {
    now := time.Now()
    return map[string]time.Time{
        "active_processes":   now.Add(s.config.ActiveProcessInterval),
        "urgent_processes":   now.Add(s.config.UrgentProcessInterval),
        "archived_processes": now.Add(s.config.ArchivedProcessInterval),
    }
}

// getStatistics retorna estatísticas do polling
func (s *PollingScheduler) getStatistics(ctx context.Context) map[string]interface{} {
    stats := map[string]interface{}{
        "total_processes":    s.processRepository.CountProcesses(ctx),
        "active_processes":   s.processRepository.CountActiveProcesses(ctx),
        "urgent_processes":   s.processRepository.CountUrgentProcesses(ctx),
        "archived_processes": s.processRepository.CountArchivedProcesses(ctx),
        "last_24h_updates":   s.processRepository.CountRecentUpdates(ctx, 24*time.Hour),
    }
    
    return stats
}
```

### **🎯 Funcionalidades por Plano**

```go
// internal/domain/plan_features.go
package domain

// PlanFeatures funcionalidades por plano
type PlanFeatures struct {
    PlanName              string        `json:"plan_name"`
    MaxProcesses          int           `json:"max_processes"`
    PollingInterval       time.Duration `json:"polling_interval"`
    RealTimeNotifications bool          `json:"real_time_notifications"`
    CustomSchedule        bool          `json:"custom_schedule"`
    PrioritySupport       bool          `json:"priority_support"`
    AdvancedReports       bool          `json:"advanced_reports"`
    AIAssistant           bool          `json:"ai_assistant"`
    JurisprudenceSearch   bool          `json:"jurisprudence_search"`
    BulkOperations        bool          `json:"bulk_operations"`
    WhatsAppIntegration   bool          `json:"whatsapp_integration"`
    TelegramIntegration   bool          `json:"telegram_integration"`
    EmailNotifications    bool          `json:"email_notifications"`
    MobileApp             bool          `json:"mobile_app"`
    APIAccess             bool          `json:"api_access"`
    DataExport            bool          `json:"data_export"`
    CustomDashboard       bool          `json:"custom_dashboard"`
    MultiUser             bool          `json:"multi_user"`
    AuditLog              bool          `json:"audit_log"`
    LGPD                  bool          `json:"lgpd_compliance"`
}

// GetPlanFeatures retorna funcionalidades por plano
func GetPlanFeatures() map[string]PlanFeatures {
    return map[string]PlanFeatures{
        "starter": {
            PlanName:              "starter",
            MaxProcesses:          50,
            PollingInterval:       60 * time.Minute,
            RealTimeNotifications: false,
            CustomSchedule:        false,
            PrioritySupport:       false,
            AdvancedReports:       false,
            AIAssistant:           false,
            JurisprudenceSearch:   false,
            BulkOperations:        false,
            WhatsAppIntegration:   true,
            TelegramIntegration:   false,
            EmailNotifications:    true,
            MobileApp:             false,
            APIAccess:             false,
            DataExport:            false,
            CustomDashboard:       false,
            MultiUser:             false,
            AuditLog:              false,
            LGPD:                  true,
        },
        "professional": {
            PlanName:              "professional",
            MaxProcesses:          200,
            PollingInterval:       30 * time.Minute,
            RealTimeNotifications: true,
            CustomSchedule:        false,
            PrioritySupport:       true,
            AdvancedReports:       true,
            AIAssistant:           true,
            JurisprudenceSearch:   true,
            BulkOperations:        true,
            WhatsAppIntegration:   true,
            TelegramIntegration:   true,
            EmailNotifications:    true,
            MobileApp:             true,
            APIAccess:             false,
            DataExport:            true,
            CustomDashboard:       true,
            MultiUser:             true,
            AuditLog:              true,
            LGPD:                  true,
        },
        "business": {
            PlanName:              "business",
            MaxProcesses:          500,
            PollingInterval:       15 * time.Minute,
            RealTimeNotifications: true,
            CustomSchedule:        true,
            PrioritySupport:       true,
            AdvancedReports:       true,
            AIAssistant:           true,
            JurisprudenceSearch:   true,
            BulkOperations:        true,
            WhatsAppIntegration:   true,
            TelegramIntegration:   true,
            EmailNotifications:    true,
            MobileApp:             true,
            APIAccess:             true,
            DataExport:            true,
            CustomDashboard:       true,
            MultiUser:             true,
            AuditLog:              true,
            LGPD:                  true,
        },
        "enterprise": {
            PlanName:              "enterprise",
            MaxProcesses:          -1, // Ilimitado
            PollingInterval:       10 * time.Minute,
            RealTimeNotifications: true,
            CustomSchedule:        true,
            PrioritySupport:       true,
            AdvancedReports:       true,
            AIAssistant:           true,
            JurisprudenceSearch:   true,
            BulkOperations:        true,
            WhatsAppIntegration:   true,
            TelegramIntegration:   true,
            EmailNotifications:    true,
            MobileApp:             true,
            APIAccess:             true,
            DataExport:            true,
            CustomDashboard:       true,
            MultiUser:             true,
            AuditLog:              true,
            LGPD:                  true,
        },
    }
}
```

---

## 6. ✅ **FUNCIONALIDADES SAAS COMPLETAS**

### **📋 Checklist de Funcionalidades SaaS**

```yaml
AUTENTICAÇÃO & AUTORIZAÇÃO:
✅ JWT Authentication
✅ Multi-tenant isolation
✅ Role-based access control (RBAC)
✅ Session management
✅ Password reset
✅ Two-factor authentication (2FA)
✅ Single sign-on (SSO)

GERENCIAMENTO DE USUÁRIOS:
✅ User registration
✅ User profiles
✅ User management (CRUD)
✅ User permissions
✅ User activity tracking
✅ User onboarding

PLANOS & BILLING:
✅ Subscription plans
✅ Payment processing (ASAAS)
✅ Invoice generation
✅ Usage tracking
✅ Quota management
✅ Billing history
✅ Payment notifications
✅ Trial periods
✅ Plan upgrades/downgrades

NOTIFICAÇÕES:
✅ Email notifications
✅ WhatsApp notifications
✅ Telegram notifications
✅ Push notifications
✅ SMS notifications
✅ In-app notifications
✅ Notification preferences
✅ Bulk notifications

RELATÓRIOS & ANALYTICS:
✅ Usage analytics
✅ Performance metrics
✅ Business reports
✅ Custom dashboards
✅ Data export
✅ Audit logs
✅ Real-time metrics

INTEGRAÇÕES:
✅ REST APIs
✅ Webhooks
✅ Third-party integrations
✅ API documentation
✅ SDK availability
✅ Rate limiting
✅ API versioning

SEGURANÇA:
✅ Data encryption
✅ LGPD compliance
✅ Security audit logs
✅ Input validation
✅ XSS protection
✅ CSRF protection
✅ SQL injection protection

INFRAESTRUTURA:
✅ High availability
✅ Scalability
✅ Backup & recovery
✅ Monitoring & alerting
✅ Performance optimization
✅ CDN integration
✅ Load balancing

SUPORTE:
✅ Help documentation
✅ Knowledge base
✅ Support ticket system
✅ Live chat
✅ Email support
✅ Phone support (Enterprise)
✅ Training materials

CUSTOMIZAÇÃO:
✅ White-label options
✅ Custom branding
✅ Custom domains
✅ Custom fields
✅ Custom workflows
✅ Custom reports
```

---

## 7. 🔐 **VALIDAÇÃO DE ACESSO LUXIA (SEM LOGIN ENGESSADO)**

### **🚀 Sistema de Validação Fluida**

```go
// internal/domain/chat_auth.go
package domain

import (
    "context"
    "fmt"
    "time"
)

// ChatAuthService serviço de autenticação no chat
type ChatAuthService struct {
    userRepository    UserRepository
    tokenService      TokenService
    sessionManager    SessionManager
    notificationService NotificationService
}

// AuthMethod método de autenticação
type AuthMethod string

const (
    AuthMethodPhone     AuthMethod = "phone"     // Número de telefone
    AuthMethodEmail     AuthMethod = "email"     // Email + código
    AuthMethodToken     AuthMethod = "token"     // Token temporário
    AuthMethodBiometric AuthMethod = "biometric" // Biometria (futuro)
)

// ChatSession sessão de chat
type ChatSession struct {
    ID           string            `json:"id"`
    UserID       string            `json:"user_id"`
    ChatID       string            `json:"chat_id"`
    Platform     string            `json:"platform"` // whatsapp, telegram
    PhoneNumber  string            `json:"phone_number"`
    IsAuthenticated bool           `json:"is_authenticated"`
    Role         string            `json:"role"`
    TenantID     string            `json:"tenant_id"`
    Permissions  []string          `json:"permissions"`
    ExpiresAt    time.Time         `json:"expires_at"`
    CreatedAt    time.Time         `json:"created_at"`
    LastActivity time.Time         `json:"last_activity"`
    AuthMethod   AuthMethod        `json:"auth_method"`
    Metadata     map[string]string `json:"metadata"`
}

// AuthenticateUser autentica usuário no chat
func (s *ChatAuthService) AuthenticateUser(ctx context.Context, chatID, phoneNumber, platform string) (*ChatSession, error) {
    // 1. Verificar se já existe sessão ativa
    if session := s.sessionManager.GetActiveSession(ctx, chatID); session != nil {
        // Renovar sessão se ainda válida
        if time.Now().Before(session.ExpiresAt) {
            session.LastActivity = time.Now()
            s.sessionManager.UpdateSession(ctx, session)
            return session, nil
        }
    }
    
    // 2. Buscar usuário por telefone
    user, err := s.userRepository.GetUserByPhone(ctx, phoneNumber)
    if err != nil {
        // Usuário não encontrado - iniciar fluxo de registro
        return s.initiateRegistration(ctx, chatID, phoneNumber, platform)
    }
    
    // 3. Verificar se telefone está verificado
    if !user.PhoneVerified {
        return s.initiatePhoneVerification(ctx, chatID, phoneNumber, platform)
    }
    
    // 4. Criar sessão autenticada
    session := &ChatSession{
        ID:              s.generateSessionID(),
        UserID:          user.ID,
        ChatID:          chatID,
        Platform:        platform,
        PhoneNumber:     phoneNumber,
        IsAuthenticated: true,
        Role:            user.Role,
        TenantID:        user.TenantID,
        Permissions:     user.Permissions,
        ExpiresAt:       time.Now().Add(24 * time.Hour), // 24 horas
        CreatedAt:       time.Now(),
        LastActivity:    time.Now(),
        AuthMethod:      AuthMethodPhone,
    }
    
    // 5. Salvar sessão
    if err := s.sessionManager.CreateSession(ctx, session); err != nil {
        return nil, err
    }
    
    // 6. Notificar usuário
    s.sendWelcomeMessage(ctx, session)
    
    return session, nil
}

// initiateRegistration inicia processo de registro
func (s *ChatAuthService) initiateRegistration(ctx context.Context, chatID, phoneNumber, platform string) (*ChatSession, error) {
    // Criar sessão temporária
    session := &ChatSession{
        ID:              s.generateSessionID(),
        ChatID:          chatID,
        Platform:        platform,
        PhoneNumber:     phoneNumber,
        IsAuthenticated: false,
        ExpiresAt:       time.Now().Add(30 * time.Minute), // 30 minutos
        CreatedAt:       time.Now(),
        LastActivity:    time.Now(),
        AuthMethod:      AuthMethodPhone,
        Metadata: map[string]string{
            "registration_step": "phone_verification",
        },
    }
    
    // Enviar código de verificação
    code := s.generateVerificationCode()
    s.sendVerificationCode(ctx, phoneNumber, code)
    
    // Salvar código temporariamente
    s.sessionManager.StoreVerificationCode(ctx, phoneNumber, code, 10*time.Minute)
    
    // Enviar mensagem de boas-vindas
    s.sendRegistrationMessage(ctx, session)
    
    return session, nil
}

// sendWelcomeMessage envia mensagem de boas-vindas
func (s *ChatAuthService) sendWelcomeMessage(ctx context.Context, session *ChatSession) {
    var message string
    
    switch session.Role {
    case "ADVOGADO":
        message = fmt.Sprintf(`
🎯 Olá! Sou a *Luxia*, sua assistente de IA jurídica.

✅ Você está autenticado como *Advogado*
📱 Acesso total ao sistema Direito Lux

🔧 *Comandos disponíveis:*
• "Meus processos" - Ver processos ativos
• "Relatório mensal" - Gerar relatório
• "Buscar jurisprudência" - Pesquisar casos
• "Clientes" - Gerenciar clientes
• "Notificações" - Configurar alertas

💡 *Dica:* Fale naturalmente comigo! Posso ajudar com qualquer dúvida jurídica.

Como posso ajudar você hoje?`)
        
    case "CLIENTE":
        message = fmt.Sprintf(`
👋 Olá! Sou a *Luxia*, sua assistente jurídica.

✅ Você está autenticado como *Cliente*
📋 Posso ajudar com informações dos seus processos

🔧 *O que posso fazer:*
• Mostrar status dos seus processos
• Explicar andamentos processuais
• Agendar reuniões
• Responder dúvidas jurídicas
• Notificar sobre atualizações

💡 *Dica:* Pergunte sobre qualquer processo seu ou dúvida jurídica!

Como posso ajudar você hoje?`)
        
    case "FUNCIONARIO":
        message = fmt.Sprintf(`
👨‍💼 Olá! Sou a *Luxia*, sua assistente de trabalho.

✅ Você está autenticado como *Funcionário*
📊 Acesso aos processos atribuídos a você

🔧 *Comandos disponíveis:*
• "Meus processos" - Ver processos atribuídos
• "Prazos" - Verificar prazos vencendo
• "Relatório semanal" - Gerar relatório
• "Buscar processo" - Pesquisar processos
• "Ajuda" - Lista de comandos

Como posso ajudar você hoje?`)
    }
    
    s.notificationService.SendChatMessage(ctx, session.ChatID, session.Platform, message)
}

// sendRegistrationMessage envia mensagem de registro
func (s *ChatAuthService) sendRegistrationMessage(ctx context.Context, session *ChatSession) {
    message := fmt.Sprintf(`
🎉 *Bem-vindo ao Direito Lux!*

📱 Para começar, preciso verificar seu número de telefone.

📨 Enviei um código de verificação para *%s*
⏰ O código é válido por 10 minutos

📝 *Responda com o código recebido para continuar*

Exemplo: 123456

❓ Não recebeu o código? Digite "reenviar"`, session.PhoneNumber)
    
    s.notificationService.SendChatMessage(ctx, session.ChatID, session.Platform, message)
}

// ValidateAccess valida acesso a funcionalidade
func (s *ChatAuthService) ValidateAccess(ctx context.Context, chatID, functionality string) error {
    // 1. Obter sessão
    session := s.sessionManager.GetActiveSession(ctx, chatID)
    if session == nil {
        return fmt.Errorf("usuário não autenticado")
    }
    
    // 2. Verificar se sessão expirou
    if time.Now().After(session.ExpiresAt) {
        return fmt.Errorf("sessão expirou")
    }
    
    // 3. Verificar permissões
    if !s.hasPermission(session, functionality) {
        return fmt.Errorf("usuário não tem permissão para esta funcionalidade")
    }
    
    // 4. Atualizar última atividade
    session.LastActivity = time.Now()
    s.sessionManager.UpdateSession(ctx, session)
    
    return nil
}

// hasPermission verifica se usuário tem permissão
func (s *ChatAuthService) hasPermission(session *ChatSession, functionality string) bool {
    // Definir permissões por funcionalidade
    functionalityPermissions := map[string][]string{
        "process_search": {"read_processes"},
        "process_create": {"create_processes"},
        "jurisprudence_search": {"read_jurisprudence"},
        "generate_report": {"generate_reports"},
        "client_management": {"manage_clients"},
        "user_management": {"manage_users"},
        "billing": {"view_billing"},
        "admin": {"admin_access"},
    }
    
    requiredPermissions, exists := functionalityPermissions[functionality]
    if !exists {
        return false // Funcionalidade não definida
    }
    
    // Verificar se usuário tem pelo menos uma permissão necessária
    for _, required := range requiredPermissions {
        for _, userPerm := range session.Permissions {
            if userPerm == required {
                return true
            }
        }
    }
    
    return false
}

// RefreshSession renova sessão
func (s *ChatAuthService) RefreshSession(ctx context.Context, chatID string) (*ChatSession, error) {
    session := s.sessionManager.GetActiveSession(ctx, chatID)
    if session == nil {
        return nil, fmt.Errorf("sessão não encontrada")
    }
    
    // Renovar por mais 24 horas
    session.ExpiresAt = time.Now().Add(24 * time.Hour)
    session.LastActivity = time.Now()
    
    if err := s.sessionManager.UpdateSession(ctx, session); err != nil {
        return nil, err
    }
    
    return session, nil
}

// GetSessionInfo retorna informações da sessão
func (s *ChatAuthService) GetSessionInfo(ctx context.Context, chatID string) map[string]interface{} {
    session := s.sessionManager.GetActiveSession(ctx, chatID)
    if session == nil {
        return map[string]interface{}{
            "authenticated": false,
            "message": "Usuário não autenticado",
        }
    }
    
    return map[string]interface{}{
        "authenticated": session.IsAuthenticated,
        "user_id":      session.UserID,
        "role":         session.Role,
        "tenant_id":    session.TenantID,
        "permissions":  session.Permissions,
        "expires_at":   session.ExpiresAt,
        "time_left":    time.Until(session.ExpiresAt).String(),
        "auth_method":  session.AuthMethod,
    }
}
```

### **💬 Fluxo de Autenticação Prático**

```yaml
1. PRIMEIRA MENSAGEM:
   👤 Cliente: "Oi Luxia"
   🤖 Luxia: "Olá! Para começar, preciso verificar seu número.
            Enviei um código para +55 11 99999-9999"

2. VERIFICAÇÃO:
   👤 Cliente: "123456"
   🤖 Luxia: "✅ Verificado! Bem-vindo João Silva!
            Você é cliente do escritório Costa Advogados.
            Como posso ajudar?"

3. USO NORMAL:
   👤 Cliente: "Como está meu processo?"
   🤖 Luxia: "Seu processo 1001234-56.2024.8.26.0100:
            Status: Aguardando decisão
            Última atualização: 2 dias atrás"

4. SESSÃO EXPIRA:
   👤 Cliente: "Meus processos"
   🤖 Luxia: "Sua sessão expirou. Enviando novo código..."
```

---

## ✅ **RESUMO FINAL**

### **📚 Documentação Criada:**
- ✅ **Controle de acesso MCP** - Documentado completamente
- ✅ **Gerenciamento múltiplos CNPJs** - Sistema avançado implementado
- ✅ **Search jurisprudência** - Busca semântica para embasamento
- ✅ **Controle anti-delírio IA** - Sistema de contexto restrito
- ✅ **Polling automático** - Documentado com funcionalidades por plano
- ✅ **Funcionalidades SaaS** - Checklist completo
- ✅ **Validação Luxia** - Sistema fluido sem login engessado

### **🎯 Todos os Pontos Respondidos:**
1. **✅ MCP documentado** - Controle de acesso por role
2. **✅ DataJud multi-CNPJ** - Sistema de pool inteligente
3. **✅ Search jurisprudência** - Busca semântica + embasamento
4. **✅ IA controlada** - Sistema anti-delírio jurídico
5. **✅ Polling documentado** - Funcionalidades por plano
6. **✅ SaaS completo** - Todas as funcionalidades cobertas
7. **✅ Luxia sem login** - Autenticação fluida por telefone

**🚀 Arquitetura 100% documentada e pronta para desenvolvimento!**