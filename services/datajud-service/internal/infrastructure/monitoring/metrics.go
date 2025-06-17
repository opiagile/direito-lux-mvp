package monitoring

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/google/uuid"
	"github.com/direito-lux/datajud-service/internal/domain"
)

// DataJudMetrics coleta métricas específicas do DataJud Service
type DataJudMetrics struct {
	// Métricas de requisições
	requestsTotal         *prometheus.CounterVec
	requestDuration       *prometheus.HistogramVec
	requestsInFlight      *prometheus.GaugeVec
	requestErrors         *prometheus.CounterVec
	
	// Métricas de CNPJ providers
	cnpjQuotaUsage        *prometheus.GaugeVec
	cnpjQuotaLimit        *prometheus.GaugeVec
	cnpjRequestsTotal     *prometheus.CounterVec
	cnpjErrors            *prometheus.CounterVec
	cnpjResponseTime      *prometheus.HistogramVec
	
	// Métricas de rate limiting
	rateLimitHits         *prometheus.CounterVec
	rateLimitCurrent      *prometheus.GaugeVec
	rateLimitLimit        *prometheus.GaugeVec
	
	// Métricas de circuit breaker
	circuitBreakerState   *prometheus.GaugeVec
	circuitBreakerTrips   *prometheus.CounterVec
	circuitBreakerRequests *prometheus.CounterVec
	
	// Métricas de cache
	cacheHits             *prometheus.CounterVec
	cacheMisses           *prometheus.CounterVec
	cacheSize             *prometheus.GaugeVec
	cacheEvictions        *prometheus.CounterVec
	
	// Métricas de queue
	queueSize             *prometheus.GaugeVec
	queueProcessingTime   *prometheus.HistogramVec
	queueWorkers          *prometheus.GaugeVec
	
	// Métricas específicas DataJud
	datajudAPIResponseTime *prometheus.HistogramVec
	datajudAPIErrors       *prometheus.CounterVec
	datajudAPIRateLimit    *prometheus.GaugeVec
	
	mu sync.RWMutex
}

