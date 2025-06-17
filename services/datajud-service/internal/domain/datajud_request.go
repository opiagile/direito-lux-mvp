package domain

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RequestType define os tipos de requisição da API DataJud
type RequestType string

const (
	RequestTypeProcess   RequestType = "process"    // Consulta de processo
	RequestTypeMovement  RequestType = "movement"   // Consulta de movimentações
	RequestTypeParty     RequestType = "party"      // Consulta de partes
	RequestTypeDocument  RequestType = "document"   // Consulta de documentos
	RequestTypeBulk      RequestType = "bulk"       // Consulta em lote
)

// RequestStatus define o status da requisição
type RequestStatus string

const (
	StatusPending    RequestStatus = "pending"     // Aguardando processamento
	StatusProcessing RequestStatus = "processing"  // Em processamento
	StatusCompleted  RequestStatus = "completed"   // Concluída com sucesso
	StatusFailed     RequestStatus = "failed"      // Falhou
	StatusCached     RequestStatus = "cached"      // Retornada do cache
	StatusRetrying   RequestStatus = "retrying"    // Tentando novamente
)

// RequestPriority define a prioridade da requisição
type RequestPriority int

const (
	PriorityLow    RequestPriority = 1 // Baixa prioridade
	PriorityNormal RequestPriority = 2 // Prioridade normal
	PriorityHigh   RequestPriority = 3 // Alta prioridade
	PriorityUrgent RequestPriority = 4 // Urgente
)

// DataJudRequest representa uma requisição para a API DataJud
type DataJudRequest struct {
	ID              uuid.UUID       `json:"id"`
	TenantID        uuid.UUID       `json:"tenant_id"`
	ClientID        uuid.UUID       `json:"client_id"`
	ProcessID       *uuid.UUID      `json:"process_id,omitempty"`
	Type            RequestType     `json:"type"`
	Priority        RequestPriority `json:"priority"`
	Status          RequestStatus   `json:"status"`
	CNPJProviderID  *uuid.UUID      `json:"cnpj_provider_id,omitempty"`
	
	// Parâmetros da requisição
	ProcessNumber   string                 `json:"process_number,omitempty"`
	CourtID         string                 `json:"court_id,omitempty"`
	Parameters      map[string]interface{} `json:"parameters"`
	
	// Controle de cache
	CacheKey        string    `json:"cache_key"`
	CacheTTL        int       `json:"cache_ttl"` // TTL em segundos
	UseCache        bool      `json:"use_cache"`
	
	// Controle de retry
	RetryCount      int       `json:"retry_count"`
	MaxRetries      int       `json:"max_retries"`
	RetryAfter      *time.Time `json:"retry_after,omitempty"`
	
	// Controle de circuit breaker
	CircuitBreakerKey string `json:"circuit_breaker_key"`
	
	// Timestamps
	RequestedAt     time.Time  `json:"requested_at"`
	ProcessingAt    *time.Time `json:"processing_at,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	
	// Resultado
	Response        *DataJudResponse `json:"response,omitempty"`
	ErrorMessage    string           `json:"error_message,omitempty"`
	ErrorCode       string           `json:"error_code,omitempty"`
}

// DataJudResponse representa a resposta da API DataJud
type DataJudResponse struct {
	ID          uuid.UUID   `json:"id"`
	RequestID   uuid.UUID   `json:"request_id"`
	StatusCode  int         `json:"status_code"`
	Headers     map[string]string `json:"headers"`
	Body        json.RawMessage   `json:"body"`
	Size        int64       `json:"size"`
	Duration    int64       `json:"duration_ms"`
	FromCache   bool        `json:"from_cache"`
	ReceivedAt  time.Time   `json:"received_at"`
	
	// Dados estruturados (opcional)
	ProcessData   *ProcessResponseData   `json:"process_data,omitempty"`
	MovementData  *MovementResponseData  `json:"movement_data,omitempty"`
	PartyData     *PartyResponseData     `json:"party_data,omitempty"`
}

// ProcessResponseData dados estruturados de processo
type ProcessResponseData struct {
	Number      string                 `json:"number"`
	Title       string                 `json:"title"`
	Subject     map[string]interface{} `json:"subject"`
	Court       string                 `json:"court"`
	Status      string                 `json:"status"`
	Stage       string                 `json:"stage"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	Parties     []PartyData           `json:"parties"`
	Movements   []MovementData        `json:"movements"`
}

// MovementResponseData dados estruturados de movimentação
type MovementResponseData struct {
	Movements []MovementData `json:"movements"`
	Total     int           `json:"total"`
	Page      int           `json:"page"`
	PageSize  int           `json:"page_size"`
}

// PartyResponseData dados estruturados de partes
type PartyResponseData struct {
	Parties []PartyData `json:"parties"`
	Total   int        `json:"total"`
}

