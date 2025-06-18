package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config configuração principal do serviço
type Config struct {
	ServiceName string `json:"service_name"`
	Version     string `json:"version"`
	Server ServerConfig `json:"server"`
	Database DatabaseConfig `json:"database"`
	RabbitMQ RabbitMQConfig `json:"rabbitmq"`
	Events EventsConfig `json:"events"`
	Metrics MetricsConfig `json:"metrics"`
	Jaeger JaegerConfig `json:"jaeger"`
	Environment string `json:"environment"`
}

// ServerConfig configurações do servidor HTTP
type ServerConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

// DatabaseConfig configurações do PostgreSQL
type DatabaseConfig struct {
	Host            string `json:"host"`
	Port            int    `json:"port"`
	User            string `json:"user"`
	Password        string `json:"password"`
	Name            string `json:"name"`
	SSLMode         string `json:"ssl_mode"`
	MaxOpenConns    int    `json:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime"`
	MigrationsPath  string `json:"migrations_path"`
}

// RabbitMQConfig configurações do RabbitMQ
type RabbitMQConfig struct {
	Host           string        `json:"host"`
	Port           int           `json:"port"`
	User           string        `json:"user"`
	Password       string        `json:"password"`
	VHost          string        `json:"vhost"`
	Exchange       string        `json:"exchange"`
	URL            string        `json:"url"`
	Queue          string        `json:"queue"`
	RoutingKey     string        `json:"routing_key"`
	PrefetchCount  int           `json:"prefetch_count"`
	Durable        bool          `json:"durable"`
	AutoDelete     bool          `json:"auto_delete"`
	Exclusive      bool          `json:"exclusive"`
	NoWait         bool          `json:"no_wait"`
	ReconnectDelay time.Duration `json:"reconnect_delay"`
}

// EventsConfig configurações de eventos
type EventsConfig struct {
	AsyncEnabled  bool `json:"async_enabled"`
	BatchSize     int  `json:"batch_size"`
	FlushInterval int  `json:"flush_interval"`
}

// MetricsConfig configurações de métricas
type MetricsConfig struct {
	Enabled   bool   `json:"enabled"`
	Port      int    `json:"port"`
	Path      string `json:"path"`
	Namespace string `json:"namespace"`
}

// JaegerConfig configurações do Jaeger Tracing
type JaegerConfig struct {
	Enabled      bool   `json:"enabled"`
	ServiceName  string `json:"service_name"`
	AgentHost    string `json:"agent_host"`
	AgentPort    string `json:"agent_port"`
	Endpoint     string `json:"endpoint"`
	SamplerType  string `json:"sampler_type"`
	SamplerParam float64 `json:"sampler_param"`
}

// LoadConfig carrega configuração a partir de variáveis de ambiente
func LoadConfig() (*Config, error) {
	cfg := &Config{
		ServiceName: getEnv("SERVICE_NAME", "process-service"),
		Version:     getEnv("VERSION", "1.0.0"),
		Environment: getEnv("ENVIRONMENT", "development"),
		
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  time.Duration(getEnvInt("SERVER_READ_TIMEOUT", 30)) * time.Second,
			WriteTimeout: time.Duration(getEnvInt("SERVER_WRITE_TIMEOUT", 30)) * time.Second,
		},
		
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			Name:            getEnv("DB_NAME", "process_service"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvInt("DB_CONN_MAX_LIFETIME", 5),
			MigrationsPath:  getEnv("DB_MIGRATIONS_PATH", "./migrations"),
		},
		
		RabbitMQ: RabbitMQConfig{
			Host:           getEnv("RABBITMQ_HOST", "localhost"),
			Port:           getEnvInt("RABBITMQ_PORT", 5672),
			User:           getEnv("RABBITMQ_USER", "guest"),
			Password:       getEnv("RABBITMQ_PASSWORD", "guest"),
			VHost:          getEnv("RABBITMQ_VHOST", "/"),
			Exchange:       getEnv("RABBITMQ_EXCHANGE", "direito-lux.events"),
			URL:            getEnv("RABBITMQ_URL", ""),
			Queue:          getEnv("RABBITMQ_QUEUE", "process-service"),
			RoutingKey:     getEnv("RABBITMQ_ROUTING_KEY", "process"),
			PrefetchCount:  getEnvInt("RABBITMQ_PREFETCH_COUNT", 10),
			Durable:        getEnvBool("RABBITMQ_DURABLE", true),
			AutoDelete:     getEnvBool("RABBITMQ_AUTO_DELETE", false),
			Exclusive:      getEnvBool("RABBITMQ_EXCLUSIVE", false),
			NoWait:         getEnvBool("RABBITMQ_NO_WAIT", false),
			ReconnectDelay: time.Duration(getEnvInt("RABBITMQ_RECONNECT_DELAY", 5)) * time.Second,
		},
		
		Events: EventsConfig{
			AsyncEnabled:  getEnvBool("EVENTS_ASYNC_ENABLED", true),
			BatchSize:     getEnvInt("EVENTS_BATCH_SIZE", 10),
			FlushInterval: getEnvInt("EVENTS_FLUSH_INTERVAL", 5000),
		},

		Metrics: MetricsConfig{
			Enabled:   getEnvBool("METRICS_ENABLED", true),
			Port:      getEnvInt("METRICS_PORT", 9090),
			Path:      getEnv("METRICS_PATH", "/metrics"),
			Namespace: getEnv("METRICS_NAMESPACE", "process_service"),
		},

		Jaeger: JaegerConfig{
			Enabled:      getEnvBool("JAEGER_ENABLED", false),
			ServiceName:  getEnv("JAEGER_SERVICE_NAME", "process-service"),
			AgentHost:    getEnv("JAEGER_AGENT_HOST", "localhost"),
			AgentPort:    getEnv("JAEGER_AGENT_PORT", "6831"),
			Endpoint:     getEnv("JAEGER_ENDPOINT", ""),
			SamplerType:  getEnv("JAEGER_SAMPLER_TYPE", "const"),
			SamplerParam: getEnvFloat("JAEGER_SAMPLER_PARAM", 1.0),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

// IsDevelopment verifica se está em ambiente de desenvolvimento
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction verifica se está em ambiente de produção  
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// GetDatabaseDSN retorna a string de conexão com o banco
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}
