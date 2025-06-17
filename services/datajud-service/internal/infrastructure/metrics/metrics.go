package metrics

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/direito-lux/datajud-service/internal/infrastructure/config"
)

// Metrics estrutura que contém todas as métricas
type Metrics struct {
	config *config.Config
	logger *zap.Logger
	server *http.Server

	// HTTP Metrics
	HTTPRequestsTotal     *prometheus.CounterVec
	HTTPRequestDuration   *prometheus.HistogramVec
	HTTPRequestsInFlight  prometheus.Gauge
	HTTPResponseSize      *prometheus.HistogramVec

	// Database Metrics
	DatabaseConnections   *prometheus.GaugeVec
	DatabaseQueries       *prometheus.CounterVec
	DatabaseQueryDuration *prometheus.HistogramVec

	// Cache Metrics
	CacheOperations       *prometheus.CounterVec
	CacheHitRate          *prometheus.GaugeVec
	CacheDuration         *prometheus.HistogramVec

	// Message Queue Metrics
	MessagesSent          *prometheus.CounterVec
	MessagesReceived      *prometheus.CounterVec
	MessageProcessingTime *prometheus.HistogramVec

	// Business Metrics
	TenantOperations      *prometheus.CounterVec
	UserOperations        *prometheus.CounterVec
	ProcessOperations     *prometheus.CounterVec

	// System Metrics
	GoroutinesCount       prometheus.Gauge
	MemoryUsage          prometheus.Gauge
	CPUUsage             prometheus.Gauge
}

// NewMetrics cria uma nova instância de métricas
func NewMetrics(cfg *config.Config, logger *zap.Logger) (*Metrics, error) {
	if !cfg.Metrics.Enabled {
		logger.Info("Métricas desabilitadas")
		return &Metrics{
			config: cfg,
			logger: logger,
		}, nil
	}

	namespace := cfg.Metrics.Namespace

	m := &Metrics{
		config: cfg,
		logger: logger,

		// HTTP Metrics
		HTTPRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"method", "path", "status", "tenant_id"},
		),

		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request duration in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "path", "status", "tenant_id"},
		),

		HTTPRequestsInFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "http_requests_in_flight",
				Help:      "Number of HTTP requests currently being processed",
			},
		),

		HTTPResponseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_response_size_bytes",
				Help:      "HTTP response size in bytes",
				Buckets:   prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "path", "status"},
		),

		// Database Metrics
		DatabaseConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "database_connections",
				Help:      "Number of database connections",
			},
			[]string{"state"}, // open, idle, in_use
		),

		DatabaseQueries: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "database_queries_total",
				Help:      "Total number of database queries",
			},
			[]string{"operation", "table", "tenant_id", "status"},
		),

		DatabaseQueryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "database_query_duration_seconds",
				Help:      "Database query duration in seconds",
				Buckets:   []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5},
			},
			[]string{"operation", "table", "tenant_id"},
		),

		// Cache Metrics
		CacheOperations: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "cache_operations_total",
				Help:      "Total number of cache operations",
			},
			[]string{"operation", "status", "tenant_id"}, // get, set, delete | hit, miss, error
		),

		CacheHitRate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cache_hit_rate",
				Help:      "Cache hit rate percentage",
			},
			[]string{"tenant_id"},
		),

		CacheDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "cache_operation_duration_seconds",
				Help:      "Cache operation duration in seconds",
				Buckets:   []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1},
			},
			[]string{"operation", "tenant_id"},
		),

		// Message Queue Metrics
		MessagesSent: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "messages_sent_total",
				Help:      "Total number of messages sent",
			},
			[]string{"exchange", "routing_key", "tenant_id", "status"},
		),

		MessagesReceived: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "messages_received_total",
				Help:      "Total number of messages received",
			},
			[]string{"queue", "tenant_id", "status"},
		),

		MessageProcessingTime: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "message_processing_duration_seconds",
				Help:      "Message processing duration in seconds",
				Buckets:   []float64{0.01, 0.05, 0.1, 0.5, 1, 5, 10, 30},
			},
			[]string{"queue", "tenant_id"},
		),

		// Business Metrics
		TenantOperations: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "tenant_operations_total",
				Help:      "Total number of tenant operations",
			},
			[]string{"operation", "tenant_id", "status"},
		),

		UserOperations: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "user_operations_total",
				Help:      "Total number of user operations",
			},
			[]string{"operation", "tenant_id", "status"},
		),

		ProcessOperations: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "process_operations_total",
				Help:      "Total number of process operations",
			},
			[]string{"operation", "tenant_id", "status"},
		),

		// System Metrics
		GoroutinesCount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "goroutines_count",
				Help:      "Number of goroutines",
			},
		),

		MemoryUsage: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "memory_usage_bytes",
				Help:      "Memory usage in bytes",
			},
		),

		CPUUsage: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cpu_usage_percent",
				Help:      "CPU usage percentage",
			},
		),
	}

	// Registrar métricas
	if err := m.register(); err != nil {
		return nil, fmt.Errorf("erro ao registrar métricas: %w", err)
	}

	// Iniciar servidor de métricas
	if err := m.startServer(); err != nil {
		return nil, fmt.Errorf("erro ao iniciar servidor de métricas: %w", err)
	}

	// Iniciar coleta de métricas do sistema
	go m.collectSystemMetrics()

	logger.Info("Métricas configuradas",
		zap.String("namespace", namespace),
		zap.Int("port", cfg.Metrics.Port),
		zap.String("path", cfg.Metrics.Path),
	)

	return m, nil
}

