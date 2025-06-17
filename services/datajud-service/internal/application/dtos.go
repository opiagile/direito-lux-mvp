package application

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/direito-lux/datajud-service/internal/domain"
)

// ========================================
// REQUEST DTOS
// ========================================

// ProcessQueryRequest requisição para consulta de processo
type ProcessQueryRequest struct {
	TenantID      uuid.UUID  `json:"tenant_id" validate:"required"`
	ClientID      uuid.UUID  `json:"client_id" validate:"required"`
	ProcessID     *uuid.UUID `json:"process_id,omitempty"`
	ProcessNumber string     `json:"process_number" validate:"required"`
	CourtID       string     `json:"court_id" validate:"required"`
	UseCache      bool       `json:"use_cache"`
	Urgent        bool       `json:"urgent"`
}

// Validate valida a requisição de processo
func (r *ProcessQueryRequest) Validate() error {
	if r.TenantID == uuid.Nil {
		return errors.New("tenant_id é obrigatório")
	}
	if r.ClientID == uuid.Nil {
		return errors.New("client_id é obrigatório")
	}
	if r.ProcessNumber == "" {
		return errors.New("process_number é obrigatório")
	}
	if r.CourtID == "" {
		return errors.New("court_id é obrigatório")
	}
	return nil
}

// MovementQueryRequest requisição para consulta de movimentações
type MovementQueryRequest struct {
	TenantID      uuid.UUID  `json:"tenant_id" validate:"required"`
	ClientID      uuid.UUID  `json:"client_id" validate:"required"`
	ProcessID     *uuid.UUID `json:"process_id,omitempty"`
	ProcessNumber string     `json:"process_number" validate:"required"`
	CourtID       string     `json:"court_id" validate:"required"`
	DateFrom      *time.Time `json:"date_from,omitempty"`
	DateTo        *time.Time `json:"date_to,omitempty"`
	Page          int        `json:"page"`
	PageSize      int        `json:"page_size"`
	UseCache      bool       `json:"use_cache"`
	Urgent        bool       `json:"urgent"`
}

// Validate valida a requisição de movimentações
func (r *MovementQueryRequest) Validate() error {
	if r.TenantID == uuid.Nil {
		return errors.New("tenant_id é obrigatório")
	}
	if r.ClientID == uuid.Nil {
		return errors.New("client_id é obrigatório")
	}
	if r.ProcessNumber == "" {
		return errors.New("process_number é obrigatório")
	}
	if r.CourtID == "" {
		return errors.New("court_id é obrigatório")
	}
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 || r.PageSize > 100 {
		r.PageSize = 20
	}
	return nil
}

// PartyQueryRequest requisição para consulta de partes
type PartyQueryRequest struct {
	TenantID      uuid.UUID `json:"tenant_id" validate:"required"`
	ClientID      uuid.UUID `json:"client_id" validate:"required"`
	ProcessNumber string    `json:"process_number" validate:"required"`
	CourtID       string    `json:"court_id" validate:"required"`
	UseCache      bool      `json:"use_cache"`
	Urgent        bool      `json:"urgent"`
}

// Validate valida a requisição de partes
func (r *PartyQueryRequest) Validate() error {
	if r.TenantID == uuid.Nil {
		return errors.New("tenant_id é obrigatório")
	}
	if r.ClientID == uuid.Nil {
		return errors.New("client_id é obrigatório")
	}
	if r.ProcessNumber == "" {
		return errors.New("process_number é obrigatório")
	}
	if r.CourtID == "" {
		return errors.New("court_id é obrigatório")
	}
	return nil
}

// BulkQueryRequest requisição para consultas em lote
type BulkQueryRequest struct {
	TenantID uuid.UUID         `json:"tenant_id" validate:"required"`
	ClientID uuid.UUID         `json:"client_id" validate:"required"`
	Queries  []BulkQueryItem   `json:"queries" validate:"required,min=1,max=100"`
}

// BulkQueryItem item individual de consulta em lote
type BulkQueryItem struct {
	ProcessNumber string `json:"process_number" validate:"required"`
	CourtID       string `json:"court_id" validate:"required"`
}

// Validate valida a requisição de lote
func (r *BulkQueryRequest) Validate() error {
	if r.TenantID == uuid.Nil {
		return errors.New("tenant_id é obrigatório")
	}
	if r.ClientID == uuid.Nil {
		return errors.New("client_id é obrigatório")
	}
	if len(r.Queries) == 0 {
		return errors.New("pelo menos uma consulta é obrigatória")
	}
	if len(r.Queries) > 100 {
		return errors.New("máximo de 100 consultas por lote")
	}
	
	for i, query := range r.Queries {
		if query.ProcessNumber == "" {
			return errors.New("process_number é obrigatório na consulta " + string(rune(i+1)))
		}
		if query.CourtID == "" {
			return errors.New("court_id é obrigatório na consulta " + string(rune(i+1)))
		}
	}
	
	return nil
}

// ========================================
// RESPONSE DTOS
// ========================================

