package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config estrutura de configuração do serviço
type Config struct {
	// Aplicação
	Version     string `envconfig:"VERSION" default:"1.0.0"`
	Port        int    `envconfig:"PORT" default:"8080"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`

	// Logging
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	// Database
	Database DatabaseConfig

	// RabbitMQ  
	RabbitMQ RabbitMQConfig

	// Redis
	Redis RedisConfig

	// Metrics
	Metrics MetricsConfig
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
