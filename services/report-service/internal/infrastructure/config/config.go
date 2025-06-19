package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config representa a configuração do serviço
type Config struct {
	// Application
	ServiceName string `envconfig:"SERVICE_NAME" default:"report-service"`
	Version     string `envconfig:"VERSION" default:"1.0.0"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"info"`

	// Server
	Server ServerConfig

	// Database
	Database DatabaseConfig

	// Redis
	Redis RedisConfig

	// File Storage
	Storage StorageConfig

	// Report Generation
	Report ReportConfig

	// Tracing
	Tracing TracingConfig

	// Metrics
	Metrics MetricsConfig

	// External Services
	Services ServicesConfig
}

// ServerConfig configurações do servidor HTTP
type ServerConfig struct {
	Port         int           `envconfig:"SERVER_PORT" default:"8087"`
	ReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"30s"`
	WriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"30s"`
	IdleTimeout  time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" default:"120s"`
}

// DatabaseConfig configurações do banco de dados
type DatabaseConfig struct {
	Host            string        `envconfig:"DB_HOST" default:"localhost"`
	Port            int           `envconfig:"DB_PORT" default:"5432"`
	User            string        `envconfig:"DB_USER" default:"postgres"`
	Password        string        `envconfig:"DB_PASSWORD" required:"true"`
	Name            string        `envconfig:"DB_NAME" default:"direito_lux_dev"`
	SSLMode         string        `envconfig:"DB_SSL_MODE" default:"disable"`
	MaxOpenConns    int           `envconfig:"DB_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int           `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
	ConnMaxLifetime time.Duration `envconfig:"DB_CONN_MAX_LIFETIME" default:"5m"`
}

// RedisConfig configurações do Redis
type RedisConfig struct {
	Host         string        `envconfig:"REDIS_HOST" default:"localhost"`
	Port         int           `envconfig:"REDIS_PORT" default:"6379"`
	Password     string        `envconfig:"REDIS_PASSWORD"`
	DB           int           `envconfig:"REDIS_DB" default:"0"`
	PoolSize     int           `envconfig:"REDIS_POOL_SIZE" default:"10"`
	MaxRetries   int           `envconfig:"REDIS_MAX_RETRIES" default:"3"`
	DialTimeout  time.Duration `envconfig:"REDIS_DIAL_TIMEOUT" default:"5s"`
	ReadTimeout  time.Duration `envconfig:"REDIS_READ_TIMEOUT" default:"3s"`
	WriteTimeout time.Duration `envconfig:"REDIS_WRITE_TIMEOUT" default:"3s"`
	TTL          time.Duration `envconfig:"REDIS_TTL" default:"1h"`
}

// StorageConfig configurações de armazenamento de arquivos
type StorageConfig struct {
	Type           string        `envconfig:"STORAGE_TYPE" default:"local"` // local, gcs, s3
	LocalPath      string        `envconfig:"STORAGE_LOCAL_PATH" default:"./reports"`
	GCSBucket      string        `envconfig:"STORAGE_GCS_BUCKET"`
	GCSProject     string        `envconfig:"STORAGE_GCS_PROJECT"`
	S3Bucket       string        `envconfig:"STORAGE_S3_BUCKET"`
	S3Region       string        `envconfig:"STORAGE_S3_REGION"`
	MaxFileSize    int64         `envconfig:"STORAGE_MAX_FILE_SIZE" default:"104857600"` // 100MB
	RetentionDays  int           `envconfig:"STORAGE_RETENTION_DAYS" default:"30"`
	URLExpiration  time.Duration `envconfig:"STORAGE_URL_EXPIRATION" default:"24h"`
}

