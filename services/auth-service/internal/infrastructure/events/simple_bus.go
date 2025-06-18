package events

import (
	"go.uber.org/zap"
)

// NewEventBus cria nova instância do event bus
func NewEventBus(logger *zap.Logger) EventBus {
	return NewInMemoryEventBus(logger)
}