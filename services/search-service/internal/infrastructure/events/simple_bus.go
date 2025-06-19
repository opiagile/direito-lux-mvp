package events

import (
	"context"

	"go.uber.org/zap"
)

// Event represents a domain event
type Event interface {
	EventType() string
	AggregateID() string
}

// EventBus defines the interface for publishing events
type EventBus interface {
	Publish(ctx context.Context, event Event) error
}

// SimpleEventBus is a simple in-memory event bus
type SimpleEventBus struct {
	logger *zap.Logger
}

// NewEventBus creates a new simple event bus
func NewEventBus(logger *zap.Logger) EventBus {
	return &SimpleEventBus{
		logger: logger,
	}
}

// Publish publishes an event (simple logging implementation)
func (s *SimpleEventBus) Publish(ctx context.Context, event Event) error {
	s.logger.Info("Event published",
		zap.String("event_type", event.EventType()),
		zap.String("aggregate_id", event.AggregateID()),
	)
	return nil
}