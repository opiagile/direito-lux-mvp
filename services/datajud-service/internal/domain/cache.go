package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// CacheEntry representa uma entrada no cache
type CacheEntry struct {
	ID          uuid.UUID   `json:"id"`
	Key         string      `json:"key"`
	Value       interface{} `json:"value"`
	TTL         int         `json:"ttl"`        // TTL em segundos
	ExpiresAt   time.Time   `json:"expires_at"`
	Size        int64       `json:"size"`       // Tamanho em bytes
	HitCount    int64       `json:"hit_count"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	LastHitAt   *time.Time  `json:"last_hit_at"`
	
	// Metadados específicos do DataJud
	RequestType   RequestType `json:"request_type"`
	ProcessNumber string      `json:"process_number,omitempty"`
	CourtID       string      `json:"court_id,omitempty"`
	TenantID      uuid.UUID   `json:"tenant_id"`
}

// CacheStats estatísticas do cache
type CacheStats struct {
	TotalEntries     int64   `json:"total_entries"`
	TotalSize        int64   `json:"total_size"`
	HitCount         int64   `json:"hit_count"`
	MissCount        int64   `json:"miss_count"`
	HitRatio         float64 `json:"hit_ratio"`
	EvictionCount    int64   `json:"eviction_count"`
	ExpiredCount     int64   `json:"expired_count"`
	AverageEntrySize int64   `json:"average_entry_size"`
	OldestEntry      *time.Time `json:"oldest_entry"`
	NewestEntry      *time.Time `json:"newest_entry"`
}

// CacheRepository interface para persistência do cache
type CacheRepository interface {
	Get(key string) (*CacheEntry, error)
	Set(entry *CacheEntry) error
	Delete(key string) error
	Exists(key string) bool
	FindByTenantID(tenantID uuid.UUID) ([]*CacheEntry, error)
	FindByRequestType(requestType RequestType) ([]*CacheEntry, error)
	FindExpired() ([]*CacheEntry, error)
	CleanupExpired() (int, error)
	GetStats() (*CacheStats, error)
	Clear() error
	GetSize() (int64, error)
}

// Cache interface para operações de cache
type Cache interface {
	Get(key string) (*CacheEntry, error)
	Set(key string, value interface{}, ttl int) error
	Delete(key string) error
	Exists(key string) bool
	Clear() error
	GetStats() (*CacheStats, error)
}

// NewCacheEntry cria uma nova entrada de cache
func NewCacheEntry(key string, value interface{}, ttl int, tenantID uuid.UUID, requestType RequestType) (*CacheEntry, error) {
	now := time.Now()
	
	// Serializa o valor para calcular o tamanho
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	
	return &CacheEntry{
		ID:          uuid.New(),
		Key:         key,
		Value:       value,
		TTL:         ttl,
		ExpiresAt:   now.Add(time.Duration(ttl) * time.Second),
		Size:        int64(len(jsonValue)),
		HitCount:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
		RequestType: requestType,
		TenantID:    tenantID,
	}, nil
}

// IsExpired verifica se a entrada expirou
func (c *CacheEntry) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// Hit registra um acesso à entrada
func (c *CacheEntry) Hit() {
	c.HitCount++
	now := time.Now()
	c.LastHitAt = &now
	c.UpdatedAt = now
}

// Extend estende o TTL da entrada
func (c *CacheEntry) Extend(additionalTTL int) {
	c.TTL += additionalTTL
	c.ExpiresAt = c.ExpiresAt.Add(time.Duration(additionalTTL) * time.Second)
	c.UpdatedAt = time.Now()
}

// Refresh atualiza o valor e TTL
func (c *CacheEntry) Refresh(value interface{}, ttl int) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	c.Value = value
	c.TTL = ttl
	c.Size = int64(len(jsonValue))
	now := time.Now()
	c.ExpiresAt = now.Add(time.Duration(ttl) * time.Second)
	c.UpdatedAt = now
	
	return nil
}

// GetAge retorna a idade da entrada
func (c *CacheEntry) GetAge() time.Duration {
	return time.Since(c.CreatedAt)
}

// GetTimeToExpire retorna tempo até expirar
func (c *CacheEntry) GetTimeToExpire() time.Duration {
	if c.IsExpired() {
		return 0
	}
	return c.ExpiresAt.Sub(time.Now())
}

// SetProcessNumber define o número do processo para cache de processo
func (c *CacheEntry) SetProcessNumber(processNumber string) {
	c.ProcessNumber = processNumber
	c.UpdatedAt = time.Now()
}

// SetCourtID define o tribunal para filtros
func (c *CacheEntry) SetCourtID(courtID string) {
	c.CourtID = courtID
	c.UpdatedAt = time.Now()
}

// GetMetadata retorna metadados da entrada
func (c *CacheEntry) GetMetadata() map[string]interface{} {
	return map[string]interface{}{
		"id":             c.ID,
		"key":            c.Key,
		"ttl":            c.TTL,
		"expires_at":     c.ExpiresAt,
		"size":           c.Size,
		"hit_count":      c.HitCount,
		"age":            c.GetAge().String(),
		"time_to_expire": c.GetTimeToExpire().String(),
		"is_expired":     c.IsExpired(),
		"request_type":   c.RequestType,
		"process_number": c.ProcessNumber,
		"court_id":       c.CourtID,
		"tenant_id":      c.TenantID,
		"created_at":     c.CreatedAt,
		"updated_at":     c.UpdatedAt,
		"last_hit_at":    c.LastHitAt,
	}
}

// CacheManager gerencia o cache com políticas avançadas
type CacheManager struct {
	cache           Cache
	maxSize         int64               // Tamanho máximo em bytes
	maxEntries      int64               // Número máximo de entradas
	defaultTTL      int                 // TTL padrão em segundos
	cleanupInterval time.Duration       // Intervalo de limpeza
	hitCount        int64               // Contador de hits
	missCount       int64               // Contador de misses
	evictionCount   int64               // Contador de evictions
}

// NewCacheManager cria um novo gerenciador de cache
func NewCacheManager(cache Cache, maxSize, maxEntries int64, defaultTTL int, cleanupInterval time.Duration) *CacheManager {
	return &CacheManager{
		cache:           cache,
		maxSize:         maxSize,
		maxEntries:      maxEntries,
		defaultTTL:      defaultTTL,
		cleanupInterval: cleanupInterval,
	}
}

// Get obtém uma entrada do cache
func (cm *CacheManager) Get(key string) (*CacheEntry, error) {
	entry, err := cm.cache.Get(key)
	if err != nil {
		cm.missCount++
		return nil, err
	}
	
	if entry == nil {
		cm.missCount++
		return nil, nil
	}
	
	if entry.IsExpired() {
		cm.cache.Delete(key)
		cm.missCount++
		return nil, nil
	}
	
	entry.Hit()
	cm.hitCount++
	return entry, nil
}

// Set armazena uma entrada no cache
func (cm *CacheManager) Set(key string, value interface{}, ttl int, tenantID uuid.UUID, requestType RequestType) error {
	if ttl <= 0 {
		ttl = cm.defaultTTL
	}
	
	entry, err := NewCacheEntry(key, value, ttl, tenantID, requestType)
	if err != nil {
		return err
	}
	
	// Verifica se precisa fazer limpeza por tamanho/quantidade
	if err := cm.enforceCapacityLimits(); err != nil {
		return err
	}
	
	return cm.cache.Set(entry.Key, entry.Value, entry.TTL)
}

// Delete remove uma entrada do cache
func (cm *CacheManager) Delete(key string) error {
	return cm.cache.Delete(key)
}

// GetStats retorna estatísticas do cache
func (cm *CacheManager) GetStats() (*CacheStats, error) {
	stats, err := cm.cache.GetStats()
	if err != nil {
		return nil, err
	}
	
	// Adiciona estatísticas do manager
	stats.HitCount = cm.hitCount
	stats.MissCount = cm.missCount
	stats.EvictionCount = cm.evictionCount
	
	if cm.hitCount+cm.missCount > 0 {
		stats.HitRatio = float64(cm.hitCount) / float64(cm.hitCount+cm.missCount) * 100
	}
	
	return stats, nil
}

// CleanupExpired remove entradas expiradas
func (cm *CacheManager) CleanupExpired() (int, error) {
	// Cache em memória não tem método CleanupExpired
	// Em implementação real (Redis), isso seria executado via repository
	count := 0
	
	cm.evictionCount += int64(count)
	return count, nil
}

// enforceCapacityLimits aplica limites de capacidade
func (cm *CacheManager) enforceCapacityLimits() error {
	stats, err := cm.cache.GetStats()
	if err != nil {
		return err
	}
	
	// Limpa entradas expiradas primeiro
	if expired, err := cm.CleanupExpired(); err == nil && expired > 0 {
		// Recalcula stats após limpeza
		stats, _ = cm.cache.GetStats()
	}
	
	// Se ainda excede limites, remove entradas menos usadas
	if (cm.maxSize > 0 && stats.TotalSize > cm.maxSize) ||
	   (cm.maxEntries > 0 && stats.TotalEntries > cm.maxEntries) {
		// Implementação de LRU seria feita aqui
		// Por simplicidade, apenas limpamos entradas expiradas
		return nil
	}
	
	return nil
}

// GetHitRatio retorna a taxa de acerto do cache
func (cm *CacheManager) GetHitRatio() float64 {
	total := cm.hitCount + cm.missCount
	if total == 0 {
		return 0
	}
	return float64(cm.hitCount) / float64(total) * 100
}