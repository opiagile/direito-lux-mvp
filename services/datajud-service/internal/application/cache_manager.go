package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/direito-lux/datajud-service/internal/domain"
)

// CacheManager gerencia cache de consultas DataJud
type CacheManager struct {
	cache         domain.Cache
	repos         *domain.Repositories
	config        domain.DataJudConfig
	stats         *CacheStatistics
	cleanupTicker *time.Ticker
	stopCleanup   chan bool
	mu            sync.RWMutex
}

// CacheStatistics estatísticas detalhadas do cache
type CacheStatistics struct {
	TotalHits     int64     `json:"total_hits"`
	TotalMisses   int64     `json:"total_misses"`
	TotalSets     int64     `json:"total_sets"`
	TotalDeletes  int64     `json:"total_deletes"`
	TotalEvictions int64    `json:"total_evictions"`
	StartTime     time.Time `json:"start_time"`
	LastReset     time.Time `json:"last_reset"`
	mu            sync.RWMutex
}

// NewCacheManager cria novo gerenciador de cache
func NewCacheManager(cache domain.Cache, repos *domain.Repositories, config domain.DataJudConfig) *CacheManager {
	manager := &CacheManager{
		cache:       cache,
		repos:       repos,
		config:      config,
		stopCleanup: make(chan bool),
		stats: &CacheStatistics{
			StartTime: time.Now(),
			LastReset: time.Now(),
		},
	}

	// Iniciar limpeza periódica
	manager.startCleanupRoutine()

	return manager
}

// Get obtém entrada do cache
func (cm *CacheManager) Get(key string) (*domain.CacheEntry, error) {
	entry, err := cm.cache.Get(key)
	
	cm.stats.mu.Lock()
	if err != nil || entry == nil {
		cm.stats.TotalMisses++
	} else {
		cm.stats.TotalHits++
	}
	cm.stats.mu.Unlock()

	return entry, err
}

// Set armazena entrada no cache
func (cm *CacheManager) Set(key string, value interface{}, ttl int, tenantID uuid.UUID, requestType domain.RequestType) error {
	if ttl <= 0 {
		ttl = cm.config.DefaultCacheTTL
	}

	// Verificar limites antes de adicionar
	if err := cm.enforceCapacityLimits(); err != nil {
		return fmt.Errorf("erro ao verificar capacidade: %w", err)
	}

	entry, err := domain.NewCacheEntry(key, value, ttl, tenantID, requestType)
	if err != nil {
		return fmt.Errorf("erro ao criar entrada de cache: %w", err)
	}

	// Definir metadados específicos baseado no tipo
	cm.setCacheMetadata(entry, value)

	err = cm.cache.Set(entry)
	
	cm.stats.mu.Lock()
	if err == nil {
		cm.stats.TotalSets++
	}
	cm.stats.mu.Unlock()

	// Publicar evento de cache stored
	if err == nil {
		event := &domain.DataJudCacheStored{
			BaseEvent: domain.BaseEvent{
				ID:          uuid.New(),
				Type:        "datajud.cache.stored",
				AggregateID: tenantID,
				OccurredAt:  time.Now(),
				Version:     1,
				Metadata:    make(map[string]interface{}),
			},
			CacheKey:    key,
			RequestType: requestType,
			TTL:         ttl,
			Size:        entry.Size,
		}
		cm.repos.EventStore.SaveEvent(context.Background(), event)
	}

	return err
}

// Delete remove entrada do cache
func (cm *CacheManager) Delete(key string) error {
	err := cm.cache.Delete(key)
	
	cm.stats.mu.Lock()
	if err == nil {
		cm.stats.TotalDeletes++
	}
	cm.stats.mu.Unlock()

	return err
}

// Exists verifica se chave existe no cache
func (cm *CacheManager) Exists(key string) bool {
	return cm.cache.Exists(key)
}

// Clear limpa todo o cache
func (cm *CacheManager) Clear() error {
	return cm.cache.Clear()
}

// GetByProcessNumber obtém entradas de cache por número de processo
func (cm *CacheManager) GetByProcessNumber(processNumber string) ([]*domain.CacheEntry, error) {
	return cm.repos.Cache.FindByRequestType(domain.RequestTypeProcess)
}

// GetByTenant obtém entradas de cache por tenant
func (cm *CacheManager) GetByTenant(tenantID uuid.UUID) ([]*domain.CacheEntry, error) {
	return cm.repos.Cache.FindByTenantID(tenantID)
}

// GetByRequestType obtém entradas por tipo de requisição
func (cm *CacheManager) GetByRequestType(requestType domain.RequestType) ([]*domain.CacheEntry, error) {
	return cm.repos.Cache.FindByRequestType(requestType)
}

