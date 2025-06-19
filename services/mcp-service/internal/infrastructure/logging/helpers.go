package logging

import (
	"context"
	
	"go.uber.org/zap"
)

// Re-use the ContextKey type and constants from logger.go

// Helper functions para logging

// String cria um campo zap string
func String(key, value string) zap.Field {
	return zap.String(key, value)
}

// Int cria um campo zap int
func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// Float64 cria um campo zap float64
func Float64(key string, value float64) zap.Field {
	return zap.Float64(key, value)
}

// Bool cria um campo zap bool
func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

// Any cria um campo zap any
func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// Error cria um campo zap error
func Error(err error) zap.Field {
	return zap.Error(err)
}

// LogInfo logs an info message with context and fields
func LogInfo(ctx context.Context, logger *zap.Logger, message string, fields ...zap.Field) {
	if logger != nil {
		contextLogger := FromContext(ctx, logger)
		contextLogger.Info(message, fields...)
	}
}

// WithTraceID adds a trace ID to the context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithTenantID adds a tenant ID to the context
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

// WithUserID adds a user ID to the context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// WithOperation adds an operation to the context
func WithOperation(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, OperationKey, operation)
}