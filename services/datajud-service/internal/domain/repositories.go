package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// EventStore interface para persistência de eventos
type EventStore interface {
	SaveEvent(ctx context.Context, event DomainEvent) error
	GetEvents(ctx context.Context, aggregateID uuid.UUID) ([]DomainEvent, error)
	GetEventsByType(ctx context.Context, eventType string) ([]DomainEvent, error)
	GetEventsAfter(ctx context.Context, timestamp time.Time) ([]DomainEvent, error)
}

// UnitOfWork interface para transações
type UnitOfWork interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	IsInTransaction() bool
}

// Repositories agregador de todos os repositórios
type Repositories struct {
	CNPJProvider     CNPJProviderRepository
	CNPJPool         CNPJPoolRepository
	DataJudRequest   DataJudRequestRepository
	Cache            CacheRepository
	RateLimiter      RateLimiterRepository
	CircuitBreaker   CircuitBreakerRepository
	EventStore       EventStore
	UnitOfWork       UnitOfWork
}

// DomainService interface para serviços de domínio
type DomainService interface {
	ValidateCNPJ(cnpj string) error
	ValidateProcessNumber(processNumber string) error
	CalculateRequestPriority(requestType RequestType, urgent bool) RequestPriority
	EstimateRequestDuration(requestType RequestType) time.Duration
	ShouldUseCache(requestType RequestType, age time.Duration) bool
}

// DataJudDomainService implementação do serviço de domínio
type DataJudDomainService struct{}

// NewDataJudDomainService cria novo serviço de domínio
func NewDataJudDomainService() *DataJudDomainService {
	return &DataJudDomainService{}
}

// ValidateCNPJ valida formato e dígitos do CNPJ
func (s *DataJudDomainService) ValidateCNPJ(cnpj string) error {
	if !isValidCNPJ(cnpj) {
		return NewValidationError("CNPJ", "formato inválido")
	}
	return nil
}

// ValidateProcessNumber valida número de processo CNJ
func (s *DataJudDomainService) ValidateProcessNumber(processNumber string) error {
	// Reutiliza a função do process-service (seria importada)
	// Por simplicidade, validação básica aqui
	if len(processNumber) < 20 {
		return NewValidationError("ProcessNumber", "número de processo inválido")
	}
	return nil
}

// CalculateRequestPriority calcula prioridade baseada no tipo e urgência
func (s *DataJudDomainService) CalculateRequestPriority(requestType RequestType, urgent bool) RequestPriority {
	if urgent {
		return PriorityUrgent
	}

	switch requestType {
	case RequestTypeBulk:
		return PriorityLow
	case RequestTypeDocument:
		return PriorityNormal
	case RequestTypeMovement, RequestTypeParty:
		return PriorityNormal
	case RequestTypeProcess:
		return PriorityHigh
	default:
		return PriorityNormal
	}
}

