package events

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"github.com/streadway/amqp"
	"github.com/direito-lux/process-service/internal/domain"
)

// RabbitMQEventPublisher implementação RabbitMQ do EventPublisher
type RabbitMQEventPublisher struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	exchange   string
}

// NewRabbitMQEventPublisher cria nova instância do publisher
func NewRabbitMQEventPublisher(connection *amqp.Connection, exchange string) (*RabbitMQEventPublisher, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("erro ao criar canal: %w", err)
	}

	publisher := &RabbitMQEventPublisher{
		connection: connection,
		channel:    channel,
		exchange:   exchange,
	}

	// Configurar exchange e filas
	if err := publisher.setupExchangeAndQueues(); err != nil {
		return nil, fmt.Errorf("erro ao configurar exchange e filas: %w", err)
	}

	return publisher, nil
}

// Publish publica um evento
func (p *RabbitMQEventPublisher) Publish(event domain.DomainEvent) error {
	// Serializar evento
	eventData, err := p.serializeEvent(event)
	if err != nil {
		return fmt.Errorf("erro ao serializar evento: %w", err)
	}

	// Definir routing key baseado no tipo do evento
	routingKey := p.getRoutingKey(event.EventType())

	// Publicar evento
	err = p.channel.Publish(
		p.exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    event.Timestamp(),
			MessageId:    fmt.Sprintf("%s_%d", event.AggregateID(), time.Now().UnixNano()),
			Type:         event.EventType(),
			AppId:        "process-service",
			Headers: amqp.Table{
				"event_type":     event.EventType(),
				"event_version":  event.EventVersion(),
				"aggregate_id":   event.AggregateID(),
				"occurred_at":    event.Timestamp().Format(time.RFC3339),
			},
			Body: eventData,
		},
	)

	if err != nil {
		return fmt.Errorf("erro ao publicar evento %s: %w", event.EventType(), err)
	}

	log.Printf("Evento publicado: %s [%s]", event.EventType(), event.AggregateID())
	return nil
}

// PublishBatch publica múltiplos eventos em lote
func (p *RabbitMQEventPublisher) PublishBatch(events []domain.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	// Confirmar publicação para garantir atomicidade
	if err := p.channel.Confirm(false); err != nil {
		return fmt.Errorf("erro ao habilitar confirmação: %w", err)
	}

	confirms := p.channel.NotifyPublish(make(chan amqp.Confirmation, len(events)))

	// Publicar todos os eventos
	for _, event := range events {
		if err := p.Publish(event); err != nil {
			return fmt.Errorf("erro ao publicar evento em lote: %w", err)
		}
	}

	// Aguardar confirmações
	for i := 0; i < len(events); i++ {
		select {
		case confirm := <-confirms:
			if !confirm.Ack {
				return fmt.Errorf("evento %d não foi confirmado", confirm.DeliveryTag)
			}
		case <-time.After(30 * time.Second):
			return fmt.Errorf("timeout aguardando confirmação dos eventos")
		}
	}

	log.Printf("Lote de %d eventos publicado com sucesso", len(events))
	return nil
}

// Close fecha as conexões
func (p *RabbitMQEventPublisher) Close() error {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			return fmt.Errorf("erro ao fechar canal: %w", err)
		}
	}
	return nil
}

// setupExchangeAndQueues configura exchange e filas necessárias
func (p *RabbitMQEventPublisher) setupExchangeAndQueues() error {
	// Declarar exchange principal
	err := p.channel.ExchangeDeclare(
		p.exchange,
		"topic",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("erro ao declarar exchange: %w", err)
	}

	// Configurar filas para diferentes tipos de eventos
	queues := []QueueConfig{
		// Eventos de processo
		{
			Name:       "process.events",
			RoutingKey: "process.*",
			Durable:    true,
		},
		// Eventos de movimentação
		{
			Name:       "movement.events",
			RoutingKey: "movement.*",
			Durable:    true,
		},
		// Eventos de partes
		{
			Name:       "party.events",
			RoutingKey: "party.*",
			Durable:    true,
		},
		// Eventos de monitoramento
		{
			Name:       "monitoring.events",
			RoutingKey: "process.monitoring.*",
			Durable:    true,
		},
		// Eventos de sincronização
		{
			Name:       "sync.events",
			RoutingKey: "process.sync.*",
			Durable:    true,
		},
		// Movimentações importantes (para notificações)
		{
			Name:       "important.movements",
			RoutingKey: "movement.important.*",
			Durable:    true,
		},
		// Eventos de análise
		{
			Name:       "analysis.events",
			RoutingKey: "movement.analyzed",
			Durable:    true,
		},
	}

	// Criar filas e bindings
	for _, queueConfig := range queues {
		if err := p.createQueueAndBinding(queueConfig); err != nil {
			return fmt.Errorf("erro ao criar fila %s: %w", queueConfig.Name, err)
		}
	}

	return nil
}

