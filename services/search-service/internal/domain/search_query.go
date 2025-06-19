package domain

import (
	"time"

	"github.com/google/uuid"
)

// SearchQuery representa uma consulta de busca
type SearchQuery struct {
	ID           uuid.UUID
	Query        string
	Filters      map[string]interface{}
	Sort         []SortField
	Page         int
	Size         int
	IncludeCount bool
	Highlight    bool
	TenantID     uuid.UUID
	UserID       uuid.UUID
	CreatedAt    time.Time
}

// SortField representa um campo de ordenação
type SortField struct {
	Field string
	Order SortOrder
}

// SortOrder representa a direção da ordenação
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// SearchResult representa um resultado de busca
type SearchResult struct {
	ID          uuid.UUID
	QueryID     uuid.UUID
	Results     []SearchHit
	Total       int64
	Page        int
	Size        int
	ProcessTime time.Duration
	CacheHit    bool
	TenantID    uuid.UUID
	CreatedAt   time.Time
}

// SearchHit representa um item nos resultados de busca
type SearchHit struct {
	Index      string
	DocumentID string
	Score      float64
	Source     map[string]interface{}
	Highlights map[string][]string
}

// SearchIndex representa um índice de busca
type SearchIndex struct {
	ID            uuid.UUID
	Name          string
	Description   string
	Mapping       map[string]interface{}
	Settings      map[string]interface{}
	Aliases       []string
	IsActive      bool
	DocumentCount int64
	SizeInBytes   int64
	LastIndexedAt *time.Time
	TenantID      uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CreatedBy     uuid.UUID
	UpdatedBy     uuid.UUID
}

// IndexingOperation representa uma operação de indexação
type IndexingOperation struct {
	ID                uuid.UUID
	IndexName         string
	Operation         OperationType
	DocumentID        string
	DocumentType      string
	Status            OperationStatus
	ErrorMessage      string
	ProcessingTime    time.Duration
	DocumentsProcessed int
	DocumentData      map[string]interface{}
	TenantID          uuid.UUID
	CreatedAt         time.Time
	ProcessedAt       *time.Time
}

// OperationType representa o tipo de operação de indexação
type OperationType string

const (
	OperationIndex  OperationType = "index"
	OperationUpdate OperationType = "update"
	OperationDelete OperationType = "delete"
	OperationBulk   OperationType = "bulk"
)

// OperationStatus representa o status de uma operação
type OperationStatus string

const (
	StatusPending OperationStatus = "pending"
	StatusSuccess OperationStatus = "success"
	StatusError   OperationStatus = "error"
)

// SearchStatistics representa estatísticas de busca
type SearchStatistics struct {
	ID                  uuid.UUID
	Date                time.Time
	Hour                int
	TotalSearches       int64
	SuccessfulSearches  int64
	FailedSearches      int64
	AvgResponseTime     float64
	BasicSearches       int64
	AdvancedSearches    int64
	SuggestionsRequests int64
	IndexName           string
	TenantID            uuid.UUID
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// SearchCache representa cache de busca
type SearchCache struct {
	ID             uuid.UUID
	CacheKey       string
	QueryHash      string
	SearchQuery    map[string]interface{}
	SearchResults  map[string]interface{}
	TotalResults   int
	HitCount       int
	LastAccessedAt time.Time
	ExpiresAt      time.Time
	TenantID       uuid.UUID
	UserID         uuid.UUID
	CreatedAt      time.Time
}

// NewSearchQuery cria uma nova consulta de busca
func NewSearchQuery(query string, tenantID, userID uuid.UUID) *SearchQuery {
	return &SearchQuery{
		ID:           uuid.New(),
		Query:        query,
		Filters:      make(map[string]interface{}),
		Sort:         []SortField{},
		Page:         1,
		Size:         20,
		IncludeCount: true,
		Highlight:    true,
		TenantID:     tenantID,
		UserID:       userID,
		CreatedAt:    time.Now(),
	}
}

// NewSearchIndex cria um novo índice de busca
func NewSearchIndex(name, description string, tenantID, createdBy uuid.UUID) *SearchIndex {
	now := time.Now()
	return &SearchIndex{
		ID:            uuid.New(),
		Name:          name,
		Description:   description,
		Mapping:       make(map[string]interface{}),
		Settings:      make(map[string]interface{}),
		Aliases:       []string{},
		IsActive:      true,
		DocumentCount: 0,
		SizeInBytes:   0,
		TenantID:      tenantID,
		CreatedAt:     now,
		UpdatedAt:     now,
		CreatedBy:     createdBy,
		UpdatedBy:     createdBy,
	}
}

// NewIndexingOperation cria uma nova operação de indexação
func NewIndexingOperation(indexName string, operation OperationType, tenantID uuid.UUID) *IndexingOperation {
	return &IndexingOperation{
		ID:                 uuid.New(),
		IndexName:          indexName,
		Operation:          operation,
		Status:             StatusPending,
		DocumentsProcessed: 0,
		DocumentData:       make(map[string]interface{}),
		TenantID:           tenantID,
		CreatedAt:          time.Now(),
	}
}

// IsValid valida se a consulta de busca é válida
func (sq *SearchQuery) IsValid() bool {
	return sq.Query != "" && sq.Size > 0 && sq.Page > 0
}

// GetOffset calcula o offset baseado na página
func (sq *SearchQuery) GetOffset() int {
	return (sq.Page - 1) * sq.Size
}

// AddFilter adiciona um filtro à consulta
func (sq *SearchQuery) AddFilter(key string, value interface{}) {
	if sq.Filters == nil {
		sq.Filters = make(map[string]interface{})
	}
	sq.Filters[key] = value
}

// AddSort adiciona um campo de ordenação
func (sq *SearchQuery) AddSort(field string, order SortOrder) {
	sq.Sort = append(sq.Sort, SortField{
		Field: field,
		Order: order,
	})
}

// Complete marca uma operação como completa
func (io *IndexingOperation) Complete(processTime time.Duration, documentsProcessed int) {
	io.Status = StatusSuccess
	io.ProcessingTime = processTime
	io.DocumentsProcessed = documentsProcessed
	now := time.Now()
	io.ProcessedAt = &now
}

// Fail marca uma operação como falha
func (io *IndexingOperation) Fail(errorMessage string, processTime time.Duration) {
	io.Status = StatusError
	io.ErrorMessage = errorMessage
	io.ProcessingTime = processTime
	now := time.Now()
	io.ProcessedAt = &now
}

// IsExpired verifica se o cache expirou
func (sc *SearchCache) IsExpired() bool {
	return time.Now().After(sc.ExpiresAt)
}

// Hit incrementa o contador de hits do cache
func (sc *SearchCache) Hit() {
	sc.HitCount++
	sc.LastAccessedAt = time.Now()
}