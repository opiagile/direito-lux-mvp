package application

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/direito-lux/datajud-service/internal/domain"
)

// CNPJPoolManager gerencia pools de CNPJs por tenant
type CNPJPoolManager struct {
	repos       *domain.Repositories
	pools       map[uuid.UUID]*domain.CNPJPool // Cache de pools por tenant
	defaultPool *domain.CNPJPool               // Pool padrão global
	mu          sync.RWMutex
	config      domain.DataJudConfig
}

// NewCNPJPoolManager cria novo gerenciador de pools
func NewCNPJPoolManager(repos *domain.Repositories, config domain.DataJudConfig) *CNPJPoolManager {
	return &CNPJPoolManager{
		repos:  repos,
		pools:  make(map[uuid.UUID]*domain.CNPJPool),
		config: config,
	}
}

// Initialize inicializa o gerenciador com pools existentes
func (pm *CNPJPoolManager) Initialize(ctx context.Context) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Criar pool padrão se não existir
	if err := pm.ensureDefaultPool(ctx); err != nil {
		return fmt.Errorf("erro ao criar pool padrão: %w", err)
	}

	return nil
}

// GetTenantPool obtém o pool de um tenant específico
func (pm *CNPJPoolManager) GetTenantPool(tenantID uuid.UUID) (*domain.CNPJPool, error) {
	pm.mu.RLock()
	if pool, exists := pm.pools[tenantID]; exists {
		pm.mu.RUnlock()
		return pool, nil
	}
	pm.mu.RUnlock()

	// Buscar no banco se não estiver em cache
	pools, err := pm.repos.CNPJPool.FindByTenantID(tenantID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pools do tenant: %w", err)
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Se o tenant tem pool próprio, usar ele
	if len(pools) > 0 {
		pool := pools[0] // Usar o primeiro pool ativo
		if pool.IsActive {
			// Carregar providers do pool
			if err := pm.loadPoolProviders(pool); err != nil {
				return nil, fmt.Errorf("erro ao carregar providers do pool: %w", err)
			}
			pm.pools[tenantID] = pool
			return pool, nil
		}
	}

	// Se não tem pool próprio, criar um baseado nos CNPJs do tenant
	pool, err := pm.createTenantPool(tenantID)
	if err != nil {
		// Se falhar, usar pool padrão
		return pm.defaultPool, nil
	}

	pm.pools[tenantID] = pool
	return pool, nil
}

// CreatePool cria um novo pool para um tenant
func (pm *CNPJPoolManager) CreatePool(ctx context.Context, req *CreateCNPJPoolRequest) (*CNPJPoolResponse, error) {
	// Validar requisição
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Criar pool
	pool := domain.NewCNPJPool(req.TenantID, req.Name, req.Strategy)

	// Adicionar providers se especificados
	if len(req.ProviderIDs) > 0 {
		for _, providerID := range req.ProviderIDs {
			provider, err := pm.repos.CNPJProvider.FindByID(providerID)
			if err != nil {
				return nil, fmt.Errorf("provider %s não encontrado: %w", providerID, err)
			}

			// Verificar se o provider pertence ao tenant
			if provider.TenantID != req.TenantID {
				return nil, fmt.Errorf("provider %s não pertence ao tenant", providerID)
			}

			if err := pool.AddProvider(provider); err != nil {
				return nil, fmt.Errorf("erro ao adicionar provider ao pool: %w", err)
			}
		}
	}

	// Validar pool
	if err := pool.ValidatePool(); err != nil {
		return nil, fmt.Errorf("pool inválido: %w", err)
	}

	// Salvar no banco
	if err := pm.repos.CNPJPool.Save(pool); err != nil {
		return nil, fmt.Errorf("erro ao salvar pool: %w", err)
	}

	// Adicionar ao cache
	pm.mu.Lock()
	pm.pools[req.TenantID] = pool
	pm.mu.Unlock()

	// Publicar evento
	event := &domain.CNPJPoolCreated{
		BaseEvent: domain.BaseEvent{
			ID:          uuid.New(),
			Type:        "datajud.cnpj_pool.created",
			AggregateID: pool.ID,
			OccurredAt:  pool.CreatedAt,
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		TenantID:    req.TenantID,
		PoolName:    req.Name,
		Strategy:    req.Strategy,
		ProviderIDs: req.ProviderIDs,
	}
	pm.repos.EventStore.SaveEvent(ctx, event)

	return pm.buildPoolResponse(pool), nil
}

// AddProviderToPool adiciona um provider a um pool
func (pm *CNPJPoolManager) AddProviderToPool(ctx context.Context, poolID, providerID uuid.UUID) error {
	// Buscar pool
	pool, err := pm.repos.CNPJPool.FindByID(poolID)
	if err != nil {
		return fmt.Errorf("pool não encontrado: %w", err)
	}

	// Buscar provider
	provider, err := pm.repos.CNPJProvider.FindByID(providerID)
	if err != nil {
		return fmt.Errorf("provider não encontrado: %w", err)
	}

	// Verificar se pertencem ao mesmo tenant
	if pool.TenantID != provider.TenantID {
		return fmt.Errorf("pool e provider devem pertencer ao mesmo tenant")
	}

	// Adicionar ao pool
	if err := pool.AddProvider(provider); err != nil {
		return fmt.Errorf("erro ao adicionar provider ao pool: %w", err)
	}

	// Salvar
	if err := pm.repos.CNPJPool.Update(pool); err != nil {
		return fmt.Errorf("erro ao atualizar pool: %w", err)
	}

	// Atualizar cache
	pm.mu.Lock()
	pm.pools[pool.TenantID] = pool
	pm.mu.Unlock()

	// Publicar evento
	event := &domain.CNPJProviderAddedToPool{
		BaseEvent: domain.BaseEvent{
			ID:          uuid.New(),
			Type:        "datajud.cnpj_provider.added_to_pool",
			AggregateID: poolID,
			OccurredAt:  pool.UpdatedAt,
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		PoolID:     poolID,
		ProviderID: providerID,
		CNPJ:       provider.CNPJ,
	}
	pm.repos.EventStore.SaveEvent(ctx, event)

	return nil
}

// RemoveProviderFromPool remove um provider de um pool
func (pm *CNPJPoolManager) RemoveProviderFromPool(ctx context.Context, poolID, providerID uuid.UUID, reason string) error {
	// Buscar pool
	pool, err := pm.repos.CNPJPool.FindByID(poolID)
	if err != nil {
		return fmt.Errorf("pool não encontrado: %w", err)
	}

	// Buscar provider para obter CNPJ
	provider, err := pool.GetProvider(providerID)
	if err != nil {
		return fmt.Errorf("provider não encontrado no pool: %w", err)
	}

	// Remover do pool
	if err := pool.RemoveProvider(providerID); err != nil {
		return fmt.Errorf("erro ao remover provider do pool: %w", err)
	}

	// Salvar
	if err := pm.repos.CNPJPool.Update(pool); err != nil {
		return fmt.Errorf("erro ao atualizar pool: %w", err)
	}

	// Atualizar cache
	pm.mu.Lock()
	pm.pools[pool.TenantID] = pool
	pm.mu.Unlock()

	// Publicar evento
	event := &domain.CNPJProviderRemovedFromPool{
		BaseEvent: domain.BaseEvent{
			ID:          uuid.New(),
			Type:        "datajud.cnpj_provider.removed_from_pool",
			AggregateID: poolID,
			OccurredAt:  pool.UpdatedAt,
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		PoolID:     poolID,
		ProviderID: providerID,
		CNPJ:       provider.CNPJ,
		Reason:     reason,
	}
	pm.repos.EventStore.SaveEvent(ctx, event)

	return nil
}

// UpdatePoolStrategy atualiza a estratégia de um pool
func (pm *CNPJPoolManager) UpdatePoolStrategy(ctx context.Context, poolID uuid.UUID, newStrategy domain.CNPJPoolStrategy) error {
	// Buscar pool
	pool, err := pm.repos.CNPJPool.FindByID(poolID)
	if err != nil {
		return fmt.Errorf("pool não encontrado: %w", err)
	}

	oldStrategy := pool.Strategy
	pool.SetStrategy(newStrategy)

	// Salvar
	if err := pm.repos.CNPJPool.Update(pool); err != nil {
		return fmt.Errorf("erro ao atualizar pool: %w", err)
	}

	// Atualizar cache
	pm.mu.Lock()
	pm.pools[pool.TenantID] = pool
	pm.mu.Unlock()

	// Publicar evento
	event := &domain.CNPJPoolStrategyChanged{
		BaseEvent: domain.BaseEvent{
			ID:          uuid.New(),
			Type:        "datajud.cnpj_pool.strategy_changed",
			AggregateID: poolID,
			OccurredAt:  pool.UpdatedAt,
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		PoolID:      poolID,
		OldStrategy: oldStrategy,
		NewStrategy: newStrategy,
	}
	pm.repos.EventStore.SaveEvent(ctx, event)

	return nil
}

// GetPoolStats obtém estatísticas de um pool
func (pm *CNPJPoolManager) GetPoolStats(poolID uuid.UUID) (*domain.PoolStats, error) {
	pool, err := pm.repos.CNPJPool.FindByID(poolID)
	if err != nil {
		return nil, fmt.Errorf("pool não encontrado: %w", err)
	}

	// Carregar providers se necessário
	if len(pool.Providers) == 0 {
		if err := pm.loadPoolProviders(pool); err != nil {
			return nil, fmt.Errorf("erro ao carregar providers: %w", err)
		}
	}

	stats := pool.GetStats()
	return &stats, nil
}

// ListTenantPools lista todos os pools de um tenant
func (pm *CNPJPoolManager) ListTenantPools(tenantID uuid.UUID) ([]*CNPJPoolResponse, error) {
	pools, err := pm.repos.CNPJPool.FindByTenantID(tenantID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pools: %w", err)
	}

	responses := make([]*CNPJPoolResponse, 0, len(pools))
	for _, pool := range pools {
		// Carregar providers
		if err := pm.loadPoolProviders(pool); err != nil {
			continue // Pular pools com erro
		}
		responses = append(responses, pm.buildPoolResponse(pool))
	}

	return responses, nil
}

// ResetAllPoolsUsage reseta o uso diário de todos os pools
func (pm *CNPJPoolManager) ResetAllPoolsUsage(ctx context.Context) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for _, pool := range pm.pools {
		pool.ResetAllUsage()
		pm.repos.CNPJPool.Update(pool)
	}

	return nil
}

// ensureDefaultPool garante que existe um pool padrão
func (pm *CNPJPoolManager) ensureDefaultPool(ctx context.Context) error {
	// Buscar providers ativos sem tenant específico ou configurar um pool global
	providers, err := pm.repos.CNPJProvider.FindActiveCNPJs()
	if err != nil {
		return err
	}

	if len(providers) == 0 {
		return fmt.Errorf("nenhum CNPJ provider ativo encontrado")
	}

	// Criar pool padrão
	pm.defaultPool = domain.NewCNPJPool(
		uuid.Nil, // Tenant nil para pool global
		"default-pool",
		pm.config.DefaultPoolStrategy,
	)

	// Adicionar providers disponíveis
	for _, provider := range providers {
		if err := pm.defaultPool.AddProvider(provider); err != nil {
			continue // Pular providers com erro
		}
	}

	return nil
}

// createTenantPool cria um pool automaticamente para um tenant
func (pm *CNPJPoolManager) createTenantPool(tenantID uuid.UUID) (*domain.CNPJPool, error) {
	// Buscar CNPJs do tenant
	providers, err := pm.repos.CNPJProvider.FindByTenantID(tenantID)
	if err != nil {
		return nil, err
	}

	if len(providers) == 0 {
		return nil, fmt.Errorf("tenant não possui CNPJs cadastrados")
	}

	// Criar pool
	pool := domain.NewCNPJPool(
		tenantID,
		fmt.Sprintf("auto-pool-%s", tenantID.String()[:8]),
		pm.config.DefaultPoolStrategy,
	)

	// Adicionar providers ativos
	activeCount := 0
	for _, provider := range providers {
		if provider.IsActive {
			if err := pool.AddProvider(provider); err == nil {
				activeCount++
			}
		}
	}

	if activeCount == 0 {
		return nil, fmt.Errorf("nenhum CNPJ ativo encontrado para o tenant")
	}

	// Salvar pool
	if err := pm.repos.CNPJPool.Save(pool); err != nil {
		return nil, err
	}

	return pool, nil
}

// loadPoolProviders carrega os providers de um pool
func (pm *CNPJPoolManager) loadPoolProviders(pool *domain.CNPJPool) error {
	if len(pool.Providers) > 0 {
		return nil // Já carregado
	}

	// Buscar providers do tenant
	providers, err := pm.repos.CNPJProvider.FindByTenantID(pool.TenantID)
	if err != nil {
		return err
	}

	// Adicionar providers ativos ao pool
	for _, provider := range providers {
		if provider.IsActive {
			pool.AddProvider(provider)
		}
	}

	return nil
}

// buildPoolResponse constrói a resposta do pool
func (pm *CNPJPoolManager) buildPoolResponse(pool *domain.CNPJPool) *CNPJPoolResponse {
	stats := pool.GetStats()
	providers := make([]CNPJProviderResponse, 0, len(pool.Providers))

	for _, provider := range pool.GetAllProviders() {
		providers = append(providers, CNPJProviderResponse{
			ID:              provider.ID,
			TenantID:        provider.TenantID,
			CNPJ:            provider.CNPJ,
			Name:            provider.Name,
			Email:           provider.Email,
			DailyLimit:      provider.DailyLimit,
			DailyUsage:      provider.DailyUsage,
			AvailableQuota:  provider.GetAvailableQuota(),
			UsagePercentage: provider.GetUsagePercentage(),
			UsageResetTime:  provider.UsageResetTime,
			IsActive:        provider.IsActive,
			Priority:        provider.Priority,
			LastUsedAt:      provider.LastUsedAt,
			CreatedAt:       provider.CreatedAt,
			UpdatedAt:       provider.UpdatedAt,
		})
	}

	return &CNPJPoolResponse{
		ID:        pool.ID,
		TenantID:  pool.TenantID,
		Name:      pool.Name,
		Strategy:  pool.Strategy,
		IsActive:  pool.IsActive,
		Stats:     stats,
		Providers: providers,
		CreatedAt: pool.CreatedAt,
		UpdatedAt: pool.UpdatedAt,
	}
}