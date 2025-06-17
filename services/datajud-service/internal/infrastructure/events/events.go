package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/direito-lux/datajud-service/internal/infrastructure/logging"
)

// DomainEvent interface para eventos de domínio
type DomainEvent interface {
	GetEventID() string
	GetEventType() string
	GetAggregateID() string
	GetTenantID() string
	GetOccurredAt() time.Time
	GetVersion() int
	GetMetadata() EventMetadata
}

// BaseEvent estrutura base para eventos
type BaseEvent struct {
	EventID     string        `json:"event_id"`
	EventType   string        `json:"event_type"`
	AggregateID string        `json:"aggregate_id"`
	TenantID    string        `json:"tenant_id"`
	OccurredAt  time.Time     `json:"occurred_at"`
	Version     int           `json:"version"`
	Metadata    EventMetadata `json:"metadata"`
}

// EventMetadata metadados do evento
type EventMetadata struct {
	CausationID   string            `json:"causation_id,omitempty"`
	CorrelationID string            `json:"correlation_id,omitempty"`
	UserID        string            `json:"user_id,omitempty"`
	Source        string            `json:"source"`
	Headers       map[string]string `json:"headers,omitempty"`
}

// GetEventID implementa DomainEvent
func (e BaseEvent) GetEventID() string {
	return e.EventID
}

// GetEventType implementa DomainEvent
func (e BaseEvent) GetEventType() string {
	return e.EventType
}

// GetAggregateID implementa DomainEvent
func (e BaseEvent) GetAggregateID() string {
	return e.AggregateID
}

// GetTenantID implementa DomainEvent
func (e BaseEvent) GetTenantID() string {
	return e.TenantID
}

// GetOccurredAt implementa DomainEvent
func (e BaseEvent) GetOccurredAt() time.Time {
	return e.OccurredAt
}

// GetVersion implementa DomainEvent
func (e BaseEvent) GetVersion() int {
	return e.Version
}

// GetMetadata implementa DomainEvent
func (e BaseEvent) GetMetadata() EventMetadata {
	return e.Metadata
}

// NewBaseEvent cria um novo evento base
func NewBaseEvent(eventType, aggregateID, tenantID string, version int) BaseEvent {
	return BaseEvent{
		EventID:     uuid.New().String(),
		EventType:   eventType,
		AggregateID: aggregateID,
		TenantID:    tenantID,
		OccurredAt:  time.Now().UTC(),
		Version:     version,
		Metadata: EventMetadata{
			Source: "datajud-service",
		},
	}
}

// WithMetadata adiciona metadados ao evento
func (e BaseEvent) WithMetadata(metadata EventMetadata) BaseEvent {
	e.Metadata = metadata
	return e
}

// WithCorrelationID adiciona correlation ID ao evento
func (e BaseEvent) WithCorrelationID(correlationID string) BaseEvent {
	e.Metadata.CorrelationID = correlationID
	return e
}

// WithCausationID adiciona causation ID ao evento
func (e BaseEvent) WithCausationID(causationID string) BaseEvent {
	e.Metadata.CausationID = causationID
	return e
}

// WithUserID adiciona user ID ao evento
func (e BaseEvent) WithUserID(userID string) BaseEvent {
	e.Metadata.UserID = userID
	return e
}

// FromContext cria evento base a partir do contexto
func FromContext(ctx context.Context, eventType, aggregateID, tenantID string, version int) BaseEvent {
	event := NewBaseEvent(eventType, aggregateID, tenantID, version)
	
	// Extrair informações do contexto
	if userID := logging.GetUserID(ctx); userID != "" {
		event = event.WithUserID(userID)
	}
	
	if traceID := logging.GetTraceID(ctx); traceID != "" {
		event.Metadata.CorrelationID = traceID
	}

	return event
}

// EventHandler interface para handlers de eventos
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
	CanHandle(eventType string) bool
}

// EventBus interface para event bus
type EventBus interface {
	Publish(ctx context.Context, event DomainEvent) error
	Subscribe(eventType string, handler EventHandler) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// InMemoryEventBus implementação em memória do event bus (para testes)
type InMemoryEventBus struct {
	handlers map[string][]EventHandler
	logger   *zap.Logger
}

// NewInMemoryEventBus cria novo event bus em memória
func NewInMemoryEventBus(logger *zap.Logger) *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]EventHandler),
		logger:   logger,
	}
}

// Publish publica evento no bus em memória
func (bus *InMemoryEventBus) Publish(ctx context.Context, event DomainEvent) error {
	eventType := event.GetEventType()
	handlers, exists := bus.handlers[eventType]
	
	if !exists {
		bus.logger.Debug("Nenhum handler encontrado para evento",
			zap.String("event_type", eventType),
			zap.String("event_id", event.GetEventID()),
		)
		return nil
	}

	for _, handler := range handlers {
		go func(h EventHandler) {
			if err := h.Handle(ctx, event); err != nil {
				logging.LogError(ctx, bus.logger, "Erro ao processar evento", err,
					zap.String("event_type", eventType),
					zap.String("event_id", event.GetEventID()),
				)
			}
		}(handler)
	}

	return nil
}

