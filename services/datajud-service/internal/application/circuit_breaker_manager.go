package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/direito-lux/datajud-service/internal/domain"
)

// CircuitBreakerManager gerencia circuit breakers para diferentes recursos
type CircuitBreakerManager struct {
	repos     *domain.Repositories
	breakers  map[string]*domain.CircuitBreaker // Cache de breakers por nome
	mu        sync.RWMutex
	config    domain.DataJudConfig
	healthChecker *time.Ticker
	stopHealth    chan bool
}

// NewCircuitBreakerManager cria novo gerenciador de circuit breakers
func NewCircuitBreakerManager(repos *domain.Repositories, config domain.DataJudConfig) *CircuitBreakerManager {
	manager := &CircuitBreakerManager{
		repos:     repos,
		breakers:  make(map[string]*domain.CircuitBreaker),
		config:    config,
		stopHealth: make(chan bool),
	}

	// Iniciar verificação de saúde periódica
	manager.startHealthCheck()

	return manager
}

// Initialize carrega circuit breakers existentes
func (cbm *CircuitBreakerManager) Initialize(ctx context.Context) error {
	breakers, err := cbm.repos.CircuitBreaker.FindAll()
	if err != nil {
		return fmt.Errorf("erro ao carregar circuit breakers: %w", err)
	}

	cbm.mu.Lock()
	defer cbm.mu.Unlock()

	for _, breaker := range breakers {
		cbm.breakers[breaker.Name] = breaker
	}

	// Criar circuit breakers padrão se não existirem
	cbm.ensureDefaultBreakers(ctx)

	return nil
}

// GetOrCreate obtém ou cria um circuit breaker
func (cbm *CircuitBreakerManager) GetOrCreate(name string) *domain.CircuitBreaker {
	cbm.mu.RLock()
	if breaker, exists := cbm.breakers[name]; exists {
		cbm.mu.RUnlock()
		return breaker
	}
	cbm.mu.RUnlock()

	// Buscar no banco
	breaker, err := cbm.repos.CircuitBreaker.FindByName(name)
	if err != nil || breaker == nil {
		// Criar novo circuit breaker
		breaker = domain.NewDataJudCircuitBreaker(name)
		cbm.repos.CircuitBreaker.Save(breaker)
	}

	cbm.mu.Lock()
	cbm.breakers[name] = breaker
	cbm.mu.Unlock()

	return breaker
}

// GetBreakerByID obtém circuit breaker por ID
func (cbm *CircuitBreakerManager) GetBreakerByID(id uuid.UUID) (*domain.CircuitBreaker, error) {
	breaker, err := cbm.repos.CircuitBreaker.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("circuit breaker não encontrado: %w", err)
	}

	// Atualizar cache se necessário
	cbm.mu.Lock()
	cbm.breakers[breaker.Name] = breaker
	cbm.mu.Unlock()

	return breaker, nil
}

// CreateBreaker cria um novo circuit breaker personalizado
func (cbm *CircuitBreakerManager) CreateBreaker(ctx context.Context, name string, config domain.CircuitBreakerConfig) (*domain.CircuitBreaker, error) {
	cbm.mu.Lock()
	defer cbm.mu.Unlock()

	// Verificar se já existe
	if _, exists := cbm.breakers[name]; exists {
		return nil, fmt.Errorf("circuit breaker %s já existe", name)
	}

	// Criar novo breaker
	breaker := domain.NewCircuitBreaker(name, config)
	
	if err := cbm.repos.CircuitBreaker.Save(breaker); err != nil {
		return nil, fmt.Errorf("erro ao salvar circuit breaker: %w", err)
	}

	cbm.breakers[name] = breaker

	return breaker, nil
}

// UpdateBreakerConfig atualiza configuração de um circuit breaker
func (cbm *CircuitBreakerManager) UpdateBreakerConfig(name string, config domain.CircuitBreakerConfig) error {
	breaker := cbm.GetOrCreate(name)
	breaker.UpdateConfig(config)

	return cbm.repos.CircuitBreaker.Update(breaker)
}

