package config

import (
	"os"
	"strconv"
	"time"
)

// Config configuração principal do serviço
type Config struct {
	Server ServerConfig `json:"server"`
	Database DatabaseConfig `json:"database"`
	RabbitMQ RabbitMQConfig `json:"rabbitmq"`
	Events EventsConfig `json:"events"`
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
}

// RabbitMQConfig configurações do RabbitMQ
type RabbitMQConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	VHost    string `json:"vhost"`
	Exchange string `json:"exchange"`
}

// EventsConfig configurações de eventos
type EventsConfig struct {
	AsyncEnabled  bool `json:"async_enabled"`
	BatchSize     int  `json:"batch_size"`
	FlushInterval int  `json:"flush_interval"`
}

// LoadConfig carrega configuração a partir de variáveis de ambiente
func LoadConfig() (*Config, error) {
	cfg := &Config{
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
		},
		
		RabbitMQ: RabbitMQConfig{
			Host:     getEnv("RABBITMQ_HOST", "localhost"),
			Port:     getEnvInt("RABBITMQ_PORT", 5672),
			User:     getEnv("RABBITMQ_USER", "guest"),
			Password: getEnv("RABBITMQ_PASSWORD", "guest"),
			VHost:    getEnv("RABBITMQ_VHOST", "/"),
			Exchange: getEnv("RABBITMQ_EXCHANGE", "direito-lux.events"),
		},
		
		Events: EventsConfig{
			AsyncEnabled:  getEnvBool("EVENTS_ASYNC_ENABLED", true),
			BatchSize:     getEnvInt("EVENTS_BATCH_SIZE", 10),
			FlushInterval: getEnvInt("EVENTS_FLUSH_INTERVAL", 5000),
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
