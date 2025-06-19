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

	// Server
	Server ServerConfig

	// Database
	Database DatabaseConfig

	// RabbitMQ  
	RabbitMQ RabbitMQConfig

	// Redis
	Redis RedisConfig

	// Metrics
	Metrics MetricsConfig

	// SMTP
	SMTP SMTPConfig

	// WhatsApp
	WhatsApp WhatsAppConfig

	// Telegram
	Telegram TelegramConfig
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

// ServerConfig configurações do servidor
type ServerConfig struct {
	Port         int `envconfig:"SERVER_PORT" default:"8080"`
	ReadTimeout  int `envconfig:"SERVER_READ_TIMEOUT" default:"30"`
	WriteTimeout int `envconfig:"SERVER_WRITE_TIMEOUT" default:"30"`
	IdleTimeout  int `envconfig:"SERVER_IDLE_TIMEOUT" default:"60"`
}

// MetricsConfig configurações de métricas
type MetricsConfig struct {
	Enabled bool `envconfig:"METRICS_ENABLED" default:"true"`
	Port    int  `envconfig:"METRICS_PORT" default:"9090"`
}

// SMTPConfig configurações do SMTP
type SMTPConfig struct {
	Host        string `envconfig:"SMTP_HOST" default:"localhost"`
	Port        int    `envconfig:"SMTP_PORT" default:"587"`
	Username    string `envconfig:"SMTP_USERNAME" required:"true"`
	Password    string `envconfig:"SMTP_PASSWORD" required:"true"`
	FromEmail   string `envconfig:"SMTP_FROM_EMAIL" required:"true"`
	FromName    string `envconfig:"SMTP_FROM_NAME" default:"Direito Lux"`
	UseTLS      bool   `envconfig:"SMTP_USE_TLS" default:"false"`
	UseStartTLS bool   `envconfig:"SMTP_USE_STARTTLS" default:"true"`
}

// WhatsAppConfig configurações do WhatsApp Business API
type WhatsAppConfig struct {
	AccessToken   string `envconfig:"WHATSAPP_ACCESS_TOKEN" required:"true"`
	PhoneNumberID string `envconfig:"WHATSAPP_PHONE_NUMBER_ID" required:"true"`
	WebhookURL    string `envconfig:"WHATSAPP_WEBHOOK_URL"`
	VerifyToken   string `envconfig:"WHATSAPP_VERIFY_TOKEN" required:"true"`
}

// TelegramConfig configurações do Telegram Bot API
type TelegramConfig struct {
	BotToken      string `envconfig:"TELEGRAM_BOT_TOKEN" required:"true"`
	WebhookURL    string `envconfig:"TELEGRAM_WEBHOOK_URL"`
	WebhookSecret string `envconfig:"TELEGRAM_WEBHOOK_SECRET"`
}

// Load carrega configuração a partir de variáveis de ambiente
func Load() (*Config, error) {
	var cfg Config
	
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	
	return &cfg, nil
}

// NewConfig cria nova instância de configuração
func NewConfig() (*Config, error) {
	return Load()
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
