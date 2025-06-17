package domain

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

// CNPJPoolStrategy define as estratégias de seleção de CNPJ
type CNPJPoolStrategy string

const (
	StrategyRoundRobin   CNPJPoolStrategy = "round_robin"   // Rotação circular
	StrategyLeastUsed    CNPJPoolStrategy = "least_used"    // Menos usado
	StrategyPriority     CNPJPoolStrategy = "priority"      // Por prioridade
	StrategyAvailability CNPJPoolStrategy = "availability"  // Maior disponibilidade
)

// CNPJPool gerencia um pool de CNPJs para balanceamento de carga
type CNPJPool struct {
	ID         uuid.UUID                  `json:"id"`
	TenantID   uuid.UUID                  `json:"tenant_id"`
	Name       string                     `json:"name"`
	Strategy   CNPJPoolStrategy           `json:"strategy"`
	Providers  map[uuid.UUID]*CNPJProvider `json:"providers"`
	IsActive   bool                       `json:"is_active"`
	CreatedAt  time.Time                  `json:"created_at"`
	UpdatedAt  time.Time                  `json:"updated_at"`
	mu         sync.RWMutex               // Mutex para thread safety
	roundIndex int                        // Índice para round robin
}

// CNPJPoolRepository interface para persistência do pool
type CNPJPoolRepository interface {
	Save(pool *CNPJPool) error
	FindByID(id uuid.UUID) (*CNPJPool, error)
	FindByTenantID(tenantID uuid.UUID) ([]*CNPJPool, error)
	Update(pool *CNPJPool) error
	Delete(id uuid.UUID) error
}

// PoolStats estatísticas do pool
type PoolStats struct {
	TotalProviders     int     `json:"total_providers"`
	ActiveProviders    int     `json:"active_providers"`
	TotalDailyLimit    int     `json:"total_daily_limit"`
	TotalDailyUsage    int     `json:"total_daily_usage"`
	AvailableQuota     int     `json:"available_quota"`
	UsagePercentage    float64 `json:"usage_percentage"`
	ProvidersWithQuota int     `json:"providers_with_quota"`
}

// NewCNPJPool cria um novo pool de CNPJs
func NewCNPJPool(tenantID uuid.UUID, name string, strategy CNPJPoolStrategy) *CNPJPool {
	return &CNPJPool{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Name:      name,
		Strategy:  strategy,
		Providers: make(map[uuid.UUID]*CNPJProvider),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// AddProvider adiciona um CNPJ ao pool
func (p *CNPJPool) AddProvider(provider *CNPJProvider) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if provider == nil {
		return errors.New("provider não pode ser nil")
	}

	if _, exists := p.Providers[provider.ID]; exists {
		return errors.New("provider já existe no pool")
	}

	p.Providers[provider.ID] = provider
	p.UpdatedAt = time.Now()

	return nil
}

// RemoveProvider remove um CNPJ do pool
func (p *CNPJPool) RemoveProvider(providerID uuid.UUID) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.Providers[providerID]; !exists {
		return errors.New("provider não encontrado no pool")
	}

	delete(p.Providers, providerID)
	p.UpdatedAt = time.Now()

	return nil
}

// GetNextProvider seleciona o próximo CNPJ baseado na estratégia
func (p *CNPJPool) GetNextProvider() (*CNPJProvider, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.IsActive {
		return nil, errors.New("pool está inativo")
	}

	availableProviders := p.getAvailableProviders()
	if len(availableProviders) == 0 {
		return nil, errors.New("nenhum provider disponível")
	}

	switch p.Strategy {
	case StrategyRoundRobin:
		return p.getNextRoundRobin(availableProviders), nil
	case StrategyLeastUsed:
		return p.getLeastUsed(availableProviders), nil
	case StrategyPriority:
		return p.getByPriority(availableProviders), nil
	case StrategyAvailability:
		return p.getByAvailability(availableProviders), nil
	default:
		return p.getNextRoundRobin(availableProviders), nil
	}
}

// GetProviderWithQuota retorna um provider que tenha quota mínima disponível
func (p *CNPJPool) GetProviderWithQuota(minQuota int) (*CNPJProvider, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if !p.IsActive {
		return nil, errors.New("pool está inativo")
	}

	var bestProvider *CNPJProvider
	maxAvailable := 0

	for _, provider := range p.Providers {
		if provider.IsActive && provider.CanMakeRequest() {
			available := provider.GetAvailableQuota()
			if available >= minQuota && available > maxAvailable {
				maxAvailable = available
				bestProvider = provider
			}
		}
	}

	if bestProvider == nil {
		return nil, errors.New("nenhum provider com quota suficiente")
	}

	return bestProvider, nil
}

