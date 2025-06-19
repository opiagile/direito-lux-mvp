package events

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"go.uber.org/zap"

	"github.com/direito-lux/report-service/internal/domain"
)

// InMemoryEventBus implementação em memória do event bus
type InMemoryEventBus struct {
	handlers map[string][]domain.EventHandler
	mu       sync.RWMutex
	logger   *zap.Logger
}

// NewInMemoryEventBus cria nova instância do event bus
func NewInMemoryEventBus(logger *zap.Logger) *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]domain.EventHandler),
		logger:   logger,
	}
}

// Publish implementa domain.EventBus
func (b *InMemoryEventBus) Publish(ctx context.Context, event interface{}) error {
	eventType := b.getEventType(event)
	
	b.logger.Debug("Publishing event", 
		zap.String("event_type", eventType),
		zap.Any("event", event))

	b.mu.RLock()
	handlers, exists := b.handlers[eventType]
	b.mu.RUnlock()

	if !exists {
		b.logger.Debug("No handlers registered for event type", zap.String("event_type", eventType))
		return nil
	}

	// Executar handlers em paralelo
	var wg sync.WaitGroup
	errCh := make(chan error, len(handlers))

	for _, handler := range handlers {
		wg.Add(1)
		go func(h domain.EventHandler) {
			defer wg.Done()
			if err := h(ctx, event); err != nil {
				b.logger.Error("Event handler failed",
					zap.String("event_type", eventType),
					zap.Error(err))
				errCh <- err
			}
		}(handler)
	}

	wg.Wait()
	close(errCh)

	// Verificar se houve erros
	var errors []error
	for err := range errCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("event handlers failed: %v", errors)
	}

	return nil
}

// Subscribe implementa domain.EventBus
func (b *InMemoryEventBus) Subscribe(ctx context.Context, eventType string, handler domain.EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)
	
	b.logger.Info("Event handler subscribed", 
		zap.String("event_type", eventType),
		zap.Int("total_handlers", len(b.handlers[eventType])))

	return nil
}

// getEventType extrai o tipo do evento
func (b *InMemoryEventBus) getEventType(event interface{}) string {
	// Tentar extrair EventType do evento
	v := reflect.ValueOf(event)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Procurar campo EventType
	if v.Kind() == reflect.Struct {
		eventTypeField := v.FieldByName("EventType")
		if eventTypeField.IsValid() && eventTypeField.Kind() == reflect.String {
			return eventTypeField.String()
		}
	}

	// Fallback para o nome do tipo
	return reflect.TypeOf(event).String()
}