// createQueueAndBinding cria fila e binding
func (p *RabbitMQEventPublisher) createQueueAndBinding(config QueueConfig) error {
	// Declarar fila
	_, err := p.channel.QueueDeclare(
		config.Name,
		config.Durable,
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("erro ao declarar fila: %w", err)
	}

	// Criar binding
	err = p.channel.QueueBind(
		config.Name,
		config.RoutingKey,
		p.exchange,
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("erro ao criar binding: %w", err)
	}

	log.Printf("Fila criada: %s -> %s", config.Name, config.RoutingKey)
	return nil
}

// getRoutingKey determina routing key baseado no tipo do evento
func (p *RabbitMQEventPublisher) getRoutingKey(eventType string) string {
	// Mapear tipos de evento para routing keys
	routingKeyMap := map[string]string{
		// Eventos de processo
		"process.created":              "process.created",
		"process.updated":              "process.updated",
		"process.archived":             "process.archived",
		"process.reactivated":          "process.reactivated",
		"process.monitoring.enabled":   "process.monitoring.enabled",
		"process.monitoring.disabled":  "process.monitoring.disabled",
		"process.synced":               "process.sync.completed",
		"process.batch.sync.started":   "process.sync.batch.started",
		"process.batch.sync.completed": "process.sync.batch.completed",

		// Eventos de movimentação
		"movement.created":           "movement.created",
		"movement.analyzed":          "movement.analyzed",
		"movement.important.detected": "movement.important.detected",

		// Eventos de partes
		"party.added":   "party.added",
		"party.updated": "party.updated",
	}

	if routingKey, exists := routingKeyMap[eventType]; exists {
		return routingKey
	}

	// Fallback: usar tipo do evento como routing key
	return eventType
}

// serializeEvent serializa evento para JSON
func (p *RabbitMQEventPublisher) serializeEvent(event domain.DomainEvent) ([]byte, error) {
	// Criar envelope do evento com metadata adicional
	envelope := EventEnvelope{
		EventType:    event.EventType(),
		EventVersion: event.EventVersion(),
		AggregateID:  event.AggregateID(),
		OccurredAt:   event.Timestamp(),
		Payload:      event,
	}

	return json.Marshal(envelope)
}

// QueueConfig configuração de fila
type QueueConfig struct {
	Name       string
	RoutingKey string
	Durable    bool
}

// EventEnvelope envelope para eventos
type EventEnvelope struct {
	EventType    string                 `json:"event_type"`
	EventVersion string                 `json:"event_version"`
	AggregateID  string                 `json:"aggregate_id"`
	OccurredAt   time.Time              `json:"occurred_at"`
	Payload      domain.DomainEvent     `json:"payload"`
}

// EventMetrics métricas de eventos
type EventMetrics struct {
	TotalPublished   int64
	TotalFailed      int64
	PublishTime      time.Duration
	LastPublishedAt  time.Time
}

// MetricsCollector coleta métricas de eventos
type MetricsCollector struct {
	metrics EventMetrics
}

// NewMetricsCollector cria novo coletor de métricas
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{}
}

// RecordPublished registra evento publicado
func (m *MetricsCollector) RecordPublished(duration time.Duration) {
	m.metrics.TotalPublished++
	m.metrics.PublishTime = duration
	m.metrics.LastPublishedAt = time.Now()
}

// RecordFailed registra falha na publicação
func (m *MetricsCollector) RecordFailed() {
	m.metrics.TotalFailed++
}

