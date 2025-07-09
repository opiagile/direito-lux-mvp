package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/direito-lux/datajud-service/internal/domain"
)

// Config estrutura de configuração do serviço
type Config struct {
	// Aplicação
	ServiceName string `envconfig:"SERVICE_NAME" default:"datajud-service"`
	Version     string `envconfig:"VERSION" default:"1.0.0"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
	Port        int    `envconfig:"PORT" default:"8080"`
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
	
	// DataJud API
	DataJud DataJudConfig
}

// DatabaseConfig configuração do banco de dados
type DatabaseConfig struct {
	Host            string        `envconfig:"DB_HOST" default:"localhost"`
	Port            int           `envconfig:"DB_PORT" default:"5432"`
	Name            string        `envconfig:"DB_NAME" default:"direito_lux_dev"`
	User            string        `envconfig:"DB_USER" default:"direito_lux"`
	Password        string        `envconfig:"DB_PASSWORD" required:"true"`
	SSLMode         string        `envconfig:"DB_SSL_MODE" default:"disable"`
	MaxOpenConns    int           `envconfig:"DB_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int           `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
	ConnMaxLifetime time.Duration `envconfig:"DB_CONN_MAX_LIFETIME" default:"300s"`
	MigrationsPath  string        `envconfig:"DB_MIGRATIONS_PATH" default:"file://migrations"`
}

// RedisConfig configuração do Redis
type RedisConfig struct {
	Host         string        `envconfig:"REDIS_HOST" default:"localhost"`
	Port         int           `envconfig:"REDIS_PORT" default:"6379"`
	Password     string        `envconfig:"REDIS_PASSWORD"`
	Database     int           `envconfig:"REDIS_DATABASE" default:"0"`
	PoolSize     int           `envconfig:"REDIS_POOL_SIZE" default:"10"`
	MinIdleConns int           `envconfig:"REDIS_MIN_IDLE_CONNS" default:"5"`
	DialTimeout  time.Duration `envconfig:"REDIS_DIAL_TIMEOUT" default:"5s"`
	ReadTimeout  time.Duration `envconfig:"REDIS_READ_TIMEOUT" default:"3s"`
	WriteTimeout time.Duration `envconfig:"REDIS_WRITE_TIMEOUT" default:"3s"`
}

// RabbitMQConfig configuração do RabbitMQ
type RabbitMQConfig struct {
	URL          string        `envconfig:"RABBITMQ_URL" required:"true"`
	Exchange     string        `envconfig:"RABBITMQ_EXCHANGE" default:"direito_lux.events"`
	RoutingKey   string        `envconfig:"RABBITMQ_ROUTING_KEY" default:"template"`
	Queue        string        `envconfig:"RABBITMQ_QUEUE" default:"template.events"`
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
	ServiceName string  `envconfig:"JAEGER_SERVICE_NAME" default:"datajud-service"`
	SamplerType string  `envconfig:"JAEGER_SAMPLER_TYPE" default:"const"`
	SamplerParam float64 `envconfig:"JAEGER_SAMPLER_PARAM" default:"1"`
	Enabled     bool    `envconfig:"JAEGER_ENABLED" default:"true"`
}

