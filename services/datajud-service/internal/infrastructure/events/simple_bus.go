package events

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/direito-lux/datajud-service/internal/infrastructure/config"
)

// SimpleEventBus implementação simples do event bus
type SimpleEventBus struct {
	config   *config.Config
	logger   *zap.Logger
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewEventBus cria novo event bus simples
func NewEventBus(cfg *config.Config, logger *zap.Logger) (EventBus, error) {
	return &SimpleEventBus{
		config:   cfg,
		logger:   logger,
		handlers: make(map[string][]EventHandler),
	}, nil
}

// Publish publica evento
func (bus *SimpleEventBus) Publish(ctx context.Context, event DomainEvent) error {
	bus.logger.Info("Publicando evento",
		zap.String("event_type", event.GetEventType()),
		zap.String("aggregate_id", event.GetAggregateID()),
	)
	
	return nil
}

// Subscribe registra handler para evento
func (bus *SimpleEventBus) Subscribe(eventType string, handler EventHandler) error {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	if bus.handlers[eventType] == nil {
		bus.handlers[eventType] = make([]EventHandler, 0)
	}
	
	bus.handlers[eventType] = append(bus.handlers[eventType], handler)
	
	bus.logger.Info("Handler registrado",
		zap.String("event_type", eventType),
	)
	
	return nil
}

// Start inicia o event bus (no-op para implementação simples)
func (bus *SimpleEventBus) Start(ctx context.Context) error {
	bus.logger.Info("Event bus simples iniciado")
	return nil
}

// Stop para o event bus (no-op para implementação simples)
func (bus *SimpleEventBus) Stop(ctx context.Context) error {
	bus.logger.Info("Event bus simples parado")
	return nil
}

// StartConsuming inicia consumo (no-op para implementação simples)
func (bus *SimpleEventBus) StartConsuming(ctx context.Context) error {
	bus.logger.Info("Event bus simples consumo iniciado")
	return nil
}

// Close fecha o event bus
func (bus *SimpleEventBus) Close(ctx context.Context) error {
	bus.logger.Info("Event bus simples fechado")
	return nil
}