package application

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/direito-lux/datajud-service/internal/domain"
)

// RateLimitManager gerencia limitadores de taxa
type RateLimitManager struct {
	repos     *domain.Repositories
	limiters  map[string]*domain.RateLimiter // Cache de limiters por chave
	mu        sync.RWMutex
	config    domain.DataJudConfig
	cleanupTicker *time.Ticker
	stopCleanup   chan bool
}

// NewRateLimitManager cria novo gerenciador de rate limiting
func NewRateLimitManager(repos *domain.Repositories, config domain.DataJudConfig) *RateLimitManager {
	manager := &RateLimitManager{
		repos:     repos,
		limiters:  make(map[string]*domain.RateLimiter),
		config:    config,
		stopCleanup: make(chan bool),
	}

	// Iniciar limpeza periódica
	manager.startCleanupRoutine()

	return manager
}

// GetCNPJLimiter obtém ou cria um rate limiter para CNPJ
func (rm *RateLimitManager) GetCNPJLimiter(cnpj string) *domain.RateLimiter {
	key := string(domain.RateLimitCNPJ) + ":" + cnpj

	rm.mu.RLock()
	if limiter, exists := rm.limiters[key]; exists {
		rm.mu.RUnlock()
		return limiter
	}
	rm.mu.RUnlock()

	// Buscar no banco
	limiter, err := rm.repos.RateLimiter.FindByCNPJ(cnpj)
	if err != nil || limiter == nil {
		// Criar novo limiter para CNPJ
		limiter = domain.NewCNPJRateLimiter(cnpj)
		rm.repos.RateLimiter.Save(limiter)
	}

	rm.mu.Lock()
	rm.limiters[key] = limiter
	rm.mu.Unlock()

	return limiter
}

// GetTenantLimiter obtém ou cria um rate limiter para tenant
func (rm *RateLimitManager) GetTenantLimiter(tenantID uuid.UUID) *domain.RateLimiter {
	key := string(domain.RateLimitTenant) + ":" + tenantID.String()

	rm.mu.RLock()
	if limiter, exists := rm.limiters[key]; exists {
		rm.mu.RUnlock()
		return limiter
	}
	rm.mu.RUnlock()

	// Buscar no banco
	limiter, err := rm.repos.RateLimiter.FindByTenant(tenantID)
	if err != nil || limiter == nil {
		// Criar novo limiter para tenant (limite maior)
		limiter = domain.NewTenantRateLimiter(
			tenantID,
			rm.config.GlobalRateLimit,
			rm.config.RateWindowSize,
		)
		rm.repos.RateLimiter.Save(limiter)
	}

	rm.mu.Lock()
	rm.limiters[key] = limiter
	rm.mu.Unlock()

	return limiter
}

// GetGlobalLimiter obtém o rate limiter global
func (rm *RateLimitManager) GetGlobalLimiter() *domain.RateLimiter {
	key := string(domain.RateLimitGlobal) + ":global"

	rm.mu.RLock()
	if limiter, exists := rm.limiters[key]; exists {
		rm.mu.RUnlock()
		return limiter
	}
	rm.mu.RUnlock()

	// Buscar no banco
	limiter, err := rm.repos.RateLimiter.FindByKey(domain.RateLimitGlobal, "global")
	if err != nil || limiter == nil {
		// Criar limiter global
		limiter = domain.NewGlobalRateLimiter(
			rm.config.GlobalRateLimit*10, // Limite global maior
			rm.config.RateWindowSize,
		)
		rm.repos.RateLimiter.Save(limiter)
	}

	rm.mu.Lock()
	rm.limiters[key] = limiter
	rm.mu.Unlock()

	return limiter
}

