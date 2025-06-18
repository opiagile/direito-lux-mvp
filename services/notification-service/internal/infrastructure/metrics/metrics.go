package metrics

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/infrastructure/config"
)

// Metrics contém todas as métricas do serviço
type Metrics struct {
	// Contador de notificações enviadas
	NotificationsSent prometheus.CounterVec

	// Contador de notificações falhadas
	NotificationsFailed prometheus.CounterVec

	// Histograma de duração de processamento
	ProcessingDuration prometheus.HistogramVec

	// Gauge de conexões ativas
	ActiveConnections prometheus.Gauge

	// Gauge de memoria utilizada
	MemoryUsage prometheus.Gauge

	// Gauge de goroutines
	GoroutineCount prometheus.Gauge

	// Contador de requests HTTP
	HTTPRequests prometheus.CounterVec

	// Histograma de duração de requests HTTP
	HTTPDuration prometheus.HistogramVec

	// Gauge de health check
	HealthStatus prometheus.Gauge
}

// NewMetrics cria novas métricas
func NewMetrics() *Metrics {
	return &Metrics{
		NotificationsSent: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "notifications_sent_total",
				Help: "Total de notificações enviadas",
			},
			[]string{"type", "channel", "tenant_id"},
		),

		NotificationsFailed: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "notifications_failed_total",
				Help: "Total de notificações falhadas",
			},
			[]string{"type", "channel", "tenant_id", "error_type"},
		),

		ProcessingDuration: *promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "notification_processing_duration_seconds",
				Help:    "Duração de processamento de notificações",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"type", "channel"},
		),

		ActiveConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_connections",
				Help: "Número de conexões ativas",
			},
		),

		MemoryUsage: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "memory_usage_bytes",
				Help: "Uso de memória em bytes",
			},
		),

		GoroutineCount: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "goroutine_count",
				Help: "Número de goroutines",
			},
		),

		HTTPRequests: *promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total de requests HTTP",
			},
			[]string{"method", "path", "status"},
		),

		HTTPDuration: *promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duração de requests HTTP",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"method", "path"},
		),

		HealthStatus: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "health_status",
				Help: "Status de saúde do serviço (1=healthy, 0=unhealthy)",
			},
		),
	}
}

// StartMetricsServer inicia servidor de métricas
func StartMetricsServer(cfg *config.MetricsConfig, logger *zap.Logger) error {
	if !cfg.Enabled {
		logger.Info("Métricas Prometheus desabilitadas")
		return nil
	}

	http.Handle("/metrics", promhttp.Handler())
	
	address := fmt.Sprintf(":%d", cfg.Port)
	logger.Info("Iniciando servidor de métricas", zap.String("address", address))
	
	return http.ListenAndServe(address, nil)
}

// UpdateSystemMetrics atualiza métricas do sistema
func (m *Metrics) UpdateSystemMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.MemoryUsage.Set(float64(memStats.Alloc))
	m.GoroutineCount.Set(float64(runtime.NumGoroutine()))
}

// RecordNotificationSent registra notificação enviada
func (m *Metrics) RecordNotificationSent(notificationType, channel, tenantID string) {
	m.NotificationsSent.WithLabelValues(notificationType, channel, tenantID).Inc()
}

// RecordNotificationFailed registra notificação falhada
func (m *Metrics) RecordNotificationFailed(notificationType, channel, tenantID, errorType string) {
	m.NotificationsFailed.WithLabelValues(notificationType, channel, tenantID, errorType).Inc()
}

// RecordProcessingDuration registra duração de processamento
func (m *Metrics) RecordProcessingDuration(notificationType, channel string, duration time.Duration) {
	m.ProcessingDuration.WithLabelValues(notificationType, channel).Observe(duration.Seconds())
}

// RecordHTTPRequest registra request HTTP
func (m *Metrics) RecordHTTPRequest(method, path, status string, duration time.Duration) {
	m.HTTPRequests.WithLabelValues(method, path, status).Inc()
	m.HTTPDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

// SetHealthStatus define status de saúde
func (m *Metrics) SetHealthStatus(healthy bool) {
	if healthy {
		m.HealthStatus.Set(1)
	} else {
		m.HealthStatus.Set(0)
	}
}

// IncrementActiveConnections incrementa conexões ativas
func (m *Metrics) IncrementActiveConnections() {
	m.ActiveConnections.Inc()
}

// DecrementActiveConnections decrementa conexões ativas
func (m *Metrics) DecrementActiveConnections() {
	m.ActiveConnections.Dec()
}

// StartSystemMetricsCollector inicia coleta de métricas do sistema
func (m *Metrics) StartSystemMetricsCollector(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.UpdateSystemMetrics()
		}
	}
}