// ResetBreaker força o reset de um circuit breaker
func (cbm *CircuitBreakerManager) ResetBreaker(ctx context.Context, name string) error {
	breaker := cbm.GetOrCreate(name)
	breaker.Reset()

	if err := cbm.repos.CircuitBreaker.Update(breaker); err != nil {
		return fmt.Errorf("erro ao resetar circuit breaker: %w", err)
	}

	// Publicar evento
	event := &domain.CircuitBreakerClosed{
		BaseEvent: domain.BaseEvent{
			ID:          uuid.New(),
			Type:        "datajud.circuit_breaker.closed",
			AggregateID: breaker.ID,
			OccurredAt:  time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		BreakerName:      name,
		SuccessCount:     breaker.SuccessCount,
		SuccessThreshold: breaker.Config.SuccessThreshold,
	}
	cbm.repos.EventStore.SaveEvent(ctx, event)

	return nil
}

// ForceOpenBreaker força a abertura de um circuit breaker
func (cbm *CircuitBreakerManager) ForceOpenBreaker(ctx context.Context, name string) error {
	breaker := cbm.GetOrCreate(name)
	breaker.ForceOpen()

	if err := cbm.repos.CircuitBreaker.Update(breaker); err != nil {
		return fmt.Errorf("erro ao abrir circuit breaker: %w", err)
	}

	// Publicar evento
	event := domain.NewCircuitBreakerOpened(
		breaker.ID,
		name,
		breaker.FailureCount,
		breaker.Config.FailureThreshold,
	)
	cbm.repos.EventStore.SaveEvent(ctx, event)

	return nil
}

// GetBreakerStats obtém estatísticas de um circuit breaker
func (cbm *CircuitBreakerManager) GetBreakerStats(name string) map[string]interface{} {
	breaker := cbm.GetOrCreate(name)
	return breaker.GetStats()
}

// GetAllBreakersStats obtém estatísticas de todos os circuit breakers
func (cbm *CircuitBreakerManager) GetAllBreakersStats() map[string]map[string]interface{} {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()

	stats := make(map[string]map[string]interface{})
	for name, breaker := range cbm.breakers {
		stats[name] = breaker.GetStats()
	}

	return stats
}

// GetHealthyBreakers retorna lista de circuit breakers saudáveis
func (cbm *CircuitBreakerManager) GetHealthyBreakers() []string {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()

	healthy := make([]string, 0)
	for name, breaker := range cbm.breakers {
		if breaker.IsHealthy() {
			healthy = append(healthy, name)
		}
	}

	return healthy
}

// GetUnhealthyBreakers retorna lista de circuit breakers com problema
func (cbm *CircuitBreakerManager) GetUnhealthyBreakers() []string {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()

	unhealthy := make([]string, 0)
	for name, breaker := range cbm.breakers {
		if !breaker.IsHealthy() {
			unhealthy = append(unhealthy, name)
		}
	}

	return unhealthy
}

// ExecuteWithBreaker executa uma função com proteção de circuit breaker
func (cbm *CircuitBreakerManager) ExecuteWithBreaker(ctx context.Context, breakerName string, fn func() error) domain.ExecutionResult {
	breaker := cbm.GetOrCreate(breakerName)
	
	result := breaker.Execute(fn)

	// Salvar mudanças do breaker
	cbm.repos.CircuitBreaker.Update(breaker)

	// Publicar eventos se houve mudança de estado
	currentState := breaker.GetState()
	if result.Success && currentState == domain.StateClosed && breaker.FailureCount == 0 {
		// Circuit breaker fechou após sucesso
		event := &domain.CircuitBreakerClosed{
			BaseEvent: domain.BaseEvent{
				ID:          uuid.New(),
				Type:        "datajud.circuit_breaker.closed",
				AggregateID: breaker.ID,
				OccurredAt:  time.Now(),
				Version:     1,
				Metadata:    make(map[string]interface{}),
			},
			BreakerName:      breakerName,
			SuccessCount:     breaker.SuccessCount,
			SuccessThreshold: breaker.Config.SuccessThreshold,
		}
		cbm.repos.EventStore.SaveEvent(ctx, event)
	} else if !result.Success && currentState == domain.StateOpen {
		// Circuit breaker abriu após falha
		event := domain.NewCircuitBreakerOpened(
			breaker.ID,
			breakerName,
			breaker.FailureCount,
			breaker.Config.FailureThreshold,
		)
		cbm.repos.EventStore.SaveEvent(ctx, event)
	} else if currentState == domain.StateHalfOpen && result.Success {
		// Transição para half-open após timeout
		event := &domain.CircuitBreakerHalfOpened{
			BaseEvent: domain.BaseEvent{
				ID:          uuid.New(),
				Type:        "datajud.circuit_breaker.half_opened",
				AggregateID: breaker.ID,
				OccurredAt:  time.Now(),
				Version:     1,
				Metadata:    make(map[string]interface{}),
			},
			BreakerName: breakerName,
			Timeout:     breaker.Config.Timeout,
		}
		cbm.repos.EventStore.SaveEvent(ctx, event)
	}

	return result
}

// GetBreakersByState retorna circuit breakers em um estado específico
func (cbm *CircuitBreakerManager) GetBreakersByState(state domain.CircuitBreakerState) []string {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()

	breakers := make([]string, 0)
	for name, breaker := range cbm.breakers {
		if breaker.GetState() == state {
			breakers = append(breakers, name)
		}
	}

	return breakers
}

// GetFailureRate obtém taxa de falha geral do sistema
func (cbm *CircuitBreakerManager) GetFailureRate() float64 {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()

	totalFailures := 0
	totalRequests := 0

	for _, breaker := range cbm.breakers {
		totalFailures += breaker.FailureCount
		totalRequests += breaker.FailureCount + breaker.SuccessCount
	}

	if totalRequests == 0 {
		return 0
	}

	return float64(totalFailures) / float64(totalRequests) * 100
}

// ActivateBreaker ativa um circuit breaker
func (cbm *CircuitBreakerManager) ActivateBreaker(name string) error {
	breaker := cbm.GetOrCreate(name)
	breaker.Activate()
	return cbm.repos.CircuitBreaker.Update(breaker)
}

// DeactivateBreaker desativa um circuit breaker
func (cbm *CircuitBreakerManager) DeactivateBreaker(name string) error {
	breaker := cbm.GetOrCreate(name)
	breaker.Deactivate()
	return cbm.repos.CircuitBreaker.Update(breaker)
}

// Stop para o gerenciador e limpa recursos
func (cbm *CircuitBreakerManager) Stop() {
	if cbm.healthChecker != nil {
		cbm.healthChecker.Stop()
	}
	close(cbm.stopHealth)
}

// startHealthCheck inicia verificação de saúde periódica
func (cbm *CircuitBreakerManager) startHealthCheck() {
	cbm.healthChecker = time.NewTicker(cbm.config.HealthCheckInterval)
	
	go func() {
		for {
			select {
			case <-cbm.healthChecker.C:
				cbm.performHealthCheck()
			case <-cbm.stopHealth:
				return
			}
		}
	}()
}

// performHealthCheck verifica saúde dos circuit breakers
func (cbm *CircuitBreakerManager) performHealthCheck() {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()

	for name, breaker := range cbm.breakers {
		// Verificar se circuit breakers em half-open há muito tempo devem ser resetados
		if breaker.GetState() == domain.StateHalfOpen {
			timeSinceChange := time.Since(breaker.StateChangedAt)
			if timeSinceChange > breaker.Config.Timeout*2 {
				// Reset automático se ficou muito tempo em half-open
				breaker.Reset()
				cbm.repos.CircuitBreaker.Update(breaker)
			}
		}
		
		// Log estatísticas para monitoramento
		stats := breaker.GetStats()
		if !breaker.IsHealthy() {
			// Seria enviado para sistema de monitoramento
			_ = fmt.Sprintf("Circuit breaker unhealthy: %s, stats: %+v", name, stats)
		}
	}
}

// ensureDefaultBreakers cria circuit breakers padrão necessários
func (cbm *CircuitBreakerManager) ensureDefaultBreakers(ctx context.Context) {
	defaultBreakers := []string{
		"datajud-api",           // API principal DataJud
		"datajud-auth",          // Autenticação DataJud
		"database",              // Banco de dados
		"cache",                 // Sistema de cache
		"rabbitmq",              // Message broker
	}

	for _, name := range defaultBreakers {
		if _, exists := cbm.breakers[name]; !exists {
			config := domain.CircuitBreakerConfig{
				FailureThreshold: cbm.config.CBFailureThreshold,
				SuccessThreshold: cbm.config.CBSuccessThreshold,
				Timeout:          cbm.config.CBTimeout,
				MaxRequests:      cbm.config.CBMaxRequests,
			}
			
			breaker := domain.NewCircuitBreaker(name, config)
			cbm.repos.CircuitBreaker.Save(breaker)
			cbm.breakers[name] = breaker
		}
	}
}

// GetBreakerHistory obtém histórico de mudanças de estado (seria implementado)
func (cbm *CircuitBreakerManager) GetBreakerHistory(name string, hours int) []map[string]interface{} {
	// Implementação seria buscar eventos de mudança de estado nos últimos X horas
	// Por simplicidade, retornando vazio
	return []map[string]interface{}{}
}

// GetSystemHealth retorna saúde geral baseada nos circuit breakers
func (cbm *CircuitBreakerManager) GetSystemHealth() map[string]interface{} {
	cbm.mu.RLock()
	defer cbm.mu.RUnlock()

	totalBreakers := len(cbm.breakers)
	healthyCount := 0
	openCount := 0
	halfOpenCount := 0

	for _, breaker := range cbm.breakers {
		if breaker.IsHealthy() {
			healthyCount++
		}
		
		switch breaker.GetState() {
		case domain.StateOpen:
			openCount++
		case domain.StateHalfOpen:
			halfOpenCount++
		}
	}

	healthPercentage := float64(healthyCount) / float64(totalBreakers) * 100
	
	status := "healthy"
	if healthPercentage < 80 {
		status = "degraded"
	}
	if healthPercentage < 50 {
		status = "unhealthy"
	}

	return map[string]interface{}{
		"status":            status,
		"health_percentage": healthPercentage,
		"total_breakers":    totalBreakers,
		"healthy_breakers":  healthyCount,
		"open_breakers":     openCount,
		"half_open_breakers": halfOpenCount,
		"failure_rate":      cbm.GetFailureRate(),
		"timestamp":         time.Now(),
	}
}