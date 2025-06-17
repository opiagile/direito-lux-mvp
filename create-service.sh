#!/bin/bash

# Script para criar um novo microserviço seguindo os padrões do Direito Lux

set -e

if [ $# -ne 2 ]; then
    echo "Uso: $0 <service-name> <service-port>"
    echo "Exemplo: $0 notification-service 8085"
    exit 1
fi

SERVICE_NAME=$1
SERVICE_PORT=$2
SERVICE_TITLE=$(echo $SERVICE_NAME | sed 's/-/ /g' | sed 's/\b\w/\U&/g' | sed 's/ //g')

echo "🚀 Criando novo microserviço: $SERVICE_NAME"
echo "📋 Detalhes:"
echo "   Nome: $SERVICE_NAME"
echo "   Porta: $SERVICE_PORT" 
echo "   Título: $SERVICE_TITLE"

# Verificar se serviço já existe
if [ -d "services/$SERVICE_NAME" ]; then
    echo "❌ ERRO: Serviço $SERVICE_NAME já existe!"
    exit 1
fi

# Criar estrutura de diretórios
echo "📁 Criando estrutura de diretórios..."
mkdir -p "services/$SERVICE_NAME"
cd "services/$SERVICE_NAME"

# Estrutura padrão
mkdir -p cmd/server
mkdir -p internal/{application,domain,infrastructure/{config,database,events,http/{handlers,middleware},logging,metrics,tracing}}
mkdir -p migrations
mkdir -p docs

# Criar go.mod a partir do template
echo "📦 Criando go.mod..."
sed "s/{{SERVICE_NAME}}/$SERVICE_NAME/g" ../../templates/service-template/go.mod.template > go.mod

# Criar main.go a partir do template
echo "🔧 Criando main.go..."
sed -e "s/{{SERVICE_NAME}}/$SERVICE_NAME/g" \
    -e "s/{{SERVICE_PORT}}/$SERVICE_PORT/g" \
    -e "s/{{SERVICE_TITLE}}/$SERVICE_TITLE/g" \
    ../../templates/service-template/main.go.template > cmd/server/main.go

# Criar config básico
echo "⚙️  Criando configuração básica..."
cat > internal/infrastructure/config/config.go << 'EOF'
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
EOF

# Criar logging básico
echo "📝 Criando logging básico..."
cat > internal/infrastructure/logging/logger.go << 'EOF'
package logging

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ContextKey tipo para chaves de contexto
type ContextKey string

const (
	TraceIDKey    ContextKey = "trace_id"
	TenantIDKey   ContextKey = "tenant_id"  
	UserIDKey     ContextKey = "user_id"
	OperationKey  ContextKey = "operation"
)

// Logger wrapper para zap logger
type Logger struct {
	*zap.Logger
}

// NewLogger cria novo logger configurado
func NewLogger(level string, environment string) (*zap.Logger, error) {
	var config zap.Config
	
	switch environment {
	case "production":
		config = zap.NewProductionConfig()
		config.DisableStacktrace = true
	case "development":
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config = zap.NewProductionConfig()
	}

	// Configurar nível de log
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "fatal":
		zapLevel = zapcore.FatalLevel
	default:
		zapLevel = zapcore.InfoLevel
	}
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar logger: %w", err)
	}

	return logger, nil
}

// LogError registra erro com contexto
func LogError(ctx context.Context, logger *zap.Logger, message string, err error, fields ...zap.Field) {
	contextFields := extractContextFields(ctx)
	allFields := append(contextFields, zap.Error(err))
	allFields = append(allFields, fields...)
	logger.Error(message, allFields...)
}

// LogInfo registra informação com contexto
func LogInfo(ctx context.Context, logger *zap.Logger, message string, fields ...zap.Field) {
	contextFields := extractContextFields(ctx)
	allFields := append(contextFields, fields...)
	logger.Info(message, allFields...)
}

// extractContextFields extrai campos do contexto
func extractContextFields(ctx context.Context) []zap.Field {
	var fields []zap.Field

	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		if id, ok := traceID.(string); ok {
			fields = append(fields, zap.String("trace_id", id))
		}
	}

	if tenantID := ctx.Value(TenantIDKey); tenantID != nil {
		if id, ok := tenantID.(string); ok {
			fields = append(fields, zap.String("tenant_id", id))
		}
	}

	return fields
}
EOF

# Criar middleware básico
echo "🔒 Criando middleware básico..."
cat > internal/infrastructure/http/middleware/middleware.go << 'EOF'
package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/{{SERVICE_NAME}}/internal/infrastructure/config"
	"github.com/direito-lux/{{SERVICE_NAME}}/internal/infrastructure/logging"
)

// Logger middleware para logging de requisições
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return gin.LoggerWithWriter(os.Stdout)
}

// Recovery middleware para recuperação de panics
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.Error("Panic recuperado",
				zap.String("error", err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// RequestID middleware para adicionar ID único às requisições
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Header("X-Request-ID", requestID)
		
		// Adicionar ao contexto
		ctx := logging.WithOperation(c.Request.Context(), c.Request.Method+" "+c.Request.URL.Path)
		c.Request = c.Request.WithContext(ctx)
		
		c.Next()
	}
}

// CORS middleware
func CORS(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Em desenvolvimento, permitir qualquer origem
		allowedOrigin := "*"
		if !cfg.IsDevelopment() {
			// Em produção, verificar origens permitidas
			allowedOrigin = "https://app.direitolux.com"
		}

		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
EOF

# Substituir placeholder no middleware
sed -i '' "s/{{SERVICE_NAME}}/$SERVICE_NAME/g" internal/infrastructure/http/middleware/middleware.go

# Criar health check básico
echo "🏥 Criando health check..."
mkdir -p internal/infrastructure/http/handlers
cat > internal/infrastructure/http/handlers/health.go << 'EOF'
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler handler para health checks
type HealthHandler struct{}

// NewHealthHandler cria novo health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health retorna status de saúde do serviço
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().Unix(),
		"mode":   "full",
	})
}

// Ready retorna se o serviço está pronto
func (h *HealthHandler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"time":   time.Now().Unix(),
	})
}
EOF

echo "✅ Serviço $SERVICE_NAME criado com sucesso!"
echo ""
echo "📋 Próximos passos:"
echo "1. cd services/$SERVICE_NAME"
echo "2. go mod tidy"
echo "3. go build ./cmd/server"
echo "4. Implementar domínio e aplicação específicos"
echo "5. Adicionar ao start-services.sh"
echo ""
echo "📚 Consulte DIRETRIZES_DESENVOLVIMENTO.md para mais informações"