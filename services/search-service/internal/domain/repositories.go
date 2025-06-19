package domain

import (
	"context"

	"github.com/google/uuid"
)

// SearchQueryRepository repositório para consultas de busca
type SearchQueryRepository interface {
	Save(ctx context.Context, query *SearchQuery) error
	FindByID(ctx context.Context, id uuid.UUID) (*SearchQuery, error)
	FindByTenant(ctx context.Context, tenantID uuid.UUID, page, size int) ([]*SearchQuery, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// SearchResultRepository repositório para resultados de busca
type SearchResultRepository interface {
	Save(ctx context.Context, result *SearchResult) error
	FindByQueryID(ctx context.Context, queryID uuid.UUID) (*SearchResult, error)
	FindByTenant(ctx context.Context, tenantID uuid.UUID, page, size int) ([]*SearchResult, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// SearchIndexRepository repositório para índices de busca
type SearchIndexRepository interface {
	Save(ctx context.Context, index *SearchIndex) error
	FindByID(ctx context.Context, id uuid.UUID) (*SearchIndex, error)
	FindByName(ctx context.Context, name string) (*SearchIndex, error)
	FindByTenant(ctx context.Context, tenantID uuid.UUID) ([]*SearchIndex, error)
	FindActive(ctx context.Context) ([]*SearchIndex, error)
	Update(ctx context.Context, index *SearchIndex) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// IndexingOperationRepository repositório para operações de indexação
type IndexingOperationRepository interface {
	Save(ctx context.Context, operation *IndexingOperation) error
	FindByID(ctx context.Context, id uuid.UUID) (*IndexingOperation, error)
	FindByIndex(ctx context.Context, indexName string, page, size int) ([]*IndexingOperation, error)
	FindByStatus(ctx context.Context, status OperationStatus, page, size int) ([]*IndexingOperation, error)
	FindByTenant(ctx context.Context, tenantID uuid.UUID, page, size int) ([]*IndexingOperation, error)
	Update(ctx context.Context, operation *IndexingOperation) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByStatus(ctx context.Context, status OperationStatus) (int64, error)
}

// SearchStatisticsRepository repositório para estatísticas de busca
type SearchStatisticsRepository interface {
	Save(ctx context.Context, stats *SearchStatistics) error
	FindByID(ctx context.Context, id uuid.UUID) (*SearchStatistics, error)
	FindByDateRange(ctx context.Context, tenantID uuid.UUID, startDate, endDate string) ([]*SearchStatistics, error)
	FindByIndex(ctx context.Context, indexName string, startDate, endDate string) ([]*SearchStatistics, error)
	Update(ctx context.Context, stats *SearchStatistics) error
	Delete(ctx context.Context, id uuid.UUID) error
	AggregateByTenant(ctx context.Context, tenantID uuid.UUID, startDate, endDate string) (*SearchStatisticsAggregate, error)
	AggregateByIndex(ctx context.Context, indexName string, startDate, endDate string) (*SearchStatisticsAggregate, error)
}

// SearchCacheRepository repositório para cache de busca
type SearchCacheRepository interface {
	Save(ctx context.Context, cache *SearchCache) error
	FindByKey(ctx context.Context, cacheKey string) (*SearchCache, error)
	FindByQueryHash(ctx context.Context, queryHash string) (*SearchCache, error)
	FindByTenant(ctx context.Context, tenantID uuid.UUID, page, size int) ([]*SearchCache, error)
	Update(ctx context.Context, cache *SearchCache) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByKey(ctx context.Context, cacheKey string) error
	DeleteExpired(ctx context.Context) (int64, error)
	CountByTenant(ctx context.Context, tenantID uuid.UUID) (int64, error)
}

// ElasticsearchRepository repositório para operações do Elasticsearch
type ElasticsearchRepository interface {
	// Gerenciamento de índices
	CreateIndex(ctx context.Context, name string, mapping, settings map[string]interface{}) error
	DeleteIndex(ctx context.Context, name string) error
	IndexExists(ctx context.Context, name string) (bool, error)
	GetIndexInfo(ctx context.Context, name string) (*IndexInfo, error)
	ListIndices(ctx context.Context) ([]*IndexInfo, error)
	
	// Operações de documentos
	IndexDocument(ctx context.Context, index, id string, document map[string]interface{}) error
	UpdateDocument(ctx context.Context, index, id string, document map[string]interface{}) error
	DeleteDocument(ctx context.Context, index, id string) error
	GetDocument(ctx context.Context, index, id string) (map[string]interface{}, error)
	
	// Operações em lote
	BulkIndex(ctx context.Context, operations []BulkOperation) (*BulkResponse, error)
	
	// Busca
	Search(ctx context.Context, request *SearchRequest) (*SearchResponse, error)
	MultiSearch(ctx context.Context, requests []*SearchRequest) ([]*SearchResponse, error)
	Suggest(ctx context.Context, index, field, text string, size int) (*SuggestResponse, error)
	
	// Agregações
	Aggregate(ctx context.Context, index string, aggregations map[string]interface{}) (map[string]interface{}, error)
	
	// Monitoramento
	Health(ctx context.Context) (*HealthResponse, error)
	Stats(ctx context.Context, indices []string) (*StatsResponse, error)
}

// SearchStatisticsAggregate representa estatísticas agregadas
type SearchStatisticsAggregate struct {
	TotalSearches      int64   `json:"total_searches"`
	SuccessfulSearches int64   `json:"successful_searches"`
	FailedSearches     int64   `json:"failed_searches"`
	SuccessRate        float64 `json:"success_rate"`
	AvgResponseTime    float64 `json:"avg_response_time"`
	BasicSearches      int64   `json:"basic_searches"`
	AdvancedSearches   int64   `json:"advanced_searches"`
	SuggestionsRequests int64  `json:"suggestions_requests"`
}

// IndexInfo representa informações de um índice
type IndexInfo struct {
	Name          string                 `json:"name"`
	DocumentCount int64                  `json:"document_count"`
	SizeInBytes   int64                  `json:"size_in_bytes"`
	Mapping       map[string]interface{} `json:"mapping"`
	Settings      map[string]interface{} `json:"settings"`
	Aliases       []string               `json:"aliases"`
	Status        string                 `json:"status"`
}

// BulkOperation representa uma operação em lote
type BulkOperation struct {
	Operation    string                 `json:"operation"` // index, update, delete
	Index        string                 `json:"index"`
	DocumentID   string                 `json:"document_id"`
	Document     map[string]interface{} `json:"document,omitempty"`
}

// BulkResponse representa a resposta de uma operação em lote
type BulkResponse struct {
	Took        int                    `json:"took"`
	Errors      bool                   `json:"errors"`
	Items       []BulkResponseItem     `json:"items"`
	ProcessedCount int                 `json:"processed_count"`
	SuccessCount   int                 `json:"success_count"`
	ErrorCount     int                 `json:"error_count"`
}

// BulkResponseItem representa um item na resposta de operação em lote
type BulkResponseItem struct {
	Operation string                 `json:"operation"`
	Index     string                 `json:"index"`
	ID        string                 `json:"id"`
	Status    int                    `json:"status"`
	Error     map[string]interface{} `json:"error,omitempty"`
}

// SearchRequest representa uma requisição de busca
type SearchRequest struct {
	Index       string                 `json:"index"`
	Query       map[string]interface{} `json:"query"`
	Size        int                    `json:"size"`
	From        int                    `json:"from"`
	Sort        []map[string]interface{} `json:"sort,omitempty"`
	Highlight   map[string]interface{} `json:"highlight,omitempty"`
	Aggregations map[string]interface{} `json:"aggregations,omitempty"`
	Source      interface{}            `json:"_source,omitempty"`
}

// SearchResponse representa uma resposta de busca
type SearchResponse struct {
	Took        int                    `json:"took"`
	TimedOut    bool                   `json:"timed_out"`
	Hits        HitsResponse           `json:"hits"`
	Aggregations map[string]interface{} `json:"aggregations,omitempty"`
}

// HitsResponse representa os hits da busca
type HitsResponse struct {
	Total    TotalHits              `json:"total"`
	MaxScore float64                `json:"max_score"`
	Hits     []HitResponse          `json:"hits"`
}

// TotalHits representa o total de hits
type TotalHits struct {
	Value    int64  `json:"value"`
	Relation string `json:"relation"`
}

// HitResponse representa um hit individual
type HitResponse struct {
	Index     string                 `json:"_index"`
	ID        string                 `json:"_id"`
	Score     float64                `json:"_score"`
	Source    map[string]interface{} `json:"_source"`
	Highlight map[string][]string    `json:"highlight,omitempty"`
}

// SuggestResponse representa uma resposta de sugestão
type SuggestResponse struct {
	Suggestions []Suggestion `json:"suggestions"`
}

// Suggestion representa uma sugestão
type Suggestion struct {
	Text    string             `json:"text"`
	Offset  int                `json:"offset"`
	Length  int                `json:"length"`
	Options []SuggestionOption `json:"options"`
}

// SuggestionOption representa uma opção de sugestão
type SuggestionOption struct {
	Text  string  `json:"text"`
	Score float64 `json:"score"`
}

// HealthResponse representa uma resposta de health check
type HealthResponse struct {
	ClusterName   string `json:"cluster_name"`
	Status        string `json:"status"`
	TimedOut      bool   `json:"timed_out"`
	NumberOfNodes int    `json:"number_of_nodes"`
}

// StatsResponse representa uma resposta de estatísticas
type StatsResponse struct {
	Indices map[string]IndexStats `json:"indices"`
	Total   IndexStats            `json:"_all"`
}

// IndexStats representa estatísticas de um índice
type IndexStats struct {
	Total struct {
		Docs struct {
			Count   int64 `json:"count"`
			Deleted int64 `json:"deleted"`
		} `json:"docs"`
		Store struct {
			SizeInBytes int64 `json:"size_in_bytes"`
		} `json:"store"`
		Search struct {
			QueryTotal int64 `json:"query_total"`
			QueryTime  int64 `json:"query_time_in_millis"`
		} `json:"search"`
	} `json:"total"`
}