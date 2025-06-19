package domain

import (
	"time"

	"github.com/google/uuid"
)

// SearchEvent representa eventos relacionados a busca
type SearchEvent interface {
	EventType() string
	AggregateID() string
	TenantID() string
	Timestamp() time.Time
}

// BaseSearchEvent estrutura base para eventos de busca
type BaseSearchEvent struct {
	EventTypeValue string    `json:"event_type"`
	AggregateIDValue string  `json:"aggregate_id"`
	TenantIDValue  string    `json:"tenant_id"`
	TimestampValue time.Time `json:"timestamp"`
}

func (e BaseSearchEvent) EventType() string    { return e.EventTypeValue }
func (e BaseSearchEvent) AggregateID() string  { return e.AggregateIDValue }
func (e BaseSearchEvent) TenantID() string     { return e.TenantIDValue }
func (e BaseSearchEvent) Timestamp() time.Time { return e.TimestampValue }

// SearchQueryExecuted evento disparado quando uma busca é executada
type SearchQueryExecuted struct {
	BaseSearchEvent
	Query           string                 `json:"query"`
	Filters         map[string]interface{} `json:"filters"`
	ResultCount     int64                  `json:"result_count"`
	ProcessingTime  time.Duration          `json:"processing_time"`
	CacheHit        bool                   `json:"cache_hit"`
	UserID          string                 `json:"user_id"`
}

// NewSearchQueryExecuted cria um novo evento de busca executada
func NewSearchQueryExecuted(queryID, tenantID, userID uuid.UUID, query string, filters map[string]interface{}, resultCount int64, processTime time.Duration, cacheHit bool) *SearchQueryExecuted {
	return &SearchQueryExecuted{
		BaseSearchEvent: BaseSearchEvent{
			EventTypeValue:   "search.query.executed",
			AggregateIDValue: queryID.String(),
			TenantIDValue:    tenantID.String(),
			TimestampValue:   time.Now(),
		},
		Query:          query,
		Filters:        filters,
		ResultCount:    resultCount,
		ProcessingTime: processTime,
		CacheHit:       cacheHit,
		UserID:         userID.String(),
	}
}

// SearchQueryFailed evento disparado quando uma busca falha
type SearchQueryFailed struct {
	BaseSearchEvent
	Query        string                 `json:"query"`
	Filters      map[string]interface{} `json:"filters"`
	ErrorMessage string                 `json:"error_message"`
	ErrorCode    string                 `json:"error_code"`
	UserID       string                 `json:"user_id"`
}

// NewSearchQueryFailed cria um novo evento de busca falha
func NewSearchQueryFailed(queryID, tenantID, userID uuid.UUID, query string, filters map[string]interface{}, errorMessage, errorCode string) *SearchQueryFailed {
	return &SearchQueryFailed{
		BaseSearchEvent: BaseSearchEvent{
			EventTypeValue:   "search.query.failed",
			AggregateIDValue: queryID.String(),
			TenantIDValue:    tenantID.String(),
			TimestampValue:   time.Now(),
		},
		Query:        query,
		Filters:      filters,
		ErrorMessage: errorMessage,
		ErrorCode:    errorCode,
		UserID:       userID.String(),
	}
}

// IndexCreated evento disparado quando um índice é criado
type IndexCreated struct {
	BaseSearchEvent
	IndexName   string                 `json:"index_name"`
	Description string                 `json:"description"`
	Mapping     map[string]interface{} `json:"mapping"`
	Settings    map[string]interface{} `json:"settings"`
	CreatedBy   string                 `json:"created_by"`
}

// NewIndexCreated cria um novo evento de índice criado
func NewIndexCreated(indexID, tenantID, createdBy uuid.UUID, indexName, description string, mapping, settings map[string]interface{}) *IndexCreated {
	return &IndexCreated{
		BaseSearchEvent: BaseSearchEvent{
			EventTypeValue:   "search.index.created",
			AggregateIDValue: indexID.String(),
			TenantIDValue:    tenantID.String(),
			TimestampValue:   time.Now(),
		},
		IndexName:   indexName,
		Description: description,
		Mapping:     mapping,
		Settings:    settings,
		CreatedBy:   createdBy.String(),
	}
}

// IndexDeleted evento disparado quando um índice é deletado
type IndexDeleted struct {
	BaseSearchEvent
	IndexName string `json:"index_name"`
	DeletedBy string `json:"deleted_by"`
	Reason    string `json:"reason"`
}

// NewIndexDeleted cria um novo evento de índice deletado
func NewIndexDeleted(indexID, tenantID, deletedBy uuid.UUID, indexName, reason string) *IndexDeleted {
	return &IndexDeleted{
		BaseSearchEvent: BaseSearchEvent{
			EventTypeValue:   "search.index.deleted",
			AggregateIDValue: indexID.String(),
			TenantIDValue:    tenantID.String(),
			TimestampValue:   time.Now(),
		},
		IndexName: indexName,
		DeletedBy: deletedBy.String(),
		Reason:    reason,
	}
}

// DocumentIndexed evento disparado quando um documento é indexado
type DocumentIndexed struct {
	BaseSearchEvent
	IndexName      string                 `json:"index_name"`
	DocumentID     string                 `json:"document_id"`
	DocumentType   string                 `json:"document_type"`
	DocumentData   map[string]interface{} `json:"document_data,omitempty"`
	ProcessingTime time.Duration          `json:"processing_time"`
	OperationID    string                 `json:"operation_id"`
}