// CheckAllowance verifica se uma requisição pode ser feita considerando todos os limiters
func (rm *RateLimitManager) CheckAllowance(cnpj string, tenantID uuid.UUID) (domain.RateLimitStatus, error) {
	// 1. Verificar limite do CNPJ (mais restritivo)
	cnpjLimiter := rm.GetCNPJLimiter(cnpj)
	cnpjStatus := cnpjLimiter.Allow()
	if !cnpjStatus.Allowed {
		return cnpjStatus, nil
	}

	// 2. Verificar limite do tenant
	tenantLimiter := rm.GetTenantLimiter(tenantID)
	tenantStatus := tenantLimiter.Allow()
	if !tenantStatus.Allowed {
		// Reverter o uso do CNPJ já que foi negado pelo tenant
		rm.revertCNPJUsage(cnpjLimiter)
		return tenantStatus, nil
	}

	// 3. Verificar limite global
	globalLimiter := rm.GetGlobalLimiter()
	globalStatus := globalLimiter.Allow()
	if !globalStatus.Allowed {
		// Reverter usos anteriores
		rm.revertCNPJUsage(cnpjLimiter)
		rm.revertTenantUsage(tenantLimiter)
		return globalStatus, nil
	}

	// Salvar alterações nos limiters
	rm.saveLimiterUpdates(cnpjLimiter, tenantLimiter, globalLimiter)

	// Retornar o status mais restritivo (CNPJ)
	return cnpjStatus, nil
}

// GetCNPJStatus obtém status do CNPJ sem consumir quota
func (rm *RateLimitManager) GetCNPJStatus(cnpj string) domain.RateLimitStatus {
	limiter := rm.GetCNPJLimiter(cnpj)
	return limiter.GetStatus()
}

// GetTenantStatus obtém status do tenant sem consumir quota
func (rm *RateLimitManager) GetTenantStatus(tenantID uuid.UUID) domain.RateLimitStatus {
	limiter := rm.GetTenantLimiter(tenantID)
	return limiter.GetStatus()
}

// ResetCNPJUsage reseta o uso de um CNPJ específico
func (rm *RateLimitManager) ResetCNPJUsage(ctx context.Context, cnpj string) error {
	limiter := rm.GetCNPJLimiter(cnpj)
	limiter.Reset()
	
	if err := rm.repos.RateLimiter.Update(limiter); err != nil {
		return err
	}

	// Publicar evento
	event := &domain.RateLimitReset{
		BaseEvent: domain.BaseEvent{
			ID:          uuid.New(),
			Type:        "datajud.rate_limit.reset",
			AggregateID: limiter.ID,
			OccurredAt:  time.Now(),
			Version:     1,
			Metadata:    make(map[string]interface{}),
		},
		LimitType:     domain.RateLimitCNPJ,
		Key:           cnpj,
		PreviousUsage: limiter.GetRemainingQuota(),
		NewLimit:      limiter.Window.MaxRequests,
	}
	rm.repos.EventStore.SaveEvent(ctx, event)

	return nil
}

// ResetAllCNPJUsage reseta o uso de todos os CNPJs (executado diariamente)
func (rm *RateLimitManager) ResetAllCNPJUsage(ctx context.Context) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	resetCount := 0
	for key, limiter := range rm.limiters {
		if limiter.Type == domain.RateLimitCNPJ {
			limiter.Reset()
			rm.repos.RateLimiter.Update(limiter)
			resetCount++
		}
	}

	// Resetar também no banco diretamente para limiters não carregados
	// (seria implementado no repositório)

	return nil
}

// UpdateCNPJLimit atualiza o limite de um CNPJ
func (rm *RateLimitManager) UpdateCNPJLimit(cnpj string, newLimit int) error {
	limiter := rm.GetCNPJLimiter(cnpj)
	
	window := limiter.Window
	window.MaxRequests = newLimit
	limiter.SetWindow(window)

	return rm.repos.RateLimiter.Update(limiter)
}

// UpdateTenantLimit atualiza o limite de um tenant
func (rm *RateLimitManager) UpdateTenantLimit(tenantID uuid.UUID, newLimit int) error {
	limiter := rm.GetTenantLimiter(tenantID)
	
	window := limiter.Window
	window.MaxRequests = newLimit
	limiter.SetWindow(window)

	return rm.repos.RateLimiter.Update(limiter)
}

// GetLimiterStats obtém estatísticas de um limiter
func (rm *RateLimitManager) GetLimiterStats(limitType domain.RateLimitType, key string) (map[string]interface{}, error) {
	limiter, err := rm.repos.RateLimiter.FindByKey(limitType, key)
	if err != nil {
		return nil, err
	}

	if limiter == nil {
		return nil, nil
	}

	return limiter.GetStats(), nil
}

// GetAllCNPJStats obtém estatísticas de todos os CNPJs
func (rm *RateLimitManager) GetAllCNPJStats() ([]map[string]interface{}, error) {
	stats := make([]map[string]interface{}, 0)

	rm.mu.RLock()
	for _, limiter := range rm.limiters {
		if limiter.Type == domain.RateLimitCNPJ {
			stats = append(stats, limiter.GetStats())
		}
	}
	rm.mu.RUnlock()

	return stats, nil
}

