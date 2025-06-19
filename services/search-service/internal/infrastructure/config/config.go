package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config estrutura de configuração do serviço
type Config struct {
	// Aplicação
	ServiceName string `envconfig:"SERVICE_NAME" default:"search-service"`
	Version     string `envconfig:"VERSION" default:"1.0.0"`
	Port        int    `envconfig:"PORT" default:"8086"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`

	// Logging
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	// Database
	Database DatabaseConfig

	// Elasticsearch
	Elasticsearch ElasticsearchConfig

	// RabbitMQ  
	RabbitMQ RabbitMQConfig

	// Redis
	Redis RedisConfig

	// Metrics
	Metrics MetricsConfig

	// JWT
	JWTSecret string `envconfig:"JWT_SECRET" default:"your-secret-key-here"`

	// Search
	Search SearchConfig
}

// DatabaseConfig configurações do PostgreSQL
type DatabaseConfig struct {
	Host            string `envconfig:"DB_HOST" default:"localhost"`
	Port            int    `envconfig:"DB_PORT" default:"5432"`
	User            string `envconfig:"DB_USER" default:"postgres"`
	Password        string `envconfig:"DB_PASSWORD" required:"true"`
	Name            string `envconfig:"DB_NAME" default:"direito_lux_dev"`
	SSLMode         string `envconfig:"DB_SSL_MODE" default:"disable"`
	MaxOpenConns    int    `envconfig:"DB_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int    `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
	ConnMaxLifetime time.Duration `envconfig:"DB_CONN_MAX_LIFETIME" default:"5m"`
}

// RabbitMQConfig configurações do RabbitMQ
type RabbitMQConfig struct {
	URL      string `envconfig:"RABBITMQ_URL" required:"true"`
	Host     string `envconfig:"RABBITMQ_HOST" default:"localhost"`
	Port     int    `envconfig:"RABBITMQ_PORT" default:"5672"`
	User     string `envconfig:"RABBITMQ_USER" default:"guest"`
	Password string `envconfig:"RABBITMQ_PASSWORD" default:"guest"`
	VHost    string `envconfig:"RABBITMQ_VHOST" default:"/"`
}

// RedisConfig configurações do Redis
type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST" default:"localhost"`
	Port     int    `envconfig:"REDIS_PORT" default:"6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"REDIS_DB" default:"0"`
}

// ElasticsearchConfig configurações do Elasticsearch
type ElasticsearchConfig struct {
	URL      string `envconfig:"ELASTICSEARCH_URL" default:"http://localhost:9200"`
	Username string `envconfig:"ELASTICSEARCH_USERNAME" default:""`
	Password string `envconfig:"ELASTICSEARCH_PASSWORD" default:""`
	Index    string `envconfig:"ELASTICSEARCH_INDEX" default:"direito_lux"`
	Timeout  time.Duration `envconfig:"ELASTICSEARCH_TIMEOUT" default:"30s"`
}

// SearchConfig configurações específicas de busca
type SearchConfig struct {
	MaxResults      int           `envconfig:"MAX_SEARCH_RESULTS" default:"100"`
	DefaultPageSize int           `envconfig:"DEFAULT_PAGE_SIZE" default:"20"`
	CacheTTL        time.Duration `envconfig:"SEARCH_CACHE_TTL" default:"5m"`
	IndexBatchSize  int           `envconfig:"INDEX_BATCH_SIZE" default:"1000"`
	SearchTimeout   time.Duration `envconfig:"SEARCH_TIMEOUT" default:"30s"`
}

// MetricsConfig configurações de métricas
type MetricsConfig struct {
	Enabled bool `envconfig:"METRICS_ENABLED" default:"true"`
	Port    int  `envconfig:"METRICS_PORT" default:"9090"`
}

// Load carrega configuração a partir de variáveis de ambiente
func Load() (*Config, error) {
	var cfg Config
	
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	
	return &cfg, nil
}

// Validate valida a configuração
func (c *Config) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("porta inválida: %d", c.Port)
	}
	
	return nil
}

// IsDevelopment verifica se está em ambiente de desenvolvimento
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction verifica se está em ambiente de produção  
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