// ProcessQueryResponse resposta de consulta de processo
type ProcessQueryResponse struct {
	RequestID uuid.UUID                      `json:"request_id"`
	Status    string                         `json:"status"`
	Data      *domain.ProcessResponseData    `json:"data,omitempty"`
	Error     string                         `json:"error,omitempty"`
	FromCache bool                           `json:"from_cache"`
	CachedAt  *time.Time                     `json:"cached_at,omitempty"`
	Duration  int64                          `json:"duration_ms,omitempty"`
}

// MovementQueryResponse resposta de consulta de movimentações
type MovementQueryResponse struct {
	RequestID uuid.UUID                       `json:"request_id"`
	Status    string                          `json:"status"`
	Data      *domain.MovementResponseData    `json:"data,omitempty"`
	Error     string                          `json:"error,omitempty"`
	FromCache bool                            `json:"from_cache"`
	CachedAt  *time.Time                      `json:"cached_at,omitempty"`
	Duration  int64                           `json:"duration_ms,omitempty"`
}

// PartyQueryResponse resposta de consulta de partes
type PartyQueryResponse struct {
	RequestID uuid.UUID                   `json:"request_id"`
	Status    string                      `json:"status"`
	Data      *domain.PartyResponseData   `json:"data,omitempty"`
	Error     string                      `json:"error,omitempty"`
	FromCache bool                        `json:"from_cache"`
	CachedAt  *time.Time                  `json:"cached_at,omitempty"`
	Duration  int64                       `json:"duration_ms,omitempty"`
}