// NewDataJudMetrics cria nova instância de métricas
func NewDataJudMetrics() *DataJudMetrics {
	return &DataJudMetrics{
		// Métricas de requisições
		requestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "datajud_requests_total",
				Help: "Total number of DataJud requests",
			},
			[]string{"tenant_id", "type", "status", "court_id"},
		),
		
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "datajud_request_duration_seconds",
				Help:    "Duration of DataJud requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"tenant_id", "type", "court_id"},
		),
		
		requestsInFlight: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "datajud_requests_in_flight",
				Help: "Number of DataJud requests currently being processed",
			},
			[]string{"type"},
		),
		
		requestErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "datajud_request_errors_total",
				Help: "Total number of DataJud request errors",
			},
			[]string{"tenant_id", "type", "error_code", "court_id"},
		),
		
		// Métricas de CNPJ providers
		cnpjQuotaUsage: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "datajud_cnpj_quota_usage",
				Help: "Current quota usage for CNPJ providers",
			},
			[]string{"cnpj", "tenant_id", "provider_name"},
		),
		
		cnpjQuotaLimit: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "datajud_cnpj_quota_limit",
				Help: "Daily quota limit for CNPJ providers",
			},
			[]string{"cnpj", "tenant_id", "provider_name"},
		),
		
		cnpjRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "datajud_cnpj_requests_total",
				Help: "Total requests made by each CNPJ provider",
			},
			[]string{"cnpj", "tenant_id", "provider_name", "status"},
		),
		
		cnpjErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "datajud_cnpj_errors_total",
				Help: "Total errors for each CNPJ provider",
			},
			[]string{"cnpj", "tenant_id", "provider_name", "error_type"},
		),
		
		cnpjResponseTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "datajud_cnpj_response_time_seconds",
				Help:    "Response time for each CNPJ provider",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"cnpj", "tenant_id", "provider_name"},
		),
		
		// Métricas de rate limiting
		rateLimitHits: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "datajud_rate_limit_hits_total",
				Help: "Total rate limit hits",
			},
			[]string{"limit_type", "key"},
		),
		
		rateLimitCurrent: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "datajud_rate_limit_current",
				Help: "Current rate limit usage",
			},
			[]string{"limit_type", "key"},
		),
		
		rateLimitLimit: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "datajud_rate_limit_limit",
				Help: "Rate limit threshold",
			},
			[]string{"limit_type", "key"},
		),
		
		// Métricas de circuit breaker
		circuitBreakerState: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "datajud_circuit_breaker_state",
				Help: "Circuit breaker state (0=closed, 1=open, 2=half-open)",
			},
			[]string{"breaker_name"},
		),
		
		circuitBreakerTrips: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "datajud_circuit_breaker_trips_total",
				Help: "Total circuit breaker trips",
			},
			[]string{"breaker_name", "from_state", "to_state"},
		),
		
		circuitBreakerRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "datajud_circuit_breaker_requests_total",
				Help: "Total circuit breaker requests",
			},
			[]string{"breaker_name", "result"},
		),
		
		// Métricas de cache
		cacheHits: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "datajud_cache_hits_total",
				Help: "Total cache hits",
			},
			[]string{"request_type", "tenant_id"},
		),
		
		cacheMisses: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "datajud_cache_misses_total",
				Help: "Total cache misses",
			},
			[]string{"request_type", "tenant_id"},
		),
		
		cacheSize: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "datajud_cache_size_bytes",
				Help: "Current cache size in bytes",
			},
			[]string{"request_type"},
		),
		
		cacheEvictions: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "datajud_cache_evictions_total",
				Help: "Total cache evictions",
			},
			[]string{"request_type", "reason"},
		),
		
		// Métricas de queue
		queueSize: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "datajud_queue_size",
				Help: "Current queue size",
			},
			[]string{"priority"},
		),
		
		queueProcessingTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "datajud_queue_processing_time_seconds",
				Help:    "Time spent processing queue items",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"priority"},
		),
		
		queueWorkers: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "datajud_queue_workers",
				Help: "Number of active queue workers",
			},
			[]string{"status"},
		),
		
		// Métricas específicas DataJud API
		datajudAPIResponseTime: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "datajud_api_response_time_seconds",
				Help:    "DataJud API response time",
				Buckets: []float64{0.1, 0.5, 1.0, 2.5, 5.0, 10.0},
			},
			[]string{"endpoint", "status_code", "court_id"},
		),
		
		datajudAPIErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "datajud_api_errors_total",
				Help: "Total DataJud API errors",
			},
			[]string{"endpoint", "error_type", "status_code", "court_id"},
		),
		
		datajudAPIRateLimit: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "datajud_api_rate_limit_remaining",
				Help: "Remaining DataJud API rate limit",
			},
			[]string{"cnpj", "endpoint"},
		),
	}
}

// RecordRequest registra uma requisição DataJud
func (m *DataJudMetrics) RecordRequest(tenantID uuid.UUID, requestType domain.RequestType, status string, courtID string) {
	m.requestsTotal.WithLabelValues(tenantID.String(), string(requestType), status, courtID).Inc()
}

// RecordRequestDuration registra duração de uma requisição
func (m *DataJudMetrics) RecordRequestDuration(tenantID uuid.UUID, requestType domain.RequestType, courtID string, duration time.Duration) {
	m.requestDuration.WithLabelValues(tenantID.String(), string(requestType), courtID).Observe(duration.Seconds())
}

// IncRequestsInFlight incrementa requisições em andamento
func (m *DataJudMetrics) IncRequestsInFlight(requestType domain.RequestType) {
	m.requestsInFlight.WithLabelValues(string(requestType)).Inc()
}

// DecRequestsInFlight decrementa requisições em andamento
func (m *DataJudMetrics) DecRequestsInFlight(requestType domain.RequestType) {
	m.requestsInFlight.WithLabelValues(string(requestType)).Dec()
}

// RecordRequestError registra erro de requisição
func (m *DataJudMetrics) RecordRequestError(tenantID uuid.UUID, requestType domain.RequestType, errorCode, courtID string) {
	m.requestErrors.WithLabelValues(tenantID.String(), string(requestType), errorCode, courtID).Inc()
}

// UpdateCNPJQuota atualiza métricas de quota de CNPJ
func (m *DataJudMetrics) UpdateCNPJQuota(provider *domain.CNPJProvider) {
	labels := []string{provider.CNPJ, provider.TenantID.String(), provider.Name}
	
	m.cnpjQuotaUsage.WithLabelValues(labels...).Set(float64(provider.DailyUsage))
	m.cnpjQuotaLimit.WithLabelValues(labels...).Set(float64(provider.DailyLimit))
}

