package tracing

import (
	"context"
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"

	appconfig "github.com/direito-lux/auth-service/internal/infrastructure/config"
	"github.com/direito-lux/auth-service/internal/infrastructure/logging"
)

// Tracer wrapper para o tracer do Jaeger
type Tracer struct {
	tracer opentracing.Tracer
	closer io.Closer
	logger *zap.Logger
}

// NewTracer cria uma nova instância do tracer
func NewTracer(cfg *appconfig.Config, logger *zap.Logger) (*Tracer, error) {
	if !cfg.Jaeger.Enabled {
		logger.Info("Tracing desabilitado")
		return &Tracer{
			tracer: opentracing.NoopTracer{},
			closer: io.NopCloser(nil),
			logger: logger,
		}, nil
	}

	// Configurar Jaeger
	jaegerCfg := config.Configuration{
		ServiceName: cfg.Jaeger.ServiceName,
		Sampler: &config.SamplerConfig{
			Type:  cfg.Jaeger.SamplerType,
			Param: cfg.Jaeger.SamplerParam,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           cfg.IsDevelopment(),
			CollectorEndpoint:  cfg.Jaeger.Endpoint,
			LocalAgentHostPort: "", // Usar collector endpoint
		},
	}

	// Criar tracer
	tracer, closer, err := jaegerCfg.NewTracer(
		config.Logger(jaeger.StdLogger),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar tracer: %w", err)
	}

	// Definir como tracer global
	opentracing.SetGlobalTracer(tracer)

	logger.Info("Tracer configurado",
		zap.String("service", cfg.Jaeger.ServiceName),
		zap.String("endpoint", cfg.Jaeger.Endpoint),
		zap.String("sampler_type", cfg.Jaeger.SamplerType),
		zap.Float64("sampler_param", cfg.Jaeger.SamplerParam),
	)

	return &Tracer{
		tracer: tracer,
		closer: closer,
		logger: logger,
	}, nil
}

// Close fecha o tracer
func (t *Tracer) Close() error {
	if t.closer != nil {
		return t.closer.Close()
	}
	return nil
}

// StartSpan inicia um novo span
func (t *Tracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return t.tracer.StartSpan(operationName, opts...)
}

// StartSpanFromContext inicia um span a partir do contexto
func StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, operationName, opts...)
}

// SpanFromContext extrai o span do contexto
func SpanFromContext(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}

// ContextWithSpan adiciona o span ao contexto
func ContextWithSpan(ctx context.Context, span opentracing.Span) context.Context {
	return opentracing.ContextWithSpan(ctx, span)
}

// TraceHTTPRequest cria span para requisição HTTP
func TraceHTTPRequest(ctx context.Context, method, url string) (opentracing.Span, context.Context) {
	span, ctx := StartSpanFromContext(ctx, fmt.Sprintf("HTTP %s", method))
	
	ext.HTTPMethod.Set(span, method)
	ext.HTTPUrl.Set(span, url)
	ext.Component.Set(span, "http-client")
	
	return span, ctx
}

// TraceHTTPHandler cria span para handler HTTP
func TraceHTTPHandler(ctx context.Context, method, path string) (opentracing.Span, context.Context) {
	span, ctx := StartSpanFromContext(ctx, fmt.Sprintf("%s %s", method, path))
	
	ext.HTTPMethod.Set(span, method)
	ext.HTTPUrl.Set(span, path)
	ext.Component.Set(span, "http-server")
	ext.SpanKind.Set(span, ext.SpanKindRPCServerEnum)
	
	return span, ctx
}

// TraceDatabase cria span para operação de banco
func TraceDatabase(ctx context.Context, operation, table string) (opentracing.Span, context.Context) {
	span, ctx := StartSpanFromContext(ctx, fmt.Sprintf("DB %s %s", operation, table))
	
	ext.DBType.Set(span, "postgresql")
	ext.DBStatement.Set(span, operation)
	ext.DBInstance.Set(span, table)
	ext.Component.Set(span, "database")
	
	return span, ctx
}

// TraceCache cria span para operação de cache
func TraceCache(ctx context.Context, operation, key string) (opentracing.Span, context.Context) {
	span, ctx := StartSpanFromContext(ctx, fmt.Sprintf("CACHE %s", operation))
	
	span.SetTag("cache.operation", operation)
	span.SetTag("cache.key", key)
	ext.Component.Set(span, "redis")
	
	return span, ctx
}

