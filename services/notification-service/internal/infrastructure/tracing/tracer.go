package tracing

import (
	"context"
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

// Tracer wrapper para OpenTracing
type Tracer struct {
	tracer opentracing.Tracer
	closer io.Closer
	logger *zap.Logger
}

// NewTracer cria novo tracer Jaeger
func NewTracer(logger *zap.Logger) (*Tracer, error) {
	serviceName := "notification-service"
	jaegerEndpoint := "localhost:6831"
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1000,
			LocalAgentHostPort:  jaegerEndpoint,
		},
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar tracer Jaeger: %w", err)
	}

	opentracing.SetGlobalTracer(tracer)

	return &Tracer{
		tracer: tracer,
		closer: closer,
		logger: logger,
	}, nil
}

// StartSpan inicia novo span
func (t *Tracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return t.tracer.StartSpan(operationName, opts...)
}

// StartSpanFromContext inicia span a partir do contexto
func (t *Tracer) StartSpanFromContext(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, operationName, opts...)
}

// InjectSpanContext injeta contexto do span
func (t *Tracer) InjectSpanContext(span opentracing.Span, format interface{}, carrier interface{}) error {
	return t.tracer.Inject(span.Context(), format, carrier)
}

// ExtractSpanContext extrai contexto do span
func (t *Tracer) ExtractSpanContext(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return t.tracer.Extract(format, carrier)
}

// Close fecha o tracer
func (t *Tracer) Close() error {
	if t.closer != nil {
		return t.closer.Close()
	}
	return nil
}

// WithTracing adiciona tracing ao contexto
func WithTracing(ctx context.Context, span opentracing.Span) context.Context {
	return opentracing.ContextWithSpan(ctx, span)
}

// SpanFromContext obtém span do contexto
func SpanFromContext(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}

// LogError registra erro no span
func LogError(span opentracing.Span, err error) {
	if span != nil && err != nil {
		span.SetTag("error", true)
		span.LogKV("event", "error", "message", err.Error())
	}
}

// LogInfo registra informação no span
func LogInfo(span opentracing.Span, message string) {
	if span != nil {
		span.LogKV("event", "info", "message", message)
	}
}