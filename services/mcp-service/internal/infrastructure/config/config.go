package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config estrutura de configuração do MCP Service
type Config struct {
	// Aplicação
	ServiceName string `envconfig:"SERVICE_NAME" default:"mcp-service"`
	Version     string `envconfig:"VERSION" default:"1.0.0"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
	Port        int    `envconfig:"PORT" default:"8088"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"info"`

	// Database
	Database DatabaseConfig

	// Redis
	Redis RedisConfig

	// RabbitMQ
	RabbitMQ RabbitMQConfig

	// Tracing
	Jaeger JaegerConfig

	// Metrics
	Metrics MetricsConfig

	// HTTP
	HTTP HTTPConfig

	// External Services
	Services ExternalServicesConfig

	// Claude API Configuration
	Claude ClaudeConfig

	// Bot Configurations
	WhatsApp WhatsAppConfig
	Telegram TelegramConfig
	Slack    SlackConfig

	// MCP Specific
	MCP MCPConfig
}

// DatabaseConfig configuração do banco de dados
type DatabaseConfig struct {
	Host            string        `envconfig:"DB_HOST" default:"localhost"`
	Port            int           `envconfig:"DB_PORT" default:"5434"`
	Name            string        `envconfig:"DB_NAME" default:"direito_lux_mcp"`
	User            string        `envconfig:"DB_USER" default:"mcp_user"`
	Password        string        `envconfig:"DB_PASSWORD" default:"mcp_pass_dev"`
	SSLMode         string        `envconfig:"DB_SSL_MODE" default:"disable"`
	MaxOpenConns    int           `envconfig:"DB_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int           `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
	ConnMaxLifetime time.Duration `envconfig:"DB_CONN_MAX_LIFETIME" default:"300s"`
	MigrationsPath  string        `envconfig:"DB_MIGRATIONS_PATH" default:"file://migrations"`
}

// RedisConfig configuração do Redis
type RedisConfig struct {
	Host         string        `envconfig:"REDIS_HOST" default:"localhost"`
	Port         int           `envconfig:"REDIS_PORT" default:"6380"`
	Password     string        `envconfig:"REDIS_PASSWORD" default:"redis_pass_dev"`
	Database     int           `envconfig:"REDIS_DATABASE" default:"0"`
	PoolSize     int           `envconfig:"REDIS_POOL_SIZE" default:"10"`
	MinIdleConns int           `envconfig:"REDIS_MIN_IDLE_CONNS" default:"5"`
	DialTimeout  time.Duration `envconfig:"REDIS_DIAL_TIMEOUT" default:"5s"`
	ReadTimeout  time.Duration `envconfig:"REDIS_READ_TIMEOUT" default:"3s"`
	WriteTimeout time.Duration `envconfig:"REDIS_WRITE_TIMEOUT" default:"3s"`
}

// RabbitMQConfig configuração do RabbitMQ
type RabbitMQConfig struct {
	URL          string        `envconfig:"RABBITMQ_URL" default:"amqp://mcp_user:rabbit_pass_dev@localhost:5673/mcp_vhost"`
	Exchange     string        `envconfig:"RABBITMQ_EXCHANGE" default:"mcp.events"`
	RoutingKey   string        `envconfig:"RABBITMQ_ROUTING_KEY" default:"mcp"`
	Queue        string        `envconfig:"RABBITMQ_QUEUE" default:"mcp.events"`
	Durable      bool          `envconfig:"RABBITMQ_DURABLE" default:"true"`
	AutoDelete   bool          `envconfig:"RABBITMQ_AUTO_DELETE" default:"false"`
	Exclusive    bool          `envconfig:"RABBITMQ_EXCLUSIVE" default:"false"`
	NoWait       bool          `envconfig:"RABBITMQ_NO_WAIT" default:"false"`
	PrefetchCount int          `envconfig:"RABBITMQ_PREFETCH_COUNT" default:"10"`
	ReconnectDelay time.Duration `envconfig:"RABBITMQ_RECONNECT_DELAY" default:"5s"`
}

// JaegerConfig configuração do Jaeger
type JaegerConfig struct {
	Endpoint    string  `envconfig:"JAEGER_ENDPOINT" default:"http://localhost:14268/api/traces"`
	ServiceName string  `envconfig:"JAEGER_SERVICE_NAME" default:"mcp-service"`
	SamplerType string  `envconfig:"JAEGER_SAMPLER_TYPE" default:"const"`
	SamplerParam float64 `envconfig:"JAEGER_SAMPLER_PARAM" default:"1"`
	Enabled     bool    `envconfig:"JAEGER_ENABLED" default:"true"`
}

// MetricsConfig configuração de métricas
type MetricsConfig struct {
	Enabled   bool   `envconfig:"METRICS_ENABLED" default:"true"`
	Port      int    `envconfig:"METRICS_PORT" default:"9090"`
	Path      string `envconfig:"METRICS_PATH" default:"/metrics"`
	Namespace string `envconfig:"METRICS_NAMESPACE" default:"mcp_service"`
}

