package events

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/direito-lux/mcp-service/internal/domain"
	"github.com/direito-lux/mcp-service/internal/infrastructure/logging"
	"github.com/direito-lux/mcp-service/internal/infrastructure/messaging"
	"github.com/direito-lux/mcp-service/internal/infrastructure/metrics"
)

// EventHandler representa um handler para eventos
type EventHandler func(ctx context.Context, event interface{}) error

// EventBus gerencia a publicação e subscrição de eventos
type EventBus struct {
	logger     *zap.Logger
	messaging  *messaging.RabbitMQConnection
	metrics    *metrics.Metrics
	handlers   map[string][]EventHandler
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewEventBus cria uma nova instância do EventBus
func NewEventBus(
	logger *zap.Logger,
	messaging *messaging.RabbitMQConnection,
	metrics *metrics.Metrics,
) *EventBus {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &EventBus{
		logger:    logger,
		messaging: messaging,
		metrics:   metrics,
		handlers:  make(map[string][]EventHandler),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Subscribe subscreve um handler para um tipo de evento
func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	
	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	
	eb.logger.Debug("Handler subscrito",
		zap.String("event_type", eventType),
		zap.Int("total_handlers", len(eb.handlers[eventType])),
	)
}

// Publish publica um evento
func (eb *EventBus) Publish(ctx context.Context, event interface{}) error {
	start := time.Now()
	eventType := getEventType(event)
	
	contextLogger := logging.FromContext(ctx, eb.logger)
	
	// Publicar localmente
	if err := eb.publishLocal(ctx, eventType, event); err != nil {
		eb.recordMetrics(eventType, "local", time.Since(start), err)
		return fmt.Errorf("erro ao publicar evento local: %w", err)
	}
	
	// Publicar via RabbitMQ se disponível
	if eb.messaging != nil {
		if err := eb.publishRemote(ctx, eventType, event); err != nil {
			contextLogger.Warn("Erro ao publicar evento remoto", zap.Error(err))
			// Não falha se RabbitMQ não estiver disponível
		}
	}
	
	eb.recordMetrics(eventType, "published", time.Since(start), nil)
	
	contextLogger.Debug("Evento publicado",
		zap.String("event_type", eventType),
		zap.Duration("duration", time.Since(start)),
	)
	
	return nil
}

// publishLocal publica evento para handlers locais
func (eb *EventBus) publishLocal(ctx context.Context, eventType string, event interface{}) error {
	eb.mu.RLock()
	handlers := eb.handlers[eventType]
	eb.mu.RUnlock()
	
	for _, handler := range handlers {
		go func(h EventHandler) {
			defer func() {
				if r := recover(); r != nil {
					eb.logger.Error("Panic no handler de evento",
						zap.String("event_type", eventType),
						zap.Any("panic", r),
					)
				}
			}()
			
			if err := h(ctx, event); err != nil {
				eb.logger.Error("Erro no handler de evento",
					zap.String("event_type", eventType),
					zap.Error(err),
				)
			}
		}(handler)
	}
	
	return nil
}

// publishRemote publica evento via RabbitMQ
func (eb *EventBus) publishRemote(ctx context.Context, eventType string, event interface{}) error {
	// Serializar evento
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("erro ao serializar evento: %w", err)
	}
	
	// Criar envelope do evento
	envelope := map[string]interface{}{
		"type":      eventType,
		"data":      event,
		"timestamp": time.Now().UTC(),
		"trace_id":  logging.GetTraceID(ctx),
		"tenant_id": logging.GetTenantID(ctx),
	}
	
	envelopeData, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("erro ao serializar envelope: %w", err)
	}
	
	// Publicar no RabbitMQ
	return eb.messaging.Publish(ctx, "mcp.events", eventType, envelopeData)
}

// PublishMCPEvent publica eventos específicos do MCP
func (eb *EventBus) PublishMCPEvent(ctx context.Context, event domain.MCPEvent) error {
	return eb.Publish(ctx, event)
}

// PublishSessionEvent publica eventos de sessão
func (eb *EventBus) PublishSessionEvent(ctx context.Context, event domain.SessionEvent) error {
	return eb.Publish(ctx, event)
}

// PublishToolEvent publica eventos de ferramentas
func (eb *EventBus) PublishToolEvent(ctx context.Context, event domain.ToolEvent) error {
	return eb.Publish(ctx, event)
}

// recordMetrics registra métricas do evento
func (eb *EventBus) recordMetrics(eventType, operation string, duration time.Duration, err error) {
	if eb.metrics != nil {
		eb.metrics.RecordEvent(eventType, operation, duration, err)
	}
}

// getEventType extrai o tipo do evento
func getEventType(event interface{}) string {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

// Close fecha o EventBus
func (eb *EventBus) Close() error {
	eb.cancel()
	return nil
}

// Health verifica a saúde do EventBus
func (eb *EventBus) Health(ctx context.Context) error {
	select {
	case <-eb.ctx.Done():
		return fmt.Errorf("event bus está fechado")
	default:
		return nil
	}
}