// GetStats obtém estatísticas do cache
func (cm *CacheManager) GetStats() (*domain.CacheStats, error) {
	// Estatísticas do cache base
	baseStats, err := cm.cache.GetStats()
	if err != nil {
		return nil, err
	}

	// Adicionar estatísticas do manager
	cm.stats.mu.RLock()
	baseStats.HitCount = cm.stats.TotalHits
	baseStats.MissCount = cm.stats.TotalMisses
	baseStats.EvictionCount = cm.stats.TotalEvictions
	
	if cm.stats.TotalHits+cm.stats.TotalMisses > 0 {
		baseStats.HitRatio = float64(cm.stats.TotalHits) / float64(cm.stats.TotalHits+cm.stats.TotalMisses) * 100
	}
	cm.stats.mu.RUnlock()

	return baseStats, nil
}

// GetDetailedStats obtém estatísticas detalhadas
func (cm *CacheManager) GetDetailedStats() map[string]interface{} {
	stats, _ := cm.GetStats()
	
	cm.stats.mu.RLock()
	detailed := map[string]interface{}{
		"base_stats":      stats,
		"total_hits":      cm.stats.TotalHits,
		"total_misses":    cm.stats.TotalMisses,
		"total_sets":      cm.stats.TotalSets,
		"total_deletes":   cm.stats.TotalDeletes,
		"total_evictions": cm.stats.TotalEvictions,
		"uptime":          time.Since(cm.stats.StartTime).String(),
		"last_reset":      cm.stats.LastReset,
		"hit_ratio":       cm.getHitRatio(),
		"miss_ratio":      cm.getMissRatio(),
		"operations_per_second": cm.getOperationsPerSecond(),
	}
	cm.stats.mu.RUnlock()

	return detailed
}

// GetStatsByTenant obtém estatísticas por tenant
func (cm *CacheManager) GetStatsByTenant(tenantID uuid.UUID) (map[string]interface{}, error) {
	entries, err := cm.GetByTenant(tenantID)
	if err != nil {
		return nil, err
	}

	var totalSize int64
	var totalHits int64
	requestTypes := make(map[domain.RequestType]int)

	for _, entry := range entries {
		totalSize += entry.Size
		totalHits += entry.HitCount
		requestTypes[entry.RequestType]++
	}

	return map[string]interface{}{
		"tenant_id":       tenantID,
		"total_entries":   len(entries),
		"total_size":      totalSize,
		"total_hits":      totalHits,
		"by_request_type": requestTypes,
	}, nil
}

// GetStatsByRequestType obtém estatísticas por tipo de requisição
func (cm *CacheManager) GetStatsByRequestType() (map[domain.RequestType]map[string]interface{}, error) {
	stats := make(map[domain.RequestType]map[string]interface{})
	
	requestTypes := []domain.RequestType{
		domain.RequestTypeProcess,
		domain.RequestTypeMovement,
		domain.RequestTypeParty,
		domain.RequestTypeDocument,
		domain.RequestTypeBulk,
	}

	for _, reqType := range requestTypes {
		entries, err := cm.GetByRequestType(reqType)
		if err != nil {
			continue
		}

		var totalSize int64
		var totalHits int64
		var avgAge time.Duration

		for _, entry := range entries {
			totalSize += entry.Size
			totalHits += entry.HitCount
			avgAge += entry.GetAge()
		}

		if len(entries) > 0 {
			avgAge = avgAge / time.Duration(len(entries))
		}

		stats[reqType] = map[string]interface{}{
			"total_entries": len(entries),
			"total_size":    totalSize,
			"total_hits":    totalHits,
			"average_age":   avgAge.String(),
		}
	}

	return stats, nil
}

// CleanupExpired remove entradas expiradas
func (cm *CacheManager) CleanupExpired() (int, error) {
	count, err := cm.cache.CleanupExpired()
	
	cm.stats.mu.Lock()
	cm.stats.TotalEvictions += int64(count)
	cm.stats.mu.Unlock()

	return count, err
}

// CleanupByTenant remove entradas de um tenant específico
func (cm *CacheManager) CleanupByTenant(ctx context.Context, tenantID uuid.UUID) (int, error) {
	entries, err := cm.GetByTenant(tenantID)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, entry := range entries {
		if err := cm.Delete(entry.Key); err == nil {
			count++
		}
	}

	// Publicar evento de limpeza
	if count > 0 {
		event := &domain.DataJudCacheEvicted{
			BaseEvent: domain.BaseEvent{
				ID:          uuid.New(),
				Type:        "datajud.cache.tenant_cleanup",
				AggregateID: tenantID,
				OccurredAt:  time.Now(),
				Version:     1,
				Metadata: map[string]interface{}{
					"entries_removed": count,
				},
			},
			CacheKey:    fmt.Sprintf("tenant:%s", tenantID.String()),
			RequestType: "cleanup",
			Reason:      "tenant_cleanup",
		}
		cm.repos.EventStore.SaveEvent(ctx, event)
	}

	return count, nil
}

// CleanupByAge remove entradas mais antigas que a idade especificada
func (cm *CacheManager) CleanupByAge(maxAge time.Duration) (int, error) {
	// Esta funcionalidade seria implementada no repositório
	// Por simplicidade, usar cleanup de expirados
	return cm.CleanupExpired()
}