// ReportConfig configurações de geração de relatórios
type ReportConfig struct {
	MaxConcurrent       int           `envconfig:"REPORT_MAX_CONCURRENT" default:"10"`
	DefaultTimeout      time.Duration `envconfig:"REPORT_DEFAULT_TIMEOUT" default:"5m"`
	MaxRetries          int           `envconfig:"REPORT_MAX_RETRIES" default:"3"`
	RetryDelay          time.Duration `envconfig:"REPORT_RETRY_DELAY" default:"30s"`
	PDFTimeout          time.Duration `envconfig:"REPORT_PDF_TIMEOUT" default:"2m"`
	ExcelTimeout        time.Duration `envconfig:"REPORT_EXCEL_TIMEOUT" default:"3m"`
	MaxRowsPerSheet     int           `envconfig:"REPORT_MAX_ROWS_PER_SHEET" default:"1000000"`
	ChartCacheTTL       time.Duration `envconfig:"REPORT_CHART_CACHE_TTL" default:"5m"`
	SchedulerInterval   time.Duration `envconfig:"REPORT_SCHEDULER_INTERVAL" default:"1m"`
	CleanupInterval     time.Duration `envconfig:"REPORT_CLEANUP_INTERVAL" default:"1h"`
	
	// Plan limits
	StarterMonthlyLimit      int `envconfig:"REPORT_STARTER_MONTHLY_LIMIT" default:"10"`
	ProfessionalMonthlyLimit int `envconfig:"REPORT_PROFESSIONAL_MONTHLY_LIMIT" default:"100"`
	BusinessMonthlyLimit     int `envconfig:"REPORT_BUSINESS_MONTHLY_LIMIT" default:"500"`
	// Enterprise is unlimited
}

// TracingConfig configurações de tracing
type TracingConfig struct {
	Enabled        bool    `envconfig:"TRACING_ENABLED" default:"true"`
	JaegerEndpoint string  `envconfig:"JAEGER_ENDPOINT" default:"http://localhost:14268/api/traces"`
	SampleRate     float64 `envconfig:"TRACING_SAMPLE_RATE" default:"1.0"`
}

// MetricsConfig configurações de métricas
type MetricsConfig struct {
	Enabled bool   `envconfig:"METRICS_ENABLED" default:"true"`
	Port    int    `envconfig:"METRICS_PORT" default:"9097"`
	Path    string `envconfig:"METRICS_PATH" default:"/metrics"`
}

// ServicesConfig configurações de serviços externos
type ServicesConfig struct {
	ProcessServiceURL       string `envconfig:"PROCESS_SERVICE_URL" default:"http://localhost:8083"`
	TenantServiceURL        string `envconfig:"TENANT_SERVICE_URL" default:"http://localhost:8082"`
	NotificationServiceURL  string `envconfig:"NOTIFICATION_SERVICE_URL" default:"http://localhost:8085"`
	AIServiceURL            string `envconfig:"AI_SERVICE_URL" default:"http://localhost:8000"`
	SearchServiceURL        string `envconfig:"SEARCH_SERVICE_URL" default:"http://localhost:8086"`
}

// Load carrega a configuração a partir de variáveis de ambiente
func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// GetDatabaseURL retorna a URL de conexão do banco
func (c *Config) GetDatabaseURL() string {
	return "postgres://" + c.Database.User + ":" + c.Database.Password +
		"@" + c.Database.Host + ":" + string(rune(c.Database.Port)) +
		"/" + c.Database.Name + "?sslmode=" + c.Database.SSLMode
}

// GetRedisAddr retorna o endereço do Redis
func (c *Config) GetRedisAddr() string {
	return c.Redis.Host + ":" + string(rune(c.Redis.Port))
}

// IsDevelopment verifica se está em ambiente de desenvolvimento
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction verifica se está em ambiente de produção
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// GetReportLimitByPlan retorna o limite mensal de relatórios por plano
func (c *Config) GetReportLimitByPlan(plan string) int {
	switch plan {
	case "starter":
		return c.Report.StarterMonthlyLimit
	case "professional":
		return c.Report.ProfessionalMonthlyLimit
	case "business":
		return c.Report.BusinessMonthlyLimit
	case "enterprise":
		return -1 // unlimited
	default:
		return 0
	}
}