// GetMetrics retorna métricas atuais
func (m *MetricsCollector) GetMetrics() EventMetrics {
	return m.metrics
}

// InstrumentedEventPublisher publisher com instrumentação
type InstrumentedEventPublisher struct {
	publisher domain.EventPublisher
	metrics   *MetricsCollector
}

// NewInstrumentedEventPublisher cria publisher instrumentado
func NewInstrumentedEventPublisher(publisher domain.EventPublisher) *InstrumentedEventPublisher {
	return &InstrumentedEventPublisher{
		publisher: publisher,
		metrics:   NewMetricsCollector(),
	}
}

// Publish publica evento com instrumentação
func (p *InstrumentedEventPublisher) Publish(event domain.DomainEvent) error {
	startTime := time.Now()
	
	err := p.publisher.Publish(event)
	
	duration := time.Since(startTime)
	
	if err != nil {
		p.metrics.RecordFailed()
		return err
	}
	
	p.metrics.RecordPublished(duration)
	return nil
}

// PublishBatch publica lote com instrumentação
func (p *InstrumentedEventPublisher) PublishBatch(events []domain.DomainEvent) error {
	startTime := time.Now()
	
	err := p.publisher.PublishBatch(events)
	
	duration := time.Since(startTime)
	
	if err != nil {
		p.metrics.RecordFailed()
		return err
	}
	
	for range events {
		p.metrics.RecordPublished(duration / time.Duration(len(events)))
	}
	
	return nil
}

// GetMetrics retorna métricas
func (p *InstrumentedEventPublisher) GetMetrics() EventMetrics {
	return p.metrics.GetMetrics()
}

// AsyncEventPublisher publisher assíncrono
type AsyncEventPublisher struct {
	publisher  domain.EventPublisher
	eventQueue chan domain.DomainEvent
	batchSize  int
	flushInterval time.Duration
	done       chan struct{}
}

// NewAsyncEventPublisher cria publisher assíncrono
func NewAsyncEventPublisher(publisher domain.EventPublisher, batchSize int, flushInterval time.Duration) *AsyncEventPublisher {
	async := &AsyncEventPublisher{
		publisher:     publisher,
		eventQueue:    make(chan domain.DomainEvent, batchSize*2),
		batchSize:     batchSize,
		flushInterval: flushInterval,
		done:          make(chan struct{}),
	}

	go async.processEvents()
	return async
}

// Publish adiciona evento à fila assíncrona
func (p *AsyncEventPublisher) Publish(event domain.DomainEvent) error {
	select {
	case p.eventQueue <- event:
		return nil
	default:
		// Fila cheia, publicar diretamente
		return p.publisher.Publish(event)
	}
}

// PublishBatch não implementado para versão assíncrona
func (p *AsyncEventPublisher) PublishBatch(events []domain.DomainEvent) error {
	for _, event := range events {
		if err := p.Publish(event); err != nil {
			return err
		}
	}
	return nil
}

// processEvents processa eventos em lotes
func (p *AsyncEventPublisher) processEvents() {
	batch := make([]domain.DomainEvent, 0, p.batchSize)
	ticker := time.NewTicker(p.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case event := <-p.eventQueue:
			batch = append(batch, event)
			
			if len(batch) >= p.batchSize {
				p.flushBatch(batch)
				batch = batch[:0] // reset slice
			}

		case <-ticker.C:
			if len(batch) > 0 {
				p.flushBatch(batch)
				batch = batch[:0] // reset slice
			}

		case <-p.done:
			// Processar eventos restantes
			if len(batch) > 0 {
				p.flushBatch(batch)
			}
			return
		}
	}
}

// flushBatch publica lote de eventos
func (p *AsyncEventPublisher) flushBatch(batch []domain.DomainEvent) {
	if err := p.publisher.PublishBatch(batch); err != nil {
		log.Printf("Erro ao publicar lote de eventos: %v", err)
		
		// Tentar publicar individualmente como fallback
		for _, event := range batch {
			if err := p.publisher.Publish(event); err != nil {
				log.Printf("Erro ao publicar evento individual %s: %v", event.EventType(), err)
			}
		}
	}
}

// Close para o publisher assíncrono
func (p *AsyncEventPublisher) Close() error {
	close(p.done)
	return nil
}