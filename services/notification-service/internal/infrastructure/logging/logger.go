package logging

import (
	"context"
	"fmt"

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

	if operation := ctx.Value(OperationKey); operation != nil {
		if op, ok := operation.(string); ok {
			fields = append(fields, zap.String("operation", op))
		}
	}

	return fields
}

// WithOperation adiciona operação ao contexto
func WithOperation(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, OperationKey, operation)
}

// WithTenantID adiciona tenant ID ao contexto
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

// WithTraceID adiciona trace ID ao contexto
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}
