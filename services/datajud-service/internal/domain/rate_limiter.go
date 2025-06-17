package domain

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// RateLimitType define os tipos de rate limiting
type RateLimitType string

const (
	RateLimitCNPJ   RateLimitType = "cnpj"    // Limite por CNPJ
	RateLimitTenant RateLimitType = "tenant"  // Limite por tenant
	RateLimitGlobal RateLimitType = "global"  // Limite global
)

// RateLimitWindow janela de tempo para rate limiting
type RateLimitWindow struct {
	WindowSize time.Duration `json:"window_size"`
	MaxRequests int          `json:"max_requests"`
}

// RateLimiter controla a taxa de requisições
type RateLimiter struct {
	ID          uuid.UUID       `json:"id"`
	Type        RateLimitType   `json:"type"`
	Key         string          `json:"key"`          // Chave identificadora (CNPJ, tenant_id, etc)
	Window      RateLimitWindow `json:"window"`
	Requests    []time.Time     `json:"requests"`     // Timestamps das requisições
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	LastUsedAt  *time.Time      `json:"last_used_at"`
	mu          sync.RWMutex    // Mutex para thread safety
}

// RateLimitStatus status atual do rate limiter
type RateLimitStatus struct {
	Allowed        bool          `json:"allowed"`
	RequestsUsed   int           `json:"requests_used"`
	RequestsLimit  int           `json:"requests_limit"`
	ResetTime      time.Time     `json:"reset_time"`
	RetryAfter     time.Duration `json:"retry_after"`
	WindowSize     time.Duration `json:"window_size"`
}

// RateLimiterRepository interface para persistência
type RateLimiterRepository interface {
	Save(limiter *RateLimiter) error
	FindByKey(limitType RateLimitType, key string) (*RateLimiter, error)
	FindByCNPJ(cnpj string) (*RateLimiter, error)
	FindByTenant(tenantID uuid.UUID) (*RateLimiter, error)
	Update(limiter *RateLimiter) error
	CleanupExpired() (int, error)
	GetStats(limitType RateLimitType) (map[string]interface{}, error)
}

// NewRateLimiter cria um novo rate limiter
func NewRateLimiter(limitType RateLimitType, key string, window RateLimitWindow) *RateLimiter {
	return &RateLimiter{
		ID:        uuid.New(),
		Type:      limitType,
		Key:       key,
		Window:    window,
		Requests:  make([]time.Time, 0),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewCNPJRateLimiter cria um rate limiter para CNPJ (10k/dia)
func NewCNPJRateLimiter(cnpj string) *RateLimiter {
	return NewRateLimiter(
		RateLimitCNPJ,
		cnpj,
		RateLimitWindow{
			WindowSize:  24 * time.Hour,
			MaxRequests: 10000,
		},
	)
}

// NewTenantRateLimiter cria um rate limiter para tenant
func NewTenantRateLimiter(tenantID uuid.UUID, maxRequests int, windowSize time.Duration) *RateLimiter {
	return NewRateLimiter(
		RateLimitTenant,
		tenantID.String(),
		RateLimitWindow{
			WindowSize:  windowSize,
			MaxRequests: maxRequests,
		},
	)
}

// NewGlobalRateLimiter cria um rate limiter global
func NewGlobalRateLimiter(maxRequests int, windowSize time.Duration) *RateLimiter {
	return NewRateLimiter(
		RateLimitGlobal,
		"global",
		RateLimitWindow{
			WindowSize:  windowSize,
			MaxRequests: maxRequests,
		},
	)
}

// Allow verifica se uma requisição pode ser feita
func (rl *RateLimiter) Allow() RateLimitStatus {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if !rl.IsActive {
		return RateLimitStatus{
			Allowed:       false,
			RequestsUsed:  0,
			RequestsLimit: rl.Window.MaxRequests,
			ResetTime:     time.Now().Add(rl.Window.WindowSize),
			RetryAfter:    rl.Window.WindowSize,
			WindowSize:    rl.Window.WindowSize,
		}
	}

	now := time.Now()
	windowStart := now.Add(-rl.Window.WindowSize)

	// Remove requisições antigas (fora da janela)
	rl.cleanupOldRequests(windowStart)

	requestsInWindow := len(rl.Requests)
	allowed := requestsInWindow < rl.Window.MaxRequests

	if allowed {
		// Adiciona a requisição atual
		rl.Requests = append(rl.Requests, now)
		rl.LastUsedAt = &now
		rl.UpdatedAt = now
		requestsInWindow++
	}

	// Calcula quando o limite será resetado
	var resetTime time.Time
	var retryAfter time.Duration

	if len(rl.Requests) > 0 {
		oldestRequest := rl.Requests[0]
		resetTime = oldestRequest.Add(rl.Window.WindowSize)
		if resetTime.Before(now) {
			resetTime = now.Add(rl.Window.WindowSize)
		}
		retryAfter = resetTime.Sub(now)
	} else {
		resetTime = now.Add(rl.Window.WindowSize)
		retryAfter = 0
	}

	return RateLimitStatus{
		Allowed:       allowed,
		RequestsUsed:  requestsInWindow,
		RequestsLimit: rl.Window.MaxRequests,
		ResetTime:     resetTime,
		RetryAfter:    retryAfter,
		WindowSize:    rl.Window.WindowSize,
	}
}

// GetStatus retorna o status atual sem consumir quota
func (rl *RateLimiter) GetStatus() RateLimitStatus {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	now := time.Now()
	windowStart := now.Add(-rl.Window.WindowSize)

	// Conta requisições na janela atual (sem remover)
	requestsInWindow := 0
	for _, reqTime := range rl.Requests {
		if reqTime.After(windowStart) {
			requestsInWindow++
		}
	}

	allowed := requestsInWindow < rl.Window.MaxRequests && rl.IsActive

	// Calcula resetTime
	var resetTime time.Time
	var retryAfter time.Duration

	if len(rl.Requests) > 0 {
		// Encontra a requisição mais antiga na janela
		for _, reqTime := range rl.Requests {
			if reqTime.After(windowStart) {
				resetTime = reqTime.Add(rl.Window.WindowSize)
				break
			}
		}
		if resetTime.IsZero() || resetTime.Before(now) {
			resetTime = now.Add(rl.Window.WindowSize)
		}
		retryAfter = resetTime.Sub(now)
	} else {
		resetTime = now.Add(rl.Window.WindowSize)
		retryAfter = 0
	}

	return RateLimitStatus{
		Allowed:       allowed,
		RequestsUsed:  requestsInWindow,
		RequestsLimit: rl.Window.MaxRequests,
		ResetTime:     resetTime,
		RetryAfter:    retryAfter,
		WindowSize:    rl.Window.WindowSize,
	}
}

// Reset limpa todas as requisições
func (rl *RateLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.Requests = make([]time.Time, 0)
	rl.UpdatedAt = time.Now()
}

// SetWindow atualiza a janela de rate limiting
func (rl *RateLimiter) SetWindow(window RateLimitWindow) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.Window = window
	rl.UpdatedAt = time.Now()

	// Limpa requisições antigas com a nova janela
	windowStart := time.Now().Add(-window.WindowSize)
	rl.cleanupOldRequests(windowStart)
}

// Activate ativa o rate limiter
func (rl *RateLimiter) Activate() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.IsActive = true
	rl.UpdatedAt = time.Now()
}

