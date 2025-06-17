package domain

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent interface base para eventos de domínio
type DomainEvent interface {
	GetID() uuid.UUID
	GetType() string
	GetAggregateID() uuid.UUID
	GetOccurredAt() time.Time
	GetVersion() int
	GetMetadata() map[string]interface{}
}

// BaseEvent estrutura base para eventos
type BaseEvent struct {
	ID          uuid.UUID              `json:"id"`
	Type        string                 `json:"type"`
	AggregateID uuid.UUID              `json:"aggregate_id"`
	OccurredAt  time.Time              `json:"occurred_at"`
	Version     int                    `json:"version"`
	Metadata    map[string]interface{} `json:"metadata"`
}

func (e BaseEvent) GetID() uuid.UUID                  { return e.ID }
func (e BaseEvent) GetType() string                   { return e.Type }
func (e BaseEvent) GetAggregateID() uuid.UUID         { return e.AggregateID }
func (e BaseEvent) GetOccurredAt() time.Time          { return e.OccurredAt }
func (e BaseEvent) GetVersion() int                   { return e.Version }
func (e BaseEvent) GetMetadata() map[string]interface{} { return e.Metadata }

// ========================================
// EVENTOS DE CNPJ PROVIDER
// ========================================

// CNPJProviderCreated evento de criação de provider CNPJ
type CNPJProviderCreated struct {
	BaseEvent
	TenantID   uuid.UUID `json:"tenant_id"`
	CNPJ       string    `json:"cnpj"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	DailyLimit int       `json:"daily_limit"`
	Priority   int       `json:"priority"`
}

// CNPJProviderActivated evento de ativação de provider
type CNPJProviderActivated struct {
	BaseEvent
	CNPJ string `json:"cnpj"`
}

// CNPJProviderDeactivated evento de desativação de provider
type CNPJProviderDeactivated struct {
	BaseEvent
	CNPJ   string `json:"cnpj"`
	Reason string `json:"reason"`
}

// CNPJQuotaExhausted evento de quota esgotada
type CNPJQuotaExhausted struct {
	BaseEvent
	CNPJ        string    `json:"cnpj"`
	DailyUsage  int       `json:"daily_usage"`
	DailyLimit  int       `json:"daily_limit"`
	ResetTime   time.Time `json:"reset_time"`
}

// CNPJQuotaReset evento de reset de quota diária
type CNPJQuotaReset struct {
	BaseEvent
	CNPJ          string `json:"cnpj"`
	PreviousUsage int    `json:"previous_usage"`
	NewLimit      int    `json:"new_limit"`
}

// ========================================
// EVENTOS DE POOL
// ========================================

// CNPJPoolCreated evento de criação de pool
type CNPJPoolCreated struct {
	BaseEvent
	TenantID     uuid.UUID        `json:"tenant_id"`
	PoolName     string           `json:"pool_name"`
	Strategy     CNPJPoolStrategy `json:"strategy"`
	ProviderIDs  []uuid.UUID      `json:"provider_ids"`
}

// CNPJProviderAddedToPool evento de adição ao pool
type CNPJProviderAddedToPool struct {
	BaseEvent
	PoolID     uuid.UUID `json:"pool_id"`
	ProviderID uuid.UUID `json:"provider_id"`
	CNPJ       string    `json:"cnpj"`
}

// CNPJProviderRemovedFromPool evento de remoção do pool
type CNPJProviderRemovedFromPool struct {
	BaseEvent
	PoolID     uuid.UUID `json:"pool_id"`
	ProviderID uuid.UUID `json:"provider_id"`
	CNPJ       string    `json:"cnpj"`
	Reason     string    `json:"reason"`
}

// CNPJPoolStrategyChanged evento de mudança de estratégia
type CNPJPoolStrategyChanged struct {
	BaseEvent
	PoolID      uuid.UUID        `json:"pool_id"`
	OldStrategy CNPJPoolStrategy `json:"old_strategy"`
	NewStrategy CNPJPoolStrategy `json:"new_strategy"`
}

// ========================================
// EVENTOS DE REQUISIÇÃO DATAJUD
// ========================================

// DataJudRequestCreated evento de criação de requisição
type DataJudRequestCreated struct {
	BaseEvent
	TenantID      uuid.UUID       `json:"tenant_id"`
	ClientID      uuid.UUID       `json:"client_id"`
	ProcessID     *uuid.UUID      `json:"process_id,omitempty"`
	RequestType   RequestType     `json:"request_type"`
	Priority      RequestPriority `json:"priority"`
	ProcessNumber string          `json:"process_number,omitempty"`
	CourtID       string          `json:"court_id,omitempty"`
	UseCache      bool            `json:"use_cache"`
}

// DataJudRequestStarted evento de início de processamento
type DataJudRequestStarted struct {
	BaseEvent
	CNPJProviderID uuid.UUID `json:"cnpj_provider_id"`
	CNPJ           string    `json:"cnpj"`
	ProcessingAt   time.Time `json:"processing_at"`
}

// DataJudRequestCompleted evento de conclusão com sucesso
type DataJudRequestCompleted struct {
	BaseEvent
	CNPJProviderID uuid.UUID     `json:"cnpj_provider_id"`
	StatusCode     int           `json:"status_code"`
	ResponseSize   int64         `json:"response_size"`
	Duration       time.Duration `json:"duration"`
	FromCache      bool          `json:"from_cache"`
}

// DataJudRequestFailed evento de falha na requisição
type DataJudRequestFailed struct {
	BaseEvent
	CNPJProviderID *uuid.UUID `json:"cnpj_provider_id,omitempty"`
	ErrorCode      string     `json:"error_code"`
	ErrorMessage   string     `json:"error_message"`
	RetryCount     int        `json:"retry_count"`
	WillRetry      bool       `json:"will_retry"`
}

// DataJudRequestRetrying evento de tentativa de retry
type DataJudRequestRetrying struct {
	BaseEvent
	RetryCount int           `json:"retry_count"`
	MaxRetries int           `json:"max_retries"`
	RetryAfter time.Duration `json:"retry_after"`
}

// ========================================
// EVENTOS DE CACHE
// ========================================

// DataJudCacheHit evento de hit no cache
type DataJudCacheHit struct {
	BaseEvent
	CacheKey    string      `json:"cache_key"`
	RequestType RequestType `json:"request_type"`
	HitCount    int64       `json:"hit_count"`
	Age         time.Duration `json:"age"`
}

// DataJudCacheMiss evento de miss no cache
type DataJudCacheMiss struct {
	BaseEvent
	CacheKey    string      `json:"cache_key"`
	RequestType RequestType `json:"request_type"`
}

// DataJudCacheStored evento de armazenamento no cache
type DataJudCacheStored struct {
	BaseEvent
	CacheKey    string      `json:"cache_key"`
	RequestType RequestType `json:"request_type"`
	TTL         int         `json:"ttl"`
	Size        int64       `json:"size"`
}

// DataJudCacheEvicted evento de remoção do cache
type DataJudCacheEvicted struct {
	BaseEvent
	CacheKey    string      `json:"cache_key"`
	RequestType RequestType `json:"request_type"`
	Reason      string      `json:"reason"`
	Age         time.Duration `json:"age"`
}

// ========================================
// EVENTOS DE RATE LIMITING
// ========================================

// RateLimitExceeded evento de limite de taxa excedido
type RateLimitExceeded struct {
	BaseEvent
	LimitType     RateLimitType `json:"limit_type"`
	Key           string        `json:"key"`
	RequestsUsed  int           `json:"requests_used"`
	RequestsLimit int           `json:"requests_limit"`
	ResetTime     time.Time     `json:"reset_time"`
}

// RateLimitReset evento de reset de rate limit
type RateLimitReset struct {
	BaseEvent
	LimitType       RateLimitType `json:"limit_type"`
	Key             string        `json:"key"`
	PreviousUsage   int           `json:"previous_usage"`
	NewLimit        int           `json:"new_limit"`
}

// ========================================
// EVENTOS DE CIRCUIT BREAKER
// ========================================

// CircuitBreakerOpened evento de abertura do circuit breaker
type CircuitBreakerOpened struct {
	BaseEvent
	BreakerName    string `json:"breaker_name"`
	FailureCount   int    `json:"failure_count"`
	FailureThreshold int  `json:"failure_threshold"`
}

// CircuitBreakerClosed evento de fechamento do circuit breaker
type CircuitBreakerClosed struct {
	BaseEvent
	BreakerName    string `json:"breaker_name"`
	SuccessCount   int    `json:"success_count"`
	SuccessThreshold int  `json:"success_threshold"`
}

// CircuitBreakerHalfOpened evento de half-open
type CircuitBreakerHalfOpened struct {
	BaseEvent
	BreakerName string        `json:"breaker_name"`
	Timeout     time.Duration `json:"timeout"`
}

// ========================================
// EVENTOS DE SISTEMA
// ========================================

// DataJudServiceStarted evento de início do serviço
type DataJudServiceStarted struct {
	BaseEvent
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

// DataJudServiceStopped evento de parada do serviço
type DataJudServiceStopped struct {
	BaseEvent
	Reason string `json:"reason"`
}

// DataJudAPIConnectionLost evento de perda de conexão com API
type DataJudAPIConnectionLost struct {
	BaseEvent
	Endpoint    string `json:"endpoint"`
	ErrorCode   string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

// DataJudAPIConnectionRestored evento de restauração de conexão
type DataJudAPIConnectionRestored struct {
	BaseEvent
	Endpoint string        `json:"endpoint"`
	Downtime time.Duration `json:"downtime"`
}

// ========================================
// FACTORY FUNCTIONS
// ========================================

// NewCNPJProviderCreated cria evento de criação de provider
func NewCNPJProviderCreated(providerID, tenantID uuid.UUID, cnpj, name, email string, dailyLimit, priority int) *CNPJProviderCreated {
	return &CNPJProviderCreated{
		BaseEvent: BaseEvent{
			ID:          uuid.New(),
			Type:        "datajud.cnpj_provider.created",
			AggregateID: providerID,
			OccurredAt:  time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		TenantID:   tenantID,
		CNPJ:       cnpj,
		Name:       name,
		Email:      email,
		DailyLimit: dailyLimit,
		Priority:   priority,
	}
}

// NewDataJudRequestCreated cria evento de criação de requisição
func NewDataJudRequestCreated(requestID, tenantID, clientID uuid.UUID, processID *uuid.UUID, requestType RequestType, priority RequestPriority, processNumber, courtID string, useCache bool) *DataJudRequestCreated {
	return &DataJudRequestCreated{
		BaseEvent: BaseEvent{
			ID:          uuid.New(),
			Type:        "datajud.request.created",
			AggregateID: requestID,
			OccurredAt:  time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		TenantID:      tenantID,
		ClientID:      clientID,
		ProcessID:     processID,
		RequestType:   requestType,
		Priority:      priority,
		ProcessNumber: processNumber,
		CourtID:       courtID,
		UseCache:      useCache,
	}
}

// NewDataJudRequestCompleted cria evento de conclusão de requisição
func NewDataJudRequestCompleted(requestID, cnpjProviderID uuid.UUID, statusCode int, responseSize int64, duration time.Duration, fromCache bool) *DataJudRequestCompleted {
	return &DataJudRequestCompleted{
		BaseEvent: BaseEvent{
			ID:          uuid.New(),
			Type:        "datajud.request.completed",
			AggregateID: requestID,
			OccurredAt:  time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		CNPJProviderID: cnpjProviderID,
		StatusCode:     statusCode,
		ResponseSize:   responseSize,
		Duration:       duration,
		FromCache:      fromCache,
	}
}

// NewCircuitBreakerOpened cria evento de abertura do circuit breaker
func NewCircuitBreakerOpened(breakerID uuid.UUID, breakerName string, failureCount, failureThreshold int) *CircuitBreakerOpened {
	return &CircuitBreakerOpened{
		BaseEvent: BaseEvent{
			ID:          uuid.New(),
			Type:        "datajud.circuit_breaker.opened",
			AggregateID: breakerID,
			OccurredAt:  time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		BreakerName:      breakerName,
		FailureCount:     failureCount,
		FailureThreshold: failureThreshold,
	}
}