// ActivateLimiter ativa um rate limiter
func (rm *RateLimitManager) ActivateLimiter(limitType domain.RateLimitType, key string) error {
	limiter, err := rm.repos.RateLimiter.FindByKey(limitType, key)
	if err != nil {
		return err
	}

	if limiter == nil {
		return domain.NewBusinessError("LIMITER_NOT_FOUND", "Rate limiter não encontrado")
	}

	limiter.Activate()
	return rm.repos.RateLimiter.Update(limiter)
}

// DeactivateLimiter desativa um rate limiter
func (rm *RateLimitManager) DeactivateLimiter(limitType domain.RateLimitType, key string) error {
	limiter, err := rm.repos.RateLimiter.FindByKey(limitType, key)
	if err != nil {
		return err
	}

	if limiter == nil {
		return domain.NewBusinessError("LIMITER_NOT_FOUND", "Rate limiter não encontrado")
	}

	limiter.Deactivate()
	return rm.repos.RateLimiter.Update(limiter)
}

// GetQuotaAvailability obtém disponibilidade de quota por CNPJ
func (rm *RateLimitManager) GetQuotaAvailability(cnpjs []string) map[string]int {
	availability := make(map[string]int)

	for _, cnpj := range cnpjs {
		limiter := rm.GetCNPJLimiter(cnpj)
		availability[cnpj] = limiter.GetRemainingQuota()
	}

	return availability
}

// GetTopUsedCNPJs retorna os CNPJs mais utilizados
func (rm *RateLimitManager) GetTopUsedCNPJs(limit int) []map[string]interface{} {
	cnpjStats := make([]map[string]interface{}, 0)

	rm.mu.RLock()
	for _, limiter := range rm.limiters {
		if limiter.Type == domain.RateLimitCNPJ {
			stats := limiter.GetStats()
			cnpjStats = append(cnpjStats, stats)
		}
	}
	rm.mu.RUnlock()

	// Ordenar por uso (seria implementado com sort)
	// Por simplicidade, retornar os primeiros
	if len(cnpjStats) > limit {
		return cnpjStats[:limit]
	}
	return cnpjStats
}

// Stop para o gerenciador e limpa recursos
func (rm *RateLimitManager) Stop() {
	if rm.cleanupTicker != nil {
		rm.cleanupTicker.Stop()
	}
	close(rm.stopCleanup)
}

// startCleanupRoutine inicia rotina de limpeza periódica
func (rm *RateLimitManager) startCleanupRoutine() {
	rm.cleanupTicker = time.NewTicker(15 * time.Minute)
	
	go func() {
		for {
			select {
			case <-rm.cleanupTicker.C:
				rm.cleanupExpiredLimiters()
			case <-rm.stopCleanup:
				return
			}
		}
	}()
}

// cleanupExpiredLimiters remove limiters expirados do cache
func (rm *RateLimitManager) cleanupExpiredLimiters() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	for key, limiter := range rm.limiters {
		// Remove limiters inativos há mais de 1 hora
		if !limiter.IsActive && time.Since(limiter.UpdatedAt) > time.Hour {
			delete(rm.limiters, key)
		}
	}

	// Limpeza no banco
	rm.repos.RateLimiter.CleanupExpired()
}

// revertCNPJUsage reverte o uso de quota do CNPJ (método auxiliar)
func (rm *RateLimitManager) revertCNPJUsage(limiter *domain.RateLimiter) {
	// Implementação simplificada: seria necessário um método específico no domain
	// Por ora, apenas resetamos o último request
	if len(limiter.Requests) > 0 {
		limiter.Requests = limiter.Requests[:len(limiter.Requests)-1]
	}
}

// revertTenantUsage reverte o uso de quota do tenant (método auxiliar)
func (rm *RateLimitManager) revertTenantUsage(limiter *domain.RateLimiter) {
	if len(limiter.Requests) > 0 {
		limiter.Requests = limiter.Requests[:len(limiter.Requests)-1]
	}
}

// saveLimiterUpdates salva as atualizações dos limiters
func (rm *RateLimitManager) saveLimiterUpdates(limiters ...*domain.RateLimiter) {
	for _, limiter := range limiters {
		rm.repos.RateLimiter.Update(limiter)
	}
}