// TraceMessage cria span para operação de mensagem
func TraceMessage(ctx context.Context, operation, exchange, routingKey string) (opentracing.Span, context.Context) {
	span, ctx := StartSpanFromContext(ctx, fmt.Sprintf("MSG %s", operation))
	
	span.SetTag("messaging.operation", operation)
	span.SetTag("messaging.exchange", exchange)
	span.SetTag("messaging.routing_key", routingKey)
	ext.Component.Set(span, "rabbitmq")
	
	return span, ctx
}

// TraceExternalService cria span para chamada de serviço externo
func TraceExternalService(ctx context.Context, service, operation string) (opentracing.Span, context.Context) {
	span, ctx := StartSpanFromContext(ctx, fmt.Sprintf("%s %s", service, operation))
	
	span.SetTag("service.name", service)
	span.SetTag("service.operation", operation)
	ext.Component.Set(span, "external-service")
	ext.SpanKind.Set(span, ext.SpanKindRPCClientEnum)
	
	return span, ctx
}

// SetSpanError marca o span como erro
func SetSpanError(span opentracing.Span, err error) {
	if span != nil && err != nil {
		ext.Error.Set(span, true)
		span.SetTag("error.message", err.Error())
		span.LogKV("error", err.Error())
	}
}

// SetSpanTenant adiciona informações do tenant ao span
func SetSpanTenant(span opentracing.Span, tenantID string) {
	if span != nil && tenantID != "" {
		span.SetTag("tenant.id", tenantID)
	}
}

// SetSpanUser adiciona informações do usuário ao span
func SetSpanUser(span opentracing.Span, userID string) {
	if span != nil && userID != "" {
		span.SetTag("user.id", userID)
	}
}

// FinishSpan finaliza o span com informações de erro se necessário
func FinishSpan(span opentracing.Span, err error) {
	if span != nil {
		if err != nil {
			SetSpanError(span, err)
		}
		span.Finish()
	}
}

// GetTraceID obtém o trace ID do span no contexto
func GetTraceID(ctx context.Context) string {
	span := SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	// Tentar extrair trace ID do Jaeger
	if jaegerSpan, ok := span.(*jaeger.Span); ok {
		return jaegerSpan.SpanContext().TraceID().String()
	}

	return ""
}

// InjectHTTPHeaders injeta headers de tracing na requisição HTTP
func InjectHTTPHeaders(ctx context.Context, headers map[string]string) error {
	span := SpanFromContext(ctx)
	if span == nil {
		return nil
	}

	carrier := opentracing.TextMapCarrier(headers)
	return opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.TextMap,
		carrier,
	)
}

// ExtractHTTPHeaders extrai contexto de tracing dos headers HTTP
func ExtractHTTPHeaders(headers map[string]string) (opentracing.SpanContext, error) {
	carrier := opentracing.TextMapCarrier(headers)
	return opentracing.GlobalTracer().Extract(
		opentracing.TextMap,
		carrier,
	)
}

// TracedOperation executa uma operação com tracing automático
func TracedOperation(ctx context.Context, operationName string, fn func(context.Context) error) error {
	span, ctx := StartSpanFromContext(ctx, operationName)
	defer span.Finish()

	// Adicionar informações do contexto ao span
	if tenantID := logging.GetTenantID(ctx); tenantID != "" {
		SetSpanTenant(span, tenantID)
	}
	
	if userID := logging.GetUserID(ctx); userID != "" {
		SetSpanUser(span, userID)
	}

	// Executar operação
	err := fn(ctx)
	if err != nil {
		SetSpanError(span, err)
	}

	return err
}

// TracedOperationWithResult executa uma operação com tracing e retorna resultado
func TracedOperationWithResult[T any](ctx context.Context, operationName string, fn func(context.Context) (T, error)) (T, error) {
	span, ctx := StartSpanFromContext(ctx, operationName)
	defer span.Finish()

	// Adicionar informações do contexto ao span
	if tenantID := logging.GetTenantID(ctx); tenantID != "" {
		SetSpanTenant(span, tenantID)
	}
	
	if userID := logging.GetUserID(ctx); userID != "" {
		SetSpanUser(span, userID)
	}

	// Executar operação
	result, err := fn(ctx)
	if err != nil {
		SetSpanError(span, err)
	}

	return result, err
}