// MovementData representa uma movimentação
type MovementData struct {
	Sequence    int                    `json:"sequence"`
	Date        time.Time              `json:"date"`
	Type        string                 `json:"type"`
	Code        string                 `json:"code"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Content     string                 `json:"content"`
	IsPublic    bool                   `json:"is_public"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// PartyData representa uma parte do processo
type PartyData struct {
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Document   string                 `json:"document"`
	Role       string                 `json:"role"`
	Contact    map[string]interface{} `json:"contact"`
	Address    map[string]interface{} `json:"address"`
	Lawyer     map[string]interface{} `json:"lawyer"`
}

// DataJudRequestRepository interface para persistência
type DataJudRequestRepository interface {
	Save(request *DataJudRequest) error
	FindByID(id uuid.UUID) (*DataJudRequest, error)
	FindByTenantID(tenantID uuid.UUID, limit, offset int) ([]*DataJudRequest, error)
	FindPendingRequests(limit int) ([]*DataJudRequest, error)
	FindByStatus(status RequestStatus, limit, offset int) ([]*DataJudRequest, error)
	FindByProcessNumber(processNumber string) ([]*DataJudRequest, error)
	Update(request *DataJudRequest) error
	UpdateStatus(id uuid.UUID, status RequestStatus) error
	Delete(id uuid.UUID) error
	CleanupOldRequests(olderThan time.Time) (int, error)
}

// NewDataJudRequest cria uma nova requisição DataJud
func NewDataJudRequest(tenantID, clientID uuid.UUID, requestType RequestType, priority RequestPriority) *DataJudRequest {
	now := time.Now()
	
	return &DataJudRequest{
		ID:           uuid.New(),
		TenantID:     tenantID,
		ClientID:     clientID,
		Type:         requestType,
		Priority:     priority,
		Status:       StatusPending,
		Parameters:   make(map[string]interface{}),
		UseCache:     true,
		CacheTTL:     3600, // 1 hora por padrão
		MaxRetries:   3,
		RetryCount:   0,
		RequestedAt:  now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// SetProcessNumber define o número do processo
func (r *DataJudRequest) SetProcessNumber(processNumber string) {
	r.ProcessNumber = processNumber
	r.Parameters["process_number"] = processNumber
	r.generateCacheKey()
	r.UpdatedAt = time.Now()
}

// SetCourtID define o tribunal
func (r *DataJudRequest) SetCourtID(courtID string) {
	r.CourtID = courtID
	r.Parameters["court_id"] = courtID
	r.generateCacheKey()
	r.UpdatedAt = time.Now()
}

// SetParameter define um parâmetro customizado
func (r *DataJudRequest) SetParameter(key string, value interface{}) {
	r.Parameters[key] = value
	r.generateCacheKey()
	r.UpdatedAt = time.Now()
}

// SetCNPJProvider define o CNPJ provider a ser usado
func (r *DataJudRequest) SetCNPJProvider(providerID uuid.UUID) {
	r.CNPJProviderID = &providerID
	r.UpdatedAt = time.Now()
}

// StartProcessing marca a requisição como em processamento
func (r *DataJudRequest) StartProcessing() {
	r.Status = StatusProcessing
	now := time.Now()
	r.ProcessingAt = &now
	r.UpdatedAt = now
}

// Complete marca a requisição como concluída
func (r *DataJudRequest) Complete(response *DataJudResponse) {
	r.Status = StatusCompleted
	r.Response = response
	now := time.Now()
	r.CompletedAt = &now
	r.UpdatedAt = now
}

// Fail marca a requisição como falhada
func (r *DataJudRequest) Fail(errorCode, errorMessage string) {
	r.Status = StatusFailed
	r.ErrorCode = errorCode
	r.ErrorMessage = errorMessage
	r.UpdatedAt = time.Now()
}

// Retry prepara a requisição para nova tentativa
func (r *DataJudRequest) Retry(retryAfter time.Duration) error {
	if r.RetryCount >= r.MaxRetries {
		return fmt.Errorf("máximo de tentativas (%d) excedido", r.MaxRetries)
	}

	r.Status = StatusRetrying
	r.RetryCount++
	retryTime := time.Now().Add(retryAfter)
	r.RetryAfter = &retryTime
	r.UpdatedAt = time.Now()

	return nil
}

// CanRetry verifica se a requisição pode ser reexecutada
func (r *DataJudRequest) CanRetry() bool {
	if r.RetryCount >= r.MaxRetries {
		return false
	}

	if r.RetryAfter != nil && time.Now().Before(*r.RetryAfter) {
		return false
	}

	return r.Status == StatusFailed || r.Status == StatusRetrying
}

// IsExpired verifica se a requisição expirou
func (r *DataJudRequest) IsExpired(maxAge time.Duration) bool {
	return time.Since(r.CreatedAt) > maxAge
}

// GetPriorityWeight retorna o peso da prioridade para ordenação
func (r *DataJudRequest) GetPriorityWeight() int {
	return int(r.Priority)
}

// generateCacheKey gera uma chave de cache baseada nos parâmetros
func (r *DataJudRequest) generateCacheKey() {
	data := map[string]interface{}{
		"type":           r.Type,
		"process_number": r.ProcessNumber,
		"court_id":       r.CourtID,
		"parameters":     r.Parameters,
	}

	jsonData, _ := json.Marshal(data)
	hash := md5.Sum(jsonData)
	r.CacheKey = fmt.Sprintf("datajud:%s:%x", r.Type, hash)
}

// SetCircuitBreakerKey define a chave do circuit breaker
func (r *DataJudRequest) SetCircuitBreakerKey(key string) {
	r.CircuitBreakerKey = key
	r.UpdatedAt = time.Now()
}

// GetEstimatedDuration retorna duração estimada baseada no tipo
func (r *DataJudRequest) GetEstimatedDuration() time.Duration {
	switch r.Type {
	case RequestTypeProcess:
		return 5 * time.Second
	case RequestTypeMovement:
		return 3 * time.Second
	case RequestTypeParty:
		return 2 * time.Second
	case RequestTypeDocument:
		return 10 * time.Second
	case RequestTypeBulk:
		return 30 * time.Second
	default:
		return 5 * time.Second
	}
}