// Subscribe subscreve handler para tipo de evento
func (bus *InMemoryEventBus) Subscribe(eventType string, handler EventHandler) error {
	if bus.handlers[eventType] == nil {
		bus.handlers[eventType] = make([]EventHandler, 0)
	}
	
	bus.handlers[eventType] = append(bus.handlers[eventType], handler)
	
	bus.logger.Info("Handler subscrito",
		zap.String("event_type", eventType),
	)
	
	return nil
}

// Start inicia o event bus
func (bus *InMemoryEventBus) Start(ctx context.Context) error {
	bus.logger.Info("Event bus em memória iniciado")
	return nil
}

// Stop para o event bus
func (bus *InMemoryEventBus) Stop(ctx context.Context) error {
	bus.logger.Info("Event bus em memória parado")
	return nil
}

// EventStore interface para armazenamento de eventos
type EventStore interface {
	Save(ctx context.Context, aggregateID string, events []DomainEvent, expectedVersion int) error
	Load(ctx context.Context, aggregateID string) ([]DomainEvent, error)
	LoadFromVersion(ctx context.Context, aggregateID string, version int) ([]DomainEvent, error)
}

// EventProjection interface para projeções de eventos
type EventProjection interface {
	Project(ctx context.Context, event DomainEvent) error
	CanProject(eventType string) bool
}

// EventSerializer interface para serialização de eventos
type EventSerializer interface {
	Serialize(event DomainEvent) ([]byte, error)
	Deserialize(data []byte, eventType string) (DomainEvent, error)
}

// JSONEventSerializer implementação JSON do serializador
type JSONEventSerializer struct{}

// NewJSONEventSerializer cria novo serializador JSON
func NewJSONEventSerializer() *JSONEventSerializer {
	return &JSONEventSerializer{}
}

// Serialize serializa evento para JSON
func (s *JSONEventSerializer) Serialize(event DomainEvent) ([]byte, error) {
	return json.Marshal(event)
}

// Deserialize desserializa evento do JSON
func (s *JSONEventSerializer) Deserialize(data []byte, eventType string) (DomainEvent, error) {
	// Aqui você implementaria um registry de tipos de eventos
	// Por simplicidade, retornando um evento genérico
	var event BaseEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("erro ao desserializar evento: %w", err)
	}
	
	return event, nil
}

// EventRegistry registry para tipos de eventos
type EventRegistry struct {
	eventTypes map[string]func() DomainEvent
}

// NewEventRegistry cria novo registry de eventos
func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		eventTypes: make(map[string]func() DomainEvent),
	}
}

// Register registra tipo de evento
func (r *EventRegistry) Register(eventType string, factory func() DomainEvent) {
	r.eventTypes[eventType] = factory
}

// Create cria novo evento do tipo especificado
func (r *EventRegistry) Create(eventType string) (DomainEvent, error) {
	factory, exists := r.eventTypes[eventType]
	if !exists {
		return nil, fmt.Errorf("tipo de evento não registrado: %s", eventType)
	}
	
	return factory(), nil
}

// EventSourcing utilitários para event sourcing
type EventSourcing struct {
	store      EventStore
	serializer EventSerializer
	logger     *zap.Logger
}

// NewEventSourcing cria nova instância de event sourcing
func NewEventSourcing(store EventStore, serializer EventSerializer, logger *zap.Logger) *EventSourcing {
	return &EventSourcing{
		store:      store,
		serializer: serializer,
		logger:     logger,
	}
}

// SaveEvents salva eventos no store
func (es *EventSourcing) SaveEvents(ctx context.Context, aggregateID string, events []DomainEvent, expectedVersion int) error {
	if len(events) == 0 {
		return nil
	}

	// Log dos eventos sendo salvos
	for _, event := range events {
		logging.LogInfo(ctx, es.logger, "Salvando evento",
			zap.String("event_type", event.GetEventType()),
			zap.String("event_id", event.GetEventID()),
			zap.String("aggregate_id", aggregateID),
		)
	}

	return es.store.Save(ctx, aggregateID, events, expectedVersion)
}

// LoadEvents carrega eventos do store
func (es *EventSourcing) LoadEvents(ctx context.Context, aggregateID string) ([]DomainEvent, error) {
	events, err := es.store.Load(ctx, aggregateID)
	if err != nil {
		return nil, err
	}

	logging.LogInfo(ctx, es.logger, "Eventos carregados",
		zap.String("aggregate_id", aggregateID),
		zap.Int("event_count", len(events)),
	)

	return events, nil
}

// ReplayEvents reconstrói estado a partir dos eventos
func (es *EventSourcing) ReplayEvents(ctx context.Context, events []DomainEvent, applier func(DomainEvent) error) error {
	for _, event := range events {
		if err := applier(event); err != nil {
			return fmt.Errorf("erro ao aplicar evento %s: %w", event.GetEventID(), err)
		}
	}
	
	return nil
}