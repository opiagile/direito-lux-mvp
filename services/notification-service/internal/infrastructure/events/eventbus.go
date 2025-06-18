package events

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/direito-lux/notification-service/internal/infrastructure/config"
)

// EventBus interface para bus de eventos
type EventBus interface {
	Publish(ctx context.Context, event DomainEvent) error
	PublishBatch(ctx context.Context, events []DomainEvent) error
	Subscribe(eventType string, handler EventHandler) error
	Unsubscribe(eventType string, handler EventHandler) error
	Close() error
}

// EventHandler interface para handlers de eventos
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
}

// SimpleEventBus implementação simples de event bus
type SimpleEventBus struct {
	publisher Publisher
	logger    *zap.Logger
	config    *config.Config
}

// NewEventBus cria novo event bus
func NewEventBus(cfg *config.Config, logger *zap.Logger) (EventBus, error) {
	// Criar publisher RabbitMQ
	publisher, err := NewRabbitMQPublisher(&cfg.RabbitMQ, logger)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar publisher: %w", err)
	}

	return &SimpleEventBus{
		publisher: publisher,
		logger:    logger,
		config:    cfg,
	}, nil
}

// Publish publica um evento
func (b *SimpleEventBus) Publish(ctx context.Context, event DomainEvent) error {
	return b.publisher.Publish(ctx, event)
}

// PublishBatch publica múltiplos eventos
func (b *SimpleEventBus) PublishBatch(ctx context.Context, events []DomainEvent) error {
	return b.publisher.PublishBatch(ctx, events)
}

// Subscribe implementação placeholder
func (b *SimpleEventBus) Subscribe(eventType string, handler EventHandler) error {
	// TODO: Implementar subscrição de eventos
	b.logger.Info("Subscribe chamado", zap.String("eventType", eventType))
	return nil
}

// Unsubscribe implementação placeholder
func (b *SimpleEventBus) Unsubscribe(eventType string, handler EventHandler) error {
	// TODO: Implementar cancelamento de subscrição
	b.logger.Info("Unsubscribe chamado", zap.String("eventType", eventType))
	return nil
}

// Close fecha o event bus
func (b *SimpleEventBus) Close() error {
	if b.publisher != nil {
		return b.publisher.Close()
	}
	return nil
}