// HTTPConfig configuração do servidor HTTP
type HTTPConfig struct {
	ReadTimeout       time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"10s"`
	WriteTimeout      time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"10s"`
	IdleTimeout       time.Duration `envconfig:"HTTP_IDLE_TIMEOUT" default:"60s"`
	ReadHeaderTimeout time.Duration `envconfig:"HTTP_READ_HEADER_TIMEOUT" default:"5s"`
	MaxHeaderBytes    int           `envconfig:"HTTP_MAX_HEADER_BYTES" default:"1048576"`
	
	// CORS
	CORSAllowedOrigins []string `envconfig:"CORS_ALLOWED_ORIGINS" default:"*"`
	CORSAllowedMethods []string `envconfig:"CORS_ALLOWED_METHODS" default:"GET,POST,PUT,DELETE,OPTIONS"`
	CORSAllowedHeaders []string `envconfig:"CORS_ALLOWED_HEADERS" default:"Content-Type,Authorization,X-Tenant-ID"`
	
	// Rate Limiting
	RateLimitEnabled bool `envconfig:"RATE_LIMIT_ENABLED" default:"true"`
	RateLimitRPM     int  `envconfig:"RATE_LIMIT_RPM" default:"100"`
	RateLimitBurst   int  `envconfig:"RATE_LIMIT_BURST" default:"200"`
}

// ClaudeConfig configuração da API Claude
type ClaudeConfig struct {
	APIKey     string        `envconfig:"ANTHROPIC_API_KEY" default:"sk-ant-api03-test-key"`
	Model      string        `envconfig:"ANTHROPIC_MODEL" default:"claude-3-5-sonnet-20241022"`
	MaxTokens  int           `envconfig:"ANTHROPIC_MAX_TOKENS" default:"4096"`
	Timeout    time.Duration `envconfig:"ANTHROPIC_TIMEOUT" default:"30s"`
	BaseURL    string        `envconfig:"ANTHROPIC_BASE_URL" default:"https://api.anthropic.com"`
	Version    string        `envconfig:"ANTHROPIC_VERSION" default:"2023-06-01"`
	MaxRetries int           `envconfig:"ANTHROPIC_MAX_RETRIES" default:"3"`
}

// WhatsAppConfig configuração do WhatsApp Business API
type WhatsAppConfig struct {
	AccessToken     string `envconfig:"WHATSAPP_ACCESS_TOKEN" default:"test-access-token"`
	PhoneNumberID   string `envconfig:"WHATSAPP_PHONE_NUMBER_ID" default:"test-phone-id"`
	VerifyToken     string `envconfig:"WHATSAPP_VERIFY_TOKEN" default:"test-verify-token"`
	WebhookURL      string `envconfig:"WHATSAPP_WEBHOOK_URL"`
	APIVersion      string `envconfig:"WHATSAPP_API_VERSION" default:"v17.0"`
	BaseURL         string `envconfig:"WHATSAPP_BASE_URL" default:"https://graph.facebook.com"`
	MaxMessageLength int   `envconfig:"WHATSAPP_MAX_MESSAGE_LENGTH" default:"4096"`
}

// TelegramConfig configuração do Telegram Bot
type TelegramConfig struct {
	BotToken    string `envconfig:"TELEGRAM_BOT_TOKEN" default:"test-bot-token"`
	WebhookURL  string `envconfig:"TELEGRAM_WEBHOOK_URL"`
	BaseURL     string `envconfig:"TELEGRAM_BASE_URL" default:"https://api.telegram.org"`
	MaxMessageLength int `envconfig:"TELEGRAM_MAX_MESSAGE_LENGTH" default:"4096"`
	ParseMode   string `envconfig:"TELEGRAM_PARSE_MODE" default:"Markdown"`
}

// SlackConfig configuração do Slack Bot
type SlackConfig struct {
	BotToken      string `envconfig:"SLACK_BOT_TOKEN"`
	SigningSecret string `envconfig:"SLACK_SIGNING_SECRET"`
	AppToken      string `envconfig:"SLACK_APP_TOKEN"`
	BaseURL       string `envconfig:"SLACK_BASE_URL" default:"https://slack.com/api"`
	MaxMessageLength int `envconfig:"SLACK_MAX_MESSAGE_LENGTH" default:"40000"`
}