// RecordCNPJRequest registra requisição por CNPJ
func (m *DataJudMetrics) RecordCNPJRequest(provider *domain.CNPJProvider, status string) {
	labels := []string{provider.CNPJ, provider.TenantID.String(), provider.Name, status}
	m.cnpjRequestsTotal.WithLabelValues(labels...).Inc()
}

// RecordCNPJError registra erro por CNPJ
func (m *DataJudMetrics) RecordCNPJError(provider *domain.CNPJProvider, errorType string) {
	labels := []string{provider.CNPJ, provider.TenantID.String(), provider.Name, errorType}
	m.cnpjErrors.WithLabelValues(labels...).Inc()
}

// RecordCNPJResponseTime registra tempo de resposta por CNPJ
func (m *DataJudMetrics) RecordCNPJResponseTime(provider *domain.CNPJProvider, duration time.Duration) {
	labels := []string{provider.CNPJ, provider.TenantID.String(), provider.Name}
	m.cnpjResponseTime.WithLabelValues(labels...).Observe(duration.Seconds())
}

// RecordRateLimitHit registra hit de rate limit
func (m *DataJudMetrics) RecordRateLimitHit(limitType domain.RateLimitType, key string) {
	m.rateLimitHits.WithLabelValues(string(limitType), key).Inc()
}

// UpdateRateLimit atualiza métricas de rate limit
func (m *DataJudMetrics) UpdateRateLimit(limitType domain.RateLimitType, key string, current, limit int) {
	labels := []string{string(limitType), key}
	m.rateLimitCurrent.WithLabelValues(labels...).Set(float64(current))
	m.rateLimitLimit.WithLabelValues(labels...).Set(float64(limit))
}

// UpdateCircuitBreakerState atualiza estado do circuit breaker
func (m *DataJudMetrics) UpdateCircuitBreakerState(name string, state domain.CircuitBreakerState) {
	var stateValue float64
	switch state {
	case domain.StateClosed:
		stateValue = 0
	case domain.StateOpen:
		stateValue = 1
	case domain.StateHalfOpen:
		stateValue = 2
	}
	
	m.circuitBreakerState.WithLabelValues(name).Set(stateValue)
}

// RecordCircuitBreakerTrip registra mudança de estado do circuit breaker
func (m *DataJudMetrics) RecordCircuitBreakerTrip(name string, fromState, toState domain.CircuitBreakerState) {
	m.circuitBreakerTrips.WithLabelValues(name, string(fromState), string(toState)).Inc()
}

// RecordCircuitBreakerRequest registra requisição do circuit breaker
func (m *DataJudMetrics) RecordCircuitBreakerRequest(name string, success bool) {
	result := "failure"
	if success {
		result = "success"
	}
	m.circuitBreakerRequests.WithLabelValues(name, result).Inc()
}

// RecordCacheHit registra hit no cache
func (m *DataJudMetrics) RecordCacheHit(requestType domain.RequestType, tenantID uuid.UUID) {
	m.cacheHits.WithLabelValues(string(requestType), tenantID.String()).Inc()
}

// RecordCacheMiss registra miss no cache
func (m *DataJudMetrics) RecordCacheMiss(requestType domain.RequestType, tenantID uuid.UUID) {
	m.cacheMisses.WithLabelValues(string(requestType), tenantID.String()).Inc()
}

// UpdateCacheSize atualiza tamanho do cache
func (m *DataJudMetrics) UpdateCacheSize(requestType domain.RequestType, sizeBytes int64) {
	m.cacheSize.WithLabelValues(string(requestType)).Set(float64(sizeBytes))
}

// RecordCacheEviction registra eviction do cache
func (m *DataJudMetrics) RecordCacheEviction(requestType domain.RequestType, reason string) {
	m.cacheEvictions.WithLabelValues(string(requestType), reason).Inc()
}

// UpdateQueueSize atualiza tamanho da fila
func (m *DataJudMetrics) UpdateQueueSize(priority domain.RequestPriority, size int) {
	m.queueSize.WithLabelValues(priorityToString(priority)).Set(float64(size))
}

// RecordQueueProcessingTime registra tempo de processamento da fila
func (m *DataJudMetrics) RecordQueueProcessingTime(priority domain.RequestPriority, duration time.Duration) {
	m.queueProcessingTime.WithLabelValues(priorityToString(priority)).Observe(duration.Seconds())
}

// UpdateQueueWorkers atualiza número de workers
func (m *DataJudMetrics) UpdateQueueWorkers(active, idle int) {
	m.queueWorkers.WithLabelValues("active").Set(float64(active))
	m.queueWorkers.WithLabelValues("idle").Set(float64(idle))
}