// BulkQueryResponse resposta de consultas em lote
type BulkQueryResponse struct {
	RequestID   uuid.UUID         `json:"request_id"`
	Status      string            `json:"status"`
	Results     []BulkQueryResult `json:"results"`
	StartedAt   time.Time         `json:"started_at"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
	Duration    time.Duration     `json:"duration"`
}

// BulkQueryResult resultado individual de consulta em lote
type BulkQueryResult struct {
	Index         int                         `json:"index"`
	ProcessNumber string                      `json:"process_number"`
	Status        string                      `json:"status"`
	Data          *domain.ProcessResponseData `json:"data,omitempty"`
	Error         string                      `json:"error,omitempty"`
	Duration      int64                       `json:"duration_ms,omitempty"`
}

// ========================================
// CNPJ PROVIDER MANAGEMENT DTOS
// ========================================

// CreateCNPJProviderRequest requisição para criar provider CNPJ
type CreateCNPJProviderRequest struct {
	TenantID        uuid.UUID `json:"tenant_id" validate:"required"`
	CNPJ            string    `json:"cnpj" validate:"required"`
	Name            string    `json:"name" validate:"required"`
	Email           string    `json:"email" validate:"required,email"`
	APIKey          string    `json:"api_key" validate:"required"`
	Certificate     string    `json:"certificate,omitempty"`
	CertificatePass string    `json:"certificate_pass,omitempty"`
	DailyLimit      int       `json:"daily_limit,omitempty"`
	Priority        int       `json:"priority,omitempty"`
}

// Validate valida a requisição de criação de provider
func (r *CreateCNPJProviderRequest) Validate() error {
	if r.TenantID == uuid.Nil {
		return errors.New("tenant_id é obrigatório")
	}
	if r.CNPJ == "" {
		return errors.New("cnpj é obrigatório")
	}
	if r.Name == "" {
		return errors.New("name é obrigatório")
	}
	if r.Email == "" {
		return errors.New("email é obrigatório")
	}
	if r.APIKey == "" {
		return errors.New("api_key é obrigatório")
	}
	if r.DailyLimit <= 0 {
		r.DailyLimit = 10000 // Padrão DataJud
	}
	if r.Priority <= 0 {
		r.Priority = 1
	}
	return nil
}

// UpdateCNPJProviderRequest requisição para atualizar provider
type UpdateCNPJProviderRequest struct {
	Name            *string `json:"name,omitempty"`
	Email           *string `json:"email,omitempty"`
	APIKey          *string `json:"api_key,omitempty"`
	Certificate     *string `json:"certificate,omitempty"`
	CertificatePass *string `json:"certificate_pass,omitempty"`
	DailyLimit      *int    `json:"daily_limit,omitempty"`
	Priority        *int    `json:"priority,omitempty"`
	IsActive        *bool   `json:"is_active,omitempty"`
}

// CNPJProviderResponse resposta com dados do provider
type CNPJProviderResponse struct {
	ID                uuid.UUID  `json:"id"`
	TenantID          uuid.UUID  `json:"tenant_id"`
	CNPJ              string     `json:"cnpj"`
	Name              string     `json:"name"`
	Email             string     `json:"email"`
	DailyLimit        int        `json:"daily_limit"`
	DailyUsage        int        `json:"daily_usage"`
	AvailableQuota    int        `json:"available_quota"`
	UsagePercentage   float64    `json:"usage_percentage"`
	UsageResetTime    time.Time  `json:"usage_reset_time"`
	IsActive          bool       `json:"is_active"`
	Priority          int        `json:"priority"`
	LastUsedAt        *time.Time `json:"last_used_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// CreateCNPJProviderResponse resposta de criação de provider
type CreateCNPJProviderResponse struct {
	Provider *CNPJProviderResponse `json:"provider"`
	Message  string                `json:"message"`
}

// ========================================
// POOL MANAGEMENT DTOS
// ========================================

// CreateCNPJPoolRequest requisição para criar pool
type CreateCNPJPoolRequest struct {
	TenantID     uuid.UUID               `json:"tenant_id" validate:"required"`
	Name         string                  `json:"name" validate:"required"`
	Strategy     domain.CNPJPoolStrategy `json:"strategy" validate:"required"`
	ProviderIDs  []uuid.UUID             `json:"provider_ids,omitempty"`
}

// Validate valida a requisição de criação de pool
func (r *CreateCNPJPoolRequest) Validate() error {
	if r.TenantID == uuid.Nil {
		return errors.New("tenant_id é obrigatório")
	}
	if r.Name == "" {
		return errors.New("name é obrigatório")
	}
	validStrategies := []domain.CNPJPoolStrategy{
		domain.StrategyRoundRobin,
		domain.StrategyLeastUsed,
		domain.StrategyPriority,
		domain.StrategyAvailability,
	}
	valid := false
	for _, strategy := range validStrategies {
		if r.Strategy == strategy {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("strategy inválida")
	}
	return nil
}

// CNPJPoolResponse resposta com dados do pool
type CNPJPoolResponse struct {
	ID        uuid.UUID               `json:"id"`
	TenantID  uuid.UUID               `json:"tenant_id"`
	Name      string                  `json:"name"`
	Strategy  domain.CNPJPoolStrategy `json:"strategy"`
	IsActive  bool                    `json:"is_active"`
	Stats     domain.PoolStats        `json:"stats"`
	Providers []CNPJProviderResponse  `json:"providers"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
}

// ========================================
// STATISTICS AND MONITORING DTOS
// ========================================

// DataJudStatsRequest requisição para estatísticas
type DataJudStatsRequest struct {
	TenantID  *uuid.UUID `json:"tenant_id,omitempty"`
	DateFrom  *time.Time `json:"date_from,omitempty"`
	DateTo    *time.Time `json:"date_to,omitempty"`
	GroupBy   string     `json:"group_by,omitempty"` // hour, day, week, month
}

// DataJudStatsResponse resposta com estatísticas
type DataJudStatsResponse struct {
	TotalRequests    int64                  `json:"total_requests"`
	SuccessfulRequests int64                `json:"successful_requests"`
	FailedRequests   int64                  `json:"failed_requests"`
	CacheHits        int64                  `json:"cache_hits"`
	CacheMisses      int64                  `json:"cache_misses"`
	CacheHitRatio    float64                `json:"cache_hit_ratio"`
	AvgResponseTime  float64                `json:"avg_response_time_ms"`
	TotalQuotaUsed   int64                  `json:"total_quota_used"`
	ActiveCNPJs      int                    `json:"active_cnpjs"`
	CircuitBreakers  map[string]string      `json:"circuit_breakers"`
	TopProcessNumbers []ProcessNumberStat   `json:"top_process_numbers"`
	TopCourts        []CourtStat            `json:"top_courts"`
	HourlyStats      []HourlyStat           `json:"hourly_stats,omitempty"`
	DailyStats       []DailyStat            `json:"daily_stats,omitempty"`
}

// ProcessNumberStat estatística por número de processo
type ProcessNumberStat struct {
	ProcessNumber string `json:"process_number"`
	RequestCount  int64  `json:"request_count"`
	LastRequested time.Time `json:"last_requested"`
}

// CourtStat estatística por tribunal
type CourtStat struct {
	CourtID      string `json:"court_id"`
	RequestCount int64  `json:"request_count"`
	SuccessRate  float64 `json:"success_rate"`
}

// HourlyStat estatística por hora
type HourlyStat struct {
	Hour     time.Time `json:"hour"`
	Requests int64     `json:"requests"`
	Success  int64     `json:"success"`
	Failed   int64     `json:"failed"`
}

// DailyStat estatística por dia
type DailyStat struct {
	Date     time.Time `json:"date"`
	Requests int64     `json:"requests"`
	Success  int64     `json:"success"`
	Failed   int64     `json:"failed"`
	QuotaUsed int64    `json:"quota_used"`
}

// ========================================
// HEALTH CHECK DTO
// ========================================

// HealthCheckResponse resposta de health check
type HealthCheckResponse struct {
	Status        string                     `json:"status"`
	Timestamp     time.Time                  `json:"timestamp"`
	Version       string                     `json:"version"`
	Environment   string                     `json:"environment"`
	Components    map[string]ComponentHealth `json:"components"`
	ResponseTime  time.Duration              `json:"response_time"`
}

// ComponentHealth status de um componente
type ComponentHealth struct {
	Status      string                 `json:"status"`
	Message     string                 `json:"message,omitempty"`
	Metrics     map[string]interface{} `json:"metrics,omitempty"`
	LastChecked time.Time              `json:"last_checked"`
}