// Deactivate desativa o rate limiter
func (rl *RateLimiter) Deactivate() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.IsActive = false
	rl.UpdatedAt = time.Now()
}

// GetRemainingQuota retorna quantas requisições ainda podem ser feitas
func (rl *RateLimiter) GetRemainingQuota() int {
	status := rl.GetStatus()
	remaining := status.RequestsLimit - status.RequestsUsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetUsagePercentage retorna o percentual de uso da quota
func (rl *RateLimiter) GetUsagePercentage() float64 {
	status := rl.GetStatus()
	if status.RequestsLimit == 0 {
		return 0
	}
	return float64(status.RequestsUsed) / float64(status.RequestsLimit) * 100
}

// GetNextResetTime retorna quando o limite será resetado
func (rl *RateLimiter) GetNextResetTime() time.Time {
	return rl.GetStatus().ResetTime
}

// IsQuotaExhausted verifica se a quota foi esgotada
func (rl *RateLimiter) IsQuotaExhausted() bool {
	return !rl.GetStatus().Allowed
}

// cleanupOldRequests remove requisições antigas (deve ser chamado com lock)
func (rl *RateLimiter) cleanupOldRequests(windowStart time.Time) {
	validRequests := make([]time.Time, 0, len(rl.Requests))
	
	for _, reqTime := range rl.Requests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}
	
	rl.Requests = validRequests
}

// GetKey retorna a chave do rate limiter
func (rl *RateLimiter) GetKey() string {
	return fmt.Sprintf("%s:%s", rl.Type, rl.Key)
}

// GetStats retorna estatísticas do rate limiter
func (rl *RateLimiter) GetStats() map[string]interface{} {
	status := rl.GetStatus()
	
	return map[string]interface{}{
		"id":                rl.ID,
		"type":              rl.Type,
		"key":               rl.Key,
		"is_active":         rl.IsActive,
		"requests_used":     status.RequestsUsed,
		"requests_limit":    status.RequestsLimit,
		"remaining_quota":   rl.GetRemainingQuota(),
		"usage_percentage":  rl.GetUsagePercentage(),
		"window_size":       rl.Window.WindowSize.String(),
		"reset_time":        status.ResetTime,
		"last_used_at":      rl.LastUsedAt,
		"created_at":        rl.CreatedAt,
		"updated_at":        rl.UpdatedAt,
	}
}