// RecordDataJudAPIResponse registra resposta da API DataJud
func (m *DataJudMetrics) RecordDataJudAPIResponse(endpoint, statusCode, courtID string, duration time.Duration) {
	labels := []string{endpoint, statusCode, courtID}
	m.datajudAPIResponseTime.WithLabelValues(labels...).Observe(duration.Seconds())
}

// RecordDataJudAPIError registra erro da API DataJud
func (m *DataJudMetrics) RecordDataJudAPIError(endpoint, errorType, statusCode, courtID string) {
	labels := []string{endpoint, errorType, statusCode, courtID}
	m.datajudAPIErrors.WithLabelValues(labels...).Inc()
}

// UpdateDataJudAPIRateLimit atualiza rate limit da API DataJud
func (m *DataJudMetrics) UpdateDataJudAPIRateLimit(cnpj, endpoint string, remaining int) {
	m.datajudAPIRateLimit.WithLabelValues(cnpj, endpoint).Set(float64(remaining))
}

// priorityToString converte prioridade para string
func priorityToString(priority domain.RequestPriority) string {
	switch priority {
	case domain.PriorityUrgent:
		return "urgent"
	case domain.PriorityHigh:
		return "high"
	case domain.PriorityNormal:
		return "normal"
	case domain.PriorityLow:
		return "low"
	default:
		return "unknown"
	}
}

// MetricsCollector coleta métricas periodicamente
type MetricsCollector struct {
	metrics       *DataJudMetrics
	repos         *domain.Repositories
	interval      time.Duration
	stopChan      chan bool
	isRunning     bool
	mu            sync.RWMutex
}

// NewMetricsCollector cria novo coletor de métricas
func NewMetricsCollector(metrics *DataJudMetrics, repos *domain.Repositories, interval time.Duration) *MetricsCollector {
	return &MetricsCollector{
		metrics:  metrics,
		repos:    repos,
		interval: interval,
		stopChan: make(chan bool),
	}
}

// Start inicia coleta periódica de métricas
func (mc *MetricsCollector) Start(ctx context.Context) {
	mc.mu.Lock()
	if mc.isRunning {
		mc.mu.Unlock()
		return
	}
	mc.isRunning = true
	mc.mu.Unlock()

	ticker := time.NewTicker(mc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-mc.stopChan:
			return
		case <-ticker.C:
			mc.collectMetrics(ctx)
		}
	}
}

// Stop para a coleta de métricas
func (mc *MetricsCollector) Stop() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	if mc.isRunning {
		close(mc.stopChan)
		mc.isRunning = false
	}
}

// collectMetrics coleta métricas dos repositórios
func (mc *MetricsCollector) collectMetrics(ctx context.Context) {
	// Coletar métricas de CNPJ providers
	mc.collectCNPJMetrics(ctx)
	
	// Coletar métricas de rate limiting
	mc.collectRateLimitMetrics(ctx)
	
	// Coletar métricas de circuit breakers
	mc.collectCircuitBreakerMetrics(ctx)
	
	// Coletar métricas de cache
	mc.collectCacheMetrics(ctx)
}

// collectCNPJMetrics coleta métricas dos CNPJ providers
func (mc *MetricsCollector) collectCNPJMetrics(ctx context.Context) {
	providers, err := mc.repos.CNPJProvider.FindActiveCNPJs()
	if err != nil {
		return
	}

	for _, provider := range providers {
		mc.metrics.UpdateCNPJQuota(provider)
	}
}

// collectRateLimitMetrics coleta métricas de rate limiting
func (mc *MetricsCollector) collectRateLimitMetrics(ctx context.Context) {
	// Implementação específica dependeria do repositório de rate limiter
	// Por simplicidade, não implementado aqui
}

// collectCircuitBreakerMetrics coleta métricas de circuit breakers
func (mc *MetricsCollector) collectCircuitBreakerMetrics(ctx context.Context) {
	breakers, err := mc.repos.CircuitBreaker.FindAll()
	if err != nil {
		return
	}

	for _, breaker := range breakers {
		mc.metrics.UpdateCircuitBreakerState(breaker.Name, breaker.GetState())
	}
}

// collectCacheMetrics coleta métricas de cache
func (mc *MetricsCollector) collectCacheMetrics(ctx context.Context) {
	// Implementação específica dependeria do repositório de cache
	// Por simplicidade, não implementado aqui
}