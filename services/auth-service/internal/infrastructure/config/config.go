package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config estrutura de configuração do serviço
type Config struct {
	// Aplicação
	ServiceName string `envconfig:"SERVICE_NAME" default:"auth-service"`
	Version     string `envconfig:"VERSION" default:"1.0.0"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
	Port        int    `envconfig:"PORT" default:"8081"`
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
	
	// JWT Configuration
	JWT JWTConfig
	
	// Keycloak Configuration
	Keycloak KeycloakConfig
	
	// Password Security
	Security SecurityConfig
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
	RoutingKey   string        `envconfig:"RABBITMQ_ROUTING_KEY" default:"auth_service"`
	Queue        string        `envconfig:"RABBITMQ_QUEUE" default:"auth_service.events"`
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
	ServiceName string  `envconfig:"JAEGER_SERVICE_NAME" default:"auth-service"`
	SamplerType string  `envconfig:"JAEGER_SAMPLER_TYPE" default:"const"`
	SamplerParam float64 `envconfig:"JAEGER_SAMPLER_PARAM" default:"1"`
	Enabled     bool    `envconfig:"JAEGER_ENABLED" default:"true"`
}

// MetricsConfig configuração de métricas
type MetricsConfig struct {
	Enabled   bool   `envconfig:"METRICS_ENABLED" default:"true"`
	Port      int    `envconfig:"METRICS_PORT" default:"9090"`
	Path      string `envconfig:"METRICS_PATH" default:"/metrics"`
	Namespace string `envconfig:"METRICS_NAMESPACE" default:"auth_service"`
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

// JWTConfig configuração JWT
type JWTConfig struct {
	Secret           string        `envconfig:"JWT_SECRET" required:"true"`
	ExpiryHours      int           `envconfig:"JWT_EXPIRY_HOURS" default:"24"`
	RefreshExpiryDays int          `envconfig:"JWT_REFRESH_EXPIRY_DAYS" default:"30"`
	Issuer           string        `envconfig:"JWT_ISSUER" default:"direito-lux"`
	Audience         string        `envconfig:"JWT_AUDIENCE" default:"auth-service"`
}

// KeycloakConfig configuração Keycloak
type KeycloakConfig struct {
	URL          string `envconfig:"KEYCLOAK_URL" default:"http://keycloak:8080"`
	Realm        string `envconfig:"KEYCLOAK_REALM" default:"direito-lux"`
	ClientID     string `envconfig:"KEYCLOAK_CLIENT_ID" default:"auth-service"`
	ClientSecret string `envconfig:"KEYCLOAK_CLIENT_SECRET" required:"true"`
	Enabled      bool   `envconfig:"KEYCLOAK_ENABLED" default:"false"`
}

// SecurityConfig configuração de segurança
type SecurityConfig struct {
	BCryptCost           int  `envconfig:"BCRYPT_COST" default:"12"`
	PasswordMinLength    int  `envconfig:"PASSWORD_MIN_LENGTH" default:"8"`
	PasswordRequireSymbols    bool `envconfig:"PASSWORD_REQUIRE_SYMBOLS" default:"true"`
	PasswordRequireNumbers    bool `envconfig:"PASSWORD_REQUIRE_NUMBERS" default:"true"`
	PasswordRequireUppercase  bool `envconfig:"PASSWORD_REQUIRE_UPPERCASE" default:"true"`
	PasswordRequireLowercase  bool `envconfig:"PASSWORD_REQUIRE_LOWERCASE" default:"true"`
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