// MCPConfig configurações específicas do MCP
type MCPConfig struct {
	CacheTTL              time.Duration `envconfig:"MCP_CACHE_TTL" default:"300s"`
	MaxConcurrentRequests int           `envconfig:"MCP_MAX_CONCURRENT_REQUESTS" default:"100"`
	RateLimitPerUser      int           `envconfig:"MCP_RATE_LIMIT_PER_USER" default:"60"`
	SessionTimeout        time.Duration `envconfig:"MCP_SESSION_TIMEOUT" default:"30m"`
	
	// Security
	JWTSecret       string   `envconfig:"MCP_JWT_SECRET" default:"mcp-jwt-secret-dev"`
	EncryptionKey   string   `envconfig:"MCP_ENCRYPTION_KEY" default:"mcp-encryption-key-dev"`
	AllowedOrigins  []string `envconfig:"MCP_ALLOWED_ORIGINS" default:"https://app.direitolux.com"`
	
	// Features
	ToolsEnabled          bool `envconfig:"MCP_TOOLS_ENABLED" default:"true"`
	ContextMemoryEnabled  bool `envconfig:"MCP_CONTEXT_MEMORY_ENABLED" default:"true"`
	QuotaEnforcementEnabled bool `envconfig:"MCP_QUOTA_ENFORCEMENT_ENABLED" default:"true"`
	
	// Monitoring
	MetricsEnabled bool `envconfig:"MCP_METRICS_ENABLED" default:"true"`
	LoggingLevel   string `envconfig:"MCP_LOGGING_LEVEL" default:"info"`
	TracingEnabled bool `envconfig:"MCP_TRACING_ENABLED" default:"true"`
}

// ExternalServicesConfig configuração de serviços externos
type ExternalServicesConfig struct {
	AuthServiceURL         string        `envconfig:"AUTH_SERVICE_URL" default:"http://auth-service:8080"`
	TenantServiceURL       string        `envconfig:"TENANT_SERVICE_URL" default:"http://tenant-service:8080"`
	ProcessServiceURL      string        `envconfig:"PROCESS_SERVICE_URL" default:"http://process-service:8080"`
	DataJudServiceURL      string        `envconfig:"DATAJUD_SERVICE_URL" default:"http://datajud-service:8080"`
	NotificationServiceURL string        `envconfig:"NOTIFICATION_SERVICE_URL" default:"http://notification-service:8080"`
	AIServiceURL           string        `envconfig:"AI_SERVICE_URL" default:"http://ai-service:8000"`
	SearchServiceURL       string        `envconfig:"SEARCH_SERVICE_URL" default:"http://search-service:8086"`
	
	// Timeouts
	DefaultTimeout time.Duration `envconfig:"EXTERNAL_SERVICE_TIMEOUT" default:"30s"`
	
	// Circuit Breaker
	CircuitBreakerEnabled         bool          `envconfig:"CIRCUIT_BREAKER_ENABLED" default:"true"`
	CircuitBreakerMaxRequests     uint32        `envconfig:"CIRCUIT_BREAKER_MAX_REQUESTS" default:"3"`
	CircuitBreakerInterval        time.Duration `envconfig:"CIRCUIT_BREAKER_INTERVAL" default:"60s"`
	CircuitBreakerTimeout         time.Duration `envconfig:"CIRCUIT_BREAKER_TIMEOUT" default:"30s"`
	CircuitBreakerFailureThreshold uint32       `envconfig:"CIRCUIT_BREAKER_FAILURE_THRESHOLD" default:"5"`
}

// Load carrega configurações das variáveis de ambiente
func Load() (*Config, error) {
	var cfg Config
	
	// Carregar configurações das variáveis de ambiente
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	// Validações específicas
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// New alias para Load() - mantém compatibilidade com código existente
func New() (*Config, error) {
	return Load()
}

// validate valida as configurações
func (c *Config) validate() error {
	// Validar porta
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("porta inválida: %d", c.Port)
	}

	// Validar nível de log
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("nível de log inválido: %s", c.LogLevel)
	}

	// Validar environment
	validEnvironments := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
		"test":        true,
	}
	
	if !validEnvironments[c.Environment] {
		return fmt.Errorf("environment inválido: %s", c.Environment)
	}

	// Validar modelo Claude
	validModels := map[string]bool{
		"claude-3-5-sonnet-20241022": true,
		"claude-3-opus-20240229":     true,
		"claude-3-haiku-20240307":    true,
	}
	
	if !validModels[c.Claude.Model] {
		return fmt.Errorf("modelo Claude inválido: %s", c.Claude.Model)
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

// GetDatabaseDSN retorna a string de conexão com o banco
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// GetRedisAddr retorna o endereço do Redis
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// GetClaudeHeaders retorna headers para API Claude
func (c *Config) GetClaudeHeaders() map[string]string {
	return map[string]string{
		"Content-Type":      "application/json",
		"x-api-key":         c.Claude.APIKey,
		"anthropic-version": c.Claude.Version,
	}
}

// GetWhatsAppBaseURL retorna URL base do WhatsApp
func (c *Config) GetWhatsAppBaseURL() string {
	return fmt.Sprintf("%s/%s", c.WhatsApp.BaseURL, c.WhatsApp.APIVersion)
}

// GetTelegramBaseURL retorna URL base do Telegram
func (c *Config) GetTelegramBaseURL() string {
	return fmt.Sprintf("%s/bot%s", c.Telegram.BaseURL, c.Telegram.BotToken)
}