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
	// TraceIDKey chave para trace ID no contexto
	TraceIDKey ContextKey = "trace_id"
	// TenantIDKey chave para tenant ID no contexto
	TenantIDKey ContextKey = "tenant_id"
	// UserIDKey chave para user ID no contexto
	UserIDKey ContextKey = "user_id"
	// OperationKey chave para operação no contexto
	OperationKey ContextKey = "operation"
)

// Logger wrapper para zap.Logger com funcionalidades extras
type Logger struct {
	*zap.Logger
}

// NewLogger cria uma nova instância do logger
func NewLogger(level, environment string) (*zap.Logger, error) {
	var config zap.Config

	// Configurar baseado no environment
	switch environment {
	case "production":
		config = zap.NewProductionConfig()
		config.DisableStacktrace = true
	case "development", "test":
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
	config.Level.SetLevel(zapLevel)

	// Configurar campos iniciais
	config.InitialFields = map[string]interface{}{
		"service": "mcp-service",
	}

	// Configurar encoding para produção
	if environment == "production" {
		config.Encoding = "json"
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.CallerKey = "caller"
		config.EncoderConfig.MessageKey = "message"
		config.EncoderConfig.LevelKey = "level"
	}

	logger, err := config.Build(
		zap.AddCallerSkip(1), // Skip wrapper function
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar logger: %w", err)
	}

	return logger, nil
}

// FromContext extrai informações do contexto e adiciona ao logger
func FromContext(ctx context.Context, logger *zap.Logger) *zap.Logger {
	fields := []zap.Field{}

	// Extrair trace ID
	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		if id, ok := traceID.(string); ok && id != "" {
			fields = append(fields, zap.String("trace_id", id))
		}
	}

	// Extrair tenant ID
	if tenantID := ctx.Value(TenantIDKey); tenantID != nil {
		if id, ok := tenantID.(string); ok && id != "" {
			fields = append(fields, zap.String("tenant_id", id))
		}
	}

	// Extrair user ID
	if userID := ctx.Value(UserIDKey); userID != nil {
		if id, ok := userID.(string); ok && id != "" {
			fields = append(fields, zap.String("user_id", id))
		}
	}

	// Extrair operação
	if operation := ctx.Value(OperationKey); operation != nil {
		if op, ok := operation.(string); ok && op != "" {
			fields = append(fields, zap.String("operation", op))
		}
	}

	if len(fields) > 0 {
		return logger.With(fields...)
	}

	return logger
}

// GetTraceID extrai o trace ID do contexto
func GetTraceID(ctx context.Context) string {
	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	return ""
}

// GetTenantID extrai o tenant ID do contexto
func GetTenantID(ctx context.Context) string {
	if tenantID := ctx.Value(TenantIDKey); tenantID != nil {
		if id, ok := tenantID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserID extrai o user ID do contexto
func GetUserID(ctx context.Context) string {
	if userID := ctx.Value(UserIDKey); userID != nil {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// LogAudit estrutura para logs de auditoria MCP
type AuditLog struct {
	TenantID    string                 `json:"tenant_id"`
	UserID      string                 `json:"user_id"`
	SessionID   string                 `json:"session_id"`
	Channel     string                 `json:"channel"`
	ToolName    string                 `json:"tool_name"`
	Operation   string                 `json:"operation"`
	Success     bool                   `json:"success"`
	Duration    time.Duration          `json:"duration"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// LogMCPAudit registra log de auditoria MCP
func LogMCPAudit(ctx context.Context, logger *zap.Logger, audit AuditLog) {
	contextLogger := FromContext(ctx, logger)
	
	fields := []zap.Field{
		zap.String("audit_tenant_id", audit.TenantID),
		zap.String("audit_user_id", audit.UserID),
		zap.String("audit_session_id", audit.SessionID),
		zap.String("audit_channel", audit.Channel),
		zap.String("audit_tool_name", audit.ToolName),
		zap.String("audit_operation", audit.Operation),
		zap.Bool("audit_success", audit.Success),
		zap.Duration("audit_duration", audit.Duration),
	}

	if audit.Data != nil {
		fields = append(fields, zap.Any("audit_data", audit.Data))
	}

	if audit.Error != "" {
		fields = append(fields, zap.String("audit_error", audit.Error))
	}

	contextLogger.Info("MCP_AUDIT", fields...)
}