package domain

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// CircuitBreakerState define os estados do circuit breaker
type CircuitBreakerState string

const (
	StateClosed   CircuitBreakerState = "closed"    // Funcionando normalmente
	StateOpen     CircuitBreakerState = "open"      // Bloqueando requisições
	StateHalfOpen CircuitBreakerState = "half_open" // Testando se pode voltar ao normal
)

// CircuitBreakerConfig configuração do circuit breaker
type CircuitBreakerConfig struct {
	FailureThreshold int           `json:"failure_threshold"` // Número de falhas para abrir
	SuccessThreshold int           `json:"success_threshold"` // Sucessos para fechar quando half-open
	Timeout          time.Duration `json:"timeout"`           // Tempo para tentar half-open
	MaxRequests      int           `json:"max_requests"`      // Max requests em half-open
}

// CircuitBreaker implementa o padrão Circuit Breaker
type CircuitBreaker struct {
	ID              uuid.UUID            `json:"id"`
	Name            string               `json:"name"`
	Config          CircuitBreakerConfig `json:"config"`
	State           CircuitBreakerState  `json:"state"`
	FailureCount    int                  `json:"failure_count"`
	SuccessCount    int                  `json:"success_count"`
	RequestCount    int                  `json:"request_count"`
	LastFailureTime *time.Time           `json:"last_failure_time"`
	LastSuccessTime *time.Time           `json:"last_success_time"`
	StateChangedAt  time.Time            `json:"state_changed_at"`
	IsActive        bool                 `json:"is_active"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	mu              sync.RWMutex         // Mutex para thread safety
}

// CircuitBreakerRepository interface para persistência
type CircuitBreakerRepository interface {
	Save(cb *CircuitBreaker) error
	FindByName(name string) (*CircuitBreaker, error)
	FindByID(id uuid.UUID) (*CircuitBreaker, error)
	FindAll() ([]*CircuitBreaker, error)
	Update(cb *CircuitBreaker) error
	Delete(id uuid.UUID) error
	GetStats() (map[string]interface{}, error)
}

// ExecutionResult resultado da execução
type ExecutionResult struct {
	Success     bool          `json:"success"`
	Error       error         `json:"error,omitempty"`
	Duration    time.Duration `json:"duration"`
	ExecutedAt  time.Time     `json:"executed_at"`
	State       CircuitBreakerState `json:"state"`
	Allowed     bool          `json:"allowed"`
}

// NewCircuitBreaker cria um novo circuit breaker
func NewCircuitBreaker(name string, config CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		ID:             uuid.New(),
		Name:           name,
		Config:         config,
		State:          StateClosed,
		FailureCount:   0,
		SuccessCount:   0,
		RequestCount:   0,
		StateChangedAt: time.Now(),
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// NewDataJudCircuitBreaker cria um circuit breaker específico para DataJud
func NewDataJudCircuitBreaker(name string) *CircuitBreaker {
	return NewCircuitBreaker(name, CircuitBreakerConfig{
		FailureThreshold: 5,              // 5 falhas consecutivas
		SuccessThreshold: 3,              // 3 sucessos para voltar ao normal
		Timeout:          30 * time.Second, // 30 segundos para tentar novamente
		MaxRequests:      5,              // Max 5 requests em half-open
	})
}

// CanExecute verifica se uma requisição pode ser executada
func (cb *CircuitBreaker) CanExecute() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if !cb.IsActive {
		return false
	}

	switch cb.State {
	case StateClosed:
		return true
	case StateOpen:
		return cb.shouldAttemptReset()
	case StateHalfOpen:
		return cb.RequestCount < cb.Config.MaxRequests
	default:
		return false
	}
}

// Execute executa uma função com proteção do circuit breaker
func (cb *CircuitBreaker) Execute(fn func() error) ExecutionResult {
	start := time.Now()
	result := ExecutionResult{
		ExecutedAt: start,
	}

	if !cb.CanExecute() {
		result.Allowed = false
		result.Error = errors.New("circuit breaker is open")
		result.Duration = time.Since(start)
		result.State = cb.GetState()
		return result
	}

	result.Allowed = true
	
	// Executa a função
	err := fn()
	result.Duration = time.Since(start)
	result.Success = err == nil
	result.Error = err

	// Registra o resultado
	if result.Success {
		cb.OnSuccess()
	} else {
		cb.OnFailure()
	}

	result.State = cb.GetState()
	return result
}

// OnSuccess registra um sucesso
func (cb *CircuitBreaker) OnSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()
	cb.LastSuccessTime = &now
	cb.UpdatedAt = now

	switch cb.State {
	case StateClosed:
		cb.FailureCount = 0 // Reset failure count
	case StateHalfOpen:
		cb.SuccessCount++
		cb.RequestCount++
		if cb.SuccessCount >= cb.Config.SuccessThreshold {
			cb.setState(StateClosed)
		}
	}
}

// OnFailure registra uma falha
func (cb *CircuitBreaker) OnFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()
	cb.LastFailureTime = &now
	cb.UpdatedAt = now

	switch cb.State {
	case StateClosed:
		cb.FailureCount++
		if cb.FailureCount >= cb.Config.FailureThreshold {
			cb.setState(StateOpen)
		}
	case StateHalfOpen:
		cb.RequestCount++
		cb.setState(StateOpen)
	}
}

// GetState retorna o estado atual (thread-safe)
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state()
}

// state retorna o estado atual (deve ser chamado com lock)
func (cb *CircuitBreaker) state() CircuitBreakerState {
	if cb.State == StateOpen && cb.shouldAttemptReset() {
		cb.setState(StateHalfOpen)
	}
	return cb.State
}

// shouldAttemptReset verifica se deve tentar resetar (deve ser chamado com lock)
func (cb *CircuitBreaker) shouldAttemptReset() bool {
	return time.Since(cb.StateChangedAt) >= cb.Config.Timeout
}

// setState muda o estado (deve ser chamado com lock)
func (cb *CircuitBreaker) setState(newState CircuitBreakerState) {
	if cb.State != newState {
		cb.State = newState
		cb.StateChangedAt = time.Now()
		cb.FailureCount = 0
		cb.SuccessCount = 0
		cb.RequestCount = 0
	}
}

// Reset força o reset do circuit breaker
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.setState(StateClosed)
	cb.UpdatedAt = time.Now()
}

// ForceOpen força a abertura do circuit breaker
func (cb *CircuitBreaker) ForceOpen() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.setState(StateOpen)
	cb.UpdatedAt = time.Now()
}

// UpdateConfig atualiza a configuração
func (cb *CircuitBreaker) UpdateConfig(config CircuitBreakerConfig) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.Config = config
	cb.UpdatedAt = time.Now()
}

// GetStats retorna estatísticas do circuit breaker
func (cb *CircuitBreaker) GetStats() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return map[string]interface{}{
		"id":                 cb.ID,
		"name":               cb.Name,
		"state":              cb.state(),
		"is_active":          cb.IsActive,
		"failure_count":      cb.FailureCount,
		"success_count":      cb.SuccessCount,
		"request_count":      cb.RequestCount,
		"failure_threshold":  cb.Config.FailureThreshold,
		"success_threshold":  cb.Config.SuccessThreshold,
		"timeout":            cb.Config.Timeout.String(),
		"max_requests":       cb.Config.MaxRequests,
		"last_failure_time":  cb.LastFailureTime,
		"last_success_time":  cb.LastSuccessTime,
		"state_changed_at":   cb.StateChangedAt,
		"time_since_change":  time.Since(cb.StateChangedAt).String(),
		"can_execute":        cb.CanExecute(),
		"created_at":         cb.CreatedAt,
		"updated_at":         cb.UpdatedAt,
	}
}

// Activate ativa o circuit breaker
func (cb *CircuitBreaker) Activate() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.IsActive = true
	cb.UpdatedAt = time.Now()
}

// Deactivate desativa o circuit breaker
func (cb *CircuitBreaker) Deactivate() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.IsActive = false
	cb.UpdatedAt = time.Now()
}

// GetFailureRate retorna a taxa de falha
func (cb *CircuitBreaker) GetFailureRate() float64 {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	totalRequests := cb.FailureCount + cb.SuccessCount
	if totalRequests == 0 {
		return 0
	}

	return float64(cb.FailureCount) / float64(totalRequests) * 100
}

// GetUptime retorna o tempo de atividade
func (cb *CircuitBreaker) GetUptime() time.Duration {
	return time.Since(cb.CreatedAt)
}

// IsHealthy verifica se o circuit breaker está saudável
func (cb *CircuitBreaker) IsHealthy() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state() {
	case StateClosed:
		return true
	case StateHalfOpen:
		return cb.SuccessCount > 0
	case StateOpen:
		return false
	default:
		return false
	}
}

// CircuitBreakerManager gerencia múltiplos circuit breakers
type CircuitBreakerManager struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
}

// NewCircuitBreakerManager cria um novo gerenciador
func NewCircuitBreakerManager() *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// GetOrCreate obtém ou cria um circuit breaker
func (cbm *CircuitBreakerManager) GetOrCreate(name string, config CircuitBreakerConfig) *CircuitBreaker {
	cbm.mu.Lock()
	defer cbm.mu.Unlock()

	if cb, exists := cbm.breakers[name]; exists {
		return cb
	}

	cb := NewCircuitBreaker(name, config)
	cbm.breakers[name] = cb
	return cb
}

// Get obtém um circuit breaker existente
func (cbm *CircuitBreakerManager) Get(name string) *CircuitBreaker {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()

	return cbm.breakers[name]
}

// GetAll retorna todos os circuit breakers
func (cbm *CircuitBreakerManager) GetAll() map[string]*CircuitBreaker {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()

	result := make(map[string]*CircuitBreaker)
	for name, cb := range cbm.breakers {
		result[name] = cb
	}
	return result
}

// Remove remove um circuit breaker
func (cbm *CircuitBreakerManager) Remove(name string) {
	cbm.mu.Lock()
	defer cbm.mu.Unlock()

	delete(cbm.breakers, name)
}