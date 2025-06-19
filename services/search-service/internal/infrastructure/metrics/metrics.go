package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all prometheus metrics
type Metrics struct {
	// HTTP metrics
	HTTPRequests *prometheus.CounterVec
	HTTPDuration *prometheus.HistogramVec
	
	// Search metrics
	SearchRequests *prometheus.CounterVec
	SearchDuration *prometheus.HistogramVec
	IndexedDocuments prometheus.Counter
	
	// Elasticsearch metrics
	ElasticsearchRequests *prometheus.CounterVec
	ElasticsearchDuration *prometheus.HistogramVec
}

// NewMetrics creates new metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		HTTPRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "search_http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		HTTPDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "search_http_request_duration_seconds",
				Help: "HTTP request duration in seconds",
			},
			[]string{"method", "endpoint"},
		),
		SearchRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "search_requests_total",
				Help: "Total number of search requests",
			},
			[]string{"type", "status"},
		),
		SearchDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "search_duration_seconds",
				Help: "Search request duration in seconds",
			},
			[]string{"type"},
		),
		IndexedDocuments: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "search_indexed_documents_total",
				Help: "Total number of indexed documents",
			},
		),
		ElasticsearchRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "search_elasticsearch_requests_total",
				Help: "Total number of Elasticsearch requests",
			},
			[]string{"operation", "status"},
		),
		ElasticsearchDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "search_elasticsearch_request_duration_seconds",
				Help: "Elasticsearch request duration in seconds",
			},
			[]string{"operation"},
		),
	}
}