// NewDocumentIndexed cria um novo evento de documento indexado
func NewDocumentIndexed(operationID, tenantID uuid.UUID, indexName, documentID, documentType string, documentData map[string]interface{}, processTime time.Duration) *DocumentIndexed {
	return &DocumentIndexed{
		BaseSearchEvent: BaseSearchEvent{
			EventTypeValue:   "search.document.indexed",
			AggregateIDValue: operationID.String(),
			TenantIDValue:    tenantID.String(),
			TimestampValue:   time.Now(),
		},
		IndexName:      indexName,
		DocumentID:     documentID,
		DocumentType:   documentType,
		DocumentData:   documentData,
		ProcessingTime: processTime,
		OperationID:    operationID.String(),
	}
}

// DocumentIndexingFailed evento disparado quando a indexação de um documento falha
type DocumentIndexingFailed struct {
	BaseSearchEvent
	IndexName      string                 `json:"index_name"`
	DocumentID     string                 `json:"document_id"`
	DocumentType   string                 `json:"document_type"`
	DocumentData   map[string]interface{} `json:"document_data,omitempty"`
	ErrorMessage   string                 `json:"error_message"`
	ErrorCode      string                 `json:"error_code"`
	ProcessingTime time.Duration          `json:"processing_time"`
	OperationID    string                 `json:"operation_id"`
}

// NewDocumentIndexingFailed cria um novo evento de falha na indexação
func NewDocumentIndexingFailed(operationID, tenantID uuid.UUID, indexName, documentID, documentType string, documentData map[string]interface{}, errorMessage, errorCode string, processTime time.Duration) *DocumentIndexingFailed {
	return &DocumentIndexingFailed{
		BaseSearchEvent: BaseSearchEvent{
			EventTypeValue:   "search.document.indexing_failed",
			AggregateIDValue: operationID.String(),
			TenantIDValue:    tenantID.String(),
			TimestampValue:   time.Now(),
		},
		IndexName:      indexName,
		DocumentID:     documentID,
		DocumentType:   documentType,
		DocumentData:   documentData,
		ErrorMessage:   errorMessage,
		ErrorCode:      errorCode,
		ProcessingTime: processTime,
		OperationID:    operationID.String(),
	}
}

// BulkOperationCompleted evento disparado quando uma operação em lote é concluída
type BulkOperationCompleted struct {
	BaseSearchEvent
	IndexName           string        `json:"index_name"`
	OperationType       string        `json:"operation_type"`
	DocumentsProcessed  int           `json:"documents_processed"`
	DocumentsSuccessful int           `json:"documents_successful"`
	DocumentsFailed     int           `json:"documents_failed"`
	ProcessingTime      time.Duration `json:"processing_time"`
	OperationID         string        `json:"operation_id"`
}

// NewBulkOperationCompleted cria um novo evento de operação em lote concluída
func NewBulkOperationCompleted(operationID, tenantID uuid.UUID, indexName, operationType string, processed, successful, failed int, processTime time.Duration) *BulkOperationCompleted {
	return &BulkOperationCompleted{
		BaseSearchEvent: BaseSearchEvent{
			EventTypeValue:   "search.bulk_operation.completed",
			AggregateIDValue: operationID.String(),
			TenantIDValue:    tenantID.String(),
			TimestampValue:   time.Now(),
		},
		IndexName:           indexName,
		OperationType:       operationType,
		DocumentsProcessed:  processed,
		DocumentsSuccessful: successful,
		DocumentsFailed:     failed,
		ProcessingTime:      processTime,
		OperationID:         operationID.String(),
	}
}

// CacheHit evento disparado quando há um hit no cache
type CacheHit struct {
	BaseSearchEvent
	CacheKey   string `json:"cache_key"`
	QueryHash  string `json:"query_hash"`
	HitCount   int    `json:"hit_count"`
	UserID     string `json:"user_id"`
}

// NewCacheHit cria um novo evento de cache hit
func NewCacheHit(cacheID, tenantID, userID uuid.UUID, cacheKey, queryHash string, hitCount int) *CacheHit {
	return &CacheHit{
		BaseSearchEvent: BaseSearchEvent{
			EventTypeValue:   "search.cache.hit",
			AggregateIDValue: cacheID.String(),
			TenantIDValue:    tenantID.String(),
			TimestampValue:   time.Now(),
		},
		CacheKey:  cacheKey,
		QueryHash: queryHash,
		HitCount:  hitCount,
		UserID:    userID.String(),
	}
}

// CacheMiss evento disparado quando há um miss no cache
type CacheMiss struct {
	BaseSearchEvent
	CacheKey  string `json:"cache_key"`
	QueryHash string `json:"query_hash"`
	UserID    string `json:"user_id"`
}

// NewCacheMiss cria um novo evento de cache miss
func NewCacheMiss(tenantID, userID uuid.UUID, cacheKey, queryHash string) *CacheMiss {
	return &CacheMiss{
		BaseSearchEvent: BaseSearchEvent{
			EventTypeValue:   "search.cache.miss",
			AggregateIDValue: uuid.New().String(), // Gerar ID único para o evento
			TenantIDValue:    tenantID.String(),
			TimestampValue:   time.Now(),
		},
		CacheKey:  cacheKey,
		QueryHash: queryHash,
		UserID:    userID.String(),
	}
}