// GetStats retorna estatísticas do pool
func (p *CNPJPool) GetStats() PoolStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	stats := PoolStats{
		TotalProviders: len(p.Providers),
	}

	for _, provider := range p.Providers {
		if provider.IsActive {
			stats.ActiveProviders++
			stats.TotalDailyLimit += provider.DailyLimit
			stats.TotalDailyUsage += provider.DailyUsage

			if provider.CanMakeRequest() {
				stats.ProvidersWithQuota++
			}
		}
	}

	stats.AvailableQuota = stats.TotalDailyLimit - stats.TotalDailyUsage
	if stats.TotalDailyLimit > 0 {
		stats.UsagePercentage = float64(stats.TotalDailyUsage) / float64(stats.TotalDailyLimit) * 100
	}

	return stats
}

// ResetAllUsage reseta o uso diário de todos os providers
func (p *CNPJPool) ResetAllUsage() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, provider := range p.Providers {
		provider.ResetDailyUsage()
	}

	p.UpdatedAt = time.Now()
}

// Activate ativa o pool
func (p *CNPJPool) Activate() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.IsActive = true
	p.UpdatedAt = time.Now()
}

// Deactivate desativa o pool
func (p *CNPJPool) Deactivate() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.IsActive = false
	p.UpdatedAt = time.Now()
}

// SetStrategy altera a estratégia de seleção
func (p *CNPJPool) SetStrategy(strategy CNPJPoolStrategy) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.Strategy = strategy
	p.roundIndex = 0 // Reset round robin
	p.UpdatedAt = time.Now()
}

// GetProvider retorna um provider específico
func (p *CNPJPool) GetProvider(providerID uuid.UUID) (*CNPJProvider, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	provider, exists := p.Providers[providerID]
	if !exists {
		return nil, errors.New("provider não encontrado")
	}

	return provider, nil
}

// GetAllProviders retorna todos os providers
func (p *CNPJPool) GetAllProviders() []*CNPJProvider {
	p.mu.RLock()
	defer p.mu.RUnlock()

	providers := make([]*CNPJProvider, 0, len(p.Providers))
	for _, provider := range p.Providers {
		providers = append(providers, provider)
	}

	return providers
}

// getAvailableProviders retorna providers ativos que podem fazer requisições
func (p *CNPJPool) getAvailableProviders() []*CNPJProvider {
	var available []*CNPJProvider

	for _, provider := range p.Providers {
		if provider.IsActive && provider.CanMakeRequest() {
			available = append(available, provider)
		}
	}

	return available
}

// getNextRoundRobin implementa seleção round robin
func (p *CNPJPool) getNextRoundRobin(providers []*CNPJProvider) *CNPJProvider {
	if len(providers) == 0 {
		return nil
	}

	provider := providers[p.roundIndex%len(providers)]
	p.roundIndex++

	return provider
}

// getLeastUsed retorna o provider menos usado
func (p *CNPJPool) getLeastUsed(providers []*CNPJProvider) *CNPJProvider {
	if len(providers) == 0 {
		return nil
	}

	sort.Slice(providers, func(i, j int) bool {
		return providers[i].GetUsagePercentage() < providers[j].GetUsagePercentage()
	})

	return providers[0]
}

// getByPriority retorna o provider de maior prioridade
func (p *CNPJPool) getByPriority(providers []*CNPJProvider) *CNPJProvider {
	if len(providers) == 0 {
		return nil
	}

	sort.Slice(providers, func(i, j int) bool {
		// Prioridade menor = maior prioridade (1 é mais alta que 2)
		if providers[i].Priority != providers[j].Priority {
			return providers[i].Priority < providers[j].Priority
		}
		// Em caso de empate, usa o menos usado
		return providers[i].GetUsagePercentage() < providers[j].GetUsagePercentage()
	})

	return providers[0]
}

// getByAvailability retorna o provider com maior quota disponível
func (p *CNPJPool) getByAvailability(providers []*CNPJProvider) *CNPJProvider {
	if len(providers) == 0 {
		return nil
	}

	sort.Slice(providers, func(i, j int) bool {
		return providers[i].GetAvailableQuota() > providers[j].GetAvailableQuota()
	})

	return providers[0]
}

// ValidatePool valida a configuração do pool
func (p *CNPJPool) ValidatePool() error {
	if p.Name == "" {
		return errors.New("nome do pool é obrigatório")
	}

	if len(p.Providers) == 0 {
		return errors.New("pool deve ter pelo menos um provider")
	}

	activeCount := 0
	for _, provider := range p.Providers {
		if provider.IsActive {
			activeCount++
		}
	}

	if activeCount == 0 {
		return errors.New("pool deve ter pelo menos um provider ativo")
	}

	return nil
}