// WarmupCache pré-carrega cache com dados comuns
func (cm *CacheManager) WarmupCache(ctx context.Context, tenantID uuid.UUID, processNumbers []string) error {
	// Implementação de warmup seria fazer requisições em background
	// Por simplicidade, apenas log
	return nil
}

// ResetStats reseta estatísticas do cache
func (cm *CacheManager) ResetStats() {
	cm.stats.mu.Lock()
	cm.stats.TotalHits = 0
	cm.stats.TotalMisses = 0
	cm.stats.TotalSets = 0
	cm.stats.TotalDeletes = 0
	cm.stats.TotalEvictions = 0
	cm.stats.LastReset = time.Now()
	cm.stats.mu.Unlock()
}

// Stop para o gerenciador e limpa recursos
func (cm *CacheManager) Stop() {
	if cm.cleanupTicker != nil {
		cm.cleanupTicker.Stop()
	}
	close(cm.stopCleanup)
}

// GetTopHitEntries retorna entradas com mais hits
func (cm *CacheManager) GetTopHitEntries(limit int) ([]*domain.CacheEntry, error) {
	// Esta funcionalidade seria implementada no repositório com ordenação
	// Por simplicidade, retornar entradas de processo
	return cm.GetByRequestType(domain.RequestTypeProcess)
}

// GetLargestEntries retorna maiores entradas por tamanho
func (cm *CacheManager) GetLargestEntries(limit int) ([]*domain.CacheEntry, error) {
	// Implementação seria no repositório com ordenação por size
	return cm.GetByRequestType(domain.RequestTypeDocument)
}

// OptimizeCache executa otimizações no cache
func (cm *CacheManager) OptimizeCache() error {
	// 1. Limpar expirados
	expired, _ := cm.CleanupExpired()
	
	// 2. Verificar capacidade
	if err := cm.enforceCapacityLimits(); err != nil {
		return err
	}

	// 3. Log otimização
	_ = fmt.Sprintf("Cache optimized: %d expired entries removed", expired)

	return nil
}

// startCleanupRoutine inicia rotina de limpeza periódica
func (cm *CacheManager) startCleanupRoutine() {
	cm.cleanupTicker = time.NewTicker(cm.config.CacheCleanupInterval)
	
	go func() {
		for {
			select {
			case <-cm.cleanupTicker.C:
				cm.performCleanup()
			case <-cm.stopCleanup:
				return
			}
		}
	}()
}

// performCleanup executa limpeza periódica
func (cm *CacheManager) performCleanup() {
	// Limpar entradas expiradas
	expired, err := cm.CleanupExpired()
	if err != nil {
		return
	}

	// Verificar limites de capacidade
	cm.enforceCapacityLimits()

	// Log se removeu muitas entradas
	if expired > 100 {
		_ = fmt.Sprintf("Cache cleanup: %d expired entries removed", expired)
	}
}

// enforceCapacityLimits aplica limites de capacidade
func (cm *CacheManager) enforceCapacityLimits() error {
	stats, err := cm.cache.GetStats()
	if err != nil {
		return err
	}

	// Verificar limite de tamanho
	if cm.config.MaxCacheSize > 0 && stats.TotalSize > cm.config.MaxCacheSize {
		// Implementar LRU eviction
		// Por simplicidade, apenas limpar expirados
		cm.CleanupExpired()
	}

	// Verificar limite de entradas
	if cm.config.MaxCacheEntries > 0 && stats.TotalEntries > cm.config.MaxCacheEntries {
		// Implementar LRU eviction
		cm.CleanupExpired()
	}

	return nil
}

// setCacheMetadata define metadados específicos da entrada
func (cm *CacheManager) setCacheMetadata(entry *domain.CacheEntry, value interface{}) {
	switch response := value.(type) {
	case *domain.DataJudResponse:
		if response.ProcessData != nil {
			entry.SetProcessNumber(response.ProcessData.Number)
		}
	case *domain.ProcessResponseData:
		entry.SetProcessNumber(response.Number)
	}
}

// getHitRatio calcula taxa de acerto
func (cm *CacheManager) getHitRatio() float64 {
	total := cm.stats.TotalHits + cm.stats.TotalMisses
	if total == 0 {
		return 0
	}
	return float64(cm.stats.TotalHits) / float64(total) * 100
}

// getMissRatio calcula taxa de erro
func (cm *CacheManager) getMissRatio() float64 {
	total := cm.stats.TotalHits + cm.stats.TotalMisses
	if total == 0 {
		return 0
	}
	return float64(cm.stats.TotalMisses) / float64(total) * 100
}

// getOperationsPerSecond calcula operações por segundo
func (cm *CacheManager) getOperationsPerSecond() float64 {
	uptime := time.Since(cm.stats.StartTime)
	if uptime == 0 {
		return 0
	}
	
	totalOps := cm.stats.TotalHits + cm.stats.TotalMisses + cm.stats.TotalSets + cm.stats.TotalDeletes
	return float64(totalOps) / uptime.Seconds()
}