// register registra todas as métricas no Prometheus
func (m *Metrics) register() error {
	collectors := []prometheus.Collector{
		m.HTTPRequestsTotal,
		m.HTTPRequestDuration,
		m.HTTPRequestsInFlight,
		m.HTTPResponseSize,
		m.DatabaseConnections,
		m.DatabaseQueries,
		m.DatabaseQueryDuration,
		m.CacheOperations,
		m.CacheHitRate,
		m.CacheDuration,
		m.MessagesSent,
		m.MessagesReceived,
		m.MessageProcessingTime,
		m.TenantOperations,
		m.UserOperations,
		m.ProcessOperations,
		m.GoroutinesCount,
		m.MemoryUsage,
		m.CPUUsage,
	}

	for _, collector := range collectors {
		if err := prometheus.Register(collector); err != nil {
			return err
		}
	}

	return nil
}

// startServer inicia o servidor HTTP para métricas
func (m *Metrics) startServer() error {
	if !m.config.Metrics.Enabled {
		return nil
	}

	mux := http.NewServeMux()
	mux.Handle(m.config.Metrics.Path, promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	m.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", m.config.Metrics.Port),
		Handler: mux,
	}

	go func() {
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			m.logger.Error("Erro no servidor de métricas", zap.Error(err))
		}
	}()

	return nil
}

// Close fecha o servidor de métricas
func (m *Metrics) Close() error {
	if m.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return m.server.Shutdown(ctx)
	}
	return nil
}

// collectSystemMetrics coleta métricas do sistema periodicamente
func (m *Metrics) collectSystemMetrics() {
	if !m.config.Metrics.Enabled {
		return
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Coletar métricas do runtime do Go
		m.GoroutinesCount.Set(float64(runtime.NumGoroutine()))

		// Coletar métricas de memória
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		m.MemoryUsage.Set(float64(memStats.Alloc))
	}
}

// HTTPMiddleware middleware do Gin para coletar métricas HTTP
func (m *Metrics) HTTPMiddleware() gin.HandlerFunc {
	if !m.config.Metrics.Enabled {
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Incrementar requests em andamento
		m.HTTPRequestsInFlight.Inc()
		defer m.HTTPRequestsInFlight.Dec()

		// Processar request
		c.Next()

		// Coletar métricas
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		tenantID := c.GetString("tenant_id")

		// Incrementar contador total
		m.HTTPRequestsTotal.WithLabelValues(method, path, status, tenantID).Inc()

		// Registrar duração
		m.HTTPRequestDuration.WithLabelValues(method, path, status, tenantID).Observe(duration)

		// Registrar tamanho da resposta
		responseSize := float64(c.Writer.Size())
		if responseSize > 0 {
			m.HTTPResponseSize.WithLabelValues(method, path, status).Observe(responseSize)
		}
	}
}

// RecordDatabaseQuery registra métricas de query do banco
func (m *Metrics) RecordDatabaseQuery(operation, table, tenantID string, duration time.Duration, err error) {
	if !m.config.Metrics.Enabled {
		return
	}

	status := "success"
	if err != nil {
		status = "error"
	}

	m.DatabaseQueries.WithLabelValues(operation, table, tenantID, status).Inc()
	m.DatabaseQueryDuration.WithLabelValues(operation, table, tenantID).Observe(duration.Seconds())
}

// RecordCacheOperation registra métricas de operação de cache
func (m *Metrics) RecordCacheOperation(operation, tenantID string, hit bool, duration time.Duration) {
	if !m.config.Metrics.Enabled {
		return
	}

	status := "miss"
	if hit {
		status = "hit"
	}

	m.CacheOperations.WithLabelValues(operation, status, tenantID).Inc()
	m.CacheDuration.WithLabelValues(operation, tenantID).Observe(duration.Seconds())
}

// RecordMessageSent registra métrica de mensagem enviada
func (m *Metrics) RecordMessageSent(exchange, routingKey, tenantID string, success bool) {
	if !m.config.Metrics.Enabled {
		return
	}

	status := "success"
	if !success {
		status = "error"
	}

	m.MessagesSent.WithLabelValues(exchange, routingKey, tenantID, status).Inc()
}

// RecordMessageReceived registra métrica de mensagem recebida
func (m *Metrics) RecordMessageReceived(queue, tenantID string, success bool, processingTime time.Duration) {
	if !m.config.Metrics.Enabled {
		return
	}

	status := "success"
	if !success {
		status = "error"
	}

	m.MessagesReceived.WithLabelValues(queue, tenantID, status).Inc()
	m.MessageProcessingTime.WithLabelValues(queue, tenantID).Observe(processingTime.Seconds())
}

// RecordTenantOperation registra métrica de operação de tenant
func (m *Metrics) RecordTenantOperation(operation, tenantID string, success bool) {
	if !m.config.Metrics.Enabled {
		return
	}

	status := "success"
	if !success {
		status = "error"
	}

	m.TenantOperations.WithLabelValues(operation, tenantID, status).Inc()
}

// RecordUserOperation registra métrica de operação de usuário
func (m *Metrics) RecordUserOperation(operation, tenantID string, success bool) {
	if !m.config.Metrics.Enabled {
		return
	}

	status := "success"
	if !success {
		status = "error"
	}

	m.UserOperations.WithLabelValues(operation, tenantID, status).Inc()
}

// RecordProcessOperation registra métrica de operação de processo
func (m *Metrics) RecordProcessOperation(operation, tenantID string, success bool) {
	if !m.config.Metrics.Enabled {
		return
	}

	status := "success"
	if !success {
		status = "error"
	}

	m.ProcessOperations.WithLabelValues(operation, tenantID, status).Inc()
}