// EstimateRequestDuration estima duração baseada no tipo
func (s *DataJudDomainService) EstimateRequestDuration(requestType RequestType) time.Duration {
	switch requestType {
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

// ShouldUseCache determina se deve usar cache baseado no tipo e idade
func (s *DataJudDomainService) ShouldUseCache(requestType RequestType, age time.Duration) bool {
	maxAge := map[RequestType]time.Duration{
		RequestTypeProcess:  1 * time.Hour,    // Processos mudam pouco
		RequestTypeMovement: 30 * time.Minute, // Movimentações são mais dinâmicas
		RequestTypeParty:    2 * time.Hour,    // Partes raramente mudam
		RequestTypeDocument: 24 * time.Hour,   // Documentos são estáticos
		RequestTypeBulk:     15 * time.Minute, // Bulk queries são temporais
	}

	return age <= maxAge[requestType]
}

// ValidationError erro de validação de domínio
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// NewValidationError cria novo erro de validação
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// BusinessError erro de regra de negócio
type BusinessError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e BusinessError) Error() string {
	return e.Code + ": " + e.Message
}

// NewBusinessError cria novo erro de negócio
func NewBusinessError(code, message string) BusinessError {
	return BusinessError{
		Code:    code,
		Message: message,
	}
}

// Common business errors
var (
	ErrCNPJNotFound      = NewBusinessError("CNPJ_NOT_FOUND", "CNPJ provider não encontrado")
	ErrQuotaExhausted    = NewBusinessError("QUOTA_EXHAUSTED", "Quota de requisições esgotada")
	ErrCircuitBreakerOpen = NewBusinessError("CIRCUIT_BREAKER_OPEN", "Circuit breaker está aberto")
	ErrRateLimitExceeded = NewBusinessError("RATE_LIMIT_EXCEEDED", "Rate limit excedido")
	ErrInvalidRequest    = NewBusinessError("INVALID_REQUEST", "Requisição inválida")
	ErrServiceUnavailable = NewBusinessError("SERVICE_UNAVAILABLE", "Serviço temporariamente indisponível")
)

// DataJudConfig configurações específicas do DataJud
type DataJudConfig struct {
	// API DataJud
	APIBaseURL        string        `json:"api_base_url"`
	APITimeout        time.Duration `json:"api_timeout"`
	APIRetryCount     int           `json:"api_retry_count"`
	APIRetryDelay     time.Duration `json:"api_retry_delay"`
	
	// Rate Limiting
	DefaultDailyLimit int           `json:"default_daily_limit"`
	GlobalRateLimit   int           `json:"global_rate_limit"`
	RateWindowSize    time.Duration `json:"rate_window_size"`
	
	// Cache
	DefaultCacheTTL   int           `json:"default_cache_ttl"`
	MaxCacheSize      int64         `json:"max_cache_size"`
	MaxCacheEntries   int64         `json:"max_cache_entries"`
	CacheCleanupInterval time.Duration `json:"cache_cleanup_interval"`
	
	// Circuit Breaker
	CBFailureThreshold int           `json:"cb_failure_threshold"`
	CBSuccessThreshold int           `json:"cb_success_threshold"`
	CBTimeout          time.Duration `json:"cb_timeout"`
	CBMaxRequests      int           `json:"cb_max_requests"`
	
	// Pool Strategy
	DefaultPoolStrategy CNPJPoolStrategy `json:"default_pool_strategy"`
	
	// Monitoring
	MetricsEnabled    bool          `json:"metrics_enabled"`
	HealthCheckInterval time.Duration `json:"health_check_interval"`
}

// DefaultDataJudConfig retorna configuração padrão
func DefaultDataJudConfig() DataJudConfig {
	return DataJudConfig{
		// API DataJud
		APIBaseURL:    "https://api-publica.datajud.cnj.jus.br",
		APITimeout:    30 * time.Second,
		APIRetryCount: 3,
		APIRetryDelay: 5 * time.Second,
		
		// Rate Limiting
		DefaultDailyLimit: 10000,
		GlobalRateLimit:   100,
		RateWindowSize:    time.Hour,
		
		// Cache
		DefaultCacheTTL:      3600, // 1 hora
		MaxCacheSize:         1024 * 1024 * 1024, // 1GB
		MaxCacheEntries:      100000,
		CacheCleanupInterval: 15 * time.Minute,
		
		// Circuit Breaker
		CBFailureThreshold: 5,
		CBSuccessThreshold: 3,
		CBTimeout:          30 * time.Second,
		CBMaxRequests:      5,
		
		// Pool Strategy
		DefaultPoolStrategy: StrategyLeastUsed,
		
		// Monitoring
		MetricsEnabled:      true,
		HealthCheckInterval: 30 * time.Second,
	}
}

// ServiceHealth representa o status de saúde do serviço
type ServiceHealth struct {
	Status        string                 `json:"status"`        // healthy, degraded, unhealthy
	Timestamp     time.Time              `json:"timestamp"`
	Components    map[string]ComponentHealth `json:"components"`
	ResponseTime  time.Duration          `json:"response_time"`
	Version       string                 `json:"version"`
	Environment   string                 `json:"environment"`
}

// ComponentHealth status de um componente
type ComponentHealth struct {
	Status      string                 `json:"status"`
	Message     string                 `json:"message,omitempty"`
	Metrics     map[string]interface{} `json:"metrics,omitempty"`
	LastChecked time.Time              `json:"last_checked"`
}

// HealthChecker interface para health checks
type HealthChecker interface {
	CheckHealth() ServiceHealth
	CheckComponent(name string) ComponentHealth
	IsHealthy() bool
}

// MetricsCollector interface para coleta de métricas
type MetricsCollector interface {
	IncrementCounter(name string, labels map[string]string)
	RecordHistogram(name string, value float64, labels map[string]string)
	SetGauge(name string, value float64, labels map[string]string)
	GetMetrics() map[string]interface{}
}