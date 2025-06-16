package application

import (
	"context"
)

// EventBus define a interface para publicação de eventos
type EventBus interface {
	Publish(ctx context.Context, event DomainEvent) error
}

// DomainEvent define a interface para eventos de domínio
type DomainEvent interface {
	EventType() string
	AggregateID() string
	Payload() ([]byte, error)
}