// MetricsConfig configuração de métricas
type MetricsConfig struct {
	Enabled   bool   `envconfig:"METRICS_ENABLED" default:"true"`
	Port      int    `envconfig:"METRICS_PORT" default:"9090"`
	Path      string `envconfig:"METRICS_PATH" default:"/metrics"`
	Namespace string `envconfig:"METRICS_NAMESPACE" default:"template_service"`
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

// ExternalServicesConfig configuração de serviços externos
type ExternalServicesConfig struct {
	AuthServiceURL         string        `envconfig:"AUTH_SERVICE_URL" default:"http://auth-service:8080"`
	TenantServiceURL       string        `envconfig:"TENANT_SERVICE_URL" default:"http://tenant-service:8080"`
	ProcessServiceURL      string        `envconfig:"PROCESS_SERVICE_URL" default:"http://process-service:8080"`
	DataJudServiceURL      string        `envconfig:"DATAJUD_SERVICE_URL" default:"http://datajud-service:8080"`
	NotificationServiceURL string        `envconfig:"NOTIFICATION_SERVICE_URL" default:"http://notification-service:8080"`
	AIServiceURL           string        `envconfig:"AI_SERVICE_URL" default:"http://ai-service:8000"`
	
	// Timeouts
	DefaultTimeout time.Duration `envconfig:"EXTERNAL_SERVICE_TIMEOUT" default:"30s"`
	
	// Circuit Breaker
	CircuitBreakerEnabled         bool          `envconfig:"CIRCUIT_BREAKER_ENABLED" default:"true"`
	CircuitBreakerMaxRequests     uint32        `envconfig:"CIRCUIT_BREAKER_MAX_REQUESTS" default:"3"`
	CircuitBreakerInterval        time.Duration `envconfig:"CIRCUIT_BREAKER_INTERVAL" default:"60s"`
	CircuitBreakerTimeout         time.Duration `envconfig:"CIRCUIT_BREAKER_TIMEOUT" default:"30s"`
	CircuitBreakerFailureThreshold uint32       `envconfig:"CIRCUIT_BREAKER_FAILURE_THRESHOLD" default:"5"`
}

// DataJudConfig configuração da API DataJud CNJ
type DataJudConfig struct {
	// URL base da API pública
	BaseURL string `envconfig:"DATAJUD_BASE_URL" default:"https://api-publica.datajud.cnj.jus.br"`
	
	// Chave de API pública do CNJ
	APIKey string `envconfig:"DATAJUD_API_KEY" required:"true"`
	
	// Timeouts
	Timeout        time.Duration `envconfig:"DATAJUD_TIMEOUT" default:"30s"`
	ReadTimeout    time.Duration `envconfig:"DATAJUD_READ_TIMEOUT" default:"25s"`
	WriteTimeout   time.Duration `envconfig:"DATAJUD_WRITE_TIMEOUT" default:"25s"`
	IdleTimeout    time.Duration `envconfig:"DATAJUD_IDLE_TIMEOUT" default:"90s"`
	
	// Retry configuration
	RetryCount    int           `envconfig:"DATAJUD_RETRY_COUNT" default:"3"`
	RetryDelay    time.Duration `envconfig:"DATAJUD_RETRY_DELAY" default:"1s"`
	MaxRetryDelay time.Duration `envconfig:"DATAJUD_MAX_RETRY_DELAY" default:"30s"`
	
	// Rate limiting
	RateLimitEnabled    bool          `envconfig:"DATAJUD_RATE_LIMIT_ENABLED" default:"true"`
	RateLimitRPM        int           `envconfig:"DATAJUD_RATE_LIMIT_RPM" default:"120"`
	RateLimitBurst      int           `envconfig:"DATAJUD_RATE_LIMIT_BURST" default:"10"`
	RateLimitWindow     time.Duration `envconfig:"DATAJUD_RATE_LIMIT_WINDOW" default:"1m"`
	
	// Circuit breaker
	CircuitBreakerEnabled         bool          `envconfig:"DATAJUD_CIRCUIT_BREAKER_ENABLED" default:"true"`
	CircuitBreakerMaxRequests     uint32        `envconfig:"DATAJUD_CIRCUIT_BREAKER_MAX_REQUESTS" default:"3"`
	CircuitBreakerInterval        time.Duration `envconfig:"DATAJUD_CIRCUIT_BREAKER_INTERVAL" default:"60s"`
	CircuitBreakerTimeout         time.Duration `envconfig:"DATAJUD_CIRCUIT_BREAKER_TIMEOUT" default:"30s"`
	CircuitBreakerFailureThreshold uint32       `envconfig:"DATAJUD_CIRCUIT_BREAKER_FAILURE_THRESHOLD" default:"5"`
	
	// Cache
	CacheEnabled      bool          `envconfig:"DATAJUD_CACHE_ENABLED" default:"true"`
	CacheDefaultTTL   time.Duration `envconfig:"DATAJUD_CACHE_DEFAULT_TTL" default:"1h"`
	CacheProcessTTL   time.Duration `envconfig:"DATAJUD_CACHE_PROCESS_TTL" default:"24h"`
	CacheMovementTTL  time.Duration `envconfig:"DATAJUD_CACHE_MOVEMENT_TTL" default:"30m"`
	CacheMaxSize      int           `envconfig:"DATAJUD_CACHE_MAX_SIZE" default:"10000"`
	CacheMaxMemory    int64         `envconfig:"DATAJUD_CACHE_MAX_MEMORY" default:"104857600"` // 100MB
	
	// Pool de CNPJs (legacy, será removido)
	CNPJPoolEnabled   bool   `envconfig:"DATAJUD_CNPJ_POOL_ENABLED" default:"false"`
	CNPJPoolSize      int    `envconfig:"DATAJUD_CNPJ_POOL_SIZE" default:"1"`
	CNPJPoolStrategy  string `envconfig:"DATAJUD_CNPJ_POOL_STRATEGY" default:"round_robin"`
	
	// Monitoring
	MetricsEnabled       bool   `envconfig:"DATAJUD_METRICS_ENABLED" default:"true"`
	LogRequestsEnabled   bool   `envconfig:"DATAJUD_LOG_REQUESTS_ENABLED" default:"false"`
	LogResponsesEnabled  bool   `envconfig:"DATAJUD_LOG_RESPONSES_ENABLED" default:"false"`
	
	// Elasticsearch optimization
	ESScrollSize      int           `envconfig:"DATAJUD_ES_SCROLL_SIZE" default:"1000"`
	ESScrollTimeout   time.Duration `envconfig:"DATAJUD_ES_SCROLL_TIMEOUT" default:"5m"`
	ESMaxResultWindow int           `envconfig:"DATAJUD_ES_MAX_RESULT_WINDOW" default:"10000"`
	
	// User agent
	UserAgent string `envconfig:"DATAJUD_USER_AGENT" default:"Direito-Lux/1.0"`
	
	// Fallback para mock em desenvolvimento
	MockEnabled bool `envconfig:"DATAJUD_MOCK_ENABLED" default:"false"`
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

	// Validar configurações do DataJud
	if err := c.ValidateDataJudConfig(); err != nil {
		return fmt.Errorf("erro na configuração DataJud: %w", err)
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

// ValidateDataJudConfig valida configurações específicas do DataJud
func (c *Config) ValidateDataJudConfig() error {
	// Validar API Key
	if c.DataJud.APIKey == "" {
		return fmt.Errorf("DATAJUD_API_KEY é obrigatória")
	}
	
	// Validar URL base
	if c.DataJud.BaseURL == "" {
		return fmt.Errorf("DATAJUD_BASE_URL é obrigatória")
	}
	
	// Validar timeouts
	if c.DataJud.Timeout <= 0 {
		return fmt.Errorf("DATAJUD_TIMEOUT deve ser maior que 0")
	}
	
	if c.DataJud.RetryCount < 0 {
		return fmt.Errorf("DATAJUD_RETRY_COUNT deve ser maior ou igual a 0")
	}
	
	if c.DataJud.RetryCount > 10 {
		return fmt.Errorf("DATAJUD_RETRY_COUNT não deve ser maior que 10")
	}
	
	// Validar rate limiting
	if c.DataJud.RateLimitEnabled {
		if c.DataJud.RateLimitRPM <= 0 {
			return fmt.Errorf("DATAJUD_RATE_LIMIT_RPM deve ser maior que 0")
		}
		
		if c.DataJud.RateLimitBurst <= 0 {
			return fmt.Errorf("DATAJUD_RATE_LIMIT_BURST deve ser maior que 0")
		}
	}
	
	// Validar cache
	if c.DataJud.CacheEnabled {
		if c.DataJud.CacheMaxSize <= 0 {
			return fmt.Errorf("DATAJUD_CACHE_MAX_SIZE deve ser maior que 0")
		}
		
		if c.DataJud.CacheMaxMemory <= 0 {
			return fmt.Errorf("DATAJUD_CACHE_MAX_MEMORY deve ser maior que 0")
		}
	}
	
	return nil
}

// IsDataJudMockEnabled verifica se está usando mock
func (c *Config) IsDataJudMockEnabled() bool {
	return c.DataJud.MockEnabled
}

// GetDataJudClientConfig retorna configuração do cliente DataJud
func (c *Config) GetDataJudClientConfig() DataJudClientConfig {
	return DataJudClientConfig{
		BaseURL:    c.DataJud.BaseURL,
		APIKey:     c.DataJud.APIKey,
		Timeout:    c.DataJud.Timeout,
		RetryCount: c.DataJud.RetryCount,
		RetryDelay: c.DataJud.RetryDelay,
		UserAgent:  c.DataJud.UserAgent,
		MockEnabled: c.IsDataJudMockEnabled(),
	}
}

// DataJudClientConfig configuração simplificada para o cliente
type DataJudClientConfig struct {
	BaseURL     string
	APIKey      string
	Timeout     time.Duration
	RetryCount  int
	RetryDelay  time.Duration
	UserAgent   string
	MockEnabled bool
}

// GetDataJudRateLimitConfig retorna configuração de rate limiting
func (c *Config) GetDataJudRateLimitConfig() DataJudRateLimitConfig {
	return DataJudRateLimitConfig{
		Enabled:    c.DataJud.RateLimitEnabled,
		RPM:        c.DataJud.RateLimitRPM,
		Burst:      c.DataJud.RateLimitBurst,
		Window:     c.DataJud.RateLimitWindow,
	}
}

// DataJudRateLimitConfig configuração de rate limiting
type DataJudRateLimitConfig struct {
	Enabled bool
	RPM     int
	Burst   int
	Window  time.Duration
}

// GetDataJudCacheConfig retorna configuração de cache
func (c *Config) GetDataJudCacheConfig() DataJudCacheConfig {
	return DataJudCacheConfig{
		Enabled:     c.DataJud.CacheEnabled,
		DefaultTTL:  c.DataJud.CacheDefaultTTL,
		ProcessTTL:  c.DataJud.CacheProcessTTL,
		MovementTTL: c.DataJud.CacheMovementTTL,
		MaxSize:     c.DataJud.CacheMaxSize,
		MaxMemory:   c.DataJud.CacheMaxMemory,
	}
}

// DataJudCacheConfig configuração de cache
type DataJudCacheConfig struct {
	Enabled     bool
	DefaultTTL  time.Duration
	ProcessTTL  time.Duration
	MovementTTL time.Duration
	MaxSize     int
	MaxMemory   int64
}

// GetDataJudCircuitBreakerConfig retorna configuração de circuit breaker
func (c *Config) GetDataJudCircuitBreakerConfig() DataJudCircuitBreakerConfig {
	return DataJudCircuitBreakerConfig{
		Enabled:          c.DataJud.CircuitBreakerEnabled,
		MaxRequests:      c.DataJud.CircuitBreakerMaxRequests,
		Interval:         c.DataJud.CircuitBreakerInterval,
		Timeout:          c.DataJud.CircuitBreakerTimeout,
		FailureThreshold: c.DataJud.CircuitBreakerFailureThreshold,
	}
}

// GetDataJudDomainConfig retorna configuração para o domínio DataJud
func (c *Config) GetDataJudDomainConfig() domain.DataJudConfig {
	return domain.DataJudConfig{
		APIBaseURL:           c.DataJud.BaseURL,
		APITimeout:           c.DataJud.Timeout,
		APIRetryCount:        c.DataJud.RetryCount,
		APIRetryDelay:        c.DataJud.RetryDelay,
		DefaultDailyLimit:    10000, // Limite padrão CNJ
		GlobalRateLimit:      c.DataJud.RateLimitRPM,
		RateWindowSize:       c.DataJud.RateLimitWindow,
		DefaultCacheTTL:      int(c.DataJud.CacheDefaultTTL.Seconds()),
		MaxCacheSize:         c.DataJud.CacheMaxMemory,
		MaxCacheEntries:      int64(c.DataJud.CacheMaxSize),
		CacheCleanupInterval: 5 * time.Minute,
		CBFailureThreshold:   int(c.DataJud.CircuitBreakerFailureThreshold),
		CBSuccessThreshold:   3,
		CBTimeout:            c.DataJud.CircuitBreakerTimeout,
		CBMaxRequests:        int(c.DataJud.CircuitBreakerMaxRequests),
		DefaultPoolStrategy:  domain.StrategyRoundRobin,
		MetricsEnabled:       c.DataJud.MetricsEnabled,
	}
}

// DataJudCircuitBreakerConfig configuração de circuit breaker
type DataJudCircuitBreakerConfig struct {
	Enabled          bool
	MaxRequests      uint32
	Interval         time.Duration
	Timeout          time.Duration
	FailureThreshold uint32
}

